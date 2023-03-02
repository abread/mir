package modring

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type Module struct {
	ctx    context.Context
	ownID  t.ModuleID
	params *ModuleParams

	// signal worker to start processing a module
	workerProcessModule []chan *subModule

	ringController RingController
	ring           []*subModule

	outChan chan *events.EventList

	logger logging.Logger
}

type subModule struct {
	module    modules.Module
	inputChan chan *events.EventList

	tracingCtx context.Context
	span       trace.Span

	exitedChan chan struct{}
}

var name = "github.com/filecoin-project/mir/pkg/modring"

func New(ctx context.Context, ownID t.ModuleID, ringSize int, params *ModuleParams, logger logging.Logger, outChan chan *events.EventList) *Module {
	if logger == nil {
		logger = logging.ConsoleErrorLogger
	}

	m := &Module{
		ctx:    ctx,
		ownID:  ownID,
		params: params,

		workerProcessModule: make([]chan *subModule, ringSize),

		ringController: NewRingController(ringSize),
		ring:           make([]*subModule, ringSize),

		outChan: outChan,

		logger: logger,
	}

	// start workers
	for i := 0; i < ringSize; i++ {
		m.workerProcessModule[i] = make(chan *subModule)

		go m.doWork(i)
	}

	return m
}

func (m *Module) doWork(i int) {
	doneC := m.ctx.Done()
	for {
		select {
		case <-doneC:
			return
		case mod := <-m.workerProcessModule[i]:
			if activeMod, ok := mod.module.(modules.ActiveModule); ok {
				m.doWorkActiveModule(mod.tracingCtx, activeMod, mod.inputChan)
			} else {
				passiveMod := mod.module.(modules.PassiveModule)
				m.doWorkPassiveModule(passiveMod, mod.inputChan)
			}

			mod.span.End()
			close(mod.exitedChan)
		}
	}
}

func (m *Module) doWorkActiveModule(ctx context.Context, sub modules.ActiveModule, input <-chan *events.EventList) {
	doneC := m.ctx.Done()
	for {
		select {
		case <-doneC:
			return
		case evs, ok := <-input:
			if !ok {
				return
			}

			if err := sub.ApplyEvents(ctx, evs); err != nil {
				panic(fmt.Errorf("error processing events from submodule: %w", err))
			}
		}
	}
}

func (m *Module) doWorkPassiveModule(sub modules.PassiveModule, input <-chan *events.EventList) {
	doneC := m.ctx.Done()

	for {
		select {
		case <-doneC:
			return
		case evs, ok := <-input:
			if !ok {
				return // module is to be stopped
			}

			outEvs, err := sub.ApplyEvents(evs)
			if err != nil {
				panic(fmt.Errorf("error processing events from submodule: %w", err))
			}
			m.outChan <- outEvs
		}
	}
}

func (m *Module) ImplementsModule() {}

func (m *Module) EventsOut() <-chan *events.EventList {
	return m.outChan
}

func (m *Module) ApplyEvents(ctx context.Context, eventsIn *events.EventList) error {
	// Note: this method is similar to ApplyEvents, but it only parallelizes execution across different
	// modules, not all events.
	subEventsIn, pastMsgs := m.splitEventsByDest(eventsIn)

	for ringIdx := 0; ringIdx < len(m.ring); ringIdx++ {
		if subEventsIn[ringIdx].Len() == 0 {
			continue
		}

		subID := subIDInRingIdx(&m.ringController, ringIdx)
		switch m.ringController.GetSlotStatus(subID) {
		case RingSlotPast:
			return fmt.Errorf("unreachable code: past message cannot be deemed deliverable to current ring")
		case RingSlotFuture:
			if !m.ringController.MarkCurrentIfFuture(subID) {
				return fmt.Errorf("unreachable code: future slot changed state")
			}

			sub, initialEvs, err := m.createModule(subID)
			if err != nil {
				return err
			}
			m.ring[ringIdx] = sub

			// allocate worker to submodule
			m.workerProcessModule[ringIdx] <- sub

			// send initial events
			m.ring[ringIdx].inputChan <- initialEvs

			m.logger.Log(logging.LevelDebug, "module created", "id", subID)
			fallthrough
		case RingSlotCurrent:
			m.ring[ringIdx].inputChan <- &subEventsIn[ringIdx]
		default:
			return fmt.Errorf("unknown slot status")
		}
	}

	if len(pastMsgs) > 0 && m.params.PastMsgHandler != nil {
		evsOut, err := (m.params.PastMsgHandler)(pastMsgs)
		if err != nil {
			return err
		}

		m.outChan <- evsOut
	}

	return nil
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
			m.logger.Log(logging.LevelDebug, "event received for out of view submodule", "submoduleID", subID)
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

func (m *Module) createModule(id uint64) (*subModule, *events.EventList, error) {
	subFullID := m.ownID.Then(t.NewModuleIDFromInt(id))

	subCtx, subSpan := otel.Tracer(name).Start(m.ctx, fmt.Sprintf("modring:%s", subFullID))
	subMod, moreInitialEvs, err := (m.params.Generator)(subCtx, subFullID, id, m.outChan)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating submodule: %w", err)
	}

	sub := &subModule{
		module:     subMod,
		inputChan:  make(chan *events.EventList, m.params.InputQueueSize),
		tracingCtx: subCtx,
		span:       subSpan,
		exitedChan: make(chan struct{}),
	}

	initialEvs := events.ListOf(events.Init(subFullID))
	initialEvs.PushBackList(moreInitialEvs)

	return sub, initialEvs, err
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
				return fmt.Errorf("cannot advance view: %w", err)
			}
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

	if m.ringController.IsFutureSlot(id) {
		return m.ringController.MarkPast(id)
	}

	if err := m.ringController.MarkPast(id); err != nil {
		return fmt.Errorf("cannot mark submodule as past: %w", err)
	}

	ringIdx := int(id % uint64(len(m.ring)))

	close(m.ring[ringIdx].inputChan)

	// wait for worker to stop processing module
	<-m.ring[ringIdx].exitedChan

	m.ring[ringIdx] = nil

	m.logger.Log(logging.LevelDebug, "module freed", "id", id, "newMinSlot", m.ringController.minSlot)
	return nil
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
