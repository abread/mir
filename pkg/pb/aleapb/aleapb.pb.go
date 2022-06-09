// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: aleapb/aleapb.proto

package aleapb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AleaMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//	*AleaMessage_Broadcast
	//	*AleaMessage_Agreement
	Type isAleaMessage_Type `protobuf_oneof:"type"`
}

func (x *AleaMessage) Reset() {
	*x = AleaMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AleaMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AleaMessage) ProtoMessage() {}

func (x *AleaMessage) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AleaMessage.ProtoReflect.Descriptor instead.
func (*AleaMessage) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{0}
}

func (m *AleaMessage) GetType() isAleaMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *AleaMessage) GetBroadcast() *VCBC {
	if x, ok := x.GetType().(*AleaMessage_Broadcast); ok {
		return x.Broadcast
	}
	return nil
}

func (x *AleaMessage) GetAgreement() *Agreement {
	if x, ok := x.GetType().(*AleaMessage_Agreement); ok {
		return x.Agreement
	}
	return nil
}

type isAleaMessage_Type interface {
	isAleaMessage_Type()
}

type AleaMessage_Broadcast struct {
	Broadcast *VCBC `protobuf:"bytes,1,opt,name=broadcast,proto3,oneof"`
}

type AleaMessage_Agreement struct {
	Agreement *Agreement `protobuf:"bytes,2,opt,name=agreement,proto3,oneof"`
}

func (*AleaMessage_Broadcast) isAleaMessage_Type() {}

func (*AleaMessage_Agreement) isAleaMessage_Type() {}

type Agreement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//	*Agreement_Aba
	//	*Agreement_FillGap
	//	*Agreement_Filler
	Type isAgreement_Type `protobuf_oneof:"type"`
}

func (x *Agreement) Reset() {
	*x = Agreement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Agreement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Agreement) ProtoMessage() {}

func (x *Agreement) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Agreement.ProtoReflect.Descriptor instead.
func (*Agreement) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{1}
}

func (m *Agreement) GetType() isAgreement_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Agreement) GetAba() *CobaltABBA {
	if x, ok := x.GetType().(*Agreement_Aba); ok {
		return x.Aba
	}
	return nil
}

func (x *Agreement) GetFillGap() *AgreementFillGap {
	if x, ok := x.GetType().(*Agreement_FillGap); ok {
		return x.FillGap
	}
	return nil
}

func (x *Agreement) GetFiller() *AgreementFiller {
	if x, ok := x.GetType().(*Agreement_Filler); ok {
		return x.Filler
	}
	return nil
}

type isAgreement_Type interface {
	isAgreement_Type()
}

type Agreement_Aba struct {
	Aba *CobaltABBA `protobuf:"bytes,1,opt,name=aba,proto3,oneof"`
}

type Agreement_FillGap struct {
	FillGap *AgreementFillGap `protobuf:"bytes,2,opt,name=fillGap,proto3,oneof"`
}

type Agreement_Filler struct {
	Filler *AgreementFiller `protobuf:"bytes,3,opt,name=filler,proto3,oneof"`
}

func (*Agreement_Aba) isAgreement_Type() {}

func (*Agreement_FillGap) isAgreement_Type() {}

func (*Agreement_Filler) isAgreement_Type() {}

type AgreementFillGap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueueIdx uint32 `protobuf:"varint,1,opt,name=queue_idx,json=queueIdx,proto3" json:"queue_idx,omitempty"`
	Slot     uint64 `protobuf:"varint,2,opt,name=slot,proto3" json:"slot,omitempty"` // TODO: confirm type
}

func (x *AgreementFillGap) Reset() {
	*x = AgreementFillGap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgreementFillGap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgreementFillGap) ProtoMessage() {}

func (x *AgreementFillGap) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgreementFillGap.ProtoReflect.Descriptor instead.
func (*AgreementFillGap) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{2}
}

func (x *AgreementFillGap) GetQueueIdx() uint32 {
	if x != nil {
		return x.QueueIdx
	}
	return 0
}

func (x *AgreementFillGap) GetSlot() uint64 {
	if x != nil {
		return x.Slot
	}
	return 0
}

type AgreementFiller struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries []*VCBC `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *AgreementFiller) Reset() {
	*x = AgreementFiller{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgreementFiller) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgreementFiller) ProtoMessage() {}

func (x *AgreementFiller) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgreementFiller.ProtoReflect.Descriptor instead.
func (*AgreementFiller) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{3}
}

func (x *AgreementFiller) GetEntries() []*VCBC {
	if x != nil {
		return x.Entries
	}
	return nil
}

type VCBC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instance *VCBC_VCBCId `protobuf:"bytes,1,opt,name=instance,proto3" json:"instance,omitempty"`
	// Types that are assignable to Type:
	//	*VCBC_Send
	//	*VCBC_Echo
	//	*VCBC_Final
	Type isVCBC_Type `protobuf_oneof:"type"`
}

func (x *VCBC) Reset() {
	*x = VCBC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VCBC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VCBC) ProtoMessage() {}

func (x *VCBC) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VCBC.ProtoReflect.Descriptor instead.
func (*VCBC) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{4}
}

func (x *VCBC) GetInstance() *VCBC_VCBCId {
	if x != nil {
		return x.Instance
	}
	return nil
}

func (m *VCBC) GetType() isVCBC_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *VCBC) GetSend() *VCBCSend {
	if x, ok := x.GetType().(*VCBC_Send); ok {
		return x.Send
	}
	return nil
}

func (x *VCBC) GetEcho() *VCBCEcho {
	if x, ok := x.GetType().(*VCBC_Echo); ok {
		return x.Echo
	}
	return nil
}

func (x *VCBC) GetFinal() *VCBCFinal {
	if x, ok := x.GetType().(*VCBC_Final); ok {
		return x.Final
	}
	return nil
}

type isVCBC_Type interface {
	isVCBC_Type()
}

type VCBC_Send struct {
	Send *VCBCSend `protobuf:"bytes,2,opt,name=send,proto3,oneof"`
}

type VCBC_Echo struct {
	Echo *VCBCEcho `protobuf:"bytes,3,opt,name=echo,proto3,oneof"`
}

type VCBC_Final struct {
	Final *VCBCFinal `protobuf:"bytes,4,opt,name=final,proto3,oneof"`
}

func (*VCBC_Send) isVCBC_Type() {}

func (*VCBC_Echo) isVCBC_Type() {}

func (*VCBC_Final) isVCBC_Type() {}

type VCBCSend struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload *requestpb.Batch `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *VCBCSend) Reset() {
	*x = VCBCSend{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VCBCSend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VCBCSend) ProtoMessage() {}

func (x *VCBCSend) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VCBCSend.ProtoReflect.Descriptor instead.
func (*VCBCSend) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{5}
}

func (x *VCBCSend) GetPayload() *requestpb.Batch {
	if x != nil {
		return x.Payload
	}
	return nil
}

type VCBCEcho struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload        *requestpb.Batch `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	SignatureShare []byte           `protobuf:"bytes,2,opt,name=signature_share,json=signatureShare,proto3" json:"signature_share,omitempty"`
}

func (x *VCBCEcho) Reset() {
	*x = VCBCEcho{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VCBCEcho) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VCBCEcho) ProtoMessage() {}

func (x *VCBCEcho) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VCBCEcho.ProtoReflect.Descriptor instead.
func (*VCBCEcho) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{6}
}

func (x *VCBCEcho) GetPayload() *requestpb.Batch {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *VCBCEcho) GetSignatureShare() []byte {
	if x != nil {
		return x.SignatureShare
	}
	return nil
}

type VCBCFinal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload   *requestpb.Batch `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Signature []byte           `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *VCBCFinal) Reset() {
	*x = VCBCFinal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VCBCFinal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VCBCFinal) ProtoMessage() {}

func (x *VCBCFinal) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VCBCFinal.ProtoReflect.Descriptor instead.
func (*VCBCFinal) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{7}
}

func (x *VCBCFinal) GetPayload() *requestpb.Batch {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *VCBCFinal) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type CobaltABBA struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instance uint64 `protobuf:"varint,1,opt,name=instance,proto3" json:"instance,omitempty"` // TODO
}

func (x *CobaltABBA) Reset() {
	*x = CobaltABBA{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CobaltABBA) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CobaltABBA) ProtoMessage() {}

func (x *CobaltABBA) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CobaltABBA.ProtoReflect.Descriptor instead.
func (*CobaltABBA) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{8}
}

func (x *CobaltABBA) GetInstance() uint64 {
	if x != nil {
		return x.Instance
	}
	return 0
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{9}
}

type VCBC_VCBCId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InitiatorId uint64 `protobuf:"varint,1,opt,name=initiator_id,json=initiatorId,proto3" json:"initiator_id,omitempty"`
	Priority    uint64 `protobuf:"varint,2,opt,name=priority,proto3" json:"priority,omitempty"`
}

func (x *VCBC_VCBCId) Reset() {
	*x = VCBC_VCBCId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_aleapb_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VCBC_VCBCId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VCBC_VCBCId) ProtoMessage() {}

func (x *VCBC_VCBCId) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_aleapb_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VCBC_VCBCId.ProtoReflect.Descriptor instead.
func (*VCBC_VCBCId) Descriptor() ([]byte, []int) {
	return file_aleapb_aleapb_proto_rawDescGZIP(), []int{4, 0}
}

func (x *VCBC_VCBCId) GetInitiatorId() uint64 {
	if x != nil {
		return x.InitiatorId
	}
	return 0
}

func (x *VCBC_VCBCId) GetPriority() uint64 {
	if x != nil {
		return x.Priority
	}
	return 0
}

var File_aleapb_aleapb_proto protoreflect.FileDescriptor

var file_aleapb_aleapb_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x1a, 0x19, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x76, 0x0a, 0x0b, 0x41, 0x6c, 0x65, 0x61,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2c, 0x0a, 0x09, 0x62, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x61, 0x6c, 0x65,
	0x61, 0x70, 0x62, 0x2e, 0x56, 0x43, 0x42, 0x43, 0x48, 0x00, 0x52, 0x09, 0x62, 0x72, 0x6f, 0x61,
	0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x31, 0x0a, 0x09, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70,
	0x62, 0x2e, 0x41, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x09, 0x61,
	0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x22, 0xa4, 0x01, 0x0a, 0x09, 0x41, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x26,
	0x0a, 0x03, 0x61, 0x62, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x6c,
	0x65, 0x61, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x62, 0x61, 0x6c, 0x74, 0x41, 0x42, 0x42, 0x41, 0x48,
	0x00, 0x52, 0x03, 0x61, 0x62, 0x61, 0x12, 0x34, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x6c, 0x47, 0x61,
	0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62,
	0x2e, 0x41, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x6c, 0x47, 0x61,
	0x70, 0x48, 0x00, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x6c, 0x47, 0x61, 0x70, 0x12, 0x31, 0x0a, 0x06,
	0x66, 0x69, 0x6c, 0x6c, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61,
	0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x41, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x46,
	0x69, 0x6c, 0x6c, 0x65, 0x72, 0x48, 0x00, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x6c, 0x65, 0x72, 0x42,
	0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x43, 0x0a, 0x10, 0x41, 0x67, 0x72, 0x65, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x6c, 0x47, 0x61, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x5f, 0x69, 0x64, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x49, 0x64, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6c, 0x6f, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x6c, 0x6f, 0x74, 0x22, 0x39, 0x0a, 0x0f,
	0x41, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x6c, 0x65, 0x72, 0x12,
	0x26, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x56, 0x43, 0x42, 0x43, 0x52, 0x07,
	0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x83, 0x02, 0x0a, 0x04, 0x56, 0x43, 0x42, 0x43,
	0x12, 0x2f, 0x0a, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x56, 0x43, 0x42, 0x43,
	0x2e, 0x56, 0x43, 0x42, 0x43, 0x49, 0x64, 0x52, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x12, 0x26, 0x0a, 0x04, 0x73, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x56, 0x43, 0x42, 0x43, 0x53, 0x65, 0x6e,
	0x64, 0x48, 0x00, 0x52, 0x04, 0x73, 0x65, 0x6e, 0x64, 0x12, 0x26, 0x0a, 0x04, 0x65, 0x63, 0x68,
	0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62,
	0x2e, 0x56, 0x43, 0x42, 0x43, 0x45, 0x63, 0x68, 0x6f, 0x48, 0x00, 0x52, 0x04, 0x65, 0x63, 0x68,
	0x6f, 0x12, 0x29, 0x0a, 0x05, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x56, 0x43, 0x42, 0x43, 0x46, 0x69,
	0x6e, 0x61, 0x6c, 0x48, 0x00, 0x52, 0x05, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x1a, 0x47, 0x0a, 0x06,
	0x56, 0x43, 0x42, 0x43, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61,
	0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x69, 0x6e,
	0x69, 0x74, 0x69, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x79, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x36, 0x0a,
	0x08, 0x56, 0x43, 0x42, 0x43, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x2a, 0x0a, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x5f, 0x0a, 0x08, 0x56, 0x43, 0x42, 0x43, 0x45, 0x63, 0x68,
	0x6f, 0x12, 0x2a, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x27, 0x0a,
	0x0f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0e, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x53, 0x68, 0x61, 0x72, 0x65, 0x22, 0x55, 0x0a, 0x09, 0x56, 0x43, 0x42, 0x43, 0x46, 0x69,
	0x6e, 0x61, 0x6c, 0x12, 0x2a, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62,
	0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x28, 0x0a,
	0x0a, 0x43, 0x6f, 0x62, 0x61, 0x6c, 0x74, 0x41, 0x42, 0x42, 0x41, 0x12, 0x1a, 0x0a, 0x08, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x08, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x68, 0x79, 0x70, 0x65, 0x72, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2d, 0x6c, 0x61, 0x62, 0x73,
	0x2f, 0x6d, 0x69, 0x72, 0x62, 0x66, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x61,
	0x6c, 0x65, 0x61, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aleapb_aleapb_proto_rawDescOnce sync.Once
	file_aleapb_aleapb_proto_rawDescData = file_aleapb_aleapb_proto_rawDesc
)

func file_aleapb_aleapb_proto_rawDescGZIP() []byte {
	file_aleapb_aleapb_proto_rawDescOnce.Do(func() {
		file_aleapb_aleapb_proto_rawDescData = protoimpl.X.CompressGZIP(file_aleapb_aleapb_proto_rawDescData)
	})
	return file_aleapb_aleapb_proto_rawDescData
}

var file_aleapb_aleapb_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_aleapb_aleapb_proto_goTypes = []interface{}{
	(*AleaMessage)(nil),      // 0: aleapb.AleaMessage
	(*Agreement)(nil),        // 1: aleapb.Agreement
	(*AgreementFillGap)(nil), // 2: aleapb.AgreementFillGap
	(*AgreementFiller)(nil),  // 3: aleapb.AgreementFiller
	(*VCBC)(nil),             // 4: aleapb.VCBC
	(*VCBCSend)(nil),         // 5: aleapb.VCBCSend
	(*VCBCEcho)(nil),         // 6: aleapb.VCBCEcho
	(*VCBCFinal)(nil),        // 7: aleapb.VCBCFinal
	(*CobaltABBA)(nil),       // 8: aleapb.CobaltABBA
	(*Status)(nil),           // 9: aleapb.Status
	(*VCBC_VCBCId)(nil),      // 10: aleapb.VCBC.VCBCId
	(*requestpb.Batch)(nil),  // 11: requestpb.Batch
}
var file_aleapb_aleapb_proto_depIdxs = []int32{
	4,  // 0: aleapb.AleaMessage.broadcast:type_name -> aleapb.VCBC
	1,  // 1: aleapb.AleaMessage.agreement:type_name -> aleapb.Agreement
	8,  // 2: aleapb.Agreement.aba:type_name -> aleapb.CobaltABBA
	2,  // 3: aleapb.Agreement.fillGap:type_name -> aleapb.AgreementFillGap
	3,  // 4: aleapb.Agreement.filler:type_name -> aleapb.AgreementFiller
	4,  // 5: aleapb.AgreementFiller.entries:type_name -> aleapb.VCBC
	10, // 6: aleapb.VCBC.instance:type_name -> aleapb.VCBC.VCBCId
	5,  // 7: aleapb.VCBC.send:type_name -> aleapb.VCBCSend
	6,  // 8: aleapb.VCBC.echo:type_name -> aleapb.VCBCEcho
	7,  // 9: aleapb.VCBC.final:type_name -> aleapb.VCBCFinal
	11, // 10: aleapb.VCBCSend.payload:type_name -> requestpb.Batch
	11, // 11: aleapb.VCBCEcho.payload:type_name -> requestpb.Batch
	11, // 12: aleapb.VCBCFinal.payload:type_name -> requestpb.Batch
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_aleapb_aleapb_proto_init() }
func file_aleapb_aleapb_proto_init() {
	if File_aleapb_aleapb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aleapb_aleapb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AleaMessage); i {
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
		file_aleapb_aleapb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Agreement); i {
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
		file_aleapb_aleapb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgreementFillGap); i {
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
		file_aleapb_aleapb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgreementFiller); i {
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
		file_aleapb_aleapb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VCBC); i {
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
		file_aleapb_aleapb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VCBCSend); i {
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
		file_aleapb_aleapb_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VCBCEcho); i {
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
		file_aleapb_aleapb_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VCBCFinal); i {
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
		file_aleapb_aleapb_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CobaltABBA); i {
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
		file_aleapb_aleapb_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
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
		file_aleapb_aleapb_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VCBC_VCBCId); i {
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
	file_aleapb_aleapb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*AleaMessage_Broadcast)(nil),
		(*AleaMessage_Agreement)(nil),
	}
	file_aleapb_aleapb_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Agreement_Aba)(nil),
		(*Agreement_FillGap)(nil),
		(*Agreement_Filler)(nil),
	}
	file_aleapb_aleapb_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*VCBC_Send)(nil),
		(*VCBC_Echo)(nil),
		(*VCBC_Final)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aleapb_aleapb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aleapb_aleapb_proto_goTypes,
		DependencyIndexes: file_aleapb_aleapb_proto_depIdxs,
		MessageInfos:      file_aleapb_aleapb_proto_msgTypes,
	}.Build()
	File_aleapb_aleapb_proto = out.File
	file_aleapb_aleapb_proto_rawDesc = nil
	file_aleapb_aleapb_proto_goTypes = nil
	file_aleapb_aleapb_proto_depIdxs = nil
}
