package modringpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_Free) Unwrap() *FreeSubmodule {
	return w.Free
}
