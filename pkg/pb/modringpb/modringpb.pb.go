// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: modringpb/modringpb.proto

package modringpb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
	messagepb "github.com/filecoin-project/mir/pkg/pb/messagepb"
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
	//	*Event_Free
	//	*Event_Freed
	//	*Event_PastMessagesRecvd
	Type isEvent_Type `protobuf_oneof:"type"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modringpb_modringpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_modringpb_modringpb_proto_msgTypes[0]
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
	return file_modringpb_modringpb_proto_rawDescGZIP(), []int{0}
}

func (m *Event) GetType() isEvent_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Event) GetFree() *FreeSubmodule {
	if x, ok := x.GetType().(*Event_Free); ok {
		return x.Free
	}
	return nil
}

func (x *Event) GetFreed() *FreedSubmodule {
	if x, ok := x.GetType().(*Event_Freed); ok {
		return x.Freed
	}
	return nil
}

func (x *Event) GetPastMessagesRecvd() *PastMessagesRecvd {
	if x, ok := x.GetType().(*Event_PastMessagesRecvd); ok {
		return x.PastMessagesRecvd
	}
	return nil
}

type isEvent_Type interface {
	isEvent_Type()
}

type Event_Free struct {
	Free *FreeSubmodule `protobuf:"bytes,1,opt,name=free,proto3,oneof"`
}

type Event_Freed struct {
	Freed *FreedSubmodule `protobuf:"bytes,2,opt,name=freed,proto3,oneof"`
}

type Event_PastMessagesRecvd struct {
	PastMessagesRecvd *PastMessagesRecvd `protobuf:"bytes,3,opt,name=past_messages_recvd,json=pastMessagesRecvd,proto3,oneof"`
}

func (*Event_Free) isEvent_Type() {}

func (*Event_Freed) isEvent_Type() {}

func (*Event_PastMessagesRecvd) isEvent_Type() {}

type FreeSubmodule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Origin *FreeSubmoduleOrigin `protobuf:"bytes,2,opt,name=origin,proto3" json:"origin,omitempty"`
}

func (x *FreeSubmodule) Reset() {
	*x = FreeSubmodule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modringpb_modringpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FreeSubmodule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FreeSubmodule) ProtoMessage() {}

func (x *FreeSubmodule) ProtoReflect() protoreflect.Message {
	mi := &file_modringpb_modringpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FreeSubmodule.ProtoReflect.Descriptor instead.
func (*FreeSubmodule) Descriptor() ([]byte, []int) {
	return file_modringpb_modringpb_proto_rawDescGZIP(), []int{1}
}

func (x *FreeSubmodule) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *FreeSubmodule) GetOrigin() *FreeSubmoduleOrigin {
	if x != nil {
		return x.Origin
	}
	return nil
}

type FreedSubmodule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Origin *FreeSubmoduleOrigin `protobuf:"bytes,1,opt,name=origin,proto3" json:"origin,omitempty"`
}

func (x *FreedSubmodule) Reset() {
	*x = FreedSubmodule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modringpb_modringpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FreedSubmodule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FreedSubmodule) ProtoMessage() {}

func (x *FreedSubmodule) ProtoReflect() protoreflect.Message {
	mi := &file_modringpb_modringpb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FreedSubmodule.ProtoReflect.Descriptor instead.
func (*FreedSubmodule) Descriptor() ([]byte, []int) {
	return file_modringpb_modringpb_proto_rawDescGZIP(), []int{2}
}

func (x *FreedSubmodule) GetOrigin() *FreeSubmoduleOrigin {
	if x != nil {
		return x.Origin
	}
	return nil
}

type PastMessagesRecvd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []*PastMessage `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *PastMessagesRecvd) Reset() {
	*x = PastMessagesRecvd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modringpb_modringpb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PastMessagesRecvd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PastMessagesRecvd) ProtoMessage() {}

func (x *PastMessagesRecvd) ProtoReflect() protoreflect.Message {
	mi := &file_modringpb_modringpb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PastMessagesRecvd.ProtoReflect.Descriptor instead.
func (*PastMessagesRecvd) Descriptor() ([]byte, []int) {
	return file_modringpb_modringpb_proto_rawDescGZIP(), []int{3}
}

func (x *PastMessagesRecvd) GetMessages() []*PastMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

type PastMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DestId  uint64             `protobuf:"varint,1,opt,name=dest_id,json=destId,proto3" json:"dest_id,omitempty"`
	From    string             `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Message *messagepb.Message `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PastMessage) Reset() {
	*x = PastMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modringpb_modringpb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PastMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PastMessage) ProtoMessage() {}

func (x *PastMessage) ProtoReflect() protoreflect.Message {
	mi := &file_modringpb_modringpb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PastMessage.ProtoReflect.Descriptor instead.
func (*PastMessage) Descriptor() ([]byte, []int) {
	return file_modringpb_modringpb_proto_rawDescGZIP(), []int{4}
}

func (x *PastMessage) GetDestId() uint64 {
	if x != nil {
		return x.DestId
	}
	return 0
}

func (x *PastMessage) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *PastMessage) GetMessage() *messagepb.Message {
	if x != nil {
		return x.Message
	}
	return nil
}

// ============================================================
// Data structures
// ============================================================
type FreeSubmoduleOrigin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Module string `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
	// Types that are assignable to Type:
	//
	//	*FreeSubmoduleOrigin_ContextStore
	//	*FreeSubmoduleOrigin_Dsl
	Type isFreeSubmoduleOrigin_Type `protobuf_oneof:"Type"`
}

func (x *FreeSubmoduleOrigin) Reset() {
	*x = FreeSubmoduleOrigin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modringpb_modringpb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FreeSubmoduleOrigin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FreeSubmoduleOrigin) ProtoMessage() {}

func (x *FreeSubmoduleOrigin) ProtoReflect() protoreflect.Message {
	mi := &file_modringpb_modringpb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FreeSubmoduleOrigin.ProtoReflect.Descriptor instead.
func (*FreeSubmoduleOrigin) Descriptor() ([]byte, []int) {
	return file_modringpb_modringpb_proto_rawDescGZIP(), []int{5}
}

func (x *FreeSubmoduleOrigin) GetModule() string {
	if x != nil {
		return x.Module
	}
	return ""
}

func (m *FreeSubmoduleOrigin) GetType() isFreeSubmoduleOrigin_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *FreeSubmoduleOrigin) GetContextStore() *contextstorepb.Origin {
	if x, ok := x.GetType().(*FreeSubmoduleOrigin_ContextStore); ok {
		return x.ContextStore
	}
	return nil
}

func (x *FreeSubmoduleOrigin) GetDsl() *dslpb.Origin {
	if x, ok := x.GetType().(*FreeSubmoduleOrigin_Dsl); ok {
		return x.Dsl
	}
	return nil
}

type isFreeSubmoduleOrigin_Type interface {
	isFreeSubmoduleOrigin_Type()
}

type FreeSubmoduleOrigin_ContextStore struct {
	ContextStore *contextstorepb.Origin `protobuf:"bytes,2,opt,name=context_store,json=contextStore,proto3,oneof"`
}

type FreeSubmoduleOrigin_Dsl struct {
	Dsl *dslpb.Origin `protobuf:"bytes,3,opt,name=dsl,proto3,oneof"`
}

func (*FreeSubmoduleOrigin_ContextStore) isFreeSubmoduleOrigin_Type() {}

func (*FreeSubmoduleOrigin_Dsl) isFreeSubmoduleOrigin_Type() {}

var File_modringpb_modringpb_proto protoreflect.FileDescriptor

var file_modringpb_modringpb_proto_rawDesc = []byte{
	0x0a, 0x19, 0x6d, 0x6f, 0x64, 0x72, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2f, 0x6d, 0x6f, 0x64, 0x72,
	0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d, 0x6f, 0x64,
	0x72, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x1a, 0x19, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x70,
	0x62, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x23, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x70,
	0x62, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x70, 0x62,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x64, 0x73, 0x6c, 0x70, 0x62, 0x2f, 0x64, 0x73,
	0x6c, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72, 0x2f, 0x63,
	0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xce, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x2e, 0x0a, 0x04, 0x66, 0x72, 0x65, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x6d, 0x6f, 0x64, 0x72, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x46, 0x72, 0x65, 0x65,
	0x53, 0x75, 0x62, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x04, 0x66, 0x72, 0x65,
	0x65, 0x12, 0x31, 0x0a, 0x05, 0x66, 0x72, 0x65, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x6d, 0x6f, 0x64, 0x72, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x46, 0x72, 0x65,
	0x65, 0x64, 0x53, 0x75, 0x62, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x05, 0x66,
	0x72, 0x65, 0x65, 0x64, 0x12, 0x4e, 0x0a, 0x13, 0x70, 0x61, 0x73, 0x74, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x5f, 0x72, 0x65, 0x63, 0x76, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x6d, 0x6f, 0x64, 0x72, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x50, 0x61,
	0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x63, 0x76, 0x64, 0x48,
	0x00, 0x52, 0x11, 0x70, 0x61, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52,
	0x65, 0x63, 0x76, 0x64, 0x3a, 0x04, 0x90, 0xa6, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0x63, 0x0a, 0x0d, 0x46, 0x72, 0x65, 0x65,
	0x53, 0x75, 0x62, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3c, 0x0a, 0x06, 0x6f, 0x72, 0x69,
	0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6d, 0x6f, 0x64, 0x72,
	0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x46, 0x72, 0x65, 0x65, 0x53, 0x75, 0x62, 0x6d, 0x6f, 0x64,
	0x75, 0x6c, 0x65, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x42, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x52,
	0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0x54, 0x0a,
	0x0e, 0x46, 0x72, 0x65, 0x65, 0x64, 0x53, 0x75, 0x62, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x12,
	0x3c, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x6d, 0x6f, 0x64, 0x72, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x46, 0x72, 0x65, 0x65,
	0x53, 0x75, 0x62, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x42,
	0x04, 0xa0, 0xa6, 0x1d, 0x01, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x3a, 0x04, 0x98,
	0xa6, 0x1d, 0x01, 0x22, 0x4d, 0x0a, 0x11, 0x50, 0x61, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x52, 0x65, 0x63, 0x76, 0x64, 0x12, 0x32, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x6f, 0x64,
	0x72, 0x69, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x50, 0x61, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x3a, 0x04, 0x98, 0xa6,
	0x1d, 0x01, 0x22, 0xa4, 0x01, 0x0a, 0x0b, 0x50, 0x61, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x64, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x64, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x48, 0x0a, 0x04, 0x66,
	0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x34, 0x82, 0xa6, 0x1d, 0x30, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f,
	0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x52,
	0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x2c, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x3a, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0xd5, 0x01, 0x0a, 0x13, 0x46, 0x72,
	0x65, 0x65, 0x53, 0x75, 0x62, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x4f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x12, 0x4e, 0x0a, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x36, 0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x12, 0x3d, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x78, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x48, 0x00, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x12, 0x21, 0x0a, 0x03, 0x64, 0x73, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x64, 0x73, 0x6c, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x48, 0x00, 0x52, 0x03,
	0x64, 0x73, 0x6c, 0x3a, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x6d, 0x6f, 0x64, 0x72,
	0x69, 0x6e, 0x67, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_modringpb_modringpb_proto_rawDescOnce sync.Once
	file_modringpb_modringpb_proto_rawDescData = file_modringpb_modringpb_proto_rawDesc
)

func file_modringpb_modringpb_proto_rawDescGZIP() []byte {
	file_modringpb_modringpb_proto_rawDescOnce.Do(func() {
		file_modringpb_modringpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_modringpb_modringpb_proto_rawDescData)
	})
	return file_modringpb_modringpb_proto_rawDescData
}

var file_modringpb_modringpb_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_modringpb_modringpb_proto_goTypes = []interface{}{
	(*Event)(nil),                 // 0: modringpb.Event
	(*FreeSubmodule)(nil),         // 1: modringpb.FreeSubmodule
	(*FreedSubmodule)(nil),        // 2: modringpb.FreedSubmodule
	(*PastMessagesRecvd)(nil),     // 3: modringpb.PastMessagesRecvd
	(*PastMessage)(nil),           // 4: modringpb.PastMessage
	(*FreeSubmoduleOrigin)(nil),   // 5: modringpb.FreeSubmoduleOrigin
	(*messagepb.Message)(nil),     // 6: messagepb.Message
	(*contextstorepb.Origin)(nil), // 7: contextstorepb.Origin
	(*dslpb.Origin)(nil),          // 8: dslpb.Origin
}
var file_modringpb_modringpb_proto_depIdxs = []int32{
	1, // 0: modringpb.Event.free:type_name -> modringpb.FreeSubmodule
	2, // 1: modringpb.Event.freed:type_name -> modringpb.FreedSubmodule
	3, // 2: modringpb.Event.past_messages_recvd:type_name -> modringpb.PastMessagesRecvd
	5, // 3: modringpb.FreeSubmodule.origin:type_name -> modringpb.FreeSubmoduleOrigin
	5, // 4: modringpb.FreedSubmodule.origin:type_name -> modringpb.FreeSubmoduleOrigin
	4, // 5: modringpb.PastMessagesRecvd.messages:type_name -> modringpb.PastMessage
	6, // 6: modringpb.PastMessage.message:type_name -> messagepb.Message
	7, // 7: modringpb.FreeSubmoduleOrigin.context_store:type_name -> contextstorepb.Origin
	8, // 8: modringpb.FreeSubmoduleOrigin.dsl:type_name -> dslpb.Origin
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_modringpb_modringpb_proto_init() }
func file_modringpb_modringpb_proto_init() {
	if File_modringpb_modringpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_modringpb_modringpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_modringpb_modringpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FreeSubmodule); i {
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
		file_modringpb_modringpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FreedSubmodule); i {
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
		file_modringpb_modringpb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PastMessagesRecvd); i {
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
		file_modringpb_modringpb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PastMessage); i {
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
		file_modringpb_modringpb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FreeSubmoduleOrigin); i {
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
	file_modringpb_modringpb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Event_Free)(nil),
		(*Event_Freed)(nil),
		(*Event_PastMessagesRecvd)(nil),
	}
	file_modringpb_modringpb_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*FreeSubmoduleOrigin_ContextStore)(nil),
		(*FreeSubmoduleOrigin_Dsl)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_modringpb_modringpb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_modringpb_modringpb_proto_goTypes,
		DependencyIndexes: file_modringpb_modringpb_proto_depIdxs,
		MessageInfos:      file_modringpb_modringpb_proto_msgTypes,
	}.Build()
	File_modringpb_modringpb_proto = out.File
	file_modringpb_modringpb_proto_rawDesc = nil
	file_modringpb_modringpb_proto_goTypes = nil
	file_modringpb_modringpb_proto_depIdxs = nil
}
