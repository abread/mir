package reliablenetpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[Ev any] interface {
	Event_Type
	Unwrap() *Ev
}

func (p *Event_SendMessage) Unwrap() *SendMessage {
	return p.SendMessage
}

func (p *Event_Ack) Unwrap() *Ack {
	return p.Ack
}

func (p *Event_MarkRecvd) Unwrap() *MarkRecvd {
	return p.MarkRecvd
}

func (p *Event_MarkModuleMsgsRecvd) Unwrap() *MarkModuleMsgsRecvd {
	return p.MarkModuleMsgsRecvd
}

func (p *Event_RetransmitAll) Unwrap() *RetransmitAll {
	return p.RetransmitAll
}
