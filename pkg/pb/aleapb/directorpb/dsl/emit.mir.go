package directorpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func Heartbeat(m dsl.Module, destModule types.ModuleID) {
	dsl.EmitMirEvent(m, events.Heartbeat(destModule))
}
