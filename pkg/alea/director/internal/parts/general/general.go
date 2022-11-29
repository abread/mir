package general

import (
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/agreement/aagdsl"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/abcdsl"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	batchdbdsl "github.com/filecoin-project/mir/pkg/availability/batchdb/dsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	mempooldsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	mempoolevents "github.com/filecoin-project/mir/pkg/mempool/events"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	aleapbCommon "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
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
	bcOwnQueueHead                  uint64

	agQueueHeads map[uint32]uint64

	slotsReadyToDeliver   map[uint32]set[uint64]
	stalledAgreementSlot  *aleapbCommon.Slot
	stalledAgreementRound uint64
}

func Include(m dsl.Module, mc *common.ModuleConfig, params *common.ModuleParams, tunables *common.ModuleTunables, nodeID t.NodeID, logger logging.Logger) {
	state := newState(params, tunables, nodeID)
	ownQueueIdx := uint32(slices.Index(params.AllNodes, nodeID))

	// TODO: split sections into different parts (with independent state)

	// =============================================================================================
	// Delivery
	// =============================================================================================

	// upon broadcast completion, save transactions
	abcdsl.UponDeliver(m, func(slot *aleapbCommon.Slot, txIDs []t.TxID, txs []*requestpb.Request, signature []byte) error {
		batchdbdsl.StoreBatch(m, mc.BatchDB, common.FormatAleaBatchID(slot), txIDs, txs, signature, slot)
		return nil
	})
	batchdbdsl.UponBatchStored(m, func(slot *aleapbCommon.Slot) error {
		// track vcb completion (needed for other ops in agreement round control)
		state.slotsReadyToDeliver[slot.QueueIdx][slot.QueueSlot] = struct{}{}
		return nil
	})

	// upon agreement round completion, deliver if it was decided to do so
	aagdsl.UponDeliver(m, func(round uint64, decision bool) error {
		if !decision {
			// nothing to deliver
			return nil
		}

		queueIdx := uint32(round % uint64(len(params.AllNodes)))
		slot := &aleapbCommon.Slot{
			QueueIdx:  queueIdx,
			QueueSlot: state.agQueueHeads[queueIdx],
		}

		// next round won't start until we say so, and previous rounds already delivered, so we can deliver immediately
		logger.Log(logging.LevelDebug, "delivering cert", "agreementRound", round, "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
		dsl.EmitEvent(m, events.DeliverCert(mc.Consumer, t.SeqNr(round), &availabilitypb.Cert{
			Type: &availabilitypb.Cert_Alea{
				Alea: &aleapb.Cert{
					Slot: slot,
				},
			},
		}))

		// pop queue
		state.agQueueHeads[queueIdx]++

		// remove tracked slot readyness (don't want to run out of memory)
		delete(state.slotsReadyToDeliver[slot.QueueIdx], slot.QueueSlot)

		return nil
	})

	// =============================================================================================
	// Batch Cutting / Own Queue Broadcast Control
	// =============================================================================================

	abcdsl.UponDeliver(m, func(slot *aleapbCommon.Slot, _txIDs []t.TxID, _txs []*requestpb.Request, _signature []byte) error {
		if slot.QueueIdx == ownQueueIdx {
			state.unagreedBroadcastedOwnSlotCount++
		}
		return nil
	})
	aagdsl.UponDeliver(m, func(round uint64, decision bool) error {
		queueIdx := uint32(round % uint64(len(params.AllNodes)))
		queueSlot := state.agQueueHeads[queueIdx]

		if decision {
			// track unagreed own slots
			if queueIdx == ownQueueIdx {
				state.unagreedBroadcastedOwnSlotCount--
			}

			// free broadcast slots
			// TODO: only free slot when broadcast/fill-gap concludes?
			// agreement takes longer so this shouldn't be a show-stopper
			abcdsl.FreeSlot(m, mc.AleaBroadcast, &aleapbCommon.Slot{
				QueueIdx:  queueIdx,
				QueueSlot: queueSlot,
			})
		}

		return nil
	})

	// upon unagreed own slots < max, cut a new batch and broadcast it
	dsl.UponCondition(m, func() error {
		if !state.batchCutInProgress && state.unagreedBroadcastedOwnSlotCount < tunables.TargetOwnUnagreedBatchCount {
			logger.Log(logging.LevelDebug, "requesting more transactions")
			mempooldsl.RequestBatch(m, mc.Mempool, &struct{}{})
			state.batchCutInProgress = true
		}
		return nil
	})
	mempooldsl.UponNewBatch(m, func(txIDs []t.TxID, txs []*requestpb.Request, context *struct{}) error {
		if len(txs) == 0 {
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
								m.DslHandle().StoreContext(&struct{}{}),
							),
						},
					})},
					tunables.BatchCutFailRetryDelay,
				),
			)

			return nil
		}

		abcdsl.StartBroadcast(m, mc.AleaBroadcast, state.bcOwnQueueHead, txIDs, txs)
		state.bcOwnQueueHead++
		return nil
	})
	abcdsl.UponDeliver(m, func(slot *aleapbCommon.Slot, _txIDs []t.TxID, _txs []*requestpb.Request, _signature []byte) error {
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
		// prepare next round
		nextQueueIdx := uint32((round + 1) % uint64(len(params.AllNodes)))
		nextQueueSlot := state.agQueueHeads[nextQueueIdx]

		// start next round if corresponding slot's broadcast is complete (or if we have pending requests)
		// otherwise stall until vcb completes (for this slot or one of ours) or another node starts the round
		if _, bcDone := state.slotsReadyToDeliver[nextQueueIdx][nextQueueSlot]; bcDone || state.unagreedBroadcastedOwnSlotCount > 0 { // TODO: track unagreedBcOwnSlotCount separately
			logger.Log(logging.LevelDebug, "progressing to next round", "agreementRound", round+1, "input", true)
			aagdsl.InputValue(m, mc.AleaAgreement, round+1, true)
			state.stalledAgreementSlot = nil
		} else {
			logger.Log(logging.LevelDebug, "stalling next round", "agreementRound", round+1)
			state.stalledAgreementSlot = &aleapbCommon.Slot{
				QueueIdx:  nextQueueIdx,
				QueueSlot: nextQueueSlot,
			}
			state.stalledAgreementRound = round + 1
		}

		return nil
	})

	// upon vcb completion for the stalled agreement slot or one of ours, start the stalled agreement round
	abcdsl.UponDeliver(m, func(slot *aleapbCommon.Slot, _txIDs []t.TxID, _txs []*requestpb.Request, _signature []byte) error {
		if state.stalledAgreementSlot == nil {
			return nil // nothing to do
		}

		recvdStalledSlotBc := slot.QueueIdx == state.stalledAgreementSlot.QueueIdx && slot.QueueSlot == state.stalledAgreementSlot.QueueSlot

		if slot.QueueIdx == ownQueueIdx || recvdStalledSlotBc {
			// input value to kickstart stalled agreement round
			// previously we did not receive the stalled slot, so we either received now and vote
			// for delivery or vote against it
			logger.Log(logging.LevelDebug, "resuming stalled agreement round", "agreementRound", state.stalledAgreementRound, "input", recvdStalledSlotBc)
			aagdsl.InputValue(m, mc.AleaAgreement, state.stalledAgreementRound, recvdStalledSlotBc)
			state.stalledAgreementSlot = nil
		}

		return nil
	})

	// upon another node starting the agreement round, input a value to it
	aagdsl.UponRequestInput(m, func(round uint64) error {
		if round != state.stalledAgreementRound || state.stalledAgreementSlot == nil {
			return nil // out of order message
		}

		// we're here, so vcb hasn't completed for this slot yet: we must vote against delivery :(
		logger.Log(logging.LevelDebug, "resuming stalled agreement round", "agreementRound", state.stalledAgreementRound, "input", false)
		aagdsl.InputValue(m, mc.AleaAgreement, round, false)
		state.stalledAgreementSlot = nil

		return nil
	})

}

func newState(params *common.ModuleParams, tunables *common.ModuleTunables, nodeID t.NodeID) *state {
	N := len(params.AllNodes)

	state := &state{
		unagreedBroadcastedOwnSlotCount: 0,
		batchCutInProgress:              false,
		bcOwnQueueHead:                  0,

		agQueueHeads: make(map[uint32]uint64, N),

		slotsReadyToDeliver: make(map[uint32]set[uint64], N),
		stalledAgreementSlot: &aleapbCommon.Slot{
			QueueIdx:  0,
			QueueSlot: 0,
		},
		stalledAgreementRound: 0,
	}

	ownQueueIdx := uint32(slices.Index(params.AllNodes, nodeID))

	for queueIdx := uint32(0); queueIdx < uint32(N); queueIdx++ {
		state.agQueueHeads[queueIdx] = 0

		if queueIdx == ownQueueIdx {
			state.slotsReadyToDeliver[queueIdx] = make(set[uint64], tunables.TargetOwnUnagreedBatchCount)
		} else {
			state.slotsReadyToDeliver[queueIdx] = make(set[uint64], tunables.MaxConcurrentVcbPerQueue)
		}
	}

	return state
}
