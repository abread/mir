package bcpbdsl

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/events"
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func StartBroadcast(m dsl.Module, destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txs []*types1.Request) {
	dsl.EmitMirEvent(m, events.StartBroadcast(destModule, queueSlot, txs))
}

func Deliver(m dsl.Module, destModule types.ModuleID, slot *types2.Slot) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, slot))
}

func FreeSlot(m dsl.Module, destModule types.ModuleID, slot *types2.Slot) {
	dsl.EmitMirEvent(m, events.FreeSlot(destModule, slot))
}

func DoFillGap(m dsl.Module, destModule types.ModuleID, slot *types2.Slot) {
	dsl.EmitMirEvent(m, events.DoFillGap(destModule, slot))
}

func BcStarted(m dsl.Module, destModule types.ModuleID, slot *types2.Slot) {
	dsl.EmitMirEvent(m, events.BcStarted(destModule, slot))
}
