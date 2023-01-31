// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: abbapb/abbapb.proto

package abbapb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	_ "github.com/filecoin-project/mir/pkg/pb/mir"
	_ "github.com/filecoin-project/mir/pkg/pb/net"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*Event_InputValue
	//	*Event_Deliver
	Type isEvent_Type `protobuf_oneof:"type"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_abbapb_abbapb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{0}
}

func (m *Event) GetType() isEvent_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Event) GetInputValue() *InputValue {
	if x, ok := x.GetType().(*Event_InputValue); ok {
		return x.InputValue
	}
	return nil
}

func (x *Event) GetDeliver() *Deliver {
	if x, ok := x.GetType().(*Event_Deliver); ok {
		return x.Deliver
	}
	return nil
}

type isEvent_Type interface {
	isEvent_Type()
}

type Event_InputValue struct {
	InputValue *InputValue `protobuf:"bytes,1,opt,name=input_value,json=inputValue,proto3,oneof"`
}

type Event_Deliver struct {
	Deliver *Deliver `protobuf:"bytes,2,opt,name=deliver,proto3,oneof"`
}

func (*Event_InputValue) isEvent_Type() {}

func (*Event_Deliver) isEvent_Type() {}

type InputValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Input bool `protobuf:"varint,1,opt,name=input,proto3" json:"input,omitempty"`
}

func (x *InputValue) Reset() {
	*x = InputValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputValue) ProtoMessage() {}

func (x *InputValue) ProtoReflect() protoreflect.Message {
	mi := &file_abbapb_abbapb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputValue.ProtoReflect.Descriptor instead.
func (*InputValue) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{1}
}

func (x *InputValue) GetInput() bool {
	if x != nil {
		return x.Input
	}
	return false
}

type Deliver struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result       bool   `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	OriginModule string `protobuf:"bytes,2,opt,name=origin_module,json=originModule,proto3" json:"origin_module,omitempty"`
}

func (x *Deliver) Reset() {
	*x = Deliver{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deliver) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deliver) ProtoMessage() {}

func (x *Deliver) ProtoReflect() protoreflect.Message {
	mi := &file_abbapb_abbapb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deliver.ProtoReflect.Descriptor instead.
func (*Deliver) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{2}
}

func (x *Deliver) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

func (x *Deliver) GetOriginModule() string {
	if x != nil {
		return x.OriginModule
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*Message_FinishMessage
	//	*Message_InitMessage
	//	*Message_AuxMessage
	//	*Message_ConfMessage
	//	*Message_CoinMessage
	Type isMessage_Type `protobuf_oneof:"type"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_abbapb_abbapb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{3}
}

func (m *Message) GetType() isMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Message) GetFinishMessage() *FinishMessage {
	if x, ok := x.GetType().(*Message_FinishMessage); ok {
		return x.FinishMessage
	}
	return nil
}

func (x *Message) GetInitMessage() *InitMessage {
	if x, ok := x.GetType().(*Message_InitMessage); ok {
		return x.InitMessage
	}
	return nil
}

func (x *Message) GetAuxMessage() *AuxMessage {
	if x, ok := x.GetType().(*Message_AuxMessage); ok {
		return x.AuxMessage
	}
	return nil
}

func (x *Message) GetConfMessage() *ConfMessage {
	if x, ok := x.GetType().(*Message_ConfMessage); ok {
		return x.ConfMessage
	}
	return nil
}

func (x *Message) GetCoinMessage() *CoinMessage {
	if x, ok := x.GetType().(*Message_CoinMessage); ok {
		return x.CoinMessage
	}
	return nil
}

type isMessage_Type interface {
	isMessage_Type()
}

type Message_FinishMessage struct {
	FinishMessage *FinishMessage `protobuf:"bytes,1,opt,name=finish_message,json=finishMessage,proto3,oneof"`
}

type Message_InitMessage struct {
	InitMessage *InitMessage `protobuf:"bytes,2,opt,name=init_message,json=initMessage,proto3,oneof"`
}

type Message_AuxMessage struct {
	AuxMessage *AuxMessage `protobuf:"bytes,3,opt,name=aux_message,json=auxMessage,proto3,oneof"`
}

type Message_ConfMessage struct {
	ConfMessage *ConfMessage `protobuf:"bytes,4,opt,name=conf_message,json=confMessage,proto3,oneof"`
}

type Message_CoinMessage struct {
	CoinMessage *CoinMessage `protobuf:"bytes,5,opt,name=coin_message,json=coinMessage,proto3,oneof"`
}

func (*Message_FinishMessage) isMessage_Type() {}

func (*Message_InitMessage) isMessage_Type() {}

func (*Message_AuxMessage) isMessage_Type() {}

func (*Message_ConfMessage) isMessage_Type() {}

func (*Message_CoinMessage) isMessage_Type() {}

type FinishMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value bool `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *FinishMessage) Reset() {
	*x = FinishMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinishMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishMessage) ProtoMessage() {}

func (x *FinishMessage) ProtoReflect() protoreflect.Message {
	mi := &file_abbapb_abbapb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishMessage.ProtoReflect.Descriptor instead.
func (*FinishMessage) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{4}
}

func (x *FinishMessage) GetValue() bool {
	if x != nil {
		return x.Value
	}
	return false
}

type InitMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoundNumber uint64 `protobuf:"varint,1,opt,name=round_number,json=roundNumber,proto3" json:"round_number,omitempty"`
	Estimate    bool   `protobuf:"varint,2,opt,name=estimate,proto3" json:"estimate,omitempty"`
}

func (x *InitMessage) Reset() {
	*x = InitMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitMessage) ProtoMessage() {}

func (x *InitMessage) ProtoReflect() protoreflect.Message {
	mi := &file_abbapb_abbapb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitMessage.ProtoReflect.Descriptor instead.
func (*InitMessage) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{5}
}

func (x *InitMessage) GetRoundNumber() uint64 {
	if x != nil {
		return x.RoundNumber
	}
	return 0
}

func (x *InitMessage) GetEstimate() bool {
	if x != nil {
		return x.Estimate
	}
	return false
}

type AuxMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoundNumber uint64 `protobuf:"varint,1,opt,name=round_number,json=roundNumber,proto3" json:"round_number,omitempty"`
	Value       bool   `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *AuxMessage) Reset() {
	*x = AuxMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuxMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuxMessage) ProtoMessage() {}

func (x *AuxMessage) ProtoReflect() protoreflect.Message {
	mi := &file_abbapb_abbapb_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuxMessage.ProtoReflect.Descriptor instead.
func (*AuxMessage) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{6}
}

func (x *AuxMessage) GetRoundNumber() uint64 {
	if x != nil {
		return x.RoundNumber
	}
	return 0
}

func (x *AuxMessage) GetValue() bool {
	if x != nil {
		return x.Value
	}
	return false
}

type ConfMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoundNumber uint64 `protobuf:"varint,1,opt,name=round_number,json=roundNumber,proto3" json:"round_number,omitempty"`
	Values      uint32 `protobuf:"varint,2,opt,name=values,proto3" json:"values,omitempty"`
}

func (x *ConfMessage) Reset() {
	*x = ConfMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfMessage) ProtoMessage() {}

func (x *ConfMessage) ProtoReflect() protoreflect.Message {
	mi := &file_abbapb_abbapb_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfMessage.ProtoReflect.Descriptor instead.
func (*ConfMessage) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{7}
}

func (x *ConfMessage) GetRoundNumber() uint64 {
	if x != nil {
		return x.RoundNumber
	}
	return 0
}

func (x *ConfMessage) GetValues() uint32 {
	if x != nil {
		return x.Values
	}
	return 0
}

type CoinMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoundNumber uint64 `protobuf:"varint,1,opt,name=round_number,json=roundNumber,proto3" json:"round_number,omitempty"`
	CoinShare   []byte `protobuf:"bytes,2,opt,name=coin_share,json=coinShare,proto3" json:"coin_share,omitempty"`
}

func (x *CoinMessage) Reset() {
	*x = CoinMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoinMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoinMessage) ProtoMessage() {}

func (x *CoinMessage) ProtoReflect() protoreflect.Message {
	mi := &file_abbapb_abbapb_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoinMessage.ProtoReflect.Descriptor instead.
func (*CoinMessage) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{8}
}

func (x *CoinMessage) GetRoundNumber() uint64 {
	if x != nil {
		return x.RoundNumber
	}
	return 0
}

func (x *CoinMessage) GetCoinShare() []byte {
	if x != nil {
		return x.CoinShare
	}
	return nil
}

var File_abbapb_abbapb_proto protoreflect.FileDescriptor

var file_abbapb_abbapb_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2f, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x1a, 0x1c, 0x6d,
	0x69, 0x72, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6e, 0x65, 0x74,
	0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7f, 0x0a, 0x05, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x35, 0x0a, 0x0b, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62,
	0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2b, 0x0a, 0x07, 0x64, 0x65, 0x6c,
	0x69, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x62, 0x62,
	0x61, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x48, 0x00, 0x52, 0x07, 0x64,
	0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x3a, 0x04, 0x90, 0xa6, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0x28, 0x0a, 0x0a, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x3a, 0x04,
	0x98, 0xa6, 0x1d, 0x01, 0x22, 0x84, 0x01, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x5b, 0x0a, 0x0d, 0x6f, 0x72, 0x69, 0x67,
	0x69, 0x6e, 0x5f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x36, 0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x4d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0xc2, 0x02, 0x0a, 0x07,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3e, 0x0a, 0x0e, 0x66, 0x69, 0x6e, 0x69, 0x73,
	0x68, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0d, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x38, 0x0a, 0x0c, 0x69, 0x6e, 0x69, 0x74, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x69, 0x6e, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x35, 0x0a, 0x0b, 0x61, 0x75, 0x78, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e,
	0x41, 0x75, 0x78, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x61, 0x75,
	0x78, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x38, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x66, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x38, 0x0a, 0x0c, 0x63, 0x6f, 0x69, 0x6e, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70,
	0x62, 0x2e, 0x43, 0x6f, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52,
	0x0b, 0x63, 0x6f, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x3a, 0x04, 0xc8, 0xe4,
	0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0xc8, 0xe4, 0x1d, 0x01,
	0x22, 0x2b, 0x0a, 0x0d, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x22, 0x52, 0x0a,
	0x0b, 0x49, 0x6e, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x72, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0b, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x1a, 0x0a, 0x08, 0x65, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x08, 0x65, 0x73, 0x74, 0x69, 0x6d, 0x61, 0x74, 0x65, 0x3a, 0x04, 0xd0, 0xe4, 0x1d,
	0x01, 0x22, 0x4b, 0x0a, 0x0a, 0x41, 0x75, 0x78, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x22, 0x8f,
	0x01, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x66, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x21,
	0x0a, 0x0c, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x57, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x42, 0x3f, 0x82, 0xa6, 0x1d, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x62, 0x62, 0x61, 0x2f,
	0x61, 0x62, 0x62, 0x61, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53,
	0x65, 0x74, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x3a, 0x04, 0xd0, 0xe4, 0x1d, 0x01,
	0x22, 0x55, 0x0a, 0x0b, 0x43, 0x6f, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6f, 0x69, 0x6e, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x6f, 0x69, 0x6e, 0x53, 0x68, 0x61, 0x72,
	0x65, 0x3a, 0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70,
	0x62, 0x2f, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_abbapb_abbapb_proto_rawDescOnce sync.Once
	file_abbapb_abbapb_proto_rawDescData = file_abbapb_abbapb_proto_rawDesc
)

func file_abbapb_abbapb_proto_rawDescGZIP() []byte {
	file_abbapb_abbapb_proto_rawDescOnce.Do(func() {
		file_abbapb_abbapb_proto_rawDescData = protoimpl.X.CompressGZIP(file_abbapb_abbapb_proto_rawDescData)
	})
	return file_abbapb_abbapb_proto_rawDescData
}

var file_abbapb_abbapb_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_abbapb_abbapb_proto_goTypes = []interface{}{
	(*Event)(nil),         // 0: abbapb.Event
	(*InputValue)(nil),    // 1: abbapb.InputValue
	(*Deliver)(nil),       // 2: abbapb.Deliver
	(*Message)(nil),       // 3: abbapb.Message
	(*FinishMessage)(nil), // 4: abbapb.FinishMessage
	(*InitMessage)(nil),   // 5: abbapb.InitMessage
	(*AuxMessage)(nil),    // 6: abbapb.AuxMessage
	(*ConfMessage)(nil),   // 7: abbapb.ConfMessage
	(*CoinMessage)(nil),   // 8: abbapb.CoinMessage
}
var file_abbapb_abbapb_proto_depIdxs = []int32{
	1, // 0: abbapb.Event.input_value:type_name -> abbapb.InputValue
	2, // 1: abbapb.Event.deliver:type_name -> abbapb.Deliver
	4, // 2: abbapb.Message.finish_message:type_name -> abbapb.FinishMessage
	5, // 3: abbapb.Message.init_message:type_name -> abbapb.InitMessage
	6, // 4: abbapb.Message.aux_message:type_name -> abbapb.AuxMessage
	7, // 5: abbapb.Message.conf_message:type_name -> abbapb.ConfMessage
	8, // 6: abbapb.Message.coin_message:type_name -> abbapb.CoinMessage
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_abbapb_abbapb_proto_init() }
func file_abbapb_abbapb_proto_init() {
	if File_abbapb_abbapb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_abbapb_abbapb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abbapb_abbapb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InputValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abbapb_abbapb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Deliver); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abbapb_abbapb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abbapb_abbapb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinishMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abbapb_abbapb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abbapb_abbapb_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuxMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abbapb_abbapb_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abbapb_abbapb_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoinMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_abbapb_abbapb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Event_InputValue)(nil),
		(*Event_Deliver)(nil),
	}
	file_abbapb_abbapb_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Message_FinishMessage)(nil),
		(*Message_InitMessage)(nil),
		(*Message_AuxMessage)(nil),
		(*Message_ConfMessage)(nil),
		(*Message_CoinMessage)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_abbapb_abbapb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_abbapb_abbapb_proto_goTypes,
		DependencyIndexes: file_abbapb_abbapb_proto_depIdxs,
		MessageInfos:      file_abbapb_abbapb_proto_msgTypes,
	}.Build()
	File_abbapb_abbapb_proto = out.File
	file_abbapb_abbapb_proto_rawDesc = nil
	file_abbapb_abbapb_proto_goTypes = nil
	file_abbapb_abbapb_proto_depIdxs = nil
}
