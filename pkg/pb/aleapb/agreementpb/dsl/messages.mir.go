package agreementpbdsl

import (
	attribute "go.opentelemetry.io/otel/attribute"
	trace "go.opentelemetry.io/otel/trace"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/types"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_AleaAgreement](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponFinishAbbaMessageReceived(m dsl.Module, handler func(from types1.NodeID, round uint64, value bool) error) {
	UponMessageReceived[*types.Message_FinishAbba](m, func(from types1.NodeID, msg *types.FinishAbbaMessage) error {
		spanFromAttr := attribute.String("from", string(from))
		spanMsgAttr := attribute.String("message", msg.Pb().String())
		spanAttrs := trace.WithAttributes(spanFromAttr, spanMsgAttr)
		m.DslHandle().PushSpan("UponFinishAbbaMessageReceived", spanAttrs)
		defer m.DslHandle().PopSpan()

		return handler(from, msg.Round, msg.Value)
	})
}
