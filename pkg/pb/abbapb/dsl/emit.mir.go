// Code generated by Mir codegen. DO NOT EDIT.

package abbapbdsl

import (
	"time"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/abbapb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue(m dsl.Module, destModule types.ModuleID, input bool) {
	dsl.EmitEvent(m, events.InputValue(destModule, input))
}

func ContinueExecution(m dsl.Module, destModule types.ModuleID) {
	dsl.EmitEvent(m, events.ContinueExecution(destModule))
}

func Deliver(m dsl.Module, destModule types.ModuleID, result bool, srcModule types.ModuleID) {
	dsl.EmitEvent(m, events.Deliver(destModule, result, srcModule))
}

func RoundInputValue(m dsl.Module, destModule types.ModuleID, input bool) {
	dsl.EmitEvent(m, events.RoundInputValue(destModule, input))
}

func RoundContinue(m dsl.Module, destModule types.ModuleID) {
	dsl.EmitEvent(m, events.RoundContinue(destModule))
}

func RoundDeliver(m dsl.Module, destModule types.ModuleID, nextEstimate bool, roundNumber uint64, durationNoCoin time.Duration) {
	dsl.EmitEvent(m, events.RoundDeliver(destModule, nextEstimate, roundNumber, durationNoCoin))
}

func RoundFinishAll(m dsl.Module, destModule types.ModuleID, decision bool, unanimous bool) {
	dsl.EmitEvent(m, events.RoundFinishAll(destModule, decision, unanimous))
}

func Done(m dsl.Module, destModule types.ModuleID, srcModule types.ModuleID) {
	dsl.EmitEvent(m, events.Done(destModule, srcModule))
}
