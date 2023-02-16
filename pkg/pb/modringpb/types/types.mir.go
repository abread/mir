package modringpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types "github.com/filecoin-project/mir/codegen/model/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/contextstorepb/types"
	types4 "github.com/filecoin-project/mir/pkg/pb/dslpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpb "github.com/filecoin-project/mir/pkg/pb/modringpb"
	types1 "github.com/filecoin-project/mir/pkg/types"
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
	case *modringpb.Event_PastMessagesRecvd:
		return &Event_PastMessagesRecvd{PastMessagesRecvd: PastMessagesRecvdFromPb(pb.PastMessagesRecvd)}
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

type Event_PastMessagesRecvd struct {
	PastMessagesRecvd *PastMessagesRecvd
}

func (*Event_PastMessagesRecvd) isEvent_Type() {}

func (w *Event_PastMessagesRecvd) Unwrap() *PastMessagesRecvd {
	return w.PastMessagesRecvd
}

func (w *Event_PastMessagesRecvd) Pb() modringpb.Event_Type {
	return &modringpb.Event_PastMessagesRecvd{PastMessagesRecvd: (w.PastMessagesRecvd).Pb()}
}

func (*Event_PastMessagesRecvd) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.Event_PastMessagesRecvd]()}
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

type PastMessagesRecvd struct {
	Messages []*PastMessage
}

func PastMessagesRecvdFromPb(pb *modringpb.PastMessagesRecvd) *PastMessagesRecvd {
	return &PastMessagesRecvd{
		Messages: types.ConvertSlice(pb.Messages, func(t *modringpb.PastMessage) *PastMessage {
			return PastMessageFromPb(t)
		}),
	}
}

func (m *PastMessagesRecvd) Pb() *modringpb.PastMessagesRecvd {
	return &modringpb.PastMessagesRecvd{
		Messages: types.ConvertSlice(m.Messages, func(t *PastMessage) *modringpb.PastMessage {
			return (t).Pb()
		}),
	}
}

func (*PastMessagesRecvd) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.PastMessagesRecvd]()}
}

type PastMessage struct {
	DestId  uint64
	From    types1.NodeID
	Message *types2.Message
}

func PastMessageFromPb(pb *modringpb.PastMessage) *PastMessage {
	return &PastMessage{
		DestId:  pb.DestId,
		From:    (types1.NodeID)(pb.From),
		Message: types2.MessageFromPb(pb.Message),
	}
}

func (m *PastMessage) Pb() *modringpb.PastMessage {
	return &modringpb.PastMessage{
		DestId:  m.DestId,
		From:    (string)(m.From),
		Message: (m.Message).Pb(),
	}
}

func (*PastMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.PastMessage]()}
}

type FreeSubmoduleOrigin struct {
	Module types1.ModuleID
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
		return &FreeSubmoduleOrigin_ContextStore{ContextStore: types3.OriginFromPb(pb.ContextStore)}
	case *modringpb.FreeSubmoduleOrigin_Dsl:
		return &FreeSubmoduleOrigin_Dsl{Dsl: types4.OriginFromPb(pb.Dsl)}
	}
	return nil
}

type FreeSubmoduleOrigin_ContextStore struct {
	ContextStore *types3.Origin
}

func (*FreeSubmoduleOrigin_ContextStore) isFreeSubmoduleOrigin_Type() {}

func (w *FreeSubmoduleOrigin_ContextStore) Unwrap() *types3.Origin {
	return w.ContextStore
}

func (w *FreeSubmoduleOrigin_ContextStore) Pb() modringpb.FreeSubmoduleOrigin_Type {
	return &modringpb.FreeSubmoduleOrigin_ContextStore{ContextStore: (w.ContextStore).Pb()}
}

func (*FreeSubmoduleOrigin_ContextStore) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.FreeSubmoduleOrigin_ContextStore]()}
}

type FreeSubmoduleOrigin_Dsl struct {
	Dsl *types4.Origin
}

func (*FreeSubmoduleOrigin_Dsl) isFreeSubmoduleOrigin_Type() {}

func (w *FreeSubmoduleOrigin_Dsl) Unwrap() *types4.Origin {
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
		Module: (types1.ModuleID)(pb.Module),
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
