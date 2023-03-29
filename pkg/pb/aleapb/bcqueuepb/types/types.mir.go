package bcqueuepbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	bcqueuepb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
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
	case *bcqueuepb.Event_PastVcbFinal:
		return &Event_PastVcbFinal{PastVcbFinal: PastVcbFinalFromPb(pb.PastVcbFinal)}
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

type Event_PastVcbFinal struct {
	PastVcbFinal *PastVcbFinal
}

func (*Event_PastVcbFinal) isEvent_Type() {}

func (w *Event_PastVcbFinal) Unwrap() *PastVcbFinal {
	return w.PastVcbFinal
}

func (w *Event_PastVcbFinal) Pb() bcqueuepb.Event_Type {
	return &bcqueuepb.Event_PastVcbFinal{PastVcbFinal: (w.PastVcbFinal).Pb()}
}

func (*Event_PastVcbFinal) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event_PastVcbFinal]()}
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
}

func FreeSlotFromPb(pb *bcqueuepb.FreeSlot) *FreeSlot {
	return &FreeSlot{
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
	}
}

func (m *FreeSlot) Pb() *bcqueuepb.FreeSlot {
	return &bcqueuepb.FreeSlot{
		QueueSlot: (uint64)(m.QueueSlot),
	}
}

func (*FreeSlot) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.FreeSlot]()}
}

type PastVcbFinal struct {
	QueueSlot aleatypes.QueueSlot
	Txs       []*requestpb.Request
	Signature tctypes.FullSig
}

func PastVcbFinalFromPb(pb *bcqueuepb.PastVcbFinal) *PastVcbFinal {
	return &PastVcbFinal{
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
		Txs:       pb.Txs,
		Signature: (tctypes.FullSig)(pb.Signature),
	}
}

func (m *PastVcbFinal) Pb() *bcqueuepb.PastVcbFinal {
	return &bcqueuepb.PastVcbFinal{
		QueueSlot: (uint64)(m.QueueSlot),
		Txs:       m.Txs,
		Signature: ([]uint8)(m.Signature),
	}
}

func (*PastVcbFinal) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.PastVcbFinal]()}
}