package vcb

import (
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	"github.com/filecoin-project/mir/pkg/pb/vcbpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func Message(moduleID t.ModuleID, msg *vcbpb.Message) *messagepb.Message {
	return &messagepb.Message{
		DestModule: moduleID.Pb(),
		Type: &messagepb.Message_Vcb{
			Vcb: msg,
		},
	}
}

func SendMessage(moduleID t.ModuleID, data []*requestpb.Request) *messagepb.Message {
	return Message(moduleID, &vcbpb.Message{
		Type: &vcbpb.Message_SendMessage{
			SendMessage: &vcbpb.SendMessage{Data: data},
		},
	})
}

func EchoMessage(moduleID t.ModuleID, signatureShare []byte) *messagepb.Message {
	return Message(moduleID, &vcbpb.Message{
		Type: &vcbpb.Message_EchoMessage{
			EchoMessage: &vcbpb.EchoMessage{
				SignatureShare: signatureShare,
			},
		},
	})
}

func FinalMessage(moduleID t.ModuleID, data []*requestpb.Request, signature []byte) *messagepb.Message {
	return Message(moduleID, &vcbpb.Message{
		Type: &vcbpb.Message_FinalMessage{
			FinalMessage: &vcbpb.FinalMessage{
				Data:      data,
				Signature: signature,
			},
		},
	})
}
