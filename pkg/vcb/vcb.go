package vcb

import (
	"fmt"
	"runtime"

	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	rnetdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	threshDsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	vcbdsl "github.com/filecoin-project/mir/pkg/pb/vcbpb/dsl"
	vcbmsgs "github.com/filecoin-project/mir/pkg/pb/vcbpb/msgs"
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tsagg"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID
	ReliableNet  t.ModuleID
	Hasher       t.ModuleID
	ThreshCrypto t.ModuleID
	Mempool      t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         "vcb",
		Consumer:     consumer,
		ReliableNet:  "reliablenet",
		Hasher:       "hasher",
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

type state struct {
	phase vcbPhase

	payload vcbPayloadManager
	sig     []byte
}

type vcbPhase uint8

const (
	VcbPhaseAwaitingSend vcbPhase = iota
	VcbPhaseAwaitingSigData
	VcbPhaseAwaitingSigShare
	VcbPhaseAwaitingFinal
	VcbPhasePendingVerification
	VcbPhaseVerifying
	VcbPhaseAwaitingOkDelivery
	VcbPhaseAwaitingNilDelivery
	VcbPhaseDelivered
)

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, logger logging.Logger) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	state := &state{
		phase: VcbPhaseAwaitingSend,
	}
	state.payload.init(m, mc, params)

	if nodeID == params.Leader {
		setupVcbLeader(m, mc, params, nodeID, logger, state)
	} else {
		vcbdsl.UponInputValue(m, func(txs []*trantorpbtypes.Transaction) error {
			return fmt.Errorf("vcb input provided for non-leader")
		})
	}

	vcbdsl.UponFinalMessageReceived(m, func(from t.NodeID, signature tctypes.FullSig) error {
		if from != params.Leader {
			return nil // byz node // TODO: suspect?
		}

		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, FinalMsgID(), from)
		if state.phase >= VcbPhasePendingVerification {
			return nil // already received final
		}

		// logger.Log(logging.LevelDebug, "recvd FINAL", "signature", signature)

		// FINAL acts as ack for ECHO messages
		rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, EchoMsgID(), []t.NodeID{params.Leader})

		state.sig = signature

		if from == nodeID {
			// optimization: local FINAL needs no verification
			state.phase = VcbPhaseAwaitingOkDelivery
		} else {
			state.phase = VcbPhasePendingVerification
		}
		return nil
	})
	dsl.UponStateUpdates(m, func() error {
		if state.phase == VcbPhasePendingVerification && state.payload.SigData() != nil {
			threshDsl.VerifyFull[struct{}](m, mc.ThreshCrypto, state.payload.SigData(), state.sig, nil)
			state.phase = VcbPhaseVerifying
		}
		return nil
	})
	threshDsl.UponVerifyFullResult(m, func(ok bool, error string, _ctx *struct{}) error {
		if ok {
			state.phase = VcbPhaseAwaitingOkDelivery
		} else {
			state.phase = VcbPhaseAwaitingNilDelivery
		}
		return nil
	})
	dsl.UponStateUpdates(m, func() error {
		switch state.phase {
		case VcbPhaseAwaitingNilDelivery:
			logger.Log(logging.LevelWarn, "delivering failure")
			vcbdsl.Deliver(m, mc.Consumer, nil, nil, nil, mc.Self)
		case VcbPhaseAwaitingOkDelivery:
			// logger.Log(logging.LevelDebug, "delivering batch", "txs", state.payload.Txs(), "sig", state.sig)
			vcbdsl.Deliver(m, mc.Consumer,
				state.payload.Txs(),
				state.payload.TxIDs(),
				state.sig,
				mc.Self,
			)
		default:
			return nil
		}

		state.phase = VcbPhaseDelivered
		return nil
	})

	vcbdsl.UponSendMessageReceived(m, func(from t.NodeID, txs []*trantorpbtypes.Transaction) error {
		if from != params.Leader {
			return nil // byz node // TODO: suspect?
		}
		if state.phase > VcbPhaseAwaitingSend {
			return nil // already moved on from this
		}

		// logger.Log(logging.LevelDebug, "recvd SEND", "txs", txs)
		state.payload.Input(m, mc, txs)
		state.phase = VcbPhaseAwaitingSigData
		return nil
	})
	dsl.UponStateUpdates(m, func() error {
		if state.phase == VcbPhaseAwaitingSigData && state.payload.SigData() != nil {
			threshDsl.SignShare[struct{}](m, mc.ThreshCrypto, state.payload.SigData(), nil)
			state.phase = VcbPhaseAwaitingSigShare
		}
		return nil
	})
	threshDsl.UponSignShareResult(m, func(signatureShare tctypes.SigShare, _ctx *struct{}) error {
		if state.phase == VcbPhaseAwaitingSigShare {
			// logger.Log(logging.LevelDebug, "sending ECHO", "sigShare", signatureShare)
			rnetdsl.SendMessage(m, mc.ReliableNet,
				EchoMsgID(),
				vcbmsgs.EchoMessage(mc.Self, signatureShare),
				[]t.NodeID{params.Leader},
			)
		}
		return nil
	})

	return m
}

type leaderState struct {
	phase vcbLeaderPhase

	sigAgg *tsagg.ThreshSigAggregator
}

type vcbLeaderPhase uint8

const (
	VcbLeaderPhaseAwaitingInput vcbLeaderPhase = iota
	VcbLeaderPhaseAwaitingEchoes
	VcbLeaderPhaseDone
)

func setupVcbLeader(m dsl.Module, mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, logger logging.Logger, state *state) {
	leaderState := &leaderState{
		phase: VcbLeaderPhaseAwaitingInput,

		sigAgg: tsagg.New(m, &tsagg.Params{
			TCModuleID:              mc.ThreshCrypto,
			Threshold:               2*params.GetF() + 1,
			MaxVerifyShareBatchSize: runtime.NumCPU(),
			SigData:                 state.payload.SigData,
			InitialNodeCount:        params.GetN(),
		}, logging.Decorate(logger, "ThresholdSigAggregator: ")),
	}

	vcbdsl.UponEchoMessageReceived(m, func(from t.NodeID, signatureShare tctypes.SigShare) error {
		if leaderState.phase != VcbLeaderPhaseAwaitingEchoes {
			return nil // we're doing something better
		}

		if leaderState.sigAgg.Add(signatureShare, from) {
			// logger.Log(logging.LevelDebug, "recvd ECHO", "sigShare", signatureShare, "from", from)

			// ECHO message acts as acknowledgement for SEND message
			rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, SendMsgID(), []t.NodeID{from})
		}

		return nil
	})

	dsl.UponStateUpdates(m, func() error {
		if leaderState.phase == VcbLeaderPhaseAwaitingEchoes && leaderState.sigAgg.FullSig() != nil {
			// logger.Log(logging.LevelDebug, "recovered full sig. sending FINAL")

			rnetdsl.SendMessage(m, mc.ReliableNet, FinalMsgID(), vcbmsgs.FinalMessage(
				mc.Self,
				leaderState.sigAgg.FullSig(),
			), params.AllNodes)

			// all SEND messages were made redundant, we can forget about them
			rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, SendMsgID(), params.AllNodes)

			leaderState.phase = VcbLeaderPhaseDone
		}
		return nil
	})

	vcbdsl.UponInputValue(m, func(txs []*trantorpbtypes.Transaction) error {
		// logger.Log(logging.LevelDebug, "inputting value", "txs", txs)
		state.payload.Input(m, mc, txs)
		return nil
	})
	dsl.UponStateUpdates(m, func() error {
		if leaderState.phase == VcbLeaderPhaseAwaitingInput && state.payload.Txs() != nil {
			// to make things easier, we only broadcast SEND when we are ready to validate the replies (ECHO messages)
			// logger.Log(logging.LevelDebug, "sending SEND", "txs", state.payload.Txs())
			rnetdsl.SendMessage(m, mc.ReliableNet, SendMsgID(), vcbmsgs.SendMessage(
				mc.Self,
				state.payload.Txs(),
			), params.AllNodes)
			leaderState.phase = VcbLeaderPhaseAwaitingEchoes
		}
		return nil
	})
}

const (
	MsgTypeSend  = "s"
	MsgTypeEcho  = "e"
	MsgTypeFinal = "f"
)

func SendMsgID() rntypes.MsgID {
	return MsgTypeSend
}

func EchoMsgID() rntypes.MsgID {
	// each node only sends one of these messages, no other parameters are needed
	return MsgTypeEcho
}

func FinalMsgID() rntypes.MsgID {
	return MsgTypeFinal
}
