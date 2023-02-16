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

	agreementpb "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
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
	//	*Event_Round
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

func (x *Event) GetRound() *RoundEvent {
	if x, ok := x.GetType().(*Event_Round); ok {
		return x.Round
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

type Event_Round struct {
	Round *RoundEvent `protobuf:"bytes,3,opt,name=round,proto3,oneof"`
}

func (*Event_InputValue) isEvent_Type() {}

func (*Event_Deliver) isEvent_Type() {}

func (*Event_Round) isEvent_Type() {}

type InputValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Input  bool    `protobuf:"varint,1,opt,name=input,proto3" json:"input,omitempty"`
	Origin *Origin `protobuf:"bytes,2,opt,name=origin,proto3" json:"origin,omitempty"`
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

func (x *InputValue) GetOrigin() *Origin {
	if x != nil {
		return x.Origin
	}
	return nil
}

type Deliver struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result bool    `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	Origin *Origin `protobuf:"bytes,2,opt,name=origin,proto3" json:"origin,omitempty"`
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

func (x *Deliver) GetOrigin() *Origin {
	if x != nil {
		return x.Origin
	}
	return nil
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

type Origin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Module string `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
	// Types that are assignable to Type:
	//
	//	*Origin_ContextStore
	//	*Origin_Dsl
	//	*Origin_AleaAg
	Type isOrigin_Type `protobuf_oneof:"Type"`
}

func (x *Origin) Reset() {
	*x = Origin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abbapb_abbapb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Origin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Origin) ProtoMessage() {}

func (x *Origin) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Origin.ProtoReflect.Descriptor instead.
func (*Origin) Descriptor() ([]byte, []int) {
	return file_abbapb_abbapb_proto_rawDescGZIP(), []int{5}
}

func (x *Origin) GetModule() string {
	if x != nil {
		return x.Module
	}
	return ""
}

func (m *Origin) GetType() isOrigin_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Origin) GetContextStore() *contextstorepb.Origin {
	if x, ok := x.GetType().(*Origin_ContextStore); ok {
		return x.ContextStore
	}
	return nil
}

func (x *Origin) GetDsl() *dslpb.Origin {
	if x, ok := x.GetType().(*Origin_Dsl); ok {
		return x.Dsl
	}
	return nil
}

func (x *Origin) GetAleaAg() *agreementpb.AbbaOrigin {
	if x, ok := x.GetType().(*Origin_AleaAg); ok {
		return x.AleaAg
	}
	return nil
}

type isOrigin_Type interface {
	isOrigin_Type()
}

type Origin_ContextStore struct {
	ContextStore *contextstorepb.Origin `protobuf:"bytes,2,opt,name=context_store,json=contextStore,proto3,oneof"`
}

type Origin_Dsl struct {
	Dsl *dslpb.Origin `protobuf:"bytes,3,opt,name=dsl,proto3,oneof"`
}

type Origin_AleaAg struct {
	AleaAg *agreementpb.AbbaOrigin `protobuf:"bytes,4,opt,name=alea_ag,json=aleaAg,proto3,oneof"`
}

func (*Origin_ContextStore) isOrigin_Type() {}

func (*Origin_Dsl) isOrigin_Type() {}

func (*Origin_AleaAg) isOrigin_Type() {}

var File_abbapb_abbapb_proto protoreflect.FileDescriptor

var file_abbapb_abbapb_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2f, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x1a, 0x14, 0x61,
	0x62, 0x62, 0x61, 0x70, 0x62, 0x2f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x70, 0x62, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x64, 0x73, 0x6c, 0x70, 0x62, 0x2f, 0x64, 0x73, 0x6c, 0x70, 0x62,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x61, 0x6c, 0x65,
	0x61, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2f,
	0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65,
	0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x6e, 0x65, 0x74, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xab, 0x01,
	0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x35, 0x0a, 0x0b, 0x69, 0x6e, 0x70, 0x75, 0x74,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61,
	0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2b,
	0x0a, 0x07, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72,
	0x48, 0x00, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x05, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x62, 0x62,
	0x61, 0x70, 0x62, 0x2e, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00,
	0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x3a, 0x04, 0x90, 0xa6, 0x1d, 0x01, 0x42, 0x0c, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0x56, 0x0a, 0x0a, 0x49,
	0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12,
	0x2c, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x42,
	0x04, 0x98, 0xa6, 0x1d, 0x01, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x3a, 0x04, 0x98,
	0xa6, 0x1d, 0x01, 0x22, 0x55, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x2c, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e,
	0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x42, 0x04, 0xa0, 0xa6, 0x1d, 0x01, 0x52, 0x06, 0x6f, 0x72,
	0x69, 0x67, 0x69, 0x6e, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0x7c, 0x0a, 0x07, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2f, 0x0a, 0x06, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x46,
	0x69, 0x6e, 0x69, 0x73, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x06,
	0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x12, 0x2c, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x52,
	0x6f, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x05, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x3a, 0x04, 0xc8, 0xe4, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x04, 0xc8, 0xe4, 0x1d, 0x01, 0x22, 0x2b, 0x0a, 0x0d, 0x46, 0x69, 0x6e, 0x69,
	0x73, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x22, 0x83, 0x02, 0x0a, 0x06, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x12, 0x4e, 0x0a, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x36, 0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e,
	0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65,
	0x12, 0x3d, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78,
	0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x48,
	0x00, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12,
	0x21, 0x0a, 0x03, 0x64, 0x73, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x64,
	0x73, 0x6c, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x48, 0x00, 0x52, 0x03, 0x64,
	0x73, 0x6c, 0x12, 0x39, 0x0a, 0x07, 0x61, 0x6c, 0x65, 0x61, 0x5f, 0x61, 0x67, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x61, 0x67, 0x72,
	0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x41, 0x62, 0x62, 0x61, 0x4f, 0x72, 0x69,
	0x67, 0x69, 0x6e, 0x48, 0x00, 0x52, 0x06, 0x61, 0x6c, 0x65, 0x61, 0x41, 0x67, 0x3a, 0x04, 0x80,
	0xa6, 0x1d, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x42, 0x2f, 0x5a, 0x2d, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f,
	0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
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

var file_abbapb_abbapb_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_abbapb_abbapb_proto_goTypes = []interface{}{
	(*Event)(nil),                  // 0: abbapb.Event
	(*InputValue)(nil),             // 1: abbapb.InputValue
	(*Deliver)(nil),                // 2: abbapb.Deliver
	(*Message)(nil),                // 3: abbapb.Message
	(*FinishMessage)(nil),          // 4: abbapb.FinishMessage
	(*Origin)(nil),                 // 5: abbapb.Origin
	(*RoundEvent)(nil),             // 6: abbapb.RoundEvent
	(*RoundMessage)(nil),           // 7: abbapb.RoundMessage
	(*contextstorepb.Origin)(nil),  // 8: contextstorepb.Origin
	(*dslpb.Origin)(nil),           // 9: dslpb.Origin
	(*agreementpb.AbbaOrigin)(nil), // 10: aleapb.agreementpb.AbbaOrigin
}
var file_abbapb_abbapb_proto_depIdxs = []int32{
	1,  // 0: abbapb.Event.input_value:type_name -> abbapb.InputValue
	2,  // 1: abbapb.Event.deliver:type_name -> abbapb.Deliver
	6,  // 2: abbapb.Event.round:type_name -> abbapb.RoundEvent
	5,  // 3: abbapb.InputValue.origin:type_name -> abbapb.Origin
	5,  // 4: abbapb.Deliver.origin:type_name -> abbapb.Origin
	4,  // 5: abbapb.Message.finish:type_name -> abbapb.FinishMessage
	7,  // 6: abbapb.Message.round:type_name -> abbapb.RoundMessage
	8,  // 7: abbapb.Origin.context_store:type_name -> contextstorepb.Origin
	9,  // 8: abbapb.Origin.dsl:type_name -> dslpb.Origin
	10, // 9: abbapb.Origin.alea_ag:type_name -> aleapb.agreementpb.AbbaOrigin
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
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
			switch v := v.(*Origin); i {
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
		(*Event_Round)(nil),
	}
	file_abbapb_abbapb_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Message_Finish)(nil),
		(*Message_Round)(nil),
	}
	file_abbapb_abbapb_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*Origin_ContextStore)(nil),
		(*Origin_Dsl)(nil),
		(*Origin_AleaAg)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_abbapb_abbapb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
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
