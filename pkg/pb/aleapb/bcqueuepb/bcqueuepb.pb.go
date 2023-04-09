// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: aleapb/bcqueuepb/bcqueuepb.proto

package bcqueuepb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	common "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	_ "github.com/filecoin-project/mir/pkg/pb/mir"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
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
	//	*Event_FreeSlot
	//	*Event_PastVcbFinal
	//	*Event_BcStarted
	Type isEvent_Type `protobuf_oneof:"type"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[0]
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
	return file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescGZIP(), []int{0}
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

func (x *Event) GetFreeSlot() *FreeSlot {
	if x, ok := x.GetType().(*Event_FreeSlot); ok {
		return x.FreeSlot
	}
	return nil
}

func (x *Event) GetPastVcbFinal() *PastVcbFinal {
	if x, ok := x.GetType().(*Event_PastVcbFinal); ok {
		return x.PastVcbFinal
	}
	return nil
}

func (x *Event) GetBcStarted() *BcStarted {
	if x, ok := x.GetType().(*Event_BcStarted); ok {
		return x.BcStarted
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

type Event_FreeSlot struct {
	FreeSlot *FreeSlot `protobuf:"bytes,3,opt,name=free_slot,json=freeSlot,proto3,oneof"`
}

type Event_PastVcbFinal struct {
	PastVcbFinal *PastVcbFinal `protobuf:"bytes,4,opt,name=past_vcb_final,json=pastVcbFinal,proto3,oneof"`
}

type Event_BcStarted struct {
	BcStarted *BcStarted `protobuf:"bytes,5,opt,name=bc_started,json=bcStarted,proto3,oneof"`
}

func (*Event_InputValue) isEvent_Type() {}

func (*Event_Deliver) isEvent_Type() {}

func (*Event_FreeSlot) isEvent_Type() {}

func (*Event_PastVcbFinal) isEvent_Type() {}

func (*Event_BcStarted) isEvent_Type() {}

type InputValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Slot *common.Slot `protobuf:"bytes,1,opt,name=slot,proto3" json:"slot,omitempty"`
	// empty txs represents a request to listen for VCB (and deliver what was broadcast to us)
	Txs []*requestpb.Request `protobuf:"bytes,2,rep,name=txs,proto3" json:"txs,omitempty"`
}

func (x *InputValue) Reset() {
	*x = InputValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputValue) ProtoMessage() {}

func (x *InputValue) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[1]
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
	return file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescGZIP(), []int{1}
}

func (x *InputValue) GetSlot() *common.Slot {
	if x != nil {
		return x.Slot
	}
	return nil
}

func (x *InputValue) GetTxs() []*requestpb.Request {
	if x != nil {
		return x.Txs
	}
	return nil
}

type Deliver struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Slot *common.Slot `protobuf:"bytes,1,opt,name=slot,proto3" json:"slot,omitempty"`
}

func (x *Deliver) Reset() {
	*x = Deliver{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deliver) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deliver) ProtoMessage() {}

func (x *Deliver) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[2]
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
	return file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescGZIP(), []int{2}
}

func (x *Deliver) GetSlot() *common.Slot {
	if x != nil {
		return x.Slot
	}
	return nil
}

type FreeSlot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueueSlot uint64 `protobuf:"varint,1,opt,name=queue_slot,json=queueSlot,proto3" json:"queue_slot,omitempty"`
}

func (x *FreeSlot) Reset() {
	*x = FreeSlot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FreeSlot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FreeSlot) ProtoMessage() {}

func (x *FreeSlot) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FreeSlot.ProtoReflect.Descriptor instead.
func (*FreeSlot) Descriptor() ([]byte, []int) {
	return file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescGZIP(), []int{3}
}

func (x *FreeSlot) GetQueueSlot() uint64 {
	if x != nil {
		return x.QueueSlot
	}
	return 0
}

type PastVcbFinal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueueSlot uint64               `protobuf:"varint,1,opt,name=queue_slot,json=queueSlot,proto3" json:"queue_slot,omitempty"`
	Txs       []*requestpb.Request `protobuf:"bytes,2,rep,name=txs,proto3" json:"txs,omitempty"`
	Signature []byte               `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *PastVcbFinal) Reset() {
	*x = PastVcbFinal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PastVcbFinal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PastVcbFinal) ProtoMessage() {}

func (x *PastVcbFinal) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PastVcbFinal.ProtoReflect.Descriptor instead.
func (*PastVcbFinal) Descriptor() ([]byte, []int) {
	return file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescGZIP(), []int{4}
}

func (x *PastVcbFinal) GetQueueSlot() uint64 {
	if x != nil {
		return x.QueueSlot
	}
	return 0
}

func (x *PastVcbFinal) GetTxs() []*requestpb.Request {
	if x != nil {
		return x.Txs
	}
	return nil
}

func (x *PastVcbFinal) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type BcStarted struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Slot *common.Slot `protobuf:"bytes,1,opt,name=slot,proto3" json:"slot,omitempty"`
}

func (x *BcStarted) Reset() {
	*x = BcStarted{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BcStarted) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BcStarted) ProtoMessage() {}

func (x *BcStarted) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BcStarted.ProtoReflect.Descriptor instead.
func (*BcStarted) Descriptor() ([]byte, []int) {
	return file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescGZIP(), []int{5}
}

func (x *BcStarted) GetSlot() *common.Slot {
	if x != nil {
		return x.Slot
	}
	return nil
}

var File_aleapb_bcqueuepb_bcqueuepb_proto protoreflect.FileDescriptor

var file_aleapb_bcqueuepb_bcqueuepb_proto_rawDesc = []byte{
	0x0a, 0x20, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f, 0x62, 0x63, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x70, 0x62, 0x2f, 0x62, 0x63, 0x71, 0x75, 0x65, 0x75, 0x65, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x10, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x62, 0x63, 0x71, 0x75, 0x65,
	0x75, 0x65, 0x70, 0x62, 0x1a, 0x1a, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x19, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2f, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72,
	0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd4, 0x02, 0x0a, 0x05, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x3f, 0x0a, 0x0b, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70,
	0x62, 0x2e, 0x62, 0x63, 0x71, 0x75, 0x65, 0x75, 0x65, 0x70, 0x62, 0x2e, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x35, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x62,
	0x63, 0x71, 0x75, 0x65, 0x75, 0x65, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72,
	0x48, 0x00, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x09, 0x66,
	0x72, 0x65, 0x65, 0x5f, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x62, 0x63, 0x71, 0x75, 0x65, 0x75, 0x65, 0x70,
	0x62, 0x2e, 0x46, 0x72, 0x65, 0x65, 0x53, 0x6c, 0x6f, 0x74, 0x48, 0x00, 0x52, 0x08, 0x66, 0x72,
	0x65, 0x65, 0x53, 0x6c, 0x6f, 0x74, 0x12, 0x46, 0x0a, 0x0e, 0x70, 0x61, 0x73, 0x74, 0x5f, 0x76,
	0x63, 0x62, 0x5f, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x62, 0x63, 0x71, 0x75, 0x65, 0x75, 0x65, 0x70,
	0x62, 0x2e, 0x50, 0x61, 0x73, 0x74, 0x56, 0x63, 0x62, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x48, 0x00,
	0x52, 0x0c, 0x70, 0x61, 0x73, 0x74, 0x56, 0x63, 0x62, 0x46, 0x69, 0x6e, 0x61, 0x6c, 0x12, 0x3c,
	0x0a, 0x0a, 0x62, 0x63, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x62, 0x63, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x70, 0x62, 0x2e, 0x42, 0x63, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x48,
	0x00, 0x52, 0x09, 0x62, 0x63, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x3a, 0x04, 0x90, 0xa6,
	0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0x80, 0xa6, 0x1d, 0x01,
	0x22, 0x61, 0x0a, 0x0a, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x27,
	0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61,
	0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x6c, 0x6f,
	0x74, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x12, 0x24, 0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x03, 0x74, 0x78, 0x73, 0x3a, 0x04, 0x98,
	0xa6, 0x1d, 0x01, 0x22, 0x38, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x12, 0x27,
	0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61,
	0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x6c, 0x6f,
	0x74, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0x71, 0x0a,
	0x08, 0x46, 0x72, 0x65, 0x65, 0x53, 0x6c, 0x6f, 0x74, 0x12, 0x5f, 0x0a, 0x0a, 0x71, 0x75, 0x65,
	0x75, 0x65, 0x5f, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x40, 0x82,
	0xa6, 0x1d, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69,
	0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d,
	0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x6c, 0x65, 0x61, 0x2f, 0x61, 0x6c, 0x65, 0x61,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x6c, 0x6f, 0x74, 0x52,
	0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x53, 0x6c, 0x6f, 0x74, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01,
	0x22, 0xff, 0x01, 0x0a, 0x0c, 0x50, 0x61, 0x73, 0x74, 0x56, 0x63, 0x62, 0x46, 0x69, 0x6e, 0x61,
	0x6c, 0x12, 0x5f, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x73, 0x6c, 0x6f, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x40, 0x82, 0xa6, 0x1d, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61,
	0x6c, 0x65, 0x61, 0x2f, 0x61, 0x6c, 0x65, 0x61, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x53, 0x6c, 0x6f, 0x74, 0x52, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x53, 0x6c,
	0x6f, 0x74, 0x12, 0x24, 0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x52, 0x03, 0x74, 0x78, 0x73, 0x12, 0x62, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x44, 0x82, 0xa6, 0x1d,
	0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65,
	0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x63, 0x72, 0x79, 0x70, 0x74,
	0x6f, 0x2f, 0x74, 0x63, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x46, 0x75, 0x6c, 0x6c, 0x53, 0x69,
	0x67, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x3a, 0x04, 0x98, 0xa6,
	0x1d, 0x01, 0x22, 0x3a, 0x0a, 0x09, 0x42, 0x63, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x12,
	0x27, 0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x6c,
	0x6f, 0x74, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x42, 0x39,
	0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c,
	0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69,
	0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f,
	0x62, 0x63, 0x71, 0x75, 0x65, 0x75, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescOnce sync.Once
	file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescData = file_aleapb_bcqueuepb_bcqueuepb_proto_rawDesc
)

func file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescGZIP() []byte {
	file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescOnce.Do(func() {
		file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescData = protoimpl.X.CompressGZIP(file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescData)
	})
	return file_aleapb_bcqueuepb_bcqueuepb_proto_rawDescData
}

var file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_aleapb_bcqueuepb_bcqueuepb_proto_goTypes = []interface{}{
	(*Event)(nil),             // 0: aleapb.bcqueuepb.Event
	(*InputValue)(nil),        // 1: aleapb.bcqueuepb.InputValue
	(*Deliver)(nil),           // 2: aleapb.bcqueuepb.Deliver
	(*FreeSlot)(nil),          // 3: aleapb.bcqueuepb.FreeSlot
	(*PastVcbFinal)(nil),      // 4: aleapb.bcqueuepb.PastVcbFinal
	(*BcStarted)(nil),         // 5: aleapb.bcqueuepb.BcStarted
	(*common.Slot)(nil),       // 6: aleapb.common.Slot
	(*requestpb.Request)(nil), // 7: requestpb.Request
}
var file_aleapb_bcqueuepb_bcqueuepb_proto_depIdxs = []int32{
	1,  // 0: aleapb.bcqueuepb.Event.input_value:type_name -> aleapb.bcqueuepb.InputValue
	2,  // 1: aleapb.bcqueuepb.Event.deliver:type_name -> aleapb.bcqueuepb.Deliver
	3,  // 2: aleapb.bcqueuepb.Event.free_slot:type_name -> aleapb.bcqueuepb.FreeSlot
	4,  // 3: aleapb.bcqueuepb.Event.past_vcb_final:type_name -> aleapb.bcqueuepb.PastVcbFinal
	5,  // 4: aleapb.bcqueuepb.Event.bc_started:type_name -> aleapb.bcqueuepb.BcStarted
	6,  // 5: aleapb.bcqueuepb.InputValue.slot:type_name -> aleapb.common.Slot
	7,  // 6: aleapb.bcqueuepb.InputValue.txs:type_name -> requestpb.Request
	6,  // 7: aleapb.bcqueuepb.Deliver.slot:type_name -> aleapb.common.Slot
	7,  // 8: aleapb.bcqueuepb.PastVcbFinal.txs:type_name -> requestpb.Request
	6,  // 9: aleapb.bcqueuepb.BcStarted.slot:type_name -> aleapb.common.Slot
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_aleapb_bcqueuepb_bcqueuepb_proto_init() }
func file_aleapb_bcqueuepb_bcqueuepb_proto_init() {
	if File_aleapb_bcqueuepb_bcqueuepb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FreeSlot); i {
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
		file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PastVcbFinal); i {
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
		file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BcStarted); i {
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
	file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Event_InputValue)(nil),
		(*Event_Deliver)(nil),
		(*Event_FreeSlot)(nil),
		(*Event_PastVcbFinal)(nil),
		(*Event_BcStarted)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aleapb_bcqueuepb_bcqueuepb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aleapb_bcqueuepb_bcqueuepb_proto_goTypes,
		DependencyIndexes: file_aleapb_bcqueuepb_bcqueuepb_proto_depIdxs,
		MessageInfos:      file_aleapb_bcqueuepb_bcqueuepb_proto_msgTypes,
	}.Build()
	File_aleapb_bcqueuepb_bcqueuepb_proto = out.File
	file_aleapb_bcqueuepb_bcqueuepb_proto_rawDesc = nil
	file_aleapb_bcqueuepb_bcqueuepb_proto_goTypes = nil
	file_aleapb_bcqueuepb_bcqueuepb_proto_depIdxs = nil
}
