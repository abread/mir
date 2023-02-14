package messagestypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	messages "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages"
	rntypes "github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	types "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() messages.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb messages.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *messages.Message_Ack:
		return &Message_Ack{Ack: AckMessageFromPb(pb.Ack)}
	}
	return nil
}

type Message_Ack struct {
	Ack *AckMessage
}

func (*Message_Ack) isMessage_Type() {}

func (w *Message_Ack) Unwrap() *AckMessage {
	return w.Ack
}

func (w *Message_Ack) Pb() messages.Message_Type {
	return &messages.Message_Ack{Ack: (w.Ack).Pb()}
}

func (*Message_Ack) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messages.Message_Ack]()}
}

func MessageFromPb(pb *messages.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *messages.Message {
	return &messages.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messages.Message]()}
}

type AckMessage struct {
	MsgDestModule types.ModuleID
	MsgId         rntypes.MsgID
}

func AckMessageFromPb(pb *messages.AckMessage) *AckMessage {
	return &AckMessage{
		MsgDestModule: (types.ModuleID)(pb.MsgDestModule),
		MsgId:         (rntypes.MsgID)(pb.MsgId),
	}
}

func (m *AckMessage) Pb() *messages.AckMessage {
	return &messages.AckMessage{
		MsgDestModule: (string)(m.MsgDestModule),
		MsgId:         (string)(m.MsgId),
	}
}

func (*AckMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messages.AckMessage]()}
}
