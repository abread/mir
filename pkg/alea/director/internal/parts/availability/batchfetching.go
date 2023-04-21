package availability

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/alea/broadcast/bcutil"
	"github.com/filecoin-project/mir/pkg/alea/common"
	director "github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	abcdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcpbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/events"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	aleadsl "github.com/filecoin-project/mir/pkg/pb/aleapb/dsl"
	aleamsgs "github.com/filecoin-project/mir/pkg/pb/aleapb/msgs"
	batchdbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/dsl"
	adsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/dsl"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	commonpbtypes "github.com/filecoin-project/mir/pkg/pb/commonpb/types"
	eventpbevents "github.com/filecoin-project/mir/pkg/pb/eventpb/events"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	hasherpbdsl "github.com/filecoin-project/mir/pkg/pb/hasherpb/dsl"
	mempooldsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	rnetdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	requestpbtypes "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	threshDsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	transportpbdsl "github.com/filecoin-project/mir/pkg/pb/transportpb/dsl"
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	"github.com/filecoin-project/mir/pkg/timer/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
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
	ReqOrigins  []*availabilitypbtypes.RequestTransactionsOrigin
	SentFillGap bool
	Replies     map[t.NodeID]struct{}
}

// IncludeBatchFetching registers event handlers for processing availabilitypb.RequestTransactions events.
func IncludeBatchFetching(
	m dsl.Module,
	mc *director.ModuleConfig,
	params *director.ModuleParams,
	tunables *director.ModuleTunables,
	nodeID t.NodeID,
	logger logging.Logger,
) {
	state := batchFetchingState{
		RequestsState: make(map[commontypes.Slot]*RequestsState),
	}

	// When receive a request for transactions, first check the local storage.
	adsl.UponRequestTransactions(m, func(anyCert *availabilitypbtypes.Cert, origin *availabilitypbtypes.RequestTransactionsOrigin) error {
		certWrapper, present := anyCert.Type.(*availabilitypbtypes.Cert_Alea)
		if !present {
			return fmt.Errorf("unexpected certificate type. Expected: %T, got %T", certWrapper, anyCert.Type)
		}
		// NOTE: it is assumed that cert is valid.
		cert := certWrapper.Alea
		slot := cert.Slot

		reqState, present := state.RequestsState[*slot]
		if !present {
			// new request
			reqState = &RequestsState{
				ReqOrigins: make([]*availabilitypbtypes.RequestTransactionsOrigin, 0, 1),
			}
			state.RequestsState[*slot] = reqState

			// try resolving locally first
			batchdbdsl.LookupBatch(m, mc.BatchDB, common.FormatAleaBatchID(slot), slot)
		}

		// add ourselves to the list of requestors in order to receive a reply
		reqState.ReqOrigins = append(reqState.ReqOrigins, origin)

		return nil
	})

	// if broadcast delivers for a batch being requested, we can *now* resolve it locally
	abcdsl.UponDeliver(m, func(slot *commontypes.Slot) error {
		if _, present := state.RequestsState[*slot]; present {
			// TODO: avoid concurrent lookups of the same batch?
			// if bc delivers right after the transaction starts, two lookups will be performed.

			// retry locally, you will now succeed!
			batchdbdsl.LookupBatch(m, mc.BatchDB, common.FormatAleaBatchID(slot), slot)
		}

		return nil
	})

	// If the batch is present in the local storage, return it. Otherwise, ask other nodes.
	batchdbdsl.UponLookupBatchResponse(m, func(found bool, txs []*requestpbtypes.Request, metadata []byte, slot *commontypes.Slot) error {
		reqState, ok := state.RequestsState[*slot]
		if !ok {
			return nil // stale request
		}

		if found {
			for _, origin := range reqState.ReqOrigins {
				adsl.ProvideTransactions(m, origin.Module, txs, origin)
			}
			mempooldsl.MarkDelivered(m, mc.Mempool, txs)

			if reqState.SentFillGap {
				// no need for FILLER anymore
				rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, FillGapMsgID(slot), params.AllNodes)
			}

			delete(state.RequestsState, *slot)
			return nil
		}

		if reqState.SentFillGap {
			return fmt.Errorf("tried to send FILL-GAP twice")
		}
		reqState.SentFillGap = true

		// send FILL-GAP after a timeout (if request was not satisfied)
		// this way bc has more chances of completing before even trying to send a fill-gap message
		// TODO: adjust delay according to bc estimate. don't delay when bc slot was already freed
		dsl.EmitMirEvent(m, eventpbevents.TimerDelay(mc.Timer, []*eventpbtypes.Event{
			bcpbevents.DoFillGap(mc.Self, slot),
		}, types.Duration(tunables.FillGapDelay)))

		return nil
	})

	// TODO: move this event to the right component
	bcdsl.UponDoFillGap(m, func(slot *commontypes.Slot) error {
		reqState, present := state.RequestsState[*slot]
		if !present {
			return nil // already handled
		}

		reqState.Replies = make(map[t.NodeID]struct{}, len(params.AllNodes))

		// logger.Log(logging.LevelDebug, "broadcast component fell behind. requesting slot from other replicas with FILL-GAP", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot)

		// TODO: do this more inteligently: only contact some nodes, and try others if a timer expires
		// until all were tried or a response is received.
		// It would also be nice to pass a hint in the certificate that says which nodes to try first,
		// this could be provided by the agreement component based on INIT(v, 0) messages received by the abba instances.
		rnetdsl.SendMessage(m, mc.ReliableNet,
			FillGapMsgID(slot),
			aleamsgs.FillGapMessage(mc.Self, slot),
			params.AllNodes,
		)
		return nil
	})

	// When receive a request for batch from another node, lookup the batch in the local storage.
	aleadsl.UponFillGapMessageReceived(m, func(from t.NodeID, slot *commontypes.Slot) error {
		// do not ACK message - acknowledging means sending a FILLER reply

		// logger.Log(logging.LevelDebug, "satisfying FILL-GAP request", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot, "from", from)
		batchdbdsl.LookupBatch(m, mc.BatchDB, common.FormatAleaBatchID(slot), &lookupBatchOnRemoteRequestContext{from, slot})
		return nil
	})

	// If the batch is found in the local storage, send it to the requesting node.
	batchdbdsl.UponLookupBatchResponse(m, func(found bool, txs []*requestpbtypes.Request, signature []byte, context *lookupBatchOnRemoteRequestContext) error {
		if !found {
			// Ignore invalid request.
			return nil
		}

		transportpbdsl.SendMessage(m, mc.Net,
			aleamsgs.FillerMessage(mc.Self, context.slot, txs, signature),
			[]t.NodeID{context.requester},
		)
		return nil
	})

	// When receive a requested batch, compute the ids of the received transactions.
	aleadsl.UponFillerMessageReceived(m, func(from t.NodeID, slot *commontypes.Slot, txs []*requestpbtypes.Request, signature tctypes.FullSig) error {
		rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, FillGapMsgID(slot), []t.NodeID{from})

		reqState, present := state.RequestsState[*slot]
		if !present || reqState.Replies == nil {
			return nil // no request needs this message to be satisfied
		}

		if _, alreadyAnswered := reqState.Replies[from]; alreadyAnswered {
			return nil // already processed a reply from this node
		}
		reqState.Replies[from] = struct{}{}

		// logger.Log(logging.LevelDebug, "got FILLER for missing slot!", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot, "from", from)

		mempooldsl.RequestTransactionIDs(m, mc.Mempool, txs, &handleFillerContext{
			slot:      slot,
			txs:       txs,
			signature: signature,
		})
		return nil
	})

	// Compute signature data
	mempooldsl.UponTransactionIDsResponse(m, func(txIDs []tt.TxID, context *handleFillerContext) error {
		context.txIDs = txIDs
		hasherpbdsl.RequestOne(m, mc.Hasher, &commonpbtypes.HashData{
			Data: txIDs,
		}, context)
		return nil
	})
	hasherpbdsl.UponResultOne(m, func(txIDsHash []byte, context *handleFillerContext) error {
		sigData := certSigData(params, context.slot, txIDsHash)
		threshDsl.VerifyFull(m, mc.ThreshCrypto, sigData, context.signature, context)

		return nil
	})

	// Check if signature is correct
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
		batchdbdsl.StoreBatch(m, mc.BatchDB, common.FormatAleaBatchID(context.slot), context.txIDs, context.txs, context.signature /*metadata*/, context)

		// send response to requests
		// logger.Log(logging.LevelDebug, "satisfying delayed requests with FILLER", "queueIdx", context.slot.QueueIdx, "queueSlot", context.slot.QueueSlot)
		for _, origin := range requestState.ReqOrigins {
			adsl.ProvideTransactions(m, origin.Module, context.txs, origin)
		}
		mempooldsl.MarkDelivered(m, mc.Mempool, context.txs)
		delete(state.RequestsState, *context.slot)

		return nil
	})

	batchdbdsl.UponBatchStored(m, func(context *handleFillerContext) error {
		abcdsl.FreeSlot(m, mc.AleaBroadcast, context.slot)
		return nil
	})
}

func certSigData(params *director.ModuleParams, slot *commontypes.Slot, txIDsHash []byte) [][]byte {
	aleaUID := params.InstanceUID[:len(params.InstanceUID)-1]
	aleaBcInstanceUID := append(aleaUID, 'b')
	return vcb.SigData(bcutil.VCBInstanceUID(aleaBcInstanceUID, slot.QueueIdx, slot.QueueSlot), txIDsHash)
}

const (
	MsgTypeFillGap = "f"
)

func FillGapMsgID(slot *commontypes.Slot) rntypes.MsgID {
	return rntypes.MsgID(fmt.Sprintf("%s.%d.%d", MsgTypeFillGap, slot.QueueIdx, slot.QueueSlot))
}

// Context data structures

type lookupBatchOnRemoteRequestContext struct {
	requester t.NodeID
	slot      *commontypes.Slot
}

type handleFillerContext struct {
	slot      *commontypes.Slot
	txs       []*requestpbtypes.Request
	signature []byte

	txIDs []tt.TxID
}
