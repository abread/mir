package bcqueue

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	aleaCommon "github.com/filecoin-project/mir/pkg/alea/common"
	batchdbdsl "github.com/filecoin-project/mir/pkg/availability/batchdb/dsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modring"
	"github.com/filecoin-project/mir/pkg/modules"
	bcqueuedsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	bcqueuepbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	vcbpbdsl "github.com/filecoin-project/mir/pkg/pb/vcbpb/dsl"
	vcbpbevents "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	vcbpbtypes "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb"
)

func New(mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	if slices.Index(params.AllNodes, params.QueueOwner) != int(params.QueueIdx) {
		return nil, fmt.Errorf("invalid queue index/owner combination: %v - %v", params.QueueIdx, params.QueueOwner)
	}

	slots := modring.New(mc.Self, tunables.MaxConcurrentVcb, modring.ModuleParams{
		Generator: newVcbGenerator(mc, params, nodeID, logger),
	}, logging.Decorate(logger, "Modring controller: "))

	controller := newQueueController(mc, params, tunables, nodeID, logger, slots)

	return modules.RoutedModule(mc.Self, controller, slots), nil
}

func newQueueController(mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger, slots *modring.Module) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	bcqueuedsl.UponInputValue(m, func(slot *commontypes.Slot, txs []*requestpb.Request) error {
		if slot.QueueIdx != params.QueueIdx {
			return fmt.Errorf("input value given to wrong queue")
		}

		dsl.EmitMirEvent(m, vcbpbevents.InputValue(
			mc.Self.Then(t.NewModuleIDFromInt(slot.QueueSlot)),
			txs,
			&vcbpbtypes.Origin{
				Module: mc.Self,
				Type: &vcbpbtypes.Origin_AleaBcqueue{
					AleaBcqueue: &bcqueuepbtypes.VcbOrigin{
						QueueSlot: slot.QueueSlot,
					},
				},
			},
		))
		return nil
	})

	// we can't use .UponDeliver because that assumes a DSL origin
	vcbpbdsl.UponEvent[*vcbpbtypes.Event_Deliver](m, func(ev *vcbpbtypes.Deliver) error {
		origin, ok := ev.Origin.Type.(*vcbpbtypes.Origin_AleaBcqueue)
		if !ok {
			return fmt.Errorf("invalid origin for vcb instance")
		}
		slot := &commontypes.Slot{
			QueueIdx:  params.QueueIdx,
			QueueSlot: origin.AleaBcqueue.QueueSlot,
		}

		batchdbdsl.StoreBatch(m, mc.BatchDB, aleaCommon.FormatAleaBatchID(slot), ev.TxIds, ev.Txs, ev.Signature, slot)
		return nil
	})
	batchdbdsl.UponBatchStored(m, func(slot *commontypes.Slot) error {
		bcqueuedsl.Deliver(m, mc.Consumer, slot)
		return nil
	})

	bcqueuedsl.UponFreeSlot(m, func(queueSlot aleatypes.QueueSlot) error {
		if err := slots.AdvanceViewToAtLeastSubmodule(uint64(queueSlot)); err != nil {
			return fmt.Errorf("could not advance view to free queue slot: %w", err)
		}

		if err := slots.MarkSubmodulePast(uint64(queueSlot)); err != nil {
			return fmt.Errorf("failed to free queue slot: %w", err)
		}
		return nil
	})

	return m
}

func newVcbGenerator(queueMc *ModuleConfig, queueParams *ModuleParams, nodeID t.NodeID, logger logging.Logger) func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
	baseConfig := &vcb.ModuleConfig{
		Self:         "INVALID",
		ReliableNet:  queueMc.ReliableNet,
		ThreshCrypto: queueMc.ThreshCrypto,
		Mempool:      queueMc.Mempool,
	}

	params := &vcb.ModuleParams{
		InstanceUID: queueParams.InstanceUID, // TODO: review
		AllNodes:    queueParams.AllNodes,
		Leader:      queueParams.QueueOwner,
	}

	return func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
		mc := *baseConfig
		mc.Self = id

		mod := vcb.NewModule(&mc, params, nodeID, logging.Decorate(logger, "Vcb: ", "slot", idx))

		initialEvs := &events.EventList{}
		// events.Init is taken care of by modring
		if nodeID != params.Leader {
			initialEvs.PushBack(vcbpbevents.InputValue(mc.Self, nil, &vcbpbtypes.Origin{
				Module: queueMc.Self,
				Type: &vcbpbtypes.Origin_AleaBcqueue{
					AleaBcqueue: &bcqueuepbtypes.VcbOrigin{
						QueueSlot: aleatypes.QueueSlot(idx),
					},
				},
			}).Pb())
		}

		return mod, initialEvs, nil
	}
}
