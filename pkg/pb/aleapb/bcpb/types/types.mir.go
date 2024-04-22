// Code generated by Mir codegen. DO NOT EDIT.

package bcpbtypes

import (
	"time"

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
	case *bcpb.Event_RequestCert:
		return &Event_RequestCert{RequestCert: RequestCertFromPb(pb.RequestCert)}
	case *bcpb.Event_DeliverCert:
		return &Event_DeliverCert{DeliverCert: DeliverCertFromPb(pb.DeliverCert)}
	case *bcpb.Event_BcStarted:
		return &Event_BcStarted{BcStarted: BcStartedFromPb(pb.BcStarted)}
	case *bcpb.Event_FreeSlot:
		return &Event_FreeSlot{FreeSlot: FreeSlotFromPb(pb.FreeSlot)}
	case *bcpb.Event_EstimateUpdate:
		return &Event_EstimateUpdate{EstimateUpdate: EstimateUpdateFromPb(pb.EstimateUpdate)}
	case *bcpb.Event_FillGap:
		return &Event_FillGap{FillGap: DoFillGapFromPb(pb.FillGap)}
	}
	return nil
}

type Event_RequestCert struct {
	RequestCert *RequestCert
}

func (*Event_RequestCert) isEvent_Type() {}

func (w *Event_RequestCert) Unwrap() *RequestCert {
	return w.RequestCert
}

func (w *Event_RequestCert) Pb() bcpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.RequestCert == nil {
		return &bcpb.Event_RequestCert{}
	}
	return &bcpb.Event_RequestCert{RequestCert: (w.RequestCert).Pb()}
}

func (*Event_RequestCert) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_RequestCert]()}
}

type Event_DeliverCert struct {
	DeliverCert *DeliverCert
}

func (*Event_DeliverCert) isEvent_Type() {}

func (w *Event_DeliverCert) Unwrap() *DeliverCert {
	return w.DeliverCert
}

func (w *Event_DeliverCert) Pb() bcpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.DeliverCert == nil {
		return &bcpb.Event_DeliverCert{}
	}
	return &bcpb.Event_DeliverCert{DeliverCert: (w.DeliverCert).Pb()}
}

func (*Event_DeliverCert) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_DeliverCert]()}
}

type Event_BcStarted struct {
	BcStarted *BcStarted
}

func (*Event_BcStarted) isEvent_Type() {}

func (w *Event_BcStarted) Unwrap() *BcStarted {
	return w.BcStarted
}

func (w *Event_BcStarted) Pb() bcpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.BcStarted == nil {
		return &bcpb.Event_BcStarted{}
	}
	return &bcpb.Event_BcStarted{BcStarted: (w.BcStarted).Pb()}
}

func (*Event_BcStarted) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_BcStarted]()}
}

type Event_FreeSlot struct {
	FreeSlot *FreeSlot
}

func (*Event_FreeSlot) isEvent_Type() {}

func (w *Event_FreeSlot) Unwrap() *FreeSlot {
	return w.FreeSlot
}

func (w *Event_FreeSlot) Pb() bcpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.FreeSlot == nil {
		return &bcpb.Event_FreeSlot{}
	}
	return &bcpb.Event_FreeSlot{FreeSlot: (w.FreeSlot).Pb()}
}

func (*Event_FreeSlot) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_FreeSlot]()}
}

type Event_EstimateUpdate struct {
	EstimateUpdate *EstimateUpdate
}

func (*Event_EstimateUpdate) isEvent_Type() {}

func (w *Event_EstimateUpdate) Unwrap() *EstimateUpdate {
	return w.EstimateUpdate
}

func (w *Event_EstimateUpdate) Pb() bcpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.EstimateUpdate == nil {
		return &bcpb.Event_EstimateUpdate{}
	}
	return &bcpb.Event_EstimateUpdate{EstimateUpdate: (w.EstimateUpdate).Pb()}
}

func (*Event_EstimateUpdate) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Event_EstimateUpdate]()}
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

type RequestCert struct{}

func RequestCertFromPb(pb *bcpb.RequestCert) *RequestCert {
	if pb == nil {
		return nil
	}
	return &RequestCert{}
}

func (m *RequestCert) Pb() *bcpb.RequestCert {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.RequestCert{}
	{
	}

	return pbMessage
}

func (*RequestCert) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.RequestCert]()}
}

type DeliverCert struct {
	Cert *Cert
}

func DeliverCertFromPb(pb *bcpb.DeliverCert) *DeliverCert {
	if pb == nil {
		return nil
	}
	return &DeliverCert{
		Cert: CertFromPb(pb.Cert),
	}
}

func (m *DeliverCert) Pb() *bcpb.DeliverCert {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.DeliverCert{}
	{
		if m.Cert != nil {
			pbMessage.Cert = (m.Cert).Pb()
		}
	}

	return pbMessage
}

func (*DeliverCert) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.DeliverCert]()}
}

type BcStarted struct {
	Slot *Slot
}

func BcStartedFromPb(pb *bcpb.BcStarted) *BcStarted {
	if pb == nil {
		return nil
	}
	return &BcStarted{
		Slot: SlotFromPb(pb.Slot),
	}
}

func (m *BcStarted) Pb() *bcpb.BcStarted {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.BcStarted{}
	{
		if m.Slot != nil {
			pbMessage.Slot = (m.Slot).Pb()
		}
	}

	return pbMessage
}

func (*BcStarted) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.BcStarted]()}
}

type FreeSlot struct {
	Slot *Slot
}

func FreeSlotFromPb(pb *bcpb.FreeSlot) *FreeSlot {
	if pb == nil {
		return nil
	}
	return &FreeSlot{
		Slot: SlotFromPb(pb.Slot),
	}
}

func (m *FreeSlot) Pb() *bcpb.FreeSlot {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.FreeSlot{}
	{
		if m.Slot != nil {
			pbMessage.Slot = (m.Slot).Pb()
		}
	}

	return pbMessage
}

func (*FreeSlot) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.FreeSlot]()}
}

type EstimateUpdate struct {
	MaxOwnBcDuration      time.Duration
	MaxOwnBcLocalDuration time.Duration
	MaxExtBcDuration      time.Duration
	MinNetLatency         time.Duration
}

func EstimateUpdateFromPb(pb *bcpb.EstimateUpdate) *EstimateUpdate {
	if pb == nil {
		return nil
	}
	return &EstimateUpdate{
		MaxOwnBcDuration:      (time.Duration)(pb.MaxOwnBcDuration),
		MaxOwnBcLocalDuration: (time.Duration)(pb.MaxOwnBcLocalDuration),
		MaxExtBcDuration:      (time.Duration)(pb.MaxExtBcDuration),
		MinNetLatency:         (time.Duration)(pb.MinNetLatency),
	}
}

func (m *EstimateUpdate) Pb() *bcpb.EstimateUpdate {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.EstimateUpdate{}
	{
		pbMessage.MaxOwnBcDuration = (int64)(m.MaxOwnBcDuration)
		pbMessage.MaxOwnBcLocalDuration = (int64)(m.MaxOwnBcLocalDuration)
		pbMessage.MaxExtBcDuration = (int64)(m.MaxExtBcDuration)
		pbMessage.MinNetLatency = (int64)(m.MinNetLatency)
	}

	return pbMessage
}

func (*EstimateUpdate) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.EstimateUpdate]()}
}

type DoFillGap struct {
	Slot        *Slot
	NextReplica uint32
}

func DoFillGapFromPb(pb *bcpb.DoFillGap) *DoFillGap {
	if pb == nil {
		return nil
	}
	return &DoFillGap{
		Slot:        SlotFromPb(pb.Slot),
		NextReplica: pb.NextReplica,
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
		pbMessage.NextReplica = m.NextReplica
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
	Cert *Cert
	Txs  []*types.Transaction
}

func FillerMessageFromPb(pb *bcpb.FillerMessage) *FillerMessage {
	if pb == nil {
		return nil
	}
	return &FillerMessage{
		Cert: CertFromPb(pb.Cert),
		Txs: types1.ConvertSlice(pb.Txs, func(t *trantorpb.Transaction) *types.Transaction {
			return types.TransactionFromPb(t)
		}),
	}
}

func (m *FillerMessage) Pb() *bcpb.FillerMessage {
	if m == nil {
		return nil
	}
	pbMessage := &bcpb.FillerMessage{}
	{
		if m.Cert != nil {
			pbMessage.Cert = (m.Cert).Pb()
		}
		pbMessage.Txs = types1.ConvertSlice(m.Txs, func(t *types.Transaction) *trantorpb.Transaction {
			return (t).Pb()
		})
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
	Slot      *Slot
	BatchId   string
	Signature tctypes.FullSig
}

func CertFromPb(pb *bcpb.Cert) *Cert {
	if pb == nil {
		return nil
	}
	return &Cert{
		Slot:      SlotFromPb(pb.Slot),
		BatchId:   pb.BatchId,
		Signature: (tctypes.FullSig)(pb.Signature),
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
		pbMessage.BatchId = m.BatchId
		pbMessage.Signature = ([]uint8)(m.Signature)
	}

	return pbMessage
}

func (*Cert) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcpb.Cert]()}
}
