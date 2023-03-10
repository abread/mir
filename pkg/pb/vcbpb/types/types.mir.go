package vcbpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types1 "github.com/filecoin-project/mir/codegen/model/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	vcbpb "github.com/filecoin-project/mir/pkg/pb/vcbpb"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() vcbpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb vcbpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *vcbpb.Event_InputValue:
		return &Event_InputValue{InputValue: InputValueFromPb(pb.InputValue)}
	case *vcbpb.Event_Deliver:
		return &Event_Deliver{Deliver: DeliverFromPb(pb.Deliver)}
	}
	return nil
}

type Event_InputValue struct {
	InputValue *InputValue
}

func (*Event_InputValue) isEvent_Type() {}

func (w *Event_InputValue) Unwrap() *InputValue {
	return w.InputValue
}

func (w *Event_InputValue) Pb() vcbpb.Event_Type {
	return &vcbpb.Event_InputValue{InputValue: (w.InputValue).Pb()}
}

func (*Event_InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Event_InputValue]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() vcbpb.Event_Type {
	return &vcbpb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Event_Deliver]()}
}

func EventFromPb(pb *vcbpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *vcbpb.Event {
	return &vcbpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Event]()}
}

type InputValue struct {
	Txs []*requestpb.Request
}

func InputValueFromPb(pb *vcbpb.InputValue) *InputValue {
	return &InputValue{
		Txs: pb.Txs,
	}
}

func (m *InputValue) Pb() *vcbpb.InputValue {
	return &vcbpb.InputValue{
		Txs: m.Txs,
	}
}

func (*InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.InputValue]()}
}

type Deliver struct {
	Txs       []*requestpb.Request
	TxIds     []types.TxID
	Signature tctypes.FullSig
	SrcModule types.ModuleID
}

func DeliverFromPb(pb *vcbpb.Deliver) *Deliver {
	return &Deliver{
		Txs: pb.Txs,
		TxIds: types1.ConvertSlice(pb.TxIds, func(t []uint8) types.TxID {
			return (types.TxID)(t)
		}),
		Signature: (tctypes.FullSig)(pb.Signature),
		SrcModule: (types.ModuleID)(pb.SrcModule),
	}
}

func (m *Deliver) Pb() *vcbpb.Deliver {
	return &vcbpb.Deliver{
		Txs: m.Txs,
		TxIds: types1.ConvertSlice(m.TxIds, func(t types.TxID) []uint8 {
			return ([]uint8)(t)
		}),
		Signature: ([]uint8)(m.Signature),
		SrcModule: (string)(m.SrcModule),
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Deliver]()}
}

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() vcbpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb vcbpb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *vcbpb.Message_SendMessage:
		return &Message_SendMessage{SendMessage: SendMessageFromPb(pb.SendMessage)}
	case *vcbpb.Message_EchoMessage:
		return &Message_EchoMessage{EchoMessage: EchoMessageFromPb(pb.EchoMessage)}
	case *vcbpb.Message_FinalMessage:
		return &Message_FinalMessage{FinalMessage: FinalMessageFromPb(pb.FinalMessage)}
	}
	return nil
}

type Message_SendMessage struct {
	SendMessage *SendMessage
}

func (*Message_SendMessage) isMessage_Type() {}

func (w *Message_SendMessage) Unwrap() *SendMessage {
	return w.SendMessage
}

func (w *Message_SendMessage) Pb() vcbpb.Message_Type {
	return &vcbpb.Message_SendMessage{SendMessage: (w.SendMessage).Pb()}
}

func (*Message_SendMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Message_SendMessage]()}
}

type Message_EchoMessage struct {
	EchoMessage *EchoMessage
}

func (*Message_EchoMessage) isMessage_Type() {}

func (w *Message_EchoMessage) Unwrap() *EchoMessage {
	return w.EchoMessage
}

func (w *Message_EchoMessage) Pb() vcbpb.Message_Type {
	return &vcbpb.Message_EchoMessage{EchoMessage: (w.EchoMessage).Pb()}
}

func (*Message_EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Message_EchoMessage]()}
}

type Message_FinalMessage struct {
	FinalMessage *FinalMessage
}

func (*Message_FinalMessage) isMessage_Type() {}

func (w *Message_FinalMessage) Unwrap() *FinalMessage {
	return w.FinalMessage
}

func (w *Message_FinalMessage) Pb() vcbpb.Message_Type {
	return &vcbpb.Message_FinalMessage{FinalMessage: (w.FinalMessage).Pb()}
}

func (*Message_FinalMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Message_FinalMessage]()}
}

func MessageFromPb(pb *vcbpb.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *vcbpb.Message {
	return &vcbpb.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Message]()}
}

type SendMessage struct {
	Txs []*requestpb.Request
}

func SendMessageFromPb(pb *vcbpb.SendMessage) *SendMessage {
	return &SendMessage{
		Txs: pb.Txs,
	}
}

func (m *SendMessage) Pb() *vcbpb.SendMessage {
	return &vcbpb.SendMessage{
		Txs: m.Txs,
	}
}

func (*SendMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.SendMessage]()}
}

type EchoMessage struct {
	SignatureShare tctypes.SigShare
}

func EchoMessageFromPb(pb *vcbpb.EchoMessage) *EchoMessage {
	return &EchoMessage{
		SignatureShare: (tctypes.SigShare)(pb.SignatureShare),
	}
}

func (m *EchoMessage) Pb() *vcbpb.EchoMessage {
	return &vcbpb.EchoMessage{
		SignatureShare: ([]uint8)(m.SignatureShare),
	}
}

func (*EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.EchoMessage]()}
}

type FinalMessage struct {
	Txs       []*requestpb.Request
	Signature tctypes.FullSig
}

func FinalMessageFromPb(pb *vcbpb.FinalMessage) *FinalMessage {
	return &FinalMessage{
		Txs:       pb.Txs,
		Signature: (tctypes.FullSig)(pb.Signature),
	}
}

func (m *FinalMessage) Pb() *vcbpb.FinalMessage {
	return &vcbpb.FinalMessage{
		Txs:       m.Txs,
		Signature: ([]uint8)(m.Signature),
	}
}

func (*FinalMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.FinalMessage]()}
}
