package modringpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/modringpb/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func FreeSubmodule[C any](m dsl.Module, destModule types.ModuleID, id uint64, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &types1.FreeSubmoduleOrigin{
		Module: m.ModuleID(),
		Type:   &types1.FreeSubmoduleOrigin_Dsl{Dsl: dsl.MirOrigin(contextID)},
	}

	dsl.EmitMirEvent(m, events.FreeSubmodule(destModule, id, origin))
}

func FreedSubmodule(m dsl.Module, destModule types.ModuleID, origin *types1.FreeSubmoduleOrigin) {
	dsl.EmitMirEvent(m, events.FreedSubmodule(destModule, origin))
}
