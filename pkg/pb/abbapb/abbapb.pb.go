// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
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
	//	*Event_Continue
	//	*Event_Deliver
	//	*Event_Round
	//	*Event_Done
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

func (x *Event) GetContinue() *ContinueExecution {
	if x, ok := x.GetType().(*Event_Continue); ok {
		return x.Continue
	}
	return nil
}

func (x *Event) GetDeliver() *Deliver {
	if x, ok := x.GetType().(*Event_Deliver); ok {
		return x.Deliver
	}
	return nil
}

func (x *Event) GetRound() *RoundEvent {
	if x, ok := x.GetType().(*Event_Round); ok {
		return x.Round
	}
	return nil
}

func (x *Event) GetDone() *Done {
	if x, ok := x.GetType().(*Event_Done); ok {
		return x.Done
	}
	return nil
}

type isEvent_Type interface {
	isEvent_Type()
}

type Event_InputValue struct {
	InputValue *InputValue `protobuf:"bytes,1,opt,name=input_value,json=inputValue,proto3,oneof"`
}

type Event_Continue struct {
	Continue *ContinueExecution `protobuf:"bytes,2,opt,name=continue,proto3,oneof"`
}

type Event_Deliver struct {
	Deliver *Deliver `protobuf:"bytes,3,opt,name=deliver,proto3,oneof"`
}

type Event_Round struct {
	Round *RoundEvent `protobuf:"bytes,4,opt,name=round,proto3,oneof"`
}

type Event_Done struct {
	Done *Done `protobuf:"bytes,5,opt,name=done,proto3,oneof"`
}

func (*Event_InputValue) isEvent_Type() {}

func (*Event_Continue) isEvent_Type() {}

func (*Event_Deliver) isEvent_Type() {}

func (*Event_Round) isEvent_Type() {}

func (*Event_Done) isEvent_Type() {}

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

type ContinueExecution struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ContinueExecution) Reset() {
	*x = ContinueExecution{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContinueExecution) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContinueExecution) ProtoMessage() {}

func (x *ContinueExecution) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ContinueExecution.ProtoReflect.Descriptor instead.
func (*ContinueExecution) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{2}
}

type Deliver struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result    bool   `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	SrcModule string `protobuf:"bytes,2,opt,name=src_module,json=srcModule,proto3" json:"src_module,omitempty"`
}

func (x *Deliver) Reset() {
	*x = Deliver{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deliver) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deliver) ProtoMessage() {}

func (x *Deliver) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Deliver.ProtoReflect.Descriptor instead.
func (*Deliver) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{3}
}

func (x *Deliver) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

func (x *Deliver) GetSrcModule() string {
	if x != nil {
		return x.SrcModule
	}
	return ""
}

type Done struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcModule string `protobuf:"bytes,1,opt,name=src_module,json=srcModule,proto3" json:"src_module,omitempty"`
}

func (x *Done) Reset() {
	*x = Done{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Done) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Done) ProtoMessage() {}

func (x *Done) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Done.ProtoReflect.Descriptor instead.
func (*Done) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{4}
}

func (x *Done) GetSrcModule() string {
	if x != nil {
		return x.SrcModule
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*Message_Finish
	//	*Message_Round
	Type isMessage_Type `protobuf_oneof:"type"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{5}
}

func (m *Message) GetType() isMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Message) GetFinish() *FinishMessage {
	if x, ok := x.GetType().(*Message_Finish); ok {
		return x.Finish
	}
	return nil
}

func (x *Message) GetRound() *RoundMessage {
	if x, ok := x.GetType().(*Message_Round); ok {
		return x.Round
	}
	return nil
}

type isMessage_Type interface {
	isMessage_Type()
}

type Message_Finish struct {
	Finish *FinishMessage `protobuf:"bytes,1,opt,name=finish,proto3,oneof"`
}

type Message_Round struct {
	Round *RoundMessage `protobuf:"bytes,2,opt,name=round,proto3,oneof"`
}

func (*Message_Finish) isMessage_Type() {}

func (*Message_Round) isMessage_Type() {}

type FinishMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value bool `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *FinishMessage) Reset() {
	*x = FinishMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinishMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishMessage) ProtoMessage() {}

func (x *FinishMessage) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use FinishMessage.ProtoReflect.Descriptor instead.
func (*FinishMessage) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{6}
}

func (x *FinishMessage) GetValue() bool {
	if x != nil {
		return x.Value
	}
	return false
}

var File_abbapb_abbapb_proto protoreflect.FileDescriptor

var file_abbapb_abbapb_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2f, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x1a, 0x14, 0x61,
	0x62, 0x62, 0x61, 0x70, 0x62, 0x2f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x70, 0x62, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e,
	0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1c, 0x6e, 0x65, 0x74, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65,
	0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x88, 0x02, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x35, 0x0a, 0x0b, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x37, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74,
	0x69, 0x6e, 0x75, 0x65, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52,
	0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x12, 0x2b, 0x0a, 0x07, 0x64, 0x65, 0x6c,
	0x69, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x62, 0x62,
	0x61, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x48, 0x00, 0x52, 0x07, 0x64,
	0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x52,
	0x6f, 0x75, 0x6e, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x05, 0x72, 0x6f, 0x75,
	0x6e, 0x64, 0x12, 0x22, 0x0a, 0x04, 0x64, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x44, 0x6f, 0x6e, 0x65, 0x48, 0x00,
	0x52, 0x04, 0x64, 0x6f, 0x6e, 0x65, 0x3a, 0x04, 0x90, 0xa6, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0x28, 0x0a, 0x0a, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x3a, 0x04,
	0x98, 0xa6, 0x1d, 0x01, 0x22, 0x19, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65,
	0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22,
	0x7e, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x55, 0x0a, 0x0a, 0x73, 0x72, 0x63, 0x5f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x36, 0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x09,
	0x73, 0x72, 0x63, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22,
	0x63, 0x0a, 0x04, 0x44, 0x6f, 0x6e, 0x65, 0x12, 0x55, 0x0a, 0x0a, 0x73, 0x72, 0x63, 0x5f, 0x6d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x36, 0x82, 0xa6, 0x1d,
	0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65,
	0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x49, 0x44, 0x52, 0x09, 0x73, 0x72, 0x63, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x3a, 0x04,
	0x98, 0xa6, 0x1d, 0x01, 0x22, 0x7c, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x2f, 0x0a, 0x06, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x06, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68,
	0x12, 0x2c, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x3a, 0x04,
	0xc8, 0xe4, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0xc8, 0xe4,
	0x1d, 0x01, 0x22, 0x2b, 0x0a, 0x0d, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x42,
	0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69,
	0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d,
	0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_abbapb_abbapb_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_abbapb_abbapb_proto_goTypes = []interface{}{
	(*Event)(nil),             // 0: abbapb.Event
	(*InputValue)(nil),        // 1: abbapb.InputValue
	(*ContinueExecution)(nil), // 2: abbapb.ContinueExecution
	(*Deliver)(nil),           // 3: abbapb.Deliver
	(*Done)(nil),              // 4: abbapb.Done
	(*Message)(nil),           // 5: abbapb.Message
	(*FinishMessage)(nil),     // 6: abbapb.FinishMessage
	(*RoundEvent)(nil),        // 7: abbapb.RoundEvent
	(*RoundMessage)(nil),      // 8: abbapb.RoundMessage
}
var file_abbapb_abbapb_proto_depIdxs = []int32{
	1, // 0: abbapb.Event.input_value:type_name -> abbapb.InputValue
	2, // 1: abbapb.Event.continue:type_name -> abbapb.ContinueExecution
	3, // 2: abbapb.Event.deliver:type_name -> abbapb.Deliver
	7, // 3: abbapb.Event.round:type_name -> abbapb.RoundEvent
	4, // 4: abbapb.Event.done:type_name -> abbapb.Done
	6, // 5: abbapb.Message.finish:type_name -> abbapb.FinishMessage
	8, // 6: abbapb.Message.round:type_name -> abbapb.RoundMessage
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
	file_abbapb_roundpb_proto_init()
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
			switch v := v.(*ContinueExecution); i {
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
		file_abbapb_abbapb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Done); i {
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
		file_abbapb_abbapb_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
	}
	file_abbapb_abbapb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Event_InputValue)(nil),
		(*Event_Continue)(nil),
		(*Event_Deliver)(nil),
		(*Event_Round)(nil),
		(*Event_Done)(nil),
	}
	file_abbapb_abbapb_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*Message_Finish)(nil),
		(*Message_Round)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_abbapb_abbapb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
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
