package queueselectionpolicy

import (
	"sync"

	"github.com/fxamacker/cbor/v2"
	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/serializing"
)

type QueuePolicyType uint64

const (
	RoundRobin QueuePolicyType = iota
)

var encMode cbor.EncMode
var once sync.Once

func getEncMode() cbor.EncMode {
	once.Do(func() {
		encMode, _ = cbor.CoreDetEncOptions().EncMode()
	})
	return encMode
}

// A QueueSelectionPolicy implements the algorithm for selecting a set of leaders in each Alea epoch.
// It controls which slots are voted in which sequence number.
type QueueSelectionPolicy interface {
	// Slot returns the slot for the given sequence number (agreement round number), and true if it can be known.
	// If the slot for the given sequence number cannot currently be determined, the second return value is false.
	Slot(sn uint64) (bcpbtypes.Slot, bool)

	NextPossibleSnForSlot(slot bcpbtypes.Slot) (uint64, bool)

	NextSlot() bcpbtypes.Slot

	// DeliverSn updates the state of the policy object by announcing it that sequence number `sn`,
	// has been delivered, incrementing the head of the relevant queue if applicable.
	DeliverSn(sn uint64, agDecision bool) (bcpbtypes.Slot, error)

	SlotDelivered(slot bcpbtypes.Slot) bool

	SlotInQueueHead(slot bcpbtypes.Slot) bool

	QueueHead(queueIdx aleatypes.QueueIdx) aleatypes.QueueSlot

	// TODO: accept more signals to inform better policies?

	Bytes() ([]byte, error)
}

func QueuePolicyFromBytes(bytes []byte) (QueueSelectionPolicy, error) {
	queuePolicyType := serializing.Uint64FromBytes(bytes[0:8])

	switch QueuePolicyType(queuePolicyType) {
	case RoundRobin:
		return RoundRobinQueuePolicyFromBytes(bytes[8:])
	default:
		return nil, es.Errorf("invalid LeaderSelectionPolicy type: %v", queuePolicyType)
	}

}

func NewQueuePolicy(queuePolicyType QueuePolicyType, membership *trantorpbtypes.Membership) QueueSelectionPolicy {
	switch queuePolicyType {
	case RoundRobin:
		return NewRoundRobinQueuePolicy(membership)
	default:
		panic("invalid LeaderSelectionPolicy type")
	}
}
