package aleapb

type Event_Type = isEvent_Type

type Event_TypeWrapper[Ev any] interface {
	Event_Type
	Unwrap() *Ev
}

func (p *Event_AbaDeliver) Unwrap() *CobaltABBADeliver {
	return p.AbaDeliver
}

func (p *Event_AbaMessage) Unwrap() *CobaltABBAMessageRecvd {
	return p.AbaMessage
}
