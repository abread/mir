package availability

import (
	"fmt"
	"time"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	"github.com/filecoin-project/mir/pkg/alea/queueselectionpolicy"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	bcpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcpbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/events"
	bcpbmsgs "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/msgs"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	directorpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/dsl"
	apppbdsl "github.com/filecoin-project/mir/pkg/pb/apppb/dsl"
	batchdbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/dsl"
	batchdbpbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/dsl"
	availabilitypbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/dsl"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	checkpointpbtypes "github.com/filecoin-project/mir/pkg/pb/checkpointpb/types"
	eventpbevents "github.com/filecoin-project/mir/pkg/pb/eventpb/events"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	mempooldsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	rnetdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	threshDsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	transportpbdsl "github.com/filecoin-project/mir/pkg/pb/transportpb/dsl"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	"github.com/filecoin-project/mir/pkg/timer/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb"
)

// State represents the state related to this part of the module.
type batchFetchingState struct {
	epochNr tt.EpochNr

	// track which slots were stored in which epoch for retention purposes
	slotSavedEpochs     map[bcpbtypes.Slot]tt.EpochNr   // slot to epoch
	epochDeliveredSlots map[tt.EpochNr][]bcpbtypes.Slot // epoch to slots

	RequestsState map[bcpbtypes.Slot]*RequestsState
}

// RequestsState represents the state related to a request on the source node of the request.
// The node disposes of this state as soon as the request is completed.
type RequestsState struct {
	ReqOrigins  []*availabilitypbtypes.RequestTransactionsOrigin
	SentFillGap bool
	Replies     map[t.NodeID]struct{}
	EpochNr     tt.EpochNr
}

// includeBatchFetching registers event handlers for processing availabilitypb.RequestTransactions events.
func includeBatchFetching(
	m dsl.Module,
	mc ModuleConfig,
	params ModuleParams,
	logger logging.Logger,
	certDB map[bcpbtypes.Slot]*bcpbtypes.Cert,
	est *bcEstimators,
) {
	_ = logger // silence warnings

	state := batchFetchingState{
		epochNr: 0,

		slotSavedEpochs:     make(map[bcpbtypes.Slot]tt.EpochNr, int(params.RetainEpochs*params.EpochLength)),
		epochDeliveredSlots: make(map[tt.EpochNr][]bcpbtypes.Slot, int(params.RetainEpochs)),

		RequestsState: make(map[bcpbtypes.Slot]*RequestsState),
	}

	scheduleFillGap := func(slot *bcpbtypes.Slot, reqState *RequestsState) error {
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
			if cert, ok := certDB[*slot]; ok {
				batchdbdsl.LookupBatch(m, mc.BatchDB, cert.BatchId, slot)
			} else {
				if err := scheduleFillGap(slot, reqState); err != nil {
					return err
				}
			}
		}

		// add ourselves to the list of requestors in order to receive a reply
		reqState.ReqOrigins = append(reqState.ReqOrigins, origin)

		return nil
	})

	// if broadcast delivers for a batch being requested, we can *now* resolve it locally
	bcqueuepbdsl.UponDeliver(m, func(cert *bcpbtypes.Cert) error {
		if _, present := state.RequestsState[*cert.Slot]; present {
			// TODO: avoid concurrent lookups of the same batch?
			// if bc delivers right after the transaction starts, two lookups will be performed.

			// retry locally, you will now succeed!
			batchdbdsl.LookupBatch(m, mc.BatchDB, cert.BatchId, cert.Slot)
		}

		return nil
	})

	// If the batch is present in the local storage, return it. Otherwise, ask other nodes.
	batchdbdsl.UponLookupBatchResponse(m, func(found bool, txs []*trantorpbtypes.Transaction, slot *bcpbtypes.Slot) error {
		reqState, ok := state.RequestsState[*slot]
		if !ok {
			return nil // stale request
		}

		if !found {
			return es.Errorf("inconsistent batch db state: missing txs for slot %v", slot)
		}

		for _, origin := range reqState.ReqOrigins {
			availabilitypbdsl.ProvideTransactions(m, origin.Module, txs, origin)
		}

		if reqState.SentFillGap {
			// no need for FILLER anymore
			rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, FillGapMsgID(slot), params.AllNodes)
		}

		delete(state.RequestsState, *slot)
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

		if cert, ok := certDB[*slot]; ok {
			// we have it! go get the batch
			batchdbdsl.LookupBatch(m, mc.BatchDB, cert.BatchId, &lookupBatchOnRemoteRequestContext{from, cert})
		} else {
			// node is likely to be behind, send help
			// TODO: add some check to avoid sending help to nodes that are not behind (see slot data in cert)
			directorpbdsl.HelpNode(m, mc.Consumer, from)
		}
		// TODO: send indication of no-reply

		return nil
	})

	// If the batch is found in the local storage, send it to the requesting node.
	batchdbdsl.UponLookupBatchResponse(m, func(found bool, txs []*trantorpbtypes.Transaction, context *lookupBatchOnRemoteRequestContext) error {
		if !found {
			return es.Errorf("inconsistency between dbs: cert was in certdb, but no batch present")
		}

		// logger.Log(logging.LevelDebug, "satisfying FILL-GAP request", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot, "from", from)
		transportpbdsl.SendMessage(m, mc.Net,
			bcpbmsgs.FillerMessage(mc.Self, context.cert, txs),
			[]t.NodeID{context.requester},
		)
		return nil
	})

	// After receiving a Filler message, we must validate the provided information (batchID, signature)
	bcpbdsl.UponFillerMessageReceived(m, func(from t.NodeID, cert *bcpbtypes.Cert, txs []*trantorpbtypes.Transaction) error {
		rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, FillGapMsgID(cert.Slot), []t.NodeID{from})

		// TODO: do this smartly and try one node at a time instead of broadcasting FILL-GAP

		reqState, present := state.RequestsState[*cert.Slot]
		if !present || reqState.Replies == nil {
			return nil // no request needs this message to be satisfied
		}

		if _, alreadyAnswered := reqState.Replies[from]; alreadyAnswered {
			return nil // already processed a reply from this node
		}
		reqState.Replies[from] = struct{}{}

		// logger.Log(logging.LevelDebug, "got FILLER for missing slot!", "queueIdx", slot.QueueIdx, "queueSlot", slot.QueueSlot, "from", from)

		// compute tx ids in order to compute batch ID
		mempooldsl.RequestTransactionIDs(m, mc.Mempool, txs, &handleFillerContext{
			cert: cert,
			txs:  txs,
		})
		return nil
	})
	mempooldsl.UponTransactionIDsResponse(m, func(txIDs []tt.TxID, context *handleFillerContext) error {
		// compute batch ID
		mempooldsl.RequestBatchID(m, mc.Mempool, txIDs, context)
		return nil
	})
	mempooldsl.UponBatchIDResponse(m, func(batchId string, context *handleFillerContext) error {
		// check batchID
		if context.cert.BatchId != batchId {
			// TODO: report node as byz
			// TODO: do this smartly and try another node here instead of broadcasting FILL-GAP
			return nil
		}

		if _, ok := state.RequestsState[*context.cert.Slot]; !ok {
			// The request has already been completed.
			// Don't bother with verifying the signature
			return nil
		}

		// check signature
		sigData := certSigData(&params, context.cert)
		threshDsl.VerifyFull(m, mc.ThreshCrypto, sigData, context.cert.Signature, context)
		return nil
	})
	threshDsl.UponVerifyFullResult(m, func(ok bool, err string, context *handleFillerContext) error {
		if !ok {
			// TODO: report node as byz
			// TODO: do this smartly and try another node here instead of broadcasting FILL-GAP
			return nil
		}

		requestState, ok := state.RequestsState[*context.cert.Slot]
		if !ok {
			// The request has already been completed.
			return nil
		}

		// stop asking other nodes to send us stuff
		rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, FillGapMsgID(context.cert.Slot), params.AllNodes)

		// store batch/cert asynchronously
		batchdbdsl.StoreBatch(m, mc.BatchDB, context.cert.BatchId, context.txs, tt.RetentionIndex(state.epochNr), context)

		// send response to requests
		// logger.Log(logging.LevelDebug, "satisfying delayed requests with FILLER", "queueIdx", context.slot.QueueIdx, "queueSlot", context.slot.QueueSlot)
		for _, origin := range requestState.ReqOrigins {
			availabilitypbdsl.ProvideTransactions(m, origin.Module, context.txs, origin)
		}
		mempooldsl.MarkStableProposal(m, mc.Mempool, context.txs)
		delete(state.RequestsState, *context.cert.Slot)

		return nil
	})
	batchdbdsl.UponBatchStored(m, func(context *handleFillerContext) error {
		// batch is stored, we can now store the corresponding cert
		certDB[*context.cert.Slot] = context.cert

		// this means the corresponding broadcast was completed, albeit through a non-convetional path
		// we can free the vcb instance for this slot
		bcqueuepbdsl.FreeSlot(m, bccommon.BcQueueModuleID(mc.Self, context.cert.Slot.QueueIdx), context.cert.Slot.QueueSlot)
		return nil
	})

	directorpbdsl.UponNewEpoch(m, func(epoch tt.EpochNr) error {
		state.epochNr = epoch
		state.epochDeliveredSlots[epoch] = make([]bcpbtypes.Slot, 0, params.EpochLength)

		// propagate epoch changes to queues
		for queueIdx := aleatypes.QueueIdx(0); queueIdx < aleatypes.QueueIdx(len(params.AllNodes)); queueIdx++ {
			directorpbdsl.NewEpoch(m, bccommon.BcQueueModuleID(mc.Self, queueIdx), epoch)
		}
		return nil
	})

	// ensure batch retention is updated on delivery
	bcqueuepbdsl.UponBcStarted(m, func(slot *bcpbtypes.Slot, epoch tt.EpochNr) error {
		state.slotSavedEpochs[*slot] = epoch
		return nil
	})
	bcpbdsl.UponFreeSlot(m, func(slot *bcpbtypes.Slot) error {
		slotSavedEpoch := state.slotSavedEpochs[*slot]

		if slotSavedEpoch != state.epochNr {
			// update retention index in batchdb to match delivery epoch
			// this is required to avoid dropping batches too soon
			batchdbpbdsl.UpdateBatchRetention(m, mc.BatchDB, certDB[*slot].BatchId, tt.RetentionIndex(state.epochNr))
		}
		// we no longer need to worry about this slot
		delete(state.slotSavedEpochs, *slot)

		// also note when this batch was delivered
		state.epochDeliveredSlots[state.epochNr] = append(state.epochDeliveredSlots[state.epochNr], *slot)

		return nil
	})

	directorpbdsl.UponEpochCheckpointed(m, func(epoch tt.EpochNr) error {
		if epoch > tt.EpochNr(params.RetainEpochs)+1 {
			// ensure previously broadcast (but not agreed) slots don't get garbage collected too soon
			for slot, slotSavedEpoch := range state.slotSavedEpochs {
				if slotSavedEpoch < epoch-tt.EpochNr(params.RetainEpochs) {
					batchdbpbdsl.UpdateBatchRetention(m, mc.BatchDB, certDB[slot].BatchId, tt.RetentionIndex(state.epochNr))
				}
			}

			// delete old batches from db
			batchdbpbdsl.GarbageCollect(m, mc.BatchDB, tt.RetentionIndex(uint64(epoch)-params.RetainEpochs))

			// delete old availability certs from db
			for _, slot := range state.epochDeliveredSlots[epoch-tt.EpochNr(params.RetainEpochs)] {
				delete(certDB, slot)
			}
			delete(state.epochDeliveredSlots, epoch-tt.EpochNr(params.RetainEpochs))
		}

		return nil
	})

	apppbdsl.UponRestoreState(m, func(checkpoint *checkpointpbtypes.StableCheckpoint) error {
		epochNr := checkpoint.Snapshot.EpochData.EpochConfig.EpochNr

		// for batches that have not been delivered normally, set their retention index to the checkpoint epoch
		for slot := range state.slotSavedEpochs {
			batchdbpbdsl.UpdateBatchRetention(m, mc.BatchDB, certDB[slot].BatchId, tt.RetentionIndex(epochNr))
		}
		state.slotSavedEpochs = make(map[bcpbtypes.Slot]tt.EpochNr, len(state.slotSavedEpochs))

		// delete old availability certs from db
		var minEpoch tt.EpochNr
		if uint64(epochNr) > params.RetainEpochs+1 {
			minEpoch = epochNr - tt.EpochNr(params.RetainEpochs)
		} else {
			minEpoch = tt.EpochNr(0)
		}
		for epochNr, deliveredSlots := range state.epochDeliveredSlots {
			if epochNr >= minEpoch {
				continue
			}

			for _, slot := range deliveredSlots {
				delete(certDB, slot)
			}
			delete(state.epochDeliveredSlots, epochNr)
		}

		qsp, err := queueselectionpolicy.QueuePolicyFromBytes(checkpoint.Snapshot.EpochData.LeaderPolicy)
		if err != nil {
			return err
		}

		// set delivered certificate retention index to the checkpoint epoch, to ensure they are not garbage collected too soon
		// Note: this may add duplicates to state.epochDeliveredSlots, but that's ok (as long as memory is not a hard concern)
		for slot := range certDB {
			if qsp.SlotDelivered(slot) {
				state.epochDeliveredSlots[epochNr] = append(state.epochDeliveredSlots[epochNr], slot)
			}
		}

		// cancel pending fill-gap requests for slots already delivered in the checkpoint
		for slot, reqSt := range state.RequestsState {
			if qsp.SlotDelivered(slot) {
				// HACK: we emit ProvideTransactions(nil) to make sure DSL context is freed on batchfetcher
				// this is a no-op on batch-fetcher since the new epoch has already started
				for _, origin := range reqSt.ReqOrigins {
					availabilitypbdsl.ProvideTransactions(m, origin.Module, nil, origin)
				}

				delete(state.RequestsState, slot)
			}
		}

		return nil
	})
}

func certSigData(params *ModuleParams, cert *bcpbtypes.Cert) [][]byte {
	aleaUID := params.InstanceUID[:len(params.InstanceUID)-1]
	aleaBcInstanceUID := append(aleaUID, 'b')
	return vcb.SigData(bccommon.VCBInstanceUID(aleaBcInstanceUID, cert.Slot.QueueIdx, cert.Slot.QueueSlot), cert.BatchId)
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
	cert      *bcpbtypes.Cert
}

type handleFillerContext struct {
	cert *bcpbtypes.Cert
	txs  []*trantorpbtypes.Transaction
}
