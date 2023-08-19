package estimators

import (
	"math"
	"time"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/alea/util"
	"github.com/filecoin-project/mir/pkg/dsl"
	ageventsdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/dsl"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type Estimators struct {
	bcStartTimes map[commontypes.Slot]time.Time

	ownBcDuration         *util.Estimator
	ownBcQuorumDoneMargin *util.Estimator
	ownBcTotalDoneMargin  *util.Estimator

	extBcDuration   *util.ByzEstimator
	extBcDoneMargin *util.ByzEstimator

	abbaRoundNoCoinDuration *util.Estimator
}

func (e *Estimators) OwnBcMaxDurationEst() time.Duration {
	return e.ownBcDuration.MaxEstimate() + e.ownBcQuorumDoneMargin.MaxEstimate() + e.ownBcTotalDoneMargin.MaxEstimate()
}

func (e *Estimators) OwnBcMedianDurationEstNoMargin() time.Duration {
	return e.ownBcDuration.Median()
}

func (e *Estimators) ExtBcMaxDurationEst() time.Duration {
	return e.extBcDuration.MaxEstimate() + e.extBcDoneMargin.MaxEstimate()
}

func (e *Estimators) AgMinDurationEst() time.Duration {
	// the unanimity optimization lowers convergence to the time for one message broadcast (per node):
	// the INIT message, containing the initial input
	// one abba round takes roughly 3 message broadcasts (per node), +1 for the common coin
	return e.abbaRoundNoCoinDuration.MinEstimate() / 3
}

func (e *Estimators) BcRuntime(slot commontypes.Slot) (time.Duration, bool) {
	if startTime, ok := e.bcStartTimes[slot]; ok {
		return time.Since(startTime), true
	}

	return 0, false
}

func (e *Estimators) MarkBcStartedNow(slot commontypes.Slot) {
	e.bcStartTimes[slot] = time.Now()
}

func New(m dsl.Module, mc common.ModuleConfig, params common.ModuleParams, tunables common.ModuleTunables, nodeID t.NodeID) *Estimators {
	N := len(params.AllNodes)
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(params.AllNodes, nodeID))

	est := &Estimators{
		bcStartTimes: make(map[commontypes.Slot]time.Time, (N-1)*tunables.MaxConcurrentVcbPerQueue+tunables.MaxOwnUnagreedBatchCount),

		ownBcDuration:         util.NewEstimator(tunables.EstimateWindowSize),
		ownBcQuorumDoneMargin: util.NewEstimator(tunables.EstimateWindowSize),
		ownBcTotalDoneMargin:  util.NewEstimator(tunables.EstimateWindowSize),

		extBcDuration:   util.NewByzEstimator(tunables.EstimateWindowSize, N),
		extBcDoneMargin: util.NewByzEstimator(tunables.EstimateWindowSize, N),

		abbaRoundNoCoinDuration: util.NewEstimator(tunables.EstimateWindowSize),
	}

	// =============================================================================================
	// Bc Duration Estimation
	// =============================================================================================
	bcqueuepbdsl.UponBcStarted(m, func(slot *commontypes.Slot) error {
		if _, ok := est.bcStartTimes[*slot]; !ok {
			est.bcStartTimes[*slot] = time.Now()
		}

		return nil
	})
	bcqueuepbdsl.UponDeliver(m, func(slotRef *commontypes.Slot) error {
		slot := *slotRef

		startTime, ok := est.bcStartTimes[slot]
		if !ok {
			return nil // already processed
		}

		duration := time.Since(startTime)

		if slot.QueueIdx == ownQueueIdx {
			est.ownBcDuration.AddSample(duration)
		} else {
			est.extBcDuration.AddSample(int(slot.QueueIdx), duration)
		}

		delete(est.bcStartTimes, slot)

		return nil
	})

	// =============================================================================================
	// Bc Finish Duration Estimation
	// =============================================================================================
	bcqueuepbdsl.UponBcQuorumDone(m, func(slot *commontypes.Slot, deliverDelta time.Duration) error {
		// adjust own bc estimate margin
		est.ownBcQuorumDoneMargin.AddSample(deliverDelta)
		return nil
	})
	bcqueuepbdsl.UponBcAllDone(m, func(slot *commontypes.Slot, quorumDoneDelta time.Duration) error {
		// adjust own bc estimate margin
		// this pertains to the last F nodes, so we must limit it based on the quorum margin
		withoutTotal := float64(est.ownBcDuration.MaxEstimate() + est.ownBcQuorumDoneMargin.MaxEstimate())
		withTotal := withoutTotal + float64(quorumDoneDelta)
		if withTotal > withoutTotal*tunables.MaxExtSlowdownFactor {
			quorumDoneDelta = time.Duration(withoutTotal * (tunables.MaxExtSlowdownFactor - 1))
		}

		est.ownBcTotalDoneMargin.AddSample(quorumDoneDelta)
		return nil
	})
	ageventsdsl.UponDeliver(m, func(round uint64, decision bool, posQuorumWait time.Duration, posTotalDelta time.Duration) error {
		// adjust other bc estimate margins
		queueIdx := aleatypes.QueueIdx(round % uint64(N))
		if queueIdx != ownQueueIdx {
			if posQuorumWait == math.MaxInt64 {
				// failed deadline, double margin
				m := est.extBcDoneMargin.MaxEstimate()

				est.extBcDoneMargin.Clear(int(queueIdx))
				est.extBcDoneMargin.AddSample(int(queueIdx), 2*m)
			} else {
				// limit slow node influence
				if float64(posTotalDelta) > float64(posQuorumWait)*tunables.MaxExtSlowdownFactor {
					posTotalDelta = time.Duration(float64(posQuorumWait) * tunables.MaxExtSlowdownFactor)
				}

				// failed deadline, ensure margin is at least doubled
				if !decision && posTotalDelta < posQuorumWait {
					posTotalDelta = posQuorumWait
				}

				est.extBcDoneMargin.AddSample(int(queueIdx), posQuorumWait+posTotalDelta)
			}
		}
		return nil
	})

	// =============================================================================================
	// Ag Duration Estimation
	// =============================================================================================
	ageventsdsl.UponInnerAbbaRoundTime(m, func(durationNoCoin time.Duration) error {
		est.abbaRoundNoCoinDuration.AddSample(durationNoCoin)
		return nil
	})

	return est
}
