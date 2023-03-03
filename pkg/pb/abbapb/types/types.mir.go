package abbapbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	abbatypes "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	abbapb "github.com/filecoin-project/mir/pkg/pb/abbapb"
	types3 "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/contextstorepb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/dslpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type RoundEvent struct {
	Type RoundEvent_Type
}

type RoundEvent_Type interface {
	mirreflect.GeneratedType
	isRoundEvent_Type()
	Pb() abbapb.RoundEvent_Type
}

type RoundEvent_TypeWrapper[T any] interface {
	RoundEvent_Type
	Unwrap() *T
}

func RoundEvent_TypeFromPb(pb abbapb.RoundEvent_Type) RoundEvent_Type {
	switch pb := pb.(type) {
	case *abbapb.RoundEvent_InputValue:
		return &RoundEvent_InputValue{InputValue: RoundInputValueFromPb(pb.InputValue)}
	case *abbapb.RoundEvent_Deliver:
		return &RoundEvent_Deliver{Deliver: RoundDeliverFromPb(pb.Deliver)}
	case *abbapb.RoundEvent_Finish:
		return &RoundEvent_Finish{Finish: RoundFinishAllFromPb(pb.Finish)}
	}
	return nil
}

type RoundEvent_InputValue struct {
	InputValue *RoundInputValue
}

func (*RoundEvent_InputValue) isRoundEvent_Type() {}

func (w *RoundEvent_InputValue) Unwrap() *RoundInputValue {
	return w.InputValue
}

func (w *RoundEvent_InputValue) Pb() abbapb.RoundEvent_Type {
	return &abbapb.RoundEvent_InputValue{InputValue: (w.InputValue).Pb()}
}

func (*RoundEvent_InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundEvent_InputValue]()}
}

type RoundEvent_Deliver struct {
	Deliver *RoundDeliver
}

func (*RoundEvent_Deliver) isRoundEvent_Type() {}

func (w *RoundEvent_Deliver) Unwrap() *RoundDeliver {
	return w.Deliver
}

func (w *RoundEvent_Deliver) Pb() abbapb.RoundEvent_Type {
	return &abbapb.RoundEvent_Deliver{Deliver: (w.Deliver).Pb()}
}

func (*RoundEvent_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundEvent_Deliver]()}
}

type RoundEvent_Finish struct {
	Finish *RoundFinishAll
}

func (*RoundEvent_Finish) isRoundEvent_Type() {}

func (w *RoundEvent_Finish) Unwrap() *RoundFinishAll {
	return w.Finish
}

func (w *RoundEvent_Finish) Pb() abbapb.RoundEvent_Type {
	return &abbapb.RoundEvent_Finish{Finish: (w.Finish).Pb()}
}

func (*RoundEvent_Finish) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundEvent_Finish]()}
}

func RoundEventFromPb(pb *abbapb.RoundEvent) *RoundEvent {
	return &RoundEvent{
		Type: RoundEvent_TypeFromPb(pb.Type),
	}
}

func (m *RoundEvent) Pb() *abbapb.RoundEvent {
	return &abbapb.RoundEvent{
		Type: (m.Type).Pb(),
	}
}

func (*RoundEvent) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundEvent]()}
}

type RoundInputValue struct {
	Input bool
}

func RoundInputValueFromPb(pb *abbapb.RoundInputValue) *RoundInputValue {
	return &RoundInputValue{
		Input: pb.Input,
	}
}

func (m *RoundInputValue) Pb() *abbapb.RoundInputValue {
	return &abbapb.RoundInputValue{
		Input: m.Input,
	}
}

func (*RoundInputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundInputValue]()}
}

type RoundDeliver struct {
	NextEstimate bool
	RoundNumber  uint64
}

func RoundDeliverFromPb(pb *abbapb.RoundDeliver) *RoundDeliver {
	return &RoundDeliver{
		NextEstimate: pb.NextEstimate,
		RoundNumber:  pb.RoundNumber,
	}
}

func (m *RoundDeliver) Pb() *abbapb.RoundDeliver {
	return &abbapb.RoundDeliver{
		NextEstimate: m.NextEstimate,
		RoundNumber:  m.RoundNumber,
	}
}

func (*RoundDeliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundDeliver]()}
}

type RoundFinishAll struct {
	Decision bool
}

func RoundFinishAllFromPb(pb *abbapb.RoundFinishAll) *RoundFinishAll {
	return &RoundFinishAll{
		Decision: pb.Decision,
	}
}

func (m *RoundFinishAll) Pb() *abbapb.RoundFinishAll {
	return &abbapb.RoundFinishAll{
		Decision: m.Decision,
	}
}

func (*RoundFinishAll) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundFinishAll]()}
}

type RoundMessage struct {
	Type RoundMessage_Type
}

type RoundMessage_Type interface {
	mirreflect.GeneratedType
	isRoundMessage_Type()
	Pb() abbapb.RoundMessage_Type
}

type RoundMessage_TypeWrapper[T any] interface {
	RoundMessage_Type
	Unwrap() *T
}

func RoundMessage_TypeFromPb(pb abbapb.RoundMessage_Type) RoundMessage_Type {
	switch pb := pb.(type) {
	case *abbapb.RoundMessage_Init:
		return &RoundMessage_Init{Init: RoundInitMessageFromPb(pb.Init)}
	case *abbapb.RoundMessage_Aux:
		return &RoundMessage_Aux{Aux: RoundAuxMessageFromPb(pb.Aux)}
	case *abbapb.RoundMessage_Conf:
		return &RoundMessage_Conf{Conf: RoundConfMessageFromPb(pb.Conf)}
	case *abbapb.RoundMessage_Coin:
		return &RoundMessage_Coin{Coin: RoundCoinMessageFromPb(pb.Coin)}
	}
	return nil
}

type RoundMessage_Init struct {
	Init *RoundInitMessage
}

func (*RoundMessage_Init) isRoundMessage_Type() {}

func (w *RoundMessage_Init) Unwrap() *RoundInitMessage {
	return w.Init
}

func (w *RoundMessage_Init) Pb() abbapb.RoundMessage_Type {
	return &abbapb.RoundMessage_Init{Init: (w.Init).Pb()}
}

func (*RoundMessage_Init) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundMessage_Init]()}
}

type RoundMessage_Aux struct {
	Aux *RoundAuxMessage
}

func (*RoundMessage_Aux) isRoundMessage_Type() {}

func (w *RoundMessage_Aux) Unwrap() *RoundAuxMessage {
	return w.Aux
}

func (w *RoundMessage_Aux) Pb() abbapb.RoundMessage_Type {
	return &abbapb.RoundMessage_Aux{Aux: (w.Aux).Pb()}
}

func (*RoundMessage_Aux) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundMessage_Aux]()}
}

type RoundMessage_Conf struct {
	Conf *RoundConfMessage
}

func (*RoundMessage_Conf) isRoundMessage_Type() {}

func (w *RoundMessage_Conf) Unwrap() *RoundConfMessage {
	return w.Conf
}

func (w *RoundMessage_Conf) Pb() abbapb.RoundMessage_Type {
	return &abbapb.RoundMessage_Conf{Conf: (w.Conf).Pb()}
}

func (*RoundMessage_Conf) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundMessage_Conf]()}
}

type RoundMessage_Coin struct {
	Coin *RoundCoinMessage
}

func (*RoundMessage_Coin) isRoundMessage_Type() {}

func (w *RoundMessage_Coin) Unwrap() *RoundCoinMessage {
	return w.Coin
}

func (w *RoundMessage_Coin) Pb() abbapb.RoundMessage_Type {
	return &abbapb.RoundMessage_Coin{Coin: (w.Coin).Pb()}
}

func (*RoundMessage_Coin) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundMessage_Coin]()}
}

func RoundMessageFromPb(pb *abbapb.RoundMessage) *RoundMessage {
	return &RoundMessage{
		Type: RoundMessage_TypeFromPb(pb.Type),
	}
}

func (m *RoundMessage) Pb() *abbapb.RoundMessage {
	return &abbapb.RoundMessage{
		Type: (m.Type).Pb(),
	}
}

func (*RoundMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundMessage]()}
}

type RoundInitMessage struct {
	Estimate bool
}

func RoundInitMessageFromPb(pb *abbapb.RoundInitMessage) *RoundInitMessage {
	return &RoundInitMessage{
		Estimate: pb.Estimate,
	}
}

func (m *RoundInitMessage) Pb() *abbapb.RoundInitMessage {
	return &abbapb.RoundInitMessage{
		Estimate: m.Estimate,
	}
}

func (*RoundInitMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundInitMessage]()}
}

type RoundAuxMessage struct {
	Value bool
}

func RoundAuxMessageFromPb(pb *abbapb.RoundAuxMessage) *RoundAuxMessage {
	return &RoundAuxMessage{
		Value: pb.Value,
	}
}

func (m *RoundAuxMessage) Pb() *abbapb.RoundAuxMessage {
	return &abbapb.RoundAuxMessage{
		Value: m.Value,
	}
}

func (*RoundAuxMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundAuxMessage]()}
}

type RoundConfMessage struct {
	Values abbatypes.ValueSet
}

func RoundConfMessageFromPb(pb *abbapb.RoundConfMessage) *RoundConfMessage {
	return &RoundConfMessage{
		Values: (abbatypes.ValueSet)(pb.Values),
	}
}

func (m *RoundConfMessage) Pb() *abbapb.RoundConfMessage {
	return &abbapb.RoundConfMessage{
		Values: (uint32)(m.Values),
	}
}

func (*RoundConfMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundConfMessage]()}
}

type RoundCoinMessage struct {
	CoinShare tctypes.SigShare
}

func RoundCoinMessageFromPb(pb *abbapb.RoundCoinMessage) *RoundCoinMessage {
	return &RoundCoinMessage{
		CoinShare: (tctypes.SigShare)(pb.CoinShare),
	}
}

func (m *RoundCoinMessage) Pb() *abbapb.RoundCoinMessage {
	return &abbapb.RoundCoinMessage{
		CoinShare: ([]uint8)(m.CoinShare),
	}
}

func (*RoundCoinMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RoundCoinMessage]()}
}

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
	case *abbapb.Event_RequestInput:
		return &Event_RequestInput{RequestInput: RequestInputFromPb(pb.RequestInput)}
	case *abbapb.Event_InputValue:
		return &Event_InputValue{InputValue: InputValueFromPb(pb.InputValue)}
	case *abbapb.Event_Deliver:
		return &Event_Deliver{Deliver: DeliverFromPb(pb.Deliver)}
	case *abbapb.Event_Round:
		return &Event_Round{Round: RoundEventFromPb(pb.Round)}
	}
	return nil
}

type Event_RequestInput struct {
	RequestInput *RequestInput
}

func (*Event_RequestInput) isEvent_Type() {}

func (w *Event_RequestInput) Unwrap() *RequestInput {
	return w.RequestInput
}

func (w *Event_RequestInput) Pb() abbapb.Event_Type {
	return &abbapb.Event_RequestInput{RequestInput: (w.RequestInput).Pb()}
}

func (*Event_RequestInput) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Event_RequestInput]()}
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

type Event_Round struct {
	Round *RoundEvent
}

func (*Event_Round) isEvent_Type() {}

func (w *Event_Round) Unwrap() *RoundEvent {
	return w.Round
}

func (w *Event_Round) Pb() abbapb.Event_Type {
	return &abbapb.Event_Round{Round: (w.Round).Pb()}
}

func (*Event_Round) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Event_Round]()}
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

type RequestInput struct {
	Module types.ModuleID
}

func RequestInputFromPb(pb *abbapb.RequestInput) *RequestInput {
	return &RequestInput{
		Module: (types.ModuleID)(pb.Module),
	}
}

func (m *RequestInput) Pb() *abbapb.RequestInput {
	return &abbapb.RequestInput{
		Module: (string)(m.Module),
	}
}

func (*RequestInput) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.RequestInput]()}
}

type InputValue struct {
	Input  bool
	Origin *Origin
}

func InputValueFromPb(pb *abbapb.InputValue) *InputValue {
	return &InputValue{
		Input:  pb.Input,
		Origin: OriginFromPb(pb.Origin),
	}
}

func (m *InputValue) Pb() *abbapb.InputValue {
	return &abbapb.InputValue{
		Input:  m.Input,
		Origin: (m.Origin).Pb(),
	}
}

func (*InputValue) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.InputValue]()}
}

type Deliver struct {
	Result bool
	Origin *Origin
}

func DeliverFromPb(pb *abbapb.Deliver) *Deliver {
	return &Deliver{
		Result: pb.Result,
		Origin: OriginFromPb(pb.Origin),
	}
}

func (m *Deliver) Pb() *abbapb.Deliver {
	return &abbapb.Deliver{
		Result: m.Result,
		Origin: (m.Origin).Pb(),
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
	case *abbapb.Message_Finish:
		return &Message_Finish{Finish: FinishMessageFromPb(pb.Finish)}
	case *abbapb.Message_Round:
		return &Message_Round{Round: RoundMessageFromPb(pb.Round)}
	}
	return nil
}

type Message_Finish struct {
	Finish *FinishMessage
}

func (*Message_Finish) isMessage_Type() {}

func (w *Message_Finish) Unwrap() *FinishMessage {
	return w.Finish
}

func (w *Message_Finish) Pb() abbapb.Message_Type {
	return &abbapb.Message_Finish{Finish: (w.Finish).Pb()}
}

func (*Message_Finish) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Message_Finish]()}
}

type Message_Round struct {
	Round *RoundMessage
}

func (*Message_Round) isMessage_Type() {}

func (w *Message_Round) Unwrap() *RoundMessage {
	return w.Round
}

func (w *Message_Round) Pb() abbapb.Message_Type {
	return &abbapb.Message_Round{Round: (w.Round).Pb()}
}

func (*Message_Round) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Message_Round]()}
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

type Origin struct {
	Module types.ModuleID
	Type   Origin_Type
}

type Origin_Type interface {
	mirreflect.GeneratedType
	isOrigin_Type()
	Pb() abbapb.Origin_Type
}

type Origin_TypeWrapper[T any] interface {
	Origin_Type
	Unwrap() *T
}

func Origin_TypeFromPb(pb abbapb.Origin_Type) Origin_Type {
	switch pb := pb.(type) {
	case *abbapb.Origin_ContextStore:
		return &Origin_ContextStore{ContextStore: types1.OriginFromPb(pb.ContextStore)}
	case *abbapb.Origin_Dsl:
		return &Origin_Dsl{Dsl: types2.OriginFromPb(pb.Dsl)}
	case *abbapb.Origin_AleaAg:
		return &Origin_AleaAg{AleaAg: types3.AbbaOriginFromPb(pb.AleaAg)}
	}
	return nil
}

type Origin_ContextStore struct {
	ContextStore *types1.Origin
}

func (*Origin_ContextStore) isOrigin_Type() {}

func (w *Origin_ContextStore) Unwrap() *types1.Origin {
	return w.ContextStore
}

func (w *Origin_ContextStore) Pb() abbapb.Origin_Type {
	return &abbapb.Origin_ContextStore{ContextStore: (w.ContextStore).Pb()}
}

func (*Origin_ContextStore) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Origin_ContextStore]()}
}

type Origin_Dsl struct {
	Dsl *types2.Origin
}

func (*Origin_Dsl) isOrigin_Type() {}

func (w *Origin_Dsl) Unwrap() *types2.Origin {
	return w.Dsl
}

func (w *Origin_Dsl) Pb() abbapb.Origin_Type {
	return &abbapb.Origin_Dsl{Dsl: (w.Dsl).Pb()}
}

func (*Origin_Dsl) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Origin_Dsl]()}
}

type Origin_AleaAg struct {
	AleaAg *types3.AbbaOrigin
}

func (*Origin_AleaAg) isOrigin_Type() {}

func (w *Origin_AleaAg) Unwrap() *types3.AbbaOrigin {
	return w.AleaAg
}

func (w *Origin_AleaAg) Pb() abbapb.Origin_Type {
	return &abbapb.Origin_AleaAg{AleaAg: (w.AleaAg).Pb()}
}

func (*Origin_AleaAg) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Origin_AleaAg]()}
}

func OriginFromPb(pb *abbapb.Origin) *Origin {
	return &Origin{
		Module: (types.ModuleID)(pb.Module),
		Type:   Origin_TypeFromPb(pb.Type),
	}
}

func (m *Origin) Pb() *abbapb.Origin {
	return &abbapb.Origin{
		Module: (string)(m.Module),
		Type:   (m.Type).Pb(),
	}
}

func (*Origin) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*abbapb.Origin]()}
}
