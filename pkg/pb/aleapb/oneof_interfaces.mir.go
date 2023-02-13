package aleapb

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_FillGapMessage) Unwrap() *FillGapMessage {
	return w.FillGapMessage
}

func (w *Message_FillerMessage) Unwrap() *FillerMessage {
	return w.FillerMessage
}
