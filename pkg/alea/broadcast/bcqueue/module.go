package bcqueue

import (
	"strconv"
	"time"

	es "github.com/go-errors/errors"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modring"
	"github.com/filecoin-project/mir/pkg/modules"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	bcqueuepbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/events"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	reliablenetpbdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	transportpbevents "github.com/filecoin-project/mir/pkg/pb/transportpb/events"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	vcbpbdsl "github.com/filecoin-project/mir/pkg/pb/vcbpb/dsl"
	vcbpbevents "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	vcbpbmsgs "github.com/filecoin-project/mir/pkg/pb/vcbpb/msgs"
	vcbpbtypes "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
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
		Generator:      newVcbGenerator(mc, params, nodeID, logger),
		PastMsgHandler: newPastMsgHandler(mc, params),
	}, logging.Decorate(logger, "Modring controller: "))

	controller := newQueueController(mc, params, tunables, nodeID, logger, slots)

	return modules.RoutedModule(mc.Self, controller, slots), nil
}

func newQueueController(mc ModuleConfig, params ModuleParams, tunables ModuleTunables, nodeID t.NodeID, logger logging.Logger, slots *modring.Module) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	bcqueuepbdsl.UponInputValue(m, func(queueSlot aleatypes.QueueSlot, txIDs []tt.TxID, txs []*trantorpbtypes.Transaction) error {
		if len(txs) == 0 {
			return es.Errorf("cannot broadcast an empty batch")
		}

		// logger.Log(logging.LevelDebug, "starting broadcast", "queueSlot", queueSlot, "txs", txs)
		dsl.EmitEvent(m, vcbpbevents.InputValue(
			mc.Self.Then(t.NewModuleIDFromInt(queueSlot)),
			txIDs,
			txs,
		))
		return nil
	})

	// upon vcb deliver, store batch and deliver to broadcast component
	vcbpbdsl.UponDeliver(m, func(batchID string, signature tctypes.FullSig, srcModule t.ModuleID) error {
		queueSlotStr := srcModule.StripParent(mc.Self).Top()
		queueSlot, err := strconv.ParseUint(string(queueSlotStr), 10, 64)
		if err != nil {
			return es.Errorf("deliver event for invalid round: %w", err)
		}

		slot := &bcpbtypes.Slot{
			QueueIdx:  params.QueueIdx,
			QueueSlot: aleatypes.QueueSlot(queueSlot),
		}

		logger.Log(logging.LevelDebug, "delivering broadcast", "queueSlot", slot.QueueSlot)
		bcqueuepbdsl.Deliver(m, mc.Consumer, &bcpbtypes.Cert{
			Slot:      slot,
			BatchId:   batchID,
			Signature: signature,
		})
		return nil
	})

	if params.QueueOwner == nodeID {
		slotDeliverTimes := make(map[aleatypes.QueueSlot]time.Time, tunables.MaxConcurrentVcb)

		vcbpbdsl.UponDeliver(m, func(_ string, _ tctypes.FullSig, srcModule t.ModuleID) error {
			queueSlotStr := srcModule.StripParent(mc.Self).Top()
			queueSlot, err := strconv.ParseUint(string(queueSlotStr), 10, 64)
			if err != nil {
				return es.Errorf("deliver event for invalid round: %w", err)
			}

			slotDeliverTimes[aleatypes.QueueSlot(queueSlot)] = time.Now()
			return nil
		})

		vcbpbdsl.UponQuorumDone(m, func(srcModule t.ModuleID) error {
			queueSlotStr := srcModule.StripParent(mc.Self).Top()
			queueSlot, err := strconv.ParseUint(string(queueSlotStr), 10, 64)
			if err != nil {
				return es.Errorf("deliver event for invalid round: %w", err)
			}

			slot := &bcpbtypes.Slot{
				QueueIdx:  params.QueueIdx,
				QueueSlot: aleatypes.QueueSlot(queueSlot),
			}

			if startTime, ok := slotDeliverTimes[slot.QueueSlot]; ok {
				deliverDelta := time.Since(startTime)
				bcqueuepbdsl.BcQuorumDone(m, mc.Consumer, slot, deliverDelta)

				// add delta to avoid including it in the fully done delta calc
				slotDeliverTimes[slot.QueueSlot] = startTime.Add(deliverDelta)
			}

			return nil
		})

		vcbpbdsl.UponAllDone(m, func(srcModule t.ModuleID) error {
			queueSlotStr := srcModule.StripParent(mc.Self).Top()
			queueSlot, err := strconv.ParseUint(string(queueSlotStr), 10, 64)
			if err != nil {
				return es.Errorf("deliver event for invalid round: %w", err)
			}

			slot := &bcpbtypes.Slot{
				QueueIdx:  params.QueueIdx,
				QueueSlot: aleatypes.QueueSlot(queueSlot),
			}

			if startTime, ok := slotDeliverTimes[slot.QueueSlot]; ok {
				deliverDelta := time.Since(startTime)
				bcqueuepbdsl.BcAllDone(m, mc.Consumer, slot, deliverDelta)
			}

			return nil
		})

		bcqueuepbdsl.UponFreeSlot(m, func(queueSlot aleatypes.QueueSlot) error {
			delete(slotDeliverTimes, queueSlot)
			return nil
		})
	}

	bcqueuepbdsl.UponFreeSlot(m, func(queueSlot aleatypes.QueueSlot) error {
		// advance queue to not get stuck in old slots
		// we could go to the latest slot, but if we center it around the latest slot, we can still
		// recover slow broadcasts *and* accept new ones.
		var minSlot uint64
		if uint64(queueSlot) < uint64(tunables.MaxConcurrentVcb) {
			minSlot = 0
		} else {
			minSlot = uint64(queueSlot) - uint64(tunables.MaxConcurrentVcb)/2
		}
		if err := slots.AdvanceViewToAtLeastSubmodule(minSlot); err != nil {
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

func newPastMsgHandler(mc ModuleConfig, params ModuleParams) modring.PastMessagesHandler {
	return func(pastMessages []*modringpbtypes.PastMessage) (events.EventList, error) {
		evsOut := events.EmptyList()
		for _, msg := range pastMessages {
			vcbMsgW, ok := msg.Message.Type.(*messagepbtypes.Message_Vcb)
			if !ok {
				continue
			}

			_, ok = vcbMsgW.Vcb.Type.(*vcbpbtypes.Message_FinalMessage)
			if !ok {
				continue
			}

			if msg.From != params.QueueOwner {
				continue
			}

			evsOut.PushBack(transportpbevents.SendMessage(mc.Net, vcbpbmsgs.DoneMessage(msg.Message.DestModule), []t.NodeID{msg.From}))
		}

		return evsOut, nil
	}
}

func newVcbGenerator(queueMc ModuleConfig, queueParams ModuleParams, nodeID t.NodeID, logger logging.Logger) func(id t.ModuleID, idx uint64) (modules.PassiveModule, events.EventList, error) {
	baseConfig := vcb.ModuleConfig{
		Self:         "INVALID",
		Consumer:     queueMc.Self,
		Net:          queueMc.Net,
		ReliableNet:  queueMc.ReliableNet,
		ThreshCrypto: queueMc.ThreshCrypto,
		Mempool:      queueMc.Mempool,
		BatchDB:      queueMc.BatchDB,
	}

	baseParams := vcb.ModuleParams{
		InstanceUID: nil,
		EpochNr:     queueParams.EpochNr,
		AllNodes:    queueParams.AllNodes,
		Leader:      queueParams.QueueOwner,
	}

	return func(id t.ModuleID, idx uint64) (modules.PassiveModule, events.EventList, error) {
		mc := baseConfig
		mc.Self = id

		queueSlot := aleatypes.QueueSlot(idx)

		params := baseParams
		params.InstanceUID = bccommon.VCBInstanceUID(queueParams.BcInstanceUID, queueParams.QueueIdx, queueSlot)

		mod := vcb.NewModule(mc, params, nodeID, logging.Decorate(logger, "Vcb: ", "slot", idx))

		return mod, events.ListOf(
			bcqueuepbevents.BcStarted(queueMc.Consumer, &bcpbtypes.Slot{
				QueueIdx:  queueParams.QueueIdx,
				QueueSlot: queueSlot,
			}),
		), nil
	}
}
