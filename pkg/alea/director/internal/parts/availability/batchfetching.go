package availability

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/alea/broadcast"
	"github.com/filecoin-project/mir/pkg/alea/common"
	director "github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/protobuf"
	batchdbdsl "github.com/filecoin-project/mir/pkg/availability/batchdb/dsl"
	adsl "github.com/filecoin-project/mir/pkg/availability/dsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	mempooldsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	abcdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	bcdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/dsl"
	aleamsgs "github.com/filecoin-project/mir/pkg/pb/aleapb/msgs"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	rnetdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	threshDsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb"
)

// State represents the state related to this part of the module.
type batchFetchingState struct {
	RequestsState map[commontypes.Slot]*RequestsState
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
		RequestsState: make(map[commontypes.Slot]*RequestsState),
	}

	// When receive a request for transactions, first check the local storage.
	adsl.UponRequestTransactions(m, func(anyCert *apb.Cert, origin *apb.RequestTransactionsOrigin) error {
		certWrapper, ok := anyCert.Type.(*apb.Cert_Alea)
		if !ok {
			return fmt.Errorf("unexpected certificate type. Expected: %T, got %T", certWrapper, anyCert.Type)
		}
		cert := certWrapper.Alea
		slot := commontypes.SlotFromPb(cert.Slot)

		// NOTE: it is assumed that cert is valid.
		batchdbdsl.LookupBatch(m, mc.BatchDB, common.FormatAleaBatchID(slot), &lookupBatchLocallyContext{slot, origin})
		return nil
	})

	// If the batch is present in the local storage, return it. Otherwise, ask other nodes.
	batchdbdsl.UponLookupBatchResponse(m, func(found bool, txs []*requestpb.Request, metadata []byte, context *lookupBatchLocallyContext) error {
		if found {
			adsl.ProvideTransactions(m, t.ModuleID(context.origin.Module), txs, context.origin)
			return nil
		}

		slot := *context.slot
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
			FillGapMsgID(&slot),
			aleamsgs.FillGapMessage(mc.Self, &slot),
			params.AllNodes,
		)
		return nil
	})

	// When receive a request for batch from another node, lookup the batch in the local storage.
	bcdsl.UponFillGapMessageReceived(m, func(from t.NodeID, slot *commontypes.Slot) error {
		// do not ACK message - acknowledging means sending a FILLER reply

		logger.Log(logging.LevelDebug, "satisfying FILL-GAP request", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot, "from", from)
		batchdbdsl.LookupBatch(m, mc.BatchDB, common.FormatAleaBatchID(slot), &lookupBatchOnRemoteRequestContext{from, slot})
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
	bcdsl.UponFillerMessageReceived(m, func(from t.NodeID, slot *commontypes.Slot, txs []*requestpb.Request, signature tctypes.FullSig) error {
		rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, FillGapMsgID(slot), []t.NodeID{from})

		reqState, present := state.RequestsState[*slot]
		if !present {
			return nil // no request needs this message to be satisfied
		}

		if _, alreadyAnswered := reqState.Replies[from]; alreadyAnswered {
			return nil // already processed a reply from this node
		}
		reqState.Replies[from] = struct{}{}

		logger.Log(logging.LevelDebug, "got FILLER for missing slot!", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot, "from", from)

		mempooldsl.RequestTransactionIDs(m, mc.Mempool, txs, &handleFillerContext{
			slot:      slot,
			txs:       txs,
			signature: signature,
		})
		return nil
	})

	// When transaction ids are computed, check if the signature is correct
	mempooldsl.UponTransactionIDsResponse(m, func(txIDs []t.TxID, context *handleFillerContext) error {
		_, ok := state.RequestsState[*context.slot]
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

		requestState, ok := state.RequestsState[*context.slot]
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
		delete(state.RequestsState, *context.slot)

		return nil
	})

	batchdbdsl.UponBatchStored(m, func(context *handleFillerContext) error {
		abcdsl.FreeSlot(m, mc.AleaBroadcast, context.slot)
		return nil
	})
}

func certSigData(params *director.ModuleParams, slot *commontypes.Slot, txIDs []t.TxID) [][]byte {
	aleaUID := params.InstanceUID[:len(params.InstanceUID)-1]
	return vcb.SigData(broadcast.VCBInstanceUID(aleaUID, slot.QueueIdx, slot.QueueSlot), txIDs)
}

const (
	MsgTypeFillGap = "f"
)

func FillGapMsgID(slot *commontypes.Slot) rntypes.MsgID {
	return rntypes.MsgID(fmt.Sprintf("%s.%d.%d", MsgTypeFillGap, slot.QueueIdx, slot.QueueSlot))
}

// Context data structures

type lookupBatchLocallyContext struct {
	slot   *commontypes.Slot
	origin *apb.RequestTransactionsOrigin
}

type lookupBatchOnRemoteRequestContext struct {
	requester t.NodeID
	slot      *commontypes.Slot
}

type handleFillerContext struct {
	slot      *commontypes.Slot
	txs       []*requestpb.Request
	signature []byte

	txIDs   []t.TxID
	batchID t.BatchID
}
