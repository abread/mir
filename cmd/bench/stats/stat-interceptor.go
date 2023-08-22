// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package stats

import (
	"time"

	"github.com/filecoin-project/mir/pkg/events"
	ageventstypes "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/types"
	bcqueuepbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	batchfetcherpbtypes "github.com/filecoin-project/mir/pkg/pb/batchfetcherpb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	mempoolpbtypes "github.com/filecoin-project/mir/pkg/pb/mempoolpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type StatInterceptor struct {
	*Stats

	// ID of the module that is consuming the transactions.
	// Statistics will only be performed on transactions destined to this module
	// and the rest of the events will be ignored by the StatInterceptor.
	txConsumerModule t.ModuleID
}

var timeRef = time.Now()

func NewStatInterceptor(s *Stats, txConsumer t.ModuleID) *StatInterceptor {
	return &StatInterceptor{s, txConsumer}
}

func (i *StatInterceptor) Intercept(events events.EventList) error {

	// Avoid nil dereference if Intercept is called on a nil *Recorder and simply do nothing.
	// This can happen if a pointer type to *Recorder is assigned to a variable with the interface type Interceptor.
	// Mir would treat that variable as non-nil, thinking there is an interceptor, and call Intercept() on it.
	// For more explanation, see https://mangatmodi.medium.com/go-check-nil-interface-the-right-way-d142776edef1
	if i == nil {
		return nil
	}

	it := events.Iterator()
	for evt := it.Next(); evt != nil; evt = it.Next() {
		switch e := evt.Type.(type) {
		case *eventpbtypes.Event_Mempool:
			switch e := e.Mempool.Type.(type) {
			case *mempoolpbtypes.Event_NewTransactions:
				ts := time.Since(timeRef)
				for _, tx := range e.NewTransactions.Transactions {
					i.Stats.NewTX(tx, int64(ts))
				}
			case *mempoolpbtypes.Event_NewBatch:
				i.Stats.MempoolNewBatch()
			}
		case *eventpbtypes.Event_BatchFetcher:

			// Skip events destined to other modules than the one consuming the transactions.
			if evt.DestModule != i.txConsumerModule {
				continue
			}

			switch e := e.BatchFetcher.Type.(type) {
			case *batchfetcherpbtypes.Event_NewOrderedBatch:
				ts := time.Since(timeRef)
				for _, tx := range e.NewOrderedBatch.Txs {
					i.Stats.Delivered(tx, int64(ts))
				}
			}
		case *eventpbtypes.Event_AleaAgreement:
			if e2, ok := e.AleaAgreement.Type.(*ageventstypes.Event_Deliver); ok {
				i.Stats.DeliveredAgRound(e2.Deliver)
			}
		case *eventpbtypes.Event_AleaBcqueue:
			if _, ok := e.AleaBcqueue.Type.(*bcqueuepbtypes.Event_Deliver); ok {
				i.Stats.DeliveredBcSlot()
			}
		case *eventpbtypes.Event_ThreshCrypto:
			i.Stats.ThreshCryptoEvent(e.ThreshCrypto)
		}
	}

	return nil
}
