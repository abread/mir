// Code generated by Mir codegen. DO NOT EDIT.

package bcqueuepbtypes

import (
	"time"

	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types2 "github.com/filecoin-project/mir/codegen/model/types"
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	types3 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb"
	trantorpb "github.com/filecoin-project/mir/pkg/pb/trantorpb"
	types1 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/trantor/types"
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
	if pb == nil {
		return nil
	}
	switch pb := pb.(type) {
	case *bcqueuepb.Event_InputValue:
		return &Event_InputValue{InputValue: InputValueFromPb(pb.InputValue)}
	case *bcqueuepb.Event_Deliver:
		return &Event_Deliver{Deliver: DeliverFromPb(pb.Deliver)}
	case *bcqueuepb.Event_FreeSlot:
		return &Event_FreeSlot{FreeSlot: FreeSlotFromPb(pb.FreeSlot)}
	case *bcqueuepb.Event_PastVcbFinal:
		return &Event_PastVcbFinal{PastVcbFinal: PastVcbFinalFromPb(pb.PastVcbFinal)}
	case *bcqueuepb.Event_BcStarted:
		return &Event_BcStarted{BcStarted: BcStartedFromPb(pb.BcStarted)}
	case *bcqueuepb.Event_BcQuorumDone:
		return &Event_BcQuorumDone{BcQuorumDone: BcQuorumDoneFromPb(pb.BcQuorumDone)}
	case *bcqueuepb.Event_BcAllDone:
		return &Event_BcAllDone{BcAllDone: BcAllDoneFromPb(pb.BcAllDone)}
	case *bcqueuepb.Event_FreeStale:
		return &Event_FreeStale{FreeStale: FreeStaleFromPb(pb.FreeStale)}
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
	if w == nil {
		return nil
	}
	if w.InputValue == nil {
		return &bcqueuepb.Event_InputValue{}
	}
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
	if w == nil {
		return nil
	}
	if w.Deliver == nil {
		return &bcqueuepb.Event_Deliver{}
	}
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
	if w == nil {
		return nil
	}
	if w.FreeSlot == nil {
		return &bcqueuepb.Event_FreeSlot{}
	}
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
	if w == nil {
		return nil
	}
	if w.PastVcbFinal == nil {
		return &bcqueuepb.Event_PastVcbFinal{}
	}
	return &bcqueuepb.Event_PastVcbFinal{PastVcbFinal: (w.PastVcbFinal).Pb()}
}

func (*Event_PastVcbFinal) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event_PastVcbFinal]()}
}

type Event_BcStarted struct {
	BcStarted *BcStarted
}

func (*Event_BcStarted) isEvent_Type() {}

func (w *Event_BcStarted) Unwrap() *BcStarted {
	return w.BcStarted
}

func (w *Event_BcStarted) Pb() bcqueuepb.Event_Type {
	if w == nil {
		return nil
	}
	if w.BcStarted == nil {
		return &bcqueuepb.Event_BcStarted{}
	}
	return &bcqueuepb.Event_BcStarted{BcStarted: (w.BcStarted).Pb()}
}

func (*Event_BcStarted) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event_BcStarted]()}
}

type Event_BcQuorumDone struct {
	BcQuorumDone *BcQuorumDone
}

func (*Event_BcQuorumDone) isEvent_Type() {}

func (w *Event_BcQuorumDone) Unwrap() *BcQuorumDone {
	return w.BcQuorumDone
}

func (w *Event_BcQuorumDone) Pb() bcqueuepb.Event_Type {
	if w == nil {
		return nil
	}
	if w.BcQuorumDone == nil {
		return &bcqueuepb.Event_BcQuorumDone{}
	}
	return &bcqueuepb.Event_BcQuorumDone{BcQuorumDone: (w.BcQuorumDone).Pb()}
}

func (*Event_BcQuorumDone) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event_BcQuorumDone]()}
}

type Event_BcAllDone struct {
	BcAllDone *BcAllDone
}

func (*Event_BcAllDone) isEvent_Type() {}

func (w *Event_BcAllDone) Unwrap() *BcAllDone {
	return w.BcAllDone
}

func (w *Event_BcAllDone) Pb() bcqueuepb.Event_Type {
	if w == nil {
		return nil
	}
	if w.BcAllDone == nil {
		return &bcqueuepb.Event_BcAllDone{}
	}
	return &bcqueuepb.Event_BcAllDone{BcAllDone: (w.BcAllDone).Pb()}
}

func (*Event_BcAllDone) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event_BcAllDone]()}
}

type Event_FreeStale struct {
	FreeStale *FreeStale
}

func (*Event_FreeStale) isEvent_Type() {}

func (w *Event_FreeStale) Unwrap() *FreeStale {
	return w.FreeStale
}

func (w *Event_FreeStale) Pb() bcqueuepb.Event_Type {
	if w == nil {
		return nil
	}
	if w.FreeStale == nil {
		return &bcqueuepb.Event_FreeStale{}
	}
	return &bcqueuepb.Event_FreeStale{FreeStale: (w.FreeStale).Pb()}
}

func (*Event_FreeStale) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event_FreeStale]()}
}

func EventFromPb(pb *bcqueuepb.Event) *Event {
	if pb == nil {
		return nil
	}
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *bcqueuepb.Event {
	if m == nil {
		return nil
	}
	pbMessage := &bcqueuepb.Event{}
	{
		if m.Type != nil {
			pbMessage.Type = (m.Type).Pb()
		}
	}

	return pbMessage
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Event]()}
}

type InputValue struct {
	QueueSlot aleatypes.QueueSlot
	TxIds     []types.TxID
	Txs       []*types1.Transaction
}

func InputValueFromPb(pb *bcqueuepb.InputValue) *InputValue {
	if pb == nil {
		return nil
	}
	return &InputValue{
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
		TxIds: types2.ConvertSlice(pb.TxIds, func(t string) types.TxID {
			return (types.TxID)(t)
		}),
		Txs: types2.ConvertSlice(pb.Txs, func(t *trantorpb.Transaction) *types1.Transaction {
			return types1.TransactionFromPb(t)
		}),
	}
}

func (m *InputValue) Pb() *bcqueuepb.InputValue {
	if m == nil {
		return nil
	}
	pbMessage := &bcqueuepb.InputValue{}
	{
		pbMessage.QueueSlot = (uint64)(m.QueueSlot)
		pbMessage.TxIds = types2.ConvertSlice(m.TxIds, func(t types.TxID) string {
			return (string)(t)
		})
		pbMessage.Txs = types2.ConvertSlice(m.Txs, func(t *types1.Transaction) *trantorpb.Transaction {
			return (t).Pb()
		})
	}

	return pbMessage
}

func (*InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.InputValue]()}
}

type Deliver struct {
	Cert *types3.Cert
}

func DeliverFromPb(pb *bcqueuepb.Deliver) *Deliver {
	if pb == nil {
		return nil
	}
	return &Deliver{
		Cert: types3.CertFromPb(pb.Cert),
	}
}

func (m *Deliver) Pb() *bcqueuepb.Deliver {
	if m == nil {
		return nil
	}
	pbMessage := &bcqueuepb.Deliver{}
	{
		if m.Cert != nil {
			pbMessage.Cert = (m.Cert).Pb()
		}
	}

	return pbMessage
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.Deliver]()}
}

type FreeSlot struct {
	QueueSlot aleatypes.QueueSlot
}

func FreeSlotFromPb(pb *bcqueuepb.FreeSlot) *FreeSlot {
	if pb == nil {
		return nil
	}
	return &FreeSlot{
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
	}
}

func (m *FreeSlot) Pb() *bcqueuepb.FreeSlot {
	if m == nil {
		return nil
	}
	pbMessage := &bcqueuepb.FreeSlot{}
	{
		pbMessage.QueueSlot = (uint64)(m.QueueSlot)
	}

	return pbMessage
}

func (*FreeSlot) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.FreeSlot]()}
}

type PastVcbFinal struct {
	QueueSlot aleatypes.QueueSlot
	Txs       []*types1.Transaction
	Signature tctypes.FullSig
}

func PastVcbFinalFromPb(pb *bcqueuepb.PastVcbFinal) *PastVcbFinal {
	if pb == nil {
		return nil
	}
	return &PastVcbFinal{
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
		Txs: types2.ConvertSlice(pb.Txs, func(t *trantorpb.Transaction) *types1.Transaction {
			return types1.TransactionFromPb(t)
		}),
		Signature: (tctypes.FullSig)(pb.Signature),
	}
}

func (m *PastVcbFinal) Pb() *bcqueuepb.PastVcbFinal {
	if m == nil {
		return nil
	}
	pbMessage := &bcqueuepb.PastVcbFinal{}
	{
		pbMessage.QueueSlot = (uint64)(m.QueueSlot)
		pbMessage.Txs = types2.ConvertSlice(m.Txs, func(t *types1.Transaction) *trantorpb.Transaction {
			return (t).Pb()
		})
		pbMessage.Signature = ([]uint8)(m.Signature)
	}

	return pbMessage
}

func (*PastVcbFinal) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.PastVcbFinal]()}
}

type BcStarted struct {
	Slot *types3.Slot
}

func BcStartedFromPb(pb *bcqueuepb.BcStarted) *BcStarted {
	if pb == nil {
		return nil
	}
	return &BcStarted{
		Slot: types3.SlotFromPb(pb.Slot),
	}
}

func (m *BcStarted) Pb() *bcqueuepb.BcStarted {
	if m == nil {
		return nil
	}
	pbMessage := &bcqueuepb.BcStarted{}
	{
		if m.Slot != nil {
			pbMessage.Slot = (m.Slot).Pb()
		}
	}

	return pbMessage
}

func (*BcStarted) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.BcStarted]()}
}

type BcQuorumDone struct {
	Slot         *types3.Slot
	DeliverDelta time.Duration
}

func BcQuorumDoneFromPb(pb *bcqueuepb.BcQuorumDone) *BcQuorumDone {
	if pb == nil {
		return nil
	}
	return &BcQuorumDone{
		Slot:         types3.SlotFromPb(pb.Slot),
		DeliverDelta: (time.Duration)(pb.DeliverDelta),
	}
}

func (m *BcQuorumDone) Pb() *bcqueuepb.BcQuorumDone {
	if m == nil {
		return nil
	}
	pbMessage := &bcqueuepb.BcQuorumDone{}
	{
		if m.Slot != nil {
			pbMessage.Slot = (m.Slot).Pb()
		}
		pbMessage.DeliverDelta = (int64)(m.DeliverDelta)
	}

	return pbMessage
}

func (*BcQuorumDone) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.BcQuorumDone]()}
}

type BcAllDone struct {
	Slot            *types3.Slot
	QuorumDoneDelta time.Duration
}

func BcAllDoneFromPb(pb *bcqueuepb.BcAllDone) *BcAllDone {
	if pb == nil {
		return nil
	}
	return &BcAllDone{
		Slot:            types3.SlotFromPb(pb.Slot),
		QuorumDoneDelta: (time.Duration)(pb.QuorumDoneDelta),
	}
}

func (m *BcAllDone) Pb() *bcqueuepb.BcAllDone {
	if m == nil {
		return nil
	}
	pbMessage := &bcqueuepb.BcAllDone{}
	{
		if m.Slot != nil {
			pbMessage.Slot = (m.Slot).Pb()
		}
		pbMessage.QuorumDoneDelta = (int64)(m.QuorumDoneDelta)
	}

	return pbMessage
}

func (*BcAllDone) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.BcAllDone]()}
}

type FreeStale struct {
	QueueSlot aleatypes.QueueSlot
}

func FreeStaleFromPb(pb *bcqueuepb.FreeStale) *FreeStale {
	if pb == nil {
		return nil
	}
	return &FreeStale{
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
	}
}

func (m *FreeStale) Pb() *bcqueuepb.FreeStale {
	if m == nil {
		return nil
	}
	pbMessage := &bcqueuepb.FreeStale{}
	{
		pbMessage.QueueSlot = (uint64)(m.QueueSlot)
	}

	return pbMessage
}

func (*FreeStale) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcqueuepb.FreeStale]()}
}
