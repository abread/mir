package abbapb

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

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_FinishMessage) Unwrap() *FinishMessage {
	return w.FinishMessage
}

func (w *Message_InitMessage) Unwrap() *InitMessage {
	return w.InitMessage
}

func (w *Message_AuxMessage) Unwrap() *AuxMessage {
	return w.AuxMessage
}

func (w *Message_ConfMessage) Unwrap() *ConfMessage {
	return w.ConfMessage
}

func (w *Message_CoinMessage) Unwrap() *CoinMessage {
	return w.CoinMessage
}
