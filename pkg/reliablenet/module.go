package reliablenet

import (
	"fmt"
	"time"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/reliablenetpb"
	rnEvents "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/events"
	rnpbmsg "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages"
	rnetmsgs "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages/msgs"
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

type ModuleConfig struct {
	Self  t.ModuleID
	Net   t.ModuleID
	Timer t.ModuleID
}

type ModuleParams struct {
	RetransmissionLoopInterval time.Duration
	MaxRetransmissionBurst     int

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
		RetransmissionLoopInterval: 2 * time.Second,
		MaxRetransmissionBurst:     64,
		AllNodes:                   allNodes,
	}
}

type queuesMap map[t.ModuleID]map[rntypes.MsgID]*eventpb.SendMessage

type Module struct {
	config *ModuleConfig
	params *ModuleParams
	ownID  t.NodeID
	logger logging.Logger

	retransmitLoopScheduled bool

	fresherQueues queuesMap
	stalerQueues  queuesMap
}

func New(ownID t.NodeID, mc *ModuleConfig, params *ModuleParams, logger logging.Logger) (*Module, error) {
	return &Module{
		config: mc,
		params: params,
		ownID:  ownID,
		logger: logger,

		retransmitLoopScheduled: false,

		fresherQueues: make(queuesMap),
		stalerQueues:  make(queuesMap),
	}, nil
}

func (m *Module) GetPendingMessages() []*eventpb.SendMessage {
	pendingMsgs := make([]*eventpb.SendMessage, 0, len(m.fresherQueues)+len(m.stalerQueues))
	for _, queue := range m.fresherQueues {
		for _, qmsg := range queue {
			pendingMsgs = append(pendingMsgs, qmsg)
		}
	}
	for _, queue := range m.stalerQueues {
		for _, qmsg := range queue {
			pendingMsgs = append(pendingMsgs, qmsg)
		}
	}

	return pendingMsgs
}

func (m *Module) CountPendingMessages() int {
	count := 0

	for _, queue := range m.fresherQueues {
		count += len(queue)
	}
	for _, queue := range m.stalerQueues {
		count += len(queue)
	}

	return count
}

func (m *Module) ApplyEvents(evs *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsSequentially(evs, m.applyEvent)
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
		return m.SendMessage(rntypes.MsgID(ev.SendMessage.MsgId), ev.SendMessage.Msg, t.NodeIDSlice(ev.SendMessage.Destinations))
	case *reliablenetpb.Event_MarkModuleMsgsRecvd:
		return m.MarkModuleMsgsRecvd(t.ModuleID(ev.MarkModuleMsgsRecvd.DestModule), t.NodeIDSlice(ev.MarkModuleMsgsRecvd.Destinations))
	case *reliablenetpb.Event_MarkRecvd:
		return m.MarkRecvd(t.ModuleID(ev.MarkRecvd.DestModule), rntypes.MsgID(ev.MarkRecvd.MsgId), t.NodeIDSlice(ev.MarkRecvd.Destinations))
	case *reliablenetpb.Event_Ack:
		return m.SendAck(t.ModuleID(ev.Ack.DestModule), rntypes.MsgID(ev.Ack.MsgId), t.NodeID(ev.Ack.Source))
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

	return m.MarkRecvd(t.ModuleID(msg.Ack.MsgDestModule), rntypes.MsgID(msg.Ack.MsgId), []t.NodeID{from})
}

func (m *Module) SendAck(msgDestModule t.ModuleID, msgID rntypes.MsgID, msgSource t.NodeID) (*events.EventList, error) {
	if msgSource == m.ownID {
		// fast path for local messages
		return m.MarkRecvd(msgDestModule, msgID, []t.NodeID{msgSource})
	}

	return events.ListOf(
		events.SendMessage(
			m.config.Net,
			rnetmsgs.AckMessage(m.config.Self, msgDestModule, msgID).Pb(),
			[]t.NodeID{msgSource},
		),
	), nil
}

func (m *Module) retransmitAll() (*events.EventList, error) {
	evsOut := &events.EventList{}

	if len(m.stalerQueues) == 0 {
		// try to switch queues
		m.fresherQueues, m.stalerQueues = m.stalerQueues, m.fresherQueues
	}

	if len(m.stalerQueues) == 0 {
		// ok we're really out of messages
		m.retransmitLoopScheduled = false
		m.logger.Log(logging.LevelDebug, "stopping retransmission loop: no more messages")
		return evsOut, nil
	}

	nRemaining := m.params.MaxRetransmissionBurst

	for moduleID, queue := range m.stalerQueues {
		if len(queue) == 0 {
			return evsOut, fmt.Errorf("queue for module %v has no messages", moduleID)
		} else if len(queue) <= nRemaining {
			nRemaining -= len(queue)

			if err := m.retransmitWholeQueue(moduleID, evsOut); err != nil {
				return nil, err
			}
		} else if nRemaining != m.params.MaxRetransmissionBurst {
			continue // skip this one for now, we'll do it if we have no other choice
		} else {
			if err := m.retransmitQueuePartial(moduleID, evsOut, nRemaining); err != nil {
				return nil, err
			}

			break // by definition, we cannot process more messages
		}
	}

	evsOut.PushBack(events.TimerDelay(
		m.config.Timer,
		[]*eventpb.Event{rnEvents.RetransmitAll(m.config.Self).Pb()},
		t.TimeDuration(m.params.RetransmissionLoopInterval),
	))

	return evsOut, nil
}

func (m *Module) retransmitQueuePartial(moduleID t.ModuleID, evsOut *events.EventList, nRemaining int) error {
	queue := m.stalerQueues[moduleID]
	m.ensureQueueExists(m.fresherQueues, moduleID)
	fresherQueue := m.fresherQueues[moduleID]

	for msgID, qmsg := range queue {
		if nRemaining == 0 {
			return nil
		}

		if len(qmsg.Destinations) > 0 {
			m.logger.Log(logging.LevelDebug, "Retransmitting message", "msg", qmsg)
			evsOut.PushBack(
				&eventpb.Event{
					Type: &eventpb.Event_SendMessage{
						SendMessage: qmsg,
					},
					DestModule: m.config.Net.Pb(),
				})
		} else {
			return fmt.Errorf("queued message for module %v has no destinations", moduleID)
		}

		fresherQueue[msgID] = qmsg
		delete(queue, msgID)

		nRemaining--
	}

	return nil
}

func (m *Module) retransmitWholeQueue(moduleID t.ModuleID, evsOut *events.EventList) error {
	queue := m.stalerQueues[moduleID]

	for _, qmsg := range queue {
		if len(qmsg.Destinations) > 0 {
			m.logger.Log(logging.LevelDebug, "Retransmitting message", "msg", qmsg)
			evsOut.PushBack(
				&eventpb.Event{
					Type: &eventpb.Event_SendMessage{
						SendMessage: qmsg,
					},
					DestModule: m.config.Net.Pb(),
				})

		} else {
			return fmt.Errorf("queued message for module %v has no destinations", moduleID)
		}
	}

	m.fresherQueues[moduleID] = queue
	delete(m.stalerQueues, moduleID)

	return nil
}

func (m *Module) SendMessage(id rntypes.MsgID, msg *messagepb.Message, destinations []t.NodeID) (*events.EventList, error) {
	m.ensureQueueExists(m.fresherQueues, t.ModuleID(msg.DestModule))
	m.fresherQueues[t.ModuleID(msg.DestModule)][id] = &eventpb.SendMessage{
		Msg:          msg,
		Destinations: t.NodeIDSlicePb(destinations),
	}

	evsOut := events.ListOf(
		events.SendMessage(m.config.Net, msg, destinations),
	)

	if !m.retransmitLoopScheduled {
		m.logger.Log(logging.LevelDebug, "rescheduling retransmission loop")
		evsOut.PushBack(events.TimerDelay(
			m.config.Timer,
			[]*eventpb.Event{rnEvents.RetransmitAll(m.config.Self).Pb()},
			t.TimeDuration(m.params.RetransmissionLoopInterval),
		))
		m.retransmitLoopScheduled = true
	}

	return evsOut, nil
}

// Initial capacity for message queue associated to a module
// The queue is created lazily, so it will at least have one message added to it
// In a ideal scenario, are replicas are more or less up-to-date so we should only need to keep
// a few messages around.
const initialQueueCapacity = 4

func (m *Module) ensureQueueExists(qm queuesMap, destModule t.ModuleID) {
	if _, present := qm[destModule]; !present {
		qm[destModule] = make(map[rntypes.MsgID]*eventpb.SendMessage, initialQueueCapacity)
	}
}

func (m *Module) MarkModuleMsgsRecvd(destModule t.ModuleID, destinations []t.NodeID) (*events.EventList, error) {
	// optimization for all destinations
	if len(destinations) == len(m.params.AllNodes) {
		for id := range m.fresherQueues {
			if id.IsSubOf(destModule) || id == destModule {
				delete(m.fresherQueues, id)
			}
		}

		return &events.EventList{}, nil
	}

	for _, queues := range []queuesMap{m.fresherQueues, m.stalerQueues} {
		for queueID, queue := range queues {
			if !queueID.IsSubOf(destModule) && queueID != destModule {
				continue
			}

			for msgID, qmsg := range queue {
				newDests := sliceutil.Filter(qmsg.Destinations, func(_ int, d string) bool {
					return !slices.Contains(destinations, t.NodeID(d))
				})

				if len(newDests) == 0 {
					delete(queue, msgID)
				} else {
					// don't modify the underlying message to avoid data races
					queue[msgID] = &eventpb.SendMessage{
						Msg:          qmsg.Msg,
						Destinations: newDests,
					}
				}
			}

			if len(queue) == 0 {
				delete(queues, queueID)
			}
		}
	}

	return &events.EventList{}, nil
}

func (m *Module) MarkRecvd(destModule t.ModuleID, messageID rntypes.MsgID, destinations []t.NodeID) (*events.EventList, error) {
	for _, queues := range []queuesMap{m.fresherQueues, m.stalerQueues} {
		queue, queueExists := queues[destModule]
		if !queueExists {
			// queue doesn't event exist anymore, no pending messages to clean up
			return &events.EventList{}, nil
		}

		if qmsg, present := queue[messageID]; present {
			newDests := sliceutil.Filter(qmsg.Destinations, func(_ int, d string) bool {
				return !slices.Contains(destinations, t.NodeID(d))
			})

			if len(newDests) == 0 {
				delete(queue, messageID)
			} else {
				// don't modify the underlying message to avoid data races
				queue[messageID] = &eventpb.SendMessage{
					Msg:          qmsg.Msg,
					Destinations: newDests,
				}
			}
		}

		if len(queue) == 0 {
			delete(queues, destModule)
		}
	}

	return &events.EventList{}, nil
}
