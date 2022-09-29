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

func NewInfVec[T any](capacity int) InfVec[T] {
	return InfVec[T]{
		elements:   make([]T, capacity),
		slotStates: make([]slotState, capacity),
		lowerBound: 0,
	}
}

func (v *InfVec[T]) Get(n uint64) (*T, bool) {
	if v.NInOperationalRange(n) && *v.nthSlotState(n) != SLOT_FREED {
		*v.nthSlotState(n) = SLOT_OCCUPIED
		return v.nthUnchecked(n), true
	} else {
		return nil, false
	}
}

func (v *InfVec[T]) Free(n uint64) {
	if v.NInOperationalRange(n) {
		var zero T
		*v.nthUnchecked(n) = zero
		*v.nthSlotState(n) = SLOT_FREED

		// whenever possible, advance lower bound to allow free slot reuse with future elements
		for *v.nthSlotState(v.lowerBound) == SLOT_FREED {
			*v.nthSlotState(v.lowerBound) = SLOT_UNUSED
			v.lowerBound += 1
		}
	} else if n > v.UpperBound() {
		panic("cannot free a never-inserted element")
	}
}

func (v *InfVec[T]) TryPut(n uint64, t T) bool {
	if *v.nthSlotState(n) == SLOT_UNUSED {
		*v.nthUnchecked(n) = t
		return true
	} else {
		return false
	}
}

func (v *InfVec[T]) GetOrDefault(n uint64, mutator func(*T), defaultProducer func() T) (*T, bool) {

	slotState := v.nthSlotState(n)

	switch *slotState {
	case SLOT_UNUSED:
		*v.nthUnchecked(n) = defaultProducer()
		*slotState = SLOT_OCCUPIED
	case SLOT_OCCUPIED:
		mutator(v.nthUnchecked(n))
	case SLOT_FREED:
		// noop
	default:
		panic("invalid slot state in InfVec")
	}

	return v.nthUnchecked(n), *slotState != SLOT_FREED
}

func (v *InfVec[T]) LowerBound() uint64 {
	return v.lowerBound
}

func (v *InfVec[T]) UpperBound() uint64 {
	return v.lowerBound + uint64(len(v.elements))
}

func (v *InfVec[T]) Len() int {
	count := 0

	for _, slotState := range v.slotStates {
		if slotState == SLOT_OCCUPIED {
			count += 1
		}
	}

	return count
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
