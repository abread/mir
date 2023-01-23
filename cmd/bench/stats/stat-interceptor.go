// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package stats

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	bfpb "github.com/filecoin-project/mir/pkg/pb/batchfetcherpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	t "github.com/filecoin-project/mir/pkg/types"
)

type StatInterceptor struct {
	*Stats

	fakeRNet *reliablenet.Module

	// ID of the module that is consuming the transactions.
	// Statistics will only be performed on transactions destined to this module
	// and the rest of the events will be ignored by the StatInterceptor.
	txConsumerModule t.ModuleID
}

func NewStatInterceptor(s *Stats, txConsumer t.ModuleID, ownID t.NodeID, allNodes []t.NodeID) *StatInterceptor {
	fakeRNet, err := reliablenet.New(ownID, reliablenet.DefaultModuleConfig(), reliablenet.DefaultModuleParams(allNodes), logging.NilLogger)
	if err != nil {
		panic(err)
	}

	return &StatInterceptor{s, fakeRNet, txConsumer}
}

func (i *StatInterceptor) Intercept(evts *events.EventList) error {
	it := evts.Iterator()
	for evt := it.Next(); evt != nil; evt = it.Next() {

		switch e := evt.Type.(type) {
		case *eventpb.Event_NewRequests:
			for _, req := range e.NewRequests.Requests {
				i.Stats.NewRequest(req)
			}
		case *eventpb.Event_BatchFetcher:

			// Skip events destined to other modules than the one consuming the transactions.
			if t.ModuleID(evt.DestModule) != i.txConsumerModule {
				continue
			}

			switch e := e.BatchFetcher.Type.(type) {
			case *bfpb.Event_NewOrderedBatch:
				for _, req := range e.NewOrderedBatch.Txs {
					i.Stats.Delivered(req)
				}
			}
		case *eventpb.Event_ReliableNet:
			_, _ = i.fakeRNet.ApplyEvents(events.ListOf(evt))
		case *eventpb.Event_MessageReceived:
			if e.MessageReceived.Msg.DestModule == "reliablenet" {
				_, _ = i.fakeRNet.ApplyEvents(events.ListOf(evt))
			}
		}
	}
	return nil
}
