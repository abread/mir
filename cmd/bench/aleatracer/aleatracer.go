package aleatracer

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/abbapb"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb"
	"github.com/filecoin-project/mir/pkg/pb/batchfetcherpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/threshcryptopb"
	"github.com/filecoin-project/mir/pkg/pb/transportpb"
	"github.com/filecoin-project/mir/pkg/pb/vcbpb"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/localclock"
)

type AleaTracer struct {
	ownQueueIdx aleatypes.QueueIdx
	nodeCount   int

	wipBcSpan           map[commontypes.Slot]*span
	wipBcAwaitEchoSpan  map[commontypes.Slot]*span
	wipBcAwaitFinalSpan map[commontypes.Slot]*span
	wipBcModSpan        map[commontypes.Slot]*span
	wipAgSpan           map[uint64]*span
	wipAgModSpan        map[uint64]*span
	wipAbbaRoundSpan    map[abbaRoundID]*span
	wipAbbaRoundModSpan map[abbaRoundID]*span
	wipBfSpan           map[commontypes.Slot]*span
	wipBfStalledSpan    map[commontypes.Slot]*span
	wipThreshCryptoSpan map[string]*span

	bfWipSlotsCtxID map[uint64]commontypes.Slot
	bfDeliverQueue  []commontypes.Slot

	agQueueHeads            []aleatypes.QueueSlot
	agWipForQueue           map[aleatypes.QueueIdx]uint64
	agUndeliveredAbbaRounds map[uint64]map[uint64]struct{}

	inChan chan *events.EventList
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

const N = 128

func NewAleaTracer(ctx context.Context, ownQueueIdx aleatypes.QueueIdx, nodeCount int, out io.Writer) *AleaTracer {
	tracerCtx, cancel := context.WithCancel(ctx)

	tracer := &AleaTracer{
		ownQueueIdx: ownQueueIdx,
		nodeCount:   nodeCount,

		wipBcSpan:           make(map[commontypes.Slot]*span, nodeCount*N),
		wipBcAwaitEchoSpan:  make(map[commontypes.Slot]*span, nodeCount*N),
		wipBcAwaitFinalSpan: make(map[commontypes.Slot]*span, nodeCount*N),
		wipBcModSpan:        make(map[commontypes.Slot]*span, nodeCount*N),
		wipAgSpan:           make(map[uint64]*span, nodeCount*N),
		wipAgModSpan:        make(map[uint64]*span, nodeCount*N),
		wipAbbaRoundSpan:    make(map[abbaRoundID]*span, nodeCount*N),
		wipAbbaRoundModSpan: make(map[abbaRoundID]*span, nodeCount*N),
		wipBfSpan:           make(map[commontypes.Slot]*span, nodeCount*N),
		wipBfStalledSpan:    make(map[commontypes.Slot]*span, nodeCount*N),
		wipThreshCryptoSpan: make(map[string]*span, nodeCount*N),

		bfWipSlotsCtxID: make(map[uint64]commontypes.Slot, nodeCount),
		bfDeliverQueue:  nil,

		agQueueHeads:            make([]aleatypes.QueueSlot, nodeCount),
		agWipForQueue:           make(map[aleatypes.QueueIdx]uint64, nodeCount),
		agUndeliveredAbbaRounds: make(map[uint64]map[uint64]struct{}, nodeCount*N),

		inChan: make(chan *events.EventList, nodeCount*16),
		cancel: cancel,
		out:    out,
	}

	ref := localclock.RefTime().UnixNano()
	_, err := out.Write([]byte(fmt.Sprintf("class,id,start,end\nmark,,%d,%d\n", ref, ref)))
	if err != nil {
		panic(err)
	}

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
			}
		}
	}()

	return tracer
}

func (at *AleaTracer) Stop() {
	at.cancel()
	at.wg.Wait()
}

func (at *AleaTracer) Intercept(evs *events.EventList) error {
	at.inChan <- evs
	return nil
}

func (at *AleaTracer) interceptOne(event *eventpb.Event) error {
	ts := time.Duration(event.LocalTs)

	// consider all non-messagereceived events as module initialization
	if ev, ok := event.Type.(*eventpb.Event_Transport); ok {
		if _, ok := ev.Transport.Type.(*transportpb.Event_MessageReceived); !ok {
			at.registerModStart(ts, event)
		}
	} else {
		at.registerModStart(ts, event)
	}

	switch ev := event.Type.(type) {
	case *eventpb.Event_AleaBroadcast:
		switch e := ev.AleaBroadcast.Type.(type) {
		case *bcpb.Event_StartBroadcast:
			slot := commontypes.Slot{QueueIdx: at.ownQueueIdx, QueueSlot: aleatypes.QueueSlot(e.StartBroadcast.QueueSlot)}
			at.startBcSpan(ts, slot)
		case *bcpb.Event_FreeSlot:
			at.endBcModSpan(ts, *commontypes.SlotFromPb(e.FreeSlot.Slot))
		}
	case *eventpb.Event_BatchDb:
		switch e := ev.BatchDb.Type.(type) {
		case *batchdbpb.Event_Store:
			// register bc end here to capture batches obtained from FILLER messages too
			id := string(e.Store.BatchId)
			id = strings.TrimPrefix(id, "alea-")
			idParts := strings.Split(id, ":")

			queueIdx, err := strconv.ParseUint(idParts[0], 10, 64)
			if err != nil {
				return fmt.Errorf("bad queue index in batch id (%s): %w", id, err)
			}
			queueSlot, err := strconv.ParseUint(idParts[1], 10, 64)
			if err != nil {
				return fmt.Errorf("bad queue slot in batch id (%s): %w", id, err)
			}

			slot := commontypes.Slot{QueueIdx: aleatypes.QueueIdx(queueIdx), QueueSlot: aleatypes.QueueSlot(queueSlot)}
			at.endBcSpan(ts, slot)
		}
	case *eventpb.Event_AleaAgreement:
		switch e := ev.AleaAgreement.Type.(type) {
		case *agevents.Event_InputValue:
			slot := at.slotForAgRound(e.InputValue.Round)
			at.agWipForQueue[slot.QueueIdx] = e.InputValue.Round

			at.startAgSpan(ts, e.InputValue.Round)
		case *agevents.Event_Deliver:
			at.endAgSpan(ts, e.Deliver.Round)
			at.endAgModSpan(ts, e.Deliver.Round)

			slot := at.slotForAgRound(e.Deliver.Round)
			delete(at.agWipForQueue, slot.QueueIdx)
			if e.Deliver.Decision {
				at.agQueueHeads[slot.QueueIdx]++
				at.bfDeliverQueue = append(at.bfDeliverQueue, slot)
			}
		}
	case *eventpb.Event_Vcb:
		switch ev.Vcb.Type.(type) {
		case *vcbpb.Event_InputValue:
			slot := parseSlotFromModuleID(event.DestModule)
			at.startBcSpan(ts, slot)
		}
	case *eventpb.Event_Transport:
		switch e := ev.Transport.Type.(type) {
		case *transportpb.Event_SendMessage:
			switch msg := e.SendMessage.Msg.Type.(type) {
			case *messagepb.Message_Vcb:
				switch msg.Vcb.Type.(type) {
				case *vcbpb.Message_EchoMessage:
					slot := parseSlotFromModuleID(e.SendMessage.Msg.DestModule)
					at.startBcAwaitFinalSpan(ts, slot)
				case *vcbpb.Message_SendMessage:
					slot := parseSlotFromModuleID(e.SendMessage.Msg.DestModule)
					at.startBcAwaitEchoSpan(ts, slot)
				}
			}
		case *transportpb.Event_MessageReceived:
			switch msg := e.MessageReceived.Msg.Type.(type) {
			case *messagepb.Message_Vcb:
				switch msg.Vcb.Type.(type) {
				case *vcbpb.Message_FinalMessage:
					slot := parseSlotFromModuleID(event.DestModule)
					at.endBcAwaitFinalSpan(ts, slot)
				}
			}
		}
	case *eventpb.Event_Abba:
		switch e := ev.Abba.Type.(type) {
		case *abbapb.Event_Round:
			switch e2 := e.Round.Type.(type) {
			case *abbapb.RoundEvent_InputValue:
				abbaRoundID, err := at.parseAbbaRoundID(t.ModuleID(event.DestModule))
				if err != nil {
					return fmt.Errorf("invalid abba round id: %w", err)
				}
				at.startAbbaRoundSpan(ts, abbaRoundID)
			case *abbapb.RoundEvent_Deliver:
				agRoundStr := t.ModuleID(event.DestModule).StripParent("alea_ag")
				agRound, err := strconv.ParseUint(string(agRoundStr), 10, 64)
				if err != nil {
					return fmt.Errorf("invalid ag round number: %w", err)
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
	case *eventpb.Event_Availability:
		switch e := ev.Availability.Type.(type) {
		case *availabilitypb.Event_RequestTransactions:
			slot := *commontypes.SlotFromPb(e.RequestTransactions.Cert.Type.(*availabilitypb.Cert_Alea).Alea.Slot)

			at.startBfSpan(ts, slot)

			ctxID := e.RequestTransactions.Origin.GetDsl().ContextID
			at.bfWipSlotsCtxID[ctxID] = slot
		case *availabilitypb.Event_ProvideTransactions:
			ctxID := e.ProvideTransactions.Origin.GetDsl().ContextID
			slot := at.bfWipSlotsCtxID[ctxID]
			delete(at.bfWipSlotsCtxID, ctxID)

			at.endBfSpan(ts, slot)
			at.startDeliveryStalledSpan(ts, slot)
		}
	case *eventpb.Event_BatchFetcher:
		switch ev.BatchFetcher.Type.(type) {
		case *batchfetcherpb.Event_NewOrderedBatch:
			slot := at.bfDeliverQueue[0]
			at.bfDeliverQueue = at.bfDeliverQueue[1:]

			at.endDeliveryStalledSpan(ts, slot)
		}
	case *eventpb.Event_ThreshCrypto:
		switch e := ev.ThreshCrypto.Type.(type) {
		case *threshcryptopb.Event_SignShare:
			originW, ok := e.SignShare.Origin.Type.(*threshcryptopb.SignShareOrigin_Dsl)
			if !ok {
				return nil
			}

			at.startThreshCryptoSpan(ts, "tc:signShare", t.ModuleID(e.SignShare.Origin.Module), dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopb.Event_SignShareResult:
			originW, ok := e.SignShareResult.Origin.Type.(*threshcryptopb.SignShareOrigin_Dsl)
			if !ok {
				return nil
			}

			at.endThreshCryptoSpan(ts, "tc:signShare", t.ModuleID(e.SignShareResult.Origin.Module), dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopb.Event_VerifyShare:
			originW, ok := e.VerifyShare.Origin.Type.(*threshcryptopb.VerifyShareOrigin_Dsl)
			if !ok {
				return nil
			}

			at.startThreshCryptoSpan(ts, "tc:verifyShare", t.ModuleID(e.VerifyShare.Origin.Module), dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopb.Event_VerifyShareResult:
			originW, ok := e.VerifyShareResult.Origin.Type.(*threshcryptopb.VerifyShareOrigin_Dsl)
			if !ok {
				return nil
			}

			at.endThreshCryptoSpan(ts, "tc:verifyShare", t.ModuleID(e.VerifyShareResult.Origin.Module), dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopb.Event_VerifyFull:
			originW, ok := e.VerifyFull.Origin.Type.(*threshcryptopb.VerifyFullOrigin_Dsl)
			if !ok {
				return nil
			}

			at.startThreshCryptoSpan(ts, "tc:verifyFull", t.ModuleID(e.VerifyFull.Origin.Module), dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopb.Event_VerifyFullResult:
			originW, ok := e.VerifyFullResult.Origin.Type.(*threshcryptopb.VerifyFullOrigin_Dsl)
			if !ok {
				return nil
			}

			at.endThreshCryptoSpan(ts, "tc:verifyFull", t.ModuleID(e.VerifyFullResult.Origin.Module), dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopb.Event_Recover:
			originW, ok := e.Recover.Origin.Type.(*threshcryptopb.RecoverOrigin_Dsl)
			if !ok {
				return nil
			}

			at.startThreshCryptoSpan(ts, "tc:recover", t.ModuleID(e.Recover.Origin.Module), dsl.ContextID(originW.Dsl.ContextID))
		case *threshcryptopb.Event_RecoverResult:
			originW, ok := e.RecoverResult.Origin.Type.(*threshcryptopb.RecoverOrigin_Dsl)
			if !ok {
				return nil
			}

			at.endThreshCryptoSpan(ts, "tc:recover", t.ModuleID(e.RecoverResult.Origin.Module), dsl.ContextID(originW.Dsl.ContextID))

			if strings.HasPrefix(e.RecoverResult.Origin.Module, "alea_bc") {
				slot := parseSlotFromModuleID(e.RecoverResult.Origin.Module)
				at.endBcAwaitEchoSpan(ts, slot)
			}
		}
	}

	return nil
}

func (at *AleaTracer) registerModStart(ts time.Duration, event *eventpb.Event) {
	id := t.ModuleID(event.DestModule)
	if id.IsSubOf("alea_ag") {
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
	} else if id.IsSubOf("alea_bc") {
		id = id.Sub()

		queueIdx, err := strconv.ParseUint(string(id.Top()), 10, 64)
		if err != nil {
			return
		}
		queueSlot, err := strconv.ParseUint(string(id.Sub().Top()), 10, 64)
		if err != nil {
			return
		}
		at.startBcModSpan(ts, commontypes.Slot{
			QueueIdx:  aleatypes.QueueIdx(queueIdx),
			QueueSlot: aleatypes.QueueSlot(queueSlot),
		})

		// coincides with the start, more or less
		at.startBcSpan(ts, commontypes.Slot{
			QueueIdx:  aleatypes.QueueIdx(queueIdx),
			QueueSlot: aleatypes.QueueSlot(queueSlot),
		})
	}
}

func (at *AleaTracer) parseAbbaRoundID(id t.ModuleID) (abbaRoundID, error) {
	if !id.IsSubOf("alea_ag") {
		return abbaRoundID{}, fmt.Errorf("not alea_ag/*")
	}
	id = id.Sub()

	agRound, err := strconv.ParseUint(string(id.Top()), 10, 64)
	if err != nil {
		return abbaRoundID{}, fmt.Errorf("expected <ag round no.>, got %s", id.Top())
	}
	id = id.Sub()

	if id.Top() != "r" {
		return abbaRoundID{}, fmt.Errorf("expected r, got %s", id.Top())
	}

	id = id.Sub()
	abbaRound, err := strconv.ParseUint(string(id.Top()), 10, 64)
	if err != nil {
		return abbaRoundID{}, fmt.Errorf("expected <abba round number>, got %s", id.Top())
	}

	return abbaRoundID{agRound, abbaRound}, nil
}

func (at *AleaTracer) slotForAgRound(round uint64) commontypes.Slot {
	queueIdx := round % uint64(at.nodeCount)
	queueSlot := at.agQueueHeads[queueIdx]

	if wipRound, ok := at.agWipForQueue[aleatypes.QueueIdx(queueIdx)]; ok && wipRound != round {
		panic("concurrent ag rounds for same queue")
	}

	return commontypes.Slot{
		QueueIdx:  aleatypes.QueueIdx(queueIdx),
		QueueSlot: queueSlot,
	}
}

func (at *AleaTracer) _bcSpan(ts time.Duration, slot commontypes.Slot) *span {
	s, ok := at.wipBcSpan[slot]
	if !ok {
		at.wipBcSpan[slot] = &span{
			class: "bc",
			id:    fmt.Sprintf("%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBcSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBcSpan(ts time.Duration, slot commontypes.Slot) {
	at._bcSpan(ts, slot)
}
func (at *AleaTracer) endBcSpan(ts time.Duration, slot commontypes.Slot) {
	s := at._bcSpan(ts, slot)
	if s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bcAwaitEchoSpan(ts time.Duration, slot commontypes.Slot) *span {
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
func (at *AleaTracer) startBcAwaitEchoSpan(ts time.Duration, slot commontypes.Slot) {
	at._bcAwaitEchoSpan(ts, slot)
}
func (at *AleaTracer) endBcAwaitEchoSpan(ts time.Duration, slot commontypes.Slot) {
	s := at._bcAwaitEchoSpan(ts, slot)
	if s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bcAwaitFinalSpan(ts time.Duration, slot commontypes.Slot) *span {
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
func (at *AleaTracer) startBcAwaitFinalSpan(ts time.Duration, slot commontypes.Slot) {
	at._bcAwaitFinalSpan(ts, slot)
}
func (at *AleaTracer) endBcAwaitFinalSpan(ts time.Duration, slot commontypes.Slot) {
	s := at._bcAwaitFinalSpan(ts, slot)
	if s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bcModSpan(ts time.Duration, slot commontypes.Slot) *span {
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
func (at *AleaTracer) startBcModSpan(ts time.Duration, slot commontypes.Slot) {
	at._bcModSpan(ts, slot)
}
func (at *AleaTracer) endBcModSpan(ts time.Duration, slot commontypes.Slot) {
	s := at._bcModSpan(ts, slot)
	if s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _agSpan(ts time.Duration, round uint64) *span {
	slot := at.slotForAgRound(round)

	s, ok := at.wipAgSpan[round]
	if !ok {
		at.wipAgSpan[round] = &span{
			class: "ag",
			id:    fmt.Sprintf("%d/%d", slot.QueueIdx, slot.QueueSlot),
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
	s := at._agSpan(ts, round)
	if s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
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
	s := at._agModSpan(ts, round)
	if s.end == 0 {
		s.end = ts

		for abbaRound := range at.agUndeliveredAbbaRounds[round] {
			at.endAbbaRoundSpan(ts, abbaRoundID{round, abbaRound})
			at.endAbbaRoundModSpan(ts, abbaRoundID{round, abbaRound})
		}
		delete(at.agUndeliveredAbbaRounds, round)

		at.writeSpan(s)
	}
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
	s := at._abbaRoundSpan(ts, id)
	if s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
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
	s := at._abbaRoundModSpan(ts, id)
	if s.end == 0 {
		s.end = ts
		delete(at.agUndeliveredAbbaRounds[id.agRound], id.abbaRound)
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bfSpan(ts time.Duration, slot commontypes.Slot) *span {
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
func (at *AleaTracer) startBfSpan(ts time.Duration, slot commontypes.Slot) {
	at._bfSpan(ts, slot)
}
func (at *AleaTracer) endBfSpan(ts time.Duration, slot commontypes.Slot) {
	s := at._bfSpan(ts, slot)
	if s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bfStalledSpan(ts time.Duration, slot commontypes.Slot) *span {
	s, ok := at.wipBfStalledSpan[slot]
	if !ok {
		at.wipBfStalledSpan[slot] = &span{
			class: "deliveryStalled",
			id:    fmt.Sprintf("%d:%d", slot.QueueIdx, slot.QueueSlot),
			start: ts,
		}
		s = at.wipBfStalledSpan[slot]
	}
	return s
}
func (at *AleaTracer) startDeliveryStalledSpan(ts time.Duration, slot commontypes.Slot) {
	at._bfStalledSpan(ts, slot)
}
func (at *AleaTracer) endDeliveryStalledSpan(ts time.Duration, slot commontypes.Slot) {
	s := at._bfStalledSpan(ts, slot)
	if s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
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
	s := at._threshCryptoSpan(ts, class, destModule, ctxID)
	if s.end == 0 {
		s.end = ts
		at.writeSpan(s)
	}
}

func (at *AleaTracer) writeSpan(s *span) {
	_, err := at.out.Write([]byte(fmt.Sprintf("%s,%s,%d,%d\n", s.class, s.id, s.start, s.end)))
	if err != nil {
		panic(fmt.Errorf("failed to write trace span: %w", err))
	}
}

func parseSlotFromModuleID(moduleIDStr string) commontypes.Slot {
	modID := t.ModuleID(moduleIDStr)
	modID = modID.StripParent("alea_bc")
	queueIdxStr := string(modID.Top())
	queueSlotStr := string(modID.Sub().Top())

	queueIdx, err := strconv.ParseUint(queueIdxStr, 10, 32)
	if err != nil {
		panic(fmt.Errorf("failed to parse queueIdx from %s: %w", queueIdxStr, err))
	}

	queueSlot, err := strconv.ParseUint(queueSlotStr, 10, 32)
	if err != nil {
		panic(fmt.Errorf("failed to parse queueSlot from %s: %w", queueSlotStr, err))
	}

	return commontypes.Slot{
		QueueIdx:  aleatypes.QueueIdx(queueIdx),
		QueueSlot: aleatypes.QueueSlot(queueSlot),
	}
}
