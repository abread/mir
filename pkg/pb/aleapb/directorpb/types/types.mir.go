package directorpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	directorpb "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() directorpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb directorpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *directorpb.Event_Heartbeat:
		return &Event_Heartbeat{Heartbeat: HeartbeatFromPb(pb.Heartbeat)}
	}
	return nil
}

type Event_Heartbeat struct {
	Heartbeat *Heartbeat
}

func (*Event_Heartbeat) isEvent_Type() {}

func (w *Event_Heartbeat) Unwrap() *Heartbeat {
	return w.Heartbeat
}

func (w *Event_Heartbeat) Pb() directorpb.Event_Type {
	return &directorpb.Event_Heartbeat{Heartbeat: (w.Heartbeat).Pb()}
}

func (*Event_Heartbeat) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*directorpb.Event_Heartbeat]()}
}

func EventFromPb(pb *directorpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *directorpb.Event {
	return &directorpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*directorpb.Event]()}
}

type Heartbeat struct{}

func HeartbeatFromPb(pb *directorpb.Heartbeat) *Heartbeat {
	return &Heartbeat{}
}

func (m *Heartbeat) Pb() *directorpb.Heartbeat {
	return &directorpb.Heartbeat{}
}

func (*Heartbeat) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*directorpb.Heartbeat]()}
}