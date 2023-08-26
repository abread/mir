// Code generated by Mir codegen. DO NOT EDIT.

package vcbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types1 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	events "github.com/filecoin-project/mir/pkg/pb/vcbpb/events"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue(m dsl.Module, destModule types.ModuleID, txs []*types1.Transaction) {
	dsl.EmitEvent(m, events.InputValue(destModule, txs))
}

func Deliver(m dsl.Module, destModule types.ModuleID, batchId string, signature tctypes.FullSig, srcModule types.ModuleID) {
	dsl.EmitEvent(m, events.Deliver(destModule, batchId, signature, srcModule))
}

func QuorumDone(m dsl.Module, destModule types.ModuleID, srcModule types.ModuleID) {
	dsl.EmitEvent(m, events.QuorumDone(destModule, srcModule))
}

func AllDone(m dsl.Module, destModule types.ModuleID, srcModule types.ModuleID) {
	dsl.EmitEvent(m, events.AllDone(destModule, srcModule))
}
