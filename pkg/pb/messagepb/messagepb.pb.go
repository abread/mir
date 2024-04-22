//
//Copyright IBM Corp. All Rights Reserved.
//
//SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: messagepb/messagepb.proto

package messagepb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	abbapb "github.com/filecoin-project/mir/pkg/pb/abbapb"
	agreementpb "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	bcpb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	directorpb "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb"
	mscpb "github.com/filecoin-project/mir/pkg/pb/availabilitypb/mscpb"
	bcbpb "github.com/filecoin-project/mir/pkg/pb/bcbpb"
	checkpointpb "github.com/filecoin-project/mir/pkg/pb/checkpointpb"
	isspb "github.com/filecoin-project/mir/pkg/pb/isspb"
	_ "github.com/filecoin-project/mir/pkg/pb/mir"
	_ "github.com/filecoin-project/mir/pkg/pb/net"
	ordererpb "github.com/filecoin-project/mir/pkg/pb/ordererpb"
	pingpongpb "github.com/filecoin-project/mir/pkg/pb/pingpongpb"
	messages "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages"
	threshcheckpointpb "github.com/filecoin-project/mir/pkg/pb/threshcheckpointpb"
	vcbpb "github.com/filecoin-project/mir/pkg/pb/vcbpb"
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

	DestModule string `protobuf:"bytes,1,opt,name=dest_module,json=destModule,proto3" json:"dest_module,omitempty"`
	// Types that are assignable to Type:
	//
	//	*Message_Iss
	//	*Message_Bcb
	//	*Message_MultisigCollector
	//	*Message_Pingpong
	//	*Message_Checkpoint
	//	*Message_Orderer
	//	*Message_Vcb
	//	*Message_Abba
	//	*Message_AleaBroadcast
	//	*Message_AleaAgreement
	//	*Message_AleaDirector
	//	*Message_ReliableNet
	//	*Message_Threshcheckpoint
	Type isMessage_Type `protobuf_oneof:"type"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messagepb_messagepb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_messagepb_messagepb_proto_msgTypes[0]
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
	return file_messagepb_messagepb_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetDestModule() string {
	if x != nil {
		return x.DestModule
	}
	return ""
}

func (m *Message) GetType() isMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Message) GetIss() *isspb.ISSMessage {
	if x, ok := x.GetType().(*Message_Iss); ok {
		return x.Iss
	}
	return nil
}

func (x *Message) GetBcb() *bcbpb.Message {
	if x, ok := x.GetType().(*Message_Bcb); ok {
		return x.Bcb
	}
	return nil
}

func (x *Message) GetMultisigCollector() *mscpb.Message {
	if x, ok := x.GetType().(*Message_MultisigCollector); ok {
		return x.MultisigCollector
	}
	return nil
}

func (x *Message) GetPingpong() *pingpongpb.Message {
	if x, ok := x.GetType().(*Message_Pingpong); ok {
		return x.Pingpong
	}
	return nil
}

func (x *Message) GetCheckpoint() *checkpointpb.Message {
	if x, ok := x.GetType().(*Message_Checkpoint); ok {
		return x.Checkpoint
	}
	return nil
}

func (x *Message) GetOrderer() *ordererpb.Message {
	if x, ok := x.GetType().(*Message_Orderer); ok {
		return x.Orderer
	}
	return nil
}

func (x *Message) GetVcb() *vcbpb.Message {
	if x, ok := x.GetType().(*Message_Vcb); ok {
		return x.Vcb
	}
	return nil
}

func (x *Message) GetAbba() *abbapb.Message {
	if x, ok := x.GetType().(*Message_Abba); ok {
		return x.Abba
	}
	return nil
}

func (x *Message) GetAleaBroadcast() *bcpb.Message {
	if x, ok := x.GetType().(*Message_AleaBroadcast); ok {
		return x.AleaBroadcast
	}
	return nil
}

func (x *Message) GetAleaAgreement() *agreementpb.Message {
	if x, ok := x.GetType().(*Message_AleaAgreement); ok {
		return x.AleaAgreement
	}
	return nil
}

func (x *Message) GetAleaDirector() *directorpb.Message {
	if x, ok := x.GetType().(*Message_AleaDirector); ok {
		return x.AleaDirector
	}
	return nil
}

func (x *Message) GetReliableNet() *messages.Message {
	if x, ok := x.GetType().(*Message_ReliableNet); ok {
		return x.ReliableNet
	}
	return nil
}

func (x *Message) GetThreshcheckpoint() *threshcheckpointpb.Message {
	if x, ok := x.GetType().(*Message_Threshcheckpoint); ok {
		return x.Threshcheckpoint
	}
	return nil
}

type isMessage_Type interface {
	isMessage_Type()
}

type Message_Iss struct {
	Iss *isspb.ISSMessage `protobuf:"bytes,2,opt,name=iss,proto3,oneof"`
}

type Message_Bcb struct {
	Bcb *bcbpb.Message `protobuf:"bytes,3,opt,name=bcb,proto3,oneof"`
}

type Message_MultisigCollector struct {
	MultisigCollector *mscpb.Message `protobuf:"bytes,4,opt,name=multisig_collector,json=multisigCollector,proto3,oneof"`
}

type Message_Pingpong struct {
	Pingpong *pingpongpb.Message `protobuf:"bytes,5,opt,name=pingpong,proto3,oneof"`
}

type Message_Checkpoint struct {
	Checkpoint *checkpointpb.Message `protobuf:"bytes,6,opt,name=checkpoint,proto3,oneof"`
}

type Message_Orderer struct {
	Orderer *ordererpb.Message `protobuf:"bytes,7,opt,name=orderer,proto3,oneof"`
}

type Message_Vcb struct {
	Vcb *vcbpb.Message `protobuf:"bytes,8,opt,name=vcb,proto3,oneof"`
}

type Message_Abba struct {
	Abba *abbapb.Message `protobuf:"bytes,9,opt,name=abba,proto3,oneof"`
}

type Message_AleaBroadcast struct {
	AleaBroadcast *bcpb.Message `protobuf:"bytes,10,opt,name=alea_broadcast,json=aleaBroadcast,proto3,oneof"`
}

type Message_AleaAgreement struct {
	AleaAgreement *agreementpb.Message `protobuf:"bytes,11,opt,name=alea_agreement,json=aleaAgreement,proto3,oneof"`
}

type Message_AleaDirector struct {
	AleaDirector *directorpb.Message `protobuf:"bytes,14,opt,name=alea_director,json=aleaDirector,proto3,oneof"`
}

type Message_ReliableNet struct {
	ReliableNet *messages.Message `protobuf:"bytes,12,opt,name=reliable_net,json=reliableNet,proto3,oneof"`
}

type Message_Threshcheckpoint struct {
	Threshcheckpoint *threshcheckpointpb.Message `protobuf:"bytes,13,opt,name=threshcheckpoint,proto3,oneof"`
}

func (*Message_Iss) isMessage_Type() {}

func (*Message_Bcb) isMessage_Type() {}

func (*Message_MultisigCollector) isMessage_Type() {}

func (*Message_Pingpong) isMessage_Type() {}

func (*Message_Checkpoint) isMessage_Type() {}

func (*Message_Orderer) isMessage_Type() {}

func (*Message_Vcb) isMessage_Type() {}

func (*Message_Abba) isMessage_Type() {}

func (*Message_AleaBroadcast) isMessage_Type() {}

func (*Message_AleaAgreement) isMessage_Type() {}

func (*Message_AleaDirector) isMessage_Type() {}

func (*Message_ReliableNet) isMessage_Type() {}

func (*Message_Threshcheckpoint) isMessage_Type() {}

var File_messagepb_messagepb_proto protoreflect.FileDescriptor

var file_messagepb_messagepb_proto_rawDesc = []byte{
	0x0a, 0x19, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x70, 0x62, 0x2f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x70, 0x62, 0x1a, 0x11, 0x69, 0x73, 0x73, 0x70, 0x62, 0x2f, 0x69, 0x73,
	0x73, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x63, 0x62, 0x70, 0x62,
	0x2f, 0x62, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x61, 0x76,
	0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x70, 0x62, 0x2f, 0x6d, 0x73, 0x63,
	0x70, 0x62, 0x2f, 0x6d, 0x73, 0x63, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x70, 0x69, 0x6e, 0x67, 0x70, 0x6f, 0x6e, 0x67, 0x70, 0x62, 0x2f, 0x70, 0x69, 0x6e, 0x67, 0x70,
	0x6f, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x70, 0x62, 0x2f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x70, 0x62, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x65, 0x72, 0x70,
	0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2f, 0x76,
	0x63, 0x62, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x61, 0x62, 0x62, 0x61,
	0x70, 0x62, 0x2f, 0x61, 0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x16, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f, 0x62, 0x63, 0x70, 0x62, 0x2f, 0x62, 0x63, 0x70,
	0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f,
	0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x22, 0x61, 0x6c, 0x65, 0x61,
	0x70, 0x62, 0x2f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x2f, 0x64, 0x69,
	0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25,
	0x72, 0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x6e, 0x65, 0x74, 0x70, 0x62, 0x2f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2b, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x70, 0x62, 0x2f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68,
	0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f,
	0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1c, 0x6e, 0x65, 0x74, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc8,
	0x06, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x57, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x74, 0x5f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x36, 0x82, 0xa6, 0x1d, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x0a, 0x64, 0x65, 0x73, 0x74, 0x4d, 0x6f, 0x64,
	0x75, 0x6c, 0x65, 0x12, 0x25, 0x0a, 0x03, 0x69, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x69, 0x73, 0x73, 0x70, 0x62, 0x2e, 0x49, 0x53, 0x53, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x03, 0x69, 0x73, 0x73, 0x12, 0x22, 0x0a, 0x03, 0x62, 0x63,
	0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x62, 0x63, 0x62, 0x70, 0x62, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x03, 0x62, 0x63, 0x62, 0x12, 0x4e,
	0x0a, 0x12, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x73, 0x69, 0x67, 0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x61, 0x76, 0x61,
	0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x70, 0x62, 0x2e, 0x6d, 0x73, 0x63, 0x70,
	0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x11, 0x6d, 0x75, 0x6c,
	0x74, 0x69, 0x73, 0x69, 0x67, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x31,
	0x0a, 0x08, 0x70, 0x69, 0x6e, 0x67, 0x70, 0x6f, 0x6e, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x70, 0x69, 0x6e, 0x67, 0x70, 0x6f, 0x6e, 0x67, 0x70, 0x62, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x08, 0x70, 0x69, 0x6e, 0x67, 0x70, 0x6f, 0x6e,
	0x67, 0x12, 0x37, 0x0a, 0x0a, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0a,
	0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48,
	0x00, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x03, 0x76, 0x63,
	0x62, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x76, 0x63, 0x62, 0x70, 0x62, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x03, 0x76, 0x63, 0x62, 0x12, 0x25,
	0x0a, 0x04, 0x61, 0x62, 0x62, 0x61, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61,
	0x62, 0x62, 0x61, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52,
	0x04, 0x61, 0x62, 0x62, 0x61, 0x12, 0x3d, 0x0a, 0x0e, 0x61, 0x6c, 0x65, 0x61, 0x5f, 0x62, 0x72,
	0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x62, 0x63, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0d, 0x61, 0x6c, 0x65, 0x61, 0x42, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x12, 0x44, 0x0a, 0x0e, 0x61, 0x6c, 0x65, 0x61, 0x5f, 0x61, 0x67, 0x72,
	0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61,
	0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70,
	0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0d, 0x61, 0x6c, 0x65,
	0x61, 0x41, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x41, 0x0a, 0x0d, 0x61, 0x6c,
	0x65, 0x61, 0x5f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52,
	0x0c, 0x61, 0x6c, 0x65, 0x61, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x3b, 0x0a,
	0x0c, 0x72, 0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e, 0x65, 0x74, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x72, 0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x6e, 0x65,
	0x74, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x72,
	0x65, 0x6c, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x65, 0x74, 0x12, 0x49, 0x0a, 0x10, 0x74, 0x68,
	0x72, 0x65, 0x73, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x48, 0x00, 0x52, 0x10, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x3a, 0x04, 0xc0, 0xe4, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x04, 0xc8, 0xe4, 0x1d, 0x01, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e,
	0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x70, 0x62, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_messagepb_messagepb_proto_rawDescOnce sync.Once
	file_messagepb_messagepb_proto_rawDescData = file_messagepb_messagepb_proto_rawDesc
)

func file_messagepb_messagepb_proto_rawDescGZIP() []byte {
	file_messagepb_messagepb_proto_rawDescOnce.Do(func() {
		file_messagepb_messagepb_proto_rawDescData = protoimpl.X.CompressGZIP(file_messagepb_messagepb_proto_rawDescData)
	})
	return file_messagepb_messagepb_proto_rawDescData
}

var file_messagepb_messagepb_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_messagepb_messagepb_proto_goTypes = []interface{}{
	(*Message)(nil),                    // 0: messagepb.Message
	(*isspb.ISSMessage)(nil),           // 1: isspb.ISSMessage
	(*bcbpb.Message)(nil),              // 2: bcbpb.Message
	(*mscpb.Message)(nil),              // 3: availabilitypb.mscpb.Message
	(*pingpongpb.Message)(nil),         // 4: pingpongpb.Message
	(*checkpointpb.Message)(nil),       // 5: checkpointpb.Message
	(*ordererpb.Message)(nil),          // 6: ordererpb.Message
	(*vcbpb.Message)(nil),              // 7: vcbpb.Message
	(*abbapb.Message)(nil),             // 8: abbapb.Message
	(*bcpb.Message)(nil),               // 9: aleapb.bcpb.Message
	(*agreementpb.Message)(nil),        // 10: aleapb.agreementpb.Message
	(*directorpb.Message)(nil),         // 11: aleapb.directorpb.Message
	(*messages.Message)(nil),           // 12: reliablenetpb.Message
	(*threshcheckpointpb.Message)(nil), // 13: threshcheckpointpb.Message
}
var file_messagepb_messagepb_proto_depIdxs = []int32{
	1,  // 0: messagepb.Message.iss:type_name -> isspb.ISSMessage
	2,  // 1: messagepb.Message.bcb:type_name -> bcbpb.Message
	3,  // 2: messagepb.Message.multisig_collector:type_name -> availabilitypb.mscpb.Message
	4,  // 3: messagepb.Message.pingpong:type_name -> pingpongpb.Message
	5,  // 4: messagepb.Message.checkpoint:type_name -> checkpointpb.Message
	6,  // 5: messagepb.Message.orderer:type_name -> ordererpb.Message
	7,  // 6: messagepb.Message.vcb:type_name -> vcbpb.Message
	8,  // 7: messagepb.Message.abba:type_name -> abbapb.Message
	9,  // 8: messagepb.Message.alea_broadcast:type_name -> aleapb.bcpb.Message
	10, // 9: messagepb.Message.alea_agreement:type_name -> aleapb.agreementpb.Message
	11, // 10: messagepb.Message.alea_director:type_name -> aleapb.directorpb.Message
	12, // 11: messagepb.Message.reliable_net:type_name -> reliablenetpb.Message
	13, // 12: messagepb.Message.threshcheckpoint:type_name -> threshcheckpointpb.Message
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_messagepb_messagepb_proto_init() }
func file_messagepb_messagepb_proto_init() {
	if File_messagepb_messagepb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messagepb_messagepb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
	}
	file_messagepb_messagepb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Message_Iss)(nil),
		(*Message_Bcb)(nil),
		(*Message_MultisigCollector)(nil),
		(*Message_Pingpong)(nil),
		(*Message_Checkpoint)(nil),
		(*Message_Orderer)(nil),
		(*Message_Vcb)(nil),
		(*Message_Abba)(nil),
		(*Message_AleaBroadcast)(nil),
		(*Message_AleaAgreement)(nil),
		(*Message_AleaDirector)(nil),
		(*Message_ReliableNet)(nil),
		(*Message_Threshcheckpoint)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messagepb_messagepb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messagepb_messagepb_proto_goTypes,
		DependencyIndexes: file_messagepb_messagepb_proto_depIdxs,
		MessageInfos:      file_messagepb_messagepb_proto_msgTypes,
	}.Build()
	File_messagepb_messagepb_proto = out.File
	file_messagepb_messagepb_proto_rawDesc = nil
	file_messagepb_messagepb_proto_goTypes = nil
	file_messagepb_messagepb_proto_depIdxs = nil
}
