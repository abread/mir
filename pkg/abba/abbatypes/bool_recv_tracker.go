package abbatypes

import t "github.com/filecoin-project/mir/pkg/types"

type BoolRecvTrackers struct {
	falseTracker RecvTracker
	trueTracker  RecvTracker
}

func NewBoolRecvTrackers(cap int) BoolRecvTrackers {
	var res BoolRecvTrackers
	res.Reset(cap)
	return res
}

func (m *BoolRecvTrackers) Register(v bool, node t.NodeID) bool {
	if v {
		return m.trueTracker.Register(node)
	}
	return m.falseTracker.Register(node)
}

func (m *BoolRecvTrackers) Reset(cap int) {
	// the idiomatic way to this is reallocate the map
	m.trueTracker = make(RecvTracker, cap)
	m.falseTracker = make(RecvTracker, cap)
}
