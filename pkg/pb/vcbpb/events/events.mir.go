package vcbpbevents

import (
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types2 "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule types.ModuleID, txs []*requestpb.Request) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Vcb{
			Vcb: &types2.Event{
				Type: &types2.Event_InputValue{
					InputValue: &types2.InputValue{
						Txs: txs,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, txs []*requestpb.Request, txIds []types.TxID, signature tctypes.FullSig, srcModule types.ModuleID) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Vcb{
			Vcb: &types2.Event{
				Type: &types2.Event_Deliver{
					Deliver: &types2.Deliver{
						Txs:       txs,
						TxIds:     txIds,
						Signature: signature,
						SrcModule: srcModule,
					},
				},
			},
		},
	}
}
