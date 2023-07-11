package formbatchesint

import (
	"math/rand"
	"time"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/clientprogress"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool/common"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	mpdsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	mpevents "github.com/filecoin-project/mir/pkg/pb/mempoolpb/events"
	mppbtypes "github.com/filecoin-project/mir/pkg/pb/mempoolpb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/serializing"
	timertypes "github.com/filecoin-project/mir/pkg/timer/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
)

type State struct {
	*common.State

	// Transactions received, but not yet emitted in any batch.
	NewTxIDsBuckets [][]tt.TxID

	// Combined total payload size of all the transactions in the mempool.
	TotalPayloadSize int

	// Pending batch requests, i.e., batch requests that have not yet been satisfied.
	// They are indexed by order of arrival.
	PendingBatchRequests map[int]*mppbtypes.RequestBatchOrigin

	// Index of the oldest pending batch request.
	FirstPendingBatchReqID int

	// Index of the oldest pending batch request whose timer has not expired.
	// Critical requests are those whose timer has expired, but for which not enough transactions were
	// received to form a batch.
	FirstPendingNonCriticalBatchReqID int

	// Index to assign to the next new pending batch request.
	NextPendingBatchReqID int

	bucketRng      rand.Rand
	pendingTxCount int

	clientProgress *clientprogress.ClientProgress
}

// IncludeBatchCreation registers event handlers for processing NewRequests and RequestBatch events.
func IncludeBatchCreation(
	m dsl.Module,
	mc common.ModuleConfig,
	params *common.ModuleParams,
	commonState *common.State,
	logger logging.Logger,
) {
	state := &State{
		State: commonState,

		NewTxIDsBuckets: make([][]tt.TxID, params.IncomingTxBucketCount),

		TotalPayloadSize:                  0,
		PendingBatchRequests:              make(map[int]*mppbtypes.RequestBatchOrigin),
		FirstPendingBatchReqID:            0,
		FirstPendingNonCriticalBatchReqID: 0,
		NextPendingBatchReqID:             0,

		bucketRng:      *rand.New(rand.NewSource(params.RandSeed)), // nolint: gosec
		clientProgress: clientprogress.NewClientProgress(logging.NilLogger),
	}

	// cutBatch creates a new batch from whatever transactions are available in the mempool
	// (even if the batch ends up empty), and emits it as an event associated with the given origin.
	cutBatch := func(origin *mppbtypes.RequestBatchOrigin) {
		var txIDs []tt.TxID
		var txs []*trantorpbtypes.Transaction

		batchSize := 0
		txCount := 0

		bucketStartIdx := state.bucketRng.Intn(len(state.NewTxIDsBuckets))
	BatchFillLoop:
		for bucketOffset := 0; bucketOffset < len(state.NewTxIDsBuckets); bucketOffset++ {
			bucketIdx := (bucketStartIdx + bucketOffset) % len(state.NewTxIDsBuckets)
			bucket := &state.NewTxIDsBuckets[bucketIdx]
			bucketTxCount := 0

			for _, txID := range *bucket {
				tx := state.TxByID[txID]

				// Stop adding TXs if count or size limit has been reached.
				if txCount == params.MaxTransactionsInBatch || batchSize+len(tx.Data) > params.MaxPayloadInBatch {
					break BatchFillLoop
				}

				txIDs = append(txIDs, txID)
				txs = append(txs, tx)
				batchSize += len(tx.Data)
				txCount++

				bucketTxCount++
			}

			for _, txID := range (*bucket)[:bucketTxCount] {
				delete(state.TxByID, txID)
			}

			*bucket = (*bucket)[bucketTxCount:]
		}

		state.TotalPayloadSize -= batchSize
		state.pendingTxCount -= txCount

		if len(txs) < params.MinTransactionsInBatch {
			panic(es.Errorf("batch too small: %d", len(txs)))
		}

		mpdsl.NewBatch(m, origin.Module, txIDs, txs, origin)
	}

	mpdsl.UponMarkDelivered(m, func(txs []*trantorpbtypes.Transaction) error {
		for _, tx := range txs {
			state.clientProgress.Add(tx.ClientId, tx.TxNo)
		}

		mpdsl.RequestTransactionIDs[markDeliveredContext](m, mc.Self, txs, nil)
		return nil
	})
	mpdsl.UponTransactionIDsResponse(m, func(txIDs []tt.TxID, context *markDeliveredContext) error {
		for _, txID := range txIDs {
			tx, ok := state.TxByID[txID]
			if !ok {
				continue
			}

			state.TotalPayloadSize -= len(tx.Data)
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

	// storePendingRequest creates an entry for a new batch request in the pending batch request list.
	// It generates a unique ID for the pending request that can be used to associate a timeout with it.
	storePendingRequest := func(origin *mppbtypes.RequestBatchOrigin) int {
		state.PendingBatchRequests[state.NextPendingBatchReqID] = origin
		state.NextPendingBatchReqID++
		return state.NextPendingBatchReqID - 1
	}

	// Returns true if the mempool contains enough transactions for a full batch.
	haveFullBatch := func() bool {
		return state.pendingTxCount >= params.MaxTransactionsInBatch ||
			state.TotalPayloadSize >= params.MaxPayloadInBatch
	}

	// Returns true if the mempool contains enough transactions for a min-sized batch.
	haveMinBatch := func() bool {
		return state.pendingTxCount >= params.MinTransactionsInBatch
	}

	// Cuts a new batch for a batch request with the given ID and updates the corresponding internal data structures.
	servePendingReq := func(batchReqID int) {

		// Serve the given pending batch request and remove it from the pending list.
		cutBatch(state.PendingBatchRequests[batchReqID])
		delete(state.PendingBatchRequests, batchReqID)

		// Advance the pointer to the first pending batch request.
		for _, ok := state.PendingBatchRequests[state.FirstPendingBatchReqID]; !ok && state.FirstPendingBatchReqID < state.NextPendingBatchReqID; _, ok = state.PendingBatchRequests[state.FirstPendingBatchReqID] {
			state.FirstPendingBatchReqID++
		}

		if state.FirstPendingBatchReqID >= state.FirstPendingNonCriticalBatchReqID {
			state.FirstPendingNonCriticalBatchReqID = state.FirstPendingBatchReqID
		}
	}

	mpdsl.UponNewTransactions(m, func(txs []*trantorpbtypes.Transaction) error {
		filteredTxs := make([]*trantorpbtypes.Transaction, 0, len(txs))
		for _, tx := range txs {
			if state.clientProgress.CanAdd(tx.ClientId, tx.TxNo) {
				filteredTxs = append(filteredTxs, tx)
			}
		}

		mpdsl.RequestTransactionIDs(m, mc.Self, filteredTxs, &requestTxIDsContext{filteredTxs})
		return nil
	})

	mpdsl.UponTransactionIDsResponse(m, func(txIDs []tt.TxID, context *requestTxIDsContext) error {
		for i, txID := range txIDs {
			if _, ok := state.TxByID[txIDs[i]]; !ok {
				tx := context.txs[i]

				// Discard transactions with payload larger than batch limit
				// (as they would not fit in any batch, even if no other transactions were present).
				if len(tx.Data) > params.MaxPayloadInBatch {
					logger.Log(logging.LevelWarn, "Discarding transaction. Payload larger than batch limit.",
						"MaxPayloadInBatch", params.MaxPayloadInBatch, "PayloadSize", len(tx.Data))
					continue
				}

				// discard old txs again (in case they were delivered while we were computing tx IDs)
				if !state.clientProgress.CanAdd(tx.ClientId, tx.TxNo) {
					continue
				}

				state.TxByID[txID] = tx
				state.TotalPayloadSize += len(tx.Data)

				// distribute txs among buckets
				bucketIdx := txBucketIdx(len(state.NewTxIDsBuckets), txID)
				state.NewTxIDsBuckets[bucketIdx] = append(state.NewTxIDsBuckets[bucketIdx], txID)

				state.pendingTxCount++
			}
		}

		// ensure critical requests are served, even if not enough batches are yet available
		for haveMinBatch() && state.FirstPendingBatchReqID < state.FirstPendingNonCriticalBatchReqID {
			servePendingReq(state.FirstPendingBatchReqID)
		}

		for haveFullBatch() && state.FirstPendingBatchReqID < state.NextPendingBatchReqID {
			servePendingReq(state.FirstPendingBatchReqID)
		}
		return nil
	})

	mpdsl.UponRequestBatch(m, func(timeout time.Duration, origin *mppbtypes.RequestBatchOrigin) error {
		if haveFullBatch() {
			cutBatch(origin)
		} else {
			reqID := storePendingRequest(origin)
			eventpbdsl.TimerDelay(m,
				mc.Timer,
				[]*eventpbtypes.Event{mpevents.BatchTimeout(mc.Self, uint64(reqID))},
				timertypes.Duration(timeout),
			)
		}
		return nil
	})

	mpdsl.UponBatchTimeout(m, func(batchReqID uint64) error {

		reqID := int(batchReqID)

		// Load the request origin.
		_, ok := state.PendingBatchRequests[reqID]

		if ok {
			if haveMinBatch() {
				// If request is still pending, respond to it.
				servePendingReq(reqID)
			} else {
				// If a request is still pending, but we still don't have enough transactions,
				// mark the request as critical.
				if state.FirstPendingNonCriticalBatchReqID <= reqID {
					// Note: we assume all prior requests to also be critical.
					// The timeout is the same for all, so this is a safe assumption.
					state.FirstPendingNonCriticalBatchReqID = reqID + 1
				}
			}
		} else {
			// Ignore timeout if request has already been served.
			logger.Log(logging.LevelDebug, "Ignoring outdated batch timeout.",
				"batchReqID", reqID)
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
