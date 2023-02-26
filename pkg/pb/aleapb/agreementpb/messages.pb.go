// separate file to avoid circular protobuf dependencies

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: aleapb/agreementpb/messages.proto

package agreementpb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	_ "github.com/filecoin-project/mir/pkg/pb/net"
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
	//
	//	*Message_FinishAbba
	Type isMessage_Type `protobuf_oneof:"type"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_messages_proto_msgTypes[0]
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
	return file_aleapb_agreementpb_messages_proto_rawDescGZIP(), []int{0}
}

func (m *Message) GetType() isMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Message) GetFinishAbba() *FinishAbbaMessage {
	if x, ok := x.GetType().(*Message_FinishAbba); ok {
		return x.FinishAbba
	}
	return nil
}

type isMessage_Type interface {
	isMessage_Type()
}

type Message_FinishAbba struct {
	FinishAbba *FinishAbbaMessage `protobuf:"bytes,1,opt,name=finish_abba,json=finishAbba,proto3,oneof"`
}

func (*Message_FinishAbba) isMessage_Type() {}

type FinishAbbaMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Round uint64 `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	Value bool   `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *FinishAbbaMessage) Reset() {
	*x = FinishAbbaMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_agreementpb_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinishAbbaMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishAbbaMessage) ProtoMessage() {}

func (x *FinishAbbaMessage) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_agreementpb_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishAbbaMessage.ProtoReflect.Descriptor instead.
func (*FinishAbbaMessage) Descriptor() ([]byte, []int) {
	return file_aleapb_agreementpb_messages_proto_rawDescGZIP(), []int{1}
}

func (x *FinishAbbaMessage) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *FinishAbbaMessage) GetValue() bool {
	if x != nil {
		return x.Value
	}
	return false
}

var File_aleapb_agreementpb_messages_proto protoreflect.FileDescriptor

var file_aleapb_agreementpb_messages_proto_rawDesc = []byte{
	0x0a, 0x21, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x70, 0x62, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x12, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x61, 0x67, 0x72, 0x65,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x1a, 0x1c, 0x6e, 0x65, 0x74, 0x2f, 0x63, 0x6f, 0x64,
	0x65, 0x67, 0x65, 0x6e, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x67, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x48, 0x0a, 0x0b, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x5f, 0x61, 0x62, 0x62, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x61,
	0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x69, 0x73,
	0x68, 0x41, 0x62, 0x62, 0x61, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0a,
	0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x41, 0x62, 0x62, 0x61, 0x3a, 0x04, 0xc8, 0xe4, 0x1d, 0x01,
	0x42, 0x0c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0xc8, 0xe4, 0x1d, 0x01, 0x22, 0x45,
	0x0a, 0x11, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x41, 0x62, 0x62, 0x61, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x04, 0xd0, 0xe4, 0x1d, 0x01, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f,
	0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f, 0x61, 0x67, 0x72, 0x65, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aleapb_agreementpb_messages_proto_rawDescOnce sync.Once
	file_aleapb_agreementpb_messages_proto_rawDescData = file_aleapb_agreementpb_messages_proto_rawDesc
)

func file_aleapb_agreementpb_messages_proto_rawDescGZIP() []byte {
	file_aleapb_agreementpb_messages_proto_rawDescOnce.Do(func() {
		file_aleapb_agreementpb_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_aleapb_agreementpb_messages_proto_rawDescData)
	})
	return file_aleapb_agreementpb_messages_proto_rawDescData
}

var file_aleapb_agreementpb_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_aleapb_agreementpb_messages_proto_goTypes = []interface{}{
	(*Message)(nil),           // 0: aleapb.agreementpb.Message
	(*FinishAbbaMessage)(nil), // 1: aleapb.agreementpb.FinishAbbaMessage
}
var file_aleapb_agreementpb_messages_proto_depIdxs = []int32{
	1, // 0: aleapb.agreementpb.Message.finish_abba:type_name -> aleapb.agreementpb.FinishAbbaMessage
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_aleapb_agreementpb_messages_proto_init() }
func file_aleapb_agreementpb_messages_proto_init() {
	if File_aleapb_agreementpb_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aleapb_agreementpb_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_agreementpb_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinishAbbaMessage); i {
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
	file_aleapb_agreementpb_messages_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Message_FinishAbba)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aleapb_agreementpb_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aleapb_agreementpb_messages_proto_goTypes,
		DependencyIndexes: file_aleapb_agreementpb_messages_proto_depIdxs,
		MessageInfos:      file_aleapb_agreementpb_messages_proto_msgTypes,
	}.Build()
	File_aleapb_agreementpb_messages_proto = out.File
	file_aleapb_agreementpb_messages_proto_rawDesc = nil
	file_aleapb_agreementpb_messages_proto_goTypes = nil
	file_aleapb_agreementpb_messages_proto_depIdxs = nil
}