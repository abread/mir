package availability

import (
	"time"

	"golang.org/x/exp/slices"

	"github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/util"
	"github.com/filecoin-project/mir/pkg/dsl"
	bcpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

type bcEstimators struct {
	bcStartTimes map[bcpbtypes.Slot]time.Time
	tunables     ModuleTunables
	ownQueueHead uint64

	ownBcDuration         *util.Estimator
	ownBcQuorumDoneMargin *util.Estimator
	ownBcTotalDoneMargin  *util.Estimator

	extBcDuration *util.ByzEstimator

	minNetLatency *util.Estimator

	estimatesModified bool
}

func (est *bcEstimators) BcRuntime(slot bcpbtypes.Slot) (time.Duration, bool) {
	if start, ok := est.bcStartTimes[slot]; ok {
		return time.Since(start), true
	}
	return 0, false
}

func (est *bcEstimators) MaxOwnBcDuration() time.Duration {
	qMargin := est.ownBcQuorumDoneMargin.MaxEstimate()

	totalMarginIncr := est.ownBcTotalDoneMargin.MaxEstimate()
	maxTotalMarginIncr := time.Duration((est.tunables.MaxExtSlowdownFactor - 1) * float64(qMargin))
	if totalMarginIncr > maxTotalMarginIncr {
		totalMarginIncr = maxTotalMarginIncr
	}

	return est.ownBcDuration.MaxEstimate() + qMargin + totalMarginIncr
}

func (est *bcEstimators) MaxOwnBcLocalDuration() time.Duration {
	return est.ownBcDuration.MaxEstimate()
}

func (est *bcEstimators) MaxExtBcDuration() time.Duration {
	return est.extBcDuration.MaxEstimate()
}

func newBcEstimators(m dsl.Module, mc ModuleConfig, params ModuleParams, tunables ModuleTunables, nodeID t.NodeID) *bcEstimators {
	estimators := &bcEstimators{
		bcStartTimes: make(map[bcpbtypes.Slot]time.Time, len(params.AllNodes)),
		ownQueueHead: 0,

		ownBcDuration:         util.NewEstimator(tunables.EstimateWindowSize),
		ownBcQuorumDoneMargin: util.NewEstimator(tunables.EstimateWindowSize),
		ownBcTotalDoneMargin:  util.NewEstimator(tunables.EstimateWindowSize),
		extBcDuration:         util.NewByzEstimator(tunables.EstimateWindowSize, len(params.AllNodes)),
		minNetLatency:         util.NewEstimator(tunables.EstimateWindowSize),
	}
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(params.AllNodes, nodeID))

	bcqueuepbdsl.UponBcStarted(m, func(slot *bcpbtypes.Slot) error {
		if slot.QueueIdx != ownQueueIdx {
			estimators.bcStartTimes[*slot] = time.Now()
		}
		return nil
	})
	bcpbdsl.UponRequestCert(m, func() error {
		slot := bcpbtypes.Slot{
			QueueIdx:  ownQueueIdx,
			QueueSlot: aleatypes.QueueSlot(estimators.ownQueueHead),
		}
		estimators.bcStartTimes[slot] = time.Now()
		return nil
	})

	bcqueuepbdsl.UponDeliver(m, func(cert *bcpbtypes.Cert) error {
		if start, ok := estimators.bcStartTimes[*cert.Slot]; ok {
			duration := time.Since(start)

			if cert.Slot.QueueIdx == ownQueueIdx {
				estimators.ownBcDuration.AddSample(duration)
			} else {
				estimators.extBcDuration.AddSample(int(cert.Slot.QueueIdx), duration)
			}

			delete(estimators.bcStartTimes, *cert.Slot)

			estimators.estimatesModified = true
		}

		return nil
	})

	bcqueuepbdsl.UponBcQuorumDone(m, func(slot *bcpbtypes.Slot, deliverDelta time.Duration) error {
		if slot.QueueIdx != ownQueueIdx {
			return errors.Errorf("external bc cannot possibly be aware of doneness")
		}

		estimators.ownBcQuorumDoneMargin.AddSample(deliverDelta)
		estimators.estimatesModified = true
		return nil
	})

	bcqueuepbdsl.UponBcAllDone(m, func(slot *bcpbtypes.Slot, quorumDoneDelta time.Duration) error {
		if slot.QueueIdx != ownQueueIdx {
			return errors.Errorf("external bc cannot possibly be aware of doneness")
		}

		estimators.ownBcTotalDoneMargin.AddSample(quorumDoneDelta)
		estimators.estimatesModified = true
		return nil
	})

	bcqueuepbdsl.UponNetLatencyEstimate(m, func(minEstimate time.Duration) error {
		estimators.minNetLatency.AddSample(minEstimate)
		estimators.estimatesModified = true
		return nil
	})

	dsl.UponStateUpdates(m, func() error {
		if estimators.estimatesModified {
			bcpbdsl.EstimateUpdate(m, mc.AleaDirector, estimators.MaxOwnBcDuration(), estimators.MaxOwnBcLocalDuration(), estimators.MaxExtBcDuration(), estimators.minNetLatency.MinEstimate())
			estimators.estimatesModified = false
		}

		return nil
	})

	return estimators
}
