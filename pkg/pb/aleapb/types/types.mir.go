package aleapbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types2 "github.com/filecoin-project/mir/codegen/model/types"
	aleapb "github.com/filecoin-project/mir/pkg/pb/aleapb"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	trantorpb "github.com/filecoin-project/mir/pkg/pb/trantorpb"
	types1 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() aleapb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb aleapb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *aleapb.Message_FillGapMessage:
		return &Message_FillGapMessage{FillGapMessage: FillGapMessageFromPb(pb.FillGapMessage)}
	case *aleapb.Message_FillerMessage:
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

func (w *Message_FillGapMessage) Pb() aleapb.Message_Type {
	return &aleapb.Message_FillGapMessage{FillGapMessage: (w.FillGapMessage).Pb()}
}

func (*Message_FillGapMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*aleapb.Message_FillGapMessage]()}
}

type Message_FillerMessage struct {
	FillerMessage *FillerMessage
}

func (*Message_FillerMessage) isMessage_Type() {}

func (w *Message_FillerMessage) Unwrap() *FillerMessage {
	return w.FillerMessage
}

func (w *Message_FillerMessage) Pb() aleapb.Message_Type {
	return &aleapb.Message_FillerMessage{FillerMessage: (w.FillerMessage).Pb()}
}

func (*Message_FillerMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*aleapb.Message_FillerMessage]()}
}

func MessageFromPb(pb *aleapb.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *aleapb.Message {
	return &aleapb.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*aleapb.Message]()}
}

type FillGapMessage struct {
	Slot *types.Slot
}

func FillGapMessageFromPb(pb *aleapb.FillGapMessage) *FillGapMessage {
	return &FillGapMessage{
		Slot: types.SlotFromPb(pb.Slot),
	}
}

func (m *FillGapMessage) Pb() *aleapb.FillGapMessage {
	return &aleapb.FillGapMessage{
		Slot: (m.Slot).Pb(),
	}
}

func (*FillGapMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*aleapb.FillGapMessage]()}
}

type FillerMessage struct {
	Slot      *types.Slot
	Txs       []*types1.Transaction
	Signature tctypes.FullSig
}

func FillerMessageFromPb(pb *aleapb.FillerMessage) *FillerMessage {
	return &FillerMessage{
		Slot: types.SlotFromPb(pb.Slot),
		Txs: types2.ConvertSlice(pb.Txs, func(t *trantorpb.Transaction) *types1.Transaction {
			return types1.TransactionFromPb(t)
		}),
		Signature: (tctypes.FullSig)(pb.Signature),
	}
}

func (m *FillerMessage) Pb() *aleapb.FillerMessage {
	return &aleapb.FillerMessage{
		Slot: (m.Slot).Pb(),
		Txs: types2.ConvertSlice(m.Txs, func(t *types1.Transaction) *trantorpb.Transaction {
			return (t).Pb()
		}),
		Signature: ([]uint8)(m.Signature),
	}
}

func (*FillerMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*aleapb.FillerMessage]()}
}

type Cert struct {
	Slot *types.Slot
}

func CertFromPb(pb *aleapb.Cert) *Cert {
	return &Cert{
		Slot: types.SlotFromPb(pb.Slot),
	}
}

func (m *Cert) Pb() *aleapb.Cert {
	return &aleapb.Cert{
		Slot: (m.Slot).Pb(),
	}
}

func (*Cert) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*aleapb.Cert]()}
}
