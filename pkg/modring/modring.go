package modring

import (
	"fmt"
	"strconv"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type Module struct {
	ownID          t.ModuleID
	generator      ModuleGenerator
	pastMsgHandler PastMessagesHandler

	ringController RingController
	ring           []modules.PassiveModule

	logger logging.Logger
}

func New(ownID t.ModuleID, ringSize int, params ModuleParams, logger logging.Logger) *Module {
	if logger == nil {
		logger = logging.ConsoleErrorLogger
	}

	return &Module{
		ownID:          ownID,
		generator:      params.Generator,
		pastMsgHandler: params.PastMsgHandler,

		ringController: NewRingController(ringSize),
		ring:           make([]modules.PassiveModule, ringSize),

		logger: logger,
	}
}

func (m *Module) ImplementsModule() {}

func (m *Module) ApplyEvents(eventsIn *events.EventList) (*events.EventList, error) {
	// Note: this method is similar to ApplyEvents, but it only parallelizes execution across different
	// modules, not all events.
	subEventsIn, pastMsgs := m.splitEventsByDest(eventsIn)

	resultChan := make(chan *events.EventList)
	errorChan := make(chan error)

	// The event processing results will be aggregated here
	var eventsOut *events.EventList
	var firstError error

	// Start one concurrent worker for each submodule
	// Note that the processing starts concurrently for all eventlists, and only the writing of the results is synchronized.
	nSubs := 0
	for i := 0; i < len(subEventsIn); i++ {
		if subEventsIn[i].Len() == 0 {
			continue
		}

		mod, evsOut, err := m.getSubByRingIdx(i)
		if err != nil {
			firstError = err
			break
		}
		if mod == nil {
			continue // module out of view
		}

		go func(modring *Module, evsIn *events.EventList, initialEvsOut *events.EventList, j int) {
			// Apply the input event, catching potential panics.
			res, err := mod.ApplyEvents(evsIn)

			if err == nil {
				resultChan <- initialEvsOut.PushBackList(res)
			} else {
				errorChan <- fmt.Errorf("failed to process submodule #%d events: %w", j, err)
			}

		}(m, &subEventsIn[i], evsOut, i)
		nSubs++
	}

	eventsOut = &events.EventList{}

	// For each submodule, read the processing result from the common channels and aggregate it with the rest.
	for i := 0; i < nSubs; i++ {
		select {
		case evList := <-resultChan:
			eventsOut.PushBackList(evList)
		case err := <-errorChan:
			if firstError == nil {
				firstError = err
			}
		}
	}

	// Return the resulting events or an error.
	if firstError != nil {
		return nil, firstError
	}

	if len(pastMsgs) > 0 && m.pastMsgHandler != nil {
		evsOut, err := (m.pastMsgHandler)(pastMsgs)
		if err != nil {
			return nil, err
		}

		eventsOut.PushBackList(evsOut)
	}

	return eventsOut, nil
}

func (m *Module) splitEventsByDest(eventsIn *events.EventList) ([]events.EventList, []*modringpbtypes.PastMessage) {
	subEventsIn := make([]events.EventList, len(m.ring))
	pastMsgs := make([]*modringpbtypes.PastMessage, 0)

	eventsInIter := eventsIn.Iterator()
	for event := eventsInIter.Next(); event != nil; event = eventsInIter.Next() {
		subIDStr := t.ModuleID(event.DestModule).StripParent(m.ownID).Top()
		subID, err := strconv.ParseUint(string(subIDStr), 10, 64)
		if err != nil {
			m.logger.Log(logging.LevelWarn, "event received for invalid submodule index", "submoduleIdx", subIDStr)
			continue
		}

		if m.ringController.IsPastSlot(subID) {
			if ev, ok := event.Type.(*eventpb.Event_MessageReceived); ok {
				pastMsgs = append(pastMsgs, &modringpbtypes.PastMessage{
					DestId:  subID,
					From:    t.NodeID(ev.MessageReceived.From),
					Message: messagepbtypes.MessageFromPb(ev.MessageReceived.Msg),
				})
			}
			continue
		}

		if !m.ringController.IsSlotInView(subID) {
			m.logger.Log(logging.LevelWarn, "event received for out of view submodule", "submoduleID", subID)
			continue
		}

		idx := int(subID % uint64(len(m.ring)))
		subEventsIn[idx].PushBack(event)
	}

	return subEventsIn, pastMsgs
}

func (m *Module) getSubByRingIdx(ringIdx int) (modules.PassiveModule, *events.EventList, error) {
	minSlot, _ := m.ringController.GetCurrentView()
	minIdx := minSlot % uint64(len(m.ring))

	subID := m.ringController.minSlot + uint64(ringIdx) - minIdx

	switch m.ringController.GetSlotStatus(subID) {
	case RingSlotPast:
		return nil, nil, nil
	case RingSlotCurrent:
		if m.ring[ringIdx] == nil {
			return nil, nil, fmt.Errorf("module %v disappeared", subID)
		}

		return m.ring[ringIdx], &events.EventList{}, nil
	case RingSlotFuture:
		if !m.ringController.MarkCurrentIfFuture(subID) {
			return nil, nil, nil
		}

		subFullID := m.ownID.Then(t.NewModuleIDFromInt(subID))
		sub, initialEvs, err := (m.generator)(subFullID, subID)
		if err != nil {
			return nil, nil, err
		}

		m.ring[ringIdx] = sub
		initialEvs.PushBack(events.Init(subFullID))

		evsOut, err := m.ring[ringIdx].ApplyEvents(initialEvs)
		if err != nil {
			return nil, nil, err
		}

		m.logger.Log(logging.LevelDebug, "module created", "id", subID)
		return m.ring[ringIdx], evsOut, err
	default:
		return nil, nil, fmt.Errorf("unknown slot status: %v", m.ringController.GetSlotStatus(subID))
	}
}

func (m *Module) AdvanceViewToAtLeastSubmodule(id uint64) error {
	if m.ringController.IsPastSlot(id) || m.ringController.IsSlotInView(id) {
		return nil // already advanced enough
	}

	minSlot, endSlot := m.ringController.GetCurrentView()
	for slot := minSlot; slot < endSlot; slot++ {
		// free all current slots
		if err := m.MarkSubmodulePast(id); err != nil {
			return fmt.Errorf("cannot advance view: %w", err)
		}
	}

	if err := m.ringController.AdvanceViewToSlot(id); err != nil {
		return fmt.Errorf("cannot advance view: %w", err)
	}

	m.logger.Log(logging.LevelDebug, "fast-forwarded view", "id", id, "newMinSlot", m.ringController.minSlot)

	return nil
}

func (m *Module) MarkSubmodulePast(id uint64) error {
	if m.ringController.IsPastSlot(id) {
		// already in the past
		return nil
	}

	if err := m.ringController.MarkPast(id); err != nil {
		return fmt.Errorf("cannot mark submodule as past: %w", err)
	}

	// zero slot just in case (and to let the GC do its job)
	ringIdx := int(id % uint64(len(m.ring)))
	m.ring[ringIdx] = nil

	m.logger.Log(logging.LevelDebug, "module freed", "id", id, "newMinSlot", m.ringController.minSlot)
	return nil
}
