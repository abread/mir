package availability

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/alea/broadcast"
	bcdsl "github.com/filecoin-project/mir/pkg/alea/director/internal/aleadsl"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/protobuf"
	batchdbdsl "github.com/filecoin-project/mir/pkg/availability/batchdb/dsl"
	adsl "github.com/filecoin-project/mir/pkg/availability/dsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	mempooldsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	aleapbCommon "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	threshDsl "github.com/filecoin-project/mir/pkg/threshcrypto/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb"
)

// State represents the state related to this part of the module.
type batchFetchingState struct {
	RequestsState map[batchSlot]*RequestsState
}

// comparable aleapb.common.Slot variant
type batchSlot struct {
	QueueIdx  uint32
	QueueSlot uint64
}

// RequestsState represents the state related to a request on the source node of the request.
// The node disposes of this state as soon as the request is completed.
type RequestsState struct {
	ReqOrigins []*apb.RequestTransactionsOrigin
	Replies    map[t.NodeID]struct{}
}

// IncludeBatchFetching registers event handlers for processing availabilitypb.RequestTransactions events.
func IncludeBatchFetching(
	m dsl.Module,
	mc *common.ModuleConfig,
	params *common.ModuleParams,
	nodeID t.NodeID,
) {
	state := batchFetchingState{
		RequestsState: make(map[batchSlot]*RequestsState),
	}

	// When receive a request for transactions, first check the local storage.
	bcdsl.UponRequestTransactions(m, func(cert *aleapb.Cert, origin *apb.RequestTransactionsOrigin) error {

		// NOTE: it is assumed that cert is valid.
		batchdbdsl.LookupBatch(m, mc.BatchDB, formatAleaBatchID(cert.Slot), &lookupBatchLocallyContext{cert, origin})
		return nil
	})

	// If the batch is present in the local storage, return it. Otherwise, ask other nodes.
	batchdbdsl.UponLookupBatchResponse(m, func(found bool, txs []*requestpb.Request, metadata []byte, context *lookupBatchLocallyContext) error {
		if found {
			adsl.ProvideTransactions(m, t.ModuleID(context.origin.Module), txs, context.origin)
			return nil
		}

		slot := batchSlotFromPb(context.cert.Slot)
		if _, present := state.RequestsState[slot]; !present {
			state.RequestsState[slot] = &RequestsState{
				ReqOrigins: make([]*apb.RequestTransactionsOrigin, 0, 1),
				Replies:    make(map[t.NodeID]struct{}, len(params.AllNodes)),
			}
		}
		reqState := state.RequestsState[slot]
		reqState.ReqOrigins = append(reqState.ReqOrigins, context.origin)

		// TODO: do this more inteligently: only contact some nodes, and try others if a timer expires
		// until all were tried or a response is received.
		// It would also be nice to pass a hint in the certificate that says which nodes to try first,
		// this could be provided by the agreement component based on INIT(v, 0) messages received by the abba instances.
		dsl.SendMessage(m, mc.Net,
			protobuf.FillGapMessage(mc.Self, slot.Pb()),
			params.AllNodes)
		return nil
	})

	// When receive a request for batch from another node, lookup the batch in the local storage.
	bcdsl.UponFillGapMessageReceived(m, func(from t.NodeID, msg *aleapb.FillGapMessage) error {
		slot := batchSlotFromPb(msg.Slot)

		batchdbdsl.LookupBatch(m, mc.BatchDB, formatAleaBatchID(msg.Slot), &lookupBatchOnRemoteRequestContext{from, slot})
		return nil
	})

	// If the batch is found in the local storage, send it to the requesting node.
	batchdbdsl.UponLookupBatchResponse(m, func(found bool, txs []*requestpb.Request, signature []byte, context *lookupBatchOnRemoteRequestContext) error {
		if !found {
			// Ignore invalid request.
			return nil
		}

		dsl.SendMessage(m, mc.Net, protobuf.FillerMessage(mc.Self, context.slot.Pb(), txs, signature), []t.NodeID{context.requester})
		return nil
	})

	// When receive a requested batch, compute the ids of the received transactions.
	bcdsl.UponFillerMessageReceived(m, func(from t.NodeID, msg *aleapb.FillerMessage) error {
		slot := batchSlotFromPb(msg.Slot)
		reqState, present := state.RequestsState[slot]
		if !present {
			return nil // no request needs this message to be satisfied
		}

		if _, alreadyAnswered := reqState.Replies[from]; alreadyAnswered {
			return nil // already processed a reply from this node
		}
		reqState.Replies[from] = struct{}{}

		mempooldsl.RequestTransactionIDs(m, mc.Mempool, msg.Txs, &handleFillerContext{
			slot:      slot,
			txs:       msg.Txs,
			signature: msg.Signature,
		})
		return nil
	})

	// When transaction ids are computed, compute the batchID
	mempooldsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *handleFillerContext) error {
		context.txIDs = txIDs
		mempooldsl.RequestBatchID(m, mc.Mempool, txIDs, context)
		return nil
	})

	// When the id of the batch is computed, check if the signature is correct
	mempooldsl.UponBatchIDResponse(m, func(batchID t.BatchID, context *handleFillerContext) error {
		_, ok := state.RequestsState[context.slot]
		if !ok {
			// The request has already been completed.
			return nil
		}

		context.batchID = batchID
		sigData := certSigData(params, context.slot, batchID)
		threshDsl.VerifyFull(m, mc.ThreshCrypto, sigData, context.signature, context)
		return nil
	})

	threshDsl.UponVerifyFullResult(m, func(ok bool, err string, context *handleFillerContext) error {
		if !ok {
			// TODO: do this the smart way to avoid needless traffic and send requests to other nodes here
			// also go ahead and ensure that no request goes unanswered
			return nil
		}

		requestState, ok := state.RequestsState[context.slot]
		if !ok {
			// The request has already been completed.
			return nil
		}

		// store batch asynchronously
		batchdbdsl.StoreBatch(m, mc.BatchDB, context.batchID, context.txIDs, context.txs, context.signature /*metadata*/, &storeBatchContext{})

		// send response to requests
		for _, origin := range requestState.ReqOrigins {
			adsl.ProvideTransactions(m, t.ModuleID(origin.Module), context.txs, origin)
		}
		delete(state.RequestsState, context.slot)

		return nil
	})

	batchdbdsl.UponBatchStored(m, func(_ *storeBatchContext) error {
		// do nothing.
		return nil
	})
}

func formatAleaBatchID(slot *aleapbCommon.Slot) t.BatchID {
	return t.BatchID(fmt.Sprintf("alea-%d:%d", slot.QueueIdx, slot.QueueSlot))
}

func certSigData(params *common.ModuleParams, slot batchSlot, batchID t.BatchID) [][]byte {
	return vcb.SigData(broadcast.VCBInstanceUID(params.InstanceUID, slot.QueueIdx, slot.QueueSlot), batchID)
}

func batchSlotFromPb(pb *aleapbCommon.Slot) batchSlot {
	return batchSlot{
		QueueIdx:  pb.QueueIdx,
		QueueSlot: pb.QueueSlot,
	}
}

func (slot *batchSlot) Pb() *aleapbCommon.Slot {
	return &aleapbCommon.Slot{
		QueueIdx:  slot.QueueIdx,
		QueueSlot: slot.QueueSlot,
	}
}

// Context data structures

type lookupBatchLocallyContext struct {
	cert   *aleapb.Cert
	origin *apb.RequestTransactionsOrigin
}

type lookupBatchOnRemoteRequestContext struct {
	requester t.NodeID
	slot      batchSlot
}

type handleFillerContext struct {
	slot      batchSlot
	txs       []*requestpb.Request
	signature []byte

	txIDs   []t.TxID
	batchID t.BatchID
}

type storeBatchContext struct{}
