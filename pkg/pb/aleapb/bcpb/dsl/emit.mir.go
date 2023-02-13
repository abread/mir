package bcpbdsl

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func StartBroadcast(m dsl.Module, destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txIds []types.TxID, txs []*requestpb.Request) {
	dsl.EmitMirEvent(m, events.StartBroadcast(destModule, queueSlot, txIds, txs))
}

func Deliver(m dsl.Module, destModule types.ModuleID, slot *types1.Slot) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, slot))
}

func FreeSlot(m dsl.Module, destModule types.ModuleID, slot *types1.Slot) {
	dsl.EmitMirEvent(m, events.FreeSlot(destModule, slot))
}
