package bcqueue

type SlotStatus uint8

const (
	// slot not used yet
	SlotUnused SlotStatus = iota

	// slot currently occupied by data
	SlotInUse

	// slot previously used, but its data was discarded and cannot be used again
	SlotFreed
)

type WindowSizeController struct {
	minSlot uint64
	minIdx  int

	slotStatusRing []SlotStatus
}

func NewWindowSizeController(viewSize int) *WindowSizeController {
	s := &WindowSizeController{
		minSlot:        0,
		slotStatusRing: make([]SlotStatus, viewSize),
	}

	for i := range s.slotStatusRing {
		s.slotStatusRing[i] = SlotUnused
	}

	return s
}

func (s *WindowSizeController) GetSlotStatus(slot uint64) SlotStatus {
	if slot < s.minSlot {
		return SlotFreed
	} else if slot >= s.minSlot+uint64(len(s.slotStatusRing)) {
		return SlotUnused
	}

	return *s.uncheckedSlotStatusRingPtr(slot)
}

func (s *WindowSizeController) IsSlotInView(slot uint64) bool {
	return slot >= s.minSlot && slot < s.minSlot+uint64(len(s.slotStatusRing))
}

func (s *WindowSizeController) IsSlotUnused(slot uint64) bool {
	return s.GetSlotStatus(slot) == SlotUnused
}

func (s *WindowSizeController) IsSlotInUse(slot uint64) bool {
	return s.GetSlotStatus(slot) == SlotInUse
}

func (s *WindowSizeController) IsSlotFreed(slot uint64) bool {
	return s.GetSlotStatus(slot) == SlotFreed
}

func (s *WindowSizeController) TryAcquire(slot uint64) bool {
	if !(s.IsSlotInView(slot) && s.IsSlotUnused(slot)) {
		return false
	}

	*s.uncheckedSlotStatusRingPtr(slot) = SlotInUse

	return true
}

func (s *WindowSizeController) Acquire(slot uint64) {
	if !s.TryAcquire(slot) {
		panic("failed to acquire window slot")
	}
}

func (s *WindowSizeController) TryFree(slot uint64) bool {
	if slot < s.minSlot {
		return true // nothing to do, already freed by definition
	}

	if !s.IsSlotInView(slot) {
		return false
	}

	*s.uncheckedSlotStatusRingPtr(slot) = SlotFreed

	// advance view to expose unused slots
	for s.IsSlotFreed(s.minSlot) {
		*s.uncheckedSlotStatusRingPtr(s.minSlot) = SlotUnused
		s.minSlot++
		s.minIdx++
	}

	return true
}

func (s *WindowSizeController) Free(slot uint64) {
	s.TryFree(slot)
}

func (s *WindowSizeController) uncheckedSlotStatusRingPtr(slot uint64) *SlotStatus {
	return &s.slotStatusRing[(int(slot-s.minSlot)+s.minIdx)%len(s.slotStatusRing)]
}
