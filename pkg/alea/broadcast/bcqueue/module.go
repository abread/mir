package bcqueue

import (
	"sync"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	rnEvents "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/events"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	vcbpbevents "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	vcbpbtypes "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	"github.com/filecoin-project/mir/pkg/serializing"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb"
)

type QueueBcModule struct {
	config *ModuleConfig
	params *ModuleParams
	nodeID t.NodeID
	logger logging.Logger

	windowSizeCtrl *WindowSizeController
	slots          map[uint64]*queueSlot

	locker sync.Mutex
}

func New(mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) *QueueBcModule {
	return &QueueBcModule{
		config: mc,
		params: params,
		nodeID: nodeID,
		logger: logger,

		windowSizeCtrl: NewWindowSizeController(tunables.MaxConcurrentVcb),
		slots:          make(map[uint64]*queueSlot, tunables.MaxConcurrentVcb),
	}
}

type queueSlot struct {
	module modules.PassiveModule
	locker sync.Mutex
}

func (m *QueueBcModule) RouteEventToSlot(slotID uint64, event *eventpb.Event) (*events.EventList, error) {
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

func (m *QueueBcModule) getOrTryCreateQueueSlot(slotID uint64) (*queueSlot, *events.EventList, error) {
	// TODO: try to break this apart, the return type is disgusting
	m.locker.Lock()
	defer m.locker.Unlock()

	if !m.windowSizeCtrl.IsSlotInView(slotID) || m.windowSizeCtrl.IsSlotFreed(slotID) {
		return nil, nil, nil
	}

	if m.windowSizeCtrl.IsSlotUnused(slotID) {
		// we need to create the slot module first
		m.windowSizeCtrl.Acquire(slotID)

		instanceUID := append(m.params.InstanceUID, serializing.Uint64ToBytes(slotID)...)
		newSlotModuleID := m.config.Self.Then(t.NewModuleIDFromInt(slotID))

		newSlot := &queueSlot{
			module: vcb.NewModule(
				&vcb.ModuleConfig{
					Self:         newSlotModuleID,
					ReliableNet:  m.config.ReliableNet,
					ThreshCrypto: m.config.ThreshCrypto,
					Mempool:      m.config.Mempool,
				},
				&vcb.ModuleParams{
					InstanceUID: instanceUID,
					AllNodes:    m.params.AllNodes,
					Leader:      m.params.QueueOwner,
				},
				m.nodeID,
				logging.Decorate(m.logger, "Vcb: ", "queueSlot", slotID),
			),
		}

		initEvents := events.ListOf(events.Init(newSlotModuleID))
		if m.params.QueueOwner != m.nodeID {
			initEvents.PushBack(vcbpbevents.InputValue(newSlotModuleID, nil, &vcbpbtypes.Origin{
				// TODO: route slot delivery through ourselves
				Module: m.config.Consumer,
				Type: &vcbpbtypes.Origin_AleaBc{
					// TODO: use a proper origin
					AleaBc: &bcpbtypes.BcOrigin{
						Module: newSlotModuleID,
					},
				},
			}).Pb())
		}

		initOutEvents, err := newSlot.module.ApplyEvents(initEvents)
		if err != nil {
			return nil, nil, err
		}

		m.slots[slotID] = newSlot

		m.logger.Log(logging.LevelDebug, "Created slot/VCB instance", "queueSlot", slotID)
		return newSlot, initOutEvents, nil
	}

	return m.slots[slotID], &events.EventList{}, nil
}

func (m *QueueBcModule) FreeSlot(slotID uint64) (*events.EventList, error) {
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
		).Pb(),
	), nil
}

func (m *QueueBcModule) StartBroadcast(slotID uint64, txs []*requestpb.Request) (*events.EventList, error) {
	// TODO: rethink design to leave only one side responsible for window size control
	// assumes caller checked that current queue window has this slot in view (and unused)

	destModule := m.config.Self.Then(t.NewModuleIDFromInt(slotID))

	outEvent := vcbpbevents.InputValue(destModule, txs, &vcbpbtypes.Origin{
		// TODO: route slot delivery through ourselves
		Module: m.config.Consumer,
		Type: &vcbpbtypes.Origin_AleaBc{
			// TODO: use a proper origin
			AleaBc: &bcpbtypes.BcOrigin{
				Module: destModule,
			},
		},
	}).Pb()

	// routing logic will create the vcb instance for us
	m.logger.Log(logging.LevelDebug, "Starting broadcast", "queueSlot", slotID)
	return (&events.EventList{}).PushBack(outEvent), nil
}
