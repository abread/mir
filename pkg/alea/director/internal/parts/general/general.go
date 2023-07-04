package general

import (
	"math"
	"time"

	"golang.org/x/exp/slices"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bcutil"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/parts/estimators"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	aagdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/dsl"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	directorpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/dsl"
	directorpbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/events"
	aleapbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/types"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	isspbdsl "github.com/filecoin-project/mir/pkg/pb/isspb/dsl"
	mempooldsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	timert "github.com/filecoin-project/mir/pkg/timer/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

var timeRef = time.Now()

type set[T comparable] map[T]struct{}

type state struct {
	stalledBatchCut bool
	bcOwnQueueHead  aleatypes.QueueSlot

	agQueueHeads []aleatypes.QueueSlot

	slotsReadyToDeliver []set[aleatypes.QueueSlot]
	agRound             uint64
	stalledAgRound      bool

	nextCoalescedTimerDuration time.Duration
	lastWakeUp                 time.Duration
	lastScheduledWakeup        time.Duration

	targetOwnUnagreedBatches int
	ownAgConsecutiveDelivers int
}

func Include(m dsl.Module, mc common.ModuleConfig, params common.ModuleParams, tunables common.ModuleTunables, nodeID t.NodeID, logger logging.Logger, est *estimators.Estimators) {
	state := newState(params, tunables, nodeID)
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(params.AllNodes, nodeID))

	N := len(params.AllNodes)
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

	bcqueuepbdsl.UponDeliver(m, func(slot *commontypes.Slot) error {
		if slot.QueueSlot >= state.agQueueHeads[slot.QueueIdx] {
			// slot wasn't delivered yet by agreement component
			// logger.Log(logging.LevelDebug, "marking slot as ready for delivery", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
			state.slotsReadyToDeliver[slot.QueueIdx][slot.QueueSlot] = struct{}{}
		} else {
			bcqueuepbdsl.FreeSlot(m, bcutil.BcQueueModuleID(mc.BcQueuePrefix, slot.QueueIdx), slot.QueueSlot)
		}

		return nil
	})

	// upon agreement round completion, deliver if it was decided to do so
	aagdsl.UponDeliver(m, func(round uint64, decision bool, _duration time.Duration, _posQuorumWait time.Duration) error {
		if !decision {
			// nothing to deliver
			return nil
		}

		queueIdx := aleatypes.QueueIdx(round % uint64(N))
		slot := &commontypes.Slot{
			QueueIdx:  queueIdx,
			QueueSlot: state.agQueueHeads[queueIdx],
		}

		// next round won't start until we say so, and previous rounds already delivered, so we can deliver immediately
		logger.Log(logging.LevelDebug, "delivering cert", "agreementRound", round, "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
		isspbdsl.DeliverCert(m, mc.Consumer, tt.SeqNr(round), &availabilitypbtypes.Cert{
			Type: &availabilitypbtypes.Cert_Alea{
				Alea: &aleapbtypes.Cert{
					Slot: slot,
				},
			},
		}, false)

		// pop queue
		state.agQueueHeads[queueIdx]++

		// remove tracked slot readiness (don't want to run out of memory)
		// also free broadcast slot to allow broadcast component to make progress
		delete(state.slotsReadyToDeliver[slot.QueueIdx], slot.QueueSlot)
		bcqueuepbdsl.FreeSlot(m, bcutil.BcQueueModuleID(mc.BcQueuePrefix, slot.QueueIdx), slot.QueueSlot)

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

	// upon agreement round completion, prepare next round
	aagdsl.UponDeliver(m, func(round uint64, decision bool, duration time.Duration, _posQuorumWait time.Duration) error {
		state.agRound++
		state.stalledAgRound = true
		return nil
	})

	// if no agreement round is running, and any queue is able to deliver in this node, start the next round
	// a queue is able to deliver in this node if its head has been broadcast to it
	dsl.UponStateUpdates(m, func() error {
		if !state.stalledAgRound {
			return nil // nothing to do
		}

		canDeliverSomething := state.agCanDeliver(1)

		nextQueueIdx := aleatypes.QueueIdx(state.agRound % uint64(N))
		nextQueueSlot := state.agQueueHeads[nextQueueIdx]

		if canDeliverSomething {
			_, bcDone := state.slotsReadyToDeliver[nextQueueIdx][nextQueueSlot]
			if !bcDone {
				slot := commontypes.Slot{
					QueueIdx:  nextQueueIdx,
					QueueSlot: nextQueueSlot,
				}

				if bcRuntime, ok := est.BcRuntime(slot); ok {
					if nextQueueIdx == ownQueueIdx && bcRuntime < est.OwnBcMedianDurationEstNoMargin() {
						logger.Log(logging.LevelDebug, "stalling agreement input for own batch")
						return nil
					}

					maxTimeToWait := est.ExtBcMaxDurationEst() - bcRuntime

					// clamp wait time just in case
					if maxTimeToWait > tunables.MaxAgreementDelay {
						maxTimeToWait = 0
					}

					if maxTimeToWait > 0 {
						// stall agreement to allow in-flight broadcast to complete

						// schedule a timer to guarantee we reprocess the previous conditions
						// and eventually let agreement make progress, even if this broadcast
						// stalls indefinitely
						logger.Log(logging.LevelDebug, "stalling agreement input", "maxDelay", maxTimeToWait)
						state.wakeUpAfter(maxTimeToWait)

						return nil
					}
				}
			}

			// logger.Log(logging.LevelDebug, "progressing to next agreement round", "agreementRound", state.agRound, "input", bcDone)

			aagdsl.InputValue(m, mc.AleaAgreement, state.agRound, bcDone)
			state.stalledAgRound = false
		}

		return nil
	})

	// =====
	// Unagreed batch count target control
	// =====
	aagdsl.UponDeliver(m, func(round uint64, decision bool, posQuorumWait, posTotalWait time.Duration) error {
		if round%uint64(N) != uint64(ownQueueIdx) {
			return nil // not own batch
		}

		if decision {
			state.ownAgConsecutiveDelivers++
		} else {
			state.ownAgConsecutiveDelivers = 0
			// agreement is going fast, broadcast needs to catch up
			state.targetOwnUnagreedBatches++
		}

		return nil
	})
	dsl.UponStateUpdates(m, func() error {
		if state.ownAgConsecutiveDelivers > 5 {
			// broadcast caught up, we can slow it down again
			state.targetOwnUnagreedBatches /= 2

			if state.targetOwnUnagreedBatches == 0 {
				state.targetOwnUnagreedBatches = 1
			}
		}

		if state.targetOwnUnagreedBatches > tunables.MaxOwnUnagreedBatchCount {
			state.targetOwnUnagreedBatches = tunables.MaxOwnUnagreedBatchCount
		}

		return nil
	})

	// =============================================================================================
	// Batch Cutting / Own Queue Broadcast Control
	// =============================================================================================

	// upon init, cut a new batch
	dsl.UponInit(m, func() error {
		state.stalledBatchCut = false
		mempooldsl.RequestBatch[struct{}](m, mc.Mempool, nil)

		return nil
	})

	// upon nice condition (unagreed batch count < max, no batch being cut, timeToNextAgForThisNode < estBc+margin || stalled ag), cut a new batch and broadcast it
	// TODO: move to bc component
	dsl.UponStateUpdates(m, func() error {
		// bcOwnQueueHead is the next slot to be broadcast
		// agQueueHeads[ownQueueIdx] is the next slot to be agreed on
		unagreedOwnBatchCount := uint64(state.bcOwnQueueHead - state.agQueueHeads[ownQueueIdx])

		if !state.stalledBatchCut || unagreedOwnBatchCount >= uint64(state.targetOwnUnagreedBatches) {
			// batch cut in progress, or enough are cut already
			return nil
		}

		waitRoundCount := int(ownQueueIdx) - int(state.agRound%uint64(N)) - 1
		if waitRoundCount < 0 {
			waitRoundCount += N
		}

		// consider how many batches we need to deliver in each wait period
		// this corrects for bad estimates
		waitRoundCount /= state.targetOwnUnagreedBatches

		// TODO: consider progress in current round too (will mean adjustments below)
		timeToOwnQueueAgRound := est.AgMinDurationEst() * time.Duration(waitRoundCount)
		bcRuntimeEst := est.OwnBcMaxDurationEst()

		// We have a lot of time before we reach our agreement round. Let the batch fill up!
		// We must also guarantee F+1 nodes have undelivered batches, or that agreement is currently progressing,
		// otherwise an attacker can stall the system by not sending their batch to enough nodes.
		// Additionally, we don't want to delay a batch indefinitely in the presence of bad estimates: if agreement for our slot fails, we took too long!
		if timeToOwnQueueAgRound > bcRuntimeEst && state.ownAgConsecutiveDelivers > 0 && (!state.stalledAgRound || state.agCanDeliver(F+1)) {
			// ensure we are woken up to create a batch before we run out of time
			maxDelay := timeToOwnQueueAgRound - bcRuntimeEst
			logger.Log(logging.LevelDebug, "stalling batch cut", "max delay", maxDelay)
			state.wakeUpAfter(maxDelay)

			return nil
		}

		// logger.Log(logging.LevelDebug, "requesting more transactions")
		state.stalledBatchCut = false

		mempooldsl.RequestBatch[struct{}](m, mc.Mempool, nil)
		return nil
	})
	mempooldsl.UponNewBatch(m, func(txIDs []tt.TxID, txs []*trantorpbtypes.Transaction, ctx *struct{}) error {
		if len(txs) == 0 {
			return es.Errorf("empty batch. did you misconfigure your mempool?")
		}

		// logger.Log(logging.LevelDebug, "new batch", "nTransactions", len(txs))

		bcqueuepbdsl.InputValue(m, bcutil.BcQueueModuleID(mc.BcQueuePrefix, ownQueueIdx), state.bcOwnQueueHead, txs)
		state.bcOwnQueueHead++
		return nil
	})
	bcqueuepbdsl.UponDeliver(m, func(slot *commontypes.Slot) error {
		if slot.QueueIdx == ownQueueIdx && slot.QueueSlot == state.bcOwnQueueHead-1 {
			// new batch was delivered
			state.stalledBatchCut = true
		}
		return nil
	})

	// Schedule coalesced timer
	// MUST BE THE LAST UponStateUpdates HANDLER
	dsl.UponStateUpdates(m, func() error {
		now := time.Since(timeRef)

		// only schedule timer if needed (request AND there is no timer scheduled already for somewhere before this period)
		d := state.nextCoalescedTimerDuration
		if d != math.MaxInt64 && now+d > state.lastScheduledWakeup {
			eventpbdsl.TimerDelay(m, mc.Timer,
				[]*eventpbtypes.Event{
					directorpbevents.Heartbeat(mc.Self),
				},
				timert.Duration(d),
			)

			// clear coalesced timer for next batch of events
			state.nextCoalescedTimerDuration = 0
		}

		state.lastWakeUp = now
		state.lastScheduledWakeup = now + d

		return nil
	})
}

func (state *state) wakeUpAfter(d time.Duration) {
	if d < state.nextCoalescedTimerDuration {
		state.nextCoalescedTimerDuration = d
	}
}

func (state *state) agCanDeliver(min int) bool {
	nCanDeliver := 0

	for queueIdx := 0; queueIdx < len(state.slotsReadyToDeliver); queueIdx++ {
		queueSlot := state.agQueueHeads[queueIdx]
		if _, bcDone := state.slotsReadyToDeliver[queueIdx][queueSlot]; bcDone {
			nCanDeliver++

			if nCanDeliver >= min {
				return true
			}
		}
	}

	return false
}

func newState(params common.ModuleParams, tunables common.ModuleTunables, nodeID t.NodeID) *state {
	N := len(params.AllNodes)

	state := &state{
		bcOwnQueueHead: 0,

		agQueueHeads: make([]aleatypes.QueueSlot, N),

		slotsReadyToDeliver: make([]set[aleatypes.QueueSlot], N),

		nextCoalescedTimerDuration: math.MaxInt64,

		targetOwnUnagreedBatches: 1,
	}

	ownQueueIdx := uint32(slices.Index(params.AllNodes, nodeID))

	for queueIdx := uint32(0); queueIdx < uint32(N); queueIdx++ {
		state.agQueueHeads[queueIdx] = 0

		if queueIdx == ownQueueIdx {
			state.slotsReadyToDeliver[queueIdx] = make(set[aleatypes.QueueSlot], tunables.MaxOwnUnagreedBatchCount)
		} else {
			state.slotsReadyToDeliver[queueIdx] = make(set[aleatypes.QueueSlot], tunables.MaxConcurrentVcbPerQueue)
		}
	}

	return state
}
