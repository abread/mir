package modring

import "fmt"

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

func (s *RingController) IsSlotInView(slot uint64) bool {
	return slot >= s.minSlot && slot < s.minSlot+uint64(len(s.slotStatus))
}

func (s *RingController) GetCurrentView() (uint64, uint64) {
	return s.minSlot, s.minSlot + uint64(len(s.slotStatus))
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

func (s *RingController) MarkCurrentIfFuture(slot uint64) bool {
	if !(s.IsSlotInView(slot) && s.IsFutureSlot(slot)) {
		return false
	}

	*s.uncheckedSlotStatusRingPtr(slot) = RingSlotCurrent

	return true
}

func (s *RingController) MarkPast(slot uint64) error {
	if slot < s.minSlot {
		return nil // nothing to do, already freed by definition
	}

	if !s.IsSlotInView(slot) {
		return fmt.Errorf("cannot mark slot %d as past", slot)
	}

	*s.uncheckedSlotStatusRingPtr(slot) = RingSlotPast

	// advance view to expose unused slots
	for s.IsPastSlot(s.minSlot) {
		*s.uncheckedSlotStatusRingPtr(s.minSlot) = RingSlotFuture
		s.minSlot++
		s.minIdx = (s.minIdx + 1) % len(s.slotStatus)
	}

	return nil
}

func (s *RingController) AdvanceViewToSlot(slot uint64) error {
	beginSlot, endSlot := s.GetCurrentView()
	if slot < beginSlot {
		return fmt.Errorf("cannot advance view: view already advanced past slot %d", slot)
	} else if slot < endSlot {
		return nil // slot already in view
	}

	for slot := beginSlot; slot < endSlot; slot++ {
		if *s.uncheckedSlotStatusRingPtr(slot) == RingSlotCurrent {
			return fmt.Errorf("cannot advance view: slot %d still in use", slot)
		}
	}

	s.minIdx = int(slot % uint64(len(s.slotStatus)))
	s.minSlot = slot

	// all slots start out as future
	for i := 0; i < len(s.slotStatus); i++ {
		s.slotStatus[i] = RingSlotFuture
	}
	return nil
}

func (s *RingController) uncheckedSlotStatusRingPtr(slot uint64) *RingSlotStatus {
	return &s.slotStatus[int(slot%uint64(len(s.slotStatus)))]
}
