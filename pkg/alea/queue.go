package alea

type slotState uint8

const (
	SLOT_UNUSED   slotState = 0
	SLOT_OCCUPIED           = 1
	SLOT_FREED              = 2
)

type InfVec[T any] struct {
	elements   []T
	slotStates []slotState
	lowerBound uint64
}

func NewInfVec[T any](capacity int) *InfVec[T] {
	return &InfVec[T]{
		elements:   make([]T, capacity),
		slotStates: make([]slotState, capacity),
		lowerBound: 0,
	}
}

func (v *InfVec[T]) LowerBound() uint64 {
	return v.lowerBound
}

func (v *InfVec[T]) UpperBound() uint64 {
	return v.lowerBound + uint64(len(v.elements))
}

func (v *InfVec[T]) NInOperationalRange(idx uint64) bool {
	return idx >= v.LowerBound() && idx < v.UpperBound()
}

func (v *InfVec[T]) nthIdx(n uint64) int {
	return int(n % uint64(len(v.elements)))
}

func (v *InfVec[T]) nthUnchecked(n uint64) *T {
	return &v.elements[v.nthIdx(n)]
}

func (v *InfVec[T]) nthSlotState(n uint64) *slotState {
	return &v.slotStates[v.nthIdx(n)]
}

func (v *InfVec[T]) Get(n uint64) (*T, bool) {
	if v.NInOperationalRange(n) && *v.nthSlotState(n) != SLOT_FREED {
		*v.nthSlotState(n) = SLOT_OCCUPIED
		return v.nthUnchecked(n), true
	} else {
		return nil, false
	}
}

func (v *InfVec[T]) Free(idx uint64) {
	if v.NInOperationalRange(idx) {
		var zero T
		*v.nthUnchecked(idx) = zero
		*v.nthSlotState(idx) = SLOT_FREED

		// whenever possible, advance lower bound to allow free slot reuse with future elements
		for *v.nthSlotState(v.lowerBound) == SLOT_FREED {
			*v.nthSlotState(v.lowerBound) = SLOT_UNUSED
			v.lowerBound += 1
		}
	} else if idx > v.UpperBound() {
		panic("cannot free a never-inserted element")
	}
}
