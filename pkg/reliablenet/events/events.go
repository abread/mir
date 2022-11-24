package events

import (
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/reliablenetpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func SendMessage(destModule t.ModuleID, id []byte, message *messagepb.Message, destinations []t.NodeID) *eventpb.Event {
	return Event(destModule, &reliablenetpb.Event{
		Type: &reliablenetpb.Event_SendMessage{
			SendMessage: &reliablenetpb.SendMessage{
				MsgId:        id,
				Msg:          message,
				Destinations: t.NodeIDSlicePb(destinations),
			},
		},
	})
}

func Ack(destModule t.ModuleID, msgDestModule t.ModuleID, msgID []byte, msgSource t.NodeID) *eventpb.Event {
	return Event(destModule, &reliablenetpb.Event{
		Type: &reliablenetpb.Event_Ack{
			Ack: &reliablenetpb.Ack{
				DestModule: msgDestModule.Pb(),
				MsgId:      msgID,
				Source:     msgSource.Pb(),
			},
		},
	})
}

func MarkModuleMsgsRecvd(destModule t.ModuleID, msgDestModule t.ModuleID, destinations []t.NodeID) *eventpb.Event {
	return Event(destModule, &reliablenetpb.Event{
		Type: &reliablenetpb.Event_MarkModuleMsgsRecvd{
			MarkModuleMsgsRecvd: &reliablenetpb.MarkModuleMsgsRecvd{
				DestModule:   msgDestModule.Pb(),
				Destinations: t.NodeIDSlicePb(destinations),
			},
		},
	})
}

func MarkRecvd(destModule t.ModuleID, msgDestModule t.ModuleID, msgID []byte, destinations []t.NodeID) *eventpb.Event {
	return Event(destModule, &reliablenetpb.Event{
		Type: &reliablenetpb.Event_MarkRecvd{
			MarkRecvd: &reliablenetpb.MarkRecvd{
				DestModule:   msgDestModule.Pb(),
				MsgId:        msgID,
				Destinations: t.NodeIDSlicePb(destinations),
			},
		},
	})
}

func RetransmitAll(destModule t.ModuleID) *eventpb.Event {
	return Event(destModule, &reliablenetpb.Event{
		Type: &reliablenetpb.Event_RetransmitAll{
			RetransmitAll: &reliablenetpb.RetransmitAll{},
		},
	})
}

func Event(destModule t.ModuleID, ev *reliablenetpb.Event) *eventpb.Event {
	return &eventpb.Event{
		Type: &eventpb.Event_ReliableNet{
			ReliableNet: ev,
		},
		DestModule: destModule.Pb(),
	}
}
