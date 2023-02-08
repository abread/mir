package abbatypes

import t "github.com/filecoin-project/mir/pkg/types"

type RecvTracker map[t.NodeID]struct{}

func (r RecvTracker) Register(node t.NodeID) bool {
	if _, present := r[node]; present {
		return false
	}

	r[node] = struct{}{}
	return true
}
