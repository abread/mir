package agreementpb

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_FinishAbba) Unwrap() *FinishAbbaMessage {
	return w.FinishAbba
}