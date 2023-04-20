package general

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	aagdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/dsl"
	abcdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	directordsl "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/dsl"
	directorpbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/events"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	mempooldsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	requestpbtypes "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type set[T comparable] map[T]struct{}

type state struct {
	stalledBatchCut trace.Span
	bcOwnQueueHead  aleatypes.QueueSlot

	agQueueHeads []aleatypes.QueueSlot

	slotsReadyToDeliver []set[aleatypes.QueueSlot]
	agRound             uint64
	stalledRoundSpan    trace.Span

	avgAgTime   estimator
	agStartTime time.Time

	bcStartTimes map[commontypes.Slot]time.Time
	avgOwnBcTime estimator
	avgBcTime    estimator
}

func Include(m dsl.Module, mc *common.ModuleConfig, params *common.ModuleParams, tunables *common.ModuleTunables, nodeID t.NodeID, logger logging.Logger) {
	state := newState(params, tunables, nodeID)
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(params.AllNodes, nodeID))

	dsl.UponInit(m, func() error {
		dsl.EmitMirEvent(m, &eventpbtypes.Event{
			DestModule: mc.Timer,
			Type: &eventpbtypes.Event_TimerRepeat{
				TimerRepeat: &eventpb.TimerRepeat{
					Events: []*eventpb.Event{
						directorpbevents.Heartbeat(mc.Self).Pb(),
					},
					Delay:          t.TimeDuration(tunables.MaxAgreementDelay).Pb(),
					RetentionIndex: 0,
				},
			},
		})

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
		dsl.EmitEvent(m, events.DeliverCert(mc.Consumer, t.SeqNr(round), &availabilitypb.Cert{
			Type: &availabilitypb.Cert_Alea{
				Alea: &aleapb.Cert{
					Slot: slot.Pb(),
				},
			},
		}))

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
		state.stalledBatchCut = nil

		_, span := m.DslHandle().PushSpan("fetching new batch from mempool")
		mempooldsl.RequestBatch(m, mc.Mempool, &span)
		return nil
	})

	// upon nice condition (unagreed batch count < max, no batch being cut, timeToNextAgForThisNode < estBc+margin || stalled ag), cut a new batch and broadcast it
	// TODO: move to bc component
	dsl.UponCondition(m, func() error {
		// bcOwnQueueHead is the next slot to be broadcast
		// agQueueHeads[ownQueueIdx] is the next slot to be agreed on
		unagreedOwnBatchCount := uint64(state.bcOwnQueueHead - state.agQueueHeads[ownQueueIdx])

		if state.stalledBatchCut == nil || unagreedOwnBatchCount >= tunables.MaxOwnUnagreedBatchCount {
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
		timeToOwnQueueAgRound := state.avgAgTime.Estimate() * time.Duration(waitRoundCount)
		maxTimeBeforeBatch := state.avgOwnBcTime.Estimate() + tunables.BcEstimateMargin

		if state.agCanDeliver() && timeToOwnQueueAgRound > maxTimeBeforeBatch {
			// we have a lot of time before we reach our agreement round. let the batch fill up
			return nil
		}

		// logger.Log(logging.LevelDebug, "requesting more transactions")
		state.stalledBatchCut.AddEvent("resuming after unagreed own batch count got too low", trace.WithAttributes(attribute.Int64("unagreedOwnBatchCount", int64(unagreedOwnBatchCount))))

		state.stalledBatchCut.End()
		state.stalledBatchCut = nil

		_, span := m.DslHandle().PushSpan("fetching new batch from mempool")
		mempooldsl.RequestBatch(m, mc.Mempool, &span)
		return nil
	})
	mempooldsl.UponNewBatch(m, func(txIDs []t.TxID, txs []*requestpbtypes.Request, span *trace.Span) error {
		if len(txs) == 0 {
			return fmt.Errorf("empty batch. did you misconfigure your mempool?")
		}

		// logger.Log(logging.LevelDebug, "new batch", "nTransactions", len(txs))
		(*span).AddEvent("got transactions")
		(*span).End()

		abcdsl.StartBroadcast(m, mc.AleaBroadcast, state.bcOwnQueueHead, txs)
		state.bcOwnQueueHead++
		return nil
	})
	abcdsl.UponDeliver(m, func(slot *commontypes.Slot) error {
		if slot.QueueIdx == ownQueueIdx && slot.QueueSlot == state.bcOwnQueueHead-1 {
			// new batch was delivered, stall next one until ready
			_, state.stalledBatchCut = m.DslHandle().PushSpan("stalling batch cut", trace.WithAttributes(attribute.Int64("slot", int64(state.bcOwnQueueHead))))
		}
		return nil
	})

	// =============================================================================================
	// Agreement Round Control
	// =============================================================================================

	// upon init, stall agreement until a slot is deliverable
	dsl.UponInit(m, func() error {
		_, state.stalledRoundSpan = m.DslHandle().PushSpan("stalling agreement round", trace.WithAttributes(attribute.Int64("agRound", int64(state.agRound))))
		return nil
	})

	// upon agreement round completion, prepare next round
	aagdsl.UponDeliver(m, func(round uint64, decision bool) error {
		state.agRound++
		if state.stalledRoundSpan != nil {
			state.stalledRoundSpan.AddEvent("agreement delivered without local input")
			state.stalledRoundSpan.End()
		}

		state.avgAgTime.AddSample(time.Since(state.agStartTime))

		_, state.stalledRoundSpan = m.DslHandle().PushSpan("stalling agreement round", trace.WithAttributes(attribute.Int64("agRound", int64(state.agRound))))
		return nil
	})

	// if no agreement round is running, and any queue is able to deliver in this node, start the next round
	// a queue is able to deliver in this node if its head has been broadcast to it
	dsl.UponCondition(m, func() error {
		if state.stalledRoundSpan == nil {
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
					if stalledTime <= state.avgBcTime.Estimate()+tunables.BcEstimateMargin && stalledTime < tunables.MaxAgreementDelay {
						// stall agreement to allow in-flight broadcast to complete
						return nil
					}
				}
			}

			// logger.Log(logging.LevelDebug, "progressing to next agreement round", "agreementRound", state.agRound, "input", bcDone)

			state.stalledRoundSpan.End()
			state.stalledRoundSpan = nil

			aagdsl.InputValue(m, mc.AleaAgreement, state.agRound, bcDone)
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

		bcStartTimes: make(map[commontypes.Slot]time.Time, tunables.MaxConcurrentVcbPerQueue*len(params.AllNodes)),
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
