package broadcast

import (
	"context"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bcqueue"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	bcdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcqueuedsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID
	BatchDB      t.ModuleID
	Mempool      t.ModuleID
	ReliableNet  t.ModuleID
	ThreshCrypto t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         "alea_bc",
		Consumer:     "alea_dir",
		BatchDB:      "batchdb",
		Mempool:      "mempool",
		ReliableNet:  "reliablenet",
		ThreshCrypto: "threshcrypto",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // must be the alea instance uid followed by 'b'
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
}

// ModuleTunables sets the values of protocol tunables that need not be the same across all nodes.
type ModuleTunables struct {
	// Maximum number of concurrent VCB instances per queue
	// Must match the equally named tunable in the main Alea module
	// Must be at least 1
	MaxConcurrentVcbPerQueue int
}

func NewModule(ctx context.Context, mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	controller, err := newQueueController(ctx, mc, params, tunables, nodeID, logger)
	if err != nil {
		return nil, err
	}

	queues, err := createQueues(ctx, mc, params, tunables, nodeID, logger)
	if err != nil {
		return nil, err
	}

	return &bcMod{
		ownID:      mc.Self,
		controller: controller,
		queues:     queues,
		logger:     logger,
	}, nil
}

type bcMod struct {
	ownID      t.ModuleID
	controller modules.PassiveModule
	queues     []modules.PassiveModule
	logger     logging.Logger
}

func (m *bcMod) ImplementsModule() {}

func (m *bcMod) ApplyEvents(evs *events.EventList) (*events.EventList, error) {
	ownEvents, queueEvents := m.splitEvents(evs)

	evsOut, err := m.controller.ApplyEvents(ownEvents)
	if err != nil {
		return nil, err
	}

	for qIdx, qEvs := range queueEvents {
		if qEvs.Len() > 0 {
			qEvsOut, err := m.queues[qIdx].ApplyEvents(qEvs)
			if err != nil {
				return nil, err
			}

			evsOut.PushBackList(qEvsOut)
		}
	}

	return evsOut, nil
}

func (m *bcMod) splitEvents(evs *events.EventList) (*events.EventList, []*events.EventList) {
	ownEvents := events.EmptyList()
	queueEvents := make([]*events.EventList, len(m.queues))
	for i := 0; i < len(queueEvents); i++ {
		queueEvents[i] = events.EmptyList()
	}

	it := evs.Iterator()
	for ev := it.Next(); ev != nil; ev = it.Next() {
		if ev.DestModule == string(m.ownID) {
			ownEvents.PushBack(ev)
		} else {
			dest := t.ModuleID(ev.DestModule).StripParent(m.ownID).Top()
			idx, err := strconv.ParseInt(string(dest), 10, 64)
			if err != nil || int(idx) >= len(queueEvents) || idx < 0 {
				m.logger.Log(logging.LevelDebug, "event for invalid queue", "event", ev)
				continue // invalid event
			}

			queueEvents[idx].PushBack(ev)
		}
	}

	return ownEvents, queueEvents
}

func newQueueController(ctx context.Context, mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	m := dsl.NewModule(ctx, mc.Self)

	ownQueueIdx := slices.Index(params.AllNodes, nodeID)
	if ownQueueIdx == -1 {
		return nil, fmt.Errorf("own node not present in node list")
	}

	ownQueueModID := mc.Self.Then(t.NewModuleIDFromInt(ownQueueIdx))

	dsl.UponInit(m, func() error {
		// init queues
		for i := 0; i < len(params.AllNodes); i++ {
			dsl.EmitEvent(m, events.Init(mc.Self.Then(t.NewModuleIDFromInt(i))))
		}
		return nil
	})

	bcdsl.UponStartBroadcast(m, func(queueSlot aleatypes.QueueSlot, txs []*requestpb.Request) error {
		bcqueuedsl.InputValue(m, ownQueueModID, &commontypes.Slot{
			QueueIdx:  aleatypes.QueueIdx(ownQueueIdx),
			QueueSlot: queueSlot,
		}, txs)
		return nil
	})

	bcqueuedsl.UponDeliver(m, func(slot *commontypes.Slot) error {
		bcdsl.Deliver(m, mc.Consumer, slot)
		return nil
	})

	bcdsl.UponFreeSlot(m, func(slot *commontypes.Slot) error {
		bcqueuedsl.FreeSlot(m, mc.Self.Then(t.NewModuleIDFromInt(slot.QueueIdx)), slot.QueueSlot)
		return nil
	})

	return m, nil
}

func createQueues(ctx context.Context, queueMc *ModuleConfig, queueParams *ModuleParams, queueTunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) ([]modules.PassiveModule, error) {
	tunables := &bcqueue.ModuleTunables{
		MaxConcurrentVcb: queueTunables.MaxConcurrentVcbPerQueue,
	}

	queues := make([]modules.PassiveModule, len(queueParams.AllNodes))

	for idx := 0; idx < len(queueParams.AllNodes); idx++ {
		mc := &bcqueue.ModuleConfig{
			Self:         queueMc.Self.Then(t.NewModuleIDFromInt(idx)),
			Consumer:     queueMc.Self,
			BatchDB:      queueMc.BatchDB,
			Mempool:      queueMc.Mempool,
			ReliableNet:  queueMc.ReliableNet,
			ThreshCrypto: queueMc.ThreshCrypto,
		}

		params := &bcqueue.ModuleParams{
			InstanceUID: queueParams.InstanceUID, // TODO: review
			AllNodes:    queueParams.AllNodes,

			QueueIdx:   aleatypes.QueueIdx(idx),
			QueueOwner: queueParams.AllNodes[idx],
		}

		mod, err := bcqueue.New(ctx, mc, params, tunables, nodeID, logging.Decorate(logger, "BcQueue: ", "queueIdx", idx))
		if err != nil {
			return nil, err
		}

		queues[idx] = mod
	}

	return queues, nil
}
