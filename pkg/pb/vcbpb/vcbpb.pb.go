// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1
// 	protoc        v3.21.12
// source: vcbpb/vcbpb.proto

package vcbpb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	_ "github.com/filecoin-project/mir/pkg/pb/mir"
	_ "github.com/filecoin-project/mir/pkg/pb/net"
	trantorpb "github.com/filecoin-project/mir/pkg/pb/trantorpb"
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
	//	*Event_QuorumDone
	//	*Event_AllDone
	Type isEvent_Type `protobuf_oneof:"type"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vcbpb_vcbpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_vcbpb_vcbpb_proto_msgTypes[0]
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
	return file_vcbpb_vcbpb_proto_rawDescGZIP(), []int{0}
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

func (x *Event) GetQuorumDone() *QuorumDone {
	if x, ok := x.GetType().(*Event_QuorumDone); ok {
		return x.QuorumDone
	}
	return nil
}

func (x *Event) GetAllDone() *AllDone {
	if x, ok := x.GetType().(*Event_AllDone); ok {
		return x.AllDone
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

type Event_QuorumDone struct {
	QuorumDone *QuorumDone `protobuf:"bytes,3,opt,name=quorum_done,json=quorumDone,proto3,oneof"`
}

type Event_AllDone struct {
	AllDone *AllDone `protobuf:"bytes,4,opt,name=all_done,json=allDone,proto3,oneof"`
}

func (*Event_InputValue) isEvent_Type() {}

func (*Event_Deliver) isEvent_Type() {}

func (*Event_QuorumDone) isEvent_Type() {}

func (*Event_AllDone) isEvent_Type() {}

type InputValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// empty txs represents a request to listen for VCB (and deliver what was broadcast to us)
	Txs []*trantorpb.Transaction `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
}

func (x *InputValue) Reset() {
	*x = InputValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vcbpb_vcbpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputValue) ProtoMessage() {}

func (x *InputValue) ProtoReflect() protoreflect.Message {
	mi := &file_vcbpb_vcbpb_proto_msgTypes[1]
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
	return file_vcbpb_vcbpb_proto_rawDescGZIP(), []int{1}
}

func (x *InputValue) GetTxs() []*trantorpb.Transaction {
	if x != nil {
		return x.Txs
	}
	return nil
}

type Deliver struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txs       []*trantorpb.Transaction `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
	TxIds     [][]byte                 `protobuf:"bytes,2,rep,name=tx_ids,json=txIds,proto3" json:"tx_ids,omitempty"`
	Signature []byte                   `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	SrcModule string                   `protobuf:"bytes,4,opt,name=src_module,json=srcModule,proto3" json:"src_module,omitempty"`
}

func (x *Deliver) Reset() {
	*x = Deliver{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vcbpb_vcbpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deliver) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deliver) ProtoMessage() {}

func (x *Deliver) ProtoReflect() protoreflect.Message {
	mi := &file_vcbpb_vcbpb_proto_msgTypes[2]
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
	return file_vcbpb_vcbpb_proto_rawDescGZIP(), []int{2}
}

func (x *Deliver) GetTxs() []*trantorpb.Transaction {
	if x != nil {
		return x.Txs
	}
	return nil
}

func (x *Deliver) GetTxIds() [][]byte {
	if x != nil {
		return x.TxIds
	}
	return nil
}

func (x *Deliver) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Deliver) GetSrcModule() string {
	if x != nil {
		return x.SrcModule
	}
	return ""
}

type QuorumDone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcModule string `protobuf:"bytes,1,opt,name=src_module,json=srcModule,proto3" json:"src_module,omitempty"`
}

func (x *QuorumDone) Reset() {
	*x = QuorumDone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vcbpb_vcbpb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuorumDone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuorumDone) ProtoMessage() {}

func (x *QuorumDone) ProtoReflect() protoreflect.Message {
	mi := &file_vcbpb_vcbpb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuorumDone.ProtoReflect.Descriptor instead.
func (*QuorumDone) Descriptor() ([]byte, []int) {
	return file_vcbpb_vcbpb_proto_rawDescGZIP(), []int{3}
}

func (x *QuorumDone) GetSrcModule() string {
	if x != nil {
		return x.SrcModule
	}
	return ""
}

type AllDone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcModule string `protobuf:"bytes,1,opt,name=src_module,json=srcModule,proto3" json:"src_module,omitempty"`
}

func (x *AllDone) Reset() {
	*x = AllDone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vcbpb_vcbpb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllDone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllDone) ProtoMessage() {}

func (x *AllDone) ProtoReflect() protoreflect.Message {
	mi := &file_vcbpb_vcbpb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllDone.ProtoReflect.Descriptor instead.
func (*AllDone) Descriptor() ([]byte, []int) {
	return file_vcbpb_vcbpb_proto_rawDescGZIP(), []int{4}
}

func (x *AllDone) GetSrcModule() string {
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
	//	*Message_SendMessage
	//	*Message_EchoMessage
	//	*Message_FinalMessage
	//	*Message_DoneMessage
	Type isMessage_Type `protobuf_oneof:"type"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vcbpb_vcbpb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_vcbpb_vcbpb_proto_msgTypes[5]
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
	return file_vcbpb_vcbpb_proto_rawDescGZIP(), []int{5}
}

func (m *Message) GetType() isMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Message) GetSendMessage() *SendMessage {
	if x, ok := x.GetType().(*Message_SendMessage); ok {
		return x.SendMessage
	}
	return nil
}

func (x *Message) GetEchoMessage() *EchoMessage {
	if x, ok := x.GetType().(*Message_EchoMessage); ok {
		return x.EchoMessage
	}
	return nil
}

func (x *Message) GetFinalMessage() *FinalMessage {
	if x, ok := x.GetType().(*Message_FinalMessage); ok {
		return x.FinalMessage
	}
	return nil
}

func (x *Message) GetDoneMessage() *DoneMessage {
	if x, ok := x.GetType().(*Message_DoneMessage); ok {
		return x.DoneMessage
	}
	return nil
}

type isMessage_Type interface {
	isMessage_Type()
}

type Message_SendMessage struct {
	SendMessage *SendMessage `protobuf:"bytes,1,opt,name=send_message,json=sendMessage,proto3,oneof"`
}

type Message_EchoMessage struct {
	EchoMessage *EchoMessage `protobuf:"bytes,2,opt,name=echo_message,json=echoMessage,proto3,oneof"`
}

type Message_FinalMessage struct {
	FinalMessage *FinalMessage `protobuf:"bytes,3,opt,name=final_message,json=finalMessage,proto3,oneof"`
}

type Message_DoneMessage struct {
	DoneMessage *DoneMessage `protobuf:"bytes,4,opt,name=done_message,json=doneMessage,proto3,oneof"`
}

func (*Message_SendMessage) isMessage_Type() {}

func (*Message_EchoMessage) isMessage_Type() {}

func (*Message_FinalMessage) isMessage_Type() {}

func (*Message_DoneMessage) isMessage_Type() {}

type SendMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txs []*trantorpb.Transaction `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
}

func (x *SendMessage) Reset() {
	*x = SendMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vcbpb_vcbpb_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessage) ProtoMessage() {}

func (x *SendMessage) ProtoReflect() protoreflect.Message {
	mi := &file_vcbpb_vcbpb_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessage.ProtoReflect.Descriptor instead.
func (*SendMessage) Descriptor() ([]byte, []int) {
	return file_vcbpb_vcbpb_proto_rawDescGZIP(), []int{6}
}

func (x *SendMessage) GetTxs() []*trantorpb.Transaction {
	if x != nil {
		return x.Txs
	}
	return nil
}

type EchoMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SignatureShare []byte `protobuf:"bytes,1,opt,name=signature_share,json=signatureShare,proto3" json:"signature_share,omitempty"`
}

func (x *EchoMessage) Reset() {
	*x = EchoMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vcbpb_vcbpb_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoMessage) ProtoMessage() {}

func (x *EchoMessage) ProtoReflect() protoreflect.Message {
	mi := &file_vcbpb_vcbpb_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoMessage.ProtoReflect.Descriptor instead.
func (*EchoMessage) Descriptor() ([]byte, []int) {
	return file_vcbpb_vcbpb_proto_rawDescGZIP(), []int{7}
}

func (x *EchoMessage) GetSignatureShare() []byte {
	if x != nil {
		return x.SignatureShare
	}
	return nil
}

type FinalMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signature []byte `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *FinalMessage) Reset() {
	*x = FinalMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vcbpb_vcbpb_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinalMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinalMessage) ProtoMessage() {}

func (x *FinalMessage) ProtoReflect() protoreflect.Message {
	mi := &file_vcbpb_vcbpb_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinalMessage.ProtoReflect.Descriptor instead.
func (*FinalMessage) Descriptor() ([]byte, []int) {
	return file_vcbpb_vcbpb_proto_rawDescGZIP(), []int{8}
}

func (x *FinalMessage) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type DoneMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DoneMessage) Reset() {
	*x = DoneMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vcbpb_vcbpb_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DoneMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoneMessage) ProtoMessage() {}

func (x *DoneMessage) ProtoReflect() protoreflect.Message {
	mi := &file_vcbpb_vcbpb_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoneMessage.ProtoReflect.Descriptor instead.
func (*DoneMessage) Descriptor() ([]byte, []int) {
	return file_vcbpb_vcbpb_proto_rawDescGZIP(), []int{9}
}

var File_vcbpb_vcbpb_proto protoreflect.FileDescriptor

var file_vcbpb_vcbpb_proto_rawDesc = []byte{
	0x0a, 0x11, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2f, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x76, 0x63, 0x62, 0x70, 0x62, 0x1a, 0x19, 0x74, 0x72, 0x61, 0x6e,
	0x74, 0x6f, 0x72, 0x70, 0x62, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67,
	0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6e, 0x65, 0x74, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e,
	0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xe0, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x0b, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x2a, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76,
	0x65, 0x72, 0x48, 0x00, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x12, 0x34, 0x0a,
	0x0b, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x5f, 0x64, 0x6f, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75,
	0x6d, 0x44, 0x6f, 0x6e, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x71, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x44,
	0x6f, 0x6e, 0x65, 0x12, 0x2b, 0x0a, 0x08, 0x61, 0x6c, 0x6c, 0x5f, 0x64, 0x6f, 0x6e, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x41, 0x6c,
	0x6c, 0x44, 0x6f, 0x6e, 0x65, 0x48, 0x00, 0x52, 0x07, 0x61, 0x6c, 0x6c, 0x44, 0x6f, 0x6e, 0x65,
	0x3a, 0x04, 0x90, 0xa6, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x04,
	0x80, 0xa6, 0x1d, 0x01, 0x22, 0x3c, 0x0a, 0x0a, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x28, 0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x74, 0x78, 0x73, 0x3a, 0x04, 0x98, 0xa6,
	0x1d, 0x01, 0x22, 0xc7, 0x02, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x12, 0x28,
	0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x72,
	0x61, 0x6e, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x03, 0x74, 0x78, 0x73, 0x12, 0x51, 0x0a, 0x06, 0x74, 0x78, 0x5f, 0x69,
	0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x42, 0x3a, 0x82, 0xa6, 0x1d, 0x36, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69,
	0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x74, 0x6f, 0x72, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e,
	0x54, 0x78, 0x49, 0x44, 0x52, 0x05, 0x74, 0x78, 0x49, 0x64, 0x73, 0x12, 0x62, 0x0a, 0x09, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x44,
	0x82, 0xa6, 0x1d, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x6f, 0x2f, 0x74, 0x63, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x46, 0x75, 0x6c,
	0x6c, 0x53, 0x69, 0x67, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12,
	0x55, 0x0a, 0x0a, 0x73, 0x72, 0x63, 0x5f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x36, 0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x09, 0x73, 0x72, 0x63,
	0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0x69, 0x0a, 0x0a,
	0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x44, 0x6f, 0x6e, 0x65, 0x12, 0x55, 0x0a, 0x0a, 0x73, 0x72,
	0x63, 0x5f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x36,
	0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x09, 0x73, 0x72, 0x63, 0x4d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0x66, 0x0a, 0x07, 0x41, 0x6c, 0x6c, 0x44, 0x6f,
	0x6e, 0x65, 0x12, 0x55, 0x0a, 0x0a, 0x73, 0x72, 0x63, 0x5f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x36, 0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x09,
	0x73, 0x72, 0x63, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22,
	0x84, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x37, 0x0a, 0x0c, 0x73,
	0x65, 0x6e, 0x64, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x37, 0x0a, 0x0c, 0x65, 0x63, 0x68, 0x6f, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x76, 0x63, 0x62,
	0x70, 0x62, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00,
	0x52, 0x0b, 0x65, 0x63, 0x68, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3a, 0x0a,
	0x0d, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e,
	0x61, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0c, 0x66, 0x69, 0x6e,
	0x61, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x37, 0x0a, 0x0c, 0x64, 0x6f, 0x6e,
	0x65, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x44, 0x6f, 0x6e, 0x65, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x6f, 0x6e, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x3a, 0x04, 0xc8, 0xe4, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x04, 0xc8, 0xe4, 0x1d, 0x01, 0x22, 0x3d, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x2e, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x74, 0x78, 0x73, 0x3a,
	0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x22, 0x83, 0x01, 0x0a, 0x0b, 0x45, 0x63, 0x68, 0x6f, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x6e, 0x0a, 0x0f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x45,
	0x82, 0xa6, 0x1d, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x6f, 0x2f, 0x74, 0x63, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x53, 0x69, 0x67,
	0x53, 0x68, 0x61, 0x72, 0x65, 0x52, 0x0e, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x53, 0x68, 0x61, 0x72, 0x65, 0x3a, 0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x22, 0x78, 0x0a, 0x0c, 0x46,
	0x69, 0x6e, 0x61, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x62, 0x0a, 0x09, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x44,
	0x82, 0xa6, 0x1d, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x6f, 0x2f, 0x74, 0x63, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x46, 0x75, 0x6c,
	0x6c, 0x53, 0x69, 0x67, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x3a,
	0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x22, 0x13, 0x0a, 0x0b, 0x44, 0x6f, 0x6e, 0x65, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x3a, 0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69,
	0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x70, 0x62, 0x2f, 0x76, 0x63, 0x62, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_vcbpb_vcbpb_proto_rawDescOnce sync.Once
	file_vcbpb_vcbpb_proto_rawDescData = file_vcbpb_vcbpb_proto_rawDesc
)

func file_vcbpb_vcbpb_proto_rawDescGZIP() []byte {
	file_vcbpb_vcbpb_proto_rawDescOnce.Do(func() {
		file_vcbpb_vcbpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_vcbpb_vcbpb_proto_rawDescData)
	})
	return file_vcbpb_vcbpb_proto_rawDescData
}

var file_vcbpb_vcbpb_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_vcbpb_vcbpb_proto_goTypes = []interface{}{
	(*Event)(nil),                 // 0: vcbpb.Event
	(*InputValue)(nil),            // 1: vcbpb.InputValue
	(*Deliver)(nil),               // 2: vcbpb.Deliver
	(*QuorumDone)(nil),            // 3: vcbpb.QuorumDone
	(*AllDone)(nil),               // 4: vcbpb.AllDone
	(*Message)(nil),               // 5: vcbpb.Message
	(*SendMessage)(nil),           // 6: vcbpb.SendMessage
	(*EchoMessage)(nil),           // 7: vcbpb.EchoMessage
	(*FinalMessage)(nil),          // 8: vcbpb.FinalMessage
	(*DoneMessage)(nil),           // 9: vcbpb.DoneMessage
	(*trantorpb.Transaction)(nil), // 10: trantorpb.Transaction
}
var file_vcbpb_vcbpb_proto_depIdxs = []int32{
	1,  // 0: vcbpb.Event.input_value:type_name -> vcbpb.InputValue
	2,  // 1: vcbpb.Event.deliver:type_name -> vcbpb.Deliver
	3,  // 2: vcbpb.Event.quorum_done:type_name -> vcbpb.QuorumDone
	4,  // 3: vcbpb.Event.all_done:type_name -> vcbpb.AllDone
	10, // 4: vcbpb.InputValue.txs:type_name -> trantorpb.Transaction
	10, // 5: vcbpb.Deliver.txs:type_name -> trantorpb.Transaction
	6,  // 6: vcbpb.Message.send_message:type_name -> vcbpb.SendMessage
	7,  // 7: vcbpb.Message.echo_message:type_name -> vcbpb.EchoMessage
	8,  // 8: vcbpb.Message.final_message:type_name -> vcbpb.FinalMessage
	9,  // 9: vcbpb.Message.done_message:type_name -> vcbpb.DoneMessage
	10, // 10: vcbpb.SendMessage.txs:type_name -> trantorpb.Transaction
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_vcbpb_vcbpb_proto_init() }
func file_vcbpb_vcbpb_proto_init() {
	if File_vcbpb_vcbpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_vcbpb_vcbpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_vcbpb_vcbpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_vcbpb_vcbpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_vcbpb_vcbpb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuorumDone); i {
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
		file_vcbpb_vcbpb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllDone); i {
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
		file_vcbpb_vcbpb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_vcbpb_vcbpb_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessage); i {
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
		file_vcbpb_vcbpb_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoMessage); i {
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
		file_vcbpb_vcbpb_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinalMessage); i {
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
		file_vcbpb_vcbpb_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DoneMessage); i {
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
	file_vcbpb_vcbpb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Event_InputValue)(nil),
		(*Event_Deliver)(nil),
		(*Event_QuorumDone)(nil),
		(*Event_AllDone)(nil),
	}
	file_vcbpb_vcbpb_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*Message_SendMessage)(nil),
		(*Message_EchoMessage)(nil),
		(*Message_FinalMessage)(nil),
		(*Message_DoneMessage)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vcbpb_vcbpb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_vcbpb_vcbpb_proto_goTypes,
		DependencyIndexes: file_vcbpb_vcbpb_proto_depIdxs,
		MessageInfos:      file_vcbpb_vcbpb_proto_msgTypes,
	}.Build()
	File_vcbpb_vcbpb_proto = out.File
	file_vcbpb_vcbpb_proto_rawDesc = nil
	file_vcbpb_vcbpb_proto_goTypes = nil
	file_vcbpb_vcbpb_proto_depIdxs = nil
}
