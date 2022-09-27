package vcb

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	threshDsl "github.com/filecoin-project/mir/pkg/threshcrypto/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb/vcbdsl"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID // id of the module to send the "Deliver" event to
	Net          t.ModuleID
	ThreshCrypto t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         "vcb",
		Consumer:     consumer,
		Net:          "net",
		ThreshCrypto: "threshcrypto",
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
	request []byte

	sentFinal bool

	receivedEcho map[t.NodeID]struct{}
	sigShares    [][]byte
}
type vcbModuleCommonState struct {
	request []byte

	recvdSent bool
	delivered bool
}

type signSentMsgCtx struct{}
type verifyEchoMsgShareCtx struct {
	sigShare []byte
}
type recoverVcbSigCtx struct{}
type verifyFinalMsgCtx struct {
	signature []byte
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	state := vcbModuleCommonState{
		request: nil,

		recvdSent: false,
		delivered: false,
	}

	vcbdsl.UponSendMessageReceived(m, func(from t.NodeID, data []byte) error {
		if from == params.Leader && !state.recvdSent {
			state.request = data
			state.recvdSent = true
			sigMsg := dataToSignForEcho(params.InstanceUID, data)
			threshDsl.SignShare(m, mc.ThreshCrypto, sigMsg, &signSentMsgCtx{})
		}
		return nil
	})

	threshDsl.UponSignShareResult(m, func(sigShare []byte, context *signSentMsgCtx) error {
		dsl.SendMessage(m, mc.Net, EchoMessage(mc.Self, sigShare), []t.NodeID{params.Leader})
		return nil
	})

	if nodeID == params.Leader {
		setupVcbLeader(m, mc, params)
	} else {
		vcbdsl.UponFinalMessageReceived(m, func(from t.NodeID, data []byte, signature []byte) error {
			if from == params.Leader && !state.delivered {
				sigMsg := dataToSignForEcho(params.InstanceUID, data)
				threshDsl.VerifyFull(m, mc.ThreshCrypto, sigMsg, signature, &verifyFinalMsgCtx{
					signature: signature,
				})
			}

			return nil
		})

		threshDsl.UponVerifyFullResult(m, func(ok bool, err string, context *verifyFinalMsgCtx) error {
			if ok {
				state.delivered = true
				vcbdsl.Deliver(m, mc.Consumer, state.request, context.signature)
			}

			return nil
		})
	}

	return m
}

func setupVcbLeader(m dsl.Module, mc *ModuleConfig, params *ModuleParams) {
	state := vcbModuleLeaderState{
		request: nil,

		sentFinal: false,

		receivedEcho: make(map[t.NodeID]struct{}, len(params.AllNodes)),
		sigShares:    make([][]byte, params.GetN()-params.GetF()),
	}

	vcbdsl.UponBroadcastRequest(m, func(data []byte) error {
		state.request = data

		dsl.SendMessage(m, mc.Net, SendMessage(mc.Self, data), params.AllNodes)
		return nil
	})

	vcbdsl.UponEchoMessageReceived(m, func(from t.NodeID, sigShare []byte) error {
		if _, present := state.receivedEcho[from]; present {
			return nil // already received Echo from this node
		}
		state.receivedEcho[from] = struct{}{}

		threshDsl.VerifyShare(m, mc.ThreshCrypto, dataToSignForEcho(params.InstanceUID, state.request), sigShare, from, &verifyEchoMsgShareCtx{
			sigShare: sigShare,
		})

		return nil
	})

	threshDsl.UponVerifyShareResult(m, func(ok bool, err string, context *verifyEchoMsgShareCtx) error {
		if ok {
			state.sigShares = append(state.sigShares, context.sigShare)

			if len(state.sigShares) >= (2*params.GetF()+1) && !state.sentFinal {
				// TODO: avoid calling Recover while another is in progress
				threshDsl.Recover(m, mc.ThreshCrypto, dataToSignForEcho(params.InstanceUID, state.request), state.sigShares, &recoverVcbSigCtx{})
			}
		}

		return nil
	})

	threshDsl.UponRecoverResult(m, func(ok bool, fullSig []byte, err string, context *recoverVcbSigCtx) error {
		if ok && !state.sentFinal {
			state.sentFinal = true

			dsl.SendMessage(m, mc.Net, FinalMessage(mc.Self, state.request, fullSig), params.AllNodes)
			vcbdsl.Deliver(m, mc.Consumer, state.request, fullSig)
		}
		return nil
	})
}

func dataToSignForEcho(instanceUid []byte, data []byte) [][]byte {
	return [][]byte{
		instanceUid,
		[]byte("ECHO"),
		data,
	}
}
