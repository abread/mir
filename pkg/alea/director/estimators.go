package director

import (
	"math"
	"time"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/util"
	"github.com/filecoin-project/mir/pkg/dsl"
	ageventsdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/dsl"
	bcpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

type estimators struct {
	bcStartTimes map[bcpbtypes.Slot]time.Time

	// tracked by broadcast component
	maxOwnBcDuration      time.Duration
	maxOwnBcLocalDuration time.Duration
	maxExtBcDuration      time.Duration

	// tracked here, because it depends on agreement input timing
	extBcDoneMargin *util.ByzEstimator

	abbaRoundNoCoinDuration *util.Estimator
}

func (e *estimators) OwnBcMaxDurationEst() time.Duration {
	return e.maxOwnBcDuration
}

func (e *estimators) OwnBcLocalMaxDurationEst() time.Duration {
	return e.maxOwnBcLocalDuration
}

func (e *estimators) ExtBcMaxDurationEst() time.Duration {
	return e.maxExtBcDuration + e.extBcDoneMargin.MaxEstimate()
}

func (e *estimators) AgFastPathEst() time.Duration {
	// The unanimity optimization lowers convergence to the time for one message broadcast (per node):
	// the INPUT message
	// One abba round takes roughly 3 message broadcasts (per node), +1 for the common coin
	abbaEst := e.abbaRoundNoCoinDuration.MinEstimate() / 3

	// VCB is supposed to take roughly ~2 broadcasts to deliver (from the POV of the leader or follower)
	vcbEstL := e.maxOwnBcLocalDuration / 2
	vcbEstF := e.maxExtBcDuration / 2

	if abbaEst < vcbEstL && abbaEst < vcbEstF {
		return abbaEst
	} else if vcbEstF < vcbEstL {
		return vcbEstF
	}
	return vcbEstL
}

func (e *estimators) BcRuntime(slot bcpbtypes.Slot) (time.Duration, bool) {
	if startTime, ok := e.bcStartTimes[slot]; ok {
		return time.Since(startTime), true
	}

	return 0, false
}

func (e *estimators) MarkBcStartedNow(slot bcpbtypes.Slot) {
	e.bcStartTimes[slot] = time.Now()
}

func newEstimators(m dsl.Module, params ModuleParams, tunables ModuleTunables, nodeID t.NodeID) *estimators {
	allNodes := maputil.GetSortedKeys(params.Membership.Nodes)
	N := len(allNodes)
	ownQueueIdx := aleatypes.QueueIdx(slices.Index(allNodes, nodeID))

	est := &estimators{
		bcStartTimes: make(map[bcpbtypes.Slot]time.Time, (N-1)*tunables.MaxConcurrentVcbPerQueue+tunables.MaxOwnUnagreedBatchCount),

		extBcDoneMargin: util.NewByzEstimator(tunables.EstimateWindowSize, N),

		abbaRoundNoCoinDuration: util.NewEstimator(tunables.EstimateWindowSize),
	}

	bcpbdsl.UponEstimateUpdate(m, func(maxOwnBcDuration, maxOwnBcLocalDuration, maxExtBcDuration time.Duration) error {
		est.maxOwnBcDuration = maxOwnBcDuration
		est.maxOwnBcLocalDuration = maxOwnBcLocalDuration
		est.maxExtBcDuration = maxExtBcDuration
		return nil
	})

	// =============================================================================================
	// Bc Runtime Tracking
	// =============================================================================================
	bcpbdsl.UponBcStarted(m, func(slot *bcpbtypes.Slot) error {
		if _, ok := est.bcStartTimes[*slot]; !ok {
			est.bcStartTimes[*slot] = time.Now()
		}

		return nil
	})
	bcpbdsl.UponDeliverCert(m, func(cert *bcpbtypes.Cert) error {
		slot := *cert.Slot
		delete(est.bcStartTimes, slot)
		return nil
	})

	// =============================================================================================
	// External Bc Finish Done Margin Estimation
	// =============================================================================================
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
