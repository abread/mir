package vcbpbmsgs

import (
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types2 "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func SendMessage(destModule types.ModuleID, txs []*requestpb.Request) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Vcb{
			Vcb: &types2.Message{
				Type: &types2.Message_SendMessage{
					SendMessage: &types2.SendMessage{
						Txs: txs,
					},
				},
			},
		},
	}
}

func EchoMessage(destModule types.ModuleID, signatureShare tctypes.SigShare) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Vcb{
			Vcb: &types2.Message{
				Type: &types2.Message_EchoMessage{
					EchoMessage: &types2.EchoMessage{
						SignatureShare: signatureShare,
					},
				},
			},
		},
	}
}

func FinalMessage(destModule types.ModuleID, txs []*requestpb.Request, signature tctypes.FullSig) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Vcb{
			Vcb: &types2.Message{
				Type: &types2.Message_FinalMessage{
					FinalMessage: &types2.FinalMessage{
						Txs:       txs,
						Signature: signature,
					},
				},
			},
		},
	}
}
