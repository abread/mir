package reliablenetpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_SendMessage) Unwrap() *SendMessage {
	return w.SendMessage
}

func (w *Event_Ack) Unwrap() *Ack {
	return w.Ack
}

func (w *Event_MarkRecvd) Unwrap() *MarkRecvd {
	return w.MarkRecvd
}

func (w *Event_MarkModuleMsgsRecvd) Unwrap() *MarkModuleMsgsRecvd {
	return w.MarkModuleMsgsRecvd
}

func (w *Event_RetransmitAll) Unwrap() *RetransmitAll {
	return w.RetransmitAll
}
