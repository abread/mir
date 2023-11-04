package bcqueue

import (
	"math"
	"time"

	es "github.com/go-errors/errors"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/abba/abbatypes"
	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	"github.com/filecoin-project/mir/pkg/alea/queueselectionpolicy"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modring"
	"github.com/filecoin-project/mir/pkg/modules"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	bcqueuepbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/events"
	directorpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/dsl"
	apppbdsl "github.com/filecoin-project/mir/pkg/pb/apppb/dsl"
	checkpointpbtypes "github.com/filecoin-project/mir/pkg/pb/checkpointpb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	reliablenetpbevents "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/events"
	transportpbdsl "github.com/filecoin-project/mir/pkg/pb/transportpb/dsl"
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
		Generator: newVcbGenerator(mc, &params, nodeID, logger),
		CleanupHandler: func(slot uint64) (events.EventList, error) {
			return events.ListOf(
				reliablenetpbevents.MarkModuleMsgsRecvd(mc.ReliableNet, mc.Self.Then(t.NewModuleIDFromInt(slot)), params.AllNodes),
			), nil
		},
		PastMsgHandler: newPastMsgHandler(mc, params),
	}, logging.Decorate(logger, "Modring controller: "))

	controller := newQueueController(mc, &params, tunables, nodeID, logger, slots)

	var children modules.PassiveModule = slots
	if params.QueueOwner == nodeID {
		// instrument VCB for network latency estimates
		children = modules.MultiApplyModule([]modules.PassiveModule{
			slots,
			newMinNetworkLatEstimator(mc, &params, nodeID),
		})
	}

	return modules.RoutedModule(mc.Self, controller, children), nil
}

func newQueueController(mc ModuleConfig, params *ModuleParams, tunables ModuleTunables, nodeID t.NodeID, logger logging.Logger, slots *modring.Module) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	parseSlotFromModuleID := func(id t.ModuleID) (bcpbtypes.Slot, error) {
		queueSlot, err := mc.ParseQueueSlotFromModuleID(id)
		if err != nil {
			return bcpbtypes.Slot{}, err
		}

		return bcpbtypes.Slot{
			QueueIdx:  params.QueueIdx,
			QueueSlot: queueSlot,
		}, nil
	}

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
		slot, err := parseSlotFromModuleID(srcModule)
		if err != nil {
			return es.Errorf("deliver event for invalid round: %w", err)
		}

		logger.Log(logging.LevelDebug, "delivering broadcast", "queueSlot", slot.QueueSlot)
		bcqueuepbdsl.Deliver(m, mc.Consumer, &bcpbtypes.Cert{
			Slot:      &slot,
			BatchId:   batchID,
			Signature: signature,
		})
		return nil
	})

	if params.QueueOwner == nodeID {
		slotDeliverTimes := make(map[aleatypes.QueueSlot]time.Time, tunables.MaxConcurrentVcb)

		vcbpbdsl.UponDeliver(m, func(_ string, _ tctypes.FullSig, srcModule t.ModuleID) error {
			queueSlot, err := mc.ParseQueueSlotFromModuleID(srcModule)
			if err != nil {
				return es.Errorf("deliver event for invalid round: %w", err)
			}

			slotDeliverTimes[aleatypes.QueueSlot(queueSlot)] = time.Now()
			return nil
		})

		vcbpbdsl.UponQuorumDone(m, func(srcModule t.ModuleID) error {
			slot, err := parseSlotFromModuleID(srcModule)
			if err != nil {
				return es.Errorf("deliver event for invalid round: %w", err)
			}

			if startTime, ok := slotDeliverTimes[slot.QueueSlot]; ok {
				deliverDelta := time.Since(startTime)
				bcqueuepbdsl.BcQuorumDone(m, mc.Consumer, &slot, deliverDelta)

				// add delta to avoid including it in the fully done delta calc
				slotDeliverTimes[slot.QueueSlot] = startTime.Add(deliverDelta)
			}

			return nil
		})

		vcbpbdsl.UponAllDone(m, func(srcModule t.ModuleID) error {
			slot, err := parseSlotFromModuleID(srcModule)
			if err != nil {
				return es.Errorf("deliver event for invalid round: %w", err)
			}

			if startTime, ok := slotDeliverTimes[slot.QueueSlot]; ok {
				deliverDelta := time.Since(startTime)
				bcqueuepbdsl.BcAllDone(m, mc.Consumer, &slot, deliverDelta)
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

		evsOut, err := slots.MarkSubmodulePast(uint64(queueSlot))
		if err != nil {
			return es.Errorf("failed to free queue slot: %w", err)
		}
		dsl.EmitEvents(m, evsOut)
		return nil
	})

	freeSlotsBefore := func(minSlot aleatypes.QueueSlot) error {
		if err := slots.AdvanceViewToAtLeastSubmodule(uint64(minSlot)); err != nil {
			return err
		}
		evsOut, err := slots.FreePast()
		if err != nil {
			return err
		}
		dsl.EmitEvents(m, evsOut)

		return nil
	}

	// track last freed slot in each epoch to properly garbage collect them
	// note that slots are only freed after delivery
	var currentEpoch tt.EpochNr
	epochLastFreedSlot := make(map[tt.EpochNr]aleatypes.QueueSlot)
	epochLastFreedSlot[0] = aleatypes.QueueSlot(math.MaxUint64) // ensure first epoch is always initialized

	directorpbdsl.UponNewEpoch(m, func(epochNr tt.EpochNr) error {
		currentEpoch = epochNr
		epochLastFreedSlot[currentEpoch] = aleatypes.QueueSlot(math.MaxUint64)
		return nil
	})
	bcqueuepbdsl.UponFreeSlot(m, func(queueSlot aleatypes.QueueSlot) error {
		if epochLastFreedSlot[currentEpoch] < queueSlot || epochLastFreedSlot[currentEpoch] == math.MaxUint64 {
			epochLastFreedSlot[currentEpoch] = queueSlot
		}
		return nil
	})

	// garbage collect slots that are not needed anymore
	directorpbdsl.UponGCEpochs(m, func(minEpoch tt.EpochNr) error {
		maxFreedSlot := aleatypes.QueueSlot(math.MaxUint64)

		for epoch := range epochLastFreedSlot {
			if epoch >= minEpoch {
				continue // keep epoch
			}

			if epochLastFreedSlot[epoch] != math.MaxUint64 {
				// something was freed in this epoch
				if epochLastFreedSlot[epoch] > maxFreedSlot || maxFreedSlot == math.MaxUint64 {
					maxFreedSlot = epochLastFreedSlot[epoch]
				}
			}

			delete(epochLastFreedSlot, epoch)
		}

		if maxFreedSlot != math.MaxUint64 {
			// forget slots delivered in old epochs
			return freeSlotsBefore(maxFreedSlot + 1)
		}

		return nil
	})

	// when restoring from checkpoint, free all slots before the checkpoint to avoid tracking useless slots
	apppbdsl.UponRestoreState(m, func(checkpoint *checkpointpbtypes.StableCheckpoint) error {
		if params.QueueOwner == nodeID {
			// own queue does not need to restore state from the outside world
			return nil
		}

		qsp, err := queueselectionpolicy.QueuePolicyFromBytes(checkpoint.Snapshot.EpochData.LeaderPolicy)
		if err != nil {
			return err
		}

		return freeSlotsBefore(qsp.QueueHead(params.QueueIdx))
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

func newVcbGenerator(queueMc ModuleConfig, queueParams *ModuleParams, nodeID t.NodeID, logger logging.Logger) func(id t.ModuleID, idx uint64) (modules.PassiveModule, events.EventList, error) {
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
		// retention index will be updated later when the batch is fetched
		EpochNr:  tt.RetentionIndex(math.MaxUint64),
		AllNodes: queueParams.AllNodes,
		Leader:   queueParams.QueueOwner,
	}

	return func(id t.ModuleID, idx uint64) (modules.PassiveModule, events.EventList, error) {
		mc := baseConfig
		mc.Self = id

		queueSlot := aleatypes.QueueSlot(idx)

		params := baseParams
		params.InstanceUID = bccommon.VCBInstanceUID(queueParams.AleaInstanceUID, queueParams.QueueIdx, queueSlot)

		mod := vcb.NewModule(mc, params, nodeID, logging.Decorate(logger, "Vcb: ", "slot", idx))

		return mod, events.ListOf(
			bcqueuepbevents.BcStarted(queueMc.Consumer, &bcpbtypes.Slot{
				QueueIdx:  queueParams.QueueIdx,
				QueueSlot: queueSlot,
			}),
		), nil
	}
}

type netLatEstVcbInfo struct {
	ownEchoTime time.Time
	recvdEchos  abbatypes.RecvTracker // TODO: move type into separate common package
}

func newMinNetworkLatEstimator(mc ModuleConfig, params *ModuleParams, nodeID t.NodeID) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	N := len(params.AllNodes)
	F := (N - 1) / 3
	measurementThreshold := F + 2 // own node + F byz nodes + first honest node

	state := make(map[aleatypes.QueueSlot]*netLatEstVcbInfo)

	bcqueuepbdsl.UponInputValue(m, func(queueSlot aleatypes.QueueSlot, txIds []tt.TxID, txs []*trantorpbtypes.Transaction) error {
		state[queueSlot] = &netLatEstVcbInfo{}
		return nil
	})

	transportpbdsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepbtypes.Message) error {
		if vcbMsgW, ok := msg.Type.(*messagepbtypes.Message_Vcb); ok {
			if _, isEcho := vcbMsgW.Vcb.Type.(*vcbpbtypes.Message_EchoMessage); isEcho {
				queueSlot, err := mc.ParseQueueSlotFromModuleID(msg.DestModule)
				if err != nil {
					return nil // message for invalid vcb slot
				}

				vcbInfo, ok := state[queueSlot]
				if !ok {
					return nil // slot was already processed
				}

				if !vcbInfo.recvdEchos.Register(from) {
					return nil // message already processed
				}

				if from == nodeID {
					// register time for own ECHO message
					vcbInfo.ownEchoTime = time.Now()
				}

				if vcbInfo.recvdEchos.Len() == measurementThreshold {
					var emptyTime time.Time

					// we have information for an estimate!
					var minEst time.Duration
					if vcbInfo.ownEchoTime != emptyTime {
						// the leader of a VCB instance send SEND immediately upon starting the broadcast
						// afterwards, it computes a sig share and sends it to itself via an ECHO message (with ~0 latency)
						// simultaneously, followers receive the SEND message, compute their sig share, and send it via an ECHO message
						// the difference between receiving a local ECHO and remote ECHOs should be 1 RTT (~2x the network latency)
						minEst = time.Since(vcbInfo.ownEchoTime) / 2
					} else {
						// network latency is so small and signatures so fast (or processing time so slow)
						// that the followers reply was processed first!
						minEst = 0
					}

					// propagate estimate and clean up state
					bcqueuepbdsl.NetLatencyEstimate(m, mc.Consumer, minEst)
					delete(state, queueSlot)
				}
			}
		}

		return nil
	})

	dsl.UponOtherEvent(m, func(ev *eventpbtypes.Event) error {
		return nil
	})

	return m
}
