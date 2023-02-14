package vcbpbevents

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types1 "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule types.ModuleID, txs []*requestpb.Request, origin *types1.Origin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Vcb{
			Vcb: &types1.Event{
				Type: &types1.Event_InputValue{
					InputValue: &types1.InputValue{
						Txs:    txs,
						Origin: origin,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, txs []*requestpb.Request, txIds []types.TxID, signature tctypes.FullSig, origin *types1.Origin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Vcb{
			Vcb: &types1.Event{
				Type: &types1.Event_Deliver{
					Deliver: &types1.Deliver{
						Txs:       txs,
						TxIds:     txIds,
						Signature: signature,
						Origin:    origin,
					},
				},
			},
		},
	}
}
