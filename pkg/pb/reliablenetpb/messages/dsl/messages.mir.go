package messagesdsl

import (
	attribute "go.opentelemetry.io/otel/attribute"
	trace "go.opentelemetry.io/otel/trace"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages/types"
	rntypes "github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_ReliableNet](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponAckMessageReceived(m dsl.Module, handler func(from types1.NodeID, msgDestModule types1.ModuleID, msgId rntypes.MsgID) error) {
	UponMessageReceived[*types.Message_Ack](m, func(from types1.NodeID, msg *types.AckMessage) error {
		spanFromAttr := attribute.String("from", string(from))
		spanMsgAttr := attribute.String("message", msg.Pb().String())
		spanAttrs := trace.WithAttributes(spanFromAttr, spanMsgAttr)
		m.DslHandle().PushSpan("UponAckMessageReceived", spanAttrs)
		defer m.DslHandle().PopSpan()

		return handler(from, msg.MsgDestModule, msg.MsgId)
	})
}