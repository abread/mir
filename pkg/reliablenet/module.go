package reliablenet

import (
	"container/list"
	"time"

	"golang.org/x/exp/slices"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	eventpbevents "github.com/filecoin-project/mir/pkg/pb/eventpb/events"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	rnEvents "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/events"
	rnetmsgs "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages/msgs"
	rnetmsgstypes "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages/types"
	reliablenetpbtypes "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/types"
	transportpbevents "github.com/filecoin-project/mir/pkg/pb/transportpb/events"
	transportpbtypes "github.com/filecoin-project/mir/pkg/pb/transportpb/types"
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
	Msg          *messagepbtypes.Message
	Destinations []t.NodeID
	Ts           uint64 // nolint:stylecheck
}

type Module struct {
	config ModuleConfig
	params *ModuleParams
	ownID  t.NodeID
	logger logging.Logger

	retransmitLoopScheduled bool
	iterationCount          uint64

	queue          list.List
	retransmitHead *list.Element
	queueIndex     map[t.ModuleID]map[rntypes.MsgID]*list.Element
}

func New(ownID t.NodeID, mc ModuleConfig, params *ModuleParams, logger logging.Logger) (*Module, error) {
	m := &Module{
		config: mc,
		params: params,
		ownID:  ownID,
		logger: logger,

		retransmitLoopScheduled: false,

		queueIndex: make(map[t.ModuleID]map[rntypes.MsgID]*list.Element),
	}

	m.queue.Init()
	return m, nil
}

func (m *Module) GetPendingMessages() []*QueuedMessage {
	pending := make([]*QueuedMessage, 0, m.queue.Len())
	for it := m.queue.Front(); it != nil; it = it.Next() {
		qmsg := it.Value.(*QueuedMessage)
		if len(qmsg.Destinations) > 0 {
			pending = append(pending, qmsg)
		}
	}

	return pending
}

func (m *Module) CountPendingMessages() int {
	return m.queue.Len()
}

func (m *Module) ApplyEvents(evs *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsSequentially(evs, m.applyEvent)
}

func (m *Module) ImplementsModule() {}

func (m *Module) applyEvent(event *eventpbtypes.Event) (*events.EventList, error) {
	switch ev := event.Type.(type) {
	case *eventpbtypes.Event_Init:
		return &events.EventList{}, nil
	case *eventpbtypes.Event_ReliableNet:
		return m.applyRNEvent(ev.ReliableNet)
	case *eventpbtypes.Event_Transport:
		switch e := ev.Transport.Type.(type) {
		case *transportpbtypes.Event_MessageReceived:
			return m.applyMessage(e.MessageReceived.Msg, e.MessageReceived.From)
		default:
			return nil, es.Errorf("unsupported transport event type: %T", e)
		}
	default:
		return nil, es.Errorf("unsupported event type: %T", ev)
	}
}

func (m *Module) applyRNEvent(event *reliablenetpbtypes.Event) (*events.EventList, error) {
	switch ev := event.Type.(type) {
	case *reliablenetpbtypes.Event_RetransmitAll:
		return m.retransmitAll()
	case *reliablenetpbtypes.Event_SendMessage:
		return m.SendMessage(ev.SendMessage.MsgId, ev.SendMessage.Msg, ev.SendMessage.Destinations)
	case *reliablenetpbtypes.Event_MarkModuleMsgsRecvd:
		return m.MarkModuleMsgsRecvd(ev.MarkModuleMsgsRecvd.DestModule, ev.MarkModuleMsgsRecvd.Destinations)
	case *reliablenetpbtypes.Event_MarkRecvd:
		return m.MarkRecvd(ev.MarkRecvd.DestModule, ev.MarkRecvd.MsgId, ev.MarkRecvd.Destinations)
	case *reliablenetpbtypes.Event_Ack:
		return m.SendAck(ev.Ack.DestModule, ev.Ack.MsgId, ev.Ack.Source)
	default:
		return nil, es.Errorf("unsupported reliablenet event type: %T", ev)
	}
}

func (m *Module) applyMessage(message *messagepbtypes.Message, from t.NodeID) (*events.EventList, error) {
	wrappedMsg, ok := message.Type.(*messagepbtypes.Message_ReliableNet)
	if !ok {
		return &events.EventList{}, nil
	}

	msg, ok := wrappedMsg.ReliableNet.Type.(*rnetmsgstypes.Message_Ack)
	if !ok {
		return &events.EventList{}, nil
	}

	return m.MarkRecvd(msg.Ack.MsgDestModule, msg.Ack.MsgId, []t.NodeID{from})
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
		),
	), nil
}

func (m *Module) retransmitAll() (*events.EventList, error) {
	evsOut := events.EmptyList()

	if m.queue.Len() == 0 {
		m.retransmitLoopScheduled = false
		return evsOut, nil
	}

	nRemaining := m.params.MaxRetransmissionBurst
	if nRemaining > m.queue.Len() {
		nRemaining = m.queue.Len()
	}

	restarted := false

	for nRemaining > 0 && !restarted {
		if m.retransmitHead == nil || m.retransmitHead.Value.(*QueuedMessage).Ts < m.iterationCount {
			// restart queue iteration when end is reached/messages are too recent
			m.retransmitHead = m.queue.Front()
			restarted = true
		}

		for ; m.retransmitHead != nil && nRemaining > 0; m.retransmitHead = m.retransmitHead.Next() {
			qmsg := m.retransmitHead.Value.(*QueuedMessage)

			if qmsg.Ts < m.iterationCount {
				break
			}

			evsOut.PushBack(
				transportpbevents.SendMessage(m.config.Net, qmsg.Msg, qmsg.Destinations),
			)
			nRemaining--
		}
	}

	evsOut.PushBack(eventpbevents.TimerDelay(
		m.config.Timer,
		[]*eventpbtypes.Event{rnEvents.RetransmitAll(m.config.Self)},
		timert.Duration(m.params.RetransmissionLoopInterval),
	))

	m.iterationCount++
	return evsOut, nil
}

func (m *Module) SendMessage(id rntypes.MsgID, msg *messagepbtypes.Message, destinations []t.NodeID) (*events.EventList, error) {
	qmsg := &QueuedMessage{
		Msg:          msg,
		Destinations: destinations,
		Ts:           m.iterationCount,
	}

	listEl := m.queue.PushBack(qmsg)
	m.ensureSubqueueIndexExists(msg.DestModule)
	m.queueIndex[msg.DestModule][id] = listEl

	evsOut := events.ListOf(
		transportpbevents.SendMessage(m.config.Net, msg, destinations),
	)

	if !m.retransmitLoopScheduled {
		// m.logger.Log(logging.LevelDebug, "rescheduling retransmission loop")
		evsOut.PushBack(eventpbevents.TimerDelay(
			m.config.Timer,
			[]*eventpbtypes.Event{rnEvents.RetransmitAll(m.config.Self)},
			timert.Duration(m.params.RetransmissionLoopInterval),
		))
		m.retransmitLoopScheduled = true
	}

	return evsOut, nil
}

// Initial capacity for message queue associated to a module
// The queue is created lazily, so it will at least have one message added to it
// In a ideal scenario, are replicas are more or less up-to-date so we should only need to keep
// a few messages around.
const initialSubqueueIndexCapacity = 4

func (m *Module) ensureSubqueueIndexExists(destModule t.ModuleID) {
	if _, present := m.queueIndex[destModule]; !present {
		m.queueIndex[destModule] = make(map[rntypes.MsgID]*list.Element, initialSubqueueIndexCapacity)
	}
}

func (m *Module) MarkModuleMsgsRecvd(destModule t.ModuleID, destinations []t.NodeID) (*events.EventList, error) {
	// optimization for all destinations
	if len(destinations) == len(m.params.AllNodes) {
		for id, subqueueIndex := range m.queueIndex {
			if !id.IsSubOf(destModule) && id != destModule {
				continue
			}

			for _, listEl := range subqueueIndex {
				m.queue.Remove(listEl)
			}
			delete(m.queueIndex, id)
		}

		return &events.EventList{}, nil
	}

	for modID, subqueueIndex := range m.queueIndex {
		if !modID.IsSubOf(destModule) && modID != destModule {
			continue
		}

		for msgID, listEl := range subqueueIndex {
			qmsg := listEl.Value.(*QueuedMessage)
			qmsg.Destinations = sliceutil.Filter(qmsg.Destinations, func(_ int, d t.NodeID) bool {
				return !slices.Contains(destinations, d)
			})

			if len(qmsg.Destinations) == 0 {
				m.queue.Remove(listEl)
				delete(subqueueIndex, msgID)
			}
		}

		if len(subqueueIndex) == 0 {
			delete(m.queueIndex, modID)
		}
	}

	return &events.EventList{}, nil
}

func (m *Module) MarkRecvd(destModule t.ModuleID, messageID rntypes.MsgID, destinations []t.NodeID) (*events.EventList, error) {
	if subqueueIndex, present := m.queueIndex[destModule]; present {
		if listEl, present := subqueueIndex[messageID]; present {
			qmsg := listEl.Value.(*QueuedMessage)
			qmsg.Destinations = sliceutil.Filter(qmsg.Destinations, func(_ int, d t.NodeID) bool {
				return !slices.Contains(destinations, d)
			})

			if len(qmsg.Destinations) == 0 {
				m.queue.Remove(listEl)
				delete(subqueueIndex, messageID)
			}
		}

		if len(subqueueIndex) == 0 {
			delete(m.queueIndex, destModule)
		}
	}

	return &events.EventList{}, nil
}
