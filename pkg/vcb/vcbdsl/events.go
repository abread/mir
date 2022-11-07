package vcbdsl

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	"github.com/filecoin-project/mir/pkg/pb/vcbpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func Event(m dsl.Module, dest t.ModuleID, ev *vcbpb.Event) {
	dsl.EmitEvent(m, &eventpb.Event{
		DestModule: dest.Pb(),

		Type: &eventpb.Event_Vcb{
			Vcb: ev,
		},
	})
}

func Request(m dsl.Module, dest t.ModuleID, data []*requestpb.Request) {
	Event(m, dest, &vcbpb.Event{
		Type: &vcbpb.Event_Request{
			Request: &vcbpb.BroadcastRequest{
				Data: data,
			},
		},
	})
}

func Deliver(m dsl.Module, dest t.ModuleID, data []*requestpb.Request, batchID t.BatchID, signature []byte) {
	Event(m, dest, &vcbpb.Event{
		Type: &vcbpb.Event_Deliver{
			Deliver: &vcbpb.Deliver{
				Data:         data,
				BatchId:      batchID.Pb(),
				Signature:    signature,
				OriginModule: string(m.ModuleID()),
			},
		},
	})
}

// Module-specific dsl functions for processing events.

func UponEvent[EvWrapper vcbpb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*eventpb.Event_Vcb](m, func(ev *vcbpb.Event) error {
		evWrapper, ok := ev.Type.(EvWrapper)
		if !ok {
			return nil
		}
		return handler(evWrapper.Unwrap())
	})
}

func UponBroadcastRequest(m dsl.Module, handler func(data []*requestpb.Request) error) {
	UponEvent[*vcbpb.Event_Request](m, func(ev *vcbpb.BroadcastRequest) error {
		return handler(ev.Data)
	})
}

func UponDeliver(m dsl.Module, handler func(data []*requestpb.Request, batchID t.BatchID, signature []byte, from t.ModuleID) error) {
	UponEvent[*vcbpb.Event_Deliver](m, func(ev *vcbpb.Deliver) error {
		return handler(ev.Data, t.BatchID(ev.BatchId), ev.Signature, t.ModuleID(ev.OriginModule))
	})
}

func UponVcbMessageReceived(m dsl.Module, handler func(from t.NodeID, msg *vcbpb.Message) error) {
	dsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepb.Message) error {
		cbMsgWrapper, ok := msg.Type.(*messagepb.Message_Vcb)
		if !ok {
			return nil
		}

		return handler(from, cbMsgWrapper.Vcb)
	})
}

func UponSendMessageReceived(m dsl.Module, handler func(from t.NodeID, data []*requestpb.Request) error) {
	UponVcbMessageReceived(m, func(from t.NodeID, msg *vcbpb.Message) error {
		startMsgWrapper, ok := msg.Type.(*vcbpb.Message_SendMessage)
		if !ok {
			return nil
		}

		return handler(from, startMsgWrapper.SendMessage.Data)
	})
}

func UponEchoMessageReceived(m dsl.Module, handler func(from t.NodeID, signatureShare []byte) error) {
	UponVcbMessageReceived(m, func(from t.NodeID, msg *vcbpb.Message) error {
		echoMsgWrapper, ok := msg.Type.(*vcbpb.Message_EchoMessage)
		if !ok {
			return nil
		}

		return handler(from, echoMsgWrapper.EchoMessage.SignatureShare)
	})
}

func UponFinalMessageReceived(
	m dsl.Module,
	handler func(from t.NodeID, data []*requestpb.Request, signature []byte) error,
) {
	UponVcbMessageReceived(m, func(from t.NodeID, msg *vcbpb.Message) error {
		finalMsgWrapper, ok := msg.Type.(*vcbpb.Message_FinalMessage)
		if !ok {
			return nil
		}

		finalMsg := finalMsgWrapper.FinalMessage
		return handler(from, finalMsg.Data, finalMsg.Signature)
	})
}
