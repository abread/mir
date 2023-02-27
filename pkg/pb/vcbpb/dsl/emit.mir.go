package vcbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	events "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue[C any](m dsl.Module, destModule types.ModuleID, txs []*requestpb.Request, context *C) {
	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.Origin{
		Module: m.ModuleID(),
		Type:   &types1.Origin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.InputValue(destModule, txs, origin))
}

func Deliver(m dsl.Module, destModule types.ModuleID, txs []*requestpb.Request, txIds []types.TxID, signature tctypes.FullSig, origin *types1.Origin) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, txs, txIds, signature, origin))
}
