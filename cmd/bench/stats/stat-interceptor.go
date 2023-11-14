// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package stats

import (
	"github.com/filecoin-project/mir/pkg/events"
	abbapbtypes "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	ageventstypes "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/types"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	batchfetcherpbtypes "github.com/filecoin-project/mir/pkg/pb/batchfetcherpb/types"
	mempoolpbtypes "github.com/filecoin-project/mir/pkg/pb/mempoolpb/types"

	//batchfetcherpbtypes "github.com/filecoin-project/mir/pkg/pb/batchfetcherpb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	//mempoolpbtypes "github.com/filecoin-project/mir/pkg/pb/mempoolpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type StatInterceptor struct {
	*ReplicaStats
	*ClientStats
	*ClientOptLatStats

	// ID of the module that is consuming the transactions.
	// Statistics will only be performed on transactions destined to this module
	// and the rest of the events will be ignored by the StatInterceptor.
	txConsumerModule t.ModuleID

	ownClientIDPrefix string
}

func NewStatInterceptor(s *ReplicaStats, cs *ClientStats, cols *ClientOptLatStats, txConsumer t.ModuleID, ownClientIDPrefix string) *StatInterceptor {
	return &StatInterceptor{s, cs, cols, txConsumer, ownClientIDPrefix}
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
		i.interceptOne(evt)
	}

	return nil
}

func (i *StatInterceptor) interceptOne(evt *eventpbtypes.Event) {
	switch e := evt.Type.(type) {
	case *eventpbtypes.Event_Mempool:
		switch e := e.Mempool.Type.(type) {
		case *mempoolpbtypes.Event_NewTransactions:
			i.ReplicaStats.NewTransactions(e.NewTransactions.Transactions)
		case *mempoolpbtypes.Event_NewBatch:
			i.ReplicaStats.CutBatch(len(e.NewBatch.Txs))
			i.ClientOptLatStats.CutBatch(e.NewBatch.Txs)
		}
	case *eventpbtypes.Event_BatchFetcher:

		// Skip events destined to other modules than the one consuming the transactions.
		if evt.DestModule != i.txConsumerModule {
			return
		}

		switch e := e.BatchFetcher.Type.(type) {
		case *batchfetcherpbtypes.Event_NewOrderedBatch:
			i.ReplicaStats.DeliverBatch(e.NewOrderedBatch.Txs)

			for _, tx := range e.NewOrderedBatch.Txs {
				i.ClientOptLatStats.Deliver(tx)
			}
			i.ClientOptLatStats.DeliveredBatch()

			i.ClientStats.DeliveredBatch()
		}
	case *eventpbtypes.Event_AleaAgreement:
		switch e2 := e.AleaAgreement.Type.(type) {
		case *ageventstypes.Event_InputValue:
			i.ReplicaStats.AgInput(e2.InputValue)
		case *ageventstypes.Event_Deliver:
			i.ReplicaStats.AgDeliver(e2.Deliver)
		case *ageventstypes.Event_InnerAbbaRoundTime:
			i.ReplicaStats.InnerAbbaTime(e2.InnerAbbaRoundTime)
		}
	case *eventpbtypes.Event_Abba:
		switch e2 := e.Abba.Type.(type) {
		case *abbapbtypes.Event_Round:
			switch e3 := e2.Round.Type.(type) {
			case *abbapbtypes.RoundEvent_Continue:
				// TODO: don't depend on module structure details
				// assumes abba round module IDs have form <ag mod id>/<ag round id>/r/<abba round id>
				abbaRoundIDStr := string(evt.DestModule.Sub().Sub().Sub().Top())

				i.ReplicaStats.AbbaRoundContinue(abbaRoundIDStr, e3.Continue)
			}
		}
	case *eventpbtypes.Event_AleaBc:
		switch e2 := e.AleaBc.Type.(type) {
		case *bcpbtypes.Event_RequestCert:
			i.ReplicaStats.RequestCert()
		case *bcpbtypes.Event_DeliverCert:
			i.ReplicaStats.BcDeliver(e2.DeliverCert)
		}
	}
}
