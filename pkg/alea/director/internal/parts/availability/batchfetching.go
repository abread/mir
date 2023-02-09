package availability

import (
	"github.com/filecoin-project/mir/pkg/alea/broadcast"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/abcdsl"
	"github.com/filecoin-project/mir/pkg/alea/common"
	bcdsl "github.com/filecoin-project/mir/pkg/alea/director/internal/aleadsl"
	director "github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/protobuf"
	batchdbdsl "github.com/filecoin-project/mir/pkg/availability/batchdb/dsl"
	adsl "github.com/filecoin-project/mir/pkg/availability/dsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	mempooldsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	aleapbCommon "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	threshDsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	"github.com/filecoin-project/mir/pkg/reliablenet/rnetdsl"
	"github.com/filecoin-project/mir/pkg/serializing"
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
	mc *director.ModuleConfig,
	params *director.ModuleParams,
	nodeID t.NodeID,
	logger logging.Logger,
) {
	state := batchFetchingState{
		RequestsState: make(map[batchSlot]*RequestsState),
	}

	// When receive a request for transactions, first check the local storage.
	bcdsl.UponRequestTransactions(m, func(cert *aleapb.Cert, origin *apb.RequestTransactionsOrigin) error {

		// NOTE: it is assumed that cert is valid.
		batchdbdsl.LookupBatch(m, mc.BatchDB, common.FormatAleaBatchID(cert.Slot), &lookupBatchLocallyContext{cert, origin})
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
		logger.Log(logging.LevelDebug, "broadcast component fell behind. requesting slot from other replicas with FILL-GAP", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)
		rnetdsl.SendMessage(m, mc.ReliableNet,
			FillGapMsgID(slot),
			protobuf.FillGapMessage(mc.Self, slot.Pb()),
			params.AllNodes,
		)
		return nil
	})

	// When receive a request for batch from another node, lookup the batch in the local storage.
	bcdsl.UponFillGapMessageReceived(m, func(from t.NodeID, msg *aleapb.FillGapMessage) error {
		slot := batchSlotFromPb(msg.Slot)
		// do not ACK message - acknowledging means sending a FILLER reply

		logger.Log(logging.LevelDebug, "satisfying FILL-GAP request", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot, "from", from)
		batchdbdsl.LookupBatch(m, mc.BatchDB, common.FormatAleaBatchID(msg.Slot), &lookupBatchOnRemoteRequestContext{from, slot})
		return nil
	})

	// If the batch is found in the local storage, send it to the requesting node.
	batchdbdsl.UponLookupBatchResponse(m, func(found bool, txs []*requestpb.Request, signature []byte, context *lookupBatchOnRemoteRequestContext) error {
		if !found {
			// Ignore invalid request.
			return nil
		}

		dsl.SendMessage(m, mc.Net,
			protobuf.FillerMessage(mc.Self, context.slot.Pb(), txs, signature),
			[]t.NodeID{context.requester},
		)
		return nil
	})

	// When receive a requested batch, compute the ids of the received transactions.
	bcdsl.UponFillerMessageReceived(m, func(from t.NodeID, msg *aleapb.FillerMessage) error {
		slot := batchSlotFromPb(msg.Slot)
		rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, FillGapMsgID(slot), []t.NodeID{from})

		reqState, present := state.RequestsState[slot]
		if !present {
			return nil // no request needs this message to be satisfied
		}

		if _, alreadyAnswered := reqState.Replies[from]; alreadyAnswered {
			return nil // already processed a reply from this node
		}
		reqState.Replies[from] = struct{}{}

		logger.Log(logging.LevelDebug, "got FILLER for missing slot!", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot, "from", from)

		mempooldsl.RequestTransactionIDs(m, mc.Mempool, msg.Txs, &handleFillerContext{
			slot:      slot,
			txs:       msg.Txs,
			signature: msg.Signature,
		})
		return nil
	})

	// When transaction ids are computed, check if the signature is correct
	mempooldsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *handleFillerContext) error {
		_, ok := state.RequestsState[context.slot]
		if !ok {
			// The request has already been completed.
			return nil
		}

		context.txIDs = txIDs
		sigData := certSigData(params, context.slot, txIDs)
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

		// stop asking other nodes to send us stuff
		rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, FillGapMsgID(context.slot), params.AllNodes)

		// store batch asynchronously
		batchdbdsl.StoreBatch(m, mc.BatchDB, context.batchID, context.txIDs, context.txs, context.signature /*metadata*/, context)

		// send response to requests
		logger.Log(logging.LevelDebug, "satisfying delayed requests with FILLER", "queueIdx", context.slot.QueueIdx, "queueSlot", context.slot.QueueSlot)
		for _, origin := range requestState.ReqOrigins {
			adsl.ProvideTransactions(m, t.ModuleID(origin.Module), context.txs, origin)
		}
		delete(state.RequestsState, context.slot)

		return nil
	})

	batchdbdsl.UponBatchStored(m, func(context *handleFillerContext) error {
		abcdsl.FreeSlot(m, mc.AleaBroadcast, context.slot.Pb())
		return nil
	})
}

func certSigData(params *director.ModuleParams, slot batchSlot, txIDs []t.TxID) [][]byte {
	aleaUID := params.InstanceUID[:len(params.InstanceUID)-1]
	return vcb.SigData(broadcast.VCBInstanceUID(aleaUID, slot.QueueIdx, slot.QueueSlot), txIDs)
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

const (
	MsgTypeFillGap uint8 = iota
)

func FillGapMsgID(slot batchSlot) []byte {
	s := make([]byte, 0, 1+2*8)
	s = append(s, MsgTypeFillGap)
	s = append(s, serializing.Uint64ToBytes(uint64(slot.QueueIdx))...)
	s = append(s, serializing.Uint64ToBytes(slot.QueueSlot)...)
	return s
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