package reliablenetpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types2 "github.com/filecoin-project/mir/codegen/model/types"
	types "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	reliablenetpb "github.com/filecoin-project/mir/pkg/pb/reliablenetpb"
	types1 "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() reliablenetpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb reliablenetpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *reliablenetpb.Event_SendMessage:
		return &Event_SendMessage{SendMessage: SendMessageFromPb(pb.SendMessage)}
	case *reliablenetpb.Event_Ack:
		return &Event_Ack{Ack: AckFromPb(pb.Ack)}
	case *reliablenetpb.Event_MarkRecvd:
		return &Event_MarkRecvd{MarkRecvd: MarkRecvdFromPb(pb.MarkRecvd)}
	case *reliablenetpb.Event_MarkModuleMsgsRecvd:
		return &Event_MarkModuleMsgsRecvd{MarkModuleMsgsRecvd: MarkModuleMsgsRecvdFromPb(pb.MarkModuleMsgsRecvd)}
	case *reliablenetpb.Event_RetransmitAll:
		return &Event_RetransmitAll{RetransmitAll: RetransmitAllFromPb(pb.RetransmitAll)}
	}
	return nil
}

type Event_SendMessage struct {
	SendMessage *SendMessage
}

func (*Event_SendMessage) isEvent_Type() {}

func (w *Event_SendMessage) Unwrap() *SendMessage {
	return w.SendMessage
}

func (w *Event_SendMessage) Pb() reliablenetpb.Event_Type {
	return &reliablenetpb.Event_SendMessage{SendMessage: (w.SendMessage).Pb()}
}

func (*Event_SendMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.Event_SendMessage]()}
}

type Event_Ack struct {
	Ack *Ack
}

func (*Event_Ack) isEvent_Type() {}

func (w *Event_Ack) Unwrap() *Ack {
	return w.Ack
}

func (w *Event_Ack) Pb() reliablenetpb.Event_Type {
	return &reliablenetpb.Event_Ack{Ack: (w.Ack).Pb()}
}

func (*Event_Ack) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.Event_Ack]()}
}

type Event_MarkRecvd struct {
	MarkRecvd *MarkRecvd
}

func (*Event_MarkRecvd) isEvent_Type() {}

func (w *Event_MarkRecvd) Unwrap() *MarkRecvd {
	return w.MarkRecvd
}

func (w *Event_MarkRecvd) Pb() reliablenetpb.Event_Type {
	return &reliablenetpb.Event_MarkRecvd{MarkRecvd: (w.MarkRecvd).Pb()}
}

func (*Event_MarkRecvd) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.Event_MarkRecvd]()}
}

type Event_MarkModuleMsgsRecvd struct {
	MarkModuleMsgsRecvd *MarkModuleMsgsRecvd
}

func (*Event_MarkModuleMsgsRecvd) isEvent_Type() {}

func (w *Event_MarkModuleMsgsRecvd) Unwrap() *MarkModuleMsgsRecvd {
	return w.MarkModuleMsgsRecvd
}

func (w *Event_MarkModuleMsgsRecvd) Pb() reliablenetpb.Event_Type {
	return &reliablenetpb.Event_MarkModuleMsgsRecvd{MarkModuleMsgsRecvd: (w.MarkModuleMsgsRecvd).Pb()}
}

func (*Event_MarkModuleMsgsRecvd) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.Event_MarkModuleMsgsRecvd]()}
}

type Event_RetransmitAll struct {
	RetransmitAll *RetransmitAll
}

func (*Event_RetransmitAll) isEvent_Type() {}

func (w *Event_RetransmitAll) Unwrap() *RetransmitAll {
	return w.RetransmitAll
}

func (w *Event_RetransmitAll) Pb() reliablenetpb.Event_Type {
	return &reliablenetpb.Event_RetransmitAll{RetransmitAll: (w.RetransmitAll).Pb()}
}

func (*Event_RetransmitAll) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.Event_RetransmitAll]()}
}

func EventFromPb(pb *reliablenetpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *reliablenetpb.Event {
	return &reliablenetpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.Event]()}
}

type SendMessage struct {
	MsgId        []uint8
	Msg          *types.Message
	Destinations []types1.NodeID
}

func SendMessageFromPb(pb *reliablenetpb.SendMessage) *SendMessage {
	return &SendMessage{
		MsgId: pb.MsgId,
		Msg:   types.MessageFromPb(pb.Msg),
		Destinations: types2.ConvertSlice(pb.Destinations, func(t string) types1.NodeID {
			return (types1.NodeID)(t)
		}),
	}
}

func (m *SendMessage) Pb() *reliablenetpb.SendMessage {
	return &reliablenetpb.SendMessage{
		MsgId: m.MsgId,
		Msg:   (m.Msg).Pb(),
		Destinations: types2.ConvertSlice(m.Destinations, func(t types1.NodeID) string {
			return (string)(t)
		}),
	}
}

func (*SendMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.SendMessage]()}
}

type Ack struct {
	DestModule types1.ModuleID
	MsgId      []uint8
	Source     types1.NodeID
}

func AckFromPb(pb *reliablenetpb.Ack) *Ack {
	return &Ack{
		DestModule: (types1.ModuleID)(pb.DestModule),
		MsgId:      pb.MsgId,
		Source:     (types1.NodeID)(pb.Source),
	}
}

func (m *Ack) Pb() *reliablenetpb.Ack {
	return &reliablenetpb.Ack{
		DestModule: (string)(m.DestModule),
		MsgId:      m.MsgId,
		Source:     (string)(m.Source),
	}
}

func (*Ack) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.Ack]()}
}

type MarkModuleMsgsRecvd struct {
	DestModule   types1.ModuleID
	Destinations []types1.NodeID
}

func MarkModuleMsgsRecvdFromPb(pb *reliablenetpb.MarkModuleMsgsRecvd) *MarkModuleMsgsRecvd {
	return &MarkModuleMsgsRecvd{
		DestModule: (types1.ModuleID)(pb.DestModule),
		Destinations: types2.ConvertSlice(pb.Destinations, func(t string) types1.NodeID {
			return (types1.NodeID)(t)
		}),
	}
}

func (m *MarkModuleMsgsRecvd) Pb() *reliablenetpb.MarkModuleMsgsRecvd {
	return &reliablenetpb.MarkModuleMsgsRecvd{
		DestModule: (string)(m.DestModule),
		Destinations: types2.ConvertSlice(m.Destinations, func(t types1.NodeID) string {
			return (string)(t)
		}),
	}
}

func (*MarkModuleMsgsRecvd) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.MarkModuleMsgsRecvd]()}
}

type MarkRecvd struct {
	DestModule   types1.ModuleID
	MsgId        []uint8
	Destinations []types1.NodeID
}

func MarkRecvdFromPb(pb *reliablenetpb.MarkRecvd) *MarkRecvd {
	return &MarkRecvd{
		DestModule: (types1.ModuleID)(pb.DestModule),
		MsgId:      pb.MsgId,
		Destinations: types2.ConvertSlice(pb.Destinations, func(t string) types1.NodeID {
			return (types1.NodeID)(t)
		}),
	}
}

func (m *MarkRecvd) Pb() *reliablenetpb.MarkRecvd {
	return &reliablenetpb.MarkRecvd{
		DestModule: (string)(m.DestModule),
		MsgId:      m.MsgId,
		Destinations: types2.ConvertSlice(m.Destinations, func(t types1.NodeID) string {
			return (string)(t)
		}),
	}
}

func (*MarkRecvd) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.MarkRecvd]()}
}

type RetransmitAll struct{}

func RetransmitAllFromPb(pb *reliablenetpb.RetransmitAll) *RetransmitAll {
	return &RetransmitAll{}
}

func (m *RetransmitAll) Pb() *reliablenetpb.RetransmitAll {
	return &reliablenetpb.RetransmitAll{}
}

func (*RetransmitAll) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*reliablenetpb.RetransmitAll]()}
}
