package abbapbdsl

import (
	trace "go.opentelemetry.io/otel/trace"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/abbapb/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue[C any](m dsl.Module, destModule types.ModuleID, input bool, context *C) {
	kind := trace.WithSpanKind(trace.SpanKindProducer)
	m.DslHandle().PushSpan("InputValue", kind)
	defer m.DslHandle().PopSpan()

	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.Origin{
		Module: m.ModuleID(),
		Type:   &types1.Origin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.InputValue(destModule, input, origin))
}

func Deliver(m dsl.Module, destModule types.ModuleID, result bool, origin *types1.Origin) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, result, origin))
}

func RoundInputValue(m dsl.Module, destModule types.ModuleID, input bool) {
	dsl.EmitMirEvent(m, events.RoundInputValue(destModule, input))
}

func RoundDeliver(m dsl.Module, destModule types.ModuleID, nextEstimate bool, roundNumber uint64) {
	dsl.EmitMirEvent(m, events.RoundDeliver(destModule, nextEstimate, roundNumber))
}

func RoundFinishAll(m dsl.Module, destModule types.ModuleID, decision bool) {
	dsl.EmitMirEvent(m, events.RoundFinishAll(destModule, decision))
}
