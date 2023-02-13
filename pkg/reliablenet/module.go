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

	retransmitLoopScheduled bool // TODO: atomic.Bool in Go 1.19

	queues map[t.ModuleID]map[string]*eventpb.SendMessage // ModuleID -> msgID ([]byte) -> queuedMsg
}

func New(ownID t.NodeID, mc *ModuleConfig, params *ModuleParams, logger logging.Logger) (*Module, error) {
	return &Module{
		config: mc,
		params: params,
		ownID:  ownID,
		logger: logger,

		retransmitLoopScheduled: false,

		queues: make(map[t.ModuleID]map[string]*eventpb.SendMessage),
	}, nil
}

func (m *Module) GetPendingMessages() []*eventpb.SendMessage {
	pendingMsgs := make([]*eventpb.SendMessage, 0, len(m.queues))
	for _, queue := range m.queues {
		for _, qmsg := range queue {
			pendingMsgs = append(pendingMsgs, qmsg)
		}
	}

	return pendingMsgs
}

func (m *Module) CountPendingMessages() int {
	count := 0

	for _, queue := range m.queues {
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
	if msgSource == m.ownID {
		// fast path for local messages
		return m.MarkRecvd(msgDestModule, string(msgID), []t.NodeID{msgSource})
	}

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
	evsOut := &events.EventList{}

	if len(m.queues) == 0 {
		m.retransmitLoopScheduled = false
		m.logger.Log(logging.LevelDebug, "stopping retransmission loop: no more messages")
		return evsOut, nil
	}

	for moduleID, queue := range m.queues {
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
				return evsOut, fmt.Errorf("queued message for module %v has no destinations", moduleID)
			}
		}

		if len(queue) == 0 {
			return evsOut, fmt.Errorf("queue for module %v has no messages", moduleID)
		}
	}

	evsOut.PushBack(events.TimerDelay(
		m.config.Timer,
		[]*eventpb.Event{rnEvents.RetransmitAll(m.config.Self).Pb()},
		t.TimeDuration(m.params.RetransmissionLoopInterval),
	))

	return evsOut, nil
}

func (m *Module) SendMessage(id string, msg *messagepb.Message, destinations []t.NodeID) (*events.EventList, error) {
	m.ensureQueueExists(t.ModuleID(msg.DestModule))
	m.queues[t.ModuleID(msg.DestModule)][id] = &eventpb.SendMessage{
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

func (m *Module) ensureQueueExists(destModule t.ModuleID) {
	if _, present := m.queues[destModule]; !present {
		m.queues[destModule] = make(map[string]*eventpb.SendMessage, initialQueueCapacity)
	}
}

func (m *Module) MarkModuleMsgsRecvd(destModule t.ModuleID, destinations []t.NodeID) (*events.EventList, error) {
	// optimization for all destinations
	if len(destinations) == len(m.params.AllNodes) {
		for id := range m.queues {
			if id.IsSubOf(destModule) || id == destModule {
				delete(m.queues, id)
			}
		}

		return &events.EventList{}, nil
	}

	for queueID, queue := range m.queues {
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
			delete(m.queues, queueID)
		}
	}

	return &events.EventList{}, nil
}

func (m *Module) MarkRecvd(destModule t.ModuleID, messageID string, destinations []t.NodeID) (*events.EventList, error) {
	queue, queueExists := m.queues[destModule]
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
		delete(m.queues, destModule)
	}

	return &events.EventList{}, nil
}
