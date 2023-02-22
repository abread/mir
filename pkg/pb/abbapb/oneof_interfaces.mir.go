package abbapb

import (
	agreementpb "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
)

type RoundEvent_Type = isRoundEvent_Type

type RoundEvent_TypeWrapper[T any] interface {
	RoundEvent_Type
	Unwrap() *T
}

func (w *RoundEvent_InputValue) Unwrap() *RoundInputValue {
	return w.InputValue
}

func (w *RoundEvent_Deliver) Unwrap() *RoundDeliver {
	return w.Deliver
}

func (w *RoundEvent_Finish) Unwrap() *RoundFinishAll {
	return w.Finish
}

type RoundMessage_Type = isRoundMessage_Type

type RoundMessage_TypeWrapper[T any] interface {
	RoundMessage_Type
	Unwrap() *T
}

func (w *RoundMessage_Init) Unwrap() *RoundInitMessage {
	return w.Init
}

func (w *RoundMessage_Aux) Unwrap() *RoundAuxMessage {
	return w.Aux
}

func (w *RoundMessage_Conf) Unwrap() *RoundConfMessage {
	return w.Conf
}

func (w *RoundMessage_Coin) Unwrap() *RoundCoinMessage {
	return w.Coin
}

type RoundOrigin_Type = isRoundOrigin_Type

type RoundOrigin_TypeWrapper[T any] interface {
	RoundOrigin_Type
	Unwrap() *T
}

func (w *RoundOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *RoundOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}

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

func (w *Event_Round) Unwrap() *RoundEvent {
	return w.Round
}

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_Finish) Unwrap() *FinishMessage {
	return w.Finish
}

func (w *Message_Round) Unwrap() *RoundMessage {
	return w.Round
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

func (w *Origin_AleaAg) Unwrap() *agreementpb.AbbaOrigin {
	return w.AleaAg
}
