package alea

import (
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
)

// ============================
// Messages
// ============================

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
