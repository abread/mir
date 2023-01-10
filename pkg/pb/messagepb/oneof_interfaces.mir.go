package messagepb

import (
	abbapb "github.com/filecoin-project/mir/pkg/pb/abbapb"
	aleapb "github.com/filecoin-project/mir/pkg/pb/aleapb"
	agreementpb "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	mscpb "github.com/filecoin-project/mir/pkg/pb/availabilitypb/mscpb"
	bcbpb "github.com/filecoin-project/mir/pkg/pb/bcbpb"
	checkpointpb "github.com/filecoin-project/mir/pkg/pb/checkpointpb"
	isspb "github.com/filecoin-project/mir/pkg/pb/isspb"
	ordererspb "github.com/filecoin-project/mir/pkg/pb/ordererspb"
	pingpongpb "github.com/filecoin-project/mir/pkg/pb/pingpongpb"
	messages "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages"
	vcbpb "github.com/filecoin-project/mir/pkg/pb/vcbpb"
)

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_Iss) Unwrap() *isspb.ISSMessage {
	return w.Iss
}

func (w *Message_Bcb) Unwrap() *bcbpb.Message {
	return w.Bcb
}

func (w *Message_MultisigCollector) Unwrap() *mscpb.Message {
	return w.MultisigCollector
}

func (w *Message_Pingpong) Unwrap() *pingpongpb.Message {
	return w.Pingpong
}

func (w *Message_Checkpoint) Unwrap() *checkpointpb.Message {
	return w.Checkpoint
}

func (w *Message_SbMessage) Unwrap() *ordererspb.SBInstanceMessage {
	return w.SbMessage
}

func (w *Message_Vcb) Unwrap() *vcbpb.Message {
	return w.Vcb
}

func (w *Message_Abba) Unwrap() *abbapb.Message {
	return w.Abba
}

func (w *Message_Alea) Unwrap() *aleapb.Message {
	return w.Alea
}

func (w *Message_AleaAgreement) Unwrap() *agreementpb.Message {
	return w.AleaAgreement
}

func (w *Message_ReliableNet) Unwrap() *messages.Message {
	return w.ReliableNet
}
