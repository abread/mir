package modringpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types1 "github.com/filecoin-project/mir/pkg/pb/contextstorepb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/dslpb/types"
	modringpb "github.com/filecoin-project/mir/pkg/pb/modringpb"
	types "github.com/filecoin-project/mir/pkg/types"
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
	case *modringpb.Event_Freed:
		return &Event_Freed{Freed: FreedSubmoduleFromPb(pb.Freed)}
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

type Event_Freed struct {
	Freed *FreedSubmodule
}

func (*Event_Freed) isEvent_Type() {}

func (w *Event_Freed) Unwrap() *FreedSubmodule {
	return w.Freed
}

func (w *Event_Freed) Pb() modringpb.Event_Type {
	return &modringpb.Event_Freed{Freed: (w.Freed).Pb()}
}

func (*Event_Freed) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.Event_Freed]()}
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
	Id     uint64
	Origin *FreeSubmoduleOrigin
}

func FreeSubmoduleFromPb(pb *modringpb.FreeSubmodule) *FreeSubmodule {
	return &FreeSubmodule{
		Id:     pb.Id,
		Origin: FreeSubmoduleOriginFromPb(pb.Origin),
	}
}

func (m *FreeSubmodule) Pb() *modringpb.FreeSubmodule {
	return &modringpb.FreeSubmodule{
		Id:     m.Id,
		Origin: (m.Origin).Pb(),
	}
}

func (*FreeSubmodule) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.FreeSubmodule]()}
}

type FreedSubmodule struct {
	Origin *FreeSubmoduleOrigin
}

func FreedSubmoduleFromPb(pb *modringpb.FreedSubmodule) *FreedSubmodule {
	return &FreedSubmodule{
		Origin: FreeSubmoduleOriginFromPb(pb.Origin),
	}
}

func (m *FreedSubmodule) Pb() *modringpb.FreedSubmodule {
	return &modringpb.FreedSubmodule{
		Origin: (m.Origin).Pb(),
	}
}

func (*FreedSubmodule) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.FreedSubmodule]()}
}

type FreeSubmoduleOrigin struct {
	Module types.ModuleID
	Type   FreeSubmoduleOrigin_Type
}

type FreeSubmoduleOrigin_Type interface {
	mirreflect.GeneratedType
	isFreeSubmoduleOrigin_Type()
	Pb() modringpb.FreeSubmoduleOrigin_Type
}

type FreeSubmoduleOrigin_TypeWrapper[T any] interface {
	FreeSubmoduleOrigin_Type
	Unwrap() *T
}

func FreeSubmoduleOrigin_TypeFromPb(pb modringpb.FreeSubmoduleOrigin_Type) FreeSubmoduleOrigin_Type {
	switch pb := pb.(type) {
	case *modringpb.FreeSubmoduleOrigin_ContextStore:
		return &FreeSubmoduleOrigin_ContextStore{ContextStore: types1.OriginFromPb(pb.ContextStore)}
	case *modringpb.FreeSubmoduleOrigin_Dsl:
		return &FreeSubmoduleOrigin_Dsl{Dsl: types2.OriginFromPb(pb.Dsl)}
	}
	return nil
}

type FreeSubmoduleOrigin_ContextStore struct {
	ContextStore *types1.Origin
}

func (*FreeSubmoduleOrigin_ContextStore) isFreeSubmoduleOrigin_Type() {}

func (w *FreeSubmoduleOrigin_ContextStore) Unwrap() *types1.Origin {
	return w.ContextStore
}

func (w *FreeSubmoduleOrigin_ContextStore) Pb() modringpb.FreeSubmoduleOrigin_Type {
	return &modringpb.FreeSubmoduleOrigin_ContextStore{ContextStore: (w.ContextStore).Pb()}
}

func (*FreeSubmoduleOrigin_ContextStore) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.FreeSubmoduleOrigin_ContextStore]()}
}

type FreeSubmoduleOrigin_Dsl struct {
	Dsl *types2.Origin
}

func (*FreeSubmoduleOrigin_Dsl) isFreeSubmoduleOrigin_Type() {}

func (w *FreeSubmoduleOrigin_Dsl) Unwrap() *types2.Origin {
	return w.Dsl
}

func (w *FreeSubmoduleOrigin_Dsl) Pb() modringpb.FreeSubmoduleOrigin_Type {
	return &modringpb.FreeSubmoduleOrigin_Dsl{Dsl: (w.Dsl).Pb()}
}

func (*FreeSubmoduleOrigin_Dsl) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.FreeSubmoduleOrigin_Dsl]()}
}

func FreeSubmoduleOriginFromPb(pb *modringpb.FreeSubmoduleOrigin) *FreeSubmoduleOrigin {
	return &FreeSubmoduleOrigin{
		Module: (types.ModuleID)(pb.Module),
		Type:   FreeSubmoduleOrigin_TypeFromPb(pb.Type),
	}
}

func (m *FreeSubmoduleOrigin) Pb() *modringpb.FreeSubmoduleOrigin {
	return &modringpb.FreeSubmoduleOrigin{
		Module: (string)(m.Module),
		Type:   (m.Type).Pb(),
	}
}

func (*FreeSubmoduleOrigin) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.FreeSubmoduleOrigin]()}
}
