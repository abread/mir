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

	ownBcDuration util.Estimator
	extBcDuration util.ByzEstimator

	ownBcFinishMargin util.Estimator
	extBcFinishMargin util.ByzEstimator

	agDuration util.Estimator
}

func (e *Estimators) OwnBcMaxDurationEst() time.Duration {
	return e.ownBcDuration.MaxEstimate() + e.ownBcFinishMargin.MaxEstimate()
}

func (e *Estimators) OwnBcMinDurationEst() time.Duration {
	return e.ownBcDuration.MinEstimate() + e.ownBcFinishMargin.MinEstimate()
}

func (e *Estimators) ExtBcMaxDurationEst() time.Duration {
	return e.extBcDuration.MaxEstimate() + e.extBcFinishMargin.MaxEstimate()
}

func (e *Estimators) ExtBcMinDurationEst() time.Duration {
	return e.extBcDuration.MinEstimate() + e.extBcFinishMargin.MinEstimate()
}

func (e *Estimators) AgMaxDurationEst() time.Duration {
	return e.agDuration.MaxEstimate()
}

func (e *Estimators) AgMinDurationEst() time.Duration {
	return e.agDuration.MinEstimate()
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

		ownBcDuration: util.NewEstimator(tunables.EstimateWindowSize),
		extBcDuration: util.NewByzEstimator(tunables.EstimateWindowSize, N),

		ownBcFinishMargin: util.NewEstimator(tunables.EstimateWindowSize),
		extBcFinishMargin: util.NewByzEstimator(tunables.EstimateWindowSize, N),

		agDuration: util.NewEstimator(tunables.EstimateWindowSize),
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
	bcqueuepbdsl.UponBcDone(m, func(slot *commontypes.Slot, deliverDelta time.Duration) error {
		// adjust own bc estimate margin
		est.ownBcFinishMargin.AddSample(deliverDelta)
		return nil
	})
	ageventsdsl.UponDeliver(m, func(round uint64, decision bool, duration time.Duration, posQuorumWait time.Duration) error {
		// adjust other bc estimate margins
		queueIdx := aleatypes.QueueIdx(round % uint64(N))
		if queueIdx != ownQueueIdx {
			if posQuorumWait == math.MaxInt64 || !decision {
				// failed deadline, double margin
				m := est.extBcFinishMargin.MaxEstimate()

				est.extBcFinishMargin.Clear(int(queueIdx))
				est.extBcFinishMargin.AddSample(int(queueIdx), 2*m)
			} else {
				est.extBcFinishMargin.AddSample(int(queueIdx), posQuorumWait)
			}
		}
		return nil
	})

	// =============================================================================================
	// Ag Duration Estimation
	// =============================================================================================
	ageventsdsl.UponDeliver(m, func(_round uint64, _decision bool, duration time.Duration, _posQuorumWait time.Duration) error {
		est.agDuration.AddSample(duration)
		return nil
	})

	return est
}
