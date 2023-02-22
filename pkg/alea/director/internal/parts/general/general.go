package general

import (
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	mempooldsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	mempoolevents "github.com/filecoin-project/mir/pkg/mempool/events"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	aagdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/dsl"
	abcdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/mempoolpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type set[T comparable] map[T]struct{}

type state struct {
	unagreedBroadcastedOwnSlotCount int
	batchCutInProgress              bool
	bcOwnQueueHead                  aleatypes.QueueSlot

	agQueueHeads []aleatypes.QueueSlot

	slotsReadyToDeliver []set[aleatypes.QueueSlot]
	agRound             uint64
	agRoundInProgress   bool
}

func Include(m dsl.Module, mc *common.ModuleConfig, params *common.ModuleParams, tunables *common.ModuleTunables, nodeID t.NodeID, logger logging.Logger) {
	state := newState(params, tunables, nodeID)
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(params.AllNodes, nodeID))

	// TODO: split sections into different parts (with independent state)

	// =============================================================================================
	// Delivery
	// =============================================================================================

	abcdsl.UponDeliver(m, func(slot *commontypes.Slot) error {
		if slot.QueueSlot >= state.agQueueHeads[slot.QueueIdx] {
			// slot wasn't delivered yet by agreement component
			logger.Log(logging.LevelDebug, "marking slot as ready for delivery", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
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
		logger.Log(logging.LevelDebug, "delivering cert", "agreementRound", round, "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
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

	abcdsl.UponDeliver(m, func(slot *commontypes.Slot) error {
		if slot.QueueIdx == ownQueueIdx {
			state.unagreedBroadcastedOwnSlotCount++
		}

		return nil
	})
	aagdsl.UponDeliver(m, func(round uint64, decision bool) error {
		if decision {
			queueIdx := aleatypes.QueueIdx(round % uint64(len(params.AllNodes)))

			// track unagreed own slots
			if queueIdx == ownQueueIdx {
				state.unagreedBroadcastedOwnSlotCount--
			}
		}

		return nil
	})

	// upon unagreed own slots < max, cut a new batch and broadcast it
	// TODO: move to bc component
	dsl.UponCondition(m, func() error {
		if !state.batchCutInProgress && state.unagreedBroadcastedOwnSlotCount < tunables.TargetOwnUnagreedBatchCount {
			logger.Log(logging.LevelDebug, "requesting more transactions")
			mempooldsl.RequestBatch[struct{}](m, mc.Mempool, nil)
			state.batchCutInProgress = true
		}
		return nil
	})
	mempooldsl.UponNewBatch(m, func(txIDs []t.TxID, txs []*requestpb.Request, _ctx *struct{}) error {
		if len(txs) == 0 {
			var nilStructPtr *struct{}

			// batch is empty, try again after a bit
			dsl.EmitEvent(
				m,
				events.TimerDelay(
					mc.Timer,
					[]*eventpb.Event{mempoolevents.RequestBatch(mc.Mempool, &mempoolpb.RequestBatchOrigin{
						Module: mc.Self.Pb(),

						// we don't really use it, but DSL modules demand context
						Type: &mempoolpb.RequestBatchOrigin_Dsl{
							Dsl: dsl.Origin(
								m.DslHandle().StoreContext(nilStructPtr),
							),
						},
					})},
					tunables.BatchCutFailRetryDelay,
				),
			)

			return nil
		}

		abcdsl.StartBroadcast(m, mc.AleaBroadcast, state.bcOwnQueueHead, txs)
		state.bcOwnQueueHead++
		return nil
	})
	abcdsl.UponDeliver(m, func(slot *commontypes.Slot) error {
		if slot.QueueIdx == ownQueueIdx && slot.QueueSlot == state.bcOwnQueueHead-1 {
			// new batch was delivered
			state.batchCutInProgress = false
		}
		return nil
	})

	// =============================================================================================
	// Agreement Round Control
	// =============================================================================================

	// upon agreement round completion, prepare next round
	aagdsl.UponDeliver(m, func(round uint64, decision bool) error {
		state.agRoundInProgress = false
		state.agRound++
		return nil
	})

	// if no agreement round is running, and any queue is able to deliver in this node, start the next round
	// a queue is able to deliver in this node if its head has been broadcast to it
	dsl.UponCondition(m, func() error {
		if state.agRoundInProgress {
			return nil // nothing to do
		}

		shouldResumeAg := false
		for queueIdx := 0; queueIdx < len(state.slotsReadyToDeliver); queueIdx++ {
			queueSlot := state.agQueueHeads[queueIdx]
			if _, bcDone := state.slotsReadyToDeliver[queueIdx][queueSlot]; bcDone {
				shouldResumeAg = true
				break
			}
		}

		nextQueueIdx := aleatypes.QueueIdx((state.agRound + 1) % uint64(len(params.AllNodes)))
		nextQueueSlot := state.agQueueHeads[nextQueueIdx]

		if shouldResumeAg {
			state.agRoundInProgress = true

			_, bcDone := state.slotsReadyToDeliver[nextQueueIdx][nextQueueSlot]
			logger.Log(logging.LevelDebug, "progressing to next agreement round", "agreementRound", state.agRound, "input", bcDone)

			aagdsl.InputValue(m, mc.AleaAgreement, state.agRound, bcDone)
		} else {
			logger.Log(logging.LevelDebug, "stalling agreement", "agreementRound", state.agRound)
		}

		return nil
	})

	// upon another node starting the agreement round, input a value to it
	aagdsl.UponRequestInput(m, func(round uint64) error {
		if state.agRoundInProgress || state.agRound != round {
			return nil // out-of-order message
		}

		// we're here, so vcb hasn't completed for this slot yet: we must vote against delivery :(
		logger.Log(logging.LevelDebug, "resumption of stalled agreement round forced by another node", "agreementRound", state.agRound, "input", false)
		aagdsl.InputValue(m, mc.AleaAgreement, round, false)
		state.agRoundInProgress = true

		return nil
	})

}

func newState(params *common.ModuleParams, tunables *common.ModuleTunables, nodeID t.NodeID) *state {
	N := len(params.AllNodes)

	state := &state{
		unagreedBroadcastedOwnSlotCount: 0,
		batchCutInProgress:              false,
		bcOwnQueueHead:                  0,

		agQueueHeads: make([]aleatypes.QueueSlot, N),

		slotsReadyToDeliver: make([]set[aleatypes.QueueSlot], N),
	}

	ownQueueIdx := uint32(slices.Index(params.AllNodes, nodeID))

	for queueIdx := uint32(0); queueIdx < uint32(N); queueIdx++ {
		state.agQueueHeads[queueIdx] = 0

		if queueIdx == ownQueueIdx {
			state.slotsReadyToDeliver[queueIdx] = make(set[aleatypes.QueueSlot], tunables.TargetOwnUnagreedBatchCount)
		} else {
			state.slotsReadyToDeliver[queueIdx] = make(set[aleatypes.QueueSlot], tunables.MaxConcurrentVcbPerQueue)
		}
	}

	return state
}
