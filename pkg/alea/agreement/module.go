package agreement

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/abba"
	abbaEvents "github.com/filecoin-project/mir/pkg/abba/abbaevents"
	aagEvents "github.com/filecoin-project/mir/pkg/alea/agreement/aagevents"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/abbapb"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/serializing"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Alea         t.ModuleID
	Mempool      t.ModuleID
	Net          t.ModuleID
	ThreshCrypto t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         "alea_ag",
		Alea:         "alea",
		Mempool:      "mempool",
		Net:          "net",
		ThreshCrypto: "threshcrypto",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // must be the same as the one in the main and broadcast alea components
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
}

type agModule struct {
	config *ModuleConfig
	params *ModuleParams
	nodeID t.NodeID
	logger logging.Logger

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

		currentAbba: nil,

		currentRound: 0,
		inputDone:    false,
		delivered:    false,
	}
}

func (m agModule) ApplyEvents(eventsIn *events.EventList) (*events.EventList, error) {
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
		m.inputDone = false
		m.delivered = false
		return &events.EventList{}, nil
	case *eventpb.Event_AleaAgreement:
		return m.handleAgreementEvent(e.AleaAgreement)
	case *eventpb.Event_Abba:
		return m.handleABBAEvent(e.Abba)
	// TODO: WAL and app snapshot handling
	default:
		return nil, fmt.Errorf("unsupported event type: %T", e)
	}
}

func (m *agModule) proxyABBAEvent(event *eventpb.Event) (*events.EventList, error) {
	destSub := t.ModuleID(event.DestModule).StripParent(m.config.Self)

	if destSub.Top() == t.NewModuleIDFromInt(m.currentRound) {
		// all good, event can be safely forwarded
	} else if m.delivered && destSub.Top() == t.NewModuleIDFromInt(m.currentRound+1) {
		// other nodes are moving to the next agreement round, follow suit
		eventsOut, err := m.advanceRound()

		// only process the original event later, to allow the new ABBA instance to initialize
		eventsOut.PushBack(event)

		// request input value into new agreement round
		eventsOut.PushBack(aagEvents.RequestInput(m.config.Alea, m.currentRound))

		return eventsOut, err
	} else {
		// stray event (stale and we shouldn't touch it or from the future and we can't handle it yet)
		return &events.EventList{}, nil
	}

	return m.currentAbba.ApplyEvents((&events.EventList{}).PushBack(event))
}

func (m *agModule) handleAgreementEvent(event *agreementpb.Event) (*events.EventList, error) {
	evWrapped, ok := event.Type.(*agreementpb.Event_InputValue)
	if !ok {
		return nil, fmt.Errorf("unexpected agreement event: %v", event)
	}
	ev := evWrapped.Unwrap()

	if m.inputDone || m.currentRound != ev.Round {
		// stale message
		return &events.EventList{}, nil
	}

	// we skip going around the event loop, and forward this straight to the ABBA instance
	m.inputDone = true
	return m.currentAbba.ApplyEvents((&events.EventList{}).PushBack(
		abbaEvents.InputValue(m.abbaModuleID(), ev.Input),
	))
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

	// notify Alea of delivery
	m.delivered = true
	return (&events.EventList{}).PushBack(aagEvents.Deliver(m.config.Alea, m.currentRound, ev.Result)), nil
}

func (m *agModule) abbaModuleID() t.ModuleID {
	return m.config.Self.Then(t.NewModuleIDFromInt(m.currentRound))
}

func (m *agModule) advanceRound() (*events.EventList, error) {
	m.currentRound++

	m.delivered = false
	m.inputDone = false

	abbaModuleID := m.abbaModuleID()

	instanceUID := slices.Clone(m.params.InstanceUID)
	instanceUID = append(instanceUID, serializing.Uint64ToBytes(m.currentRound)...)

	m.currentAbba = abba.NewModule(&abba.ModuleConfig{
		Self:         abbaModuleID,
		Consumer:     m.config.Self,
		Net:          m.config.Net,
		ThreshCrypto: m.config.ThreshCrypto,
		Hasher:       m.config.Mempool,
	}, &abba.ModuleParams{
		InstanceUID: instanceUID,
		AllNodes:    m.params.AllNodes,
	}, m.nodeID, logging.Decorate(m.logger, "[alea-ag-abba]", "agreementRound", m.currentRound))

	return m.currentAbba.ApplyEvents((&events.EventList{}).PushBack(events.Init(abbaModuleID)))
}
