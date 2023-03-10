package abbapbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/abbapb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue(m dsl.Module, destModule types.ModuleID, input bool) {
	dsl.EmitMirEvent(m, events.InputValue(destModule, input))
}

func Deliver(m dsl.Module, destModule types.ModuleID, result bool, srcModule types.ModuleID) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, result, srcModule))
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
