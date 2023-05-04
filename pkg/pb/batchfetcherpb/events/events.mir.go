package batchfetcherpbevents

import (
	types3 "github.com/filecoin-project/mir/pkg/pb/batchfetcherpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func NewOrderedBatch(destModule types.ModuleID, txs []*types1.Transaction) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_BatchFetcher{
			BatchFetcher: &types3.Event{
				Type: &types3.Event_NewOrderedBatch{
					NewOrderedBatch: &types3.NewOrderedBatch{
						Txs: txs,
					},
				},
			},
		},
	}
}
