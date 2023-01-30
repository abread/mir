package broadcast

import (
	"fmt"
	"strconv"
	"sync"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/broadcast/abcevents"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	"github.com/filecoin-project/mir/pkg/pb/vcbpb"
	rnEvents "github.com/filecoin-project/mir/pkg/reliablenet/events"
	"github.com/filecoin-project/mir/pkg/serializing"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID
	Mempool      t.ModuleID
	ReliableNet  t.ModuleID
	ThreshCrypto t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         "alea_bc",
		Consumer:     "alea_dir",
		Mempool:      "mempool",
		ReliableNet:  "reliablenet",
		ThreshCrypto: "threshcrypto",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // must be the same as the one in the main and agreement alea components
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
}

// ModuleTunables sets the values of protocol tunables that need not be the same across all nodes.
type ModuleTunables struct {
	// Maximum number of concurrent VCB instances per queue
	// Must match the equally named tunable in the main Alea module
	// Must be at least 1
	MaxConcurrentVcbPerQueue int
}

type bcModule struct {
	config         *ModuleConfig
	ownNodeID      t.NodeID
	ownQueueIdx    int
	queueBcModules []queueBcModule

	logger logging.Logger
}

type queueBcModule struct {
	config     *ModuleConfig
	params     *ModuleParams
	queueOwner t.NodeID
	queueIdx   uint32
	nodeID     t.NodeID
	logger     logging.Logger

	windowSizeCtrl *WindowSizeController
	slots          map[uint64]*queueSlot

	locker sync.Mutex
}

type queueSlot struct {
	module modules.PassiveModule
	locker sync.Mutex
}

func NewModule(mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	queueBcModules := make([]queueBcModule, len(params.AllNodes))

	ownQueueIdx := slices.Index(params.AllNodes, nodeID)
	if ownQueueIdx == -1 {
		return nil, fmt.Errorf("nodeID not present in AllNodes")
	}

	for idx, queueOwner := range params.AllNodes {
		config := *mc
		config.Consumer = mc.Self
		config.Self = mc.Self.Then(t.NewModuleIDFromInt(idx))

		queueBcModules[idx] = queueBcModule{
			config:     &config,
			params:     params,
			queueOwner: queueOwner,
			queueIdx:   uint32(idx),
			nodeID:     nodeID,
			logger:     logging.Decorate(logger, "AleaBcQueue: ", "queueIdx", idx),

			windowSizeCtrl: NewWindowSizeController(tunables.MaxConcurrentVcbPerQueue),
			slots:          make(map[uint64]*queueSlot, tunables.MaxConcurrentVcbPerQueue),

			locker: sync.Mutex{},
		}
	}

	return &bcModule{
		config:         mc,
		ownNodeID:      nodeID,
		ownQueueIdx:    ownQueueIdx,
		queueBcModules: queueBcModules,
		logger:         logger,
	}, nil
}

func (m *bcModule) ApplyEvents(evs *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsConcurrently(evs, m.applyEvent)
}

func (m *bcModule) ImplementsModule() {}

func (m *bcModule) applyEvent(event *eventpb.Event) (*events.EventList, error) {
	if event.DestModule != string(m.config.Self) {
		return m.routeEventToQueue(event)
	}

	switch e := event.Type.(type) {
	case *eventpb.Event_Init:
		return &events.EventList{}, nil // no-op
	case *eventpb.Event_AleaBroadcast:
		return m.handleBroadcastEvent(e.AleaBroadcast)
	case *eventpb.Event_Vcb:
		return m.handleVcbEvent(e.Vcb)
	default:
		return nil, fmt.Errorf("unsupported event type: %T", e)
	}
}

func (m *bcModule) handleBroadcastEvent(event *bcpb.Event) (*events.EventList, error) {
	switch e := event.Type.(type) {
	case *bcpb.Event_StartBroadcast:
		return m.handleStartBroadcast(e.StartBroadcast)
	case *bcpb.Event_FreeSlot:
		return m.handleFreeSlot(e.FreeSlot)
	default:
		return nil, fmt.Errorf("unexpected broadcast event type: %T", e)
	}
}

func (m *bcModule) handleStartBroadcast(event *bcpb.StartBroadcast) (*events.EventList, error) {
	return m.queueBcModules[m.ownQueueIdx].StartBroadcast(event.QueueSlot, event.TxIds, event.Txs)
}

func (m *bcModule) handleFreeSlot(event *bcpb.FreeSlot) (*events.EventList, error) {
	queueIdx := int(event.Slot.QueueIdx)
	slotID := event.Slot.QueueSlot

	return m.queueBcModules[queueIdx].FreeSlot(slotID)
}

func (m *bcModule) handleVcbEvent(event *vcbpb.Event) (*events.EventList, error) {
	evWrapped, ok := event.Type.(*vcbpb.Event_Deliver)
	if !ok {
		return nil, fmt.Errorf("unexpected abba event: %v", event)
	}
	ev := evWrapped.Unwrap()

	// TODO: try to decouple this from queue internal module hierarchy
	originModuleSuffix := t.ModuleID(ev.OriginModule).StripParent(m.config.Self)
	queueIdx, err1 := strconv.ParseUint(string(originModuleSuffix.Top()), 10, 32)
	if err1 != nil {
		return nil, fmt.Errorf("could not parse queue idx: %w", err1)
	}
	queueSlot, err2 := strconv.ParseUint(string(originModuleSuffix.Sub().Top()), 10, 64)
	if err2 != nil {
		return nil, fmt.Errorf("could not parse queue slot: %w", err2)
	}

	slot := &common.Slot{
		QueueIdx:  uint32(queueIdx),
		QueueSlot: queueSlot,
	}

	m.logger.Log(logging.LevelInfo, "Delivered BC slot", "queueIdx", queueIdx, "queueSlot", queueSlot)

	return events.ListOf(
		abcevents.Deliver(m.config.Consumer, slot, t.TxIDSlice(ev.TxIds), ev.Txs, ev.Signature),
		rnEvents.MarkModuleMsgsRecvd(m.config.ReliableNet, t.ModuleID(ev.OriginModule), []t.NodeID{m.ownNodeID}),
	), nil
}

func (m *bcModule) routeEventToQueue(event *eventpb.Event) (*events.EventList, error) {
	destSub := t.ModuleID(event.DestModule).StripParent(m.config.Self)
	idx, err := strconv.Atoi(string(destSub.Top()))
	if err != nil || idx < 0 || idx >= len(m.queueBcModules) {
		// bogus message
		return &events.EventList{}, nil
	}

	slot, err := strconv.ParseUint(string(destSub.Sub().Top()), 10, 64)
	if err != nil {
		// bogus message
		return &events.EventList{}, nil
	}

	return m.queueBcModules[idx].RouteEventToSlot(slot, event)
}

func (m *queueBcModule) RouteEventToSlot(slotID uint64, event *eventpb.Event) (*events.EventList, error) {
	slot, outEvents, err := m.getOrTryCreateQueueSlot(slotID)
	if err != nil {
		return nil, err
	} else if slot == nil {
		// bogus message
		return &events.EventList{}, nil
	}

	slot.locker.Lock()
	defer slot.locker.Unlock()
	return slot.module.ApplyEvents(outEvents.PushBack(event))
}

func (m *queueBcModule) getOrTryCreateQueueSlot(slotID uint64) (*queueSlot, *events.EventList, error) {
	// TODO: try to break this apart, the return type is disgusting
	m.locker.Lock()
	defer m.locker.Unlock()

	if !m.windowSizeCtrl.IsSlotInView(slotID) || m.windowSizeCtrl.IsSlotFreed(slotID) {
		return nil, nil, nil
	}

	if m.windowSizeCtrl.IsSlotUnused(slotID) {
		// we need to create the slot module first
		m.windowSizeCtrl.Acquire(slotID)

		instanceUID := VCBInstanceUID(m.params.InstanceUID, m.queueIdx, slotID)
		newSlotModuleID := m.config.Self.Then(t.NewModuleIDFromInt(slotID))

		newSlot := &queueSlot{
			module: vcb.NewModule(
				&vcb.ModuleConfig{
					Self:         newSlotModuleID,
					Consumer:     m.config.Consumer,
					ReliableNet:  m.config.ReliableNet,
					ThreshCrypto: m.config.ThreshCrypto,
					Mempool:      m.config.Mempool,
				},
				&vcb.ModuleParams{
					InstanceUID: instanceUID,
					AllNodes:    m.params.AllNodes,
					Leader:      m.queueOwner,
				},
				m.nodeID,
				logging.Decorate(m.logger, "Vcb: ", "queueSlot", slotID),
			),
		}

		initOutEvents, err := newSlot.module.ApplyEvents(events.ListOf(events.Init(newSlotModuleID)))
		if err != nil {
			return nil, nil, err
		}

		m.slots[slotID] = newSlot

		m.logger.Log(logging.LevelDebug, "Created slot/VCB instance", "queueSlot", slotID)
		return newSlot, initOutEvents, nil
	}

	return m.slots[slotID], &events.EventList{}, nil
}

func (m *queueBcModule) FreeSlot(slotID uint64) (*events.EventList, error) {
	m.locker.Lock()
	defer m.locker.Unlock()

	m.windowSizeCtrl.Free(slotID)
	delete(m.slots, slotID)

	m.logger.Log(logging.LevelDebug, "Freed slot/VCB instance", "queueSlot", slotID)

	return events.ListOf(
		rnEvents.MarkModuleMsgsRecvd(
			m.config.ReliableNet,
			m.config.Self.Then(t.NewModuleIDFromInt(slotID)),
			m.params.AllNodes,
		),
	), nil
}

func (m *queueBcModule) StartBroadcast(slotID uint64, txIDsPb [][]byte, txs []*requestpb.Request) (*events.EventList, error) {
	// TODO: rethink design to leave only one side responsible for window size control
	// assumes caller checked that current queue window has this slot in view (and unused)

	destModule := m.config.Self.Then(t.NewModuleIDFromInt(slotID))

	// TODO: use some helper method from within vcb to create this
	outEvent := &eventpb.Event{
		DestModule: destModule.Pb(),
		Type: &eventpb.Event_Vcb{
			Vcb: &vcbpb.Event{
				Type: &vcbpb.Event_Request{
					Request: &vcbpb.BroadcastRequest{
						TxIds: txIDsPb,
						Txs:   txs,
					},
				},
			},
		},
	}

	// routing logic will create the vcb instance for us
	m.logger.Log(logging.LevelDebug, "Starting broadcast", "queueSlot", slotID)
	return (&events.EventList{}).PushBack(outEvent), nil
}

func VCBInstanceUID(bcInstanceUID []byte, queueIdx uint32, queueSlot uint64) []byte {
	uid := slices.Clone(bcInstanceUID)
	uid = append(uid, []byte("bc")...)
	uid = append(uid, serializing.Uint64ToBytes(uint64(queueIdx))...)
	uid = append(uid, serializing.Uint64ToBytes(queueSlot)...)

	return uid
}
