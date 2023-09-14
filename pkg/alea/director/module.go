package director

import (
	"math"
	"time"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	aagdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/dsl"
	bcpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	directorpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/dsl"
	directorpbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/events"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	isspbdsl "github.com/filecoin-project/mir/pkg/pb/isspb/dsl"
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
}

func NewModule(mc ModuleConfig, params ModuleParams, tunables ModuleTunables, nodeID t.NodeID, logger logging.Logger) dsl.Module { // nolint: gocyclo,gocognit
	m := dsl.NewModule(mc.Self)
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(params.AllNodes, nodeID))

	est := newEstimators(m, params, tunables, nodeID)
	state := newState(params, tunables, nodeID)

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

	bcpbdsl.UponDeliverCert(m, func(cert *bcpbtypes.Cert) error {
		slot := cert.Slot

		if slot.QueueSlot >= state.agQueueHeads[slot.QueueIdx] {
			// slot wasn't delivered yet by agreement component
			// logger.Log(logging.LevelDebug, "marking slot as ready for delivery", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
			state.slotsReadyToDeliver[slot.QueueIdx][slot.QueueSlot] = struct{}{}
		} else {
			bcpbdsl.FreeSlot(m, mc.AleaBroadcast, slot)
		}

		return nil
	})

	// upon agreement round completion, deliver if it was decided to do so
	// if not, deliver empty batch (nilCert is used because we cannot currently)
	aagdsl.UponDeliver(m, func(round uint64, decision bool, _duration time.Duration, _posQuorumWait time.Duration) error {
		if !decision {
			// deliver empty batch
			isspbdsl.DeliverCert(m, mc.Consumer, tt.SeqNr(round), nil, true)
		}

		queueIdx := aleatypes.QueueIdx(round % uint64(N))
		slot := &bcpbtypes.Slot{
			QueueIdx:  queueIdx,
			QueueSlot: state.agQueueHeads[queueIdx],
		}

		// next round won't start until we say so, and previous rounds already delivered, so we can deliver immediately
		logger.Log(logging.LevelDebug, "delivering cert", "agreementRound", round, "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
		isspbdsl.DeliverCert(m, mc.Consumer, tt.SeqNr(round), &availabilitypbtypes.Cert{
			Type: &availabilitypbtypes.Cert_Alea{
				Alea: &bcpbtypes.Cert{
					Slot: slot,
				},
			},
		}, false)

		// pop queue
		state.agQueueHeads[queueIdx]++

		// remove tracked slot readiness (don't want to run out of memory)
		// also free broadcast slot to allow broadcast component to make progress
		delete(state.slotsReadyToDeliver[slot.QueueIdx], slot.QueueSlot)
		bcqueuepbdsl.FreeSlot(m, bccommon.BcQueueModuleID(mc.AleaBroadcast, slot.QueueIdx), slot.QueueSlot)

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
		if cert.Slot.QueueSlot != state.agQueueHeads[cert.Slot.QueueIdx] {
			return nil
		}

		currentRoundQueueIdx := aleatypes.QueueIdx(state.agRound % uint64(N))
		if currentRoundQueueIdx == cert.Slot.QueueIdx {
			// this slot is for the current ag round: we need to be carefult not to input a value twice
			nextRound := state.agRound
			if state.stalledAgRound {
				logger.Log(logging.LevelDebug, "INPUT AG (BC-current)", "round", nextRound, "value", true)
				aagdsl.InputValue(m, mc.AleaAgreement, nextRound, true)
				state.stalledAgRound = false
			}
		} else {
			// this slot is not for the current ag round: we can freely input to it
			var nextRound uint64
			if currentRoundQueueIdx < cert.Slot.QueueIdx {
				// we need to go a few rounds further
				nextRound = state.agRound + uint64(cert.Slot.QueueIdx-currentRoundQueueIdx)
			} else {
				nextRound = state.agRound + uint64(N) - uint64(currentRoundQueueIdx-cert.Slot.QueueIdx)
			}

			logger.Log(logging.LevelDebug, "INPUT AG (BC-future)", "round", nextRound, "value", true)
			aagdsl.InputValue(m, mc.AleaAgreement, nextRound, true)
		}

		return nil
	})

	// Slots may be broadcast long before agreement delivers them: this event handler provides input to ag
	// for those slots that we left behind in the previous handler.
	aagdsl.UponDeliver(m, func(round uint64, decision bool, _, _ time.Duration) error {
		// if we delivered a new slot in a queue, we can input one for the next slot in the same queue
		queueIdx := aleatypes.QueueIdx(round % uint64(N))
		nextQueueSlot := state.agQueueHeads[queueIdx]
		if _, present := state.slotsReadyToDeliver[queueIdx][nextQueueSlot]; present {
			nextRound := round + uint64(N)

			logger.Log(logging.LevelDebug, "INPUT AG (AG-done)", "round", nextRound, "value", true)
			aagdsl.InputValue(m, mc.AleaAgreement, nextRound, true)
		}
		return nil
	})

	// upon init, stall agreement until a slot is deliverable
	dsl.UponInit(m, func() error {
		state.stalledAgRound = true
		return nil
	})

	aagdsl.UponDeliver(m, func(round uint64, _ bool, _, _ time.Duration) error {
		// advance to next round
		state.agRound++

		// if bc has delivered the slot for this round already, then we have already input 1 to the
		// round, and it is not stalled at all.
		nextRoundQueueIdx := aleatypes.QueueIdx(state.agRound % uint64(N))
		nextRoundQueueSlot := state.agQueueHeads[nextRoundQueueIdx]
		if _, present := state.slotsReadyToDeliver[nextRoundQueueIdx][nextRoundQueueSlot]; present {
			state.stalledAgRound = false
		} else {
			state.stalledAgRound = true
		}

		return nil
	})

	// if no agreement round is running, and any queue is able to deliver in this node, start the next round
	// a queue is able to deliver in this node if its head has been broadcast to it
	dsl.UponStateUpdates(m, func() error {
		if !state.stalledAgRound {
			return nil // nothing to do
		} else if !state.agCanDeliver(1) {
			// just continue stalling: we don't have anything to deliver yet
			return nil
		}

		nextSlotQueueIdx := aleatypes.QueueIdx(state.agRound % uint64(N))
		nextSlot := bcpbtypes.Slot{
			QueueIdx:  nextSlotQueueIdx,
			QueueSlot: state.agQueueHeads[nextSlotQueueIdx],
		}

		// entering this code path means nextSlot was not delivered yet, and we should input 0 to the
		// agreement round

		// delay inputting 0 when a broadcast is in progress
		if bcRuntime, ok := est.BcRuntime(nextSlot); ok {
			if nextSlot.QueueIdx == ownQueueIdx && bcRuntime < est.OwnBcMedianDurationEstNoMargin() {
				//logger.Log(logging.LevelDebug, "stalling agreement input for own batch")
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
				// logger.Log(logging.LevelDebug, "stalling agreement input", "maxDelay", maxTimeToWait)
				state.wakeUpAfter(maxTimeToWait)

				return nil
			}
		}

		// logger.Log(logging.LevelDebug, "progressing to next agreement round", "agreementRound", state.agRound, "input", bcDone)

		logger.Log(logging.LevelDebug, "INPUT AG (timeout)", "round", state.agRound, "value", false)
		aagdsl.InputValue(m, mc.AleaAgreement, state.agRound, false)
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
		// bcOwnQueueHead is the next slot to be broadcast
		// agQueueHeads[ownQueueIdx] is the next slot to be agreed on
		unagreedOwnBatchCount := uint64(state.bcOwnQueueHead - state.agQueueHeads[ownQueueIdx])

		if !state.stalledBatchCut || unagreedOwnBatchCount >= uint64(tunables.MaxOwnUnagreedBatchCount) {
			// batch cut in progress, or enough are cut already
			return nil
		}

		waitRoundCount := int(ownQueueIdx) - int(state.agRound%uint64(N)) - 1
		if waitRoundCount < 0 {
			waitRoundCount += N
		}

		// consider how many batches we have pending delivery
		waitRoundCount += int(unagreedOwnBatchCount) * N

		// give one ag round of leeway
		if waitRoundCount > 0 {
			waitRoundCount--
		}

		timeToOwnQueueAgRound := est.AgFastPathEst() * time.Duration(waitRoundCount)
		bcRuntimeEst := est.OwnBcMaxDurationEst()

		// We have a lot of time before we reach our agreement round. Let the batch fill up!
		// That said, if we have no unagreed batches from any correct replica, we'll fast track it.
		if timeToOwnQueueAgRound > bcRuntimeEst && (unagreedOwnBatchCount > 0 || state.agCanDeliver(F+1)) {
			// ensure we are woken up to create a batch before we run out of time
			maxDelay := timeToOwnQueueAgRound - bcRuntimeEst
			// logger.Log(logging.LevelDebug, "stalling batch cut", "max delay", maxDelay)
			state.wakeUpAfter(maxDelay)

			return nil
		}

		// logger.Log(logging.LevelDebug, "requesting more transactions")
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

			state.lastScheduledWakeup = now + d
		}

		state.lastWakeUp = now

		// clear coalesced timer for next batch of events
		state.nextCoalescedTimerDuration = math.MaxInt64

		return nil
	})

	return m
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

func newState(params ModuleParams, tunables ModuleTunables, nodeID t.NodeID) *state {
	N := len(params.AllNodes)

	state := &state{
		bcOwnQueueHead: 0,

		agQueueHeads: make([]aleatypes.QueueSlot, N),

		slotsReadyToDeliver: make([]set[aleatypes.QueueSlot], N),

		nextCoalescedTimerDuration: math.MaxInt64,
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
