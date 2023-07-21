package modring

import (
	"fmt"
	"math"
	"runtime/debug"
	"strconv"

	es "github.com/go-errors/errors"
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	transportpbtypes "github.com/filecoin-project/mir/pkg/pb/transportpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

type Module struct {
	ownID          t.ModuleID
	generator      ModuleGenerator
	pastMsgHandler PastMessagesHandler

	minID uint64
	ring  []submodule

	logger logging.Logger
}

type submodule struct {
	mod       modules.PassiveModule
	subID     uint64
	isCurrent bool
}

func New(ownID t.ModuleID, ringSize int, params ModuleParams, logger logging.Logger) *Module {
	if logger == nil {
		logger = logging.ConsoleErrorLogger
	}

	m := &Module{
		ownID:          ownID,
		generator:      params.Generator,
		pastMsgHandler: params.PastMsgHandler,

		minID: 0,
		ring:  make([]submodule, ringSize),

		logger: logger,
	}

	for i := range m.ring {
		m.ring[i].isCurrent = true
		m.ring[i].subID = uint64(i)
	}

	return m
}

func (m *Module) ImplementsModule() {}

func (m *Module) ApplyEvents(eventsIn events.EventList) (events.EventList, error) {
	// Note: this method is similar to ApplyEvents, but it only parallelizes execution across different
	// modules, not all events.
	subEventsByDest := m.splitEventsByDest(eventsIn)
	subEventsIn, pastMsgs, err := m.updateSubsAndFilterIncomingEvents(subEventsByDest)
	if err != nil {
		return events.EmptyList(), err
	}

	resultChan := make(chan events.EventList)
	errorChan := make(chan error)

	// The event processing results will be aggregated here
	eventsOut := events.EventList{}
	var firstError error

	// Start one concurrent worker for each submodule
	// Note that the processing starts concurrently for all eventlists, and only the writing of the results is synchronized.
	nSubs := 0
	for i := 0; i < len(subEventsIn); i++ {
		if subEventsIn[i].Len() == 0 {
			continue
		}

		mod, initEvs, err := m.getSubByRingIdx(i)
		if err != nil {
			firstError = err
			break
		}
		if mod == nil {
			return events.EmptyList(), fmt.Errorf("consistency error: ring idx %d is broken", i)
		}

		if initEvs.Len() > 0 {
			// created module, extract init event, and prepend it to the input event list
			oldEvsIn := subEventsIn[i]

			subEventsIn[i] = events.EmptyListWithCapacity(1 + oldEvsIn.Len())
			subEventsIn[i].PushBackList(initEvs.Head(1))
			subEventsIn[i].PushBackList(oldEvsIn)

			// queue extra init events for processing
			initEvs.RemoveFront(1) // remove init event that will be processed directly
			eventsOut.PushBackList(initEvs)
		}

		go func(modring *Module, evsIn events.EventList, j int) {
			// Apply the input event
			res, err := multiApplySafely(mod, evsIn)

			if err == nil {
				resultChan <- res
			} else {
				errorChan <- es.Errorf("failed to process submodule #%d events: %w", j, err)
			}
		}(m, subEventsIn[i], i)
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
		return events.EmptyList(), firstError
	}

	if len(pastMsgs) > 0 && m.pastMsgHandler != nil {
		evsOut, err := (m.pastMsgHandler)(pastMsgs)
		if err != nil {
			return events.EmptyList(), err
		}

		eventsOut.PushBackList(evsOut)
	}

	return eventsOut, nil
}

func (m *Module) splitEventsByDest(eventsIn events.EventList) map[uint64]*events.EventList {
	res := make(map[uint64]*events.EventList, len(m.ring))

	eventsInIter := eventsIn.Iterator()
	for event := eventsInIter.Next(); event != nil; event = eventsInIter.Next() {
		subIDStr := event.DestModule.StripParent(m.ownID).Top()
		subID, err := strconv.ParseUint(string(subIDStr), 10, 64)
		if err != nil {
			// m.logger.Log(logging.LevelWarn, "event received for invalid submodule id", "submoduleID", subIDStr)
			continue
		}

		if subID == math.MaxUint64 {
			// m.logger.Log(logging.LevelWarn, "event received for end-of-times submodule id", "submoduleID", subID)
			// TODO: report byz?
			continue
		}

		if subID > m.minID+uint64(len(m.ring))-1 {
			// m.logger.Log(logging.LevelDebug, "event received for out-of-view (future) submodule", "submoduleID", subID)
			continue
		}

		evList, ok := res[subID]
		if !ok {
			evList = &events.EventList{}
			res[subID] = evList
		}
		evList.PushBack(event)
	}

	return res
}

func (m *Module) updateSubsAndFilterIncomingEvents(eventsByDest map[uint64]*events.EventList) ([]events.EventList, []*modringpbtypes.PastMessage, error) {
	evsIn := make([]events.EventList, len(m.ring))
	var pastMsgs []*modringpbtypes.PastMessage

	// go through incoming event destinations in descending order
	// this way, we will not setup a submodule only to garbage-collect it in a later iteration
	incomingEvsSubIDs := maputil.GetKeys(eventsByDest)
	slices.SortFunc(incomingEvsSubIDs, func(a, b uint64) bool { return b < a })
	for _, subID := range incomingEvsSubIDs {
		idx := subID % uint64(len(m.ring))

		if m.ring[idx].subID > subID || (m.ring[idx].subID == subID && m.ring[idx].mod == nil && !m.ring[idx].isCurrent) {
			// events are in the past, forward received messages to past message handler
			it := eventsByDest[subID].Iterator()
			for event := it.Next(); event != nil; event = it.Next() {
				if ev, ok := event.Type.(*eventpbtypes.Event_Transport); ok {
					if e, ok := ev.Transport.Type.(*transportpbtypes.Event_MessageReceived); ok {
						pastMsgs = append(pastMsgs, &modringpbtypes.PastMessage{
							DestId:  subID,
							From:    e.MessageReceived.From,
							Message: e.MessageReceived.Msg,
						})
					}
				}
			}

			continue
		}
		// no future events reach this code, so all events here are current

		if m.ring[idx].subID < subID {
			// event is current, but must garbage collect old submodule first
			oldID := m.ring[idx].subID
			oldEvs := evsIn[idx]

			// garbage-collect old module and advance view
			m.ring[idx] = submodule{
				mod:       nil,
				subID:     subID,
				isCurrent: true,
			}
			evsIn[idx] = events.EmptyList()
			// m.logger.Log(logging.LevelDebug, "freed submodule", "submoduleID", oldID)

			// must flush previously buffered events to past message handler
			it := oldEvs.Iterator()
			for event := it.Next(); event != nil; event = it.Next() {
				if ev, ok := event.Type.(*eventpbtypes.Event_Transport); ok {
					if e, ok := ev.Transport.Type.(*transportpbtypes.Event_MessageReceived); ok {
						pastMsgs = append(pastMsgs, &modringpbtypes.PastMessage{
							DestId:  oldID,
							From:    e.MessageReceived.From,
							Message: e.MessageReceived.Msg,
						})
					}
				}
			}
		}

		if evsIn[idx].Len() > 0 {
			return nil, nil, es.Errorf("inconsistency in updateSubsAndFilterIncomingEvents: tried to assign events to the same ring position twice")
		}

		evsIn[idx] = *eventsByDest[subID]
	}

	return evsIn, pastMsgs, nil
}

func (m *Module) getSubByRingIdx(ringIdx int) (modules.PassiveModule, events.EventList, error) {
	if m.ring[ringIdx].mod != nil {
		return m.ring[ringIdx].mod, events.EmptyList(), nil
	}

	subID := m.ring[ringIdx].subID

	subFullID := m.ownID.Then(t.NewModuleIDFromInt(subID))
	sub, extraInitEvs, err := (m.generator)(subFullID, subID)
	if err != nil {
		return nil, events.EmptyList(), err
	}

	initEvs := events.EmptyListWithCapacity(extraInitEvs.Len())
	initEvs.PushBack(events.Init(subFullID))
	initEvs.PushBackList(extraInitEvs)

	m.ring[ringIdx].mod = sub

	// m.logger.Log(logging.LevelDebug, "initialized submodule", "submoduleID", subID)
	return m.ring[ringIdx].mod, initEvs, err
}

func (m *Module) AdvanceViewToAtLeastSubmodule(id uint64) error {
	if m.minID >= id || m.minID == math.MaxUint64 {
		return nil
	}

	m.minID = id

	for i := 0; i < len(m.ring); i++ {
		if m.ring[i].subID < id {
			m.ring[i].isCurrent = false
		}
	}

	// m.logger.Log(logging.LevelDebug, "fast-forwarded view", "newMinid", id)
	return nil
}

func (m *Module) MarkSubmodulePast(id uint64) error {
	if id < m.minID {
		return nil // already marked as past
	} else if id >= m.minID+uint64(len(m.ring)) {
		return es.Errorf("tried to mark far-in-the future submodule (%d) as past (minId is %d)", id, m.minID)
	} else if id == math.MaxUint64 {
		return es.Errorf("tried to mark the unreachable marker slot as past")
	}

	idx := id % uint64(len(m.ring))
	if m.ring[idx].subID == id {
		m.ring[idx].isCurrent = false
	} else {
		// garbage collect really old module
		// m.logger.Log(logging.LevelDebug, "freed submodule", "id", m.ring[idx].subID)

		m.ring[idx] = submodule{
			mod:       nil,
			subID:     id,
			isCurrent: false,
		}
	}

	// update minID
	minIdx := m.minID % uint64(len(m.ring))
	for m.ring[minIdx].subID == m.minID && !m.ring[minIdx].isCurrent {
		m.minID++
		minIdx = m.minID % uint64(len(m.ring))
	}

	// m.logger.Log(logging.LevelDebug, "module marked as past", "id", id, "newMinID", m.minID)
	return nil
}

func (m *Module) IsInView(id uint64) bool {
	if id == math.MaxUint64 {
		return false
	}

	return id >= m.minID && id < m.minID+uint64(len(m.ring))
}

func (m *Module) IsInExtendedView(id uint64) bool {
	idx := id % uint64(len(m.ring))
	return m.ring[idx].subID == id || m.IsInView(id)
}

// multiApplySafely is a wrapper around an event processing function that catches its panic and returns it as an error.
func multiApplySafely(
	module modules.PassiveModule,
	events events.EventList,
) (result events.EventList, err error) {
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

func (m *Module) MarkPastAndFreeAll() {
	m.minID = math.MaxUint64

	for i := 0; i < len(m.ring); i++ {
		m.ring[i] = submodule{
			mod:       nil,
			subID:     math.MaxUint64,
			isCurrent: false,
		}
	}

	// m.logger.Log(logging.LevelDebug, "freed all submodules, and marked everything as past")
}
