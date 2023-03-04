package aleatracer

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/abbapb"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/batchfetcherpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/mempoolpb"
	"github.com/filecoin-project/mir/pkg/pb/vcbpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type AleaTracer struct {
	refTime     time.Time
	ownQueueIdx aleatypes.QueueIdx
	nodeCount   int

	wipTxSpan                  map[string]*span
	wipTxWaitingLocalBatchSpan map[string]*span
	wipBcSpan                  map[commontypes.Slot]*span
	wipBcModSpan               map[commontypes.Slot]*span
	wipAgSpan                  map[uint64]*span
	wipAgModSpan               map[uint64]*span
	wipAbbaRoundSpan           map[abbaRoundID]*span
	wipAbbaRoundModSpan        map[abbaRoundID]*span
	wipBfSpan                  map[commontypes.Slot]*span
	wipBfStalledSpan           map[commontypes.Slot]*span

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
		refTime:     time.Now(),
		ownQueueIdx: ownQueueIdx,
		nodeCount:   nodeCount,

		wipTxSpan:                  make(map[string]*span, 1*N),
		wipTxWaitingLocalBatchSpan: make(map[string]*span, 1*N),
		wipBcSpan:                  make(map[commontypes.Slot]*span, nodeCount*N),
		wipBcModSpan:               make(map[commontypes.Slot]*span, nodeCount*N),
		wipAgSpan:                  make(map[uint64]*span, nodeCount*N),
		wipAgModSpan:               make(map[uint64]*span, nodeCount*N),
		wipAbbaRoundSpan:           make(map[abbaRoundID]*span, nodeCount*N),
		wipAbbaRoundModSpan:        make(map[abbaRoundID]*span, nodeCount*N),
		wipBfSpan:                  make(map[commontypes.Slot]*span, nodeCount*N),
		wipBfStalledSpan:           make(map[commontypes.Slot]*span, nodeCount*N),

		bfWipSlotsCtxID: make(map[uint64]commontypes.Slot, nodeCount),
		bfDeliverQueue:  nil,

		agQueueHeads:            make([]aleatypes.QueueSlot, nodeCount),
		agWipForQueue:           make(map[aleatypes.QueueIdx]uint64, nodeCount),
		agUndeliveredAbbaRounds: make(map[uint64]map[uint64]struct{}, nodeCount*N),

		inChan: make(chan *events.EventList, nodeCount*16),
		cancel: cancel,
		out:    out,
	}

	_, err := out.Write([]byte("class,id,start,end\n"))
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
	if _, ok := event.Type.(*eventpb.Event_MessageReceived); !ok {
		at.registerModStart(event)
	}

	switch ev := event.Type.(type) {
	case *eventpb.Event_NewRequests:
		for _, tx := range ev.NewRequests.Requests {
			at.startTxSpan(tx.ClientId, tx.ReqNo)
			at.startTxWaitingLocalBatchSpan(tx.ClientId, tx.ReqNo)
		}
	case *eventpb.Event_Mempool:
		switch e := ev.Mempool.Type.(type) {
		case *mempoolpb.Event_NewBatch:
			for _, tx := range e.NewBatch.Txs {
				at.endTxWaitingLocalBatchSpan(tx.ClientId, tx.ReqNo)
			}
		}
	case *eventpb.Event_AleaBroadcast:
		switch e := ev.AleaBroadcast.Type.(type) {
		case *bcpb.Event_StartBroadcast:
			slot := commontypes.Slot{QueueIdx: at.ownQueueIdx, QueueSlot: aleatypes.QueueSlot(e.StartBroadcast.QueueSlot)}
			at.startBcSpan(slot)
		case *bcpb.Event_Deliver:
			at.endBcSpan(*commontypes.SlotFromPb(e.Deliver.Slot))
		case *bcpb.Event_FreeSlot:
			at.endBcModSpan(*commontypes.SlotFromPb(e.FreeSlot.Slot))
		}
	case *eventpb.Event_AleaAgreement:
		switch e := ev.AleaAgreement.Type.(type) {
		case *agevents.Event_InputValue:
			slot := at.slotForAgRound(e.InputValue.Round)
			at.agWipForQueue[slot.QueueIdx] = e.InputValue.Round

			at.startAgSpan(e.InputValue.Round)
		case *agevents.Event_Deliver:
			at.endAgSpan(e.Deliver.Round)
			at.endAgModSpan(e.Deliver.Round)

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
			at.startBcSpan(slot)
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
				at.startAbbaRoundSpan(abbaRoundID)
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

				at.endAbbaRoundSpan(abbaRoundID)
				at.endAbbaRoundModSpan(abbaRoundID)

				delete(at.agUndeliveredAbbaRounds[abbaRoundID.agRound], abbaRoundID.abbaRound)
			}
		}
	case *eventpb.Event_Availability:
		switch e := ev.Availability.Type.(type) {
		case *availabilitypb.Event_RequestTransactions:
			slot := *commontypes.SlotFromPb(e.RequestTransactions.Cert.Type.(*availabilitypb.Cert_Alea).Alea.Slot)

			at.startBfSpan(slot)

			ctxID := e.RequestTransactions.Origin.GetDsl().ContextID
			at.bfWipSlotsCtxID[ctxID] = slot
		case *availabilitypb.Event_ProvideTransactions:
			ctxID := e.ProvideTransactions.Origin.GetDsl().ContextID
			slot := at.bfWipSlotsCtxID[ctxID]
			delete(at.bfWipSlotsCtxID, ctxID)

			at.endBfSpan(slot)
			at.startDeliveryStalledSpan(slot)
		}
	case *eventpb.Event_BatchFetcher:
		switch e := ev.BatchFetcher.Type.(type) {
		case *batchfetcherpb.Event_NewOrderedBatch:
			slot := at.bfDeliverQueue[0]
			at.bfDeliverQueue = at.bfDeliverQueue[1:]

			at.endDeliveryStalledSpan(slot)

			for _, tx := range e.NewOrderedBatch.Txs {
				at.endTxSpan(tx.ClientId, tx.ReqNo)
			}
		}
	}

	return nil
}

func (at *AleaTracer) registerModStart(event *eventpb.Event) {
	id := t.ModuleID(event.DestModule)
	if id.IsSubOf("alea_ag") {
		id = id.Sub()

		agRound, err := strconv.ParseUint(string(id.Top()), 10, 64)
		if err != nil {
			return
		}

		at.startAgModSpan(agRound)

		if id.Sub().Top() == "r" {
			id = id.Sub().Sub()
			abbaRound, err := strconv.ParseUint(string(id), 10, 64)
			if err != nil {
				return
			}

			at.startAbbaRoundModSpan(abbaRoundID{agRound, abbaRound})
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
		at.startBcModSpan(commontypes.Slot{
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

func (at *AleaTracer) _txSpan(clientID string, reqNo uint64) *span {
	id := fmt.Sprintf("%s:%d", clientID, reqNo)

	s, ok := at.wipTxSpan[id]
	if !ok {
		at.wipTxSpan[id] = &span{
			class: "tx",
			id:    id,
			start: time.Since(at.refTime),
		}
		s = at.wipTxSpan[id]
	}
	return s
}
func (at *AleaTracer) startTxSpan(clientID string, reqNo uint64) {
	at._txSpan(clientID, reqNo)
}
func (at *AleaTracer) endTxSpan(clientID string, reqNo uint64) {
	s := at._txSpan(clientID, reqNo)
	if s.end == 0 {
		s.end = time.Since(at.refTime)
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _txWaitingLocalBatchSpan(clientID string, reqNo uint64) *span {
	id := fmt.Sprintf("%s:%d", clientID, reqNo)

	s, ok := at.wipTxSpan[id]
	if !ok {
		at.wipTxSpan[id] = &span{
			class: "txWaitingLocalBatch",
			id:    id,
			start: time.Since(at.refTime),
		}
		s = at.wipTxSpan[id]
	}
	return s
}
func (at *AleaTracer) startTxWaitingLocalBatchSpan(clientID string, reqNo uint64) {
	at._txWaitingLocalBatchSpan(clientID, reqNo)
}
func (at *AleaTracer) endTxWaitingLocalBatchSpan(clientID string, reqNo uint64) {
	s := at._txWaitingLocalBatchSpan(clientID, reqNo)
	if s.end == 0 {
		s.end = time.Since(at.refTime)
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bcSpan(slot commontypes.Slot) *span {
	s, ok := at.wipBcSpan[slot]
	if !ok {
		at.wipBcSpan[slot] = &span{
			class: "bc",
			id:    fmt.Sprintf("%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: time.Since(at.refTime),
		}
		s = at.wipBcSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBcSpan(slot commontypes.Slot) {
	at._bcSpan(slot)
}
func (at *AleaTracer) endBcSpan(slot commontypes.Slot) {
	s := at._bcSpan(slot)
	if s.end == 0 {
		s.end = time.Since(at.refTime)
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bcModSpan(slot commontypes.Slot) *span {
	s, ok := at.wipBcModSpan[slot]
	if !ok {
		at.wipBcModSpan[slot] = &span{
			class: "mod",
			id:    fmt.Sprintf("alea_bc/%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: time.Since(at.refTime),
		}
		s = at.wipBcModSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBcModSpan(slot commontypes.Slot) {
	at._bcModSpan(slot)
}
func (at *AleaTracer) endBcModSpan(slot commontypes.Slot) {
	s := at._bcModSpan(slot)
	if s.end == 0 {
		s.end = time.Since(at.refTime)
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _agSpan(round uint64) *span {
	slot := at.slotForAgRound(round)

	s, ok := at.wipAgSpan[round]
	if !ok {
		at.wipAgSpan[round] = &span{
			class: "ag",
			id:    fmt.Sprintf("%d/%d", slot.QueueIdx, slot.QueueSlot),
			start: time.Since(at.refTime),
		}
		s = at.wipAgSpan[round]
	}
	return s
}
func (at *AleaTracer) startAgSpan(round uint64) {
	at._agSpan(round)
}
func (at *AleaTracer) endAgSpan(round uint64) {
	s := at._agSpan(round)
	if s.end == 0 {
		s.end = time.Since(at.refTime)
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _agModSpan(round uint64) *span {
	s, ok := at.wipAgModSpan[round]
	if !ok {
		at.wipAgModSpan[round] = &span{
			class: "mod",
			id:    fmt.Sprintf("alea_ag/%d", round),
			start: time.Since(at.refTime),
		}
		at.agUndeliveredAbbaRounds[round] = make(map[uint64]struct{}, 8)
		s = at.wipAgModSpan[round]
	}
	return s
}
func (at *AleaTracer) startAgModSpan(round uint64) {
	at._agModSpan(round)
}
func (at *AleaTracer) endAgModSpan(round uint64) {
	s := at._agModSpan(round)
	if s.end == 0 {
		s.end = time.Since(at.refTime)

		for abbaRound := range at.agUndeliveredAbbaRounds[round] {
			at.endAbbaRoundSpan(abbaRoundID{round, abbaRound})
			at.endAbbaRoundModSpan(abbaRoundID{round, abbaRound})
		}
		delete(at.agUndeliveredAbbaRounds, round)

		at.writeSpan(s)
	}
}

func (at *AleaTracer) _abbaRoundSpan(id abbaRoundID) *span {
	s, ok := at.wipAbbaRoundSpan[id]
	if !ok {
		at.wipAbbaRoundSpan[id] = &span{
			class: "abbaRound",
			id:    fmt.Sprintf("%d-%d", id.agRound, id.abbaRound),
			start: time.Since(at.refTime),
		}
		s = at.wipAbbaRoundSpan[id]
	}
	return s
}
func (at *AleaTracer) startAbbaRoundSpan(id abbaRoundID) {
	at._abbaRoundSpan(id)
}
func (at *AleaTracer) endAbbaRoundSpan(id abbaRoundID) {
	s := at._abbaRoundSpan(id)
	if s.end == 0 {
		s.end = time.Since(at.refTime)
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _abbaRoundModSpan(id abbaRoundID) *span {
	s, ok := at.wipAbbaRoundModSpan[id]
	if !ok {
		at.wipAbbaRoundModSpan[id] = &span{
			class: "mod",
			id:    fmt.Sprintf("alea_ag/%d/%d", id.agRound, id.abbaRound),
			start: time.Since(at.refTime),
		}
		at.agUndeliveredAbbaRounds[id.agRound][id.abbaRound] = struct{}{}
		s = at.wipAbbaRoundModSpan[id]
	}
	return s
}
func (at *AleaTracer) startAbbaRoundModSpan(id abbaRoundID) {
	at._abbaRoundModSpan(id)
}
func (at *AleaTracer) endAbbaRoundModSpan(id abbaRoundID) {
	s := at._abbaRoundModSpan(id)
	if s.end == 0 {
		s.end = time.Since(at.refTime)
		delete(at.agUndeliveredAbbaRounds[id.agRound], id.abbaRound)
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bfSpan(slot commontypes.Slot) *span {
	s, ok := at.wipBfSpan[slot]
	if !ok {
		at.wipBfSpan[slot] = &span{
			class: "bf",
			id:    fmt.Sprintf("%d:%d", slot.QueueIdx, slot.QueueSlot),
			start: time.Since(at.refTime),
		}
		s = at.wipBfSpan[slot]
	}
	return s
}
func (at *AleaTracer) startBfSpan(slot commontypes.Slot) {
	at._bfSpan(slot)
}
func (at *AleaTracer) endBfSpan(slot commontypes.Slot) {
	s := at._bfSpan(slot)
	if s.end == 0 {
		s.end = time.Since(at.refTime)
		at.writeSpan(s)
	}
}

func (at *AleaTracer) _bfStalledSpan(slot commontypes.Slot) *span {
	s, ok := at.wipBfStalledSpan[slot]
	if !ok {
		at.wipBfStalledSpan[slot] = &span{
			class: "deliveryStalled",
			id:    fmt.Sprintf("%d:%d", slot.QueueIdx, slot.QueueSlot),
			start: time.Since(at.refTime),
		}
		s = at.wipBfStalledSpan[slot]
	}
	return s
}
func (at *AleaTracer) startDeliveryStalledSpan(slot commontypes.Slot) {
	at._bfStalledSpan(slot)
}
func (at *AleaTracer) endDeliveryStalledSpan(slot commontypes.Slot) {
	s := at._bfStalledSpan(slot)
	if s.end == 0 {
		s.end = time.Since(at.refTime)
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
	queueSlotStr := string(modID.Sub())

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
