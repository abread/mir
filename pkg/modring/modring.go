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

func (m *Module) ApplyEvents(evts *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsConcurrently(evts, m.applyEvent)
}

func (m *Module) applyEvent(event *eventpb.Event) (*events.EventList, error) {
	if t.ModuleID(event.DestModule) == m.ownID {
		switch e := event.Type.(type) {
		case *eventpb.Event_Init:
			return events.EmptyList(), nil // Nothing to do at initialization.
		case *eventpb.Event_Modring:
			return m.applyModringEvent(e.Modring)
		default:
			return nil, fmt.Errorf("unsupported event type for modring module: %T", e)
		}
	}
	return m.forwardEvent(event)
}

func (m *Module) forwardEvent(event *eventpb.Event) (*events.EventList, error) {
	subID := t.ModuleID(event.DestModule).StripParent(m.ownID).Top()

	subIdx, err := strconv.ParseUint(string(subID), 10, 64)
	if err != nil {
		m.logger.Log(logging.LevelDebug, "message for bad submodule", "subID", subID)
		return &events.EventList{}, nil
	}

	sub, evsOut, err := m.getSub(subIdx)
	if err != nil {
		return nil, err
	}

	moreEvsOut, err := sub.ApplyEvents(events.ListOf(event))
	if err != nil {
		return nil, err
	}

	return evsOut.PushBackList(moreEvsOut), nil
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
