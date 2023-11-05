package aleatracer

import (
	"context"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	es "github.com/go-errors/errors"
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/clientprogress"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	abbapbtypes "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	ageventstypes "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/types"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	batchfetcherpbtypes "github.com/filecoin-project/mir/pkg/pb/batchfetcherpb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	mempoolpbtypes "github.com/filecoin-project/mir/pkg/pb/mempoolpb/types"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	threshcryptopbtypes "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/types"
	transportpbtypes "github.com/filecoin-project/mir/pkg/pb/transportpb/types"
	vcbpbtypes "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	"github.com/filecoin-project/mir/pkg/trantor"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type AleaTracer struct {
	ownQueueIdx aleatypes.QueueIdx
	params      *trantor.Params

	wipTxSSpan               map[txID]*span
	wipTxFSpan               map[txID]*span
	wipBcSpan                map[bcpbtypes.Slot]*span
	wipBcAwaitEchoSpan       map[bcpbtypes.Slot]*span
	wipBcAwaitFinalSpan      map[bcpbtypes.Slot]*span
	wipBcComputeSigDataSpan  map[bcpbtypes.Slot]*span
	wipBcAwaitQuorumDoneSpan map[bcpbtypes.Slot]*span
	wipBcAwaitAllDoneSpan    map[bcpbtypes.Slot]*span
	wipBcModSpan             map[bcpbtypes.Slot]*span
	wipAgSpan                map[uint64]*span
	wipAgModSpan             map[uint64]*span
	wipAbbaRoundSpan         map[abbaRoundID]*span
	wipAbbaRoundModSpan      map[abbaRoundID]*span
	wipBfSpan                map[bcpbtypes.Slot]*span
	wipBfStalledSpan         map[bcpbtypes.Slot]*span
	wipThreshCryptoSpan      map[string]*span

	// stall tracker
	nextBatchToCut       aleatypes.QueueSlot
	nextAgRound          uint64
	agRunning            bool
	wipBatchCutStallSpan *span
	wipAgStallSpan       *span
	unagreedSlots        map[aleatypes.QueueIdx]map[aleatypes.QueueSlot]struct{}

	bfWipSlotsCtxID map[uint64]bcpbtypes.Slot
	bfDeliverQueue  []bcpbtypes.Slot

	agQueueHeads            []aleatypes.QueueSlot
	agUndeliveredAbbaRounds map[uint64]map[uint64]struct{}

	txTracker *clientprogress.ClientProgress

	inChan chan events.EventList
	cancel context.CancelFunc
	wg     sync.WaitGroup
	out    io.Writer
}

type abbaRoundID struct {
	agRound   uint64
	abbaRound uint64
}

type span struct {
	class string
	id    string

	start time.Duration
	end   time.Duration
}

type txID struct {
	clientID tt.ClientID
	txNo     tt.TxNo
}

const N = 128

var timeRef = time.Now()

func NewAleaTracer(ctx context.Context, params *trantor.Params, nodeID t.NodeID, out io.Writer) *AleaTracer {
	tracerCtx, cancel := context.WithCancel(ctx)

	N := len(params.Alea.Membership.Nodes)

	tracer := &AleaTracer{
		ownQueueIdx: aleatypes.QueueIdx(slices.Index(params.Alea.AllNodes(), nodeID)),
		params:      params,

		wipTxSSpan:               make(map[txID]*span, 256),
		wipTxFSpan:               make(map[txID]*span, 256),
		wipBcSpan:                make(map[bcpbtypes.Slot]*span, N),
		wipBcAwaitEchoSpan:       make(map[bcpbtypes.Slot]*span, N),
		wipBcAwaitFinalSpan:      make(map[bcpbtypes.Slot]*span, N),
		wipBcComputeSigDataSpan:  make(map[bcpbtypes.Slot]*span, N),
		wipBcAwaitQuorumDoneSpan: make(map[bcpbtypes.Slot]*span, N),
		wipBcAwaitAllDoneSpan:    make(map[bcpbtypes.Slot]*span, N),
		wipBcModSpan:             make(map[bcpbtypes.Slot]*span, N),
		wipAgSpan:                make(map[uint64]*span, N),
		wipAgModSpan:             make(map[uint64]*span, N),
		wipAbbaRoundSpan:         make(map[abbaRoundID]*span, N*16),
		wipAbbaRoundModSpan:      make(map[abbaRoundID]*span, N*16),
		wipBfSpan:                make(map[bcpbtypes.Slot]*span, N),
		wipBfStalledSpan:         make(map[bcpbtypes.Slot]*span, N),
		wipThreshCryptoSpan:      make(map[string]*span, N*32),

		nextBatchToCut:       0,
		nextAgRound:          0,
		wipBatchCutStallSpan: nil,
		wipAgStallSpan:       nil,
		unagreedSlots:        make(map[aleatypes.QueueIdx]map[aleatypes.QueueSlot]struct{}, N),

		bfWipSlotsCtxID: make(map[uint64]bcpbtypes.Slot, N),
		bfDeliverQueue:  nil,

		agQueueHeads:            make([]aleatypes.QueueSlot, N),
		agUndeliveredAbbaRounds: make(map[uint64]map[uint64]struct{}, N),

		txTracker: clientprogress.NewClientProgress(),

		inChan: make(chan events.EventList, N*16),
		cancel: cancel,
		out:    out,
	}

	for i := aleatypes.QueueIdx(0); i < aleatypes.QueueIdx(N); i++ {
		tracer.unagreedSlots[i] = make(map[aleatypes.QueueSlot]struct{}, 32)
	}

	ref := timeRef.UnixNano()
	_, err := out.Write([]byte(fmt.Sprintf("class,id,start,end\nmark,,%d,%d\n", ref, ref)))
	if err != nil {
		panic(err)
	}

	tracer.startBatchCutStallSpan(time.Duration(0), aleatypes.QueueSlot(0))

	tracer.wg.Add(1)
	doneC := tracerCtx.Done()
	go func() {
		defer tracer.wg.Done()

		for {
			select {
			case <-doneC:
				return
			case evs := <-tracer.inChan:
				it := evs.Iterator()
				for ev := it.Next(); ev != nil; ev = it.Next() {
					if err := tracer.interceptOne(ev); err != nil {
						panic(err)
					}
				}

				tracer.txTracker.GarbageCollect()
			}
		}
	}()

	return tracer
}

func (at *AleaTracer) nodeCount() int {
	return len(at.params.Alea.Membership.Nodes)
}

func (at *AleaTracer) Stop() {
	at.cancel()
	at.wg.Wait()
}

func (at *AleaTracer) Intercept(evs events.EventList) error {
	at.inChan <- evs
	return nil
}

var EmptySlot = bcpbtypes.Slot{QueueIdx: math.MaxUint32, QueueSlot: math.MaxUint64}

func (at *AleaTracer) interceptOne(event *eventpbtypes.Event) error { // nolint: gocognit,gocyclo
	ts := time.Since(timeRef)

	// consider all non-messagereceived events as module initialization
	if ev, ok := event.Type.(*eventpbtypes.Event_Transport); ok {
		if _, ok := ev.Transport.Type.(*transportpbtypes.Event_MessageReceived); !ok {
			at.registerModStart(ts, event)
		}
	} else {
		at.registerModStart(ts, event)
	}

	switch ev := event.Type.(type) {
	case *eventpbtypes.Event_Mempool:
		switch e := ev.Mempool.Type.(type) {
		case *mempoolpbtypes.Event_NewTransactions:
			for _, tx := range e.NewTransactions.Transactions {
				at.startTxFSpan(ts, tx.ClientId, tx.TxNo)
			}
		case *mempoolpbtypes.Event_RequestBatch:
			at.endBatchCutStallSpan(ts, at.nextBatchToCut)
		case *mempoolpbtypes.Event_NewBatch:
			for _, tx := range e.NewBatch.Txs {
				at.startTxSSpan(ts, tx.ClientId, tx.TxNo)
			}
		case *mempoolpbtypes.Event_BatchIdResponse:
			if strings.HasPrefix(string(event.DestModule), "availability/") {
				slot := parseSlotFromModuleID(event.DestModule)
				at.endBcComputeSigDataSpan(ts, slot)
			}
		}
	case *eventpbtypes.Event_AleaBcqueue:
		switch e := ev.AleaBcqueue.Type.(type) {
		case *bcqueuepbtypes.Event_InputValue:
			queueIdx := parseQueueIdxFromModuleID(event.DestModule)
			slot := bcpbtypes.Slot{
				QueueIdx:  queueIdx,
				QueueSlot: e.InputValue.QueueSlot,
			}
			at.startBcSpan(ts, slot)
		case *bcqueuepbtypes.Event_FreeSlot:
			queueIdx := parseQueueIdxFromModuleID(event.DestModule)
			slot := bcpbtypes.Slot{
				QueueIdx:  queueIdx,
				QueueSlot: e.FreeSlot.QueueSlot,
			}
			at.endBcSpan(ts, slot)
		case *bcqueuepbtypes.Event_Deliver:
			slot := e.Deliver.Cert.Slot

			if slot.QueueIdx == at.ownQueueIdx && slot.QueueSlot == at.nextBatchToCut {
				at.nextBatchToCut++
				at.startBatchCutStallSpan(ts, at.nextBatchToCut)
			}

			if slot.QueueSlot >= at.agQueueHeads[slot.QueueIdx] {
				at.unagreedSlots[slot.QueueIdx][slot.QueueSlot] = struct{}{}
			}
		}
	case *eventpbtypes.Event_AleaAgreement:
		switch e := ev.AleaAgreement.Type.(type) {
		case *ageventstypes.Event_Deliver:
			at.endAgSpan(ts, e.Deliver.Round)
			at.endAgModSpan(ts, e.Deliver.Round)

			if e.Deliver.Decision {
				slot := at.slotForAgRound(e.Deliver.Round)
				at.bfDeliverQueue = append(at.bfDeliverQueue, slot)
				at.agQueueHeads[slot.QueueIdx]++
				delete(at.unagreedSlots[slot.QueueIdx], slot.QueueSlot)
			} else {
				at.bfDeliverQueue = append(at.bfDeliverQueue, EmptySlot)
			}

			at.nextAgRound = e.Deliver.Round + 1
			nextQueueIdx := aleatypes.QueueIdx(at.nextAgRound % uint64(at.nodeCount()))
			if _, ok := at.unagreedSlots[nextQueueIdx][at.agQueueHeads[nextQueueIdx]]; ok {
				// slot is present, ag will start immediately
				at.startAgStallSpan(ts, at.nextAgRound)
				at.endAgStallSpan(ts, at.nextAgRound)
				at.startAgSpan(ts, at.nextAgRound)
			} else {
				// slot not present, ag will potentially stall a bit
				at.agRunning = false
				at.startAgStallSpan(ts, at.nextAgRound)
			}
		}
	case *eventpbtypes.Event_Vcb:
		switch e := ev.Vcb.Type.(type) {
		case *vcbpbtypes.Event_InputValue:
			slot := parseSlotFromModuleID(event.DestModule)
			at.startBcSpan(ts, slot)
			at.startBcComputeSigDataSpan(ts, slot)
		case *vcbpbtypes.Event_QuorumDone:
			slot := parseSlotFromModuleID(e.QuorumDone.SrcModule)
			at.endBcAwaitQuorumDoneSpan(ts, slot)
			at.startBcAwaitAllDoneSpan(ts, slot)
		case *vcbpbtypes.Event_AllDone:
			slot := parseSlotFromModuleID(e.AllDone.SrcModule)
			at.endBcAwaitAllDoneSpan(ts, slot)
			at.endBcSpan(ts, slot)
		}
	case *eventpbtypes.Event_Transport:
		switch e := ev.Transport.Type.(type) {
		case *transportpbtypes.Event_SendMessage:
			switch msg := e.SendMessage.Msg.Type.(type) {
			case *messagepbtypes.Message_Vcb:
				slot := parseSlotFromModuleID(e.SendMessage.Msg.DestModule)
				switch msg.Vcb.Type.(type) {
				case *vcbpbtypes.Message_EchoMessage:
					at.startBcAwaitFinalSpan(ts, slot)
				case *vcbpbtypes.Message_SendMessage:
					at.startBcAwaitEchoSpan(ts, slot)
				case *vcbpbtypes.Message_FinalMessage:
					at.startBcAwaitQuorumDoneSpan(ts, slot)
				case *vcbpbtypes.Message_DoneMessage:
					if slot.QueueIdx != at.ownQueueIdx {
						at.endBcSpan(ts, slot)
						at.endBcModSpan(ts, slot)
					}
					// for own queue, wait for slot to be freed by agreement or all VCBs completing
				}
			}
		case *transportpbtypes.Event_MessageReceived:
			switch msg := e.MessageReceived.Msg.Type.(type) {
			case *messagepbtypes.Message_Vcb:
				slot := parseSlotFromModuleID(event.DestModule)
				switch msg.Vcb.Type.(type) {
				case *vcbpbtypes.Message_SendMessage:
					at.startBcSpan(ts, slot)
					at.startBcComputeSigDataSpan(ts, slot)
				case *vcbpbtypes.Message_FinalMessage:
					at.endBcAwaitFinalSpan(ts, slot)
				}
			}
		}
	case *eventpbtypes.Event_Abba:
		switch e := ev.Abba.Type.(type) {
		case *abbapbtypes.Event_Continue:
			agRoundID, err := at.parseAgRoundID(event.DestModule)
			if err != nil {
				return es.Errorf("invalid ag round id: %w", err)
			}
			at.endAgStallSpan(ts, agRoundID)
			at.startAgSpan(ts, agRoundID)

			at.startAbbaRoundSpan(ts, abbaRoundID{agRoundID, 0})
			at.agRunning = true
		case *abbapbtypes.Event_Round:
			switch e2 := e.Round.Type.(type) {
			case *abbapbtypes.RoundEvent_InputValue:
				abbaRoundID, err := at.parseAbbaRoundID(event.DestModule)
				if err != nil {
					return es.Errorf("invalid abba round id: %w", err)
				}

				if abbaRoundID.abbaRound != 0 || (abbaRoundID.agRound == at.nextAgRound && at.agRunning) {
					at.startAbbaRoundSpan(ts, abbaRoundID)
				}
			case *abbapbtypes.RoundEvent_Deliver:
				agRoundStr := event.DestModule.StripParent("ordering")
				agRound, err := strconv.ParseUint(string(agRoundStr), 10, 64)
				if err != nil {
					return es.Errorf("invalid ag round number: %w", err)
				}

				abbaRoundID := abbaRoundID{
					agRound:   agRound,
					abbaRound: e2.Deliver.RoundNumber,
				}

				at.endAbbaRoundSpan(ts, abbaRoundID)
				at.endAbbaRoundModSpan(ts, abbaRoundID)

				delete(at.agUndeliveredAbbaRounds[abbaRoundID.agRound], abbaRoundID.abbaRound)
			}
		}
	case *eventpbtypes.Event_Availability:
		switch e := ev.Availability.Type.(type) {
		case *availabilitypbtypes.Event_RequestTransactions:
			slot := *e.RequestTransactions.Cert.Type.(*availabilitypbtypes.Cert_Alea).Alea.Slot

			at.startBfSpan(ts, slot)

			ctxID := e.RequestTransactions.Origin.Type.(*availabilitypbtypes.RequestTransactionsOrigin_Dsl).Dsl.ContextID
			at.bfWipSlotsCtxID[ctxID] = slot
		case *availabilitypbtypes.Event_ProvideTransactions:
			ctxID := e.ProvideTransactions.Origin.Type.(*availabilitypbtypes.RequestTransactionsOrigin_Dsl).Dsl.ContextID
			slot := at.bfWipSlotsCtxID[ctxID]
			delete(at.bfWipSlotsCtxID, ctxID)

			at.endBfSpan(ts, slot)
			at.startDeliveryStalledSpan(ts, slot)

			// register bc end here to capture batches obtained from FILLER messages too
			at.endBcSpan(ts, slot)
		}
	case *eventpbtypes.Event_BatchFetcher:
		switch e := ev.BatchFetcher.Type.(type) {
		case *batchfetcherpbtypes.Event_NewOrderedBatch:
			slot := at.bfDeliverQueue[0]
			at.bfDeliverQueue = at.bfDeliverQueue[1:]
			if slot != EmptySlot {
				at.endDeliveryStalledSpan(ts, slot)

				for _, tx := range e.NewOrderedBatch.Txs {
					at.endTxSpan(ts, tx.ClientId, tx.TxNo)
				}
			}
		}
	case *eventpbtypes.Event_ThreshCrypto:
		switch e := ev.ThreshCrypto.Type.(type) {
		case *threshcryptopbtypes.Event_SignShare:
			originW, ok := e.SignShare.Origin.Type.(*threshcryptopbtypes.SignShareOrigin_Dsl)
			if !ok {
				return nil
			}

			at.startThreshCryptoSpan(ts, "tc:signShare", e.SignShare.Origin.Module, dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopbtypes.Event_SignShareResult:
			originW, ok := e.SignShareResult.Origin.Type.(*threshcryptopbtypes.SignShareOrigin_Dsl)
			if !ok {
				return nil
			}

			at.endThreshCryptoSpan(ts, "tc:signShare", e.SignShareResult.Origin.Module, dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopbtypes.Event_VerifyShare:
			originW, ok := e.VerifyShare.Origin.Type.(*threshcryptopbtypes.VerifyShareOrigin_Dsl)
			if !ok {
				return nil
			}

			at.startThreshCryptoSpan(ts, "tc:verifyShare", e.VerifyShare.Origin.Module, dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopbtypes.Event_VerifyShareResult:
			originW, ok := e.VerifyShareResult.Origin.Type.(*threshcryptopbtypes.VerifyShareOrigin_Dsl)
			if !ok {
				return nil
			}

			at.endThreshCryptoSpan(ts, "tc:verifyShare", e.VerifyShareResult.Origin.Module, dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopbtypes.Event_VerifyFull:
			originW, ok := e.VerifyFull.Origin.Type.(*threshcryptopbtypes.VerifyFullOrigin_Dsl)
			if !ok {
				return nil
			}

			at.startThreshCryptoSpan(ts, "tc:verifyFull", e.VerifyFull.Origin.Module, dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopbtypes.Event_VerifyFullResult:
			originW, ok := e.VerifyFullResult.Origin.Type.(*threshcryptopbtypes.VerifyFullOrigin_Dsl)
			if !ok {
				return nil
			}

			at.endThreshCryptoSpan(ts, "tc:verifyFull", e.VerifyFullResult.Origin.Module, dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopbtypes.Event_Recover:
			originW, ok := e.Recover.Origin.Type.(*threshcryptopbtypes.RecoverOrigin_Dsl)
			if !ok {
				return nil
			}
			if strings.HasPrefix(string(e.Recover.Origin.Module), "availability/") {
				slot := parseSlotFromModuleID(e.Recover.Origin.Module)
				if slot.QueueIdx == at.ownQueueIdx {
					at.endBcAwaitEchoSpan(ts, slot)
				}
			}

			at.startThreshCryptoSpan(ts, "tc:recover", e.Recover.Origin.Module, dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopbtypes.Event_RecoverResult:
			originW, ok := e.RecoverResult.Origin.Type.(*threshcryptopbtypes.RecoverOrigin_Dsl)
			if !ok {
				return nil
			}

			at.endThreshCryptoSpan(ts, "tc:recover", e.RecoverResult.Origin.Module, dsl.ContextID(originW.Dsl.ContextID))
		}
	}

	return nil
}

func (at *AleaTracer) registerModStart(ts time.Duration, event *eventpbtypes.Event) {
	id := event.DestModule
	if id.IsSubOf("ordering") {
		id = id.Sub()

		agRound, err := strconv.ParseUint(string(id.Top()), 10, 64)
		if err != nil {
			return
		}

		at.startAgModSpan(ts, agRound)

		if id.Sub().Top() == "r" {
			id = id.Sub().Sub()
			abbaRound, err := strconv.ParseUint(string(id), 10, 64)
			if err != nil {
				return
			}

			at.startAbbaRoundModSpan(ts, abbaRoundID{agRound, abbaRound})
		}
	} else if strings.HasPrefix(string(id), "availability/") {
		id = t.ModuleID(strings.TrimPrefix(string(id), "availability/"))

		queueIdx, err := strconv.ParseUint(string(id.Top()), 10, 64)
		if err != nil {
			return
		}
		queueSlot, err := strconv.ParseUint(string(id.Sub().Top()), 10, 64)
		if err != nil {
			return
		}

		// end prev mod span
		if queueSlot > uint64(at.params.Alea.MaxConcurrentVcbPerQueue) {
			at.endBcModSpan(ts, bcpbtypes.Slot{
				QueueIdx:  aleatypes.QueueIdx(queueIdx),
				QueueSlot: aleatypes.QueueSlot(queueSlot - uint64(at.params.Alea.MaxConcurrentVcbPerQueue)),
			})
		}
		at.startBcModSpan(ts, bcpbtypes.Slot{
			QueueIdx:  aleatypes.QueueIdx(queueIdx),
			QueueSlot: aleatypes.QueueSlot(queueSlot),
		})

		// coincides with the start, more or less
		at.startBcSpan(ts, bcpbtypes.Slot{
			QueueIdx:  aleatypes.QueueIdx(queueIdx),
			QueueSlot: aleatypes.QueueSlot(queueSlot),
		})
	}
}

func (at *AleaTracer) parseAgRoundID(id t.ModuleID) (uint64, error) {
	if !id.IsSubOf("ordering") {
		return 0, es.Errorf("not ordering/*")
	}
	id = id.Sub()

	agRound, err := strconv.ParseUint(string(id.Top()), 10, 64)
	if err != nil {
		return 0, es.Errorf("expected <ag round no.>, got %s", id.Top())
	}

	return agRound, nil
}

func (at *AleaTracer) parseAbbaRoundID(id t.ModuleID) (abbaRoundID, error) {
	if !id.IsSubOf("ordering") {
		return abbaRoundID{}, es.Errorf("not ordering/*")
	}
	id = id.Sub()

	agRound, err := strconv.ParseUint(string(id.Top()), 10, 64)
	if err != nil {
		return abbaRoundID{}, es.Errorf("expected <ag round no.>, got %s", id.Top())
	}
	id = id.Sub()

	if id.Top() != "r" {
		return abbaRoundID{}, es.Errorf("expected r, got %s", id.Top())
	}

	id = id.Sub()
	abbaRound, err := strconv.ParseUint(string(id.Top()), 10, 64)
	if err != nil {
		return abbaRoundID{}, es.Errorf("expected <abba round number>, got %s", id.Top())
	}

	return abbaRoundID{agRound, abbaRound}, nil
}

func (at *AleaTracer) slotForAgRound(round uint64) bcpbtypes.Slot {
	queueIdx := round % uint64(at.nodeCount())
	queueSlot := at.agQueueHeads[queueIdx]

	return bcpbtypes.Slot{
		QueueIdx:  aleatypes.QueueIdx(queueIdx),
		QueueSlot: queueSlot,
	}
}

func (at *AleaTracer) _bcSpan(ts time.Duration, slot bcpbtypes.Slot) *span {
	s, ok := at.wipBcSpan[slot]
	if !ok && at.agQueueHeads[slot.QueueIdx] <= slot.QueueSlot {
		at.wipBcSpan[slot] = &span{
			class: "bc",
			id:    fmt.Sprintf("%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBcSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBcSpan(ts time.Duration, slot bcpbtypes.Slot) {
	at._bcSpan(ts, slot)
}
func (at *AleaTracer) endBcSpan(ts time.Duration, slot bcpbtypes.Slot) {
	s, ok := at.wipBcSpan[slot]
	if !ok {
		return
	}
	s.end = ts
	at.writeSpan(s)

	delete(at.wipBcSpan, slot)

	if slot.QueueIdx == at.ownQueueIdx {
		at.endBcAwaitEchoSpan(ts, slot)
	}
	at.endBcAwaitFinalSpan(ts, slot)
	at.endBcComputeSigDataSpan(ts, slot)
	delete(at.wipBcAwaitEchoSpan, slot)
	delete(at.wipBcAwaitFinalSpan, slot)
	delete(at.wipBcComputeSigDataSpan, slot)
}

func (at *AleaTracer) startTxSSpan(ts time.Duration, clientID tt.ClientID, txNo tt.TxNo) {
	id := txID{clientID, txNo}
	at.wipTxSSpan[id] = &span{
		class: "tx:bestLat",
		id:    fmt.Sprintf("%s:%d", clientID, txNo),
		start: ts,
	}
}
func (at *AleaTracer) _txFSpan(ts time.Duration, clientID tt.ClientID, txNo tt.TxNo) *span {
	id := txID{clientID, txNo}

	s, ok := at.wipTxFSpan[id]
	if !ok {
		at.wipTxFSpan[id] = &span{
			class: "tx",
			id:    fmt.Sprintf("%s:%d", clientID, txNo),
			start: ts,
		}
		s = at.wipTxFSpan[id]
	}
	return s
}
func (at *AleaTracer) startTxFSpan(ts time.Duration, clientID tt.ClientID, txNo tt.TxNo) {
	if at.txTracker.Add(clientID, txNo) {
		at._txFSpan(ts, clientID, txNo)
	}
}
func (at *AleaTracer) endTxSpan(ts time.Duration, clientID tt.ClientID, txNo tt.TxNo) {
	id := txID{clientID, txNo}

	s, ok := at.wipTxFSpan[id]
	if ok {
		s.end = ts
		at.writeSpan(s)

		delete(at.wipTxFSpan, id)
	}

	s, ok = at.wipTxSSpan[id]
	if ok {
		s.end = ts
		at.writeSpan(s)

		delete(at.wipTxSSpan, id)
	}
}

func (at *AleaTracer) _batchCutStallSpan(ts time.Duration, queueSlot aleatypes.QueueSlot) *span {
	if at.wipBatchCutStallSpan == nil {
		at.wipBatchCutStallSpan = &span{
			class: "stall:batchCut",
			id:    fmt.Sprintf("%d/%d", at.ownQueueIdx, queueSlot),
			start: ts,
		}
	}
	return at.wipBatchCutStallSpan
}
func (at *AleaTracer) startBatchCutStallSpan(ts time.Duration, queueSlot aleatypes.QueueSlot) {
	at._batchCutStallSpan(ts, queueSlot)
}
func (at *AleaTracer) endBatchCutStallSpan(ts time.Duration, qs aleatypes.QueueSlot) bool {
	if at.wipBatchCutStallSpan == nil {
		return false
	}

	at.wipBatchCutStallSpan.end = ts
	if qs > 0 {
		at.writeSpan(at.wipBatchCutStallSpan)
	}
	at.wipBatchCutStallSpan = nil
	return true
}

func (at *AleaTracer) _agStallSpan(ts time.Duration, round uint64) *span {
	if at.wipAgStallSpan == nil {
		at.wipAgStallSpan = &span{
			class: "stall:ag",
			id:    strconv.FormatUint(round, 10),
			start: ts,
		}
	}
	return at.wipAgStallSpan
}
func (at *AleaTracer) startAgStallSpan(ts time.Duration, round uint64) {
	at._agStallSpan(ts, round)
}
func (at *AleaTracer) endAgStallSpan(ts time.Duration, _ uint64) {
	if at.wipAgStallSpan == nil {
		return
	}

	at.wipAgStallSpan.end = ts
	at.writeSpan(at.wipAgStallSpan)
	at.wipAgStallSpan = nil
}

func (at *AleaTracer) _bcAwaitEchoSpan(ts time.Duration, slot bcpbtypes.Slot) *span {
	s, ok := at.wipBcAwaitEchoSpan[slot]
	if !ok {
		at.wipBcAwaitEchoSpan[slot] = &span{
			class: "bc:waitEcho",
			id:    fmt.Sprintf("%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBcAwaitEchoSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBcAwaitEchoSpan(ts time.Duration, slot bcpbtypes.Slot) {
	at._bcAwaitEchoSpan(ts, slot)
}
func (at *AleaTracer) endBcAwaitEchoSpan(ts time.Duration, slot bcpbtypes.Slot) {
	s, ok := at.wipBcAwaitEchoSpan[slot]
	if ok && s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bcAwaitFinalSpan(ts time.Duration, slot bcpbtypes.Slot) *span {
	s, ok := at.wipBcAwaitFinalSpan[slot]
	if !ok {
		at.wipBcAwaitFinalSpan[slot] = &span{
			class: "bc:waitFinal",
			id:    fmt.Sprintf("%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBcAwaitFinalSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBcAwaitFinalSpan(ts time.Duration, slot bcpbtypes.Slot) {
	at._bcAwaitFinalSpan(ts, slot)
}
func (at *AleaTracer) endBcAwaitFinalSpan(ts time.Duration, slot bcpbtypes.Slot) {
	s, ok := at.wipBcAwaitFinalSpan[slot]
	if ok && s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bcComputeSigDataSpan(ts time.Duration, slot bcpbtypes.Slot) *span {
	s, ok := at.wipBcComputeSigDataSpan[slot]
	if !ok {
		at.wipBcComputeSigDataSpan[slot] = &span{
			class: "bc:cSigData",
			id:    fmt.Sprintf("%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBcComputeSigDataSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBcComputeSigDataSpan(ts time.Duration, slot bcpbtypes.Slot) {
	at._bcComputeSigDataSpan(ts, slot)
}
func (at *AleaTracer) endBcComputeSigDataSpan(ts time.Duration, slot bcpbtypes.Slot) {
	s, ok := at.wipBcComputeSigDataSpan[slot]
	if ok && s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bcAwaitQuorumDoneSpan(ts time.Duration, slot bcpbtypes.Slot) *span {
	s, ok := at.wipBcAwaitQuorumDoneSpan[slot]
	if !ok {
		_, modOk := at.wipBcModSpan[slot]
		if !modOk {
			return nil
		}

		at.wipBcAwaitQuorumDoneSpan[slot] = &span{
			class: "bc:awaitQuorum",
			id:    fmt.Sprintf("%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBcAwaitQuorumDoneSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBcAwaitQuorumDoneSpan(ts time.Duration, slot bcpbtypes.Slot) {
	at._bcAwaitQuorumDoneSpan(ts, slot)
}
func (at *AleaTracer) endBcAwaitQuorumDoneSpan(ts time.Duration, slot bcpbtypes.Slot) {
	s := at._bcAwaitQuorumDoneSpan(ts, slot)
	if s != nil && s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bcAwaitAllDoneSpan(ts time.Duration, slot bcpbtypes.Slot) *span {
	s, ok := at.wipBcAwaitAllDoneSpan[slot]
	if !ok {
		_, modOk := at.wipBcModSpan[slot]
		if !modOk {
			return nil
		}

		at.wipBcAwaitAllDoneSpan[slot] = &span{
			class: "bc:awaitAll",
			id:    fmt.Sprintf("%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBcAwaitAllDoneSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBcAwaitAllDoneSpan(ts time.Duration, slot bcpbtypes.Slot) {
	at._bcAwaitAllDoneSpan(ts, slot)
}
func (at *AleaTracer) endBcAwaitAllDoneSpan(ts time.Duration, slot bcpbtypes.Slot) {
	s := at._bcAwaitAllDoneSpan(ts, slot)
	if s != nil && s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bcModSpan(ts time.Duration, slot bcpbtypes.Slot) *span {
	s, ok := at.wipBcModSpan[slot]
	if !ok {
		at.wipBcModSpan[slot] = &span{
			class: "mod",
			id:    fmt.Sprintf("alea_bc/%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBcModSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBcModSpan(ts time.Duration, slot bcpbtypes.Slot) {
	at._bcModSpan(ts, slot)
}
func (at *AleaTracer) endBcModSpan(ts time.Duration, slot bcpbtypes.Slot) {
	s, ok := at.wipBcModSpan[slot]
	if !ok {
		return
	}
	s.end = ts
	at.writeSpan(s)

	delete(at.wipBcModSpan, slot)
	delete(at.wipBcAwaitQuorumDoneSpan, slot)
	delete(at.wipBcAwaitAllDoneSpan, slot)
}

func (at *AleaTracer) _agSpan(ts time.Duration, round uint64) *span {
	s, ok := at.wipAgSpan[round]
	if !ok {
		at.wipAgSpan[round] = &span{
			class: "ag",
			id:    fmt.Sprintf("%d", round),
			start: ts,
		}
		s = at.wipAgSpan[round]
	}
	return s
}
func (at *AleaTracer) startAgSpan(ts time.Duration, round uint64) {
	at._agSpan(ts, round)
}
func (at *AleaTracer) endAgSpan(ts time.Duration, round uint64) {
	s, ok := at.wipAgSpan[round]
	if !ok {
		return
	}
	s.end = ts
	at.writeSpan(s)

	delete(at.wipAgSpan, round)
}

func (at *AleaTracer) _agModSpan(ts time.Duration, round uint64) *span {
	s, ok := at.wipAgModSpan[round]
	if !ok {
		at.wipAgModSpan[round] = &span{
			class: "mod",
			id:    fmt.Sprintf("alea_ag/%d", round),
			start: ts,
		}
		at.agUndeliveredAbbaRounds[round] = make(map[uint64]struct{}, 8)
		s = at.wipAgModSpan[round]
	}
	return s
}
func (at *AleaTracer) startAgModSpan(ts time.Duration, round uint64) {
	at._agModSpan(ts, round)
}
func (at *AleaTracer) endAgModSpan(ts time.Duration, round uint64) {
	s, ok := at.wipAgModSpan[round]
	if !ok {
		return
	}
	s.end = ts

	if undeliveredRounds, ok := at.agUndeliveredAbbaRounds[round]; ok {
		for abbaRound := range undeliveredRounds {
			at.endAbbaRoundSpan(ts, abbaRoundID{round, abbaRound})
			at.endAbbaRoundModSpan(ts, abbaRoundID{round, abbaRound})
		}

		delete(at.agUndeliveredAbbaRounds, round)
	}

	at.writeSpan(s)
	delete(at.wipAgModSpan, round)
}

func (at *AleaTracer) _abbaRoundSpan(ts time.Duration, id abbaRoundID) *span {
	s, ok := at.wipAbbaRoundSpan[id]
	if !ok {
		at.wipAbbaRoundSpan[id] = &span{
			class: "abbaRound",
			id:    fmt.Sprintf("%d-%d", id.agRound, id.abbaRound),
			start: ts,
		}
		s = at.wipAbbaRoundSpan[id]
	}
	return s
}
func (at *AleaTracer) startAbbaRoundSpan(ts time.Duration, id abbaRoundID) {
	at._abbaRoundSpan(ts, id)
}
func (at *AleaTracer) endAbbaRoundSpan(ts time.Duration, id abbaRoundID) {
	s, ok := at.wipAbbaRoundSpan[id]
	if !ok {
		return
	}
	s.end = ts
	at.writeSpan(s)

	delete(at.wipAbbaRoundSpan, id)
}

func (at *AleaTracer) _abbaRoundModSpan(ts time.Duration, id abbaRoundID) *span {
	s, ok := at.wipAbbaRoundModSpan[id]
	if !ok {
		at.wipAbbaRoundModSpan[id] = &span{
			class: "mod",
			id:    fmt.Sprintf("alea_ag/%d/%d", id.agRound, id.abbaRound),
			start: ts,
		}
		at.agUndeliveredAbbaRounds[id.agRound][id.abbaRound] = struct{}{}
		s = at.wipAbbaRoundModSpan[id]
	}
	return s
}
func (at *AleaTracer) startAbbaRoundModSpan(ts time.Duration, id abbaRoundID) {
	at._abbaRoundModSpan(ts, id)
}
func (at *AleaTracer) endAbbaRoundModSpan(ts time.Duration, id abbaRoundID) {
	s, ok := at.wipAbbaRoundModSpan[id]
	if !ok {
		return
	}
	s.end = ts
	delete(at.agUndeliveredAbbaRounds[id.agRound], id.abbaRound)
	at.writeSpan(s)

	delete(at.wipAbbaRoundModSpan, id)
}

func (at *AleaTracer) _bfSpan(ts time.Duration, slot bcpbtypes.Slot) *span {
	s, ok := at.wipBfSpan[slot]
	if !ok {
		at.wipBfSpan[slot] = &span{
			class: "bf",
			id:    fmt.Sprintf("%d:%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBfSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBfSpan(ts time.Duration, slot bcpbtypes.Slot) {
	at._bfSpan(ts, slot)
}
func (at *AleaTracer) endBfSpan(ts time.Duration, slot bcpbtypes.Slot) {
	s, ok := at.wipBfSpan[slot]
	if !ok {
		return
	}

	s.end = ts
	at.writeSpan(s)

	delete(at.wipBfSpan, slot)
}

func (at *AleaTracer) _bfStalledSpan(ts time.Duration, slot bcpbtypes.Slot) *span {
	s, ok := at.wipBfStalledSpan[slot]
	if !ok {
		at.wipBfStalledSpan[slot] = &span{
			class: "stall:bf",
			id:    fmt.Sprintf("%d:%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBfStalledSpan[slot]
	}
	return s
}
func (at *AleaTracer) startDeliveryStalledSpan(ts time.Duration, slot bcpbtypes.Slot) {
	at._bfStalledSpan(ts, slot)
}
func (at *AleaTracer) endDeliveryStalledSpan(ts time.Duration, slot bcpbtypes.Slot) {
	s, ok := at.wipBfStalledSpan[slot]
	if !ok {
		return
	}

	s.end = ts
	at.writeSpan(s)

	delete(at.wipBfStalledSpan, slot)
}

func (at *AleaTracer) _threshCryptoSpan(ts time.Duration, class string, destModule t.ModuleID, ctxID dsl.ContextID) *span {
	id := fmt.Sprintf("%s(%v)", destModule, ctxID)
	uid := fmt.Sprintf("%s %s", class, id)
	s, ok := at.wipThreshCryptoSpan[uid]
	if !ok {
		at.wipThreshCryptoSpan[uid] = &span{
			class: class,
			id:    id,
			start: ts,
		}
		s = at.wipThreshCryptoSpan[uid]
	}
	return s
}
func (at *AleaTracer) startThreshCryptoSpan(ts time.Duration, class string, destModule t.ModuleID, ctxID dsl.ContextID) {
	at._threshCryptoSpan(ts, class, destModule, ctxID)
}
func (at *AleaTracer) endThreshCryptoSpan(ts time.Duration, class string, destModule t.ModuleID, ctxID dsl.ContextID) {
	id := fmt.Sprintf("%s(%v)", destModule, ctxID)
	uid := fmt.Sprintf("%s %s", class, id)
	s, ok := at.wipThreshCryptoSpan[uid]
	if !ok {
		return
	}

	s.end = ts
	at.writeSpan(s)
	delete(at.wipThreshCryptoSpan, uid)
}

func (at *AleaTracer) writeSpan(s *span) {
	_, err := at.out.Write([]byte(fmt.Sprintf("%s,%s,%d,%d\n", s.class, s.id, s.start, s.end)))
	if err != nil {
		panic(es.Errorf("failed to write trace span: %w", err))
	}
}

func parseSlotFromModuleID(moduleID t.ModuleID) bcpbtypes.Slot {
	if !strings.HasPrefix(string(moduleID), "availability/") {
		panic(es.Errorf("id is not from a bcqueue: %s", moduleID))
	}

	modID := t.ModuleID(strings.TrimPrefix(string(moduleID), "availability/"))
	queueIdxStr := string(modID.Top())
	queueSlotStr := string(modID.Sub().Top())

	queueIdx, err := strconv.ParseUint(queueIdxStr, 10, 32)
	if err != nil {
		panic(es.Errorf("failed to parse queueIdx from %s: %w", queueIdxStr, err))
	}

	queueSlot, err := strconv.ParseUint(queueSlotStr, 10, 32)
	if err != nil {
		panic(es.Errorf("failed to parse queueSlot from %s: %w", queueSlotStr, err))
	}

	return bcpbtypes.Slot{
		QueueIdx:  aleatypes.QueueIdx(queueIdx),
		QueueSlot: aleatypes.QueueSlot(queueSlot),
	}
}

func parseQueueIdxFromModuleID(moduleID t.ModuleID) aleatypes.QueueIdx {
	if !strings.HasPrefix(string(moduleID), "availability/") {
		panic(es.Errorf("id is not from a bcqueue: %s", moduleID))
	}

	modID := t.ModuleID(strings.TrimPrefix(string(moduleID), "availability/"))
	queueIdxStr := string(modID.Top())

	queueIdx, err := strconv.ParseUint(queueIdxStr, 10, 32)
	if err != nil {
		panic(es.Errorf("failed to parse queueIdx from %s: %w", queueIdxStr, err))
	}
	return aleatypes.QueueIdx(queueIdx)
}
