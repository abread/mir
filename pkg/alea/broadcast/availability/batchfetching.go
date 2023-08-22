package availability

import (
	"fmt"
	"time"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	"github.com/filecoin-project/mir/pkg/alea/util"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	bcpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcpbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/events"
	bcpbmsgs "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/msgs"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	batchdbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/dsl"
	availabilitypbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/dsl"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	eventpbevents "github.com/filecoin-project/mir/pkg/pb/eventpb/events"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	hasherpbdsl "github.com/filecoin-project/mir/pkg/pb/hasherpb/dsl"
	hasherpbtypes "github.com/filecoin-project/mir/pkg/pb/hasherpb/types"
	mempooldsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	rnetdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	threshDsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	transportpbdsl "github.com/filecoin-project/mir/pkg/pb/transportpb/dsl"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	"github.com/filecoin-project/mir/pkg/timer/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
	"github.com/filecoin-project/mir/pkg/vcb"
)

// State represents the state related to this part of the module.
type batchFetchingState struct {
	RequestsState map[bcpbtypes.Slot]*RequestsState
}

// RequestsState represents the state related to a request on the source node of the request.
// The node disposes of this state as soon as the request is completed.
type RequestsState struct {
	ReqOrigins  []*availabilitypbtypes.RequestTransactionsOrigin
	SentFillGap bool
	Replies     map[t.NodeID]struct{}
}

// includeBatchFetching registers event handlers for processing availabilitypb.RequestTransactions events.
func includeBatchFetching(
	m dsl.Module,
	mc bccommon.ModuleConfig,
	params bccommon.ModuleParams,
	logger logging.Logger,
	est *bcEstimators,
) {
	_ = logger // silence warnings

	state := batchFetchingState{
		RequestsState: make(map[bcpbtypes.Slot]*RequestsState),
	}

	// When receive a request for transactions, first check the local storage.
	availabilitypbdsl.UponRequestTransactions(m, func(anyCert *availabilitypbtypes.Cert, origin *availabilitypbtypes.RequestTransactionsOrigin) error {
		certWrapper, present := anyCert.Type.(*availabilitypbtypes.Cert_Alea)
		if !present {
			return es.Errorf("unexpected certificate type. Expected: %T, got %T", certWrapper, anyCert.Type)
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
			batchdbdsl.LookupBatch(m, mc.BatchDB, util.FormatAleaBatchID(slot), slot)
		}

		// add ourselves to the list of requestors in order to receive a reply
		reqState.ReqOrigins = append(reqState.ReqOrigins, origin)

		return nil
	})

	// if broadcast delivers for a batch being requested, we can *now* resolve it locally
	bcqueuepbdsl.UponDeliver(m, func(slot *bcpbtypes.Slot) error {
		if _, present := state.RequestsState[*slot]; present {
			// TODO: avoid concurrent lookups of the same batch?
			// if bc delivers right after the transaction starts, two lookups will be performed.

			// retry locally, you will now succeed!
			batchdbdsl.LookupBatch(m, mc.BatchDB, util.FormatAleaBatchID(slot), slot)
		}

		return nil
	})

	// If the batch is present in the local storage, return it. Otherwise, ask other nodes.
	batchdbdsl.UponLookupBatchResponse(m, func(found bool, txs []*trantorpbtypes.Transaction, metadata []byte, slot *bcpbtypes.Slot) error {
		reqState, ok := state.RequestsState[*slot]
		if !ok {
			return nil // stale request
		}

		if found {
			for _, origin := range reqState.ReqOrigins {
				availabilitypbdsl.ProvideTransactions(m, origin.Module, txs, origin)
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
			return es.Errorf("tried to send FILL-GAP twice")
		}
		reqState.SentFillGap = true

		// send FILL-GAP after a timeout (if request was not satisfied)
		// this way bc has more chances of completing before even trying to send a fill-gap message
		delay := time.Duration(0)

		if bcRuntime, ok := est.BcRuntime(*slot); ok {
			delay = est.MaxExtBcDuration() - bcRuntime
			if delay < 0 {
				delay = 0
			}
		}

		// TODO: adjust delay according to bc estimate. don't delay when bc slot was already freed
		dsl.EmitEvent(m, eventpbevents.TimerDelay(mc.Timer, []*eventpbtypes.Event{
			bcpbevents.DoFillGap(mc.Self, slot),
		}, types.Duration(delay)))

		return nil
	})

	bcpbdsl.UponDoFillGap(m, func(slot *bcpbtypes.Slot) error {
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
			bcpbmsgs.FillGapMessage(mc.Self, slot),
			params.AllNodes,
		)
		return nil
	})

	// When receive a request for batch from another node, lookup the batch in the local storage.
	bcpbdsl.UponFillGapMessageReceived(m, func(from t.NodeID, slot *bcpbtypes.Slot) error {
		// do not ACK message - acknowledging means sending a FILLER reply

		// logger.Log(logging.LevelDebug, "satisfying FILL-GAP request", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot, "from", from)
		batchdbdsl.LookupBatch(m, mc.BatchDB, util.FormatAleaBatchID(slot), &lookupBatchOnRemoteRequestContext{from, slot})
		return nil
	})

	// If the batch is found in the local storage, send it to the requesting node.
	batchdbdsl.UponLookupBatchResponse(m, func(found bool, txs []*trantorpbtypes.Transaction, signature []byte, context *lookupBatchOnRemoteRequestContext) error {
		if !found {
			// Ignore invalid request.
			return nil
		}

		transportpbdsl.SendMessage(m, mc.Net,
			bcpbmsgs.FillerMessage(mc.Self, context.slot, txs, signature),
			[]t.NodeID{context.requester},
		)
		return nil
	})

	// When receive a requested batch, compute the ids of the received transactions.
	bcpbdsl.UponFillerMessageReceived(m, func(from t.NodeID, slot *bcpbtypes.Slot, txs []*trantorpbtypes.Transaction, signature tctypes.FullSig) error {
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
		hasherpbdsl.RequestOne(m, mc.Hasher, &hasherpbtypes.HashData{
			Data: sliceutil.Transform(txIDs, func(i int, txID tt.TxID) []byte {
				return []byte(txID)
			}),
		}, context)
		return nil
	})
	hasherpbdsl.UponResultOne(m, func(txIDsHash []byte, context *handleFillerContext) error {
		sigData := certSigData(params.InstanceUID, context.slot, txIDsHash)
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
		// TODO: proper epochs (retention index)
		batchdbdsl.StoreBatch(m, mc.BatchDB, util.FormatAleaBatchID(context.slot), context.txs, tt.RetentionIndex(0), context.signature /*metadata*/, context)

		// send response to requests
		// logger.Log(logging.LevelDebug, "satisfying delayed requests with FILLER", "queueIdx", context.slot.QueueIdx, "queueSlot", context.slot.QueueSlot)
		for _, origin := range requestState.ReqOrigins {
			availabilitypbdsl.ProvideTransactions(m, origin.Module, context.txs, origin)
		}
		mempooldsl.MarkDelivered(m, mc.Mempool, context.txs)
		delete(state.RequestsState, *context.slot)

		return nil
	})

	batchdbdsl.UponBatchStored(m, func(context *handleFillerContext) error {
		bcqueuepbdsl.FreeSlot(m, bccommon.BcQueueModuleID(mc.Self, context.slot.QueueIdx), context.slot.QueueSlot)
		return nil
	})
}

func certSigData(instanceUID []byte, slot *bcpbtypes.Slot, txIDsHash []byte) [][]byte {
	aleaUID := instanceUID[:len(instanceUID)-1]
	aleaBcInstanceUID := append(aleaUID, 'b')
	return vcb.SigData(bccommon.VCBInstanceUID(aleaBcInstanceUID, slot.QueueIdx, slot.QueueSlot), txIDsHash)
}

const (
	MsgTypeFillGap = "f"
)

func FillGapMsgID(slot *bcpbtypes.Slot) rntypes.MsgID {
	return rntypes.MsgID(fmt.Sprintf("%s.%d.%d", MsgTypeFillGap, slot.QueueIdx, slot.QueueSlot))
}

// Context data structures

type lookupBatchOnRemoteRequestContext struct {
	requester t.NodeID
	slot      *bcpbtypes.Slot
}

type handleFillerContext struct {
	slot      *bcpbtypes.Slot
	txs       []*trantorpbtypes.Transaction
	signature []byte

	txIDs []tt.TxID
}
