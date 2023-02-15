package modringpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/modringpb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func FreeSubmodule(m dsl.Module, destModule types.ModuleID, id uint64) {
	dsl.EmitMirEvent(m, events.FreeSubmodule(destModule, id))
}
