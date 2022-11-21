package aagevents

import (
	"github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func RequestInput(destModule t.ModuleID, round uint64) *eventpb.Event {
	return Event(destModule, &agreementpb.Event{
		Type: &agreementpb.Event_RequestInput{
			RequestInput: &agreementpb.RequestInput{
				Round: round,
			},
		},
	})
}

func InputValue(destModule t.ModuleID, round uint64, input bool) *eventpb.Event {
	return Event(destModule, &agreementpb.Event{
		Type: &agreementpb.Event_InputValue{
			InputValue: &agreementpb.InputValue{
				Round: round,
				Input: input,
			},
		},
	})
}

func Deliver(destModule t.ModuleID, round uint64, decision bool) *eventpb.Event {
	return Event(destModule, &agreementpb.Event{
		Type: &agreementpb.Event_Deliver{
			Deliver: &agreementpb.Deliver{
				Round:    round,
				Decision: decision,
			},
		},
	})
}

func Event(destModule t.ModuleID, ev *agreementpb.Event) *eventpb.Event {
	return &eventpb.Event{
		DestModule: destModule.Pb(),
		Type: &eventpb.Event_AleaAgreement{
			AleaAgreement: ev,
		},
	}
}

func AbbaFinishMessage(destModule t.ModuleID, agreementRound uint64, value bool) *messagepb.Message {
	return Message(destModule, &agreementpb.Message{
		Type: &agreementpb.Message_FinishAbba{
			FinishAbba: &agreementpb.FinishAbbaMessage{
				Round: agreementRound,
				Value: value,
			},
		},
	})
}

func Message(destModule t.ModuleID, msg *agreementpb.Message) *messagepb.Message {
	return &messagepb.Message{
		DestModule: destModule.Pb(),
		Type: &messagepb.Message_AleaAgreement{
			AleaAgreement: msg,
		},
	}
}
