package aba

import (
	abaConfig "github.com/filecoin-project/mir/pkg/alea/aba/config"
	t "github.com/filecoin-project/mir/pkg/types"
)

type MsgSupportTrackerBool struct {
	selfID       t.NodeID
	receivedFrom map[t.NodeID]void
	valueCount   map[bool]int
}

func newMsgSupportTrackerBool(config *abaConfig.Config) *MsgSupportTrackerBool {
	tracker := &MsgSupportTrackerBool{
		selfID:       config.SelfNodeID(),
		receivedFrom: make(map[t.NodeID]void, len(config.Members)),
		valueCount:   make(map[bool]int, 2),
	}

	tracker.valueCount[false] = 0
	tracker.valueCount[true] = 0

	return tracker
}

func (tracker *MsgSupportTrackerBool) RegisterMsg(from t.NodeID, value bool) bool {
	if _, alreadyRegistered := tracker.receivedFrom[from]; !alreadyRegistered {
		tracker.receivedFrom[from] = void{}
		tracker.valueCount[value]++

		return true
	}

	return false
}

func (tracker *MsgSupportTrackerBool) SelfAlreadySentAnything() bool {
	_, alreadyRegistered := tracker.receivedFrom[tracker.selfID]
	return alreadyRegistered
}

func (tracker *MsgSupportTrackerBool) ValueCountFor(value bool) int {
	return tracker.valueCount[value]
}
