package bcqueuepbdsl

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/events"
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue(m dsl.Module, destModule types.ModuleID, slot *types1.Slot, txs []*requestpb.Request) {
	dsl.EmitMirEvent(m, events.InputValue(destModule, slot, txs))
}

func Deliver(m dsl.Module, destModule types.ModuleID, slot *types1.Slot) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, slot))
}

func FreeSlot[C any](m dsl.Module, destModule types.ModuleID, queueSlot aleatypes.QueueSlot, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &types2.FreeSlotOrigin{
		Module: m.ModuleID(),
		Type:   &types2.FreeSlotOrigin_Dsl{Dsl: dsl.MirOrigin(contextID)},
	}

	dsl.EmitMirEvent(m, events.FreeSlot(destModule, queueSlot, origin))
}

func SlotFreed(m dsl.Module, destModule types.ModuleID, origin *types2.FreeSlotOrigin) {
	dsl.EmitMirEvent(m, events.SlotFreed(destModule, origin))
}
