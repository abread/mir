package broadcast

import (
	"sync"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	"github.com/filecoin-project/mir/pkg/pb/vcbpb"
	rnEvents "github.com/filecoin-project/mir/pkg/reliablenet/events"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb"
)

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
