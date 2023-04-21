package transportpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types2 "github.com/filecoin-project/mir/codegen/model/types"
	types "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	transportpb "github.com/filecoin-project/mir/pkg/pb/transportpb"
	types1 "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() transportpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb transportpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *transportpb.Event_SendMessage:
		return &Event_SendMessage{SendMessage: SendMessageFromPb(pb.SendMessage)}
	case *transportpb.Event_MessageReceived:
		return &Event_MessageReceived{MessageReceived: MessageReceivedFromPb(pb.MessageReceived)}
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

func (w *Event_SendMessage) Pb() transportpb.Event_Type {
	return &transportpb.Event_SendMessage{SendMessage: (w.SendMessage).Pb()}
}

func (*Event_SendMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*transportpb.Event_SendMessage]()}
}

type Event_MessageReceived struct {
	MessageReceived *MessageReceived
}

func (*Event_MessageReceived) isEvent_Type() {}

func (w *Event_MessageReceived) Unwrap() *MessageReceived {
	return w.MessageReceived
}

func (w *Event_MessageReceived) Pb() transportpb.Event_Type {
	return &transportpb.Event_MessageReceived{MessageReceived: (w.MessageReceived).Pb()}
}

func (*Event_MessageReceived) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*transportpb.Event_MessageReceived]()}
}

func EventFromPb(pb *transportpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *transportpb.Event {
	return &transportpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*transportpb.Event]()}
}

type SendMessage struct {
	Msg          *types.Message
	Destinations []types1.NodeID
}

func SendMessageFromPb(pb *transportpb.SendMessage) *SendMessage {
	return &SendMessage{
		Msg: types.MessageFromPb(pb.Msg),
		Destinations: types2.ConvertSlice(pb.Destinations, func(t string) types1.NodeID {
			return (types1.NodeID)(t)
		}),
	}
}

func (m *SendMessage) Pb() *transportpb.SendMessage {
	return &transportpb.SendMessage{
		Msg: (m.Msg).Pb(),
		Destinations: types2.ConvertSlice(m.Destinations, func(t types1.NodeID) string {
			return (string)(t)
		}),
	}
}

func (*SendMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*transportpb.SendMessage]()}
}

type MessageReceived struct {
	From types1.NodeID
	Msg  *types.Message
}

func MessageReceivedFromPb(pb *transportpb.MessageReceived) *MessageReceived {
	return &MessageReceived{
		From: (types1.NodeID)(pb.From),
		Msg:  types.MessageFromPb(pb.Msg),
	}
}

func (m *MessageReceived) Pb() *transportpb.MessageReceived {
	return &transportpb.MessageReceived{
		From: (string)(m.From),
		Msg:  (m.Msg).Pb(),
	}
}

func (*MessageReceived) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*transportpb.MessageReceived]()}
}
