package abba

import t "github.com/filecoin-project/mir/pkg/types"

type recvTracker map[t.NodeID]struct{}

func (r recvTracker) Register(node t.NodeID) bool {
	if _, present := r[node]; present {
		return false
	}

	r[node] = struct{}{}
	return true
}
