// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package stats

import (
	"github.com/filecoin-project/mir/pkg/events"
	ageventstypes "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/types"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	batchfetcherpbtypes "github.com/filecoin-project/mir/pkg/pb/batchfetcherpb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	mempoolpbtypes "github.com/filecoin-project/mir/pkg/pb/mempoolpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type StatInterceptor struct {
	*LiveStats

	// ID of the module that is consuming the transactions.
	// Statistics will only be performed on transactions destined to this module
	// and the rest of the events will be ignored by the StatInterceptor.
	txConsumerModule t.ModuleID
}

func NewStatInterceptor(s *LiveStats, txConsumer t.ModuleID) *StatInterceptor {
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
				for _, tx := range e.NewTransactions.Transactions {
					i.LiveStats.Submit(tx)
				}
			}
		case *eventpbtypes.Event_BatchFetcher:

			// Skip events destined to other modules than the one consuming the transactions.
			if evt.DestModule != i.txConsumerModule {
				continue
			}

			switch e := e.BatchFetcher.Type.(type) {
			case *batchfetcherpbtypes.Event_NewOrderedBatch:
				for _, tx := range e.NewOrderedBatch.Txs {
					i.LiveStats.Deliver(tx)
				}
			}
		case *eventpbtypes.Event_AleaAgreement:
			switch e2 := e.AleaAgreement.Type.(type) {
			case *ageventstypes.Event_InputValue:
				i.LiveStats.AgInput(e2.InputValue)
			case *ageventstypes.Event_Deliver:
				i.LiveStats.AgDeliver(e2.Deliver)
			case *ageventstypes.Event_InnerAbbaRoundTime:
				i.LiveStats.InnerAbbaTime(e2.InnerAbbaRoundTime)
			}
		case *eventpbtypes.Event_AleaBc:
			if e2, ok := e.AleaBc.Type.(*bcpbtypes.Event_DeliverCert); ok {
				i.LiveStats.BcDeliver(e2.DeliverCert)
			}
		}
	}

	return nil
}
