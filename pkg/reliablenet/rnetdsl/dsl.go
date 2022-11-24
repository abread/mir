package rnetdsl

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/reliablenetpb"
	rnetpbMessages "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages"
	rnetEvents "github.com/filecoin-project/mir/pkg/reliablenet/events"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func SendMessage(m dsl.Module, destModule t.ModuleID, id []byte, msg *messagepb.Message, destinations []t.NodeID) {
	dsl.EmitEvent(m, rnetEvents.SendMessage(destModule, id, msg, destinations))
}

func MarkModuleMsgsRecvd(m dsl.Module, destModule t.ModuleID, msgDestModule t.ModuleID, destinations []t.NodeID) {
	dsl.EmitEvent(m, rnetEvents.MarkModuleMsgsRecvd(destModule, msgDestModule, destinations))
}

func MarkRecvd(m dsl.Module, destModule t.ModuleID, msgDestModule t.ModuleID, msgID []byte, destinations []t.NodeID) {
	dsl.EmitEvent(m, rnetEvents.MarkRecvd(destModule, msgDestModule, msgID, destinations))
}

func RetransmitAll(m dsl.Module, destModule t.ModuleID) {
	dsl.EmitEvent(m, rnetEvents.RetransmitAll(destModule))
}

func Ack(m dsl.Module, destModule t.ModuleID, msgDestModule t.ModuleID, msgID []byte, msgSource t.NodeID) {
	dsl.EmitEvent(m, rnetEvents.Ack(destModule, msgDestModule, msgID, msgSource))
}

// Module-specific dsl functions for processing events.

func UponEvent[EvWrapper reliablenetpb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*eventpb.Event_ReliableNet](m, func(ev *reliablenetpb.Event) error {
		evWrapper, ok := ev.Type.(EvWrapper)
		if !ok {
			return nil
		}
		return handler(evWrapper.Unwrap())
	})
}

func UponSendMessage(m dsl.Module, handler func(id []byte, msg *messagepb.Message, destinations []t.NodeID) error) {
	UponEvent[*reliablenetpb.Event_SendMessage](m, func(ev *reliablenetpb.SendMessage) error {
		return handler(ev.MsgId, ev.Msg, t.NodeIDSlice(ev.Destinations))
	})
}

func UponAck(m dsl.Module, handler func(msgDestModule t.ModuleID, msgID []byte, msgSource t.NodeID) error) {
	UponEvent[*reliablenetpb.Event_Ack](m, func(ev *reliablenetpb.Ack) error {
		return handler(t.ModuleID(ev.DestModule), ev.MsgId, t.NodeID(ev.Source))
	})
}

func UponMarkModuleMsgsRecvd(m dsl.Module, handler func(msgDestModule t.ModuleID, destinations []t.NodeID) error) {
	UponEvent[*reliablenetpb.Event_MarkModuleMsgsRecvd](m, func(ev *reliablenetpb.MarkModuleMsgsRecvd) error {
		return handler(t.ModuleID(ev.DestModule), t.NodeIDSlice(ev.Destinations))
	})
}

func UponMarkRecvd(m dsl.Module, handler func(msgDestModule t.ModuleID, msgId []byte, destination []t.NodeID) error) {
	UponEvent[*reliablenetpb.Event_MarkRecvd](m, func(ev *reliablenetpb.MarkRecvd) error {
		return handler(t.ModuleID(ev.DestModule), ev.MsgId, t.NodeIDSlice(ev.Destinations))
	})
}

func UponRetransmitAll(m dsl.Module, handler func() error) {
	UponEvent[*reliablenetpb.Event_RetransmitAll](m, func(ev *reliablenetpb.RetransmitAll) error {
		return handler()
	})
}

func UponMessageReceived(m dsl.Module, handler func(from t.NodeID, msg *rnetpbMessages.Message) error) {
	dsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepb.Message) error {
		cbMsgWrapper, ok := msg.Type.(*messagepb.Message_ReliableNet)
		if !ok {
			return nil
		}

		return handler(from, cbMsgWrapper.ReliableNet)
	})
}

func UponAckMessageReceived(m dsl.Module, handler func(from t.NodeID, msgDestModule t.ModuleID, msgId []byte) error) {
	UponMessageReceived(m, func(from t.NodeID, msg *rnetpbMessages.Message) error {
		msgWrapper, ok := msg.Type.(*rnetpbMessages.Message_Ack)
		if !ok {
			return nil
		}

		return handler(from, t.ModuleID(msgWrapper.Ack.MsgDestModule), msgWrapper.Ack.MsgId)
	})
}
