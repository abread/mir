package abbapbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	abbatypes "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	abbapb "github.com/filecoin-project/mir/pkg/pb/abbapb"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() abbapb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb abbapb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *abbapb.Event_InputValue:
		return &Event_InputValue{InputValue: InputValueFromPb(pb.InputValue)}
	case *abbapb.Event_Deliver:
		return &Event_Deliver{Deliver: DeliverFromPb(pb.Deliver)}
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

func (w *Event_InputValue) Pb() abbapb.Event_Type {
	return &abbapb.Event_InputValue{InputValue: (w.InputValue).Pb()}
}

func (*Event_InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Event_InputValue]()}
}

type Event_Deliver struct {
	Deliver *Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() abbapb.Event_Type {
	return &abbapb.Event_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Event_Deliver]()}
}

func EventFromPb(pb *abbapb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *abbapb.Event {
	return &abbapb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Event]()}
}

type InputValue struct {
	Input bool
}

func InputValueFromPb(pb *abbapb.InputValue) *InputValue {
	return &InputValue{
		Input: pb.Input,
	}
}

func (m *InputValue) Pb() *abbapb.InputValue {
	return &abbapb.InputValue{
		Input: m.Input,
	}
}

func (*InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.InputValue]()}
}

type Deliver struct {
	Result       bool
	OriginModule types.ModuleID
}

func DeliverFromPb(pb *abbapb.Deliver) *Deliver {
	return &Deliver{
		Result:       pb.Result,
		OriginModule: (types.ModuleID)(pb.OriginModule),
	}
}

func (m *Deliver) Pb() *abbapb.Deliver {
	return &abbapb.Deliver{
		Result:       m.Result,
		OriginModule: (string)(m.OriginModule),
	}
}

func (*Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Deliver]()}
}

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() abbapb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb abbapb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *abbapb.Message_FinishMessage:
		return &Message_FinishMessage{FinishMessage: FinishMessageFromPb(pb.FinishMessage)}
	case *abbapb.Message_InitMessage:
		return &Message_InitMessage{InitMessage: InitMessageFromPb(pb.InitMessage)}
	case *abbapb.Message_AuxMessage:
		return &Message_AuxMessage{AuxMessage: AuxMessageFromPb(pb.AuxMessage)}
	case *abbapb.Message_ConfMessage:
		return &Message_ConfMessage{ConfMessage: ConfMessageFromPb(pb.ConfMessage)}
	case *abbapb.Message_CoinMessage:
		return &Message_CoinMessage{CoinMessage: CoinMessageFromPb(pb.CoinMessage)}
	}
	return nil
}

type Message_FinishMessage struct {
	FinishMessage *FinishMessage
}

func (*Message_FinishMessage) isMessage_Type() {}

func (w *Message_FinishMessage) Unwrap() *FinishMessage {
	return w.FinishMessage
}

func (w *Message_FinishMessage) Pb() abbapb.Message_Type {
	return &abbapb.Message_FinishMessage{FinishMessage: (w.FinishMessage).Pb()}
}

func (*Message_FinishMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Message_FinishMessage]()}
}

type Message_InitMessage struct {
	InitMessage *InitMessage
}

func (*Message_InitMessage) isMessage_Type() {}

func (w *Message_InitMessage) Unwrap() *InitMessage {
	return w.InitMessage
}

func (w *Message_InitMessage) Pb() abbapb.Message_Type {
	return &abbapb.Message_InitMessage{InitMessage: (w.InitMessage).Pb()}
}

func (*Message_InitMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Message_InitMessage]()}
}

type Message_AuxMessage struct {
	AuxMessage *AuxMessage
}

func (*Message_AuxMessage) isMessage_Type() {}

func (w *Message_AuxMessage) Unwrap() *AuxMessage {
	return w.AuxMessage
}

func (w *Message_AuxMessage) Pb() abbapb.Message_Type {
	return &abbapb.Message_AuxMessage{AuxMessage: (w.AuxMessage).Pb()}
}

func (*Message_AuxMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Message_AuxMessage]()}
}

type Message_ConfMessage struct {
	ConfMessage *ConfMessage
}

func (*Message_ConfMessage) isMessage_Type() {}

func (w *Message_ConfMessage) Unwrap() *ConfMessage {
	return w.ConfMessage
}

func (w *Message_ConfMessage) Pb() abbapb.Message_Type {
	return &abbapb.Message_ConfMessage{ConfMessage: (w.ConfMessage).Pb()}
}

func (*Message_ConfMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Message_ConfMessage]()}
}

type Message_CoinMessage struct {
	CoinMessage *CoinMessage
}

func (*Message_CoinMessage) isMessage_Type() {}

func (w *Message_CoinMessage) Unwrap() *CoinMessage {
	return w.CoinMessage
}

func (w *Message_CoinMessage) Pb() abbapb.Message_Type {
	return &abbapb.Message_CoinMessage{CoinMessage: (w.CoinMessage).Pb()}
}

func (*Message_CoinMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Message_CoinMessage]()}
}

func MessageFromPb(pb *abbapb.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *abbapb.Message {
	return &abbapb.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Message]()}
}

type FinishMessage struct {
	Value bool
}

func FinishMessageFromPb(pb *abbapb.FinishMessage) *FinishMessage {
	return &FinishMessage{
		Value: pb.Value,
	}
}

func (m *FinishMessage) Pb() *abbapb.FinishMessage {
	return &abbapb.FinishMessage{
		Value: m.Value,
	}
}

func (*FinishMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.FinishMessage]()}
}

type InitMessage struct {
	RoundNumber uint64
	Estimate    bool
}

func InitMessageFromPb(pb *abbapb.InitMessage) *InitMessage {
	return &InitMessage{
		RoundNumber: pb.RoundNumber,
		Estimate:    pb.Estimate,
	}
}

func (m *InitMessage) Pb() *abbapb.InitMessage {
	return &abbapb.InitMessage{
		RoundNumber: m.RoundNumber,
		Estimate:    m.Estimate,
	}
}

func (*InitMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.InitMessage]()}
}

type AuxMessage struct {
	RoundNumber uint64
	Value       bool
}

func AuxMessageFromPb(pb *abbapb.AuxMessage) *AuxMessage {
	return &AuxMessage{
		RoundNumber: pb.RoundNumber,
		Value:       pb.Value,
	}
}

func (m *AuxMessage) Pb() *abbapb.AuxMessage {
	return &abbapb.AuxMessage{
		RoundNumber: m.RoundNumber,
		Value:       m.Value,
	}
}

func (*AuxMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.AuxMessage]()}
}

type ConfMessage struct {
	RoundNumber uint64
	Values      abbatypes.ValueSet
}

func ConfMessageFromPb(pb *abbapb.ConfMessage) *ConfMessage {
	return &ConfMessage{
		RoundNumber: pb.RoundNumber,
		Values:      (abbatypes.ValueSet)(pb.Values),
	}
}

func (m *ConfMessage) Pb() *abbapb.ConfMessage {
	return &abbapb.ConfMessage{
		RoundNumber: m.RoundNumber,
		Values:      (uint32)(m.Values),
	}
}

func (*ConfMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.ConfMessage]()}
}

type CoinMessage struct {
	RoundNumber uint64
	CoinShare   tctypes.SigShare
}

func CoinMessageFromPb(pb *abbapb.CoinMessage) *CoinMessage {
	return &CoinMessage{
		RoundNumber: pb.RoundNumber,
		CoinShare:   (tctypes.SigShare)(pb.CoinShare),
	}
}

func (m *CoinMessage) Pb() *abbapb.CoinMessage {
	return &abbapb.CoinMessage{
		RoundNumber: m.RoundNumber,
		CoinShare:   ([]uint8)(m.CoinShare),
	}
}

func (*CoinMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.CoinMessage]()}
}
