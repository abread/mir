package bcpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_StartBroadcast) Unwrap() *StartBroadcast {
	return w.StartBroadcast
}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_FreeSlot) Unwrap() *FreeSlot {
	return w.FreeSlot
}

func (w *Event_FillGap) Unwrap() *DoFillGap {
	return w.FillGap
}
