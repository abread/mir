// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: availabilitypb/batchdbpb/batchdbpb.proto

package batchdbpb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
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
	//	*Event_Lookup
	//	*Event_LookupResponse
	//	*Event_Store
	//	*Event_Stored
	Type isEvent_Type `protobuf_oneof:"Type"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[0]
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
	return file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescGZIP(), []int{0}
}

func (m *Event) GetType() isEvent_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Event) GetLookup() *LookupBatch {
	if x, ok := x.GetType().(*Event_Lookup); ok {
		return x.Lookup
	}
	return nil
}

func (x *Event) GetLookupResponse() *LookupBatchResponse {
	if x, ok := x.GetType().(*Event_LookupResponse); ok {
		return x.LookupResponse
	}
	return nil
}

func (x *Event) GetStore() *StoreBatch {
	if x, ok := x.GetType().(*Event_Store); ok {
		return x.Store
	}
	return nil
}

func (x *Event) GetStored() *BatchStored {
	if x, ok := x.GetType().(*Event_Stored); ok {
		return x.Stored
	}
	return nil
}

type isEvent_Type interface {
	isEvent_Type()
}

type Event_Lookup struct {
	Lookup *LookupBatch `protobuf:"bytes,1,opt,name=lookup,proto3,oneof"`
}

type Event_LookupResponse struct {
	LookupResponse *LookupBatchResponse `protobuf:"bytes,2,opt,name=lookup_response,json=lookupResponse,proto3,oneof"`
}

type Event_Store struct {
	Store *StoreBatch `protobuf:"bytes,3,opt,name=store,proto3,oneof"`
}

type Event_Stored struct {
	Stored *BatchStored `protobuf:"bytes,4,opt,name=stored,proto3,oneof"`
}

func (*Event_Lookup) isEvent_Type() {}

func (*Event_LookupResponse) isEvent_Type() {}

func (*Event_Store) isEvent_Type() {}

func (*Event_Stored) isEvent_Type() {}

// LookupBatch is used to pull a batch with its metadata from the local batch database.
type LookupBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BatchId []byte             `protobuf:"bytes,1,opt,name=batch_id,json=batchId,proto3" json:"batch_id,omitempty"`
	Origin  *LookupBatchOrigin `protobuf:"bytes,2,opt,name=origin,proto3" json:"origin,omitempty"`
}

func (x *LookupBatch) Reset() {
	*x = LookupBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LookupBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupBatch) ProtoMessage() {}

func (x *LookupBatch) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupBatch.ProtoReflect.Descriptor instead.
func (*LookupBatch) Descriptor() ([]byte, []int) {
	return file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescGZIP(), []int{1}
}

func (x *LookupBatch) GetBatchId() []byte {
	if x != nil {
		return x.BatchId
	}
	return nil
}

func (x *LookupBatch) GetOrigin() *LookupBatchOrigin {
	if x != nil {
		return x.Origin
	}
	return nil
}

// LookupBatchResponse is a response to a LookupBatch event.
type LookupBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Found    bool                 `protobuf:"varint,1,opt,name=found,proto3" json:"found,omitempty"`
	Txs      []*requestpb.Request `protobuf:"bytes,2,rep,name=txs,proto3" json:"txs,omitempty"`
	Metadata []byte               `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Origin   *LookupBatchOrigin   `protobuf:"bytes,4,opt,name=origin,proto3" json:"origin,omitempty"`
}

func (x *LookupBatchResponse) Reset() {
	*x = LookupBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LookupBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupBatchResponse) ProtoMessage() {}

func (x *LookupBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupBatchResponse.ProtoReflect.Descriptor instead.
func (*LookupBatchResponse) Descriptor() ([]byte, []int) {
	return file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescGZIP(), []int{2}
}

func (x *LookupBatchResponse) GetFound() bool {
	if x != nil {
		return x.Found
	}
	return false
}

func (x *LookupBatchResponse) GetTxs() []*requestpb.Request {
	if x != nil {
		return x.Txs
	}
	return nil
}

func (x *LookupBatchResponse) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *LookupBatchResponse) GetOrigin() *LookupBatchOrigin {
	if x != nil {
		return x.Origin
	}
	return nil
}

// StoreBatch is used to store a new batch in the local batch database.
type StoreBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BatchId  []byte               `protobuf:"bytes,1,opt,name=batch_id,json=batchId,proto3" json:"batch_id,omitempty"`
	TxIds    [][]byte             `protobuf:"bytes,2,rep,name=tx_ids,json=txIds,proto3" json:"tx_ids,omitempty"`
	Txs      []*requestpb.Request `protobuf:"bytes,3,rep,name=txs,proto3" json:"txs,omitempty"`
	Metadata []byte               `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Origin   *StoreBatchOrigin    `protobuf:"bytes,5,opt,name=origin,proto3" json:"origin,omitempty"`
}

func (x *StoreBatch) Reset() {
	*x = StoreBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreBatch) ProtoMessage() {}

func (x *StoreBatch) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreBatch.ProtoReflect.Descriptor instead.
func (*StoreBatch) Descriptor() ([]byte, []int) {
	return file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescGZIP(), []int{3}
}

func (x *StoreBatch) GetBatchId() []byte {
	if x != nil {
		return x.BatchId
	}
	return nil
}

func (x *StoreBatch) GetTxIds() [][]byte {
	if x != nil {
		return x.TxIds
	}
	return nil
}

func (x *StoreBatch) GetTxs() []*requestpb.Request {
	if x != nil {
		return x.Txs
	}
	return nil
}

func (x *StoreBatch) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *StoreBatch) GetOrigin() *StoreBatchOrigin {
	if x != nil {
		return x.Origin
	}
	return nil
}

// BatchStored is a response to a VerifyCert event.
type BatchStored struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Origin *StoreBatchOrigin `protobuf:"bytes,1,opt,name=origin,proto3" json:"origin,omitempty"`
}

func (x *BatchStored) Reset() {
	*x = BatchStored{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchStored) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchStored) ProtoMessage() {}

func (x *BatchStored) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchStored.ProtoReflect.Descriptor instead.
func (*BatchStored) Descriptor() ([]byte, []int) {
	return file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescGZIP(), []int{4}
}

func (x *BatchStored) GetOrigin() *StoreBatchOrigin {
	if x != nil {
		return x.Origin
	}
	return nil
}

type LookupBatchOrigin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Module string `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
	// Types that are assignable to Type:
	//
	//	*LookupBatchOrigin_ContextStore
	//	*LookupBatchOrigin_Dsl
	Type isLookupBatchOrigin_Type `protobuf_oneof:"Type"`
}

func (x *LookupBatchOrigin) Reset() {
	*x = LookupBatchOrigin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LookupBatchOrigin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupBatchOrigin) ProtoMessage() {}

func (x *LookupBatchOrigin) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupBatchOrigin.ProtoReflect.Descriptor instead.
func (*LookupBatchOrigin) Descriptor() ([]byte, []int) {
	return file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescGZIP(), []int{5}
}

func (x *LookupBatchOrigin) GetModule() string {
	if x != nil {
		return x.Module
	}
	return ""
}

func (m *LookupBatchOrigin) GetType() isLookupBatchOrigin_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *LookupBatchOrigin) GetContextStore() *contextstorepb.Origin {
	if x, ok := x.GetType().(*LookupBatchOrigin_ContextStore); ok {
		return x.ContextStore
	}
	return nil
}

func (x *LookupBatchOrigin) GetDsl() *dslpb.Origin {
	if x, ok := x.GetType().(*LookupBatchOrigin_Dsl); ok {
		return x.Dsl
	}
	return nil
}

type isLookupBatchOrigin_Type interface {
	isLookupBatchOrigin_Type()
}

type LookupBatchOrigin_ContextStore struct {
	ContextStore *contextstorepb.Origin `protobuf:"bytes,2,opt,name=context_store,json=contextStore,proto3,oneof"`
}

type LookupBatchOrigin_Dsl struct {
	Dsl *dslpb.Origin `protobuf:"bytes,3,opt,name=dsl,proto3,oneof"`
}

func (*LookupBatchOrigin_ContextStore) isLookupBatchOrigin_Type() {}

func (*LookupBatchOrigin_Dsl) isLookupBatchOrigin_Type() {}

type StoreBatchOrigin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Module string `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
	// Types that are assignable to Type:
	//
	//	*StoreBatchOrigin_ContextStore
	//	*StoreBatchOrigin_Dsl
	Type isStoreBatchOrigin_Type `protobuf_oneof:"Type"`
}

func (x *StoreBatchOrigin) Reset() {
	*x = StoreBatchOrigin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreBatchOrigin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreBatchOrigin) ProtoMessage() {}

func (x *StoreBatchOrigin) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreBatchOrigin.ProtoReflect.Descriptor instead.
func (*StoreBatchOrigin) Descriptor() ([]byte, []int) {
	return file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescGZIP(), []int{6}
}

func (x *StoreBatchOrigin) GetModule() string {
	if x != nil {
		return x.Module
	}
	return ""
}

func (m *StoreBatchOrigin) GetType() isStoreBatchOrigin_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *StoreBatchOrigin) GetContextStore() *contextstorepb.Origin {
	if x, ok := x.GetType().(*StoreBatchOrigin_ContextStore); ok {
		return x.ContextStore
	}
	return nil
}

func (x *StoreBatchOrigin) GetDsl() *dslpb.Origin {
	if x, ok := x.GetType().(*StoreBatchOrigin_Dsl); ok {
		return x.Dsl
	}
	return nil
}

type isStoreBatchOrigin_Type interface {
	isStoreBatchOrigin_Type()
}

type StoreBatchOrigin_ContextStore struct {
	ContextStore *contextstorepb.Origin `protobuf:"bytes,2,opt,name=context_store,json=contextStore,proto3,oneof"`
}

type StoreBatchOrigin_Dsl struct {
	Dsl *dslpb.Origin `protobuf:"bytes,3,opt,name=dsl,proto3,oneof"`
}

func (*StoreBatchOrigin_ContextStore) isStoreBatchOrigin_Type() {}

func (*StoreBatchOrigin_Dsl) isStoreBatchOrigin_Type() {}

var File_availabilitypb_batchdbpb_batchdbpb_proto protoreflect.FileDescriptor

var file_availabilitypb_batchdbpb_batchdbpb_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x70, 0x62,
	0x2f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x64, 0x62, 0x70, 0x62, 0x2f, 0x62, 0x61, 0x74, 0x63, 0x68,
	0x64, 0x62, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x62, 0x61, 0x74, 0x63,
	0x68, 0x64, 0x62, 0x70, 0x62, 0x1a, 0x23, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x64, 0x73, 0x6c, 0x70,
	0x62, 0x2f, 0x64, 0x73, 0x6c, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72, 0x2f, 0x63, 0x6f,
	0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf9, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x30, 0x0a, 0x06, 0x6c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x64, 0x62, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x6f,
	0x6b, 0x75, 0x70, 0x42, 0x61, 0x74, 0x63, 0x68, 0x48, 0x00, 0x52, 0x06, 0x6c, 0x6f, 0x6f, 0x6b,
	0x75, 0x70, 0x12, 0x49, 0x0a, 0x0f, 0x6c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x5f, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x62, 0x61,
	0x74, 0x63, 0x68, 0x64, 0x62, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x0e, 0x6c,
	0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a,
	0x05, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x64, 0x62, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x48, 0x00, 0x52, 0x05, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x30, 0x0a, 0x06,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x64, 0x62, 0x70, 0x62, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74,
	0x6f, 0x72, 0x65, 0x64, 0x48, 0x00, 0x52, 0x06, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x3a, 0x04,
	0x90, 0xa6, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x04, 0x80, 0xa6,
	0x1d, 0x01, 0x22, 0xc0, 0x01, 0x0a, 0x0b, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x42, 0x61, 0x74,
	0x63, 0x68, 0x12, 0x6f, 0x0a, 0x08, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x42, 0x54, 0x82, 0xa6, 0x1d, 0x50, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x76,
	0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x69,
	0x73, 0x69, 0x67, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x49, 0x44, 0x52, 0x07, 0x62, 0x61, 0x74, 0x63,
	0x68, 0x49, 0x64, 0x12, 0x3a, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x64, 0x62, 0x70, 0x62, 0x2e,
	0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x42, 0x61, 0x74, 0x63, 0x68, 0x4f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x42, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x3a,
	0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0xaf, 0x01, 0x0a, 0x13, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x66, 0x6f,
	0x75, 0x6e, 0x64, 0x12, 0x24, 0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x52, 0x03, 0x74, 0x78, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x3a, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x64, 0x62, 0x70,
	0x62, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x42, 0x61, 0x74, 0x63, 0x68, 0x4f, 0x72, 0x69,
	0x67, 0x69, 0x6e, 0x42, 0x04, 0xa0, 0xa6, 0x1d, 0x01, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0xd3, 0x02, 0x0a, 0x0a, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x6f, 0x0a, 0x08, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x54, 0x82, 0xa6, 0x1d, 0x50, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69,
	0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2f, 0x6d,
	0x75, 0x6c, 0x74, 0x69, 0x73, 0x69, 0x67, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x49, 0x44, 0x52, 0x07,
	0x62, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x12, 0x51, 0x0a, 0x06, 0x74, 0x78, 0x5f, 0x69, 0x64,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x42, 0x3a, 0x82, 0xa6, 0x1d, 0x36, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e,
	0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x74, 0x72, 0x61, 0x6e, 0x74, 0x6f, 0x72, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x54,
	0x78, 0x49, 0x44, 0x52, 0x05, 0x74, 0x78, 0x49, 0x64, 0x73, 0x12, 0x24, 0x0a, 0x03, 0x74, 0x78,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x03, 0x74, 0x78, 0x73,
	0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x39, 0x0a, 0x06,
	0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x64, 0x62, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x42, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x52,
	0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0x4e, 0x0a,
	0x0b, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x12, 0x39, 0x0a, 0x06,
	0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x64, 0x62, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x42, 0x04, 0xa0, 0xa6, 0x1d, 0x01, 0x52,
	0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0xd9, 0x01,
	0x0a, 0x11, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x42, 0x61, 0x74, 0x63, 0x68, 0x4f, 0x72, 0x69,
	0x67, 0x69, 0x6e, 0x12, 0x4e, 0x0a, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x36, 0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x06, 0x6d, 0x6f, 0x64,
	0x75, 0x6c, 0x65, 0x12, 0x3d, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x78, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x69, 0x67,
	0x69, 0x6e, 0x48, 0x00, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x12, 0x21, 0x0a, 0x03, 0x64, 0x73, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0d, 0x2e, 0x64, 0x73, 0x6c, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x48, 0x00,
	0x52, 0x03, 0x64, 0x73, 0x6c, 0x3a, 0x04, 0x90, 0xa6, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0xd8, 0x01, 0x0a, 0x10, 0x53, 0x74,
	0x6f, 0x72, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x4e,
	0x0a, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x36,
	0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x3d,
	0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x70, 0x62, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x48, 0x00, 0x52,
	0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x21, 0x0a,
	0x03, 0x64, 0x73, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x64, 0x73, 0x6c,
	0x70, 0x62, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x48, 0x00, 0x52, 0x03, 0x64, 0x73, 0x6c,
	0x3a, 0x04, 0x90, 0xa6, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x04,
	0x80, 0xa6, 0x1d, 0x01, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x61,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x70, 0x62, 0x2f, 0x62, 0x61,
	0x74, 0x63, 0x68, 0x64, 0x62, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescOnce sync.Once
	file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescData = file_availabilitypb_batchdbpb_batchdbpb_proto_rawDesc
)

func file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescGZIP() []byte {
	file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescOnce.Do(func() {
		file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescData)
	})
	return file_availabilitypb_batchdbpb_batchdbpb_proto_rawDescData
}

var file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_availabilitypb_batchdbpb_batchdbpb_proto_goTypes = []interface{}{
	(*Event)(nil),                 // 0: batchdbpb.Event
	(*LookupBatch)(nil),           // 1: batchdbpb.LookupBatch
	(*LookupBatchResponse)(nil),   // 2: batchdbpb.LookupBatchResponse
	(*StoreBatch)(nil),            // 3: batchdbpb.StoreBatch
	(*BatchStored)(nil),           // 4: batchdbpb.BatchStored
	(*LookupBatchOrigin)(nil),     // 5: batchdbpb.LookupBatchOrigin
	(*StoreBatchOrigin)(nil),      // 6: batchdbpb.StoreBatchOrigin
	(*requestpb.Request)(nil),     // 7: requestpb.Request
	(*contextstorepb.Origin)(nil), // 8: contextstorepb.Origin
	(*dslpb.Origin)(nil),          // 9: dslpb.Origin
}
var file_availabilitypb_batchdbpb_batchdbpb_proto_depIdxs = []int32{
	1,  // 0: batchdbpb.Event.lookup:type_name -> batchdbpb.LookupBatch
	2,  // 1: batchdbpb.Event.lookup_response:type_name -> batchdbpb.LookupBatchResponse
	3,  // 2: batchdbpb.Event.store:type_name -> batchdbpb.StoreBatch
	4,  // 3: batchdbpb.Event.stored:type_name -> batchdbpb.BatchStored
	5,  // 4: batchdbpb.LookupBatch.origin:type_name -> batchdbpb.LookupBatchOrigin
	7,  // 5: batchdbpb.LookupBatchResponse.txs:type_name -> requestpb.Request
	5,  // 6: batchdbpb.LookupBatchResponse.origin:type_name -> batchdbpb.LookupBatchOrigin
	7,  // 7: batchdbpb.StoreBatch.txs:type_name -> requestpb.Request
	6,  // 8: batchdbpb.StoreBatch.origin:type_name -> batchdbpb.StoreBatchOrigin
	6,  // 9: batchdbpb.BatchStored.origin:type_name -> batchdbpb.StoreBatchOrigin
	8,  // 10: batchdbpb.LookupBatchOrigin.context_store:type_name -> contextstorepb.Origin
	9,  // 11: batchdbpb.LookupBatchOrigin.dsl:type_name -> dslpb.Origin
	8,  // 12: batchdbpb.StoreBatchOrigin.context_store:type_name -> contextstorepb.Origin
	9,  // 13: batchdbpb.StoreBatchOrigin.dsl:type_name -> dslpb.Origin
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_availabilitypb_batchdbpb_batchdbpb_proto_init() }
func file_availabilitypb_batchdbpb_batchdbpb_proto_init() {
	if File_availabilitypb_batchdbpb_batchdbpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LookupBatch); i {
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
		file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LookupBatchResponse); i {
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
		file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreBatch); i {
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
		file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchStored); i {
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
		file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LookupBatchOrigin); i {
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
		file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreBatchOrigin); i {
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
	file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Event_Lookup)(nil),
		(*Event_LookupResponse)(nil),
		(*Event_Store)(nil),
		(*Event_Stored)(nil),
	}
	file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*LookupBatchOrigin_ContextStore)(nil),
		(*LookupBatchOrigin_Dsl)(nil),
	}
	file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*StoreBatchOrigin_ContextStore)(nil),
		(*StoreBatchOrigin_Dsl)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_availabilitypb_batchdbpb_batchdbpb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_availabilitypb_batchdbpb_batchdbpb_proto_goTypes,
		DependencyIndexes: file_availabilitypb_batchdbpb_batchdbpb_proto_depIdxs,
		MessageInfos:      file_availabilitypb_batchdbpb_batchdbpb_proto_msgTypes,
	}.Build()
	File_availabilitypb_batchdbpb_batchdbpb_proto = out.File
	file_availabilitypb_batchdbpb_batchdbpb_proto_rawDesc = nil
	file_availabilitypb_batchdbpb_batchdbpb_proto_goTypes = nil
	file_availabilitypb_batchdbpb_batchdbpb_proto_depIdxs = nil
}
