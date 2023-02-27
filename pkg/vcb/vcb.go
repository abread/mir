package vcb

import (
	"context"
	"fmt"

	"github.com/filecoin-project/mir/pkg/abba/abbatypes"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	rnetdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	threshDsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	threshcryptopbtypes "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/types"
	vcbdsl "github.com/filecoin-project/mir/pkg/pb/vcbpb/dsl"
	vcbmsgs "github.com/filecoin-project/mir/pkg/pb/vcbpb/msgs"
	vcbtypes "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	ReliableNet  t.ModuleID
	ThreshCrypto t.ModuleID
	Mempool      t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self:         "vcb",
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

type state struct {
	phase vcbPhase

	payload *vcbPayloadManager
	sig     []byte
	origin  *vcbtypes.Origin
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

type leaderState struct {
	phase vcbLeaderPhase

	receivedEcho           abbatypes.RecvTracker // TODO: move to common package
	sigShares              []tctypes.SigShare
	lastCombineAttemptSize int
}

type vcbLeaderPhase uint8

const (
	VcbLeaderPhaseAwaitingInput vcbLeaderPhase = iota
	VcbLeaderPhaseAwaitingEchoes
	VcbLeaderPhaseCombining
	VcbLeaderPhaseDone
)

func NewModule(ctx context.Context, mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, logger logging.Logger) modules.PassiveModule {
	m := dsl.NewModule(ctx, mc.Self)

	state := &state{
		payload: newVcbPayloadManager(m, mc, params),

		phase: VcbPhaseAwaitingSend,
	}

	if nodeID == params.Leader {
		setupVcbLeader(m, mc, params, nodeID, logger, state)
	}

	vcbdsl.UponInputValue(m, func(txs []*requestpb.Request, origin *vcbtypes.Origin) error {
		if state.origin != nil {
			return fmt.Errorf("duplicate input value in vcb")
		}
		state.origin = origin
		return nil
	})

	vcbdsl.UponFinalMessageReceived(m, func(from t.NodeID, txs []*requestpb.Request, signature tctypes.FullSig) error {
		if from != params.Leader {
			return nil // byz node // TODO: suspect?
		}

		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, FinalMsgID(), from)

		if state.phase >= VcbPhaseDelivered {
			return nil // already returned
		}

		logger.Log(logging.LevelTrace, "recvd FINAL", "txs", txs, "signature", signature)

		// FINAL acts as ack for ECHO messages
		rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, EchoMsgID(), []t.NodeID{params.Leader})

		state.payload.Input(txs)
		state.sig = signature
		state.phase = VcbPhasePendingVerification
		return nil
	})
	dsl.UponCondition(m, func() error {
		if state.phase == VcbPhasePendingVerification && state.payload.SigData() != nil && state.origin != nil {
			var nilStructPtr *struct{}

			if nodeID == params.Leader {
				// bypass verification, it came from ourselves
				contextID := m.DslHandle().StoreContext(nilStructPtr)
				traceCtx := m.DslHandle().TraceContextAsMap()
				origin := &threshcryptopbtypes.VerifyFullOrigin{
					Module: mc.ThreshCrypto,
					Type:   &threshcryptopbtypes.VerifyFullOrigin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
				}
				threshDsl.VerifyFullResult(m, mc.Self, true, "", origin)
			} else {
				threshDsl.VerifyFull[struct{}](m, mc.ThreshCrypto, state.payload.SigData(), state.sig, nil)
			}

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
	dsl.UponCondition(m, func() error {
		if state.origin == nil {
			return nil
		}

		switch state.phase {
		case VcbPhaseAwaitingNilDelivery:
			logger.Log(logging.LevelWarn, "delivering failure")
			vcbdsl.Deliver(m, state.origin.Module, nil, nil, nil, state.origin)
		case VcbPhaseAwaitingOkDelivery:
			logger.Log(logging.LevelTrace, "delivering batch", "txs", state.payload.Txs(), "sig", state.sig)
			vcbdsl.Deliver(m, state.origin.Module,
				state.payload.Txs(),
				state.payload.TxIDs(),
				state.sig,
				state.origin,
			)
		default:
			return nil
		}

		state.phase = VcbPhaseDelivered
		return nil
	})

	vcbdsl.UponSendMessageReceived(m, func(from t.NodeID, txs []*requestpb.Request) error {
		if from != params.Leader {
			return nil // byz node // TODO: suspect?
		}
		if state.phase > VcbPhaseAwaitingSend {
			return nil // already moved on from this
		}

		logger.Log(logging.LevelTrace, "recvd SEND", "txs", txs)
		state.payload.Input(txs)
		state.phase = VcbPhaseAwaitingSigData
		return nil
	})
	dsl.UponCondition(m, func() error {
		if state.phase == VcbPhaseAwaitingSigData && state.payload.SigData() != nil {
			threshDsl.SignShare[struct{}](m, mc.ThreshCrypto, state.payload.SigData(), nil)
			state.phase = VcbPhaseAwaitingSigShare
		}
		return nil
	})
	threshDsl.UponSignShareResult(m, func(signatureShare tctypes.SigShare, _ctx *struct{}) error {
		if state.phase == VcbPhaseAwaitingSigShare {
			logger.Log(logging.LevelTrace, "sending ECHO", "sigShare", signatureShare)
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

func setupVcbLeader(m dsl.Module, mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, logger logging.Logger, state *state) {
	leaderState := &leaderState{
		receivedEcho: make(abbatypes.RecvTracker, len(params.AllNodes)),
		sigShares:    make([]tctypes.SigShare, 0, len(params.AllNodes)),

		phase: VcbLeaderPhaseAwaitingInput,
	}

	vcbdsl.UponEchoMessageReceived(m, func(from t.NodeID, signatureShare tctypes.SigShare) error {
		if leaderState.phase != VcbLeaderPhaseAwaitingEchoes {
			return nil // we're doing something better
		}

		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, EchoMsgID(), from)
		if !leaderState.receivedEcho.Register(from) {
			return nil // already received
		}
		logger.Log(logging.LevelTrace, "recvd ECHO", "sigShare", signatureShare, "from", from)

		// ECHO message acts as acknowledgement for SEND message
		rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, SendMsgID(), []t.NodeID{from})

		threshDsl.VerifyShare(m, mc.ThreshCrypto, state.payload.SigData(), signatureShare, from, &signatureShare)
		return nil
	})
	threshDsl.UponVerifyShareResult(m, func(ok bool, error string, sigShare *tctypes.SigShare) error {
		// TODO: report byz nodes?
		if ok {
			leaderState.sigShares = append(leaderState.sigShares, *sigShare)
		}
		return nil
	})

	dsl.UponCondition(m, func() error {
		if leaderState.phase == VcbLeaderPhaseAwaitingEchoes && len(leaderState.sigShares) >= 2*params.GetF()+1 && len(leaderState.sigShares) > leaderState.lastCombineAttemptSize {
			logger.Log(logging.LevelTrace, "trying to combine full signature")

			threshDsl.Recover[struct{}](m, mc.ThreshCrypto, state.payload.SigData(), leaderState.sigShares, nil)
			leaderState.phase = VcbLeaderPhaseCombining
			leaderState.lastCombineAttemptSize = len(leaderState.sigShares)
		}
		return nil
	})
	threshDsl.UponRecoverResult(m, func(fullSignature tctypes.FullSig, ok bool, error string, _ctx *struct{}) error {
		if ok {
			logger.Log(logging.LevelTrace, "recovered full sig. sending FINAL")
			rnetdsl.SendMessage(m, mc.ReliableNet, FinalMsgID(), vcbmsgs.FinalMessage(
				mc.Self,
				state.payload.Txs(),
				fullSignature,
			), params.AllNodes)

			// all SEND messages were made redundant, we can forget about them
			rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, SendMsgID(), params.AllNodes)

			leaderState.phase = VcbLeaderPhaseDone
		} else {
			logger.Log(logging.LevelWarn, "full signature combination failed, retrying after receiving more ECHOes")
			leaderState.phase = VcbLeaderPhaseAwaitingEchoes
		}
		return nil
	})

	vcbdsl.UponInputValue(m, func(txs []*requestpb.Request, origin *vcbtypes.Origin) error {
		if nodeID == params.Leader {
			logger.Log(logging.LevelTrace, "inputting value", "txs", txs)
			state.payload.Input(txs)
		}
		return nil
	})
	dsl.UponCondition(m, func() error {
		if leaderState.phase == VcbLeaderPhaseAwaitingInput && state.payload.SigData() != nil {
			// to make things easier, we only broadcast SEND when we are ready to validate the replies (ECHO messages)
			logger.Log(logging.LevelTrace, "sending SEND", "txs", state.payload.Txs())
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
