package bcbpbdsl

import (
	attribute "go.opentelemetry.io/otel/attribute"
	trace "go.opentelemetry.io/otel/trace"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/bcbpb/types"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_Bcb](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponStartMessageReceived(m dsl.Module, handler func(from types1.NodeID, data []uint8) error) {
	UponMessageReceived[*types.Message_StartMessage](m, func(from types1.NodeID, msg *types.StartMessage) error {
		spanFromAttr := attribute.String("from", string(from))
		spanMsgAttr := attribute.String("message", msg.Pb().String())
		spanAttrs := trace.WithAttributes(spanFromAttr, spanMsgAttr)
		m.DslHandle().PushSpan("UponStartMessageReceived", spanAttrs)
		defer m.DslHandle().PopSpan()

		return handler(from, msg.Data)
	})
}

func UponEchoMessageReceived(m dsl.Module, handler func(from types1.NodeID, signature []uint8) error) {
	UponMessageReceived[*types.Message_EchoMessage](m, func(from types1.NodeID, msg *types.EchoMessage) error {
		spanFromAttr := attribute.String("from", string(from))
		spanMsgAttr := attribute.String("message", msg.Pb().String())
		spanAttrs := trace.WithAttributes(spanFromAttr, spanMsgAttr)
		m.DslHandle().PushSpan("UponEchoMessageReceived", spanAttrs)
		defer m.DslHandle().PopSpan()

		return handler(from, msg.Signature)
	})
}

func UponFinalMessageReceived(m dsl.Module, handler func(from types1.NodeID, data []uint8, signers []types1.NodeID, signatures [][]uint8) error) {
	UponMessageReceived[*types.Message_FinalMessage](m, func(from types1.NodeID, msg *types.FinalMessage) error {
		spanFromAttr := attribute.String("from", string(from))
		spanMsgAttr := attribute.String("message", msg.Pb().String())
		spanAttrs := trace.WithAttributes(spanFromAttr, spanMsgAttr)
		m.DslHandle().PushSpan("UponFinalMessageReceived", spanAttrs)
		defer m.DslHandle().PopSpan()

		return handler(from, msg.Data, msg.Signers, msg.Signatures)
	})
}
