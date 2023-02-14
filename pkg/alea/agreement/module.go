package agreement

import (
	"fmt"
	"strconv"

	"github.com/filecoin-project/mir/pkg/abba"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/abbapb"
	abbaEvents "github.com/filecoin-project/mir/pkg/pb/abbapb/events"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	aagEvents "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/events"
	aagMsgs "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/msgs"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	rnEvents "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/events"
	"github.com/filecoin-project/mir/pkg/serializing"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID
	Hasher       t.ModuleID
	ReliableNet  t.ModuleID
	Net          t.ModuleID
	ThreshCrypto t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         "alea_ag",
		Consumer:     "alea_dir",
		Hasher:       "hasher",
		ReliableNet:  "reliablenet",
		Net:          "net",
		ThreshCrypto: "threshcrypto",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // must be the alea ID followed by 'a'
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
}

type agModule struct {
	config *ModuleConfig
	params *ModuleParams
	nodeID t.NodeID
	logger logging.Logger

	roundDecisionHistory []bool

	currentAbba modules.PassiveModule

	currentRound uint64
	inputDone    bool
	delivered    bool
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, logger logging.Logger) modules.PassiveModule {
	return &agModule{
		config: mc,
		params: params,
		nodeID: nodeID,
		logger: logger,

		roundDecisionHistory: nil, // TODO: consider using a bitvec (will reduce space and batch reallocations as a bonus)

		currentAbba: nil,

		currentRound: 0,
		inputDone:    false,
		delivered:    false,
	}
}

func (m *agModule) ApplyEvents(eventsIn *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsSequentially(eventsIn, m.applyEvent)
}

func (m *agModule) ImplementsModule() {}

func (m *agModule) applyEvent(event *eventpb.Event) (*events.EventList, error) {
	if event.DestModule != m.config.Self.Pb() {
		return m.proxyABBAEvent(event)
	}

	switch e := event.Type.(type) {
	case *eventpb.Event_Init:
		m.currentRound = 0
		return m.initializeRound() // first round is not lazy, for simplicity
	case *eventpb.Event_MessageReceived:
		return m.handleMessageReceived(e.MessageReceived.Msg, t.NodeID(e.MessageReceived.From))
	case *eventpb.Event_AleaAgreement:
		return m.handleAgreementEvent(e.AleaAgreement)
	case *eventpb.Event_Abba:
		return m.handleABBAEvent(e.Abba)
	// TODO: WAL and app snapshot handling
	default:
		return nil, fmt.Errorf("unsupported event type: %T", e)
	}
}

func (m *agModule) handleMessageReceived(message *messagepb.Message, from t.NodeID) (*events.EventList, error) {
	msgWrapper, ok := message.Type.(*messagepb.Message_AleaAgreement)
	if !ok {
		m.logger.Log(logging.LevelDebug, "unknown message type")
		return &events.EventList{}, nil
	}

	msg, ok := msgWrapper.AleaAgreement.Type.(*agreementpb.Message_FinishAbba)
	if !ok {
		m.logger.Log(logging.LevelDebug, "unknown agreement message type")
		return &events.EventList{}, nil
	}

	if msg.FinishAbba.Round == m.currentRound {
		// we've fallen back, and this replica is trying to help us catch up
		destModuleID := m.abbaModuleID()

		return m.currentAbba.ApplyEvents(events.ListOf(
			events.MessageReceived(
				destModuleID,
				from,
				abba.FinishMessage(destModuleID, msg.FinishAbba.Value).Pb(),
			),
		))
	}
	// message is stale

	return &events.EventList{}, nil
}

func (m *agModule) proxyABBAEvent(event *eventpb.Event) (*events.EventList, error) {
	destSub := t.ModuleID(event.DestModule).StripParent(m.config.Self)
	r, ok := strconv.ParseUint(string(destSub.Top()), 10, 64)
	if ok != nil {
		return &events.EventList{}, nil // bogus message
	}

	if r == m.currentRound {
		// all good, event can be safely forwarded
	} else if ev, isMsg := event.Type.(*eventpb.Event_MessageReceived); r < m.currentRound && isMsg {
		// old round: can be a replica trying to help us catch up, or one stalled in the past.
		// we'll inform its agreement component that we're done, and help propagate the final decision
		// see agModule::handleMessageReceived

		from := t.NodeID(ev.MessageReceived.From)
		return events.ListOf(
			events.SendMessage(m.config.Net, aagMsgs.FinishAbbaMessage(
				m.config.Self,
				r,
				m.roundDecisionHistory[r],
			).Pb(), []t.NodeID{from}),

			// ABBA does not mark FINISH(_) messages as received, so we must acknowledge any we may receive
			rnEvents.Ack(
				m.config.ReliableNet,
				m.config.Self.Then(t.NewModuleIDFromInt(r)).Then(t.ModuleID("__global")), // TODO: use constant from abba
				abba.FinishMsgID(),
				from,
			).Pb(),
		), nil
	} else if m.delivered && r == m.currentRound+1 {
		// other nodes are moving to the next agreement round, follow suit
		eventsOut, err := m.advanceRound()

		// only process the original event later, to allow the new ABBA instance to initialize
		eventsOut.PushBack(event)

		return eventsOut, err
	} else {
		// stray event (stale and we shouldn't touch it or from the future and we can't handle it yet)
		return &events.EventList{}, nil
	}

	evsOut, err := m.currentAbba.ApplyEvents((&events.EventList{}).PushBack(event))

	if err == nil && !m.inputDone {
		// request input value into new agreement round
		// may end up happening multiple times.
		evsOut.PushBack(aagEvents.RequestInput(m.config.Consumer, m.currentRound).Pb())
	}

	return evsOut, err
}

func (m *agModule) handleAgreementEvent(event *agreementpb.Event) (*events.EventList, error) {
	evWrapped, ok := event.Type.(*agreementpb.Event_InputValue)
	if !ok {
		return nil, fmt.Errorf("unexpected agreement event: %v", event)
	}
	ev := evWrapped.Unwrap()

	evsOut := &events.EventList{}

	if ev.Round == m.currentRound+1 && m.delivered {
		// time for a new round
		evsOut, err := m.advanceRound()
		m.logger.Log(logging.LevelDebug, "Advanced agreement round", "agreementRound", m.currentRound)
		if err != nil {
			return evsOut, err
		}
	} else if m.inputDone || m.currentRound != ev.Round {
		// stale message
		m.logger.Log(logging.LevelDebug, "discarding stale InputValue event", "agreementRound", ev.Round, "input", ev.Input)
		return &events.EventList{}, nil
	}

	m.inputDone = true

	moreEvsOut, err := m.currentAbba.ApplyEvents((&events.EventList{}).PushBack(
		abbaEvents.InputValue(m.abbaModuleID(), ev.Input).Pb(),
	))
	if err != nil {
		return evsOut, err
	}

	return evsOut.PushBackList(moreEvsOut), nil
}

func (m *agModule) handleABBAEvent(event *abbapb.Event) (*events.EventList, error) {
	evWrapped, ok := event.Type.(*abbapb.Event_Deliver)
	if !ok {
		return nil, fmt.Errorf("unexpected abba event: %v", event)
	}
	ev := evWrapped.Unwrap()

	if m.delivered {
		return nil, fmt.Errorf("abba module double-delivered result")
	}

	evsOut := &events.EventList{}

	// notify Alea of delivery
	m.delivered = true
	evsOut.PushBack(aagEvents.Deliver(m.config.Consumer, m.currentRound, ev.Result).Pb())

	// free up message queues for past round (we'll just re-send the FINAL message later if needed)
	if len(m.roundDecisionHistory) != int(m.currentRound) {
		return nil, fmt.Errorf("inconsistent agreement historical data")
	}
	m.roundDecisionHistory = append(m.roundDecisionHistory, ev.Result)

	return evsOut, nil
}

func (m *agModule) abbaModuleID() t.ModuleID {
	return m.config.Self.Then(t.NewModuleIDFromInt(m.currentRound))
}

func (m *agModule) advanceRound() (*events.EventList, error) {
	clearOldMessageQueue := rnEvents.MarkModuleMsgsRecvd(
		m.config.ReliableNet,
		m.config.Self.Then(t.NewModuleIDFromInt(m.currentRound)),
		m.params.AllNodes,
	).Pb()

	m.currentRound++

	evsOut, err := m.initializeRound()
	evsOut.PushBack(clearOldMessageQueue)

	return evsOut, err
}

func (m *agModule) initializeRound() (*events.EventList, error) {
	m.delivered = false
	m.inputDone = false

	abbaModuleID := m.abbaModuleID()

	instanceUID := make([]byte, len(m.params.InstanceUID)+8)
	instanceUID = append(instanceUID, m.params.InstanceUID...)
	instanceUID = append(instanceUID, serializing.Uint64ToBytes(m.currentRound)...)

	m.currentAbba = abba.NewModule(&abba.ModuleConfig{
		Self:         abbaModuleID,
		Consumer:     m.config.Self,
		ReliableNet:  m.config.ReliableNet,
		ThreshCrypto: m.config.ThreshCrypto,
		Hasher:       m.config.Hasher,
	}, &abba.ModuleParams{
		InstanceUID: instanceUID,
		AllNodes:    m.params.AllNodes,
	}, m.nodeID, logging.Decorate(m.logger, "Abba: ", "agreementRound", m.currentRound))

	return m.currentAbba.ApplyEvents((&events.EventList{}).PushBack(events.Init(abbaModuleID)))
}
