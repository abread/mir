package vcbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	events "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func BroadcastRequest(m dsl.Module, destModule types.ModuleID, txIds []types.TxID, txs []*requestpb.Request) {
	dsl.EmitMirEvent(m, events.BroadcastRequest(destModule, txIds, txs))
}

func Deliver(m dsl.Module, destModule types.ModuleID, txs []*requestpb.Request, txIds []types.TxID, signature []uint8, originModule types.ModuleID) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, txs, txIds, signature, originModule))
}
