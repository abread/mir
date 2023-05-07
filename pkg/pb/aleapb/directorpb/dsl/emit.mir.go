package directorpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types1 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	events "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func Heartbeat(m dsl.Module, destModule types.ModuleID) {
	dsl.EmitMirEvent(m, events.Heartbeat(destModule))
}

func DoFillGap(m dsl.Module, destModule types.ModuleID, slot *types1.Slot) {
	dsl.EmitMirEvent(m, events.DoFillGap(destModule, slot))
}
