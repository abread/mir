package reliablenet

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/reliablenetpb"
	rnpbmsg "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages"
	rnEvents "github.com/filecoin-project/mir/pkg/reliablenet/events"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleConfig struct {
	Self  t.ModuleID
	Net   t.ModuleID
	Timer t.ModuleID
}

type ModuleParams struct {
	RetransmissionLoopInterval time.Duration

	AllNodes []t.NodeID
}

func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self:  "reliablenet",
		Net:   "net",
		Timer: "timer",
	}
}

func DefaultModuleParams(allNodes []t.NodeID) *ModuleParams {
	return &ModuleParams{
		RetransmissionLoopInterval: 500 * time.Millisecond,
		AllNodes:                   allNodes,
	}
}

type Module struct {
	config *ModuleConfig
	params *ModuleParams
	ownID  t.NodeID
	logger logging.Logger

	retransmitLoopScheduled uint32 // TODO: atomic.Bool in Go 1.19

	queues map[t.ModuleID]*sync.Map // ModuleID -> msgID ([]byte) -> queuedMsg
	locker sync.RWMutex
}

type queuedMsg struct {
	msg          *messagepb.Message
	destinations *sync.Map // NodeID -> struct{}
}

func New(ownID t.NodeID, mc *ModuleConfig, params *ModuleParams, logger logging.Logger) (*Module, error) {
	return &Module{
		config: mc,
		params: params,
		ownID:  ownID,
		logger: logger,

		retransmitLoopScheduled: 0,

		queues: make(map[t.ModuleID]*sync.Map),
	}, nil
}

func (m *Module) GetPendingMessages() []*eventpb.SendMessage {
	m.locker.Lock()
	defer m.locker.Unlock()

	pendingMsgs := make([]*eventpb.SendMessage, 0)
	for _, queue := range m.queues {
		queue.Range(func(key, value any) bool {
			qmsg := value.(queuedMsg)
			dests := make([]t.NodeID, 0)
			qmsg.destinations.Range(func(key, value any) bool {
				d := key.(t.NodeID)
				dests = append(dests, d)
				return true
			})

			pendingMsgs = append(pendingMsgs, &eventpb.SendMessage{
				Msg:          qmsg.msg,
				Destinations: t.NodeIDSlicePb(dests),
			})

			return true
		})
	}

	return pendingMsgs
}

func (m *Module) ApplyEvents(evs *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsConcurrently(evs, m.applyEvent)
}

func (m *Module) ImplementsModule() {}

func (m *Module) applyEvent(event *eventpb.Event) (*events.EventList, error) {
	switch ev := event.Type.(type) {
	case *eventpb.Event_Init:
		return &events.EventList{}, nil
	case *eventpb.Event_ReliableNet:
		return m.applyRNEvent(ev.ReliableNet)
	case *eventpb.Event_MessageReceived:
		return m.applyMessage(ev.MessageReceived.Msg, t.NodeID(ev.MessageReceived.From))
	default:
		return nil, fmt.Errorf("unsupported event type: %T", ev)
	}
}

func (m *Module) applyRNEvent(event *reliablenetpb.Event) (*events.EventList, error) {
	switch ev := event.Type.(type) {
	case *reliablenetpb.Event_RetransmitAll:
		return m.retransmitAll()
	case *reliablenetpb.Event_SendMessage:
		return m.SendMessage(string(ev.SendMessage.MsgId), ev.SendMessage.Msg, t.NodeIDSlice(ev.SendMessage.Destinations))
	case *reliablenetpb.Event_MarkModuleMsgsRecvd:
		return m.MarkModuleMsgsRecvd(t.ModuleID(ev.MarkModuleMsgsRecvd.DestModule), t.NodeIDSlice(ev.MarkModuleMsgsRecvd.Destinations))
	case *reliablenetpb.Event_MarkRecvd:
		return m.MarkRecvd(t.ModuleID(ev.MarkRecvd.DestModule), string(ev.MarkRecvd.MsgId), t.NodeIDSlice(ev.MarkRecvd.Destinations))
	case *reliablenetpb.Event_Ack:
		return m.SendAck(t.ModuleID(ev.Ack.DestModule), ev.Ack.MsgId, t.NodeID(ev.Ack.Source))
	default:
		return nil, fmt.Errorf("unsupported reliablenet event type: %T", ev)
	}
}

func (m *Module) applyMessage(message *messagepb.Message, from t.NodeID) (*events.EventList, error) {
	wrappedMsg, ok := message.Type.(*messagepb.Message_ReliableNet)
	if !ok {
		return &events.EventList{}, nil
	}

	msg, ok := wrappedMsg.ReliableNet.Type.(*rnpbmsg.Message_Ack)
	if !ok {
		return &events.EventList{}, nil
	}

	return m.MarkRecvd(t.ModuleID(msg.Ack.MsgDestModule), string(msg.Ack.MsgId), []t.NodeID{from})
}

func (m *Module) SendAck(msgDestModule t.ModuleID, msgID []byte, msgSource t.NodeID) (*events.EventList, error) {
	return events.ListOf(
		events.SendMessage(
			m.config.Net,
			&messagepb.Message{
				Type: &messagepb.Message_ReliableNet{
					ReliableNet: &rnpbmsg.Message{
						Type: &rnpbmsg.Message_Ack{
							Ack: &rnpbmsg.AckMessage{
								MsgDestModule: msgDestModule.Pb(),
								MsgId:         msgID,
							},
						},
					},
				},
				DestModule: m.config.Self.Pb(),
			},
			[]t.NodeID{msgSource},
		),
	), nil
}

func (m *Module) retransmitAll() (*events.EventList, error) {
	evsOut, probablyEmptyQueues, err := m.retransmitNoDeleteQueues()
	shouldScheduleLoop := true

	// clean up empty queues
	if err != nil && len(probablyEmptyQueues) > 0 {
		m.locker.Lock()
		defer m.locker.Unlock()

		for _, moduleID := range probablyEmptyQueues {
			msgCount := 0
			m.queues[moduleID].Range(func(_, _ any) bool {
				msgCount++
				return true
			})

			if msgCount == 0 {
				delete(m.queues, moduleID)
			}
		}

		if len(m.queues) == 0 {
			atomic.StoreUint32(&m.retransmitLoopScheduled, 0)
			shouldScheduleLoop = false
		}
	}

	if shouldScheduleLoop {
		evsOut.PushBack(events.TimerDelay(
			m.config.Timer,
			[]*eventpb.Event{rnEvents.RetransmitAll(m.config.Self)},
			t.TimeDuration(m.params.RetransmissionLoopInterval),
		))
	}

	return evsOut, err
}

func (m *Module) retransmitNoDeleteQueues() (*events.EventList, []t.ModuleID, error) { // TODO: name
	m.locker.RLock()
	defer m.locker.RUnlock()

	evsOut := &events.EventList{}
	probablyEmptyQueues := make([]t.ModuleID, 0)

	for moduleID, queue := range m.queues {
		msgCount := 0

		var err error
		queue.Range(func(key, value any) bool {
			qmsg, ok := value.(queuedMsg)
			if !ok {
				err = fmt.Errorf("inconsistent message queue (expected queuedMsg, found something else)")
				return false
			}

			dests := make([]t.NodeID, 0, len(m.params.AllNodes))
			qmsg.destinations.Range(func(key, _ any) bool {
				dest, ok := key.(t.NodeID)
				if !ok {
					err = fmt.Errorf("inconsistent message queue (expected t.NodeID, found something else)")
					return false
				}

				dests = append(dests, dest)
				return true
			})
			if err != nil {
				return false
			}

			if len(dests) > 0 {
				m.logger.Log(logging.LevelDebug, "Retransmitting message", "dests", dests, "msg", qmsg.msg)
				evsOut.PushBack(events.SendMessage(m.config.Net, qmsg.msg, dests))
				msgCount++
			} else {
				queue.Delete(key)
			}
			return true
		})

		if msgCount == 0 {
			probablyEmptyQueues = append(probablyEmptyQueues, moduleID)
		}
	}

	return evsOut, probablyEmptyQueues, nil
}

func (m *Module) SendMessage(id string, msg *messagepb.Message, destinations []t.NodeID) (*events.EventList, error) {
	destinationSet := &sync.Map{}
	for _, nodeID := range destinations {
		if nodeID != m.ownID {
			destinationSet.Store(nodeID, struct{}{})
		}
	}

	qmsg := queuedMsg{
		msg:          msg,
		destinations: destinationSet,
	}

	locker := m.ensureQueueExistsAndLock(t.ModuleID(msg.DestModule))
	m.queues[t.ModuleID(msg.DestModule)].Store(id, qmsg)
	shouldScheduleLoop := atomic.CompareAndSwapUint32(&m.retransmitLoopScheduled, 0, 1)
	locker.Unlock()

	evsOut := events.ListOf(
		events.SendMessage(m.config.Net, msg, destinations),
	)

	if shouldScheduleLoop {
		evsOut.PushBack(events.TimerDelay(
			m.config.Timer,
			[]*eventpb.Event{rnEvents.RetransmitAll(m.config.Self)},
			t.TimeDuration(m.params.RetransmissionLoopInterval),
		))
	}

	return evsOut, nil
}

func (m *Module) ensureQueueExistsAndLock(destModule t.ModuleID) sync.Locker {
	locker := m.locker.RLocker()
	locker.Lock()

	if _, present := m.queues[destModule]; present {
		return locker
	}

	// does not exist, reacquire lock with write privileges
	locker.Unlock()
	locker = &m.locker
	locker.Lock()

	// we dropped the lock, the queue may have been created in the meantime
	if _, present := m.queues[destModule]; present {
		return locker
	}

	// must create queue
	m.queues[destModule] = &sync.Map{}

	// we could try to downgrade the lock to a read lock, but another routine could delete the queue
	// so it could lead to a livelock situation
	// the critical section will be small anyway
	return locker
}

func (m *Module) MarkModuleMsgsRecvd(destModule t.ModuleID, destinations []t.NodeID) (*events.EventList, error) {
	// optimization for all destinations
	if len(destinations) == len(m.params.AllNodes) {
		m.locker.Lock()
		defer m.locker.Unlock()

		delete(m.queues, destModule)
		return &events.EventList{}, nil
	}

	m.locker.RLock()
	defer m.locker.RUnlock()

	queue, queueExists := m.queues[destModule]
	if !queueExists {
		// queue doesn't event exist anymore, no pending messages to clean up
		return &events.EventList{}, nil
	}

	var err error
	queue.Range(func(key, value any) bool {
		qmsg, ok := value.(queuedMsg)
		if !ok {
			err = fmt.Errorf("inconsistent message queue (expected queuedMsg, found something else)")
			return false
		}

		for _, nodeID := range destinations {
			qmsg.destinations.Delete(nodeID)
		}

		return true
	})

	return &events.EventList{}, err
}

func (m *Module) MarkRecvd(destModule t.ModuleID, id string, destinations []t.NodeID) (*events.EventList, error) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	queue, queueExists := m.queues[destModule]
	if !queueExists {
		// queue doesn't event exist anymore, no pending messages to clean up
		return &events.EventList{}, nil
	}

	if val, present := queue.Load(id); present {
		qmsg, ok := val.(queuedMsg)
		if !ok {
			return nil, fmt.Errorf("inconsistent message queue (expected queueMsg, found something else)")
		}

		for _, d := range destinations {
			qmsg.destinations.Delete(d)
		}
	}

	return &events.EventList{}, nil
}
