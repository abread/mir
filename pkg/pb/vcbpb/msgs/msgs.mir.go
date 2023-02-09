package vcbpbmsgs

import (
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types2 "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
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

func EchoMessage(destModule types.ModuleID, signatureShare []uint8) *types1.Message {
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

func FinalMessage(destModule types.ModuleID, txs []*requestpb.Request, signature []uint8) *types1.Message {
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
