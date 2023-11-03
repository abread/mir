package director

import (
	"fmt"
	"math"
	"time"

	"golang.org/x/exp/slices"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	"github.com/filecoin-project/mir/pkg/alea/queueselectionpolicy"
	"github.com/filecoin-project/mir/pkg/checkpoint"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	aagdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/dsl"
	bcpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	directorpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/dsl"
	directorpbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/events"
	directorpbmsgs "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/msgs"
	apppbdsl "github.com/filecoin-project/mir/pkg/pb/apppb/dsl"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	checkpointpbdsl "github.com/filecoin-project/mir/pkg/pb/checkpointpb/dsl"
	checkpointpbtypes "github.com/filecoin-project/mir/pkg/pb/checkpointpb/types"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	factorypbdsl "github.com/filecoin-project/mir/pkg/pb/factorypb/dsl"
	factorypbtypes "github.com/filecoin-project/mir/pkg/pb/factorypb/types"
	isspbdsl "github.com/filecoin-project/mir/pkg/pb/isspb/dsl"
	reliablenetpbdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	threshcheckpointpbdsl "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb/dsl"
	threshchkpvalidatorpbdsl "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb/threshchkpvalidatorpb/dsl"
	threshcheckpointpbtypes "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb/types"
	transportpbdsl "github.com/filecoin-project/mir/pkg/pb/transportpb/dsl"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	timert "github.com/filecoin-project/mir/pkg/timer/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

var timeRef = time.Now()

type set[T comparable] map[T]struct{}

type state struct {
	stalledBatchCut bool
	bcOwnQueueHead  aleatypes.QueueSlot

	queueSelectionPolicy queueselectionpolicy.QueueSelectionPolicy

	slotsReadyToDeliver set[bcpbtypes.Slot]
	agRound             uint64
	stalledAgRound      bool
	lastAgInput         time.Duration

	minAgRound            uint64
	nodeEpochMap          map[t.NodeID]tt.EpochNr
	lastStableCheckpoint  *threshcheckpointpbtypes.StableCheckpoint
	helpedNodes           map[t.NodeID]struct{}
	liveStableCheckpoints map[tt.SeqNr]struct{}

	nextCoalescedTimerDuration time.Duration
	lastWakeUp                 time.Duration
	lastScheduledWakeup        time.Duration
}

func NewModule( // nolint: gocyclo,gocognit
	mc ModuleConfig,
	params ModuleParams,
	tunables ModuleTunables,
	startingChkp *checkpoint.StableCheckpoint,
	nodeID t.NodeID,
	logger logging.Logger,
) (dsl.Module, error) {
	if startingChkp == nil {
		return nil, es.Errorf("missing initial checkpoint")
	}
	initialQsp, err := queueselectionpolicy.QueuePolicyFromBytes(startingChkp.Snapshot.EpochData.LeaderPolicy)
	if err != nil {
		return nil, err
	}

	m := dsl.NewModule(mc.Self)
	allNodes := maputil.GetSortedKeys(params.Membership.Nodes)
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(allNodes, nodeID))

	est := newEstimators(m, params, tunables, nodeID)
	state := newState(params, tunables, initialQsp)

	N := len(allNodes)
	F := (N - 1) / 3

	dsl.UponInit(m, func() error {
		return nil
	})

	directorpbdsl.UponHeartbeat(m, func() error {
		// UponStateUpdate(s) code will run
		return nil
	})

	// TODO: split sections into different parts (with independent state)

	// =============================================================================================
	// Delivery
	// =============================================================================================

	bcpbdsl.UponDeliverCert(m, func(cert *bcpbtypes.Cert) error {
		slot := *cert.Slot

		if !state.queueSelectionPolicy.SlotDelivered(slot) {
			// slot wasn't delivered yet by agreement component
			// logger.Log(logging.LevelDebug, "marking slot as ready for delivery", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
			state.slotsReadyToDeliver[slot] = struct{}{}
		} else {
			bcpbdsl.FreeSlot(m, mc.AleaBroadcast, cert.Slot)
		}

		return nil
	})

	// upon agreement round completion, deliver if it was decided to do so
	// if not, deliver empty batch (nilCert is used because we cannot currently)
	aagdsl.UponDeliver(m, func(round uint64, decision bool, _duration time.Duration, _posQuorumWait time.Duration) error {
		if round < state.minAgRound {
			// stale agreement round, from older epoch
			return nil
		}

		// pop queue
		slot, err := state.queueSelectionPolicy.DeliverSn(round, decision)
		if err != nil {
			return es.Errorf("could not deliver round %d: %w", round, err)
		}

		if decision {
			// agreement component delivers in-order, so we can deliver immediately
			logger.Log(logging.LevelDebug, "delivering cert", "agreementRound", round, "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
			isspbdsl.DeliverCert(m, mc.BatchFetcher, tt.SeqNr(round), &availabilitypbtypes.Cert{
				Type: &availabilitypbtypes.Cert_Alea{
					Alea: &bcpbtypes.Cert{
						Slot: &slot,
					},
				},
			}, false)

			// remove tracked slot readiness (don't want to run out of memory)
			// also free broadcast slot to allow broadcast component to make progress
			delete(state.slotsReadyToDeliver, slot)
			bcqueuepbdsl.FreeSlot(m, bccommon.BcQueueModuleID(mc.AleaBroadcast, slot.QueueIdx), slot.QueueSlot)
		} else {
			// deliver empty batch
			isspbdsl.DeliverCert(m, mc.BatchFetcher, tt.SeqNr(round), nil, true)
			logger.Log(logging.LevelDebug, "delivering empty cert", "agreementRound", round)
			return nil
		}

		return nil
	})

	// =============================================================================================
	// Agreement Round Control
	// =============================================================================================

	// upon init, stall agreement until a slot is deliverable
	dsl.UponInit(m, func() error {
		state.stalledAgRound = true
		return nil
	})

	// eagerly input 1 into rounds as soon as possible
	bcpbdsl.UponDeliverCert(m, func(cert *bcpbtypes.Cert) error {
		// If we delivered the next bc slot to be ordered, we can vote 1 for it in ag.
		// Note: we can't vote for slots delivered further ahead: we do not know in which ag round they
		// will be voted
		if !state.queueSelectionPolicy.SlotInQueueHead(*cert.Slot) {
			return nil
		}

		currentRoundQueueIdx := state.queueSelectionPolicy.NextSlot().QueueIdx
		if currentRoundQueueIdx == cert.Slot.QueueIdx {
			// this slot is for the current ag round: we need to be carefult not to input a value twice
			nextRound := state.agRound
			if state.stalledAgRound {
				logger.Log(logging.LevelDebug, "INPUT AG (BC-current)", "round", nextRound, "value", true)
				aagdsl.InputValue(m, mc.AleaAgreement, nextRound, true)
				state.lastAgInput = time.Since(timeRef)
				state.stalledAgRound = false
			}
		} else {
			// this slot is not for the current ag round: we can freely input to it
			// TODO: move this logic to queue selection policy
			var nextRound uint64
			if currentRoundQueueIdx < cert.Slot.QueueIdx {
				// we need to go a few rounds further
				nextRound = state.agRound + uint64(cert.Slot.QueueIdx-currentRoundQueueIdx)
			} else {
				nextRound = state.agRound + uint64(N) - uint64(currentRoundQueueIdx-cert.Slot.QueueIdx)
			}

			logger.Log(logging.LevelDebug, "INPUT AG (BC-future)", "round", nextRound, "value", true)
			aagdsl.InputValue(m, mc.AleaAgreement, nextRound, true)
			state.lastAgInput = time.Since(timeRef)
		}

		return nil
	})

	// Slots may be broadcast long before agreement delivers them: this event handler provides input to ag
	// for those slots that we left behind in the previous handler.
	aagdsl.UponDeliver(m, func(round uint64, _ bool, _, _ time.Duration) error {
		if round < state.minAgRound {
			// stale agreement round, from older epoch
			return nil
		}

		// if we delivered a new slot in a queue, we can input one for the next slot in the same queue
		// TODO: move this logic to queue selection policy
		nextRoundSameQueue := round + uint64(N)

		slot, ok := state.queueSelectionPolicy.Slot(nextRoundSameQueue)
		if !ok {
			return es.Errorf("could not find slot for round %d", nextRoundSameQueue)
		}

		if _, present := state.slotsReadyToDeliver[slot]; present {
			logger.Log(logging.LevelDebug, "INPUT AG (AG-done)", "round", nextRoundSameQueue, "value", true)
			aagdsl.InputValue(m, mc.AleaAgreement, nextRoundSameQueue, true)
			state.lastAgInput = time.Since(timeRef)

			bcpbdsl.MarkStableProposal(m, mc.AleaBroadcast, &slot)
		}
		return nil
	})

	// upon init, stall agreement until a slot is deliverable
	dsl.UponInit(m, func() error {
		state.stalledAgRound = true
		return nil
	})

	aagdsl.UponDeliver(m, func(round uint64, _ bool, _, _ time.Duration) error {
		if round < state.minAgRound {
			// stale agreement round, from older epoch
			return nil
		}

		// advance to next round
		state.agRound++

		// if bc has delivered the slot for this round already, then we have already input 1 to the
		// round, and it is not stalled at all.
		nextSlot, ok := state.queueSelectionPolicy.Slot(state.agRound)
		if !ok {
			return es.Errorf("bad queue selection policy: current ag round has no assigned slot")
		}
		if _, present := state.slotsReadyToDeliver[nextSlot]; present {
			state.stalledAgRound = false
		} else {
			state.stalledAgRound = true
		}

		return nil
	})

	agCanDeliverCount := func() int {
		nCanDeliver := 0

		for round := state.agRound; round < state.agRound+uint64(N); round++ {
			slot, ok := state.queueSelectionPolicy.Slot(round)
			if !ok {
				continue
			}

			if _, bcDone := state.slotsReadyToDeliver[slot]; bcDone {
				nCanDeliver++
			}
		}

		return nCanDeliver
	}

	agCanDeliverK := func(k int) bool {
		nCanDeliver := 0

		for round := state.agRound; round < state.agRound+uint64(N); round++ {
			slot, ok := state.queueSelectionPolicy.Slot(round)
			if !ok {
				continue
			}

			if _, bcDone := state.slotsReadyToDeliver[slot]; bcDone {
				nCanDeliver++
				if nCanDeliver >= k {
					return true
				}
			}
		}

		return nCanDeliver >= k
	}

	// if no agreement round is running, and any queue is able to deliver in this node, start the next round
	// a queue is able to deliver in this node if its head has been broadcast to it
	dsl.UponStateUpdates(m, func() error {
		if !state.stalledAgRound {
			return nil // nothing to do
		} else if agCanDeliverCount() == 0 && tunables.MaxAgStall > (time.Since(timeRef)-state.lastAgInput) {
			// just continue stalling for a while more: we don't have anything to deliver yet
			state.wakeUpAfter(tunables.MaxAgStall - (time.Since(timeRef) - state.lastAgInput))
			return nil
		}

		nextSlot, ok := state.queueSelectionPolicy.Slot(state.agRound)
		if !ok {
			return es.Errorf("bad queue selection policy: current ag round has no assigned slot")
		}

		// entering this code path means nextSlot was not delivered yet, and we should input 0 to the
		// agreement round

		// delay inputting 0 when a broadcast is in progress
		if bcRuntime, ok := est.BcRuntime(nextSlot); ok {
			if nextSlot.QueueIdx == ownQueueIdx && bcRuntime < est.OwnBcLocalMaxDurationEst() && bcRuntime < tunables.MaxAgStall {
				//logger.Log(logging.LevelDebug, "stalling agreement input for own batch")
				state.wakeUpAfter(tunables.MaxAgStall - bcRuntime)
				return nil
			}

			maxTimeToWait := est.ExtBcMaxDurationEst() - bcRuntime

			if maxTimeToWait > 0 {
				// stall agreement to allow in-flight broadcast to complete

				// schedule a timer to guarantee we reprocess the previous conditions
				// and eventually let agreement make progress, even if this broadcast
				// stalls indefinitely
				// logger.Log(logging.LevelDebug, "stalling agreement input", "maxDelay", maxTimeToWait)
				state.wakeUpAfter(maxTimeToWait)

				return nil
			}
		}

		// logger.Log(logging.LevelDebug, "progressing to next agreement round", "agreementRound", state.agRound, "input", bcDone)

		logger.Log(logging.LevelDebug, "INPUT AG (timeout)", "round", state.agRound, "value", false)
		aagdsl.InputValue(m, mc.AleaAgreement, state.agRound, false)
		state.lastAgInput = time.Since(timeRef)
		state.stalledAgRound = false

		return nil
	})

	// =============================================================================================
	// Batch Cutting / Own Queue Broadcast Control
	// =============================================================================================

	// upon init, cut a new batch
	dsl.UponInit(m, func() error {
		state.stalledBatchCut = false
		est.MarkBcStartedNow(bcpbtypes.Slot{QueueIdx: ownQueueIdx, QueueSlot: 0})

		bcpbdsl.RequestCert(m, mc.AleaBroadcast)
		return nil
	})

	// upon nice condition (unagreed batch count < max, no batch being cut, timeToNextAgForThisNode < estBc+margin || stalled ag), cut a new batch and broadcast it
	// TODO: move to bc component
	dsl.UponStateUpdates(m, func() error {
		if !state.stalledBatchCut {
			// no need to cut another batch while another is already being cut/broadcast
			return nil
		}

		if state.queueSelectionPolicy.QueueHead(ownQueueIdx) > state.bcOwnQueueHead {
			return es.Errorf("inconsistency: queue policy is ahead of bcOwnQueueHead")
		}

		// bcOwnQueueHead is the next slot to be broadcast
		// agQueueHeads[ownQueueIdx] is the next slot to be agreed on
		unagreedOwnBatchCount := uint64(state.bcOwnQueueHead - state.queueSelectionPolicy.QueueHead(ownQueueIdx))

		if unagreedOwnBatchCount >= uint64(tunables.MaxOwnUnagreedBatchCount) {
			// enough are cut already
			return nil
		}

		waitRoundCount := int(ownQueueIdx) - int(state.agRound%uint64(N)) - 1
		if waitRoundCount < 0 {
			waitRoundCount += N
		}

		// consider how many batches we have pending delivery
		waitRoundCount += int(unagreedOwnBatchCount) * N

		// consider how many batches can be instantly delivered
		// TODO: may increase tx duplication when clients multicast txs and byz nodes exist
		waitRoundCount -= agCanDeliverCount()

		agRoundTime := est.AgFastPathEst()
		timeToOwnQueueAgRound := agRoundTime * time.Duration(waitRoundCount)
		bcRuntimeEst := est.OwnBcMaxDurationEst()

		// We have a lot of time before we reach our agreement round. Let the batch fill up!
		// That said, if we have no unagreed batches from any correct replica, we'll fast track it.
		if timeToOwnQueueAgRound > bcRuntimeEst && (unagreedOwnBatchCount > 0 || agCanDeliverK(F+1)) {
			// ensure we are woken up to create a batch before we run out of time
			maxDelay := timeToOwnQueueAgRound - bcRuntimeEst

			if maxDelay > agRoundTime { // don't delay batch cut for very short periods of time
				// logger.Log(logging.LevelDebug, "stalling batch cut", "max delay", maxDelay)
				state.wakeUpAfter(maxDelay)

				return nil
			}
		}

		logger.Log(logging.LevelDebug, "requesting more transactions")
		state.stalledBatchCut = false

		est.MarkBcStartedNow(bcpbtypes.Slot{QueueIdx: ownQueueIdx, QueueSlot: state.bcOwnQueueHead})

		bcpbdsl.RequestCert(m, mc.AleaBroadcast)
		return nil
	})
	bcpbdsl.UponDeliverCert(m, func(cert *bcpbtypes.Cert) error {
		if cert.Slot.QueueIdx == ownQueueIdx && cert.Slot.QueueSlot == state.bcOwnQueueHead {
			// new batch was delivered
			state.stalledBatchCut = true
			state.bcOwnQueueHead++
		}

		return nil
	})

	memberships := make([]*trantorpbtypes.Membership, 1)
	memberships[0] = params.Membership

	// =============================================================================================
	// Checkpointing
	// =============================================================================================

	// checkpoint every <epoch length> rounds
	advanceEpoch := func(newEpoch tt.EpochNr) {
		// inform components of the new epoch
		for _, mod := range []t.ModuleID{mc.AleaAgreement, mc.AleaBroadcast} {
			directorpbdsl.NewEpoch(m, mod, newEpoch)
		}
		apppbdsl.NewEpoch(m, mc.BatchFetcher, newEpoch, mc.Self)
		logger.Log(logging.LevelInfo, "Advanced to new epoch", "epoch", newEpoch)
	}
	checkpointAndAdvanceEpoch := func(newEpoch tt.EpochNr) error {
		advanceEpoch(newEpoch)
		nextSeqNr := tt.SeqNr(uint64(newEpoch) * params.EpochLength)

		logger.Log(logging.LevelInfo, "Starting checkpoint", "epoch", newEpoch)
		// Create a new checkpoint tracker to start the checkpointing protocol.
		// nextDeliveredSN is the first sequence number *not* included in the checkpoint,
		// i.e., as sequence numbers start at 0, the checkpoint includes the first nextDeliveredSN sequence numbers.
		// TODO: support membership changes
		chkpModuleID := mc.Checkpoint.Then(t.NewModuleIDFromInt(newEpoch))
		serializedPolicy, err := state.queueSelectionPolicy.Bytes()
		if err != nil {
			return err
		}
		factorypbdsl.NewModule(m,
			mc.Checkpoint,
			chkpModuleID,
			tt.RetentionIndex(newEpoch),
			&factorypbtypes.GeneratorParams{
				Type: &factorypbtypes.GeneratorParams_ThreshCheckpoint{
					ThreshCheckpoint: &threshcheckpointpbtypes.InstanceParams{
						Membership:       params.Membership,
						LeaderPolicyData: serializedPolicy,
						EpochConfig: &trantorpbtypes.EpochConfig{
							EpochNr:     newEpoch,
							FirstSn:     nextSeqNr,
							Length:      params.EpochLength,
							Memberships: memberships,
						},
						Threshold: uint64(2*F + 1),
					},
				},
			},
		)

		// Ask the application for a state snapshot and have it send the result directly to the checkpoint module.
		// Note that the new instance of the checkpoint protocol is not yet created at this moment,
		// but it is guaranteed to be created before the application's response.
		// This is because the NewModule event will already be enqueued for the checkpoint factory
		// when the application receives the snapshot request.
		apppbdsl.SnapshotRequest(m, mc.BatchFetcher, chkpModuleID)

		state.liveStableCheckpoints[nextSeqNr] = struct{}{}

		return nil
	}
	aagdsl.UponDeliver(m, func(round uint64, _ bool, _, _ time.Duration) error {
		if round < state.minAgRound {
			// stale agreement round, from older epoch
			return nil
		}

		state.minAgRound = round + 1

		if (round+1)%params.EpochLength == 0 {
			epochNr := tt.EpochNr(round / params.EpochLength)
			return checkpointAndAdvanceEpoch(epochNr + 1)
		}

		return nil
	})
	saveLatestStableCheckpoint := func(checkpoint *threshcheckpointpbtypes.StableCheckpoint) {
		state.lastStableCheckpoint = checkpoint
		logger.Log(logging.LevelInfo, "Updated stable checkpoint", "epoch", checkpoint.Snapshot.EpochData.EpochConfig.EpochNr)

		epochNr := checkpoint.Snapshot.EpochData.EpochConfig.EpochNr

		if uint64(epochNr) > params.RetainEpochs+1 {
			firstRetainedEpoch := epochNr - tt.EpochNr(params.RetainEpochs)

			// inform components of the checkpoint stabilization to clear old data
			for _, mod := range []t.ModuleID{mc.AleaAgreement, mc.AleaBroadcast} {
				directorpbdsl.GCEpochs(m, mod, firstRetainedEpoch)
			}

			// cleanup old checkpointing instances
			factorypbdsl.GarbageCollect(m, mc.Checkpoint, tt.RetentionIndex(firstRetainedEpoch))

			// ensure we don't retain stale checkpoint messages as well
			for sn := range state.liveStableCheckpoints {
				if sn < tt.SeqNr(uint64(firstRetainedEpoch)*params.EpochLength) {
					chkpEpochNr := uint64(sn) / params.EpochLength

					// TODO: make it work with membership changes
					reliablenetpbdsl.MarkModuleMsgsRecvd(m, mc.ReliableNet, mc.Checkpoint.Then(t.NewModuleIDFromInt(chkpEpochNr)), allNodes)
					delete(state.liveStableCheckpoints, sn)
				}
			}
		}
	}
	threshcheckpointpbdsl.UponStableCheckpoint(m, func(sn tt.SeqNr, snapshot *trantorpbtypes.StateSnapshot, signature tctypes.FullSig) error {
		// we got it, no longer pending
		delete(state.liveStableCheckpoints, sn)

		if state.lastStableCheckpoint == nil || sn > state.lastStableCheckpoint.Sn {
			saveLatestStableCheckpoint(&threshcheckpointpbtypes.StableCheckpoint{
				Sn:        sn,
				Snapshot:  snapshot,
				Signature: signature,
			})

			// preemptively help nodes that are very far behind, to ensure liveness
			currentEpochNr := state.lastStableCheckpoint.Snapshot.EpochData.EpochConfig.EpochNr
			for nodeID, epochNr := range state.nodeEpochMap {
				// Note: we add the current node idx to avoid overloading the remote node with snapshots
				if epochNr+tt.EpochNr(params.RetainEpochs)+tt.EpochNr(ownQueueIdx) < currentEpochNr {
					directorpbdsl.HelpNode(m, mc.Self, nodeID)
					transportpbdsl.ForceSendMessage(m, mc.Net,
						directorpbmsgs.StableCheckpoint(
							mc.Self,
							state.lastStableCheckpoint,
						),
						[]t.NodeID{nodeID},
					)
				}
			}
		}

		return nil
	})
	isspbdsl.UponNewConfig(m, func(epochNr tt.EpochNr, membership *trantorpbtypes.Membership) error {
		// TODO: support membership changes
		return nil
	})

	// help stale replicas catch up, when hinted by other components
	checkpointpbdsl.UponEpochProgress(m, func(nodeId t.NodeID, epoch tt.EpochNr) error {
		if epoch > state.nodeEpochMap[nodeId] {
			state.nodeEpochMap[nodeId] = epoch

			// TODO: don't send events to reliablenet if the node was not helped (volume should be low anyway)
			reliablenetpbdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, rntypes.MsgID(fmt.Sprintf("s%d", epoch)), []t.NodeID{nodeId})
		}
		return nil
	})
	directorpbdsl.UponHelpNode(m, func(nodeID t.NodeID) error {
		//logger.Log(logging.LevelDebug, "wanted to help", "node", nodeID)
		if state.lastStableCheckpoint == nil {
			// can't help if we don't have a checkpoint yet
			return nil
		}

		epochNr := state.lastStableCheckpoint.Snapshot.EpochData.EpochConfig.EpochNr
		if state.nodeEpochMap[nodeID] < epochNr+tt.EpochNr(params.RetainEpochs) {
			if _, ok := state.helpedNodes[nodeID]; ok {
				return nil // already helped
			}
			state.nodeEpochMap[nodeID] = epochNr
			state.helpedNodes[nodeID] = struct{}{}

			logger.Log(logging.LevelDebug, "Helping node catch up", "node", nodeID, "chkpEpoch", epochNr, "theirEpoch", state.nodeEpochMap[nodeID])
			reliablenetpbdsl.ForceSendMessage(m, mc.ReliableNet, rntypes.MsgID(fmt.Sprintf("s%d", epochNr)),
				directorpbmsgs.StableCheckpoint(
					mc.Self,
					state.lastStableCheckpoint,
				),
				[]t.NodeID{nodeID},
			)

			// clear previous help attempts
			reliablenetpbdsl.MarkModuleMsgsRecvd(m, mc.ReliableNet, mc.Self, []t.NodeID{nodeID})
		}
		return nil
	})

	// allow help by others
	directorpbdsl.UponStableCheckpointReceived(m, func(from t.NodeID, checkpoint *threshcheckpointpbtypes.StableCheckpoint) error {
		epochNr := checkpoint.Snapshot.EpochData.EpochConfig.EpochNr
		if epochNr > state.nodeEpochMap[from] {
			// even if the node is byz and this is bs, it's fine, we'll just not help it as much
			state.nodeEpochMap[from] = epochNr
		}

		if checkpoint.Sn < tt.SeqNr(state.agRound) {
			// stale checkpoint
			return nil
		}

		// TODO: use go 1.21's clear method
		state.helpedNodes = make(map[t.NodeID]struct{})

		threshchkpvalidatorpbdsl.ValidateCheckpoint(m, mc.ChkpValidator, checkpoint, 0, memberships, checkpoint)
		return nil
	})
	threshchkpvalidatorpbdsl.UponCheckpointValidated(m, func(err error, checkpoint *threshcheckpointpbtypes.StableCheckpoint) error {
		if err != nil {
			// bad sig on checkpoint
			// TODO: report byz node
			logger.Log(logging.LevelError, "bad checkpoint signature", "err", err)
			return nil
		}

		if checkpoint.Sn <= tt.SeqNr(state.agRound) {
			// stale checkpoint
			return nil
		}

		logger.Log(logging.LevelWarn, "Fell behind. Restoring from received checkpoint.", "chkpSn", checkpoint.Sn, "ourSn", state.agRound)

		// HACK: make a threshcheckpoint into a checkpoint, deterministically
		regularCheckpoint := &checkpointpbtypes.StableCheckpoint{
			Sn:       checkpoint.Sn,
			Snapshot: checkpoint.Snapshot,
			Cert: map[t.NodeID][]uint8{
				allNodes[0]: checkpoint.Signature,
			},
		}

		// apply data from checkpoint across all relevant components
		for _, mod := range []t.ModuleID{mc.BatchFetcher, mc.AleaAgreement, mc.AleaBroadcast} {
			apppbdsl.RestoreState(m, mod, regularCheckpoint)
		}

		// also apply it to the director
		var errQs error
		state.queueSelectionPolicy, errQs = queueselectionpolicy.QueuePolicyFromBytes(checkpoint.Snapshot.EpochData.LeaderPolicy)
		if errQs != nil {
			return errQs
		}

		// current round is the checkpoint round, and it starts waiting for input
		state.minAgRound = uint64(checkpoint.Sn)
		state.agRound = uint64(checkpoint.Sn)
		state.stalledAgRound = true

		// clear slots delivered indirectly through the checkpoint
		for slot := range state.slotsReadyToDeliver {
			if state.queueSelectionPolicy.SlotDelivered(slot) {
				delete(state.slotsReadyToDeliver, slot)
			}
		}

		// then advance to new epoch, and store the received checkpoint
		advanceEpoch(checkpoint.Snapshot.EpochData.EpochConfig.EpochNr + 1)
		saveLatestStableCheckpoint(checkpoint)

		// note: no need to create a new checkpoint module, since the checkpoint was already created and validated
		return nil
	})

	// Schedule coalesced timer
	// MUST BE THE LAST UponStateUpdates HANDLER
	dsl.UponStateUpdates(m, func() error {
		now := time.Since(timeRef)

		// only schedule timer if needed (request AND there is no timer scheduled already for somewhere before this period)
		// Note: when now+d==lastScheduledWakeup, we must schedule a timer, since monotonic clocks may return the same time twice.
		d := state.nextCoalescedTimerDuration
		if d != math.MaxInt64 && now+d >= state.lastScheduledWakeup {
			eventpbdsl.TimerDelay(m, mc.Timer,
				[]*eventpbtypes.Event{
					directorpbevents.Heartbeat(mc.Self),
				},
				timert.Duration(d),
			)

			state.lastScheduledWakeup = now + d
		}

		state.lastWakeUp = now

		// clear coalesced timer for next batch of events
		state.nextCoalescedTimerDuration = math.MaxInt64

		return nil
	})

	return m, nil
}

func (state *state) wakeUpAfter(d time.Duration) {
	if d < state.nextCoalescedTimerDuration {
		state.nextCoalescedTimerDuration = d
	}
}

func newState(params ModuleParams, tunables ModuleTunables, qsp queueselectionpolicy.QueueSelectionPolicy) *state {
	N := len(params.Membership.Nodes)

	state := &state{
		bcOwnQueueHead:        0,
		queueSelectionPolicy:  qsp,
		slotsReadyToDeliver:   make(set[bcpbtypes.Slot], N*tunables.MaxConcurrentVcbPerQueue),
		nodeEpochMap:          make(map[t.NodeID]tt.EpochNr, N),
		helpedNodes:           make(map[t.NodeID]struct{}),
		liveStableCheckpoints: make(map[tt.SeqNr]struct{}),

		nextCoalescedTimerDuration: math.MaxInt64,
	}

	for nodeID := range params.Membership.Nodes {
		state.nodeEpochMap[nodeID] = 0
	}

	return state
}
