package modring

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/modringpb"
	modringpbevents "github.com/filecoin-project/mir/pkg/pb/modringpb/events"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type Module struct {
	ownID     t.ModuleID
	generator ModuleGenerator

	lock           sync.Mutex
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

	results := make([]chan *events.EventList, len(subEventsIn))
	errors := make([]chan error, len(subEventsIn))
	for i := 0; i < eventsIn.Len(); i++ {
		results[i] = make(chan *events.EventList)
		errors[i] = make(chan error)
	}

	// Start one concurrent worker for each submodule
	// Note that the processing starts concurrently for all eventlists, and only the writing of the results is synchronized.
	for i := 0; i < len(subEventsIn); i++ {
		go func(mod modules.PassiveModule, evs *events.EventList, j int) {

			// Apply the input event, catching potential panics.
			res, err := mod.ApplyEvents(evs)

			// Write processing results to the output channels.
			// Attention: Those (unbuffered) channels must be read by the aggregator in the same order
			//            as they are being written here, otherwise the system gets stuck.
			results[j] <- res
			errors[j] <- err

		}(m.ring[i], &subEventsIn[i], i)
		i++
	}

	// The event processing results will be aggregated here.
	var firstError error
	var eventsOut *events.EventList

	// the modring's events will be handled in the current go routine
	eventsOut, firstError = modules.ApplyEventsSequentially(ownEventsIn, m.applyEvent)

	// For each submodule, read the processing result from the common channels and aggregate it with the rest.
	for i := 0; i < len(subEventsIn); i++ {

		// Attention: The (unbuffered) errors and results channels must be read in the same order
		//            as they are being written by the worker goroutines, otherwise the system gets stuck.

		// Read results from common channel and add it to the accumulator.
		if evList := <-results[i]; evList != nil {
			eventsOut.PushBackList(evList)
		}

		// Read error from common channel.
		// We only consider the first error, as ApplyEventsConcurrently only returns a single error.
		// TODO: Explore possibilities of aggregating multiple errors in one.
		if err := <-errors[i]; err != nil && firstError == nil {
			firstError = err
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
			idxStr := t.ModuleID(event.DestModule).StripParent(m.ownID).Top()
			idx, err := strconv.ParseUint(string(idxStr), 10, 64)
			if err != nil {
				m.logger.Log(logging.LevelWarn, "event received for invalid submodule index", "submoduleIdx", idxStr)
			}

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

func (m *Module) getSub(subIdx uint64) (modules.PassiveModule, *events.EventList, error) {
	ringIdx := int(subIdx % uint64(len(m.ring)))

	m.lock.Lock()
	defer m.lock.Unlock()

	switch m.ringController.GetSlotStatus(subIdx) {
	case RingSlotPast:
		return nil, nil, nil
	case RingSlotCurrent:
		return m.ring[ringIdx], &events.EventList{}, nil
	case RingSlotFuture:
		if !m.ringController.TryAcquire(subIdx) {
			return nil, nil, nil
		}

		subFullID := m.ownID.Then(t.NewModuleIDFromInt(subIdx))
		sub, err := (m.generator)(subFullID)
		if err != nil {
			return nil, nil, err
		}

		m.ring[ringIdx] = sub

		evsOut, err := m.ring[ringIdx].ApplyEvents(events.ListOf(events.Init(subFullID)))
		if err != nil {
			return nil, nil, err
		}

		return m.ring[ringIdx], evsOut, err
	default:
		return nil, nil, fmt.Errorf("unknown slot status: %v", m.ringController.GetSlotStatus(subIdx))
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
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.ringController.TryFree(event.Id) {
		// zero slot just in case (and to let the GC do its job)
		ringIdx := int(event.Id % uint64(len(m.ring)))
		m.ring[ringIdx] = nil

		return events.ListOf(
			modringpbevents.FreedSubmodule(
				t.ModuleID(event.Origin.Module),
				modringpbtypes.FreeSubmoduleOriginFromPb(event.Origin),
			).Pb(),
		), nil
	} else if m.ringController.IsSlotUnused(event.Id) {
		return nil, fmt.Errorf("cannot free future slot (id=%v)", event.Id)
	} else {
		// already freed, but that's alright
		return &events.EventList{}, nil
	}
}
