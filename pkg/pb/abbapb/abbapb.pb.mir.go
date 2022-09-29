package abbapb

type Event_Type = isEvent_Type

type Event_TypeWrapper[Ev any] interface {
	Event_Type
	Unwrap() *Ev
}

func (p *Event_InputValue) Unwrap() *InputValue {
	return p.InputValue
}

func (p *Event_Deliver) Unwrap() *Deliver {
	return p.Deliver
}
