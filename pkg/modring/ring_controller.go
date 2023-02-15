package modring

type RingSlotStatus uint8

const (
	// slot not used yet
	RingSlotFuture RingSlotStatus = iota

	// slot currently occupied by data
	RingSlotCurrent

	// slot previously used, but its data was discarded and cannot be used again
	RingSlotPast
)

type RingController struct {
	minSlot uint64
	minIdx  int

	slotStatus []RingSlotStatus
}

func NewRingController(viewSize int) RingController {
	s := RingController{
		minSlot:    0,
		slotStatus: make([]RingSlotStatus, viewSize),
	}

	for i := range s.slotStatus {
		s.slotStatus[i] = RingSlotFuture
	}

	return s
}

func (s *RingController) GetSlotStatus(slot uint64) RingSlotStatus {
	if slot < s.minSlot {
		return RingSlotPast
	} else if slot >= s.minSlot+uint64(len(s.slotStatus)) {
		return RingSlotFuture
	}

	return *s.uncheckedSlotStatusRingPtr(slot)
}

func (s *RingController) isSlotInView(slot uint64) bool {
	return slot >= s.minSlot && slot < s.minSlot+uint64(len(s.slotStatus))
}

func (s *RingController) IsFutureSlot(slot uint64) bool {
	return s.GetSlotStatus(slot) == RingSlotFuture
}

func (s *RingController) IsCurrentSlot(slot uint64) bool {
	return s.GetSlotStatus(slot) == RingSlotCurrent
}

func (s *RingController) IsPastSlot(slot uint64) bool {
	return s.GetSlotStatus(slot) == RingSlotPast
}

func (s *RingController) TryAcquire(slot uint64) bool {
	if !(s.isSlotInView(slot) && s.IsFutureSlot(slot)) {
		return false
	}

	*s.uncheckedSlotStatusRingPtr(slot) = RingSlotCurrent

	return true
}

func (s *RingController) Acquire(slot uint64) {
	if !s.TryAcquire(slot) {
		panic("failed to acquire window slot")
	}
}

func (s *RingController) TryFree(slot uint64) bool {
	if slot < s.minSlot {
		return true // nothing to do, already freed by definition
	}

	if !s.isSlotInView(slot) {
		return false
	}

	*s.uncheckedSlotStatusRingPtr(slot) = RingSlotPast

	// advance view to expose unused slots
	for s.IsPastSlot(s.minSlot) {
		*s.uncheckedSlotStatusRingPtr(s.minSlot) = RingSlotFuture
		s.minSlot++
		s.minIdx = (s.minIdx + 1) % len(s.slotStatus)
	}

	return true
}

func (s *RingController) Free(slot uint64) {
	s.TryFree(slot)
}

func (s *RingController) uncheckedSlotStatusRingPtr(slot uint64) *RingSlotStatus {
	return &s.slotStatus[(int(slot-s.minSlot)+s.minIdx)%len(s.slotStatus)]
}
