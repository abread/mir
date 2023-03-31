package bcpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	bcpb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() bcpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb bcpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *bcpb.Event_StartBroadcast:
		return &Event_StartBroadcast{StartBroadcast: StartBroadcastFromPb(pb.StartBroadcast)}
	case *bcpb.Event_Deliver:
		return &Event_Deliver{Deliver: DeliverFromPb(pb.Deliver)}
	case *bcpb.Event_FreeSlot:
		return &Event_FreeSlot{FreeSlot: FreeSlotFromPb(pb.FreeSlot)}
	case *bcpb.Event_FillGap:
		return &Event_FillGap{FillGap: DoFillGapFromPb(pb.FillGap)}
	case *bcpb.Event_BcStarted:
		return &Event_BcStarted{BcStarted: BcStartedFromPb(pb.BcStarted)}
	}
	return nil
}

type Event_StartBroadcast struct {
	StartBroadcast *StartBroadcast
}

func (*Event_StartBroadcast) isEvent_Type() {}

func (w *Event_StartBroadcast) Unwrap() *StartBroadcast {
	return w.StartBroadcast
}

func (w *Event_StartBroadcast) Pb() bcpb.Event_Type {
	return &bcpb.Event_StartBroadcast{StartBroadcast: (w.StartBroadcast).Pb()}
}

func (*Event_StartBroadcast) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_StartBroadcast]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() bcpb.Event_Type {
	return &bcpb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_Deliver]()}
}

type Event_FreeSlot struct {
	FreeSlot *FreeSlot
}

func (*Event_FreeSlot) isEvent_Type() {}

func (w *Event_FreeSlot) Unwrap() *FreeSlot {
	return w.FreeSlot
}

func (w *Event_FreeSlot) Pb() bcpb.Event_Type {
	return &bcpb.Event_FreeSlot{FreeSlot: (w.FreeSlot).Pb()}
}

func (*Event_FreeSlot) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_FreeSlot]()}
}

type Event_FillGap struct {
	FillGap *DoFillGap
}

func (*Event_FillGap) isEvent_Type() {}

func (w *Event_FillGap) Unwrap() *DoFillGap {
	return w.FillGap
}

func (w *Event_FillGap) Pb() bcpb.Event_Type {
	return &bcpb.Event_FillGap{FillGap: (w.FillGap).Pb()}
}

func (*Event_FillGap) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_FillGap]()}
}

type Event_BcStarted struct {
	BcStarted *BcStarted
}

func (*Event_BcStarted) isEvent_Type() {}

func (w *Event_BcStarted) Unwrap() *BcStarted {
	return w.BcStarted
}

func (w *Event_BcStarted) Pb() bcpb.Event_Type {
	return &bcpb.Event_BcStarted{BcStarted: (w.BcStarted).Pb()}
}

func (*Event_BcStarted) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_BcStarted]()}
}

func EventFromPb(pb *bcpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *bcpb.Event {
	return &bcpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event]()}
}

type StartBroadcast struct {
	QueueSlot aleatypes.QueueSlot
	Txs       []*requestpb.Request
}

func StartBroadcastFromPb(pb *bcpb.StartBroadcast) *StartBroadcast {
	return &StartBroadcast{
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
		Txs:       pb.Txs,
	}
}

func (m *StartBroadcast) Pb() *bcpb.StartBroadcast {
	return &bcpb.StartBroadcast{
		QueueSlot: (uint64)(m.QueueSlot),
		Txs:       m.Txs,
	}
}

func (*StartBroadcast) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.StartBroadcast]()}
}

type Deliver struct {
	Slot *types.Slot
}

func DeliverFromPb(pb *bcpb.Deliver) *Deliver {
	return &Deliver{
		Slot: types.SlotFromPb(pb.Slot),
	}
}

func (m *Deliver) Pb() *bcpb.Deliver {
	return &bcpb.Deliver{
		Slot: (m.Slot).Pb(),
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Deliver]()}
}

type FreeSlot struct {
	Slot *types.Slot
}

func FreeSlotFromPb(pb *bcpb.FreeSlot) *FreeSlot {
	return &FreeSlot{
		Slot: types.SlotFromPb(pb.Slot),
	}
}

func (m *FreeSlot) Pb() *bcpb.FreeSlot {
	return &bcpb.FreeSlot{
		Slot: (m.Slot).Pb(),
	}
}

func (*FreeSlot) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.FreeSlot]()}
}

type DoFillGap struct {
	Slot *types.Slot
}

func DoFillGapFromPb(pb *bcpb.DoFillGap) *DoFillGap {
	return &DoFillGap{
		Slot: types.SlotFromPb(pb.Slot),
	}
}

func (m *DoFillGap) Pb() *bcpb.DoFillGap {
	return &bcpb.DoFillGap{
		Slot: (m.Slot).Pb(),
	}
}

func (*DoFillGap) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.DoFillGap]()}
}

type BcStarted struct {
	Slot *types.Slot
}

func BcStartedFromPb(pb *bcpb.BcStarted) *BcStarted {
	return &BcStarted{
		Slot: types.SlotFromPb(pb.Slot),
	}
}

func (m *BcStarted) Pb() *bcpb.BcStarted {
	return &bcpb.BcStarted{
		Slot: (m.Slot).Pb(),
	}
}

func (*BcStarted) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.BcStarted]()}
}
