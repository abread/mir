package vcbpb

import (
	bcpb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
)

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_InputValue) Unwrap() *InputValue {
	return w.InputValue
}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_SendMessage) Unwrap() *SendMessage {
	return w.SendMessage
}

func (w *Message_EchoMessage) Unwrap() *EchoMessage {
	return w.EchoMessage
}

func (w *Message_FinalMessage) Unwrap() *FinalMessage {
	return w.FinalMessage
}

type Origin_Type = isOrigin_Type

type Origin_TypeWrapper[T any] interface {
	Origin_Type
	Unwrap() *T
}

func (w *Origin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *Origin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}

func (w *Origin_AleaBc) Unwrap() *bcpb.BcOrigin {
	return w.AleaBc
}
