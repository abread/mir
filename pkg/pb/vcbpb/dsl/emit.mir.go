package vcbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	events "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue(m dsl.Module, destModule types.ModuleID, txs []*requestpb.Request) {
	dsl.EmitMirEvent(m, events.InputValue(destModule, txs))
}

func Deliver(m dsl.Module, destModule types.ModuleID, txs []*requestpb.Request, txIds []types.TxID, signature tctypes.FullSig, srcModule types.ModuleID) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, txs, txIds, signature, srcModule))
}
