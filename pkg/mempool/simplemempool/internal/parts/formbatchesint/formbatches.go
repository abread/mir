package formbatchesint

import (
	"math/rand"

	"github.com/filecoin-project/mir/pkg/clientprogress"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool/common"
	mpdsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	mppbtypes "github.com/filecoin-project/mir/pkg/pb/mempoolpb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/serializing"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
)

type State struct {
	*common.State

	NewTxIDsBuckets [][]tt.TxID
	bucketRng       rand.Rand
	pendingTxCount  int

	pendingBatchRequests []*mppbtypes.RequestBatchOrigin

	clientProgress *clientprogress.ClientProgress
}

// IncludeBatchCreation registers event handlers for processing NewRequests and RequestBatch events.
func IncludeBatchCreation(
	m dsl.Module,
	mc common.ModuleConfig,
	params *common.ModuleParams,
	commonState *common.State,
) {
	state := &State{
		State: commonState,

		NewTxIDsBuckets: make([][]tt.TxID, params.IncomingTxBucketCount),
		bucketRng:       *rand.New(rand.NewSource(params.RandSeed)), // nolint: gosec

		pendingBatchRequests: nil,
		clientProgress:       clientprogress.NewClientProgress(logging.NilLogger),
	}

	mpdsl.UponNewTransactions(m, func(txs []*trantorpbtypes.Transaction) error {
		filteredTxs := make([]*trantorpbtypes.Transaction, 0, len(txs))
		for _, tx := range txs {
			// TODO: can we use Add here? depends on whether we can trust incoming transactions to be valid
			if state.clientProgress.CanAdd(tx.ClientId, tx.TxNo) {
				filteredTxs = append(filteredTxs, tx)
			}
		}

		mpdsl.RequestTransactionIDs(m, mc.Self, filteredTxs, &requestTxIDsContext{filteredTxs})
		return nil
	})

	mpdsl.UponTransactionIDsResponse(m, func(txIDs []tt.TxID, context *requestTxIDsContext) error {
		for i := range txIDs {
			if _, ok := state.TxByID[txIDs[i]]; !ok {
				id := txIDs[i]

				state.TxByID[id] = context.txs[i]

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

	dsl.UponStateUpdates(m, func() error {
		for len(state.pendingBatchRequests) > 0 && state.pendingTxCount >= params.MinTransactionsInBatch {
			var txIDs []tt.TxID
			var txs []*trantorpbtypes.Transaction
			batchTxCount := 0

			startingIdx := state.bucketRng.Intn(len(state.NewTxIDsBuckets))
			for i := 0; i < len(state.NewTxIDsBuckets) && batchTxCount < params.MaxTransactionsInBatch; i++ {
				bucket := &state.NewTxIDsBuckets[(startingIdx+i)%len(state.NewTxIDsBuckets)]

				txCount := 0
				for _, txID := range *bucket {
					tx := state.TxByID[txID]

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

	mpdsl.UponMarkDelivered(m, func(txs []*trantorpbtypes.Transaction) error {
		for _, tx := range txs {
			state.clientProgress.Add(tx.ClientId, tx.TxNo)
		}

		mpdsl.RequestTransactionIDs[markDeliveredContext](m, mc.Self, txs, nil)
		return nil
	})
	mpdsl.UponTransactionIDsResponse(m, func(txIDs []tt.TxID, context *markDeliveredContext) error {
		for _, txID := range txIDs {
			if _, ok := state.TxByID[txID]; !ok {
				continue
			}

			state.pendingTxCount--
			delete(state.TxByID, txID)

			bucketIdx := txBucketIdx(len(state.NewTxIDsBuckets), txID)
			for i, id := range state.NewTxIDsBuckets[bucketIdx] {
				if txID == id {
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
	txs []*trantorpbtypes.Transaction
}

type markDeliveredContext struct{}

func txBucketIdx(nBuckets int, txID tt.TxID) int {
	return int(serializing.Uint64FromBytes([]byte(txID[len(txID)-8:])) % uint64(nBuckets))
}