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
	ID           rntypes.MsgID
	Msg          *messagepb.Message
	Destinations []t.NodeID
	Ts           uint64
}

type Module struct {
	config ModuleConfig
	params *ModuleParams
	ownID  t.NodeID
	logger logging.Logger

	retransmitLoopScheduled bool
	iterationCount          uint64

	queue      []*QueuedMessage
	queueHead  int
	queueIndex map[t.ModuleID]map[rntypes.MsgID]*QueuedMessage
}

func New(ownID t.NodeID, mc ModuleConfig, params *ModuleParams, logger logging.Logger) (*Module, error) {
	return &Module{
		config: mc,
		params: params,
		ownID:  ownID,
		logger: logger,

		retransmitLoopScheduled: false,

		queueIndex: make(map[t.ModuleID]map[rntypes.MsgID]*QueuedMessage),
	}, nil
}

func (m *Module) GetPendingMessages() []*QueuedMessage {
	pending := make([]*QueuedMessage, 0, len(m.queue))
	for _, qmsg := range m.queue {
		if len(qmsg.Destinations) > 0 {
			pending = append(pending, qmsg)
		}
	}

	return pending
}

func (m *Module) CountPendingMessages() int {
	return len(m.GetPendingMessages())
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

	// clean up ack-ed msgs
	for len(m.queue) > 0 && len(m.queue[0].Destinations) > 0 {
		m.queue = m.queue[1:]
		if m.queueHead > 0 {
			m.queueHead--
		}
	}

	if len(m.queue) == 0 {
		m.retransmitLoopScheduled = false
		return evsOut, nil
	}

	nRemaining := m.params.MaxRetransmissionBurst
	if nRemaining > len(m.queue) {
		nRemaining = len(m.queue)
	}

	restarted := m.queueHead == 0
	origQueueHead := m.queueHead

	for !restarted && nRemaining > 0 {
		if m.queue[0].Ts < m.iterationCount {
			// messages are too recent, bring back queue head
			m.queueHead = 0
			nRemaining = origQueueHead // don't retransmit twice
			restarted = true
		}

		for nRemaining > 0 && m.queue[m.queueHead].Ts < m.iterationCount {
			qmsg := m.queue[m.queueHead]
			if len(qmsg.Destinations) == 0 {
				continue
			}

			evsOut.PushBack(
				transportpbevents.SendMessage(m.config.Net, messagepbtypes.MessageFromPb(qmsg.Msg), qmsg.Destinations).Pb(),
			)

			m.queueHead++
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

func (m *Module) SendMessage(id rntypes.MsgID, msg *messagepb.Message, destinations []t.NodeID) (*events.EventList, error) {
	qmsg := &QueuedMessage{
		ID:           id,
		Msg:          msg,
		Destinations: destinations,
		Ts:           m.iterationCount,
	}

	m.ensureSubqueueIndexExists(t.ModuleID(msg.DestModule))
	m.queueIndex[t.ModuleID(msg.DestModule)][id] = qmsg
	m.queue = append(m.queue, qmsg)

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
const initialSubqueueIndexCapacity = 4

func (m *Module) ensureSubqueueIndexExists(destModule t.ModuleID) {
	if _, present := m.queueIndex[destModule]; !present {
		m.queueIndex[destModule] = make(map[rntypes.MsgID]*QueuedMessage, initialSubqueueIndexCapacity)
	}
}

func (m *Module) MarkModuleMsgsRecvd(destModule t.ModuleID, destinations []t.NodeID) (*events.EventList, error) {
	// optimization for all destinations
	if len(destinations) == len(m.params.AllNodes) {
		for id, subqueueIndex := range m.queueIndex {
			if !id.IsSubOf(destModule) && id != destModule {
				continue
			}

			for _, qmsg := range subqueueIndex {
				qmsg.Destinations = nil
			}
			delete(m.queueIndex, id)
		}

		return &events.EventList{}, nil
	}

	for modID, subqueueIndex := range m.queueIndex {
		if !modID.IsSubOf(destModule) && modID != destModule {
			continue
		}

		for msgID, qmsg := range subqueueIndex {
			qmsg.Destinations = sliceutil.Filter(qmsg.Destinations, func(_ int, d t.NodeID) bool {
				return !slices.Contains(destinations, d)
			})
			subqueueIndex[msgID] = qmsg

			if len(qmsg.Destinations) == 0 {
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
		if qmsg, present := subqueueIndex[messageID]; present {
			qmsg.Destinations = sliceutil.Filter(qmsg.Destinations, func(_ int, d t.NodeID) bool {
				return !slices.Contains(destinations, d)
			})
			subqueueIndex[messageID] = qmsg

			if len(qmsg.Destinations) == 0 {
				delete(subqueueIndex, messageID)
			}
		}

		if len(subqueueIndex) == 0 {
			delete(m.queueIndex, destModule)
		}
	}

	return &events.EventList{}, nil
}
