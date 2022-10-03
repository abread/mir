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
	threshDsl "github.com/filecoin-project/mir/pkg/threshcrypto/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
	"github.com/filecoin-project/mir/pkg/vcb/vcbdsl"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID // id of the module to send the "Deliver" event to
	Net          t.ModuleID
	ThreshCrypto t.ModuleID
	Mempool      t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         "vcb",
		Consumer:     consumer,
		Net:          "net",
		ThreshCrypto: "threshcrypto",
		Mempool:      "mempool",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
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
	data    []*requestpb.Request
	sigData [][]byte
	batchID t.BatchID

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

type precomputeSigDataCtx struct{}

func SigData(instanceUid []byte, batchID t.BatchID) [][]byte {
	return [][]byte{
		[]byte("github.com/filecoin-project/mir/pkg/vbc"),
		instanceUid,
		[]byte(batchID),
	}
}

func NewReconfigurableModule(mc *ModuleConfig, nodeID t.NodeID, logger logging.Logger) modules.PassiveModule {
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
				)
				return inst, nil
			},
		),
		logger,
	)
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	state := vcbModuleCommonState{
		data:    nil,
		sigData: nil, // cache sigData

		recvdSent:  false,
		recvdFinal: false,
		delivered:  false,
	}

	vcbdsl.UponSendMessageReceived(m, func(from t.NodeID, data []*requestpb.Request) error {
		if from == params.Leader && !state.recvdSent && !state.delivered {
			state.data = data
			state.recvdSent = true

			mpdsl.RequestTransactionIDs(m, mc.Mempool, data, &handleSendCtx{})
		}
		return nil
	})

	mpdsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *handleSendCtx) error {
		mpdsl.RequestBatchID(m, mc.Mempool, txIDs, context)
		return nil
	})

	mpdsl.UponBatchIDResponse(m, func(batchID t.BatchID, context *handleSendCtx) error {
		state.sigData = SigData(params.InstanceUID, batchID)
		state.batchID = batchID
		threshDsl.SignShare(m, mc.ThreshCrypto, state.sigData, context)
		return nil
	})

	threshDsl.UponSignShareResult(m, func(sigShare []byte, context *handleSendCtx) error {
		dsl.SendMessage(m, mc.Net, EchoMessage(mc.Self, sigShare), []t.NodeID{params.Leader})
		return nil
	})

	if nodeID == params.Leader {
		setupVcbLeader(m, mc, params, &state)
	} else {
		vcbdsl.UponFinalMessageReceived(m, func(from t.NodeID, data []*requestpb.Request, signature []byte) error {
			if from == params.Leader && !state.delivered && !state.recvdFinal {
				state.recvdFinal = true
				ctx := &handleFinalCtx{
					signature: signature,
				}

				if state.sigData != nil {
					threshDsl.VerifyFull(m, mc.ThreshCrypto, state.sigData, signature, ctx)
				} else {
					state.data = data
					mpdsl.RequestTransactionIDs(m, mc.Mempool, data, ctx)
				}
			}

			return nil
		})

		mpdsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *handleFinalCtx) error {
			mpdsl.RequestBatchID(m, mc.Mempool, txIDs, context)
			return nil
		})

		mpdsl.UponBatchIDResponse(m, func(batchID t.BatchID, context *handleFinalCtx) error {
			state.sigData = SigData(params.InstanceUID, batchID)
			state.batchID = batchID
			threshDsl.VerifyFull(m, mc.ThreshCrypto, state.sigData, context.signature, context)
			return nil
		})

		threshDsl.UponVerifyFullResult(m, func(ok bool, err string, context *handleFinalCtx) error {
			if ok {
				state.delivered = true
				vcbdsl.Deliver(m, mc.Consumer, state.data, state.batchID, context.signature)
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
		sigShares:    make([][]byte, params.GetN()-params.GetF()),
	}

	vcbdsl.UponBroadcastRequest(m, func(data []*requestpb.Request) error {
		if commonState.data != nil {
			return fmt.Errorf("cannot vcb-broadcast more than once in same instance")
		} else if data == nil {
			return fmt.Errorf("cannot vcb-broadcast nil")
		}

		commonState.data = data

		// pre-compute sigData before broadcasting SEND(m) to simplify further code
		ctx := &precomputeSigDataCtx{}
		mpdsl.RequestTransactionIDs(m, mc.Mempool, data, ctx)
		return nil
	})

	mpdsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *precomputeSigDataCtx) error {
		mpdsl.RequestBatchID(m, mc.Mempool, txIDs, context)
		return nil
	})

	mpdsl.UponBatchIDResponse(m, func(batchID t.BatchID, context *precomputeSigDataCtx) error {
		commonState.sigData = SigData(params.InstanceUID, batchID)
		commonState.batchID = batchID
		dsl.SendMessage(m, mc.Net, SendMessage(mc.Self, commonState.data), params.AllNodes)
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

			dsl.SendMessage(m, mc.Net, FinalMessage(mc.Self, commonState.data, fullSig), params.AllNodes)
			vcbdsl.Deliver(m, mc.Consumer, commonState.data, commonState.batchID, fullSig)
			commonState.delivered = true
		}
		return nil
	})
}
