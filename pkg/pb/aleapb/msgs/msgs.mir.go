package aleapbmsgs

import (
	types1 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/aleapb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types4 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func FillGapMessage(destModule types.ModuleID, slot *types1.Slot) *types2.Message {
	return &types2.Message{
		DestModule: destModule,
		Type: &types2.Message_Alea{
			Alea: &types3.Message{
				Type: &types3.Message_FillGapMessage{
					FillGapMessage: &types3.FillGapMessage{
						Slot: slot,
					},
				},
			},
		},
	}
}

func FillerMessage(destModule types.ModuleID, slot *types1.Slot, txs []*types4.Transaction, signature tctypes.FullSig) *types2.Message {
	return &types2.Message{
		DestModule: destModule,
		Type: &types2.Message_Alea{
			Alea: &types3.Message{
				Type: &types3.Message_FillerMessage{
					FillerMessage: &types3.FillerMessage{
						Slot:      slot,
						Txs:       txs,
						Signature: signature,
					},
				},
			},
		},
	}
}
