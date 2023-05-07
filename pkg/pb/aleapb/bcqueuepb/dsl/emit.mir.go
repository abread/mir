package bcqueuepbdsl

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/events"
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue(m dsl.Module, destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txs []*types1.Transaction) {
	dsl.EmitMirEvent(m, events.InputValue(destModule, queueSlot, txs))
}

func Deliver(m dsl.Module, destModule types.ModuleID, slot *types2.Slot) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, slot))
}

func FreeSlot(m dsl.Module, destModule types.ModuleID, queueSlot aleatypes.QueueSlot) {
	dsl.EmitMirEvent(m, events.FreeSlot(destModule, queueSlot))
}

func PastVcbFinal(m dsl.Module, destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txs []*types1.Transaction, signature tctypes.FullSig) {
	dsl.EmitMirEvent(m, events.PastVcbFinal(destModule, queueSlot, txs, signature))
}

func BcStarted(m dsl.Module, destModule types.ModuleID, slot *types2.Slot) {
	dsl.EmitMirEvent(m, events.BcStarted(destModule, slot))
}
