package aleapbdsl

import (
	attribute "go.opentelemetry.io/otel/attribute"
	trace "go.opentelemetry.io/otel/trace"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types3 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/types"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types4 "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_Alea](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponFillGapMessageReceived(m dsl.Module, handler func(from types1.NodeID, slot *types3.Slot) error) {
	UponMessageReceived[*types.Message_FillGapMessage](m, func(from types1.NodeID, msg *types.FillGapMessage) error {
		spanFromAttr := attribute.String("from", string(from))
		spanMsgAttr := attribute.String("message", msg.Pb().String())
		spanAttrs := trace.WithAttributes(spanFromAttr, spanMsgAttr)
		m.DslHandle().PushSpan("UponFillGapMessageReceived", spanAttrs)
		defer m.DslHandle().PopSpan()

		return handler(from, msg.Slot)
	})
}

func UponFillerMessageReceived(m dsl.Module, handler func(from types1.NodeID, slot *types3.Slot, txs []*types4.Request, signature tctypes.FullSig) error) {
	UponMessageReceived[*types.Message_FillerMessage](m, func(from types1.NodeID, msg *types.FillerMessage) error {
		spanFromAttr := attribute.String("from", string(from))
		spanMsgAttr := attribute.String("message", msg.Pb().String())
		spanAttrs := trace.WithAttributes(spanFromAttr, spanMsgAttr)
		m.DslHandle().PushSpan("UponFillerMessageReceived", spanAttrs)
		defer m.DslHandle().PopSpan()

		return handler(from, msg.Slot, msg.Txs, msg.Signature)
	})
}