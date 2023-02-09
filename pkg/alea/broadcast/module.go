package broadcast

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/broadcast/abcevents"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bcqueue"
	"github.com/filecoin-project/mir/pkg/alea/common"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	aleaCommonPb "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb"
	batchdbpbevents "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/events"
	batchdbpbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/types"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/vcbpb"
	"github.com/filecoin-project/mir/pkg/serializing"
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

type bcModule struct {
	config         *ModuleConfig
	ownNodeID      t.NodeID
	ownQueueIdx    int
	queueBcModules []*bcqueue.QueueBcModule

	logger logging.Logger
}

func NewModule(mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	queueBcModules := make([]*bcqueue.QueueBcModule, len(params.AllNodes))

	ownQueueIdx := slices.Index(params.AllNodes, nodeID)
	if ownQueueIdx == -1 {
		return nil, fmt.Errorf("nodeID not present in AllNodes")
	}

	bcqueueConfig := bcqueue.ModuleConfig{
		Self:         mc.Self,
		Consumer:     mc.Self,
		BatchDB:      mc.BatchDB,
		Mempool:      mc.Mempool,
		ReliableNet:  mc.ReliableNet,
		ThreshCrypto: mc.ThreshCrypto,
	}
	bcqueueTunables := bcqueue.ModuleTunables{
		MaxConcurrentVcb: tunables.MaxConcurrentVcbPerQueue,
	}
	for idx, queueOwner := range params.AllNodes {
		config := bcqueueConfig
		config.Self = mc.Self.Then(t.NewModuleIDFromInt(idx))

		uid := make([]byte, len(params.InstanceUID)+4)
		uid = append(uid, params.InstanceUID...)
		uid = append(uid, serializing.Uint32ToBytes(uint32(idx))...)

		queueBcModules[idx] = bcqueue.New(
			&config,
			&bcqueue.ModuleParams{
				InstanceUID: uid,
				AllNodes:    params.AllNodes,
				QueueOwner:  queueOwner,
			},
			&bcqueueTunables,
			nodeID,
			logging.Decorate(logger, "AleaBcQueue: ", "queueIdx", idx),
		)
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
	case *eventpb.Event_BatchDb:
		return m.handleBatchDBEvent(e.BatchDb)
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

	slot := &aleaCommonPb.Slot{
		QueueIdx:  uint32(queueIdx),
		QueueSlot: queueSlot,
	}

	// store delivered slot
	m.logger.Log(logging.LevelDebug, "storing batch", "queueIdx", queueIdx, "queueSlot", queueSlot)
	return events.ListOf(
		batchdbpbevents.StoreBatch(
			m.config.BatchDB,
			common.FormatAleaBatchID(slot).Pb(),
			ev.TxIds,
			ev.Txs,
			ev.Signature,
			AleaBcStoreBatchOrigin(m.config.Self, slot),
		).Pb(),
	), nil
}

func AleaBcStoreBatchOrigin(module t.ModuleID, slot *aleaCommonPb.Slot) *batchdbpbtypes.StoreBatchOrigin {
	return &batchdbpbtypes.StoreBatchOrigin{
		Module: module,
		Type: &batchdbpbtypes.StoreBatchOrigin_AleaBroadcast{
			AleaBroadcast: &bcpb.StoreBatchOrigin{
				Slot: slot,
			},
		},
	}
}

func (m *bcModule) handleBatchDBEvent(event *batchdbpb.Event) (*events.EventList, error) {
	evWrapped, ok := event.Type.(*batchdbpb.Event_Stored)
	if !ok {
		return nil, fmt.Errorf("unexpected batchdb event: %v", event)
	}
	ev := evWrapped.Unwrap()

	origin, ok := ev.Origin.Type.(*batchdbpb.StoreBatchOrigin_AleaBroadcast)
	if !ok {
		return nil, fmt.Errorf("unexpected batchdb event (unknown origin): %v", event)
	}

	slot := origin.AleaBroadcast.Slot

	m.logger.Log(logging.LevelInfo, "Delivered BC slot", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)

	return events.ListOf(
		abcevents.Deliver(m.config.Consumer, slot),
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
