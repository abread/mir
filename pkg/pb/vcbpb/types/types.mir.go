// Code generated by Mir codegen. DO NOT EDIT.

package vcbpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types1 "github.com/filecoin-project/mir/codegen/model/types"
	trantorpb "github.com/filecoin-project/mir/pkg/pb/trantorpb"
	types "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	vcbpb "github.com/filecoin-project/mir/pkg/pb/vcbpb"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types2 "github.com/filecoin-project/mir/pkg/trantor/types"
	types3 "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() vcbpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb vcbpb.Event_Type) Event_Type {
	if pb == nil {
		return nil
	}
	switch pb := pb.(type) {
	case *vcbpb.Event_InputValue:
		return &Event_InputValue{InputValue: InputValueFromPb(pb.InputValue)}
	case *vcbpb.Event_Deliver:
		return &Event_Deliver{Deliver: DeliverFromPb(pb.Deliver)}
	case *vcbpb.Event_QuorumDone:
		return &Event_QuorumDone{QuorumDone: QuorumDoneFromPb(pb.QuorumDone)}
	case *vcbpb.Event_FullyDone:
		return &Event_FullyDone{FullyDone: FullyDoneFromPb(pb.FullyDone)}
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

func (w *Event_InputValue) Pb() vcbpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.InputValue == nil {
		return &vcbpb.Event_InputValue{}
	}
	return &vcbpb.Event_InputValue{InputValue: (w.InputValue).Pb()}
}

func (*Event_InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Event_InputValue]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() vcbpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.Deliver == nil {
		return &vcbpb.Event_Deliver{}
	}
	return &vcbpb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Event_Deliver]()}
}

type Event_QuorumDone struct {
	QuorumDone *QuorumDone
}

func (*Event_QuorumDone) isEvent_Type() {}

func (w *Event_QuorumDone) Unwrap() *QuorumDone {
	return w.QuorumDone
}

func (w *Event_QuorumDone) Pb() vcbpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.QuorumDone == nil {
		return &vcbpb.Event_QuorumDone{}
	}
	return &vcbpb.Event_QuorumDone{QuorumDone: (w.QuorumDone).Pb()}
}

func (*Event_QuorumDone) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Event_QuorumDone]()}
}

type Event_FullyDone struct {
	FullyDone *FullyDone
}

func (*Event_FullyDone) isEvent_Type() {}

func (w *Event_FullyDone) Unwrap() *FullyDone {
	return w.FullyDone
}

func (w *Event_FullyDone) Pb() vcbpb.Event_Type {
	if w == nil {
		return nil
	}
	if w.FullyDone == nil {
		return &vcbpb.Event_FullyDone{}
	}
	return &vcbpb.Event_FullyDone{FullyDone: (w.FullyDone).Pb()}
}

func (*Event_FullyDone) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Event_FullyDone]()}
}

func EventFromPb(pb *vcbpb.Event) *Event {
	if pb == nil {
		return nil
	}
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *vcbpb.Event {
	if m == nil {
		return nil
	}
	pbMessage := &vcbpb.Event{}
	{
		if m.Type != nil {
			pbMessage.Type = (m.Type).Pb()
		}
	}

	return pbMessage
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Event]()}
}

type InputValue struct {
	Txs []*types.Transaction
}

func InputValueFromPb(pb *vcbpb.InputValue) *InputValue {
	if pb == nil {
		return nil
	}
	return &InputValue{
		Txs: types1.ConvertSlice(pb.Txs, func(t *trantorpb.Transaction) *types.Transaction {
			return types.TransactionFromPb(t)
		}),
	}
}

func (m *InputValue) Pb() *vcbpb.InputValue {
	if m == nil {
		return nil
	}
	pbMessage := &vcbpb.InputValue{}
	{
		pbMessage.Txs = types1.ConvertSlice(m.Txs, func(t *types.Transaction) *trantorpb.Transaction {
			return (t).Pb()
		})
	}

	return pbMessage
}

func (*InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.InputValue]()}
}

type Deliver struct {
	Txs       []*types.Transaction
	TxIds     []types2.TxID
	Signature tctypes.FullSig
	SrcModule types3.ModuleID
}

func DeliverFromPb(pb *vcbpb.Deliver) *Deliver {
	if pb == nil {
		return nil
	}
	return &Deliver{
		Txs: types1.ConvertSlice(pb.Txs, func(t *trantorpb.Transaction) *types.Transaction {
			return types.TransactionFromPb(t)
		}),
		TxIds: types1.ConvertSlice(pb.TxIds, func(t []uint8) types2.TxID {
			return (types2.TxID)(t)
		}),
		Signature: (tctypes.FullSig)(pb.Signature),
		SrcModule: (types3.ModuleID)(pb.SrcModule),
	}
}

func (m *Deliver) Pb() *vcbpb.Deliver {
	if m == nil {
		return nil
	}
	pbMessage := &vcbpb.Deliver{}
	{
		pbMessage.Txs = types1.ConvertSlice(m.Txs, func(t *types.Transaction) *trantorpb.Transaction {
			return (t).Pb()
		})
		pbMessage.TxIds = types1.ConvertSlice(m.TxIds, func(t types2.TxID) []uint8 {
			return ([]uint8)(t)
		})
		pbMessage.Signature = ([]uint8)(m.Signature)
		pbMessage.SrcModule = (string)(m.SrcModule)
	}

	return pbMessage
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Deliver]()}
}

type QuorumDone struct {
	SrcModule types3.ModuleID
}

func QuorumDoneFromPb(pb *vcbpb.QuorumDone) *QuorumDone {
	if pb == nil {
		return nil
	}
	return &QuorumDone{
		SrcModule: (types3.ModuleID)(pb.SrcModule),
	}
}

func (m *QuorumDone) Pb() *vcbpb.QuorumDone {
	if m == nil {
		return nil
	}
	pbMessage := &vcbpb.QuorumDone{}
	{
		pbMessage.SrcModule = (string)(m.SrcModule)
	}

	return pbMessage
}

func (*QuorumDone) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.QuorumDone]()}
}

type FullyDone struct {
	SrcModule types3.ModuleID
}

func FullyDoneFromPb(pb *vcbpb.FullyDone) *FullyDone {
	if pb == nil {
		return nil
	}
	return &FullyDone{
		SrcModule: (types3.ModuleID)(pb.SrcModule),
	}
}

func (m *FullyDone) Pb() *vcbpb.FullyDone {
	if m == nil {
		return nil
	}
	pbMessage := &vcbpb.FullyDone{}
	{
		pbMessage.SrcModule = (string)(m.SrcModule)
	}

	return pbMessage
}

func (*FullyDone) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.FullyDone]()}
}

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() vcbpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb vcbpb.Message_Type) Message_Type {
	if pb == nil {
		return nil
	}
	switch pb := pb.(type) {
	case *vcbpb.Message_SendMessage:
		return &Message_SendMessage{SendMessage: SendMessageFromPb(pb.SendMessage)}
	case *vcbpb.Message_EchoMessage:
		return &Message_EchoMessage{EchoMessage: EchoMessageFromPb(pb.EchoMessage)}
	case *vcbpb.Message_FinalMessage:
		return &Message_FinalMessage{FinalMessage: FinalMessageFromPb(pb.FinalMessage)}
	case *vcbpb.Message_DoneMessage:
		return &Message_DoneMessage{DoneMessage: DoneMessageFromPb(pb.DoneMessage)}
	}
	return nil
}

type Message_SendMessage struct {
	SendMessage *SendMessage
}

func (*Message_SendMessage) isMessage_Type() {}

func (w *Message_SendMessage) Unwrap() *SendMessage {
	return w.SendMessage
}

func (w *Message_SendMessage) Pb() vcbpb.Message_Type {
	if w == nil {
		return nil
	}
	if w.SendMessage == nil {
		return &vcbpb.Message_SendMessage{}
	}
	return &vcbpb.Message_SendMessage{SendMessage: (w.SendMessage).Pb()}
}

func (*Message_SendMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Message_SendMessage]()}
}

type Message_EchoMessage struct {
	EchoMessage *EchoMessage
}

func (*Message_EchoMessage) isMessage_Type() {}

func (w *Message_EchoMessage) Unwrap() *EchoMessage {
	return w.EchoMessage
}

func (w *Message_EchoMessage) Pb() vcbpb.Message_Type {
	if w == nil {
		return nil
	}
	if w.EchoMessage == nil {
		return &vcbpb.Message_EchoMessage{}
	}
	return &vcbpb.Message_EchoMessage{EchoMessage: (w.EchoMessage).Pb()}
}

func (*Message_EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Message_EchoMessage]()}
}

type Message_FinalMessage struct {
	FinalMessage *FinalMessage
}

func (*Message_FinalMessage) isMessage_Type() {}

func (w *Message_FinalMessage) Unwrap() *FinalMessage {
	return w.FinalMessage
}

func (w *Message_FinalMessage) Pb() vcbpb.Message_Type {
	if w == nil {
		return nil
	}
	if w.FinalMessage == nil {
		return &vcbpb.Message_FinalMessage{}
	}
	return &vcbpb.Message_FinalMessage{FinalMessage: (w.FinalMessage).Pb()}
}

func (*Message_FinalMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Message_FinalMessage]()}
}

type Message_DoneMessage struct {
	DoneMessage *DoneMessage
}

func (*Message_DoneMessage) isMessage_Type() {}

func (w *Message_DoneMessage) Unwrap() *DoneMessage {
	return w.DoneMessage
}

func (w *Message_DoneMessage) Pb() vcbpb.Message_Type {
	if w == nil {
		return nil
	}
	if w.DoneMessage == nil {
		return &vcbpb.Message_DoneMessage{}
	}
	return &vcbpb.Message_DoneMessage{DoneMessage: (w.DoneMessage).Pb()}
}

func (*Message_DoneMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Message_DoneMessage]()}
}

func MessageFromPb(pb *vcbpb.Message) *Message {
	if pb == nil {
		return nil
	}
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *vcbpb.Message {
	if m == nil {
		return nil
	}
	pbMessage := &vcbpb.Message{}
	{
		if m.Type != nil {
			pbMessage.Type = (m.Type).Pb()
		}
	}

	return pbMessage
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.Message]()}
}

type SendMessage struct {
	Txs []*types.Transaction
}

func SendMessageFromPb(pb *vcbpb.SendMessage) *SendMessage {
	if pb == nil {
		return nil
	}
	return &SendMessage{
		Txs: types1.ConvertSlice(pb.Txs, func(t *trantorpb.Transaction) *types.Transaction {
			return types.TransactionFromPb(t)
		}),
	}
}

func (m *SendMessage) Pb() *vcbpb.SendMessage {
	if m == nil {
		return nil
	}
	pbMessage := &vcbpb.SendMessage{}
	{
		pbMessage.Txs = types1.ConvertSlice(m.Txs, func(t *types.Transaction) *trantorpb.Transaction {
			return (t).Pb()
		})
	}

	return pbMessage
}

func (*SendMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.SendMessage]()}
}

type EchoMessage struct {
	SignatureShare tctypes.SigShare
}

func EchoMessageFromPb(pb *vcbpb.EchoMessage) *EchoMessage {
	if pb == nil {
		return nil
	}
	return &EchoMessage{
		SignatureShare: (tctypes.SigShare)(pb.SignatureShare),
	}
}

func (m *EchoMessage) Pb() *vcbpb.EchoMessage {
	if m == nil {
		return nil
	}
	pbMessage := &vcbpb.EchoMessage{}
	{
		pbMessage.SignatureShare = ([]uint8)(m.SignatureShare)
	}

	return pbMessage
}

func (*EchoMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.EchoMessage]()}
}

type FinalMessage struct {
	Signature tctypes.FullSig
}

func FinalMessageFromPb(pb *vcbpb.FinalMessage) *FinalMessage {
	if pb == nil {
		return nil
	}
	return &FinalMessage{
		Signature: (tctypes.FullSig)(pb.Signature),
	}
}

func (m *FinalMessage) Pb() *vcbpb.FinalMessage {
	if m == nil {
		return nil
	}
	pbMessage := &vcbpb.FinalMessage{}
	{
		pbMessage.Signature = ([]uint8)(m.Signature)
	}

	return pbMessage
}

func (*FinalMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.FinalMessage]()}
}

type DoneMessage struct{}

func DoneMessageFromPb(pb *vcbpb.DoneMessage) *DoneMessage {
	if pb == nil {
		return nil
	}
	return &DoneMessage{}
}

func (m *DoneMessage) Pb() *vcbpb.DoneMessage {
	if m == nil {
		return nil
	}
	pbMessage := &vcbpb.DoneMessage{}
	{
	}

	return pbMessage
}

func (*DoneMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*vcbpb.DoneMessage]()}
}
