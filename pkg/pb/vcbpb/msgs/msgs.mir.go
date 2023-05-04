package vcbpbmsgs

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func SendMessage(destModule types.ModuleID, txs []*types1.Transaction) *types2.Message {
	return &types2.Message{
		DestModule: destModule,
		Type: &types2.Message_Vcb{
			Vcb: &types3.Message{
				Type: &types3.Message_SendMessage{
					SendMessage: &types3.SendMessage{
						Txs: txs,
					},
				},
			},
		},
	}
}

func EchoMessage(destModule types.ModuleID, signatureShare tctypes.SigShare) *types2.Message {
	return &types2.Message{
		DestModule: destModule,
		Type: &types2.Message_Vcb{
			Vcb: &types3.Message{
				Type: &types3.Message_EchoMessage{
					EchoMessage: &types3.EchoMessage{
						SignatureShare: signatureShare,
					},
				},
			},
		},
	}
}

func FinalMessage(destModule types.ModuleID, signature tctypes.FullSig) *types2.Message {
	return &types2.Message{
		DestModule: destModule,
		Type: &types2.Message_Vcb{
			Vcb: &types3.Message{
				Type: &types3.Message_FinalMessage{
					FinalMessage: &types3.FinalMessage{
						Signature: signature,
					},
				},
			},
		},
	}
}
