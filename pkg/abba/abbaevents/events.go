package abbaevents

import (
	"github.com/filecoin-project/mir/pkg/pb/abbapb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule t.ModuleID, input bool) *eventpb.Event {
	return Event(destModule, &abbapb.Event{
		Type: &abbapb.Event_InputValue{
			InputValue: &abbapb.InputValue{
				Input: input,
			},
		},
	})
}

func Deliver(destModule t.ModuleID, result bool) *eventpb.Event {
	return Event(destModule, &abbapb.Event{
		Type: &abbapb.Event_Deliver{
			Deliver: &abbapb.Deliver{
				Result: result,
			},
		},
	})
}

func Event(destModule t.ModuleID, ev *abbapb.Event) *eventpb.Event {
	return &eventpb.Event{
		DestModule: destModule.Pb(),

		Type: &eventpb.Event_Abba{
			Abba: ev,
		},
	}
}
