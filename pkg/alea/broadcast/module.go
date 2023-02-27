package broadcast

import (
	"context"
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bcqueue"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modring"
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

	queues := modring.New(mc.Self, len(params.AllNodes), modring.ModuleParams{
		Generator: newBcQueueGenerator(mc, params, tunables, nodeID, logger),
	}, logging.Decorate(logger, "Modring controller: "))

	return modules.RoutedModule(mc.Self, controller, queues), nil
}

func newQueueController(ctx context.Context, mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	m := dsl.NewModule(ctx, mc.Self)

	ownQueueIdx := slices.Index(params.AllNodes, nodeID)
	if ownQueueIdx == -1 {
		return nil, fmt.Errorf("own node not present in node list")
	}

	ownQueueModID := mc.Self.Then(t.NewModuleIDFromInt(ownQueueIdx))

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

func newBcQueueGenerator(queueMc *ModuleConfig, queueParams *ModuleParams, queueTunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
	tunables := &bcqueue.ModuleTunables{
		MaxConcurrentVcb: queueTunables.MaxConcurrentVcbPerQueue,
	}

	return func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
		mc := &bcqueue.ModuleConfig{
			Self:         id,
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
			QueueOwner: queueParams.AllNodes[int(idx)],
		}

		mod, err := bcqueue.New(context.TODO(), mc, params, tunables, nodeID, logging.Decorate(logger, "BcQueue: ", "queueIdx", idx))
		if err != nil {
			return nil, nil, err
		}

		return mod, events.EmptyList(), nil
	}
}
