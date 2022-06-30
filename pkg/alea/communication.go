package alea

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func (alea *Alea) sendMessage(msg *messagepb.Message, destination t.NodeID) *events.EventList {
	return (&events.EventList{}).PushBack(events.SendMessage(netModuleName, msg, []t.NodeID{destination}))
}

func (alea *Alea) broadcastMessage(msg *messagepb.Message) *events.EventList {
	return (&events.EventList{}).PushBack(events.SendMessage(netModuleName, msg, alea.config.Membership))
}

// TODO: implement retry & ack logic
