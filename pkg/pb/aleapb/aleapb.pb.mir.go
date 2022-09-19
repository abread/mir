package aleapb

type Event_Type = isEvent_Type

type Event_TypeWrapper[Ev any] interface {
	Event_Type
	Unwrap() *Ev
}

func (p *Event_VcbcDeliver) Unwrap() *VCBCDeliver {
	return p.VcbcDeliver
}

func (p *Event_VcbcMessage) Unwrap() *VCBCMessageRecvd {
	return p.VcbcMessage
}

func (p *Event_AbaDeliver) Unwrap() *CobaltABBADeliver {
	return p.AbaDeliver
}
