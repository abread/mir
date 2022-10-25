package bcpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[Ev any] interface {
	Event_Type
	Unwrap() *Ev
}

func (p *Event_StartBroadcast) Unwrap() *StartBroadcast {
	return p.StartBroadcast
}

func (p *Event_Deliver) Unwrap() *Deliver {
	return p.Deliver
}

func (p *Event_FreeSlot) Unwrap() *FreeSlot {
	return p.FreeSlot
}
