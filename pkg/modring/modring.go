package modring

import (
	"fmt"
	"runtime/debug"
	"strconv"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/modringpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type Module struct {
	ownID     t.ModuleID
	generator ModuleGenerator

	ringController RingController
	ring           []modules.PassiveModule

	logger logging.Logger
}

func New(ownID t.ModuleID, ringSize int, params ModuleParams, logger logging.Logger) *Module {
	if logger == nil {
		logger = logging.ConsoleErrorLogger
	}

	return &Module{
		ownID:     ownID,
		generator: params.Generator,

		ringController: NewRingController(ringSize),
		ring:           make([]modules.PassiveModule, ringSize),

		logger: logger,
	}
}

func (m *Module) ImplementsModule() {}

func (m *Module) ApplyEvents(eventsIn *events.EventList) (*events.EventList, error) {
	// Note: this method is similar to ApplyEvents, but it only parallelizes execution across diferent
	// modules, not all events.
	ownEventsIn, subEventsIn := m.splitEventsByDest(eventsIn)

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
		if mod == nil {
			continue // module already out of view
			// TODO: allow custom handler for unhandled events?
		}
		if err != nil {
			firstError = err
			break
		}

		go func(modring *Module, evsIn *events.EventList, initialEvsOut *events.EventList, j int) {
			// Apply the input event, catching potential panics.
			res, err := mod.ApplyEvents(evsIn)

			if err == nil {
				resultChan <- initialEvsOut.PushBackList(res)
			} else {
				errorChan <- err
			}

		}(m, &subEventsIn[i], evsOut, i)
		nSubs++
	}

	// the modring's events will be handled in the current go routine
	ownEventsOut, ownError := modules.ApplyEventsSequentially(ownEventsIn, m.applyEvent)
	if ownError == nil {
		eventsOut = ownEventsOut
	} else {
		if firstError == nil {
			firstError = ownError
		}
		eventsOut = &events.EventList{}
	}

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
	return eventsOut, nil
}

func (m *Module) splitEventsByDest(eventsIn *events.EventList) (*events.EventList, []events.EventList) {
	ownEventsIn := &events.EventList{}
	subEventsIn := make([]events.EventList, len(m.ring))

	eventsInIter := eventsIn.Iterator()
	for event := eventsInIter.Next(); event != nil; event = eventsInIter.Next() {
		if event.DestModule == m.ownID.Pb() {
			ownEventsIn.PushBack(event)
		} else {
			subIDStr := t.ModuleID(event.DestModule).StripParent(m.ownID).Top()
			subID, err := strconv.ParseUint(string(subIDStr), 10, 64)
			if err != nil {
				m.logger.Log(logging.LevelWarn, "event received for invalid submodule index", "submoduleIdx", subIDStr)
			}

			idx := int(subID % uint64(len(m.ring)))
			subEventsIn[idx].PushBack(event)
		}
	}

	return ownEventsIn, subEventsIn
}

func (m *Module) applyEvent(event *eventpb.Event) (*events.EventList, error) {
	switch e := event.Type.(type) {
	case *eventpb.Event_Init:
		return events.EmptyList(), nil // Nothing to do at initialization.
	case *eventpb.Event_Modring:
		return m.applyModringEvent(e.Modring)
	default:
		return nil, fmt.Errorf("unsupported event type for modring module: %T", e)
	}
}

func (m *Module) getSubByRingIdx(ringIdx int) (modules.PassiveModule, *events.EventList, error) {
	subID := m.ringController.minSlot + uint64(ringIdx) - uint64(m.ringController.minIdx)

	switch m.ringController.GetSlotStatus(subID) {
	case RingSlotPast:
		return nil, nil, nil
	case RingSlotCurrent:
		if m.ring[ringIdx] == nil {
			return nil, nil, fmt.Errorf("module %v disappeared", subID)
		}

		return m.ring[ringIdx], &events.EventList{}, nil
	case RingSlotFuture:
		if !m.ringController.TryAcquire(subID) {
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

		return m.ring[ringIdx], evsOut, err
	default:
		return nil, nil, fmt.Errorf("unknown slot status: %v", m.ringController.GetSlotStatus(subID))
	}
}

func (m *Module) applyModringEvent(event *modringpb.Event) (*events.EventList, error) {
	switch e := event.Type.(type) {
	case *modringpb.Event_Free:
		return m.applyFreeSubmodule(e.Free)
	default:
		return nil, fmt.Errorf("unsupported event type for modring module: %T", e)
	}
}

func (m *Module) applyFreeSubmodule(event *modringpb.FreeSubmodule) (*events.EventList, error) {
	if m.ringController.TryFree(event.Id) {
		// zero slot just in case (and to let the GC do its job)
		ringIdx := int(event.Id % uint64(len(m.ring)))
		m.ring[ringIdx] = nil

		m.logger.Log(logging.LevelDebug, "module freed", "id", event.Id, "ringIdx", ringIdx)
		return &events.EventList{}, nil
	} else if m.ringController.IsFutureSlot(event.Id) {
		return nil, fmt.Errorf("cannot free future slot (id=%v)", event.Id)
	} else {
		m.logger.Log(logging.LevelDebug, "tried to double-free module", "id", event.Id)
		// already freed, but that's alright
		return &events.EventList{}, nil
	}
}

func applyAllSafely(m modules.PassiveModule, evs *events.EventList) (result *events.EventList, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rErr, ok := r.(error); ok {
				err = fmt.Errorf("event application panicked: %w\nStack trace:\n%s", rErr, string(debug.Stack()))
			} else {
				err = fmt.Errorf("event application panicked: %v\nStack trace:\n%s", r, string(debug.Stack()))
			}
		}
	}()

	return m.ApplyEvents(evs)
}
