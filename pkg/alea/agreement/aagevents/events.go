package aagevents

import (
	"github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
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
