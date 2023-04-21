package formbatchesint

import (
	"math/rand"

	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool/internal/common"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	mpdsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	mppbtypes "github.com/filecoin-project/mir/pkg/pb/mempoolpb/types"
	requestpbtypes "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type State struct {
	*common.State

	NewTxIDsBuckets [][]t.TxID
	bucketRng       rand.Rand
	pendingTxCount  int

	pendingBatchRequests []*mppbtypes.RequestBatchOrigin
}

// IncludeBatchCreation registers event handlers for processing NewRequests and RequestBatch events.
func IncludeBatchCreation(
	m dsl.Module,
	mc *common.ModuleConfig,
	params *common.ModuleParams,
	commonState *common.State,
) {
	state := &State{
		State: commonState,

		NewTxIDsBuckets: make([][]t.TxID, params.IncomingTxBucketCount),
		bucketRng:       *rand.New(rand.NewSource(params.RandSeed)),

		pendingBatchRequests: nil,
	}

	eventpbdsl.UponNewRequests(m, func(txs []*requestpbtypes.Request) error {
		mpdsl.RequestTransactionIDs(m, mc.Self, txs, &requestTxIDsContext{txs})
		return nil
	})

	mpdsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *requestTxIDsContext) error {
		for i := range txIDs {
			if _, ok := state.TxByID[string(txIDs[i])]; !ok {
				id := txIDs[i]

				state.TxByID[string(id)] = context.txs[i]

				// distribute txs among buckets
				bucketIdx := txBucketIdx(len(state.NewTxIDsBuckets), id)
				state.NewTxIDsBuckets[bucketIdx] = append(state.NewTxIDsBuckets[bucketIdx], id)

				state.pendingTxCount++
			}
		}
		return nil
	})

	mpdsl.UponRequestBatch(m, func(origin *mppbtypes.RequestBatchOrigin) error {
		state.pendingBatchRequests = append(state.pendingBatchRequests, origin)
		return nil
	})

	dsl.UponCondition(m, func() error {
		for len(state.pendingBatchRequests) > 0 && state.pendingTxCount >= params.MinTransactionsInBatch {
			var txIDs []t.TxID
			var txs []*requestpbtypes.Request
			batchTxCount := 0

			startingIdx := state.bucketRng.Intn(len(state.NewTxIDsBuckets))
			for i := 0; i < len(state.NewTxIDsBuckets) && batchTxCount < params.MaxTransactionsInBatch; i++ {
				bucket := &state.NewTxIDsBuckets[(startingIdx+i)%len(state.NewTxIDsBuckets)]

				txCount := 0
				for _, txID := range *bucket {
					tx := state.TxByID[string(txID)]

					// TODO: add other limitations (if any) here.
					if batchTxCount == params.MaxTransactionsInBatch {
						break
					}

					txIDs = append(txIDs, txID)
					txs = append(txs, tx)
					batchTxCount++
					txCount++
				}

				for _, txID := range (*bucket)[:txCount] {
					delete(state.TxByID, string(txID))
				}

				*bucket = (*bucket)[txCount:]
			}

			origin := state.pendingBatchRequests[0]
			state.pendingBatchRequests = state.pendingBatchRequests[1:]
			state.pendingTxCount -= batchTxCount

			mpdsl.NewBatch(m, origin.Module, txIDs, txs, origin)
		}

		return nil
	})

	mpdsl.UponMarkDelivered(m, func(txs []*requestpbtypes.Request) error {
		mpdsl.RequestTransactionIDs[markDeliveredContext](m, mc.Self, txs, nil)
		return nil
	})
	mpdsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *markDeliveredContext) error {
		for _, txID := range txIDs {
			if _, ok := state.TxByID[string(txID)]; !ok {
				continue
			}

			state.pendingTxCount--
			delete(state.TxByID, string(txID))

			bucketIdx := txBucketIdx(len(state.NewTxIDsBuckets), txID)
			for i, id := range state.NewTxIDsBuckets[bucketIdx] {
				if slices.Equal(txID, id) {
					bucket := state.NewTxIDsBuckets[bucketIdx]
					// TODO: may be better served by a different data structure
					state.NewTxIDsBuckets[bucketIdx] = append(bucket[:i], bucket[i+1:]...)
				}
			}

		}

		return nil
	})
}

// Context data structures

type requestTxIDsContext struct {
	txs []*requestpbtypes.Request
}

type markDeliveredContext struct{}

func txBucketIdx(nBuckets int, txID t.TxID) int {
	return int(t.Uint64FromBytes(txID[len(txID)-8:]) % uint64(nBuckets))
}
