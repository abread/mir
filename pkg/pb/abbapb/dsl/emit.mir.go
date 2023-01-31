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

func Deliver(m dsl.Module, destModule types.ModuleID, result bool, originModule types.ModuleID) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, result, originModule))
}
