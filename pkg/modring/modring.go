package modring

import (
	"math"
	"runtime/debug"
	"strconv"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	"github.com/filecoin-project/mir/pkg/pb/transportpb"
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
	eventsOut := &events.EventList{}
	var firstError error

	// Start one concurrent worker for each submodule
	// Note that the processing starts concurrently for all eventlists, and only the writing of the results is synchronized.
	nSubs := 0
	for i := 0; i < len(subEventsIn); i++ {
		if subEventsIn[i].Len() == 0 {
			continue
		}

		mod, initEvent, err := m.getSubByRingIdx(i)
		if err != nil {
			firstError = err
			break
		}
		if mod == nil {
			continue // module out of view
		}

		if initEvent != nil {
			// created module, loop events in the system again

			// TODO: INIT event may be processed *after* other events directed at this module
			// possible solution: track which modules are pending initialization, and loop events back in
			// until INIT is processed
			initEvent.Next = append(initEvent.Next, subEventsIn[i].Slice()...)
			eventsOut.PushBack(initEvent)

			continue
		}

		go func(modring *Module, evsIn *events.EventList, j int) {
			// Apply the input event
			res, err := multiApplySafely(mod, evsIn)

			if err == nil {
				resultChan <- res
			} else {
				errorChan <- es.Errorf("failed to process submodule #%d events: %w", j, err)
			}
		}(m, &subEventsIn[i], i)
		nSubs++
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
			// m.logger.Log(logging.LevelWarn, "event received for invalid submodule index", "submoduleIdx", subIDStr)
			continue
		}

		if m.ringController.IsPastSlot(subID) {
			if ev, ok := event.Type.(*eventpb.Event_Transport); ok {
				if e, ok := ev.Transport.Type.(*transportpb.Event_MessageReceived); ok {
					pastMsgs = append(pastMsgs, &modringpbtypes.PastMessage{
						DestId:  subID,
						From:    t.NodeID(e.MessageReceived.From),
						Message: messagepbtypes.MessageFromPb(e.MessageReceived.Msg),
					})
				}
			}
			continue
		}

		if !m.ringController.IsSlotInView(subID) {
			// m.logger.Log(logging.LevelDebug, "event received for out of view submodule", "submoduleID", subID)
			continue
		}

		idx := int(subID % uint64(len(m.ring)))
		subEventsIn[idx].PushBack(event)
	}

	return subEventsIn, pastMsgs
}

func subIDInRingIdx(rc *RingController, ringIdx int) uint64 {
	minSlot, endSlot := rc.GetCurrentView()
	ringSize := endSlot - minSlot + 1
	minIdx := int(minSlot % ringSize)

	if ringIdx >= minIdx {
		return minSlot + uint64(ringIdx-minIdx)
	}

	return minSlot + uint64(int(ringSize)-minIdx) + uint64(ringIdx)
}

func (m *Module) getSubByRingIdx(ringIdx int) (modules.PassiveModule, *eventpb.Event, error) {
	subID := subIDInRingIdx(&m.ringController, ringIdx)

	switch m.ringController.GetSlotStatus(subID) {
	case RingSlotPast:
		return nil, nil, nil
	case RingSlotCurrent:
		if m.ring[ringIdx] == nil {
			return nil, nil, es.Errorf("module %v disappeared", subID)
		}

		return m.ring[ringIdx], nil, nil
	case RingSlotFuture:
		if !m.ringController.MarkCurrentIfFuture(subID) {
			return nil, nil, nil
		}

		subFullID := m.ownID.Then(t.NewModuleIDFromInt(subID))
		sub, initialEvs, err := (m.generator)(subFullID, subID)
		if err != nil {
			return nil, nil, err
		}

		// construct init event
		initEvent := events.Init(subFullID)
		initEvent.Next = append(initEvent.Next, initialEvs.Slice()...)

		m.ring[ringIdx] = sub
		// m.logger.Log(logging.LevelDebug, "module created", "id", subID)

		return m.ring[ringIdx], initEvent, err
	default:
		return nil, nil, es.Errorf("unknown slot status: %v", m.ringController.GetSlotStatus(subID))
	}
}

func (m *Module) AdvanceViewToAtLeastSubmodule(id uint64) error {
	if m.ringController.IsPastSlot(id) || m.ringController.IsSlotInView(id) {
		return nil // already advanced enough
	}

	// free all current slots
	minSlot, endSlot := m.ringController.GetCurrentView()
	for slot := minSlot; slot <= endSlot; slot++ {
		if m.ringController.IsCurrentSlot(slot) {
			if err := m.MarkSubmodulePast(slot); err != nil {
				return es.Errorf("cannot advance view: %w", err)
			}
		}
	}

	if err := m.ringController.AdvanceViewToSlot(id); err != nil {
		return es.Errorf("cannot advance view: %w", err)
	}

	// m.logger.Log(logging.LevelDebug, "fast-forwarded view", "id", id, "newMinSlot", m.ringController.minSlot)

	return nil
}

func (m *Module) MarkSubmodulePast(id uint64) error {
	if m.ringController.IsPastSlot(id) {
		// already in the past
		return nil
	}

	if m.ringController.IsFutureSlot(id) {
		return m.ringController.MarkPast(id)
	}

	if err := m.ringController.MarkPast(id); err != nil {
		return es.Errorf("cannot mark submodule as past: %w", err)
	}

	// zero slot just in case (and to let the GC do its job)
	ringIdx := int(id % uint64(len(m.ring)))
	m.ring[ringIdx] = nil

	m.ring[ringIdx] = nil

	// m.logger.Log(logging.LevelDebug, "module freed", "id", id, "newMinSlot", m.ringController.minSlot)
	return nil
}

// multiApplySafely is a wrapper around an event processing function that catches its panic and returns it as an error.
func multiApplySafely(
	module modules.PassiveModule,
	events *events.EventList,
) (result *events.EventList, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rErr, ok := r.(error); ok {
				err = es.Errorf("event application panicked: %w\nStack trace:\n%s", rErr, string(debug.Stack()))
			} else {
				err = es.Errorf("event application panicked: %v\nStack trace:\n%s", r, string(debug.Stack()))
			}
		}
	}()

	return module.ApplyEvents(events)
}

func (m *Module) MarkAllPast() error {
	if err := m.AdvanceViewToAtLeastSubmodule(math.MaxUint64); err != nil {
		return err
	}

	if err := m.MarkSubmodulePast(math.MaxUint64); err != nil {
		return err
	}
	return m.MarkSubmodulePast(math.MaxUint64)
}
