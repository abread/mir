// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: availabilitypb/mscpb/mscpb.proto

package mscpb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	commonpb "github.com/filecoin-project/mir/pkg/pb/commonpb"
	_ "github.com/filecoin-project/mir/pkg/pb/mir"
	_ "github.com/filecoin-project/mir/pkg/pb/net"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//	*Message_RequestSig
	//	*Message_Sig
	//	*Message_RequestBatch
	//	*Message_ProvideBatch
	Type isMessage_Type `protobuf_oneof:"type"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[0]
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
	return file_availabilitypb_mscpb_mscpb_proto_rawDescGZIP(), []int{0}
}

func (m *Message) GetType() isMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Message) GetRequestSig() *RequestSigMessage {
	if x, ok := x.GetType().(*Message_RequestSig); ok {
		return x.RequestSig
	}
	return nil
}

func (x *Message) GetSig() *SigMessage {
	if x, ok := x.GetType().(*Message_Sig); ok {
		return x.Sig
	}
	return nil
}

func (x *Message) GetRequestBatch() *RequestBatchMessage {
	if x, ok := x.GetType().(*Message_RequestBatch); ok {
		return x.RequestBatch
	}
	return nil
}

func (x *Message) GetProvideBatch() *ProvideBatchMessage {
	if x, ok := x.GetType().(*Message_ProvideBatch); ok {
		return x.ProvideBatch
	}
	return nil
}

type isMessage_Type interface {
	isMessage_Type()
}

type Message_RequestSig struct {
	RequestSig *RequestSigMessage `protobuf:"bytes,1,opt,name=request_sig,json=requestSig,proto3,oneof"`
}

type Message_Sig struct {
	Sig *SigMessage `protobuf:"bytes,2,opt,name=sig,proto3,oneof"`
}

type Message_RequestBatch struct {
	RequestBatch *RequestBatchMessage `protobuf:"bytes,3,opt,name=request_batch,json=requestBatch,proto3,oneof"`
}

type Message_ProvideBatch struct {
	ProvideBatch *ProvideBatchMessage `protobuf:"bytes,4,opt,name=provide_batch,json=provideBatch,proto3,oneof"`
}

func (*Message_RequestSig) isMessage_Type() {}

func (*Message_Sig) isMessage_Type() {}

func (*Message_RequestBatch) isMessage_Type() {}

func (*Message_ProvideBatch) isMessage_Type() {}

type RequestSigMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txs   []*requestpb.Request `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
	ReqId uint64               `protobuf:"varint,2,opt,name=req_id,json=reqId,proto3" json:"req_id,omitempty"`
}

func (x *RequestSigMessage) Reset() {
	*x = RequestSigMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestSigMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestSigMessage) ProtoMessage() {}

func (x *RequestSigMessage) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestSigMessage.ProtoReflect.Descriptor instead.
func (*RequestSigMessage) Descriptor() ([]byte, []int) {
	return file_availabilitypb_mscpb_mscpb_proto_rawDescGZIP(), []int{1}
}

func (x *RequestSigMessage) GetTxs() []*requestpb.Request {
	if x != nil {
		return x.Txs
	}
	return nil
}

func (x *RequestSigMessage) GetReqId() uint64 {
	if x != nil {
		return x.ReqId
	}
	return 0
}

type SigMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signature []byte `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	ReqId     uint64 `protobuf:"varint,2,opt,name=req_id,json=reqId,proto3" json:"req_id,omitempty"`
}

func (x *SigMessage) Reset() {
	*x = SigMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SigMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SigMessage) ProtoMessage() {}

func (x *SigMessage) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SigMessage.ProtoReflect.Descriptor instead.
func (*SigMessage) Descriptor() ([]byte, []int) {
	return file_availabilitypb_mscpb_mscpb_proto_rawDescGZIP(), []int{2}
}

func (x *SigMessage) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *SigMessage) GetReqId() uint64 {
	if x != nil {
		return x.ReqId
	}
	return 0
}

type RequestBatchMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BatchId []byte `protobuf:"bytes,1,opt,name=batch_id,json=batchId,proto3" json:"batch_id,omitempty"`
	ReqId   uint64 `protobuf:"varint,2,opt,name=req_id,json=reqId,proto3" json:"req_id,omitempty"`
}

func (x *RequestBatchMessage) Reset() {
	*x = RequestBatchMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestBatchMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestBatchMessage) ProtoMessage() {}

func (x *RequestBatchMessage) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestBatchMessage.ProtoReflect.Descriptor instead.
func (*RequestBatchMessage) Descriptor() ([]byte, []int) {
	return file_availabilitypb_mscpb_mscpb_proto_rawDescGZIP(), []int{3}
}

func (x *RequestBatchMessage) GetBatchId() []byte {
	if x != nil {
		return x.BatchId
	}
	return nil
}

func (x *RequestBatchMessage) GetReqId() uint64 {
	if x != nil {
		return x.ReqId
	}
	return 0
}

type ProvideBatchMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txs     []*requestpb.Request `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
	ReqId   uint64               `protobuf:"varint,2,opt,name=req_id,json=reqId,proto3" json:"req_id,omitempty"`
	BatchId []byte               `protobuf:"bytes,3,opt,name=batch_id,json=batchId,proto3" json:"batch_id,omitempty"`
}

func (x *ProvideBatchMessage) Reset() {
	*x = ProvideBatchMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProvideBatchMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProvideBatchMessage) ProtoMessage() {}

func (x *ProvideBatchMessage) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProvideBatchMessage.ProtoReflect.Descriptor instead.
func (*ProvideBatchMessage) Descriptor() ([]byte, []int) {
	return file_availabilitypb_mscpb_mscpb_proto_rawDescGZIP(), []int{4}
}

func (x *ProvideBatchMessage) GetTxs() []*requestpb.Request {
	if x != nil {
		return x.Txs
	}
	return nil
}

func (x *ProvideBatchMessage) GetReqId() uint64 {
	if x != nil {
		return x.ReqId
	}
	return 0
}

func (x *ProvideBatchMessage) GetBatchId() []byte {
	if x != nil {
		return x.BatchId
	}
	return nil
}

// ============================================================
// Data structures
// ============================================================
type Cert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BatchId []byte `protobuf:"bytes,1,opt,name=batch_id,json=batchId,proto3" json:"batch_id,omitempty"`
	// TODO: can be encoded as n bits
	Signers    []string `protobuf:"bytes,2,rep,name=signers,proto3" json:"signers,omitempty"`
	Signatures [][]byte `protobuf:"bytes,3,rep,name=signatures,proto3" json:"signatures,omitempty"`
}

func (x *Cert) Reset() {
	*x = Cert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cert) ProtoMessage() {}

func (x *Cert) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cert.ProtoReflect.Descriptor instead.
func (*Cert) Descriptor() ([]byte, []int) {
	return file_availabilitypb_mscpb_mscpb_proto_rawDescGZIP(), []int{5}
}

func (x *Cert) GetBatchId() []byte {
	if x != nil {
		return x.BatchId
	}
	return nil
}

func (x *Cert) GetSigners() []string {
	if x != nil {
		return x.Signers
	}
	return nil
}

func (x *Cert) GetSignatures() [][]byte {
	if x != nil {
		return x.Signatures
	}
	return nil
}

type Certs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Certs []*Cert `protobuf:"bytes,1,rep,name=certs,proto3" json:"certs,omitempty"`
}

func (x *Certs) Reset() {
	*x = Certs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Certs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Certs) ProtoMessage() {}

func (x *Certs) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Certs.ProtoReflect.Descriptor instead.
func (*Certs) Descriptor() ([]byte, []int) {
	return file_availabilitypb_mscpb_mscpb_proto_rawDescGZIP(), []int{6}
}

func (x *Certs) GetCerts() []*Cert {
	if x != nil {
		return x.Certs
	}
	return nil
}

type InstanceParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Membership  *commonpb.Membership `protobuf:"bytes,1,opt,name=membership,proto3" json:"membership,omitempty"`
	Limit       uint64               `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	MaxRequests uint64               `protobuf:"varint,3,opt,name=max_requests,json=maxRequests,proto3" json:"max_requests,omitempty"`
}

func (x *InstanceParams) Reset() {
	*x = InstanceParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceParams) ProtoMessage() {}

func (x *InstanceParams) ProtoReflect() protoreflect.Message {
	mi := &file_availabilitypb_mscpb_mscpb_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceParams.ProtoReflect.Descriptor instead.
func (*InstanceParams) Descriptor() ([]byte, []int) {
	return file_availabilitypb_mscpb_mscpb_proto_rawDescGZIP(), []int{7}
}

func (x *InstanceParams) GetMembership() *commonpb.Membership {
	if x != nil {
		return x.Membership
	}
	return nil
}

func (x *InstanceParams) GetLimit() uint64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *InstanceParams) GetMaxRequests() uint64 {
	if x != nil {
		return x.MaxRequests
	}
	return 0
}

var File_availabilitypb_mscpb_mscpb_proto protoreflect.FileDescriptor

var file_availabilitypb_mscpb_mscpb_proto_rawDesc = []byte{
	0x0a, 0x20, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x70, 0x62,
	0x2f, 0x6d, 0x73, 0x63, 0x70, 0x62, 0x2f, 0x6d, 0x73, 0x63, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x14, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x70, 0x62, 0x2e, 0x6d, 0x73, 0x63, 0x70, 0x62, 0x1a, 0x19, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x70, 0x62, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x62, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6e, 0x65,
	0x74, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72, 0x2f,
	0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc3, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x4a, 0x0a, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f,
	0x73, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x61, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x70, 0x62, 0x2e, 0x6d, 0x73, 0x63, 0x70, 0x62,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x69, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x69, 0x67,
	0x12, 0x34, 0x0a, 0x03, 0x73, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e,
	0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x70, 0x62, 0x2e, 0x6d,
	0x73, 0x63, 0x70, 0x62, 0x2e, 0x53, 0x69, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48,
	0x00, 0x52, 0x03, 0x73, 0x69, 0x67, 0x12, 0x50, 0x0a, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x5f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e,
	0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x70, 0x62, 0x2e, 0x6d,
	0x73, 0x63, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0c, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x50, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x76,
	0x69, 0x64, 0x65, 0x5f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x29, 0x2e, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x70, 0x62,
	0x2e, 0x6d, 0x73, 0x63, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0c, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x3a, 0x04, 0xc8, 0xe4, 0x1d, 0x01,
	0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0xc8, 0xe4, 0x1d, 0x01, 0x22, 0x56,
	0x0a, 0x11, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x69, 0x67, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x52, 0x03, 0x74, 0x78, 0x73, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x65, 0x71,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x65, 0x71, 0x49, 0x64,
	0x3a, 0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x22, 0x47, 0x0a, 0x0a, 0x53, 0x69, 0x67, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x65, 0x71, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x72, 0x65, 0x71, 0x49, 0x64, 0x3a, 0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x22,
	0x4d, 0x0a, 0x13, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x62, 0x61, 0x74, 0x63, 0x68, 0x49,
	0x64, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x65, 0x71, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x72, 0x65, 0x71, 0x49, 0x64, 0x3a, 0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x22, 0x73,
	0x0a, 0x13, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x03, 0x74, 0x78, 0x73, 0x12, 0x15, 0x0a, 0x06, 0x72,
	0x65, 0x71, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x65, 0x71,
	0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x62, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x3a, 0x04, 0xd0,
	0xe4, 0x1d, 0x01, 0x22, 0x97, 0x01, 0x0a, 0x04, 0x43, 0x65, 0x72, 0x74, 0x12, 0x19, 0x0a, 0x08,
	0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07,
	0x62, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x12, 0x4e, 0x0a, 0x07, 0x73, 0x69, 0x67, 0x6e, 0x65,
	0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x42, 0x34, 0x82, 0xa6, 0x1d, 0x30, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69,
	0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x52, 0x07,
	0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0a, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x3a, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0x3f, 0x0a,
	0x05, 0x43, 0x65, 0x72, 0x74, 0x73, 0x12, 0x30, 0x0a, 0x05, 0x63, 0x65, 0x72, 0x74, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x79, 0x70, 0x62, 0x2e, 0x6d, 0x73, 0x63, 0x70, 0x62, 0x2e, 0x43, 0x65, 0x72,
	0x74, 0x52, 0x05, 0x63, 0x65, 0x72, 0x74, 0x73, 0x3a, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0x85,
	0x01, 0x0a, 0x0e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x12, 0x34, 0x0a, 0x0a, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x70, 0x62,
	0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x52, 0x0a, 0x6d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x21, 0x0a,
	0x0c, 0x6d, 0x61, 0x78, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x6d, 0x61, 0x78, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73,
	0x3a, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x42, 0x3d, 0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62,
	0x2f, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x70, 0x62, 0x2f,
	0x6d, 0x73, 0x63, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_availabilitypb_mscpb_mscpb_proto_rawDescOnce sync.Once
	file_availabilitypb_mscpb_mscpb_proto_rawDescData = file_availabilitypb_mscpb_mscpb_proto_rawDesc
)

func file_availabilitypb_mscpb_mscpb_proto_rawDescGZIP() []byte {
	file_availabilitypb_mscpb_mscpb_proto_rawDescOnce.Do(func() {
		file_availabilitypb_mscpb_mscpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_availabilitypb_mscpb_mscpb_proto_rawDescData)
	})
	return file_availabilitypb_mscpb_mscpb_proto_rawDescData
}

var file_availabilitypb_mscpb_mscpb_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_availabilitypb_mscpb_mscpb_proto_goTypes = []interface{}{
	(*Message)(nil),             // 0: availabilitypb.mscpb.Message
	(*RequestSigMessage)(nil),   // 1: availabilitypb.mscpb.RequestSigMessage
	(*SigMessage)(nil),          // 2: availabilitypb.mscpb.SigMessage
	(*RequestBatchMessage)(nil), // 3: availabilitypb.mscpb.RequestBatchMessage
	(*ProvideBatchMessage)(nil), // 4: availabilitypb.mscpb.ProvideBatchMessage
	(*Cert)(nil),                // 5: availabilitypb.mscpb.Cert
	(*Certs)(nil),               // 6: availabilitypb.mscpb.Certs
	(*InstanceParams)(nil),      // 7: availabilitypb.mscpb.InstanceParams
	(*requestpb.Request)(nil),   // 8: requestpb.Request
	(*commonpb.Membership)(nil), // 9: commonpb.Membership
}
var file_availabilitypb_mscpb_mscpb_proto_depIdxs = []int32{
	1, // 0: availabilitypb.mscpb.Message.request_sig:type_name -> availabilitypb.mscpb.RequestSigMessage
	2, // 1: availabilitypb.mscpb.Message.sig:type_name -> availabilitypb.mscpb.SigMessage
	3, // 2: availabilitypb.mscpb.Message.request_batch:type_name -> availabilitypb.mscpb.RequestBatchMessage
	4, // 3: availabilitypb.mscpb.Message.provide_batch:type_name -> availabilitypb.mscpb.ProvideBatchMessage
	8, // 4: availabilitypb.mscpb.RequestSigMessage.txs:type_name -> requestpb.Request
	8, // 5: availabilitypb.mscpb.ProvideBatchMessage.txs:type_name -> requestpb.Request
	5, // 6: availabilitypb.mscpb.Certs.certs:type_name -> availabilitypb.mscpb.Cert
	9, // 7: availabilitypb.mscpb.InstanceParams.membership:type_name -> commonpb.Membership
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_availabilitypb_mscpb_mscpb_proto_init() }
func file_availabilitypb_mscpb_mscpb_proto_init() {
	if File_availabilitypb_mscpb_mscpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_availabilitypb_mscpb_mscpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_availabilitypb_mscpb_mscpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestSigMessage); i {
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
		file_availabilitypb_mscpb_mscpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SigMessage); i {
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
		file_availabilitypb_mscpb_mscpb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestBatchMessage); i {
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
		file_availabilitypb_mscpb_mscpb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProvideBatchMessage); i {
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
		file_availabilitypb_mscpb_mscpb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cert); i {
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
		file_availabilitypb_mscpb_mscpb_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Certs); i {
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
		file_availabilitypb_mscpb_mscpb_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceParams); i {
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
	file_availabilitypb_mscpb_mscpb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Message_RequestSig)(nil),
		(*Message_Sig)(nil),
		(*Message_RequestBatch)(nil),
		(*Message_ProvideBatch)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_availabilitypb_mscpb_mscpb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_availabilitypb_mscpb_mscpb_proto_goTypes,
		DependencyIndexes: file_availabilitypb_mscpb_mscpb_proto_depIdxs,
		MessageInfos:      file_availabilitypb_mscpb_mscpb_proto_msgTypes,
	}.Build()
	File_availabilitypb_mscpb_mscpb_proto = out.File
	file_availabilitypb_mscpb_mscpb_proto_rawDesc = nil
	file_availabilitypb_mscpb_mscpb_proto_goTypes = nil
	file_availabilitypb_mscpb_mscpb_proto_depIdxs = nil
}
