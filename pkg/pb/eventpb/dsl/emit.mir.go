package eventpbdsl

import (
	trace "go.opentelemetry.io/otel/trace"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/eventpb/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func SignRequest[C any](m dsl.Module, destModule types.ModuleID, data [][]uint8, context *C) {
	kind := trace.WithSpanKind(trace.SpanKindProducer)
	m.DslHandle().PushSpan("SignRequest", kind)
	defer m.DslHandle().PopSpan()

	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.SignOrigin{
		Module: m.ModuleID(),
		Type:   &types1.SignOrigin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.SignRequest(destModule, data, origin))
}

func SignResult(m dsl.Module, destModule types.ModuleID, signature []uint8, origin *types1.SignOrigin) {
	dsl.EmitMirEvent(m, events.SignResult(destModule, signature, origin))
}

func VerifyNodeSigs[C any](m dsl.Module, destModule types.ModuleID, data []*types1.SigVerData, signatures [][]uint8, nodeIds []types.NodeID, context *C) {
	kind := trace.WithSpanKind(trace.SpanKindProducer)
	m.DslHandle().PushSpan("VerifyNodeSigs", kind)
	defer m.DslHandle().PopSpan()

	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.SigVerOrigin{
		Module: m.ModuleID(),
		Type:   &types1.SigVerOrigin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.VerifyNodeSigs(destModule, data, signatures, origin, nodeIds))
}

func NodeSigsVerified(m dsl.Module, destModule types.ModuleID, origin *types1.SigVerOrigin, nodeIds []types.NodeID, valid []bool, errors []error, allOk bool) {
	dsl.EmitMirEvent(m, events.NodeSigsVerified(destModule, origin, nodeIds, valid, errors, allOk))
}

func SendMessage(m dsl.Module, destModule types.ModuleID, msg *types2.Message, destinations []types.NodeID) {
	dsl.EmitMirEvent(m, events.SendMessage(destModule, msg, destinations))
}

func MessageReceived(m dsl.Module, destModule types.ModuleID, from types.NodeID, msg *types2.Message) {
	dsl.EmitMirEvent(m, events.MessageReceived(destModule, from, msg))
}
