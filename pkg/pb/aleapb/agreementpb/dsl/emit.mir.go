package agreementpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func RequestInput(m dsl.Module, destModule types.ModuleID, round uint64) {
	dsl.EmitMirEvent(m, events.RequestInput(destModule, round))
}

func InputValue(m dsl.Module, destModule types.ModuleID, round uint64, input bool) {
	dsl.EmitMirEvent(m, events.InputValue(destModule, round, input))
}

func Deliver(m dsl.Module, destModule types.ModuleID, round uint64, decision bool) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, round, decision))
}
