package formbatches

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	mpdsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool/internal/common"
	mppb "github.com/filecoin-project/mir/pkg/pb/mempoolpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type State struct {
	*common.State
	NewTxIDs []t.TxID

	pendingBatchRequests []*mppb.RequestBatchOrigin
}

// IncludeBatchCreation registers event handlers for processing NewRequests and RequestBatch events.
func IncludeBatchCreation(
	m dsl.Module,
	mc *common.ModuleConfig,
	params *common.ModuleParams,
	commonState *common.State,
) {
	state := &State{
		State:    commonState,
		NewTxIDs: nil,

		pendingBatchRequests: nil,
	}

	dsl.UponNewRequests(m, func(txs []*requestpb.Request) error {
		mpdsl.RequestTransactionIDs(m, mc.Self, txs, &requestTxIDsContext{txs})
		return nil
	})

	mpdsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *requestTxIDsContext) error {
		for i := range txIDs {
			state.TxByID[txIDs[i]] = context.txs[i]
		}
		state.NewTxIDs = append(state.NewTxIDs, txIDs...)
		return nil
	})

	mpdsl.UponRequestBatch(m, func(origin *mppb.RequestBatchOrigin) error {
		state.pendingBatchRequests = append(state.pendingBatchRequests, origin)
		return nil
	})

	dsl.UponCondition(m, func() error {
		for len(state.pendingBatchRequests) > 0 && len(state.NewTxIDs) >= params.MinTransactionsInBatch {
			var txIDs []t.TxID
			var txs []*requestpb.Request
			batchSize := 0

			txCount := 0
			for _, txID := range state.NewTxIDs {
				tx := state.TxByID[txID]

				// TODO: add other limitations (if any) here.
				if txCount == params.MaxTransactionsInBatch {
					break
				}

				txIDs = append(txIDs, txID)
				txs = append(txs, tx)
				batchSize += len(tx.Data)
				txCount++
			}

			state.NewTxIDs = state.NewTxIDs[txCount:]

			origin := state.pendingBatchRequests[0]
			state.pendingBatchRequests = state.pendingBatchRequests[1:]

			mpdsl.NewBatch(m, t.ModuleID(origin.Module), txIDs, txs, origin)
		}

		return nil
	})
}

// Context data structures

type requestTxIDsContext struct {
	txs []*requestpb.Request
}
