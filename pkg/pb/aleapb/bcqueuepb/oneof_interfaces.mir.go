package bcqueuepb

import (
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

func (w *Event_FreeSlot) Unwrap() *FreeSlot {
	return w.FreeSlot
}

func (w *Event_SlotFreed) Unwrap() *SlotFreed {
	return w.SlotFreed
}

type FreeSlotOrigin_Type = isFreeSlotOrigin_Type

type FreeSlotOrigin_TypeWrapper[T any] interface {
	FreeSlotOrigin_Type
	Unwrap() *T
}

func (w *FreeSlotOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *FreeSlotOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}
