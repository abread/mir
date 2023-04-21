// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: factorypb/factorypb.proto

package factorypb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	mscpb "github.com/filecoin-project/mir/pkg/pb/availabilitypb/mscpb"
	checkpointpb "github.com/filecoin-project/mir/pkg/pb/checkpointpb"
	_ "github.com/filecoin-project/mir/pkg/pb/mir"
	ordererpb "github.com/filecoin-project/mir/pkg/pb/ordererpb"
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
	//	*Event_NewModule
	//	*Event_GarbageCollect
	Type isEvent_Type `protobuf_oneof:"type"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_factorypb_factorypb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_factorypb_factorypb_proto_msgTypes[0]
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
	return file_factorypb_factorypb_proto_rawDescGZIP(), []int{0}
}

func (m *Event) GetType() isEvent_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Event) GetNewModule() *NewModule {
	if x, ok := x.GetType().(*Event_NewModule); ok {
		return x.NewModule
	}
	return nil
}

func (x *Event) GetGarbageCollect() *GarbageCollect {
	if x, ok := x.GetType().(*Event_GarbageCollect); ok {
		return x.GarbageCollect
	}
	return nil
}

type isEvent_Type interface {
	isEvent_Type()
}

type Event_NewModule struct {
	NewModule *NewModule `protobuf:"bytes,1,opt,name=new_module,json=newModule,proto3,oneof"`
}

type Event_GarbageCollect struct {
	GarbageCollect *GarbageCollect `protobuf:"bytes,2,opt,name=garbage_collect,json=garbageCollect,proto3,oneof"`
}

func (*Event_NewModule) isEvent_Type() {}

func (*Event_GarbageCollect) isEvent_Type() {}

// NewModule is a request to create a new module.
// module_id must be prefixed with the factory module's ID followed by "/".
// retention_index is used for garbage collection (see: GarbageCollect).
type NewModule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ModuleId       string           `protobuf:"bytes,1,opt,name=module_id,json=moduleId,proto3" json:"module_id,omitempty"`
	RetentionIndex uint64           `protobuf:"varint,2,opt,name=retention_index,json=retentionIndex,proto3" json:"retention_index,omitempty"`
	Params         *GeneratorParams `protobuf:"bytes,3,opt,name=params,proto3" json:"params,omitempty"`
}

func (x *NewModule) Reset() {
	*x = NewModule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_factorypb_factorypb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewModule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewModule) ProtoMessage() {}

func (x *NewModule) ProtoReflect() protoreflect.Message {
	mi := &file_factorypb_factorypb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewModule.ProtoReflect.Descriptor instead.
func (*NewModule) Descriptor() ([]byte, []int) {
	return file_factorypb_factorypb_proto_rawDescGZIP(), []int{1}
}

func (x *NewModule) GetModuleId() string {
	if x != nil {
		return x.ModuleId
	}
	return ""
}

func (x *NewModule) GetRetentionIndex() uint64 {
	if x != nil {
		return x.RetentionIndex
	}
	return 0
}

func (x *NewModule) GetParams() *GeneratorParams {
	if x != nil {
		return x.Params
	}
	return nil
}

// GarbageCollect is a request to destroy all modules that have retention index smaller than the given index.
type GarbageCollect struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RetentionIndex uint64 `protobuf:"varint,1,opt,name=retention_index,json=retentionIndex,proto3" json:"retention_index,omitempty"`
}

func (x *GarbageCollect) Reset() {
	*x = GarbageCollect{}
	if protoimpl.UnsafeEnabled {
		mi := &file_factorypb_factorypb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GarbageCollect) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GarbageCollect) ProtoMessage() {}

func (x *GarbageCollect) ProtoReflect() protoreflect.Message {
	mi := &file_factorypb_factorypb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GarbageCollect.ProtoReflect.Descriptor instead.
func (*GarbageCollect) Descriptor() ([]byte, []int) {
	return file_factorypb_factorypb_proto_rawDescGZIP(), []int{2}
}

func (x *GarbageCollect) GetRetentionIndex() uint64 {
	if x != nil {
		return x.RetentionIndex
	}
	return 0
}

// GeneratorParams is a union of all possible types of module parameters
// for modules that can be created by the factory module.
type GeneratorParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*GeneratorParams_MultisigCollector
	//	*GeneratorParams_Checkpoint
	//	*GeneratorParams_EchoTestModule
	//	*GeneratorParams_PbftModule
	Type isGeneratorParams_Type `protobuf_oneof:"type"`
}

func (x *GeneratorParams) Reset() {
	*x = GeneratorParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_factorypb_factorypb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeneratorParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeneratorParams) ProtoMessage() {}

func (x *GeneratorParams) ProtoReflect() protoreflect.Message {
	mi := &file_factorypb_factorypb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeneratorParams.ProtoReflect.Descriptor instead.
func (*GeneratorParams) Descriptor() ([]byte, []int) {
	return file_factorypb_factorypb_proto_rawDescGZIP(), []int{3}
}

func (m *GeneratorParams) GetType() isGeneratorParams_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *GeneratorParams) GetMultisigCollector() *mscpb.InstanceParams {
	if x, ok := x.GetType().(*GeneratorParams_MultisigCollector); ok {
		return x.MultisigCollector
	}
	return nil
}

func (x *GeneratorParams) GetCheckpoint() *checkpointpb.InstanceParams {
	if x, ok := x.GetType().(*GeneratorParams_Checkpoint); ok {
		return x.Checkpoint
	}
	return nil
}

func (x *GeneratorParams) GetEchoTestModule() *EchoModuleParams {
	if x, ok := x.GetType().(*GeneratorParams_EchoTestModule); ok {
		return x.EchoTestModule
	}
	return nil
}

func (x *GeneratorParams) GetPbftModule() *ordererpb.PBFTModule {
	if x, ok := x.GetType().(*GeneratorParams_PbftModule); ok {
		return x.PbftModule
	}
	return nil
}

type isGeneratorParams_Type interface {
	isGeneratorParams_Type()
}

type GeneratorParams_MultisigCollector struct {
	MultisigCollector *mscpb.InstanceParams `protobuf:"bytes,1,opt,name=multisig_collector,json=multisigCollector,proto3,oneof"`
}

type GeneratorParams_Checkpoint struct {
	Checkpoint *checkpointpb.InstanceParams `protobuf:"bytes,2,opt,name=checkpoint,proto3,oneof"`
}

type GeneratorParams_EchoTestModule struct {
	EchoTestModule *EchoModuleParams `protobuf:"bytes,3,opt,name=echo_test_module,json=echoTestModule,proto3,oneof"`
}

type GeneratorParams_PbftModule struct {
	PbftModule *ordererpb.PBFTModule `protobuf:"bytes,4,opt,name=pbft_module,json=pbftModule,proto3,oneof"`
}

func (*GeneratorParams_MultisigCollector) isGeneratorParams_Type() {}

func (*GeneratorParams_Checkpoint) isGeneratorParams_Type() {}

func (*GeneratorParams_EchoTestModule) isGeneratorParams_Type() {}

func (*GeneratorParams_PbftModule) isGeneratorParams_Type() {}

// Used only for unit tests.
type EchoModuleParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prefix string `protobuf:"bytes,1,opt,name=prefix,proto3" json:"prefix,omitempty"` // This prefix is prepended to all strings the module echoes.
}

func (x *EchoModuleParams) Reset() {
	*x = EchoModuleParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_factorypb_factorypb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoModuleParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoModuleParams) ProtoMessage() {}

func (x *EchoModuleParams) ProtoReflect() protoreflect.Message {
	mi := &file_factorypb_factorypb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoModuleParams.ProtoReflect.Descriptor instead.
func (*EchoModuleParams) Descriptor() ([]byte, []int) {
	return file_factorypb_factorypb_proto_rawDescGZIP(), []int{4}
}

func (x *EchoModuleParams) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

var File_factorypb_factorypb_proto protoreflect.FileDescriptor

var file_factorypb_factorypb_proto_rawDesc = []byte{
	0x0a, 0x19, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x70, 0x62, 0x2f, 0x66, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x79, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x66, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x79, 0x70, 0x62, 0x1a, 0x20, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x79, 0x70, 0x62, 0x2f, 0x6d, 0x73, 0x63, 0x70, 0x62, 0x2f, 0x6d, 0x73, 0x63,
	0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x70, 0x62, 0x2f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x65, 0x72, 0x70, 0x62, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65,
	0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x98, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x35, 0x0a, 0x0a,
	0x6e, 0x65, 0x77, 0x5f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x70, 0x62, 0x2e, 0x4e, 0x65, 0x77,
	0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x09, 0x6e, 0x65, 0x77, 0x4d, 0x6f, 0x64,
	0x75, 0x6c, 0x65, 0x12, 0x44, 0x0a, 0x0f, 0x67, 0x61, 0x72, 0x62, 0x61, 0x67, 0x65, 0x5f, 0x63,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x66,
	0x61, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x70, 0x62, 0x2e, 0x47, 0x61, 0x72, 0x62, 0x61, 0x67, 0x65,
	0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x48, 0x00, 0x52, 0x0e, 0x67, 0x61, 0x72, 0x62, 0x61,
	0x67, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x3a, 0x04, 0x90, 0xa6, 0x1d, 0x01, 0x42,
	0x0c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0x89, 0x02,
	0x0a, 0x09, 0x4e, 0x65, 0x77, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x53, 0x0a, 0x09, 0x6d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x36,
	0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x08, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x49, 0x64,
	0x12, 0x6d, 0x0a, 0x0f, 0x72, 0x65, 0x74, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x42, 0x44, 0x82, 0xa6, 0x1d, 0x40, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f,
	0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x74, 0x6f, 0x72, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x52, 0x65, 0x74, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52,
	0x0e, 0x72, 0x65, 0x74, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12,
	0x32, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x52, 0x06, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0x85, 0x01, 0x0a, 0x0e, 0x47, 0x61,
	0x72, 0x62, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x12, 0x6d, 0x0a, 0x0f,
	0x72, 0x65, 0x74, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x44, 0x82, 0xa6, 0x1d, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74,
	0x72, 0x61, 0x6e, 0x74, 0x6f, 0x72, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x74,
	0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x0e, 0x72, 0x65, 0x74,
	0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x3a, 0x04, 0x98, 0xa6, 0x1d,
	0x01, 0x22, 0xb9, 0x02, 0x0a, 0x0f, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x55, 0x0a, 0x12, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x73, 0x69,
	0x67, 0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x24, 0x2e, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x70, 0x62, 0x2e, 0x6d, 0x73, 0x63, 0x70, 0x62, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x48, 0x00, 0x52, 0x11, 0x6d, 0x75, 0x6c, 0x74, 0x69,
	0x73, 0x69, 0x67, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x3e, 0x0a, 0x0a,
	0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x70, 0x62, 0x2e,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x48, 0x00,
	0x52, 0x0a, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x47, 0x0a, 0x10,
	0x65, 0x63, 0x68, 0x6f, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x70, 0x62, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x48, 0x00, 0x52, 0x0e, 0x65, 0x63, 0x68, 0x6f, 0x54, 0x65, 0x73, 0x74, 0x4d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x38, 0x0a, 0x0b, 0x70, 0x62, 0x66, 0x74, 0x5f, 0x6d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x50, 0x42, 0x46, 0x54, 0x4d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x62, 0x66, 0x74, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x3a,
	0x04, 0x80, 0xa6, 0x1d, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x30, 0x0a,
	0x10, 0x45, 0x63, 0x68, 0x6f, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x3a, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x42,
	0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69,
	0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d,
	0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72,
	0x79, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_factorypb_factorypb_proto_rawDescOnce sync.Once
	file_factorypb_factorypb_proto_rawDescData = file_factorypb_factorypb_proto_rawDesc
)

func file_factorypb_factorypb_proto_rawDescGZIP() []byte {
	file_factorypb_factorypb_proto_rawDescOnce.Do(func() {
		file_factorypb_factorypb_proto_rawDescData = protoimpl.X.CompressGZIP(file_factorypb_factorypb_proto_rawDescData)
	})
	return file_factorypb_factorypb_proto_rawDescData
}

var file_factorypb_factorypb_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_factorypb_factorypb_proto_goTypes = []interface{}{
	(*Event)(nil),                       // 0: factorypb.Event
	(*NewModule)(nil),                   // 1: factorypb.NewModule
	(*GarbageCollect)(nil),              // 2: factorypb.GarbageCollect
	(*GeneratorParams)(nil),             // 3: factorypb.GeneratorParams
	(*EchoModuleParams)(nil),            // 4: factorypb.EchoModuleParams
	(*mscpb.InstanceParams)(nil),        // 5: availabilitypb.mscpb.InstanceParams
	(*checkpointpb.InstanceParams)(nil), // 6: checkpointpb.InstanceParams
	(*ordererpb.PBFTModule)(nil),        // 7: ordererpb.PBFTModule
}
var file_factorypb_factorypb_proto_depIdxs = []int32{
	1, // 0: factorypb.Event.new_module:type_name -> factorypb.NewModule
	2, // 1: factorypb.Event.garbage_collect:type_name -> factorypb.GarbageCollect
	3, // 2: factorypb.NewModule.params:type_name -> factorypb.GeneratorParams
	5, // 3: factorypb.GeneratorParams.multisig_collector:type_name -> availabilitypb.mscpb.InstanceParams
	6, // 4: factorypb.GeneratorParams.checkpoint:type_name -> checkpointpb.InstanceParams
	4, // 5: factorypb.GeneratorParams.echo_test_module:type_name -> factorypb.EchoModuleParams
	7, // 6: factorypb.GeneratorParams.pbft_module:type_name -> ordererpb.PBFTModule
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_factorypb_factorypb_proto_init() }
func file_factorypb_factorypb_proto_init() {
	if File_factorypb_factorypb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_factorypb_factorypb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_factorypb_factorypb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewModule); i {
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
		file_factorypb_factorypb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GarbageCollect); i {
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
		file_factorypb_factorypb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GeneratorParams); i {
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
		file_factorypb_factorypb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoModuleParams); i {
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
	file_factorypb_factorypb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Event_NewModule)(nil),
		(*Event_GarbageCollect)(nil),
	}
	file_factorypb_factorypb_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*GeneratorParams_MultisigCollector)(nil),
		(*GeneratorParams_Checkpoint)(nil),
		(*GeneratorParams_EchoTestModule)(nil),
		(*GeneratorParams_PbftModule)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_factorypb_factorypb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_factorypb_factorypb_proto_goTypes,
		DependencyIndexes: file_factorypb_factorypb_proto_depIdxs,
		MessageInfos:      file_factorypb_factorypb_proto_msgTypes,
	}.Build()
	File_factorypb_factorypb_proto = out.File
	file_factorypb_factorypb_proto_rawDesc = nil
	file_factorypb_factorypb_proto_goTypes = nil
	file_factorypb_factorypb_proto_depIdxs = nil
}
