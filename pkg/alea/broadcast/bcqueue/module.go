package bcqueue

import (
	"context"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bcutil"
	aleaCommon "github.com/filecoin-project/mir/pkg/alea/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modring"
	"github.com/filecoin-project/mir/pkg/modules"
	bcqueuedsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	bcqueuepbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/events"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	batchdbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/dsl"
	reliablenetpbdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	requestpbtypes "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	vcbpbdsl "github.com/filecoin-project/mir/pkg/pb/vcbpb/dsl"
	vcbpbevents "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	vcbpbtypes "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb"
)

func New(ctx context.Context, mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	if slices.Index(params.AllNodes, params.QueueOwner) != int(params.QueueIdx) {
		return nil, fmt.Errorf("invalid queue index/owner combination: %v - %v", params.QueueIdx, params.QueueOwner)
	}

	slots := modring.New(ctx, mc.Self, tunables.MaxConcurrentVcb, modring.ModuleParams{
		Generator: newVcbGenerator(mc, params, nodeID, logger),
	}, logging.Decorate(logger, "Modring controller: "))

	controller := newQueueController(ctx, mc, params, tunables, nodeID, logger, slots)

	return modules.RoutedModule(mc.Self, controller, slots), nil
}

func newQueueController(ctx context.Context, mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger, slots *modring.Module) modules.PassiveModule {
	m := dsl.NewModule(ctx, mc.Self)

	bcqueuedsl.UponInputValue(m, func(slot *commontypes.Slot, txs []*requestpbtypes.Request) error {
		if slot.QueueIdx != params.QueueIdx {
			return fmt.Errorf("input value given to wrong queue")
		}

		if len(txs) == 0 {
			return fmt.Errorf("cannot broadcast an empty batch")
		}

		// logger.Log(logging.LevelDebug, "starting broadcast", "queueSlot", slot.QueueSlot, "txs", txs)
		dsl.EmitMirEvent(m, vcbpbevents.InputValue(
			mc.Self.Then(t.NewModuleIDFromInt(slot.QueueSlot)),
			txs,
		))
		return nil
	})

	// we can't use .UponDeliver because that assumes a DSL origin
	// upon vcb deliver, store batch and deliver to broadcast component
	vcbpbdsl.UponEvent[*vcbpbtypes.Event_Deliver](m, func(ev *vcbpbtypes.Deliver) error {
		queueSlotStr := ev.SrcModule.StripParent(mc.Self).Top()
		queueSlot, err := strconv.ParseUint(string(queueSlotStr), 10, 64)
		if err != nil {
			return fmt.Errorf("deliver event for invalid round: %w", err)
		}

		slot := &commontypes.Slot{
			QueueIdx:  params.QueueIdx,
			QueueSlot: aleatypes.QueueSlot(queueSlot),
		}

		batchdbdsl.StoreBatch(m, mc.BatchDB, aleaCommon.FormatAleaBatchID(slot), ev.TxIds, ev.Txs, ev.Signature, slot)
		return nil
	})

	batchdbdsl.UponBatchStored(m, func(slot *commontypes.Slot) error {
		// logger.Log(logging.LevelDebug, "delivering broadcast", "queueSlot", slot.QueueSlot)
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

		// clean up old messages
		reliablenetpbdsl.MarkModuleMsgsRecvd(m, mc.ReliableNet, mc.Self.Then(t.NewModuleIDFromInt(queueSlot)), params.AllNodes)

		return nil
	})

	return m
}

func newVcbGenerator(queueMc *ModuleConfig, queueParams *ModuleParams, nodeID t.NodeID, logger logging.Logger) func(ctx context.Context, id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
	baseConfig := &vcb.ModuleConfig{
		Self:         "INVALID",
		Consumer:     queueMc.Self,
		ReliableNet:  queueMc.ReliableNet,
		Hasher:       queueMc.Hasher,
		ThreshCrypto: queueMc.ThreshCrypto,
		Mempool:      queueMc.Mempool,
	}

	baseParams := &vcb.ModuleParams{
		InstanceUID: nil,
		AllNodes:    queueParams.AllNodes,
		Leader:      queueParams.QueueOwner,
	}

	return func(ctx context.Context, id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
		mc := *baseConfig
		mc.Self = id

		queueSlot := aleatypes.QueueSlot(idx)

		params := *baseParams
		params.InstanceUID = bcutil.VCBInstanceUID(queueParams.BcInstanceUID, queueParams.QueueIdx, queueSlot)

		mod := vcb.NewModule(ctx, &mc, &params, nodeID, logging.Decorate(logger, "Vcb: ", "slot", idx))

		return mod, events.ListOf(
			bcqueuepbevents.BcStarted(queueMc.Consumer, &commontypes.Slot{
				QueueIdx:  queueParams.QueueIdx,
				QueueSlot: queueSlot,
			}).Pb(),
		), nil
	}
}
