package directorpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
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
	case *directorpb.Event_FillGap:
		return &Event_FillGap{FillGap: DoFillGapFromPb(pb.FillGap)}
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

type Event_FillGap struct {
	FillGap *DoFillGap
}

func (*Event_FillGap) isEvent_Type() {}

func (w *Event_FillGap) Unwrap() *DoFillGap {
	return w.FillGap
}

func (w *Event_FillGap) Pb() directorpb.Event_Type {
	return &directorpb.Event_FillGap{FillGap: (w.FillGap).Pb()}
}

func (*Event_FillGap) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*directorpb.Event_FillGap]()}
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

type DoFillGap struct {
	Slot *types.Slot
}

func DoFillGapFromPb(pb *directorpb.DoFillGap) *DoFillGap {
	return &DoFillGap{
		Slot: types.SlotFromPb(pb.Slot),
	}
}

func (m *DoFillGap) Pb() *directorpb.DoFillGap {
	return &directorpb.DoFillGap{
		Slot: (m.Slot).Pb(),
	}
}

func (*DoFillGap) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*directorpb.DoFillGap]()}
}
