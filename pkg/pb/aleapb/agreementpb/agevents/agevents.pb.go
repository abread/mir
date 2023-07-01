// separate file to avoid circular protobuf dependencies

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1
// 	protoc        v3.21.12
// source: aleapb/agreementpb/agevents/agevents.proto

package agevents

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	_ "github.com/filecoin-project/mir/pkg/pb/mir"
	modringpb "github.com/filecoin-project/mir/pkg/pb/modringpb"
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
	//	*Event_StaleMsgsRevcd
	//	*Event_InnerAbbaRoundTime
	Type isEvent_Type `protobuf_oneof:"type"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[0]
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
	return file_aleapb_agreementpb_agevents_agevents_proto_rawDescGZIP(), []int{0}
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

func (x *Event) GetStaleMsgsRevcd() *StaleMsgsRecvd {
	if x, ok := x.GetType().(*Event_StaleMsgsRevcd); ok {
		return x.StaleMsgsRevcd
	}
	return nil
}

func (x *Event) GetInnerAbbaRoundTime() *InnerAbbaRoundTime {
	if x, ok := x.GetType().(*Event_InnerAbbaRoundTime); ok {
		return x.InnerAbbaRoundTime
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

type Event_StaleMsgsRevcd struct {
	StaleMsgsRevcd *StaleMsgsRecvd `protobuf:"bytes,3,opt,name=stale_msgs_revcd,json=staleMsgsRevcd,proto3,oneof"`
}

type Event_InnerAbbaRoundTime struct {
	InnerAbbaRoundTime *InnerAbbaRoundTime `protobuf:"bytes,4,opt,name=inner_abba_round_time,json=innerAbbaRoundTime,proto3,oneof"`
}

func (*Event_InputValue) isEvent_Type() {}

func (*Event_Deliver) isEvent_Type() {}

func (*Event_StaleMsgsRevcd) isEvent_Type() {}

func (*Event_InnerAbbaRoundTime) isEvent_Type() {}

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
		mi := &file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputValue) ProtoMessage() {}

func (x *InputValue) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[1]
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
	return file_aleapb_agreementpb_agevents_agevents_proto_rawDescGZIP(), []int{1}
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
	// Amount of time from input to a quorum(2F+1) of inputs with value 1.
	// Can be negative, when quorum is reached before input.
	// Set to math.MaxInt64 when quorum is never reached.
	PosQuorumWait int64 `protobuf:"varint,3,opt,name=pos_quorum_wait,json=posQuorumWait,proto3" json:"pos_quorum_wait,omitempty"`
	// Amount of time from a quorum(2F+1) of 1-valued inputs to all nodes giving 1-valued inputs.
	// Set to math.MaxInt64 when quroum is never reached.
	PosTotalWait int64 `protobuf:"varint,4,opt,name=pos_total_wait,json=posTotalWait,proto3" json:"pos_total_wait,omitempty"`
}

func (x *Deliver) Reset() {
	*x = Deliver{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deliver) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deliver) ProtoMessage() {}

func (x *Deliver) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[2]
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
	return file_aleapb_agreementpb_agevents_agevents_proto_rawDescGZIP(), []int{2}
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

func (x *Deliver) GetPosQuorumWait() int64 {
	if x != nil {
		return x.PosQuorumWait
	}
	return 0
}

func (x *Deliver) GetPosTotalWait() int64 {
	if x != nil {
		return x.PosTotalWait
	}
	return 0
}

type InnerAbbaRoundTime struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Duration of an "inner" ABBA round, excluding coin messages and their processing.
	DurationNoCoin int64 `protobuf:"varint,1,opt,name=duration_no_coin,json=durationNoCoin,proto3" json:"duration_no_coin,omitempty"`
}

func (x *InnerAbbaRoundTime) Reset() {
	*x = InnerAbbaRoundTime{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InnerAbbaRoundTime) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InnerAbbaRoundTime) ProtoMessage() {}

func (x *InnerAbbaRoundTime) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InnerAbbaRoundTime.ProtoReflect.Descriptor instead.
func (*InnerAbbaRoundTime) Descriptor() ([]byte, []int) {
	return file_aleapb_agreementpb_agevents_agevents_proto_rawDescGZIP(), []int{3}
}

func (x *InnerAbbaRoundTime) GetDurationNoCoin() int64 {
	if x != nil {
		return x.DurationNoCoin
	}
	return 0
}

type StaleMsgsRecvd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []*modringpb.PastMessage `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *StaleMsgsRecvd) Reset() {
	*x = StaleMsgsRecvd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StaleMsgsRecvd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StaleMsgsRecvd) ProtoMessage() {}

func (x *StaleMsgsRecvd) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StaleMsgsRecvd.ProtoReflect.Descriptor instead.
func (*StaleMsgsRecvd) Descriptor() ([]byte, []int) {
	return file_aleapb_agreementpb_agevents_agevents_proto_rawDescGZIP(), []int{4}
}

func (x *StaleMsgsRecvd) GetMessages() []*modringpb.PastMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

var File_aleapb_agreementpb_agevents_agevents_proto protoreflect.FileDescriptor

var file_aleapb_agreementpb_agevents_agevents_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x61, 0x67,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x61, 0x6c,
	0x65, 0x61, 0x70, 0x62, 0x2e, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62,
	0x1a, 0x19, 0x6d, 0x6f, 0x64, 0x72, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2f, 0x6d, 0x6f, 0x64, 0x72,
	0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72,
	0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc4, 0x02, 0x0a, 0x05, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x41, 0x0a, 0x0b, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70,
	0x62, 0x2e, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x37, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62,
	0x2e, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c,
	0x69, 0x76, 0x65, 0x72, 0x48, 0x00, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x12,
	0x4e, 0x0a, 0x10, 0x73, 0x74, 0x61, 0x6c, 0x65, 0x5f, 0x6d, 0x73, 0x67, 0x73, 0x5f, 0x72, 0x65,
	0x76, 0x63, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61, 0x6c, 0x65, 0x61,
	0x70, 0x62, 0x2e, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x53,
	0x74, 0x61, 0x6c, 0x65, 0x4d, 0x73, 0x67, 0x73, 0x52, 0x65, 0x63, 0x76, 0x64, 0x48, 0x00, 0x52,
	0x0e, 0x73, 0x74, 0x61, 0x6c, 0x65, 0x4d, 0x73, 0x67, 0x73, 0x52, 0x65, 0x76, 0x63, 0x64, 0x12,
	0x5b, 0x0a, 0x15, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x61, 0x62, 0x62, 0x61, 0x5f, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26,
	0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x70, 0x62, 0x2e, 0x49, 0x6e, 0x6e, 0x65, 0x72, 0x41, 0x62, 0x62, 0x61, 0x52, 0x6f, 0x75,
	0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x48, 0x00, 0x52, 0x12, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x41,
	0x62, 0x62, 0x61, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x3a, 0x04, 0x90, 0xa6,
	0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0x80, 0xa6, 0x1d, 0x01,
	0x22, 0x3e, 0x0a, 0x0a, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01,
	0x22, 0xb5, 0x01, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05,
	0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75,
	0x6e, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x39,
	0x0a, 0x0f, 0x70, 0x6f, 0x73, 0x5f, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x5f, 0x77, 0x61, 0x69,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x82, 0xa6, 0x1d, 0x0d, 0x74, 0x69, 0x6d,
	0x65, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x70, 0x6f, 0x73, 0x51,
	0x75, 0x6f, 0x72, 0x75, 0x6d, 0x57, 0x61, 0x69, 0x74, 0x12, 0x37, 0x0a, 0x0e, 0x70, 0x6f, 0x73,
	0x5f, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x77, 0x61, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x03, 0x42, 0x11, 0x82, 0xa6, 0x1d, 0x0d, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x70, 0x6f, 0x73, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x57, 0x61,
	0x69, 0x74, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0x57, 0x0a, 0x12, 0x49, 0x6e, 0x6e, 0x65,
	0x72, 0x41, 0x62, 0x62, 0x61, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3b,
	0x0a, 0x10, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x6f, 0x5f, 0x63, 0x6f,
	0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x82, 0xa6, 0x1d, 0x0d, 0x74, 0x69,
	0x6d, 0x65, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0e, 0x64, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x43, 0x6f, 0x69, 0x6e, 0x3a, 0x04, 0x98, 0xa6, 0x1d,
	0x01, 0x22, 0x4a, 0x0a, 0x0e, 0x53, 0x74, 0x61, 0x6c, 0x65, 0x4d, 0x73, 0x67, 0x73, 0x52, 0x65,
	0x63, 0x76, 0x64, 0x12, 0x32, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x6f, 0x64, 0x72, 0x69, 0x6e, 0x67, 0x70,
	0x62, 0x2e, 0x50, 0x61, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x42, 0x44, 0x5a,
	0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65,
	0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f, 0x61,
	0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aleapb_agreementpb_agevents_agevents_proto_rawDescOnce sync.Once
	file_aleapb_agreementpb_agevents_agevents_proto_rawDescData = file_aleapb_agreementpb_agevents_agevents_proto_rawDesc
)

func file_aleapb_agreementpb_agevents_agevents_proto_rawDescGZIP() []byte {
	file_aleapb_agreementpb_agevents_agevents_proto_rawDescOnce.Do(func() {
		file_aleapb_agreementpb_agevents_agevents_proto_rawDescData = protoimpl.X.CompressGZIP(file_aleapb_agreementpb_agevents_agevents_proto_rawDescData)
	})
	return file_aleapb_agreementpb_agevents_agevents_proto_rawDescData
}

var file_aleapb_agreementpb_agevents_agevents_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_aleapb_agreementpb_agevents_agevents_proto_goTypes = []interface{}{
	(*Event)(nil),                 // 0: aleapb.agreementpb.Event
	(*InputValue)(nil),            // 1: aleapb.agreementpb.InputValue
	(*Deliver)(nil),               // 2: aleapb.agreementpb.Deliver
	(*InnerAbbaRoundTime)(nil),    // 3: aleapb.agreementpb.InnerAbbaRoundTime
	(*StaleMsgsRecvd)(nil),        // 4: aleapb.agreementpb.StaleMsgsRecvd
	(*modringpb.PastMessage)(nil), // 5: modringpb.PastMessage
}
var file_aleapb_agreementpb_agevents_agevents_proto_depIdxs = []int32{
	1, // 0: aleapb.agreementpb.Event.input_value:type_name -> aleapb.agreementpb.InputValue
	2, // 1: aleapb.agreementpb.Event.deliver:type_name -> aleapb.agreementpb.Deliver
	4, // 2: aleapb.agreementpb.Event.stale_msgs_revcd:type_name -> aleapb.agreementpb.StaleMsgsRecvd
	3, // 3: aleapb.agreementpb.Event.inner_abba_round_time:type_name -> aleapb.agreementpb.InnerAbbaRoundTime
	5, // 4: aleapb.agreementpb.StaleMsgsRecvd.messages:type_name -> modringpb.PastMessage
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_aleapb_agreementpb_agevents_agevents_proto_init() }
func file_aleapb_agreementpb_agevents_agevents_proto_init() {
	if File_aleapb_agreementpb_agevents_agevents_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InnerAbbaRoundTime); i {
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
		file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StaleMsgsRecvd); i {
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
	file_aleapb_agreementpb_agevents_agevents_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Event_InputValue)(nil),
		(*Event_Deliver)(nil),
		(*Event_StaleMsgsRevcd)(nil),
		(*Event_InnerAbbaRoundTime)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aleapb_agreementpb_agevents_agevents_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aleapb_agreementpb_agevents_agevents_proto_goTypes,
		DependencyIndexes: file_aleapb_agreementpb_agevents_agevents_proto_depIdxs,
		MessageInfos:      file_aleapb_agreementpb_agevents_agevents_proto_msgTypes,
	}.Build()
	File_aleapb_agreementpb_agevents_agevents_proto = out.File
	file_aleapb_agreementpb_agevents_agevents_proto_rawDesc = nil
	file_aleapb_agreementpb_agevents_agevents_proto_goTypes = nil
	file_aleapb_agreementpb_agevents_agevents_proto_depIdxs = nil
}
