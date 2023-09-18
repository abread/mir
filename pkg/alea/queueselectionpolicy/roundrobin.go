package queueselectionpolicy

import (
	"github.com/fxamacker/cbor/v2"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/serializing"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

// The RoundRobinQueuePolicy is a trivial slot selection policy, which round-robins between all
// nodes in each ag round and advances through each queue sequentially, not deciding a slot until the
// previous one in the same queue has been decided.
type RoundRobinQueuePolicy struct {
	Membership   []t.NodeID
	AgQueueHeads []aleatypes.QueueSlot
	AgRound      uint64
}

func NewRoundRobinQueuePolicy(membership *trantorpbtypes.Membership) *RoundRobinQueuePolicy {
	return &RoundRobinQueuePolicy{
		Membership:   maputil.GetSortedKeys(membership.Nodes),
		AgQueueHeads: make([]aleatypes.QueueSlot, len(membership.Nodes)),
	}
}

func (pol *RoundRobinQueuePolicy) Slot(sn uint64) (bcpbtypes.Slot, bool) {
	if sn >= pol.AgRound+uint64(len(pol.Membership)) {
		// sn is too far in the future
		return bcpbtypes.Slot{}, false
	}

	queueIdx := sn % uint64(len(pol.Membership))
	return bcpbtypes.Slot{
		QueueIdx:  aleatypes.QueueIdx(queueIdx),
		QueueSlot: pol.AgQueueHeads[int(queueIdx)],
	}, true
}

func (pol *RoundRobinQueuePolicy) NextPossibleSnForSlot(slot bcpbtypes.Slot) (uint64, bool) {
	if slot.QueueSlot < pol.AgQueueHeads[slot.QueueIdx] {
		// slot was already delivered
		return 0, false
	} else if slot.QueueSlot > pol.AgQueueHeads[slot.QueueIdx] {
		// slot is too far in the future
		return 0, false
	}

	currentRoundQueue := aleatypes.QueueIdx(pol.AgRound % uint64(len(pol.Membership)))
	if slot.QueueIdx < currentRoundQueue {
		return pol.AgRound + uint64(aleatypes.QueueIdx(len(pol.Membership))+currentRoundQueue-slot.QueueIdx), true
	} else if slot.QueueIdx > currentRoundQueue {
		return pol.AgRound + uint64(slot.QueueIdx-currentRoundQueue), true
	} else {
		return pol.AgRound, true
	}
}

func (pol *RoundRobinQueuePolicy) NextSlot() bcpbtypes.Slot {
	queueIdx := pol.AgRound % uint64(len(pol.Membership))
	return bcpbtypes.Slot{
		QueueIdx:  aleatypes.QueueIdx(queueIdx),
		QueueSlot: pol.AgQueueHeads[int(queueIdx)],
	}
}

func (pol *RoundRobinQueuePolicy) QueueHead(queueIdx aleatypes.QueueIdx) aleatypes.QueueSlot {
	return pol.AgQueueHeads[int(queueIdx)]
}

func (pol *RoundRobinQueuePolicy) SlotInQueueHead(slot bcpbtypes.Slot) bool {
	return pol.AgQueueHeads[slot.QueueIdx] == slot.QueueSlot
}

func (pol *RoundRobinQueuePolicy) SlotDelivered(slot bcpbtypes.Slot) bool {
	return pol.AgQueueHeads[slot.QueueIdx] > slot.QueueSlot
}

func (pol *RoundRobinQueuePolicy) DeliverSn(sn uint64, agDecision bool) (bcpbtypes.Slot, error) {
	if sn != pol.AgRound {
		return bcpbtypes.Slot{}, es.Errorf("roundrobinqueuepolicy: received sn out of order")
	}

	queueIdx := pol.AgRound % uint64(len(pol.Membership))
	slot := bcpbtypes.Slot{
		QueueIdx:  aleatypes.QueueIdx(queueIdx),
		QueueSlot: pol.AgQueueHeads[int(queueIdx)],
	}

	if agDecision {
		// delivered slot, update queue head
		pol.AgQueueHeads[queueIdx]++
	}

	// update round number
	pol.AgRound++

	return slot, nil
}

func (pol *RoundRobinQueuePolicy) Bytes() ([]byte, error) {
	ser, err := getEncMode().Marshal(pol)
	if err != nil {
		return nil, err
	}
	out := serializing.Uint64ToBytes(uint64(RoundRobin))
	out = append(out, ser...)
	return out, nil
}

func RoundRobinQueuePolicyFromBytes(data []byte) (*RoundRobinQueuePolicy, error) {
	b := &RoundRobinQueuePolicy{}
	err := cbor.Unmarshal(data, b)
	return b, err
}
