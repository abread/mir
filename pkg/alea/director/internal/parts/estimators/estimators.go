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

	ownBcDuration           util.Estimator
	ownBcQuorumFinishMargin util.Estimator
	ownBcTotalFinishMargin  util.Estimator

	extBcDuration     util.ByzEstimator
	extBcFinishMargin util.ByzEstimator

	abbaRoundDuration util.Estimator
}

func (e *Estimators) OwnBcMaxDurationEst() time.Duration {
	return e.ownBcDuration.MaxEstimate() + e.ownBcQuorumFinishMargin.MaxEstimate() + e.ownBcTotalFinishMargin.MaxEstimate()
}

func (e *Estimators) ExtBcMedianDurationEst() time.Duration {
	return e.extBcDuration.Median() + e.extBcFinishMargin.Median()
}

func (e *Estimators) ExtBcMaxDurationEst() time.Duration {
	return e.extBcDuration.MaxEstimate() + e.extBcFinishMargin.MaxEstimate()
}

func (e *Estimators) AgMinDurationEst() time.Duration {
	// the unanimity optimization lowers convergence to the time for one message broadcast (per node):
	// the INIT message, containing the initial input
	// one abba round takes roughly 5 message broadcasts (per node)
	return e.abbaRoundDuration.MinEstimate() / 5
}

func (e *Estimators) BcRuntime(slot commontypes.Slot) (time.Duration, bool) {
	if startTime, ok := e.bcStartTimes[slot]; ok {
		return time.Since(startTime), true
	}

	return 0, false
}

func New(m dsl.Module, params common.ModuleParams, tunables common.ModuleTunables, nodeID t.NodeID) *Estimators {
	N := len(params.AllNodes)
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(params.AllNodes, nodeID))

	est := &Estimators{
		bcStartTimes: make(map[commontypes.Slot]time.Time, (N-1)*tunables.MaxConcurrentVcbPerQueue+int(tunables.MaxOwnUnagreedBatchCount)),

		ownBcDuration:           util.NewEstimator(tunables.EstimateWindowSize),
		ownBcQuorumFinishMargin: util.NewEstimator(tunables.EstimateWindowSize),
		ownBcTotalFinishMargin:  util.NewEstimator(tunables.EstimateWindowSize),

		extBcDuration:     util.NewByzEstimator(tunables.EstimateWindowSize, N),
		extBcFinishMargin: util.NewByzEstimator(tunables.EstimateWindowSize, N),

		abbaRoundDuration: util.NewEstimator(tunables.EstimateWindowSize),
	}

	// =============================================================================================
	// Bc Duration Estimation
	// =============================================================================================
	bcqueuepbdsl.UponBcStarted(m, func(slot *commontypes.Slot) error {
		est.bcStartTimes[*slot] = time.Now()

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
		est.ownBcQuorumFinishMargin.AddSample(deliverDelta)
		return nil
	})
	bcqueuepbdsl.UponBcAllDone(m, func(slot *commontypes.Slot, quorumDoneDelta time.Duration) error {
		// adjust own bc estimate margin
		// this pertains to the last F nodes, so we must limit it based on the quorum margin
		quorumMargin := est.ownBcQuorumFinishMargin.MaxEstimate()
		if float64(quorumDoneDelta) > float64(quorumMargin)*tunables.MaxOwnBcTotalFinishSlowdownFactor {
			quorumDoneDelta = quorumMargin
		}

		est.ownBcTotalFinishMargin.AddSample(quorumDoneDelta)
		return nil
	})
	ageventsdsl.UponDeliver(m, func(round uint64, decision bool, posQuorumWait time.Duration, posTotalDelta time.Duration) error {
		// adjust other bc estimate margins
		queueIdx := aleatypes.QueueIdx(round % uint64(N))
		if queueIdx != ownQueueIdx {
			if posQuorumWait == math.MaxInt64 {
				// failed deadline, double margin
				m := est.extBcFinishMargin.MaxEstimate()

				est.extBcFinishMargin.Clear(int(queueIdx))
				est.extBcFinishMargin.AddSample(int(queueIdx), 2*m)
			} else {
				// limit slow node influence
				if float64(posTotalDelta) > float64(posQuorumWait)*tunables.MaxOwnBcTotalFinishSlowdownFactor {
					posTotalDelta = time.Duration(float64(posQuorumWait) * tunables.MaxOwnBcTotalFinishSlowdownFactor)
				}

				// failed deadline, ensure margin is at least doubled
				if !decision && posTotalDelta < posQuorumWait {
					posTotalDelta = posQuorumWait
				}

				est.extBcFinishMargin.AddSample(int(queueIdx), posQuorumWait+posTotalDelta)
			}
		}
		return nil
	})

	// =============================================================================================
	// Ag Duration Estimation
	// =============================================================================================
	ageventsdsl.UponInnerAbbaRoundTime(m, func(duration time.Duration) error {
		est.abbaRoundDuration.AddSample(duration)
		return nil
	})

	return est
}
