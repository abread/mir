package messages

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_Ack) Unwrap() *AckMessage {
	return w.Ack
}