package directorpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_Heartbeat) Unwrap() *Heartbeat {
	return w.Heartbeat
}

func (w *Event_FillGap) Unwrap() *DoFillGap {
	return w.FillGap
}
