package vcbpbevents

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule types.ModuleID, txs []*types1.Request) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Vcb{
			Vcb: &types3.Event{
				Type: &types3.Event_InputValue{
					InputValue: &types3.InputValue{
						Txs: txs,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, txs []*types1.Request, txIds []types.TxID, signature tctypes.FullSig, srcModule types.ModuleID) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Vcb{
			Vcb: &types3.Event{
				Type: &types3.Event_Deliver{
					Deliver: &types3.Deliver{
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