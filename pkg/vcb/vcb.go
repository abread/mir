package vcb

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/factorymodule"
	"github.com/filecoin-project/mir/pkg/logging"
	mpdsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/factorymodulepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	"github.com/filecoin-project/mir/pkg/reliablenet/rnetdsl"
	"github.com/filecoin-project/mir/pkg/serializing"
	threshDsl "github.com/filecoin-project/mir/pkg/threshcrypto/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
	"github.com/filecoin-project/mir/pkg/vcb/vcbdsl"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID // id of the module to send the "Deliver" event to
	ReliableNet  t.ModuleID
	ThreshCrypto t.ModuleID
	Mempool      t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         "vcb",
		Consumer:     consumer,
		ReliableNet:  "reliablenet",
		ThreshCrypto: "threshcrypto",
		Mempool:      "mempool",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters, apart from the Origin.
type ModuleParams struct {
	InstanceUID []byte     // unique identifier for this instance of VCB
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
	Leader      t.NodeID   // the id of the leader of the instance
}

// GetN returns the total number of nodes.
func (params *ModuleParams) GetN() int {
	return len(params.AllNodes)
}

// GetF returns the maximum tolerated number of faulty nodes.
func (params *ModuleParams) GetF() int {
	return (params.GetN() - 1) / 3
}

type vcbModuleLeaderState struct {
	sentFinal bool

	receivedEcho map[t.NodeID]struct{}
	sigShares    [][]byte
}
type vcbModuleCommonState struct {
	txs     []*requestpb.Request
	sigData [][]byte
	txIDs   []t.TxID

	recvdSent  bool
	recvdFinal bool
	delivered  bool
}

type handleSendCtx struct{}

type verifyEchoMsgShareCtx struct {
	sigShare []byte
}
type recoverVcbSigCtx struct{}
type handleFinalCtx struct {
	signature []byte
}

func SigData(instanceUID []byte, txIDs []t.TxID) [][]byte {
	res := make([][]byte, 0, len(txIDs)+3)

	res = append(res, []byte("github.com/filecoin-project/mir/pkg/vcb"))
	res = append(res, instanceUID)
	res = append(res, serializing.Uint64ToBytes(uint64(len(txIDs))))

	for _, txID := range txIDs {
		res = append(res, txID.Pb())
	}

	return res
}

func NewReconfigurableModule(mc *ModuleConfig, nodeID t.NodeID, logger logging.Logger) *factorymodule.FactoryModule {
	return factorymodule.New(
		mc.Self,
		factorymodule.DefaultParams(

			// This function will be called whenever the factory module
			// is asked to create a new instance of the vcb protocol.
			func(vcbID t.ModuleID, genericParams *factorymodulepb.GeneratorParams) (modules.PassiveModule, error) {
				params := genericParams.Type.(*factorymodulepb.GeneratorParams_Vcb).Vcb

				// Extract the IDs of the nodes in the membership associated with this instance
				allNodes := maputil.GetSortedKeys(t.Membership(params.Membership))
				leader := t.NodeID(params.LeaderId)

				if !slices.Contains(allNodes, leader) {
					return nil, fmt.Errorf("leader must be part of the node set")
				}

				// Crate a copy of basic module config with an adapted ID for the submodule.
				submc := *mc
				submc.Self = vcbID

				// Create a new instance of the vcb protocol.
				inst := NewModule(
					&submc,
					&ModuleParams{
						// TODO: Use InstanceUIDs properly.
						//       (E.g., concatenate this with the instantiating protocol's InstanceUID when introduced.)
						InstanceUID: []byte(vcbID),
						AllNodes:    allNodes,
						Leader:      leader,
					},
					nodeID,
					logging.Decorate(logger, "Vcb: ", "leader", params.LeaderId, "id", vcbID),
				)
				return inst, nil
			},
		),
		logger,
	)
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, logger logging.Logger) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	state := vcbModuleCommonState{
		txs:     nil,
		sigData: nil, // cache sigData

		recvdSent:  false,
		recvdFinal: false,
		delivered:  false,
	}

	vcbdsl.UponSendMessageReceived(m, func(from t.NodeID, txs []*requestpb.Request) error {
		if from == params.Leader && !state.recvdSent && !state.delivered {
			state.txs = txs
			state.recvdSent = true

			mpdsl.RequestTransactionIDs(m, mc.Mempool, txs, &handleSendCtx{})
		}
		return nil
	})

	mpdsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *handleSendCtx) error {
		state.txIDs = txIDs
		state.sigData = SigData(params.InstanceUID, txIDs)
		threshDsl.SignShare(m, mc.ThreshCrypto, state.sigData, context)
		return nil
	})

	threshDsl.UponSignShareResult(m, func(sigShare []byte, context *handleSendCtx) error {
		rnetdsl.SendMessage(m, mc.ReliableNet,
			EchoMsgID(),
			EchoMessage(mc.Self, sigShare),
			[]t.NodeID{params.Leader},
		)
		return nil
	})

	if nodeID == params.Leader {
		setupVcbLeader(m, mc, params, &state)
	} else {
		vcbdsl.UponFinalMessageReceived(m, func(from t.NodeID, txs []*requestpb.Request, signature []byte) error {
			if from == params.Leader && !state.delivered && !state.recvdFinal {
				state.recvdFinal = true
				rnetdsl.MarkModuleMsgsRecvd(m, mc.ReliableNet, mc.Self, params.AllNodes)
				rnetdsl.Ack(m, mc.ReliableNet, mc.Self, FinalMsgID(), from)

				ctx := &handleFinalCtx{
					signature: signature,
				}

				if state.sigData != nil {
					threshDsl.VerifyFull(m, mc.ThreshCrypto, state.sigData, signature, ctx)
				} else {
					state.txs = txs
					mpdsl.RequestTransactionIDs(m, mc.Mempool, txs, ctx)
				}
			}

			return nil
		})

		mpdsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *handleFinalCtx) error {
			state.txIDs = txIDs
			state.sigData = SigData(params.InstanceUID, txIDs)
			threshDsl.VerifyFull(m, mc.ThreshCrypto, state.sigData, context.signature, context)
			return nil
		})

		threshDsl.UponVerifyFullResult(m, func(ok bool, err string, context *handleFinalCtx) error {
			if ok {
				state.delivered = true
				vcbdsl.Deliver(m, mc.Consumer, state.txIDs, state.txs, context.signature)
			}

			return nil
		})
	}

	return m
}

func setupVcbLeader(m dsl.Module, mc *ModuleConfig, params *ModuleParams, commonState *vcbModuleCommonState) {
	state := vcbModuleLeaderState{
		sentFinal: false,

		receivedEcho: make(map[t.NodeID]struct{}, len(params.AllNodes)),
		sigShares:    make([][]byte, 0, params.GetN()-params.GetF()),
	}

	vcbdsl.UponBroadcastRequest(m, func(txIDs []t.TxID, txs []*requestpb.Request) error {
		if commonState.txs != nil {
			return fmt.Errorf("cannot vcb-broadcast more than once in same instance")
		} else if len(txs) == 0 {
			return fmt.Errorf("cannot vcb-broadcast an empty batch")
		}

		commonState.txIDs = txIDs
		commonState.txs = txs
		commonState.sigData = SigData(params.InstanceUID, txIDs)

		rnetdsl.SendMessage(m, mc.ReliableNet,
			SendMsgID(),
			SendMessage(mc.Self, commonState.txs),
			params.AllNodes,
		)
		return nil
	})

	vcbdsl.UponEchoMessageReceived(m, func(from t.NodeID, sigShare []byte) error {
		if commonState.delivered {
			return nil
		}

		if _, present := state.receivedEcho[from]; present {
			return nil // already received Echo from this node
		}
		state.receivedEcho[from] = struct{}{}
		rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, SendMsgID(), []t.NodeID{from})

		threshDsl.VerifyShare(m, mc.ThreshCrypto, commonState.sigData, sigShare, from, &verifyEchoMsgShareCtx{
			sigShare: sigShare,
		})

		return nil
	})

	threshDsl.UponVerifyShareResult(m, func(ok bool, err string, context *verifyEchoMsgShareCtx) error {
		if ok {
			state.sigShares = append(state.sigShares, context.sigShare)

			if len(state.sigShares) >= (2*params.GetF()+1) && !state.sentFinal {
				// TODO: avoid calling Recover while another is in progress
				threshDsl.Recover(m, mc.ThreshCrypto, commonState.sigData, state.sigShares, &recoverVcbSigCtx{})
			}
		}

		return nil
	})

	threshDsl.UponRecoverResult(m, func(ok bool, fullSig []byte, err string, context *recoverVcbSigCtx) error {
		if ok && !state.sentFinal {
			state.sentFinal = true

			rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, SendMsgID(), params.AllNodes)
			rnetdsl.SendMessage(m, mc.ReliableNet,
				FinalMsgID(),
				FinalMessage(mc.Self, commonState.txs, fullSig),
				params.AllNodes,
			)
			vcbdsl.Deliver(m, mc.Consumer, commonState.txIDs, commonState.txs, fullSig)
			commonState.delivered = true
		}
		return nil
	})
}

const (
	MSG_TYPE_SEND uint8 = iota
	MSG_TYPE_ECHO
	MSG_TYPE_FINAL
)

func SendMsgID() []byte {
	return []byte{MSG_TYPE_SEND}
}

func EchoMsgID() []byte {
	// each node only sends one of these messages, no other parameters are needed
	return []byte{MSG_TYPE_ECHO}
}

func FinalMsgID() []byte {
	return []byte{MSG_TYPE_FINAL}
}
