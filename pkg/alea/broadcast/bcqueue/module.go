package bcqueue

import (
	"strconv"

	es "github.com/go-errors/errors"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bcutil"
	aleaCommon "github.com/filecoin-project/mir/pkg/alea/util"
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
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	vcbpbdsl "github.com/filecoin-project/mir/pkg/pb/vcbpb/dsl"
	vcbpbevents "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb"
)

func New(mc ModuleConfig, params ModuleParams, tunables ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	if slices.Index(params.AllNodes, params.QueueOwner) != int(params.QueueIdx) {
		return nil, es.Errorf("invalid queue index/owner combination: %v - %v", params.QueueIdx, params.QueueOwner)
	}

	slots := modring.New(mc.Self, tunables.MaxConcurrentVcb, modring.ModuleParams{
		Generator: newVcbGenerator(mc, params, nodeID, logger),
	}, logging.Decorate(logger, "Modring controller: "))

	controller := newQueueController(mc, params, logger, slots)

	return modules.RoutedModule(mc.Self, controller, slots), nil
}

func newQueueController(mc ModuleConfig, params ModuleParams, logger logging.Logger, slots *modring.Module) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	bcqueuedsl.UponInputValue(m, func(queueSlot aleatypes.QueueSlot, txs []*trantorpbtypes.Transaction) error {
		if len(txs) == 0 {
			return es.Errorf("cannot broadcast an empty batch")
		}

		// logger.Log(logging.LevelDebug, "starting broadcast", "queueSlot", queueSlot, "txs", txs)
		dsl.EmitMirEvent(m, vcbpbevents.InputValue(
			mc.Self.Then(t.NewModuleIDFromInt(queueSlot)),
			txs,
		))
		return nil
	})

	// upon vcb deliver, store batch and deliver to broadcast component
	vcbpbdsl.UponDeliver(m, func(txs []*trantorpbtypes.Transaction, txIds []tt.TxID, signature tctypes.FullSig, srcModule t.ModuleID) error {
		queueSlotStr := srcModule.StripParent(mc.Self).Top()
		queueSlot, err := strconv.ParseUint(string(queueSlotStr), 10, 64)
		if err != nil {
			return es.Errorf("deliver event for invalid round: %w", err)
		}

		slot := &commontypes.Slot{
			QueueIdx:  params.QueueIdx,
			QueueSlot: aleatypes.QueueSlot(queueSlot),
		}

		batchdbdsl.StoreBatch(m, mc.BatchDB, aleaCommon.FormatAleaBatchID(slot), txIds, txs, signature, slot)
		return nil
	})

	batchdbdsl.UponBatchStored(m, func(slot *commontypes.Slot) error {
		logger.Log(logging.LevelDebug, "delivering broadcast", "queueSlot", slot.QueueSlot)
		bcqueuedsl.Deliver(m, mc.Consumer, slot)
		return nil
	})

	bcqueuedsl.UponFreeSlot(m, func(queueSlot aleatypes.QueueSlot) error {
		if err := slots.AdvanceViewToAtLeastSubmodule(uint64(queueSlot)); err != nil {
			return es.Errorf("could not advance view to free queue slot: %w", err)
		}

		if err := slots.MarkSubmodulePast(uint64(queueSlot)); err != nil {
			return es.Errorf("failed to free queue slot: %w", err)
		}

		// clean up old messages
		reliablenetpbdsl.MarkModuleMsgsRecvd(m, mc.ReliableNet, mc.Self.Then(t.NewModuleIDFromInt(queueSlot)), params.AllNodes)

		return nil
	})

	return m
}

func newVcbGenerator(queueMc ModuleConfig, queueParams ModuleParams, nodeID t.NodeID, logger logging.Logger) func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
	baseConfig := vcb.ModuleConfig{
		Self:         "INVALID",
		Consumer:     queueMc.Self,
		Net:          queueMc.Net,
		ReliableNet:  queueMc.ReliableNet,
		Hasher:       queueMc.Hasher,
		ThreshCrypto: queueMc.ThreshCrypto,
		Mempool:      queueMc.Mempool,
	}

	baseParams := vcb.ModuleParams{
		InstanceUID: nil,
		AllNodes:    queueParams.AllNodes,
		Leader:      queueParams.QueueOwner,
	}

	return func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
		mc := baseConfig
		mc.Self = id

		queueSlot := aleatypes.QueueSlot(idx)

		params := baseParams
		params.InstanceUID = bcutil.VCBInstanceUID(queueParams.BcInstanceUID, queueParams.QueueIdx, queueSlot)

		mod := vcb.NewModule(mc, &params, nodeID, logging.Decorate(logger, "Vcb: ", "slot", idx))

		return mod, events.ListOf(
			bcqueuepbevents.BcStarted(queueMc.Consumer, &commontypes.Slot{
				QueueIdx:  queueParams.QueueIdx,
				QueueSlot: queueSlot,
			}).Pb(),
		), nil
	}
}
