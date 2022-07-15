package alea

import (
	"fmt"

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

type pendingMessageTree struct {
	message              *messagepb.Message
	destinations         []t.NodeID
	timeSinceLastAttempt uint64
	children             map[uint64]*pendingMessageTree
}

func NewPendingMessageTree() pendingMessageTree {
	return pendingMessageTree{
		message:      nil,
		destinations: nil,
		children:     make(map[uint64]*pendingMessageTree, 0),
	}
}

func (tree *pendingMessageTree) Put(id []uint64, destinations []t.NodeID, msg *messagepb.Message) error {
	if len(id) == 0 {
		if tree.message != nil {
			return fmt.Errorf("Message already present: %+v for nodes %v", tree.message, tree.destinations)
		} else {
			tree.message = msg
			tree.destinations = destinations
			return nil
		}
	} else {
		return tree.subtree(id[0]).Put(id[1:], destinations, msg)
	}
}

func (tree *pendingMessageTree) MarkReceivedAll(id []uint64) {
	if len(id) == 0 {
		tree.message = nil
		tree.destinations = nil
		tree.children = nil
	} else if len(id) == 1 {
		delete(tree.children, id[0])
	} else if st, exists := tree.children[id[0]]; exists {
		st.MarkReceivedAll(id[1:])
	}
}

func (tree *pendingMessageTree) MarkReceivedOne(id []uint64, n t.NodeID) {
	if len(id) == 0 {
		tree.destinations = DeleteSliceEl(tree.destinations, n)

		for _, st := range tree.children {
			st.MarkReceivedOne(id, n)
		}
	} else {
		if st, exists := tree.children[id[0]]; exists {
			st.MarkReceivedOne(id[1:], n)
		}
	}
}

func (tree *pendingMessageTree) ToEventList() *events.EventList {
	list := &events.EventList{}

	// TODO: track message timers individually to prevent a thundering herd of retransmissions (also allows for smarter timeouts)
	if tree.message != nil && len(tree.destinations) > 0 {
		list.PushBack(events.SendMessage(netModuleName, tree.message, tree.destinations))
	}

	for _, st := range tree.children {
		list.PushBackList(st.ToEventList())
	}

	return list
}

// clear empty tree nodes
func (tree *pendingMessageTree) Gc() bool {
	for idx, st := range tree.children {
		if st.Gc() {
			delete(tree.children, idx)
		}
	}

	return tree.message == nil && len(tree.children) == 0
}

func (tree *pendingMessageTree) subtree(idPart uint64) *pendingMessageTree {
	if st, exists := tree.children[idPart]; exists {
		return st
	} else {
		st := NewPendingMessageTree()
		tree.children[idPart] = &st
		return &st
	}
}

func DeleteSliceEl(s []t.NodeID, toDelete t.NodeID) []t.NodeID {
	for idx, v := range s {
		if v == toDelete {
			s[idx] = s[len(s)-1]
			return s[:len(s)-1]
		}
	}

	return s
}

// TODO: implement retry & ack logic
