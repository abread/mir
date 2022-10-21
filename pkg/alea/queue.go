package alea

type slotState uint8

const (
	UnusedSlot slotState = iota
	OccupiedSlot
	FreedSlot
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
	if v.NInOperationalRange(n) && *v.nthSlotState(n) != FreedSlot {
		*v.nthSlotState(n) = OccupiedSlot
		return v.nthUnchecked(n), true
	}

	return nil, false
}

func (v *InfVec[T]) Free(n uint64) {
	if v.NInOperationalRange(n) {
		var zero T
		*v.nthUnchecked(n) = zero
		*v.nthSlotState(n) = FreedSlot

		// whenever possible, advance lower bound to allow free slot reuse with future elements
		for *v.nthSlotState(v.lowerBound) == FreedSlot {
			*v.nthSlotState(v.lowerBound) = UnusedSlot
			v.lowerBound++
		}
	} else if n > v.UpperBound() {
		panic("cannot free a never-inserted element")
	}
}

func (v *InfVec[T]) TryPut(n uint64, t T) bool {
	if *v.nthSlotState(n) == UnusedSlot {
		*v.nthUnchecked(n) = t
		return true
	}

	return false
}

func (v *InfVec[T]) GetOrDefault(n uint64, mutator func(*T), defaultProducer func() T) (*T, bool) {

	slotState := v.nthSlotState(n)

	switch *slotState {
	case UnusedSlot:
		*v.nthUnchecked(n) = defaultProducer()
		*slotState = OccupiedSlot
	case OccupiedSlot:
		mutator(v.nthUnchecked(n))
	case FreedSlot:
		// noop
	default:
		panic("invalid slot state in InfVec")
	}

	return v.nthUnchecked(n), *slotState != FreedSlot
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
		if slotState == OccupiedSlot {
			count++
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
