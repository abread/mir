package dsl

import (
	abaConfig "github.com/filecoin-project/mir/pkg/alea/aba/config"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func BroadcastFinish(m dsl.Module, config *abaConfig.Config, value bool) {
	broadcastMessage(m, config, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_Finish{
			Finish: &aleapb.CobaltFinish{
				Value: value,
			},
		},
	})
}

func BroadcastInit(m dsl.Module, config *abaConfig.Config, roundNumber uint64, estimate bool) {
	broadcastRoundMessage(m, config, &aleapb.CobaltRoundMessage{
		RoundNumber: roundNumber,
		Type: &aleapb.CobaltRoundMessage_Init{
			Init: &aleapb.CobaltInit{
				Estimate: estimate,
			},
		},
	})
}

func BroadcastAux(m dsl.Module, config *abaConfig.Config, roundNumber uint64, value bool) {
	broadcastRoundMessage(m, config, &aleapb.CobaltRoundMessage{
		RoundNumber: roundNumber,
		Type: &aleapb.CobaltRoundMessage_Aux{
			Aux: &aleapb.CobaltAux{
				Value: value,
			},
		},
	})
}

func BroadcastConf(m dsl.Module, config *abaConfig.Config, roundNumber uint64, values aleapb.CobaltValueSet) {
	broadcastRoundMessage(m, config, &aleapb.CobaltRoundMessage{
		RoundNumber: roundNumber,
		Type: &aleapb.CobaltRoundMessage_Conf{
			Conf: &aleapb.CobaltConf{
				Values: values,
			},
		},
	})
}

func BroadcastCoinTossShare(m dsl.Module, config *abaConfig.Config, roundNumber uint64, coinShare []byte) {
	broadcastRoundMessage(m, config, &aleapb.CobaltRoundMessage{
		RoundNumber: roundNumber,
		Type: &aleapb.CobaltRoundMessage_CoinTossShare{
			CoinTossShare: &aleapb.CobaltCoinTossShare{
				CoinShare: coinShare,
			},
		},
	})
}

func UponFinishMessage(m dsl.Module, config *abaConfig.Config, handler func(from t.NodeID, value bool) error) {
	UponMessage(m, config, func(from t.NodeID, msg *aleapb.CobaltABBA) error {
		mWrapper, ok := msg.Type.(*aleapb.CobaltABBA_Finish)
		if !ok {
			return nil
		}

		m := mWrapper.Finish
		return handler(from, m.Value)
	})
}

func UponRoundMessage(m dsl.Module, config *abaConfig.Config, handler func(from t.NodeID, msg *aleapb.CobaltRoundMessage) error) {
	UponMessage(m, config, func(from t.NodeID, msg *aleapb.CobaltABBA) error {
		mWrapper, ok := msg.Type.(*aleapb.CobaltABBA_RoundMessage)
		if !ok {
			return nil
		}

		m := mWrapper.RoundMessage
		return handler(from, m)
	})
}

func UponMessage(m dsl.Module, config *abaConfig.Config, handler func(from t.NodeID, msg *aleapb.CobaltABBA) error) {
	dsl.UponEvent[*eventpb.Event_Alea](m, func(ev *aleapb.Event) error {
		evWrapper, ok := ev.Type.(aleapb.Event_TypeWrapper[aleapb.CobaltABBAMessageRecvd])
		if !ok {
			return nil
		}

		unwrapped := evWrapper.Unwrap()
		if unwrapped.InstanceId != config.InstanceId {
			panic("ABA message routed to wrong ABA instance")
		}

		return handler(t.NodeID(unwrapped.From), unwrapped.Message)
	})
}

func broadcastMessage(m dsl.Module, config *abaConfig.Config, msg *aleapb.CobaltABBA) {
	msgPrepared := prepareMessage(config, msg)
	// TODO: use better broadcast
	dsl.SendMessage(m, config.NetModuleID, msgPrepared, config.Members)
}

func broadcastRoundMessage(m dsl.Module, config *abaConfig.Config, msg *aleapb.CobaltRoundMessage) {
	msgPrepared := prepareRoundMessage(config, msg)
	// TODO: use better broadcast
	dsl.SendMessage(m, config.NetModuleID, msgPrepared, config.Members)
}

func prepareRoundMessage(config *abaConfig.Config, msg *aleapb.CobaltRoundMessage) *messagepb.Message {
	return prepareMessage(config, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_RoundMessage{
			RoundMessage: msg,
		},
	})
}

func prepareMessage(config *abaConfig.Config, msg *aleapb.CobaltABBA) *messagepb.Message {
	return &messagepb.Message{
		Type: &messagepb.Message_Alea{
			Alea: &aleapb.AleaMessage{
				Type: &aleapb.AleaMessage_Agreement{
					Agreement: &aleapb.AgreementMsg{
						InstanceId: config.InstanceId,
						Message:    msg,
					},
				},
			},
		},
	}
}
