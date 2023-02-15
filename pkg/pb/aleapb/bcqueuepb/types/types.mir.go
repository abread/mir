package bcqueuepbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	bcqueuepb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/contextstorepb/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/dslpb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types1 "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() bcqueuepb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb bcqueuepb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *bcqueuepb.Event_InputValue:
		return &Event_InputValue{InputValue: InputValueFromPb(pb.InputValue)}
	case *bcqueuepb.Event_Deliver:
		return &Event_Deliver{Deliver: DeliverFromPb(pb.Deliver)}
	case *bcqueuepb.Event_FreeSlot:
		return &Event_FreeSlot{FreeSlot: FreeSlotFromPb(pb.FreeSlot)}
	case *bcqueuepb.Event_SlotFreed:
		return &Event_SlotFreed{SlotFreed: SlotFreedFromPb(pb.SlotFreed)}
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

func (w *Event_InputValue) Pb() bcqueuepb.Event_Type {
	return &bcqueuepb.Event_InputValue{InputValue: (w.InputValue).Pb()}
}

func (*Event_InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event_InputValue]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() bcqueuepb.Event_Type {
	return &bcqueuepb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event_Deliver]()}
}

type Event_FreeSlot struct {
	FreeSlot *FreeSlot
}

func (*Event_FreeSlot) isEvent_Type() {}

func (w *Event_FreeSlot) Unwrap() *FreeSlot {
	return w.FreeSlot
}

func (w *Event_FreeSlot) Pb() bcqueuepb.Event_Type {
	return &bcqueuepb.Event_FreeSlot{FreeSlot: (w.FreeSlot).Pb()}
}

func (*Event_FreeSlot) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event_FreeSlot]()}
}

type Event_SlotFreed struct {
	SlotFreed *SlotFreed
}

func (*Event_SlotFreed) isEvent_Type() {}

func (w *Event_SlotFreed) Unwrap() *SlotFreed {
	return w.SlotFreed
}

func (w *Event_SlotFreed) Pb() bcqueuepb.Event_Type {
	return &bcqueuepb.Event_SlotFreed{SlotFreed: (w.SlotFreed).Pb()}
}

func (*Event_SlotFreed) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event_SlotFreed]()}
}

func EventFromPb(pb *bcqueuepb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *bcqueuepb.Event {
	return &bcqueuepb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event]()}
}

type InputValue struct {
	Slot *types.Slot
	Txs  []*requestpb.Request
}

func InputValueFromPb(pb *bcqueuepb.InputValue) *InputValue {
	return &InputValue{
		Slot: types.SlotFromPb(pb.Slot),
		Txs:  pb.Txs,
	}
}

func (m *InputValue) Pb() *bcqueuepb.InputValue {
	return &bcqueuepb.InputValue{
		Slot: (m.Slot).Pb(),
		Txs:  m.Txs,
	}
}

func (*InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.InputValue]()}
}

type Deliver struct {
	Slot *types.Slot
}

func DeliverFromPb(pb *bcqueuepb.Deliver) *Deliver {
	return &Deliver{
		Slot: types.SlotFromPb(pb.Slot),
	}
}

func (m *Deliver) Pb() *bcqueuepb.Deliver {
	return &bcqueuepb.Deliver{
		Slot: (m.Slot).Pb(),
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Deliver]()}
}

type FreeSlot struct {
	QueueSlot aleatypes.QueueSlot
	Origin    *FreeSlotOrigin
}

func FreeSlotFromPb(pb *bcqueuepb.FreeSlot) *FreeSlot {
	return &FreeSlot{
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
		Origin:    FreeSlotOriginFromPb(pb.Origin),
	}
}

func (m *FreeSlot) Pb() *bcqueuepb.FreeSlot {
	return &bcqueuepb.FreeSlot{
		QueueSlot: (uint64)(m.QueueSlot),
		Origin:    (m.Origin).Pb(),
	}
}

func (*FreeSlot) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.FreeSlot]()}
}

type SlotFreed struct {
	Origin *FreeSlotOrigin
}

func SlotFreedFromPb(pb *bcqueuepb.SlotFreed) *SlotFreed {
	return &SlotFreed{
		Origin: FreeSlotOriginFromPb(pb.Origin),
	}
}

func (m *SlotFreed) Pb() *bcqueuepb.SlotFreed {
	return &bcqueuepb.SlotFreed{
		Origin: (m.Origin).Pb(),
	}
}

func (*SlotFreed) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.SlotFreed]()}
}

type FreeSlotOrigin struct {
	Module types1.ModuleID
	Type   FreeSlotOrigin_Type
}

type FreeSlotOrigin_Type interface {
	mirreflect.GeneratedType
	isFreeSlotOrigin_Type()
	Pb() bcqueuepb.FreeSlotOrigin_Type
}

type FreeSlotOrigin_TypeWrapper[T any] interface {
	FreeSlotOrigin_Type
	Unwrap() *T
}

func FreeSlotOrigin_TypeFromPb(pb bcqueuepb.FreeSlotOrigin_Type) FreeSlotOrigin_Type {
	switch pb := pb.(type) {
	case *bcqueuepb.FreeSlotOrigin_ContextStore:
		return &FreeSlotOrigin_ContextStore{ContextStore: types2.OriginFromPb(pb.ContextStore)}
	case *bcqueuepb.FreeSlotOrigin_Dsl:
		return &FreeSlotOrigin_Dsl{Dsl: types3.OriginFromPb(pb.Dsl)}
	}
	return nil
}

type FreeSlotOrigin_ContextStore struct {
	ContextStore *types2.Origin
}

func (*FreeSlotOrigin_ContextStore) isFreeSlotOrigin_Type() {}

func (w *FreeSlotOrigin_ContextStore) Unwrap() *types2.Origin {
	return w.ContextStore
}

func (w *FreeSlotOrigin_ContextStore) Pb() bcqueuepb.FreeSlotOrigin_Type {
	return &bcqueuepb.FreeSlotOrigin_ContextStore{ContextStore: (w.ContextStore).Pb()}
}

func (*FreeSlotOrigin_ContextStore) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.FreeSlotOrigin_ContextStore]()}
}

type FreeSlotOrigin_Dsl struct {
	Dsl *types3.Origin
}

func (*FreeSlotOrigin_Dsl) isFreeSlotOrigin_Type() {}

func (w *FreeSlotOrigin_Dsl) Unwrap() *types3.Origin {
	return w.Dsl
}

func (w *FreeSlotOrigin_Dsl) Pb() bcqueuepb.FreeSlotOrigin_Type {
	return &bcqueuepb.FreeSlotOrigin_Dsl{Dsl: (w.Dsl).Pb()}
}

func (*FreeSlotOrigin_Dsl) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.FreeSlotOrigin_Dsl]()}
}

func FreeSlotOriginFromPb(pb *bcqueuepb.FreeSlotOrigin) *FreeSlotOrigin {
	return &FreeSlotOrigin{
		Module: (types1.ModuleID)(pb.Module),
		Type:   FreeSlotOrigin_TypeFromPb(pb.Type),
	}
}

func (m *FreeSlotOrigin) Pb() *bcqueuepb.FreeSlotOrigin {
	return &bcqueuepb.FreeSlotOrigin{
		Module: (string)(m.Module),
		Type:   (m.Type).Pb(),
	}
}

func (*FreeSlotOrigin) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.FreeSlotOrigin]()}
}

type VcbOrigin struct {
	QueueSlot aleatypes.QueueSlot
}

func VcbOriginFromPb(pb *bcqueuepb.VcbOrigin) *VcbOrigin {
	return &VcbOrigin{
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
	}
}

func (m *VcbOrigin) Pb() *bcqueuepb.VcbOrigin {
	return &bcqueuepb.VcbOrigin{
		QueueSlot: (uint64)(m.QueueSlot),
	}
}

func (*VcbOrigin) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.VcbOrigin]()}
}
