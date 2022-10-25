package aleapb

import (
	agreementpb "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	bcpb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
)

type Event_Type = isEvent_Type

type Event_TypeWrapper[Ev any] interface {
	Event_Type
	Unwrap() *Ev
}

func (p *Event_Broadcast) Unwrap() *bcpb.Event {
	return p.Broadcast
}

func (p *Event_Agreement) Unwrap() *agreementpb.Event {
	return p.Agreement
}
