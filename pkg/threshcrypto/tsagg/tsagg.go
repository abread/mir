package tsagg

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	threshcryptopbdsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	t "github.com/filecoin-project/mir/pkg/types"
)

// TODO: document
type ThreshSigAggregator struct {
	m      dsl.Module
	params *Params

	recvdShare map[t.NodeID]struct{}

	unknownShares       []tctypes.SigShare
	unknownShareSources []t.NodeID

	okShares []tctypes.SigShare
	fullSig  tctypes.FullSig

	recoveryInProgress bool
}

func (agg *ThreshSigAggregator) FullSig() tctypes.FullSig {
	return agg.fullSig
}

func (agg *ThreshSigAggregator) Add(share tctypes.SigShare, from t.NodeID) bool {
	if _, alreadyProcessed := agg.recvdShare[from]; alreadyProcessed {
		return false
	}
	agg.recvdShare[from] = struct{}{}

	if agg.fullSig != nil {
		return true
	}

	if agg.inSlowPath() {
		if from == agg.params.OwnNodeID {
			// fast path for share verification: we trust our own signatures
			agg.okShares = append(agg.okShares, share)
			return true
		} else if data := agg.params.SigData(); data != nil {
			// sign data available: we can verify this share here
			threshcryptopbdsl.VerifyShare(agg.m, agg.params.TCModuleID, data, share, from, &verifyShareCtx{share})
			return true
		}
	}

	// save share and source for later use/verification
	agg.unknownShares = append(agg.unknownShares, share)
	agg.unknownShareSources = append(agg.unknownShareSources, from)

	return true
}

func (agg *ThreshSigAggregator) inFastPath() bool {
	return agg.okShares == nil && agg.fullSig == nil
}

func (agg *ThreshSigAggregator) inSlowPath() bool {
	return agg.okShares != nil && agg.fullSig == nil
}

type Params struct {
	OwnNodeID               t.NodeID
	TCModuleID              t.ModuleID
	Threshold               int
	MaxVerifyShareBatchSize int

	SigData func() [][]byte

	InitialNodeCount int
}

func New(m dsl.Module, params *Params, logger logging.Logger) *ThreshSigAggregator {
	agg := &ThreshSigAggregator{
		m:      m,
		params: params,

		recvdShare:    make(map[t.NodeID]struct{}, params.InitialNodeCount),
		unknownShares: make([]tctypes.SigShare, 0, params.InitialNodeCount),
	}

	dsl.UponCondition(m, func() error {
		if agg.recoveryInProgress || agg.fullSig != nil {
			return nil
		}

		data := params.SigData()
		if data == nil {
			// can't do anything yet
			return nil
		}

		if agg.inFastPath() {
			// fast path attempts recovery without verifying shares
			if len(agg.unknownShares) >= params.Threshold {
				threshcryptopbdsl.Recover(m, params.TCModuleID, data, agg.unknownShares, &fastPathRecoverCtx{})
				agg.recoveryInProgress = true
			}
		} else {
			if len(agg.okShares) >= params.Threshold {
				// enough verified shares, attempt recovery
				threshcryptopbdsl.Recover(m, params.TCModuleID, data, agg.okShares, &slowPathRecoverCtx{})
				agg.recoveryInProgress = true
			} else if len(agg.unknownShares) > 0 {
				// verify previously unverified shares
				for i := 0; i < params.MaxVerifyShareBatchSize && len(agg.unknownShares) > 0; i++ {
					share := agg.unknownShares[0]
					from := agg.unknownShareSources[0]

					if from == params.OwnNodeID {
						// fast path for verify share: we trust our own signatures
						agg.okShares = append(agg.okShares, share)
						i-- // this does not count for the batch
					} else {
						threshcryptopbdsl.VerifyShare(agg.m, agg.params.TCModuleID, data, share, from, &verifyShareCtx{share})
					}

					agg.unknownShares = agg.unknownShares[1:]
					agg.unknownShareSources = agg.unknownShareSources[1:]
				}
			}
		}

		return nil
	})

	// fast path
	threshcryptopbdsl.UponRecoverResult(m, func(fullSignature tctypes.FullSig, ok bool, error string, context *fastPathRecoverCtx) error {
		if ok {
			// must verify recovered signature
			context.fullSig = fullSignature
			threshcryptopbdsl.VerifyFull(m, params.TCModuleID, params.SigData(), fullSignature, context)
		} else {
			// move to slow path
			logger.Log(logging.LevelDebug, "Recover failed, moving to slow path", "error", error, "shareSources", agg.unknownShareSources)

			agg.okShares = make([]tctypes.SigShare, 0, params.InitialNodeCount)
			agg.recoveryInProgress = false
		}

		return nil
	})
	threshcryptopbdsl.UponVerifyFullResult(m, func(ok bool, error string, context *fastPathRecoverCtx) error {
		if ok {
			// we're done!
			agg.fullSig = context.fullSig

			// let shares be GC-ed
			agg.unknownShareSources = nil
			agg.unknownShares = nil
		} else {
			// move to slow path
			logger.Log(logging.LevelDebug, "Verification failed, moving to slow path", "error", error, "shareSources", agg.unknownShareSources)
			agg.okShares = make([]tctypes.SigShare, 0, params.InitialNodeCount)
		}

		agg.recoveryInProgress = false
		return nil
	})

	// slow path
	threshcryptopbdsl.UponVerifyShareResult(m, func(ok bool, error string, context *verifyShareCtx) error {
		if ok && agg.okShares != nil {
			agg.okShares = append(agg.okShares, context.share)
		}

		// TODO: report byz behavior? (when !ok)
		return nil
	})
	threshcryptopbdsl.UponRecoverResult(m, func(fullSignature tctypes.FullSig, ok bool, error string, context *slowPathRecoverCtx) error {
		if ok {
			agg.fullSig = fullSignature

			// let shares be GC-ed
			agg.unknownShares = nil
			agg.unknownShareSources = nil
			agg.okShares = nil
		} else {
			return fmt.Errorf("enough valid shares from different sources but could not recover signatures. isn't this impossible?")
		}

		agg.recoveryInProgress = false
		return nil
	})

	return agg
}

type fastPathRecoverCtx struct {
	fullSig tctypes.FullSig
}

type slowPathRecoverCtx struct{}

type verifyShareCtx struct {
	share tctypes.SigShare
}
