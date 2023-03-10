package agreementpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	agreementpb "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() agreementpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb agreementpb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *agreementpb.Message_FinishAbba:
		return &Message_FinishAbba{FinishAbba: FinishAbbaMessageFromPb(pb.FinishAbba)}
	}
	return nil
}

type Message_FinishAbba struct {
	FinishAbba *FinishAbbaMessage
}

func (*Message_FinishAbba) isMessage_Type() {}

func (w *Message_FinishAbba) Unwrap() *FinishAbbaMessage {
	return w.FinishAbba
}

func (w *Message_FinishAbba) Pb() agreementpb.Message_Type {
	return &agreementpb.Message_FinishAbba{FinishAbba: (w.FinishAbba).Pb()}
}

func (*Message_FinishAbba) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agreementpb.Message_FinishAbba]()}
}

func MessageFromPb(pb *agreementpb.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *agreementpb.Message {
	return &agreementpb.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agreementpb.Message]()}
}

type FinishAbbaMessage struct {
	Round uint64
	Value bool
}

func FinishAbbaMessageFromPb(pb *agreementpb.FinishAbbaMessage) *FinishAbbaMessage {
	return &FinishAbbaMessage{
		Round: pb.Round,
		Value: pb.Value,
	}
}

func (m *FinishAbbaMessage) Pb() *agreementpb.FinishAbbaMessage {
	return &agreementpb.FinishAbbaMessage{
		Round: m.Round,
		Value: m.Value,
	}
}

func (*FinishAbbaMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agreementpb.FinishAbbaMessage]()}
}
