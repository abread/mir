// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: aleapb/agreementpb/agreementpb.proto

package agreementpb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	_ "github.com/filecoin-project/mir/pkg/pb/mir"
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
	//	*Event_RequestInput
	//	*Event_InputValue
	//	*Event_Deliver
	Type isEvent_Type `protobuf_oneof:"type"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[0]
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
	return file_aleapb_agreementpb_agreementpb_proto_rawDescGZIP(), []int{0}
}

func (m *Event) GetType() isEvent_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Event) GetRequestInput() *RequestInput {
	if x, ok := x.GetType().(*Event_RequestInput); ok {
		return x.RequestInput
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

type Event_RequestInput struct {
	RequestInput *RequestInput `protobuf:"bytes,1,opt,name=request_input,json=requestInput,proto3,oneof"`
}

type Event_InputValue struct {
	InputValue *InputValue `protobuf:"bytes,2,opt,name=input_value,json=inputValue,proto3,oneof"`
}

type Event_Deliver struct {
	Deliver *Deliver `protobuf:"bytes,3,opt,name=deliver,proto3,oneof"`
}

func (*Event_RequestInput) isEvent_Type() {}

func (*Event_InputValue) isEvent_Type() {}

func (*Event_Deliver) isEvent_Type() {}

type RequestInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Round uint64 `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
}

func (x *RequestInput) Reset() {
	*x = RequestInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestInput) ProtoMessage() {}

func (x *RequestInput) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestInput.ProtoReflect.Descriptor instead.
func (*RequestInput) Descriptor() ([]byte, []int) {
	return file_aleapb_agreementpb_agreementpb_proto_rawDescGZIP(), []int{1}
}

func (x *RequestInput) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

type InputValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Round uint64 `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	Input bool   `protobuf:"varint,2,opt,name=input,proto3" json:"input,omitempty"`
}

func (x *InputValue) Reset() {
	*x = InputValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputValue) ProtoMessage() {}

func (x *InputValue) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[2]
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
	return file_aleapb_agreementpb_agreementpb_proto_rawDescGZIP(), []int{2}
}

func (x *InputValue) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
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

	Round    uint64 `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	Decision bool   `protobuf:"varint,2,opt,name=decision,proto3" json:"decision,omitempty"`
}

func (x *Deliver) Reset() {
	*x = Deliver{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deliver) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deliver) ProtoMessage() {}

func (x *Deliver) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[3]
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
	return file_aleapb_agreementpb_agreementpb_proto_rawDescGZIP(), []int{3}
}

func (x *Deliver) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *Deliver) GetDecision() bool {
	if x != nil {
		return x.Decision
	}
	return false
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*Message_FinishAbba
	Type isMessage_Type `protobuf_oneof:"type"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[4]
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
	return file_aleapb_agreementpb_agreementpb_proto_rawDescGZIP(), []int{4}
}

func (m *Message) GetType() isMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Message) GetFinishAbba() *FinishAbbaMessage {
	if x, ok := x.GetType().(*Message_FinishAbba); ok {
		return x.FinishAbba
	}
	return nil
}

type isMessage_Type interface {
	isMessage_Type()
}

type Message_FinishAbba struct {
	FinishAbba *FinishAbbaMessage `protobuf:"bytes,1,opt,name=finish_abba,json=finishAbba,proto3,oneof"`
}

func (*Message_FinishAbba) isMessage_Type() {}

type FinishAbbaMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Round uint64 `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	Value bool   `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *FinishAbbaMessage) Reset() {
	*x = FinishAbbaMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinishAbbaMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishAbbaMessage) ProtoMessage() {}

func (x *FinishAbbaMessage) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agreementpb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishAbbaMessage.ProtoReflect.Descriptor instead.
func (*FinishAbbaMessage) Descriptor() ([]byte, []int) {
	return file_aleapb_agreementpb_agreementpb_proto_rawDescGZIP(), []int{5}
}

func (x *FinishAbbaMessage) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *FinishAbbaMessage) GetValue() bool {
	if x != nil {
		return x.Value
	}
	return false
}

var File_aleapb_agreementpb_agreementpb_proto protoreflect.FileDescriptor

var file_aleapb_agreementpb_agreementpb_proto_rawDesc = []byte{
	0x0a, 0x24, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x61,
	0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x1a, 0x1c, 0x6d, 0x69, 0x72, 0x2f,
	0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xda, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x47, 0x0a, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x6c, 0x65, 0x61,
	0x70, 0x62, 0x2e, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x48, 0x00, 0x52, 0x0c, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x41, 0x0a, 0x0b, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1e, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x37,
	0x0a, 0x07, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x48, 0x00, 0x52, 0x07,
	0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0x24, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x22, 0x38, 0x0a, 0x0a, 0x49,
	0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75,
	0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x22, 0x3b, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72,
	0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x22, 0x5b, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x48, 0x0a,
	0x0b, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x5f, 0x61, 0x62, 0x62, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x61, 0x67, 0x72, 0x65,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x41, 0x62,
	0x62, 0x61, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x66, 0x69, 0x6e,
	0x69, 0x73, 0x68, 0x41, 0x62, 0x62, 0x61, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22,
	0x3f, 0x0a, 0x11, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x41, 0x62, 0x62, 0x61, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x61, 0x6c, 0x65, 0x61, 0x70,
	0x62, 0x2f, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aleapb_agreementpb_agreementpb_proto_rawDescOnce sync.Once
	file_aleapb_agreementpb_agreementpb_proto_rawDescData = file_aleapb_agreementpb_agreementpb_proto_rawDesc
)

func file_aleapb_agreementpb_agreementpb_proto_rawDescGZIP() []byte {
	file_aleapb_agreementpb_agreementpb_proto_rawDescOnce.Do(func() {
		file_aleapb_agreementpb_agreementpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_aleapb_agreementpb_agreementpb_proto_rawDescData)
	})
	return file_aleapb_agreementpb_agreementpb_proto_rawDescData
}

var file_aleapb_agreementpb_agreementpb_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_aleapb_agreementpb_agreementpb_proto_goTypes = []interface{}{
	(*Event)(nil),             // 0: aleapb.agreementpb.Event
	(*RequestInput)(nil),      // 1: aleapb.agreementpb.RequestInput
	(*InputValue)(nil),        // 2: aleapb.agreementpb.InputValue
	(*Deliver)(nil),           // 3: aleapb.agreementpb.Deliver
	(*Message)(nil),           // 4: aleapb.agreementpb.Message
	(*FinishAbbaMessage)(nil), // 5: aleapb.agreementpb.FinishAbbaMessage
}
var file_aleapb_agreementpb_agreementpb_proto_depIdxs = []int32{
	1, // 0: aleapb.agreementpb.Event.request_input:type_name -> aleapb.agreementpb.RequestInput
	2, // 1: aleapb.agreementpb.Event.input_value:type_name -> aleapb.agreementpb.InputValue
	3, // 2: aleapb.agreementpb.Event.deliver:type_name -> aleapb.agreementpb.Deliver
	5, // 3: aleapb.agreementpb.Message.finish_abba:type_name -> aleapb.agreementpb.FinishAbbaMessage
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_aleapb_agreementpb_agreementpb_proto_init() }
func file_aleapb_agreementpb_agreementpb_proto_init() {
	if File_aleapb_agreementpb_agreementpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aleapb_agreementpb_agreementpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_agreementpb_agreementpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestInput); i {
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
		file_aleapb_agreementpb_agreementpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_agreementpb_agreementpb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_agreementpb_agreementpb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_agreementpb_agreementpb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinishAbbaMessage); i {
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
	file_aleapb_agreementpb_agreementpb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Event_RequestInput)(nil),
		(*Event_InputValue)(nil),
		(*Event_Deliver)(nil),
	}
	file_aleapb_agreementpb_agreementpb_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*Message_FinishAbba)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aleapb_agreementpb_agreementpb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aleapb_agreementpb_agreementpb_proto_goTypes,
		DependencyIndexes: file_aleapb_agreementpb_agreementpb_proto_depIdxs,
		MessageInfos:      file_aleapb_agreementpb_agreementpb_proto_msgTypes,
	}.Build()
	File_aleapb_agreementpb_agreementpb_proto = out.File
	file_aleapb_agreementpb_agreementpb_proto_rawDesc = nil
	file_aleapb_agreementpb_agreementpb_proto_goTypes = nil
	file_aleapb_agreementpb_agreementpb_proto_depIdxs = nil
}