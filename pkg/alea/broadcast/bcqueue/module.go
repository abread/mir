package bcqueue

import (
	"context"
	"fmt"
	"strconv"

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
	bcqueuepbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/events"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	mempooldsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	reliablenetpbdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	threshcryptopbdsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	vcbpbdsl "github.com/filecoin-project/mir/pkg/pb/vcbpb/dsl"
	vcbpbevents "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	vcbpbtypes "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb"
)

func New(ctx context.Context, mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	if slices.Index(params.AllNodes, params.QueueOwner) != int(params.QueueIdx) {
		return nil, fmt.Errorf("invalid queue index/owner combination: %v - %v", params.QueueIdx, params.QueueOwner)
	}

	slots := modring.New(ctx, mc.Self, tunables.MaxConcurrentVcb, modring.ModuleParams{
		Generator:      newVcbGenerator(mc, params, nodeID, logger),
		PastMsgHandler: newPastMsgHandler(mc, params),
	}, logging.Decorate(logger, "Modring controller: "))

	controller := newQueueController(ctx, mc, params, tunables, nodeID, logger, slots)

	return modules.RoutedModule(mc.Self, controller, slots), nil
}

func newQueueController(ctx context.Context, mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger, slots *modring.Module) modules.PassiveModule {
	m := dsl.NewModule(ctx, mc.Self)

	bcqueuedsl.UponInputValue(m, func(slot *commontypes.Slot, txs []*requestpb.Request) error {
		if slot.QueueIdx != params.QueueIdx {
			return fmt.Errorf("input value given to wrong queue")
		}

		if len(txs) == 0 {
			return fmt.Errorf("cannot broadcast an empty batch")
		}

		logger.Log(logging.LevelDebug, "starting broadcast", "queueSlot", slot.QueueSlot, "txs", txs)
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

	// upon stale vcb final message, validate signature, store and deliver batch
	bcqueuedsl.UponPastVcbFinal(m, func(queueSlot aleatypes.QueueSlot, txs []*requestpb.Request, signature tctypes.FullSig) error {
		logger.Log(logging.LevelDebug, "processing out-of-order VCB final message", "queueSlot", queueSlot)

		context := &processPastVcbCtx{
			queueSlot: queueSlot,
			txs:       txs,
			signature: signature,
		}
		mempooldsl.RequestTransactionIDs(m, mc.Mempool, txs, context)

		reliablenetpbdsl.Ack(m, mc.ReliableNet, mc.Self.Then(t.NewModuleIDFromInt(queueSlot)), vcb.FinalMsgID(), params.QueueOwner)

		return nil
	})
	mempooldsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *processPastVcbCtx) error {
		context.txIDs = txIDs
		threshcryptopbdsl.VerifyFull(m, mc.ThreshCrypto, vcb.SigData(params.InstanceUID, txIDs), context.signature, context)
		return nil
	})
	threshcryptopbdsl.UponVerifyFullResult(m, func(ok bool, error string, context *processPastVcbCtx) error {
		if !ok {
			return nil // bad batch
			// TODO: report byz node?
		}

		logger.Log(logging.LevelDebug, "out-of-order VCB Final is valid. storing batch for delivery", "queueIdx", params.QueueIdx, "queueSlot", context.queueSlot)

		slot := &commontypes.Slot{
			QueueIdx:  params.QueueIdx,
			QueueSlot: context.queueSlot,
		}
		batchdbdsl.StoreBatch(m, mc.BatchDB, aleaCommon.FormatAleaBatchID(slot), context.txIDs, context.txs, context.signature, slot)
		return nil
	})

	batchdbdsl.UponBatchStored(m, func(slot *commontypes.Slot) error {
		logger.Log(logging.LevelDebug, "delivering broadcast", "queueSlot", slot.QueueSlot)
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

		// do not mark VCB FINAL messages as received, as they can avoid FILL-GAP messages

		return nil
	})

	return m
}

func newVcbGenerator(queueMc *ModuleConfig, queueParams *ModuleParams, nodeID t.NodeID, logger logging.Logger) func(ctx context.Context, id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
	baseConfig := &vcb.ModuleConfig{
		Self:         "INVALID",
		Consumer:     queueMc.Self,
		ReliableNet:  queueMc.ReliableNet,
		ThreshCrypto: queueMc.ThreshCrypto,
		Mempool:      queueMc.Mempool,
	}

	params := &vcb.ModuleParams{
		InstanceUID: queueParams.InstanceUID, // TODO: review
		AllNodes:    queueParams.AllNodes,
		Leader:      queueParams.QueueOwner,
	}

	return func(ctx context.Context, id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
		mc := *baseConfig
		mc.Self = id

		mod := vcb.NewModule(ctx, &mc, params, nodeID, logging.Decorate(logger, "Vcb: ", "slot", idx))

		return mod, &events.EventList{}, nil
	}
}

func newPastMsgHandler(mc *ModuleConfig, params *ModuleParams) func([]*modringpbtypes.PastMessage) (*events.EventList, error) {
	return func(pastMessages []*modringpbtypes.PastMessage) (*events.EventList, error) {
		evsOut := &events.EventList{}

		for _, msg := range pastMessages {
			if msg.From != params.QueueOwner {
				continue
			}

			vcbMsgWrapper, ok := msg.Message.Type.(*messagepbtypes.Message_Vcb)
			if !ok {
				continue
			}

			vcbFinalMsgWrapper, ok := vcbMsgWrapper.Unwrap().Type.(*vcbpbtypes.Message_FinalMessage)
			if !ok {
				continue
			}
			vcbFinalMsg := vcbFinalMsgWrapper.Unwrap()

			evsOut.PushBack(bcqueuepbevents.PastVcbFinal(
				mc.Self,
				aleatypes.QueueSlot(msg.DestId),
				vcbFinalMsg.Txs,
				vcbFinalMsg.Signature,
			).Pb())
		}

		return evsOut, nil
	}
}

type processPastVcbCtx struct {
	queueSlot aleatypes.QueueSlot
	txs       []*requestpb.Request
	txIDs     []t.TxID
	signature tctypes.FullSig
}
