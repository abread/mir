package modringpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	modringpb "github.com/filecoin-project/mir/pkg/pb/modringpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() modringpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb modringpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *modringpb.Event_Free:
		return &Event_Free{Free: FreeSubmoduleFromPb(pb.Free)}
	}
	return nil
}

type Event_Free struct {
	Free *FreeSubmodule
}

func (*Event_Free) isEvent_Type() {}

func (w *Event_Free) Unwrap() *FreeSubmodule {
	return w.Free
}

func (w *Event_Free) Pb() modringpb.Event_Type {
	return &modringpb.Event_Free{Free: (w.Free).Pb()}
}

func (*Event_Free) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.Event_Free]()}
}

func EventFromPb(pb *modringpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *modringpb.Event {
	return &modringpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.Event]()}
}

type FreeSubmodule struct {
	Id uint64
}

func FreeSubmoduleFromPb(pb *modringpb.FreeSubmodule) *FreeSubmodule {
	return &FreeSubmodule{
		Id: pb.Id,
	}
}

func (m *FreeSubmodule) Pb() *modringpb.FreeSubmodule {
	return &modringpb.FreeSubmodule{
		Id: m.Id,
	}
}

func (*FreeSubmodule) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.FreeSubmodule]()}
}
