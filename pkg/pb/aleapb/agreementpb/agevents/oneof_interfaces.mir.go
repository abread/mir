// Code generated by Mir codegen. DO NOT EDIT.

package agevents

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

func (w *Event_StaleMsgsRevcd) Unwrap() *StaleMsgsRecvd {
	return w.StaleMsgsRevcd
}

func (w *Event_InnerAbbaRoundTime) Unwrap() *InnerAbbaRoundTime {
	return w.InnerAbbaRoundTime
}
