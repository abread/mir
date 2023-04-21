package vcbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types1 "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	events "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types2 "github.com/filecoin-project/mir/pkg/trantor/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue(m dsl.Module, destModule types.ModuleID, txs []*types1.Request) {
	dsl.EmitMirEvent(m, events.InputValue(destModule, txs))
}

func Deliver(m dsl.Module, destModule types.ModuleID, txs []*types1.Request, txIds []types2.TxID, signature tctypes.FullSig, srcModule types.ModuleID) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, txs, txIds, signature, srcModule))
}
