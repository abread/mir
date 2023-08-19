// Code generated by Mir codegen. DO NOT EDIT.

package bcpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types1 "github.com/filecoin-project/mir/codegen/model/types"
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	bcpb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	trantorpb "github.com/filecoin-project/mir/pkg/pb/trantorpb"
	types "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
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
	if pb == nil {
		return nil
	}
	switch pb := pb.(type) {
	case *bcpb.Event_FillGap:
		return &Event_FillGap{FillGap: DoFillGapFromPb(pb.FillGap)}
	}
	return nil
}

type Event_FillGap struct {
	FillGap *DoFillGap
}

func (*Event_FillGap) isEvent_Type() {}

func (w *Event_FillGap) Unwrap() *DoFillGap {
	return w.FillGap
}

func (w *Event_FillGap) Pb() bcpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.FillGap == nil {
		return &bcpb.Event_FillGap{}
	}
	return &bcpb.Event_FillGap{FillGap: (w.FillGap).Pb()}
}

func (*Event_FillGap) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_FillGap]()}
}

func EventFromPb(pb *bcpb.Event) *Event {
	if pb == nil {
		return nil
	}
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *bcpb.Event {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.Event{}
	{
		if m.Type != nil {
			pbMessage.Type = (m.Type).Pb()
		}
	}

	return pbMessage
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event]()}
}

type DoFillGap struct {
	Slot *Slot
}

func DoFillGapFromPb(pb *bcpb.DoFillGap) *DoFillGap {
	if pb == nil {
		return nil
	}
	return &DoFillGap{
		Slot: SlotFromPb(pb.Slot),
	}
}

func (m *DoFillGap) Pb() *bcpb.DoFillGap {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.DoFillGap{}
	{
		if m.Slot != nil {
			pbMessage.Slot = (m.Slot).Pb()
		}
	}

	return pbMessage
}

func (*DoFillGap) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.DoFillGap]()}
}

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() bcpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb bcpb.Message_Type) Message_Type {
	if pb == nil {
		return nil
	}
	switch pb := pb.(type) {
	case *bcpb.Message_FillGapMessage:
		return &Message_FillGapMessage{FillGapMessage: FillGapMessageFromPb(pb.FillGapMessage)}
	case *bcpb.Message_FillerMessage:
		return &Message_FillerMessage{FillerMessage: FillerMessageFromPb(pb.FillerMessage)}
	}
	return nil
}

type Message_FillGapMessage struct {
	FillGapMessage *FillGapMessage
}

func (*Message_FillGapMessage) isMessage_Type() {}

func (w *Message_FillGapMessage) Unwrap() *FillGapMessage {
	return w.FillGapMessage
}

func (w *Message_FillGapMessage) Pb() bcpb.Message_Type {
	if w == nil {
		return nil
	}
	if w.FillGapMessage == nil {
		return &bcpb.Message_FillGapMessage{}
	}
	return &bcpb.Message_FillGapMessage{FillGapMessage: (w.FillGapMessage).Pb()}
}

func (*Message_FillGapMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Message_FillGapMessage]()}
}

type Message_FillerMessage struct {
	FillerMessage *FillerMessage
}

func (*Message_FillerMessage) isMessage_Type() {}

func (w *Message_FillerMessage) Unwrap() *FillerMessage {
	return w.FillerMessage
}

func (w *Message_FillerMessage) Pb() bcpb.Message_Type {
	if w == nil {
		return nil
	}
	if w.FillerMessage == nil {
		return &bcpb.Message_FillerMessage{}
	}
	return &bcpb.Message_FillerMessage{FillerMessage: (w.FillerMessage).Pb()}
}

func (*Message_FillerMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Message_FillerMessage]()}
}

func MessageFromPb(pb *bcpb.Message) *Message {
	if pb == nil {
		return nil
	}
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *bcpb.Message {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.Message{}
	{
		if m.Type != nil {
			pbMessage.Type = (m.Type).Pb()
		}
	}

	return pbMessage
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Message]()}
}

type FillGapMessage struct {
	Slot *Slot
}

func FillGapMessageFromPb(pb *bcpb.FillGapMessage) *FillGapMessage {
	if pb == nil {
		return nil
	}
	return &FillGapMessage{
		Slot: SlotFromPb(pb.Slot),
	}
}

func (m *FillGapMessage) Pb() *bcpb.FillGapMessage {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.FillGapMessage{}
	{
		if m.Slot != nil {
			pbMessage.Slot = (m.Slot).Pb()
		}
	}

	return pbMessage
}

func (*FillGapMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.FillGapMessage]()}
}

type FillerMessage struct {
	Slot      *Slot
	Txs       []*types.Transaction
	Signature tctypes.FullSig
}

func FillerMessageFromPb(pb *bcpb.FillerMessage) *FillerMessage {
	if pb == nil {
		return nil
	}
	return &FillerMessage{
		Slot: SlotFromPb(pb.Slot),
		Txs: types1.ConvertSlice(pb.Txs, func(t *trantorpb.Transaction) *types.Transaction {
			return types.TransactionFromPb(t)
		}),
		Signature: (tctypes.FullSig)(pb.Signature),
	}
}

func (m *FillerMessage) Pb() *bcpb.FillerMessage {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.FillerMessage{}
	{
		if m.Slot != nil {
			pbMessage.Slot = (m.Slot).Pb()
		}
		pbMessage.Txs = types1.ConvertSlice(m.Txs, func(t *types.Transaction) *trantorpb.Transaction {
			return (t).Pb()
		})
		pbMessage.Signature = ([]uint8)(m.Signature)
	}

	return pbMessage
}

func (*FillerMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.FillerMessage]()}
}

type Slot struct {
	QueueIdx  aleatypes.QueueIdx
	QueueSlot aleatypes.QueueSlot
}

func SlotFromPb(pb *bcpb.Slot) *Slot {
	if pb == nil {
		return nil
	}
	return &Slot{
		QueueIdx:  (aleatypes.QueueIdx)(pb.QueueIdx),
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
	}
}

func (m *Slot) Pb() *bcpb.Slot {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.Slot{}
	{
		pbMessage.QueueIdx = (uint32)(m.QueueIdx)
		pbMessage.QueueSlot = (uint64)(m.QueueSlot)
	}

	return pbMessage
}

func (*Slot) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Slot]()}
}

type Cert struct {
	Slot        *Slot
	BatchDigest []uint8
	Signature   tctypes.FullSig
}

func CertFromPb(pb *bcpb.Cert) *Cert {
	if pb == nil {
		return nil
	}
	return &Cert{
		Slot:        SlotFromPb(pb.Slot),
		BatchDigest: pb.BatchDigest,
		Signature:   (tctypes.FullSig)(pb.Signature),
	}
}

func (m *Cert) Pb() *bcpb.Cert {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.Cert{}
	{
		if m.Slot != nil {
			pbMessage.Slot = (m.Slot).Pb()
		}
		pbMessage.BatchDigest = m.BatchDigest
		pbMessage.Signature = ([]uint8)(m.Signature)
	}

	return pbMessage
}

func (*Cert) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Cert]()}
}
