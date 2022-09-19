package alea

import (
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
)

// ============================
// Messages
// ============================

func CobaltAbbaFinish(id uint64, v bit) *messagepb.Message {
	return prepareAbaMessage(id, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_Finish{
			Finish: &aleapb.CobaltFinish{
				Value: v.Pb(),
			},
		},
	})
}

func CobaltAbbaInit(id uint64, r uint64, v bit) *messagepb.Message {
	return prepareAbaMessage(id, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_Init{
			Init: &aleapb.CobaltInit{
				Round: r,
				Value: v.Pb(),
			},
		},
	})
}

func CobaltAbbaAux(id uint64, r uint64, v bit) *messagepb.Message {
	return prepareAbaMessage(id, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_Aux{
			Aux: &aleapb.CobaltAux{
				Round: r,
				Value: v.Pb(),
			},
		},
	})
}

func CobaltAbbaConf(id uint64, r uint64, c BitSet) *messagepb.Message {
	return prepareAbaMessage(id, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_Conf{
			Conf: &aleapb.CobaltConf{
				Round:  r,
				Config: c.Pb(),
			},
		},
	})
}

func CobaltAbbaCoinShare(id uint64, r uint64, share threshSigShare) *messagepb.Message {
	return prepareAbaMessage(id, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_CoinShare{
			CoinShare: &aleapb.CobaltCoinShare{
				Round:     r,
				CoinShare: share.Pb(),
			},
		},
	})
}

func prepareAbaMessage(id uint64, msg *aleapb.CobaltABBA) *messagepb.Message {
	return prepareAleaMessage(&aleapb.AleaMessage{
		Type: &aleapb.AleaMessage_Agreement{
			Agreement: &aleapb.Agreement{
				AgreementRound: id,
				Msg:            msg,
			},
		},
	})
}

func prepareAleaMessage(msg *aleapb.AleaMessage) *messagepb.Message {
	return &messagepb.Message{
		Type: &messagepb.Message_Alea{
			Alea: msg,
		},
	}
}

// ============================
// Events
// ============================

func VCBCDeliver(id MsgId) *eventpb.Event {
	return prepareAleaEvent(&aleapb.AleaEvent{
		Type: &aleapb.AleaEvent_VcbcDeliver{
			VcbcDeliver: &aleapb.VCBCDeliver{
				Id: id.Pb(),
			},
		},
	})
}

func AbaDeliver(agreementRound uint64, result bit) *eventpb.Event {
	return prepareAleaEvent(&aleapb.AleaEvent{
		Type: &aleapb.AleaEvent_AbaDeliver{
			AbaDeliver: &aleapb.CobaltABBADeliver{
				AgreementRound: agreementRound,
				Result:         result.Pb(),
			},
		},
	})
}

func prepareAleaEvent(ev *aleapb.AleaEvent) *eventpb.Event {
	return &eventpb.Event{
		Type: &eventpb.Event_Alea{
			Alea: ev,
		},
	}
}
