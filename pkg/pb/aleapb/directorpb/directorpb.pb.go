// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1
// 	protoc        v3.21.12
// source: aleapb/directorpb/directorpb.proto

package directorpb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	common "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
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
	//	*Event_Heartbeat
	//	*Event_FillGap
	//	*Event_Stats
	Type isEvent_Type `protobuf_oneof:"type"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_directorpb_directorpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_directorpb_directorpb_proto_msgTypes[0]
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
	return file_aleapb_directorpb_directorpb_proto_rawDescGZIP(), []int{0}
}

func (m *Event) GetType() isEvent_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Event) GetHeartbeat() *Heartbeat {
	if x, ok := x.GetType().(*Event_Heartbeat); ok {
		return x.Heartbeat
	}
	return nil
}

func (x *Event) GetFillGap() *DoFillGap {
	if x, ok := x.GetType().(*Event_FillGap); ok {
		return x.FillGap
	}
	return nil
}

func (x *Event) GetStats() *Stats {
	if x, ok := x.GetType().(*Event_Stats); ok {
		return x.Stats
	}
	return nil
}

type isEvent_Type interface {
	isEvent_Type()
}

type Event_Heartbeat struct {
	Heartbeat *Heartbeat `protobuf:"bytes,1,opt,name=heartbeat,proto3,oneof"`
}

type Event_FillGap struct {
	FillGap *DoFillGap `protobuf:"bytes,2,opt,name=fill_gap,json=fillGap,proto3,oneof"`
}

type Event_Stats struct {
	Stats *Stats `protobuf:"bytes,3,opt,name=stats,proto3,oneof"`
}

func (*Event_Heartbeat) isEvent_Type() {}

func (*Event_FillGap) isEvent_Type() {}

func (*Event_Stats) isEvent_Type() {}

type Heartbeat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Heartbeat) Reset() {
	*x = Heartbeat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_directorpb_directorpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Heartbeat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Heartbeat) ProtoMessage() {}

func (x *Heartbeat) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_directorpb_directorpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Heartbeat.ProtoReflect.Descriptor instead.
func (*Heartbeat) Descriptor() ([]byte, []int) {
	return file_aleapb_directorpb_directorpb_proto_rawDescGZIP(), []int{1}
}

type DoFillGap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Slot *common.Slot `protobuf:"bytes,1,opt,name=slot,proto3" json:"slot,omitempty"`
}

func (x *DoFillGap) Reset() {
	*x = DoFillGap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_directorpb_directorpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DoFillGap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoFillGap) ProtoMessage() {}

func (x *DoFillGap) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_directorpb_directorpb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoFillGap.ProtoReflect.Descriptor instead.
func (*DoFillGap) Descriptor() ([]byte, []int) {
	return file_aleapb_directorpb_directorpb_proto_rawDescGZIP(), []int{2}
}

func (x *DoFillGap) GetSlot() *common.Slot {
	if x != nil {
		return x.Slot
	}
	return nil
}

type Stats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SlotsWaitingDelivery uint64 `protobuf:"varint,1,opt,name=slots_waiting_delivery,json=slotsWaitingDelivery,proto3" json:"slots_waiting_delivery,omitempty"`
	MinAgDurationEst     int64  `protobuf:"varint,2,opt,name=min_ag_duration_est,json=minAgDurationEst,proto3" json:"min_ag_duration_est,omitempty"`
	MaxAgDurationEst     int64  `protobuf:"varint,3,opt,name=max_ag_duration_est,json=maxAgDurationEst,proto3" json:"max_ag_duration_est,omitempty"`
	MinBcDurationEst     int64  `protobuf:"varint,4,opt,name=min_bc_duration_est,json=minBcDurationEst,proto3" json:"min_bc_duration_est,omitempty"`
	MaxBcDurationEst     int64  `protobuf:"varint,5,opt,name=max_bc_duration_est,json=maxBcDurationEst,proto3" json:"max_bc_duration_est,omitempty"`
	MinOwnBcDurationEst  int64  `protobuf:"varint,6,opt,name=min_own_bc_duration_est,json=minOwnBcDurationEst,proto3" json:"min_own_bc_duration_est,omitempty"`
	MaxOwnBcDurationEst  int64  `protobuf:"varint,7,opt,name=max_own_bc_duration_est,json=maxOwnBcDurationEst,proto3" json:"max_own_bc_duration_est,omitempty"`
	BcEstMargin          int64  `protobuf:"varint,8,opt,name=bc_est_margin,json=bcEstMargin,proto3" json:"bc_est_margin,omitempty"`
}

func (x *Stats) Reset() {
	*x = Stats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aleapb_directorpb_directorpb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stats) ProtoMessage() {}

func (x *Stats) ProtoReflect() protoreflect.Message {
	mi := &file_aleapb_directorpb_directorpb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stats.ProtoReflect.Descriptor instead.
func (*Stats) Descriptor() ([]byte, []int) {
	return file_aleapb_directorpb_directorpb_proto_rawDescGZIP(), []int{3}
}

func (x *Stats) GetSlotsWaitingDelivery() uint64 {
	if x != nil {
		return x.SlotsWaitingDelivery
	}
	return 0
}

func (x *Stats) GetMinAgDurationEst() int64 {
	if x != nil {
		return x.MinAgDurationEst
	}
	return 0
}

func (x *Stats) GetMaxAgDurationEst() int64 {
	if x != nil {
		return x.MaxAgDurationEst
	}
	return 0
}

func (x *Stats) GetMinBcDurationEst() int64 {
	if x != nil {
		return x.MinBcDurationEst
	}
	return 0
}

func (x *Stats) GetMaxBcDurationEst() int64 {
	if x != nil {
		return x.MaxBcDurationEst
	}
	return 0
}

func (x *Stats) GetMinOwnBcDurationEst() int64 {
	if x != nil {
		return x.MinOwnBcDurationEst
	}
	return 0
}

func (x *Stats) GetMaxOwnBcDurationEst() int64 {
	if x != nil {
		return x.MaxOwnBcDurationEst
	}
	return 0
}

func (x *Stats) GetBcEstMargin() int64 {
	if x != nil {
		return x.BcEstMargin
	}
	return 0
}

var File_aleapb_directorpb_directorpb_proto protoreflect.FileDescriptor

var file_aleapb_directorpb_directorpb_proto_rawDesc = []byte{
	0x0a, 0x22, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x70, 0x62, 0x2f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x1a, 0x1a, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6d, 0x69, 0x72, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65, 0x6e,
	0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xc6, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x3c, 0x0a, 0x09, 0x68,
	0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x70, 0x62, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x48, 0x00, 0x52, 0x09,
	0x68, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x39, 0x0a, 0x08, 0x66, 0x69, 0x6c,
	0x6c, 0x5f, 0x67, 0x61, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x6c,
	0x65, 0x61, 0x70, 0x62, 0x2e, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x2e,
	0x44, 0x6f, 0x46, 0x69, 0x6c, 0x6c, 0x47, 0x61, 0x70, 0x48, 0x00, 0x52, 0x07, 0x66, 0x69, 0x6c,
	0x6c, 0x47, 0x61, 0x70, 0x12, 0x30, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70, 0x62, 0x2e, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x48, 0x00, 0x52,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x3a, 0x04, 0x90, 0xa6, 0x1d, 0x01, 0x42, 0x0c, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x04, 0x80, 0xa6, 0x1d, 0x01, 0x22, 0x11, 0x0a, 0x09, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0x3a, 0x0a,
	0x09, 0x44, 0x6f, 0x46, 0x69, 0x6c, 0x6c, 0x47, 0x61, 0x70, 0x12, 0x27, 0x0a, 0x04, 0x73, 0x6c,
	0x6f, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61, 0x6c, 0x65, 0x61, 0x70,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x6c, 0x6f, 0x74, 0x52, 0x04, 0x73,
	0x6c, 0x6f, 0x74, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01, 0x22, 0x94, 0x04, 0x0a, 0x05, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x12, 0x34, 0x0a, 0x16, 0x73, 0x6c, 0x6f, 0x74, 0x73, 0x5f, 0x77, 0x61, 0x69,
	0x74, 0x69, 0x6e, 0x67, 0x5f, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x14, 0x73, 0x6c, 0x6f, 0x74, 0x73, 0x57, 0x61, 0x69, 0x74, 0x69, 0x6e,
	0x67, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x12, 0x40, 0x0a, 0x13, 0x6d, 0x69, 0x6e,
	0x5f, 0x61, 0x67, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x73, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x82, 0xa6, 0x1d, 0x0d, 0x74, 0x69, 0x6d, 0x65,
	0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x6d, 0x69, 0x6e, 0x41, 0x67,
	0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x73, 0x74, 0x12, 0x40, 0x0a, 0x13, 0x6d,
	0x61, 0x78, 0x5f, 0x61, 0x67, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65,
	0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x82, 0xa6, 0x1d, 0x0d, 0x74, 0x69,
	0x6d, 0x65, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x6d, 0x61, 0x78,
	0x41, 0x67, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x73, 0x74, 0x12, 0x40, 0x0a,
	0x13, 0x6d, 0x69, 0x6e, 0x5f, 0x62, 0x63, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x65, 0x73, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x82, 0xa6, 0x1d, 0x0d,
	0x74, 0x69, 0x6d, 0x65, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x6d,
	0x69, 0x6e, 0x42, 0x63, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x73, 0x74, 0x12,
	0x40, 0x0a, 0x13, 0x6d, 0x61, 0x78, 0x5f, 0x62, 0x63, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x65, 0x73, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x82, 0xa6,
	0x1d, 0x0d, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x10, 0x6d, 0x61, 0x78, 0x42, 0x63, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x73,
	0x74, 0x12, 0x47, 0x0a, 0x17, 0x6d, 0x69, 0x6e, 0x5f, 0x6f, 0x77, 0x6e, 0x5f, 0x62, 0x63, 0x5f,
	0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x73, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x42, 0x11, 0x82, 0xa6, 0x1d, 0x0d, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x44, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x13, 0x6d, 0x69, 0x6e, 0x4f, 0x77, 0x6e, 0x42, 0x63, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x73, 0x74, 0x12, 0x47, 0x0a, 0x17, 0x6d, 0x61,
	0x78, 0x5f, 0x6f, 0x77, 0x6e, 0x5f, 0x62, 0x63, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x65, 0x73, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x82, 0xa6, 0x1d,
	0x0d, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x13,
	0x6d, 0x61, 0x78, 0x4f, 0x77, 0x6e, 0x42, 0x63, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x45, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x0d, 0x62, 0x63, 0x5f, 0x65, 0x73, 0x74, 0x5f, 0x6d, 0x61,
	0x72, 0x67, 0x69, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x82, 0xa6, 0x1d, 0x0d,
	0x74, 0x69, 0x6d, 0x65, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x62,
	0x63, 0x45, 0x73, 0x74, 0x4d, 0x61, 0x72, 0x67, 0x69, 0x6e, 0x3a, 0x04, 0x98, 0xa6, 0x1d, 0x01,
	0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x61, 0x6c, 0x65, 0x61, 0x70,
	0x62, 0x2f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aleapb_directorpb_directorpb_proto_rawDescOnce sync.Once
	file_aleapb_directorpb_directorpb_proto_rawDescData = file_aleapb_directorpb_directorpb_proto_rawDesc
)

func file_aleapb_directorpb_directorpb_proto_rawDescGZIP() []byte {
	file_aleapb_directorpb_directorpb_proto_rawDescOnce.Do(func() {
		file_aleapb_directorpb_directorpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_aleapb_directorpb_directorpb_proto_rawDescData)
	})
	return file_aleapb_directorpb_directorpb_proto_rawDescData
}

var file_aleapb_directorpb_directorpb_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_aleapb_directorpb_directorpb_proto_goTypes = []interface{}{
	(*Event)(nil),       // 0: aleapb.directorpb.Event
	(*Heartbeat)(nil),   // 1: aleapb.directorpb.Heartbeat
	(*DoFillGap)(nil),   // 2: aleapb.directorpb.DoFillGap
	(*Stats)(nil),       // 3: aleapb.directorpb.Stats
	(*common.Slot)(nil), // 4: aleapb.common.Slot
}
var file_aleapb_directorpb_directorpb_proto_depIdxs = []int32{
	1, // 0: aleapb.directorpb.Event.heartbeat:type_name -> aleapb.directorpb.Heartbeat
	2, // 1: aleapb.directorpb.Event.fill_gap:type_name -> aleapb.directorpb.DoFillGap
	3, // 2: aleapb.directorpb.Event.stats:type_name -> aleapb.directorpb.Stats
	4, // 3: aleapb.directorpb.DoFillGap.slot:type_name -> aleapb.common.Slot
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_aleapb_directorpb_directorpb_proto_init() }
func file_aleapb_directorpb_directorpb_proto_init() {
	if File_aleapb_directorpb_directorpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aleapb_directorpb_directorpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_aleapb_directorpb_directorpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Heartbeat); i {
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
		file_aleapb_directorpb_directorpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DoFillGap); i {
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
		file_aleapb_directorpb_directorpb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stats); i {
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
	file_aleapb_directorpb_directorpb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Event_Heartbeat)(nil),
		(*Event_FillGap)(nil),
		(*Event_Stats)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_aleapb_directorpb_directorpb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aleapb_directorpb_directorpb_proto_goTypes,
		DependencyIndexes: file_aleapb_directorpb_directorpb_proto_depIdxs,
		MessageInfos:      file_aleapb_directorpb_directorpb_proto_msgTypes,
	}.Build()
	File_aleapb_directorpb_directorpb_proto = out.File
	file_aleapb_directorpb_directorpb_proto_rawDesc = nil
	file_aleapb_directorpb_directorpb_proto_goTypes = nil
	file_aleapb_directorpb_directorpb_proto_depIdxs = nil
}
