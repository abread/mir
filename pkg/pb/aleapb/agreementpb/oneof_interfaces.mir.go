package agreementpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_RequestInput) Unwrap() *RequestInput {
	return w.RequestInput
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

func (w *Message_FinishAbba) Unwrap() *FinishAbbaMessage {
	return w.FinishAbba
}
