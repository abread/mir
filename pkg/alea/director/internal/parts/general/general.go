package general

import (
	"fmt"
	"time"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	aagdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/dsl"
	abcdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	directordsl "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/dsl"
	directorpbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/events"
	aleapbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/types"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	isspbdsl "github.com/filecoin-project/mir/pkg/pb/isspb/dsl"
	mempooldsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	requestpbtypes "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	timert "github.com/filecoin-project/mir/pkg/timer/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type set[T comparable] map[T]struct{}

type state struct {
	stalledBatchCut bool
	bcOwnQueueHead  aleatypes.QueueSlot

	agQueueHeads []aleatypes.QueueSlot

	slotsReadyToDeliver []set[aleatypes.QueueSlot]
	agRound             uint64
	stalledAgRound      bool

	avgAgTime   *estimator
	agStartTime time.Time

	bcStartTimes map[commontypes.Slot]time.Time
	avgOwnBcTime *estimator
	avgBcTime    *estimator
}

func Include(m dsl.Module, mc *common.ModuleConfig, params *common.ModuleParams, tunables *common.ModuleTunables, nodeID t.NodeID, logger logging.Logger) {
	state := newState(params, tunables, nodeID)
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(params.AllNodes, nodeID))

	dsl.UponInit(m, func() error {
		eventpbdsl.TimerRepeat(m, mc.Timer,
			[]*eventpbtypes.Event{
				directorpbevents.Heartbeat(mc.Self),
			},
			timert.Duration(tunables.MaxAgreementDelay),
			tt.RetentionIndex(0),
		)

		return nil
	})

	directordsl.UponHeartbeat(m, func() error {
		return nil // no-op, we just need our UponCondition code to run
	})

	// =============================================================================================
	// Broadcast duration estimation
	// =============================================================================================
	abcdsl.UponBcStarted(m, func(slot *commontypes.Slot) error {
		state.bcStartTimes[*slot] = time.Now()

		return nil
	})
	abcdsl.UponDeliver(m, func(slotRef *commontypes.Slot) error {
		slot := *slotRef

		startTime, ok := state.bcStartTimes[slot]
		if !ok {
			return nil // already processed
		}

		duration := time.Since(startTime)

		state.avgBcTime.AddSample(duration)
		if slot.QueueIdx == ownQueueIdx {
			state.avgOwnBcTime.AddSample(duration)
		}

		delete(state.bcStartTimes, slot)

		return nil
	})

	// TODO: split sections into different parts (with independent state)

	// =============================================================================================
	// Delivery
	// =============================================================================================

	abcdsl.UponDeliver(m, func(slot *commontypes.Slot) error {
		if slot.QueueSlot >= state.agQueueHeads[slot.QueueIdx] {
			// slot wasn't delivered yet by agreement component
			// logger.Log(logging.LevelDebug, "marking slot as ready for delivery", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
			state.slotsReadyToDeliver[slot.QueueIdx][slot.QueueSlot] = struct{}{}
		}

		return nil
	})

	// upon agreement round completion, deliver if it was decided to do so
	aagdsl.UponDeliver(m, func(round uint64, decision bool) error {
		if !decision {
			// nothing to deliver
			return nil
		}

		queueIdx := aleatypes.QueueIdx(round % uint64(len(params.AllNodes)))
		slot := &commontypes.Slot{
			QueueIdx:  queueIdx,
			QueueSlot: state.agQueueHeads[queueIdx],
		}

		// next round won't start until we say so, and previous rounds already delivered, so we can deliver immediately
		// logger.Log(logging.LevelDebug, "delivering cert", "agreementRound", round, "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
		isspbdsl.DeliverCert(m, mc.Consumer, tt.SeqNr(round), &availabilitypbtypes.Cert{
			Type: &availabilitypbtypes.Cert_Alea{
				Alea: &aleapbtypes.Cert{
					Slot: slot,
				},
			},
		})

		// pop queue
		state.agQueueHeads[queueIdx]++

		// remove tracked slot readyness (don't want to run out of memory)
		// also free broadcast slot to allow broadcast component to make progress
		delete(state.slotsReadyToDeliver[slot.QueueIdx], slot.QueueSlot)
		abcdsl.FreeSlot(m, mc.AleaBroadcast, slot)

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

		if !state.stalledBatchCut || unagreedOwnBatchCount >= tunables.MaxOwnUnagreedBatchCount {
			// batch cut in progress, or enough are cut already
			return nil
		}

		waitRoundCount := int(ownQueueIdx) - int(state.agRound%uint64(len(params.AllNodes))) - 1
		if waitRoundCount == -1 {
			waitRoundCount = len(params.AllNodes)
		} else if waitRoundCount < 0 {
			waitRoundCount += len(params.AllNodes)
		}

		// TODO: consider progress in current round too (will mean adjustments below)
		timeToOwnQueueAgRound := state.avgAgTime.MinEstimate() * time.Duration(waitRoundCount)
		maxTimeBeforeBatch := state.avgOwnBcTime.MaxEstimate() + tunables.BcEstimateMargin

		if state.agCanDeliver() && timeToOwnQueueAgRound > maxTimeBeforeBatch {
			// we have a lot of time before we reach our agreement round. let the batch fill up
			return nil
		}

		// logger.Log(logging.LevelDebug, "requesting more transactions")
		state.stalledBatchCut = false

		mempooldsl.RequestBatch[struct{}](m, mc.Mempool, nil)
		return nil
	})
	mempooldsl.UponNewBatch(m, func(txIDs []tt.TxID, txs []*requestpbtypes.Request, ctx *struct{}) error {
		if len(txs) == 0 {
			return fmt.Errorf("empty batch. did you misconfigure your mempool?")
		}

		// logger.Log(logging.LevelDebug, "new batch", "nTransactions", len(txs))

		abcdsl.StartBroadcast(m, mc.AleaBroadcast, state.bcOwnQueueHead, txs)
		state.bcOwnQueueHead++
		return nil
	})
	abcdsl.UponDeliver(m, func(slot *commontypes.Slot) error {
		if slot.QueueIdx == ownQueueIdx && slot.QueueSlot == state.bcOwnQueueHead-1 {
			// new batch was delivered
			state.stalledBatchCut = true
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

	// upon agreement round completion, prepare next round
	aagdsl.UponDeliver(m, func(round uint64, decision bool) error {
		state.agRound++
		state.avgAgTime.AddSample(time.Since(state.agStartTime))

		state.stalledAgRound = true
		return nil
	})

	// if no agreement round is running, and any queue is able to deliver in this node, start the next round
	// a queue is able to deliver in this node if its head has been broadcast to it
	dsl.UponStateUpdates(m, func() error {
		if !state.stalledAgRound {
			return nil // nothing to do
		}

		canDeliverSomething := state.agCanDeliver()

		nextQueueIdx := aleatypes.QueueIdx(state.agRound % uint64(len(params.AllNodes)))
		nextQueueSlot := state.agQueueHeads[nextQueueIdx]

		if canDeliverSomething {
			_, bcDone := state.slotsReadyToDeliver[nextQueueIdx][nextQueueSlot]
			if !bcDone {
				slot := commontypes.Slot{
					QueueIdx:  nextQueueIdx,
					QueueSlot: nextQueueSlot,
				}

				if startTime, ok := state.bcStartTimes[slot]; ok {
					stalledTime := time.Since(startTime)
					if stalledTime <= state.avgBcTime.MaxEstimate()+tunables.BcEstimateMargin && stalledTime < tunables.MaxAgreementDelay {
						// stall agreement to allow in-flight broadcast to complete
						return nil
					}
				}
			}

			// logger.Log(logging.LevelDebug, "progressing to next agreement round", "agreementRound", state.agRound, "input", bcDone)

			aagdsl.InputValue(m, mc.AleaAgreement, state.agRound, bcDone)
			state.stalledAgRound = false
			state.agStartTime = time.Now()
		}

		return nil
	})
}

func (state *state) agCanDeliver() bool {
	agCanDeliver := false

	for queueIdx := 0; queueIdx < len(state.slotsReadyToDeliver); queueIdx++ {
		queueSlot := state.agQueueHeads[queueIdx]
		if _, bcDone := state.slotsReadyToDeliver[queueIdx][queueSlot]; bcDone {
			agCanDeliver = true
			break
		}
	}

	return agCanDeliver
}

func newState(params *common.ModuleParams, tunables *common.ModuleTunables, nodeID t.NodeID) *state {
	N := len(params.AllNodes)

	state := &state{
		bcOwnQueueHead: 0,

		agQueueHeads: make([]aleatypes.QueueSlot, N),

		slotsReadyToDeliver: make([]set[aleatypes.QueueSlot], N),

		avgAgTime: NewEstimator(N * tunables.MaxConcurrentVcbPerQueue),

		bcStartTimes: make(map[commontypes.Slot]time.Time, tunables.MaxConcurrentVcbPerQueue*len(params.AllNodes)),
		avgOwnBcTime: NewEstimator(N * tunables.MaxConcurrentVcbPerQueue),
		avgBcTime:    NewEstimator(N * tunables.MaxConcurrentVcbPerQueue),
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
