package ageventstypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types1 "github.com/filecoin-project/mir/codegen/model/types"
	agevents "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents"
	modringpb "github.com/filecoin-project/mir/pkg/pb/modringpb"
	types "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() agevents.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb agevents.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *agevents.Event_InputValue:
		return &Event_InputValue{InputValue: InputValueFromPb(pb.InputValue)}
	case *agevents.Event_Deliver:
		return &Event_Deliver{Deliver: DeliverFromPb(pb.Deliver)}
	case *agevents.Event_StaleMsgsRevcd:
		return &Event_StaleMsgsRevcd{StaleMsgsRevcd: StaleMsgsRecvdFromPb(pb.StaleMsgsRevcd)}
	}
	return nil
}

type Event_InputValue struct {
	InputValue *InputValue
}

func (*Event_InputValue) isEvent_Type() {}

func (w *Event_InputValue) Unwrap() *InputValue {
	return w.InputValue
}

func (w *Event_InputValue) Pb() agevents.Event_Type {
	return &agevents.Event_InputValue{InputValue: (w.InputValue).Pb()}
}

func (*Event_InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agevents.Event_InputValue]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() agevents.Event_Type {
	return &agevents.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agevents.Event_Deliver]()}
}

type Event_StaleMsgsRevcd struct {
	StaleMsgsRevcd *StaleMsgsRecvd
}

func (*Event_StaleMsgsRevcd) isEvent_Type() {}

func (w *Event_StaleMsgsRevcd) Unwrap() *StaleMsgsRecvd {
	return w.StaleMsgsRevcd
}

func (w *Event_StaleMsgsRevcd) Pb() agevents.Event_Type {
	return &agevents.Event_StaleMsgsRevcd{StaleMsgsRevcd: (w.StaleMsgsRevcd).Pb()}
}

func (*Event_StaleMsgsRevcd) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agevents.Event_StaleMsgsRevcd]()}
}

func EventFromPb(pb *agevents.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *agevents.Event {
	return &agevents.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agevents.Event]()}
}

type InputValue struct {
	Round uint64
	Input bool
}

func InputValueFromPb(pb *agevents.InputValue) *InputValue {
	return &InputValue{
		Round: pb.Round,
		Input: pb.Input,
	}
}

func (m *InputValue) Pb() *agevents.InputValue {
	return &agevents.InputValue{
		Round: m.Round,
		Input: m.Input,
	}
}

func (*InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agevents.InputValue]()}
}

type Deliver struct {
	Round    uint64
	Decision bool
}

func DeliverFromPb(pb *agevents.Deliver) *Deliver {
	return &Deliver{
		Round:    pb.Round,
		Decision: pb.Decision,
	}
}

func (m *Deliver) Pb() *agevents.Deliver {
	return &agevents.Deliver{
		Round:    m.Round,
		Decision: m.Decision,
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agevents.Deliver]()}
}

type StaleMsgsRecvd struct {
	Messages []*types.PastMessage
}

func StaleMsgsRecvdFromPb(pb *agevents.StaleMsgsRecvd) *StaleMsgsRecvd {
	return &StaleMsgsRecvd{
		Messages: types1.ConvertSlice(pb.Messages, func(t *modringpb.PastMessage) *types.PastMessage {
			return types.PastMessageFromPb(t)
		}),
	}
}

func (m *StaleMsgsRecvd) Pb() *agevents.StaleMsgsRecvd {
	return &agevents.StaleMsgsRecvd{
		Messages: types1.ConvertSlice(m.Messages, func(t *types.PastMessage) *modringpb.PastMessage {
			return (t).Pb()
		}),
	}
}

func (*StaleMsgsRecvd) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*agevents.StaleMsgsRecvd]()}
}
