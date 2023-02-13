package agreementpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	agreementpb "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() agreementpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb agreementpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *agreementpb.Event_RequestInput:
		return &Event_RequestInput{RequestInput: RequestInputFromPb(pb.RequestInput)}
	case *agreementpb.Event_InputValue:
		return &Event_InputValue{InputValue: InputValueFromPb(pb.InputValue)}
	case *agreementpb.Event_Deliver:
		return &Event_Deliver{Deliver: DeliverFromPb(pb.Deliver)}
	}
	return nil
}

type Event_RequestInput struct {
	RequestInput *RequestInput
}

func (*Event_RequestInput) isEvent_Type() {}

func (w *Event_RequestInput) Unwrap() *RequestInput {
	return w.RequestInput
}

func (w *Event_RequestInput) Pb() agreementpb.Event_Type {
	return &agreementpb.Event_RequestInput{RequestInput: (w.RequestInput).Pb()}
}

func (*Event_RequestInput) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agreementpb.Event_RequestInput]()}
}

type Event_InputValue struct {
	InputValue *InputValue
}

func (*Event_InputValue) isEvent_Type() {}

func (w *Event_InputValue) Unwrap() *InputValue {
	return w.InputValue
}

func (w *Event_InputValue) Pb() agreementpb.Event_Type {
	return &agreementpb.Event_InputValue{InputValue: (w.InputValue).Pb()}
}

func (*Event_InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agreementpb.Event_InputValue]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() agreementpb.Event_Type {
	return &agreementpb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agreementpb.Event_Deliver]()}
}

func EventFromPb(pb *agreementpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *agreementpb.Event {
	return &agreementpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agreementpb.Event]()}
}

type RequestInput struct {
	Round uint64
}

func RequestInputFromPb(pb *agreementpb.RequestInput) *RequestInput {
	return &RequestInput{
		Round: pb.Round,
	}
}

func (m *RequestInput) Pb() *agreementpb.RequestInput {
	return &agreementpb.RequestInput{
		Round: m.Round,
	}
}

func (*RequestInput) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agreementpb.RequestInput]()}
}

type InputValue struct {
	Round uint64
	Input bool
}

func InputValueFromPb(pb *agreementpb.InputValue) *InputValue {
	return &InputValue{
		Round: pb.Round,
		Input: pb.Input,
	}
}

func (m *InputValue) Pb() *agreementpb.InputValue {
	return &agreementpb.InputValue{
		Round: m.Round,
		Input: m.Input,
	}
}

func (*InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agreementpb.InputValue]()}
}

type Deliver struct {
	Round    uint64
	Decision bool
}

func DeliverFromPb(pb *agreementpb.Deliver) *Deliver {
	return &Deliver{
		Round:    pb.Round,
		Decision: pb.Decision,
	}
}

func (m *Deliver) Pb() *agreementpb.Deliver {
	return &agreementpb.Deliver{
		Round:    m.Round,
		Decision: m.Decision,
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agreementpb.Deliver]()}
}

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
