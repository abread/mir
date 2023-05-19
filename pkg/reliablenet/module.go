package reliablenet

import (
	"fmt"
	"time"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	eventpbevents "github.com/filecoin-project/mir/pkg/pb/eventpb/events"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	"github.com/filecoin-project/mir/pkg/pb/reliablenetpb"
	rnEvents "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/events"
	rnpbmsg "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages"
	rnetmsgs "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages/msgs"
	"github.com/filecoin-project/mir/pkg/pb/transportpb"
	transportpbevents "github.com/filecoin-project/mir/pkg/pb/transportpb/events"
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	timert "github.com/filecoin-project/mir/pkg/timer/types"
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

func DefaultModuleParams(allNodes []t.NodeID) *ModuleParams {
	return &ModuleParams{
		RetransmissionLoopInterval: 5000 * time.Millisecond,
		MaxRetransmissionBurst:     64, // half the default for the net transport
		AllNodes:                   allNodes,
	}
}

type QueuedMessage struct {
	Msg          *messagepb.Message
	Destinations []t.NodeID
	Ts           uint64
}

type queuesMap map[t.ModuleID]map[rntypes.MsgID]QueuedMessage

type Module struct {
	config ModuleConfig
	params *ModuleParams
	ownID  t.NodeID
	logger logging.Logger

	retransmitLoopScheduled bool
	iterationCount          uint64

	fresherQueues queuesMap
	stalerQueues  queuesMap
}

func New(ownID t.NodeID, mc ModuleConfig, params *ModuleParams, logger logging.Logger) (*Module, error) {
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

func (m *Module) GetPendingMessages() []QueuedMessage {
	pendingMsgs := make([]QueuedMessage, 0, len(m.fresherQueues)+len(m.stalerQueues))
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
	case *eventpb.Event_Transport:
		switch e := ev.Transport.Type.(type) {
		case *transportpb.Event_MessageReceived:
			return m.applyMessage(e.MessageReceived.Msg, t.NodeID(e.MessageReceived.From))
		default:
			return nil, fmt.Errorf("unsupported transport event type: %T", e)
		}
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
		transportpbevents.SendMessage(
			m.config.Net,
			rnetmsgs.AckMessage(m.config.Self, msgDestModule, msgID),
			[]t.NodeID{msgSource},
		).Pb(),
	), nil
}

func (m *Module) retransmitAll() (*events.EventList, error) {
	evsOut := events.EmptyList()

	if len(m.stalerQueues) == 0 {
		// try to switch queues
		m.fresherQueues, m.stalerQueues = m.stalerQueues, m.fresherQueues
	}

	if len(m.stalerQueues) == 0 {
		// ok we're really out of messages
		m.retransmitLoopScheduled = false
		// m.logger.Log(logging.LevelDebug, "stopping retransmission loop: no more messages")
		return evsOut, nil
	}

	nRemaining := m.params.MaxRetransmissionBurst

	for moduleID, queue := range m.stalerQueues {
		if len(queue) == 0 {
			return evsOut, fmt.Errorf("queue for module %v has no messages", moduleID)
		}

		if err := m.retransmitQueue(moduleID, evsOut, &nRemaining); err != nil {
			return nil, err
		}

		if nRemaining == 0 {
			break // by definition, we cannot process more messages
		}
	}

	evsOut.PushBack(eventpbevents.TimerDelay(
		m.config.Timer,
		[]*eventpbtypes.Event{rnEvents.RetransmitAll(m.config.Self)},
		timert.Duration(m.params.RetransmissionLoopInterval),
	).Pb())

	m.iterationCount++
	return evsOut, nil
}

func (m *Module) retransmitQueue(moduleID t.ModuleID, evsOut *events.EventList, nRemaining *int) error {
	queue := m.stalerQueues[moduleID]

	for msgID, qmsg := range queue {
		if *nRemaining == 0 {
			return nil
		}
		if qmsg.Ts <= m.iterationCount {
			continue
		}

		if len(qmsg.Destinations) > 0 {
			// m.logger.Log(logging.LevelDebug, "Retransmitting message", "msg", qmsg)
			evsOut.PushBack(
				transportpbevents.SendMessage(m.config.Net, messagepbtypes.MessageFromPb(qmsg.Msg), qmsg.Destinations).Pb(),
			)
		} else {
			return fmt.Errorf("queued message for module %v has no destinations", moduleID)
		}

		m.ensureQueueExists(m.fresherQueues, moduleID)
		m.fresherQueues[moduleID][msgID] = qmsg
		delete(queue, msgID)

		*nRemaining--
	}

	if len(queue) == 0 {
		delete(m.stalerQueues, moduleID)
	}

	return nil
}

func (m *Module) SendMessage(id rntypes.MsgID, msg *messagepb.Message, destinations []t.NodeID) (*events.EventList, error) {
	m.ensureQueueExists(m.fresherQueues, t.ModuleID(msg.DestModule))
	m.fresherQueues[t.ModuleID(msg.DestModule)][id] = QueuedMessage{
		Msg:          msg,
		Destinations: destinations,
		Ts:           m.iterationCount,
	}

	evsOut := events.ListOf(
		transportpbevents.SendMessage(m.config.Net, messagepbtypes.MessageFromPb(msg), destinations).Pb(),
	)

	if !m.retransmitLoopScheduled {
		// m.logger.Log(logging.LevelDebug, "rescheduling retransmission loop")
		evsOut.PushBack(eventpbevents.TimerDelay(
			m.config.Timer,
			[]*eventpbtypes.Event{rnEvents.RetransmitAll(m.config.Self)},
			timert.Duration(m.params.RetransmissionLoopInterval),
		).Pb())
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
		qm[destModule] = make(map[rntypes.MsgID]QueuedMessage, initialQueueCapacity)
	}
}

func (m *Module) MarkModuleMsgsRecvd(destModule t.ModuleID, destinations []t.NodeID) (*events.EventList, error) {
	// optimization for all destinations
	if len(destinations) == len(m.params.AllNodes) {
		for _, queues := range []queuesMap{m.fresherQueues, m.stalerQueues} {
			for id := range queues {
				if id.IsSubOf(destModule) || id == destModule {
					delete(m.fresherQueues, id)
				}
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
				qmsg.Destinations = sliceutil.Filter(qmsg.Destinations, func(_ int, d t.NodeID) bool {
					return !slices.Contains(destinations, d)
				})
				queue[msgID] = qmsg

				if len(qmsg.Destinations) == 0 {
					delete(queue, msgID)
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
			continue
		}

		if qmsg, present := queue[messageID]; present {
			qmsg.Destinations = sliceutil.Filter(qmsg.Destinations, func(_ int, d t.NodeID) bool {
				return !slices.Contains(destinations, d)
			})
			queue[messageID] = qmsg

			if len(qmsg.Destinations) == 0 {
				delete(queue, messageID)
			}
		}

		if len(queue) == 0 {
			delete(queues, destModule)
		}
	}

	return &events.EventList{}, nil
}
