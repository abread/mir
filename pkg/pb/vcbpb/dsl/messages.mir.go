package vcbpbdsl

import (
	attribute "go.opentelemetry.io/otel/attribute"
	trace "go.opentelemetry.io/otel/trace"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	types "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_Vcb](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponSendMessageReceived(m dsl.Module, handler func(from types1.NodeID, txs []*types3.Request) error) {
	UponMessageReceived[*types.Message_SendMessage](m, func(from types1.NodeID, msg *types.SendMessage) error {
		spanFromAttr := attribute.String("from", string(from))
		spanMsgAttr := attribute.String("message", msg.Pb().String())
		spanAttrs := trace.WithAttributes(spanFromAttr, spanMsgAttr)
		m.DslHandle().PushSpan("UponSendMessageReceived", spanAttrs)
		defer m.DslHandle().PopSpan()

		return handler(from, msg.Txs)
	})
}

func UponEchoMessageReceived(m dsl.Module, handler func(from types1.NodeID, signatureShare tctypes.SigShare) error) {
	UponMessageReceived[*types.Message_EchoMessage](m, func(from types1.NodeID, msg *types.EchoMessage) error {
		spanFromAttr := attribute.String("from", string(from))
		spanMsgAttr := attribute.String("message", msg.Pb().String())
		spanAttrs := trace.WithAttributes(spanFromAttr, spanMsgAttr)
		m.DslHandle().PushSpan("UponEchoMessageReceived", spanAttrs)
		defer m.DslHandle().PopSpan()

		return handler(from, msg.SignatureShare)
	})
}

func UponFinalMessageReceived(m dsl.Module, handler func(from types1.NodeID, signature tctypes.FullSig) error) {
	UponMessageReceived[*types.Message_FinalMessage](m, func(from types1.NodeID, msg *types.FinalMessage) error {
		spanFromAttr := attribute.String("from", string(from))
		spanMsgAttr := attribute.String("message", msg.Pb().String())
		spanAttrs := trace.WithAttributes(spanFromAttr, spanMsgAttr)
		m.DslHandle().PushSpan("UponFinalMessageReceived", spanAttrs)
		defer m.DslHandle().PopSpan()

		return handler(from, msg.Signature)
	})
}
