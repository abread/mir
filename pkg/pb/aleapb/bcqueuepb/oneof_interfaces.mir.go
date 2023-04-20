package bcqueuepb

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

func (w *Event_PastVcbFinal) Unwrap() *PastVcbFinal {
	return w.PastVcbFinal
}

func (w *Event_BcStarted) Unwrap() *BcStarted {
	return w.BcStarted
}