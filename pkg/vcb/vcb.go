package vcb

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	mpdsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	threshDsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	vcbdsl "github.com/filecoin-project/mir/pkg/pb/vcbpb/dsl"
	"github.com/filecoin-project/mir/pkg/reliablenet/rnetdsl"
	"github.com/filecoin-project/mir/pkg/serializing"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	t "github.com/filecoin-project/mir/pkg/types"
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
	sigShares    []tctypes.SigShare
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

	threshDsl.UponSignShareResult(m, func(sigShare tctypes.SigShare, context *handleSendCtx) error {
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
		vcbdsl.UponFinalMessageReceived(m, func(from t.NodeID, txs []*requestpb.Request, signature tctypes.FullSig) error {
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
				vcbdsl.Deliver(m, mc.Consumer, state.txs, state.txIDs, context.signature, mc.Self)
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
		sigShares:    make([]tctypes.SigShare, 0, params.GetN()),
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

	vcbdsl.UponEchoMessageReceived(m, func(from t.NodeID, sigShare tctypes.SigShare) error {
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

	threshDsl.UponRecoverResult(m, func(fullSig tctypes.FullSig, ok bool, err string, context *recoverVcbSigCtx) error {
		if ok && !state.sentFinal {
			state.sentFinal = true

			rnetdsl.SendMessage(m, mc.ReliableNet,
				FinalMsgID(),
				FinalMessage(mc.Self, commonState.txs, fullSig),
				params.AllNodes,
			)

			// no need to send SEND messages anymore (or any message to ourselves)
			rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, SendMsgID(), params.AllNodes)
			rnetdsl.MarkModuleMsgsRecvd(m, mc.ReliableNet, mc.Self, []t.NodeID{params.Leader})
			// this is running concurrently with the SendMessage above, so the FINISH() message may remain in queue

			vcbdsl.Deliver(m, mc.Consumer, commonState.txs, commonState.txIDs, fullSig, mc.Self)
			commonState.delivered = true
		}
		return nil
	})

	vcbdsl.UponFinalMessageReceived(m, func(from t.NodeID, data []*requestpb.Request, signature tctypes.FullSig) error {
		if from == params.Leader {
			// ack ourselves to free up the queue
			rnetdsl.Ack(m, mc.ReliableNet, mc.Self, FinalMsgID(), from)
		}
		return nil
	})
}

const (
	MsgTypeSend uint8 = iota
	MsgTypeEcho
	MsgTypeFinal
)

func SendMsgID() []byte {
	return []byte{MsgTypeSend}
}

func EchoMsgID() []byte {
	// each node only sends one of these messages, no other parameters are needed
	return []byte{MsgTypeEcho}
}

func FinalMsgID() []byte {
	return []byte{MsgTypeFinal}
}
