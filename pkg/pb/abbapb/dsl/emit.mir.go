package abbapbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/abbapb/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue[C any](m dsl.Module, destModule types.ModuleID, input bool, context *C) {
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

func RoundInputValue[C any](m dsl.Module, destModule types.ModuleID, input bool, context *C) {
	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.RoundOrigin{
		Module: m.ModuleID(),
		Type:   &types1.RoundOrigin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.RoundInputValue(destModule, input, origin))
}

func RoundDeliver(m dsl.Module, destModule types.ModuleID, nextEstimate bool, origin *types1.RoundOrigin) {
	dsl.EmitMirEvent(m, events.RoundDeliver(destModule, nextEstimate, origin))
}

func RoundFinishAll(m dsl.Module, destModule types.ModuleID, decision bool) {
	dsl.EmitMirEvent(m, events.RoundFinishAll(destModule, decision))
}
