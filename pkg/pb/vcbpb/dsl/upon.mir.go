// Code generated by Mir codegen. DO NOT EDIT.

package vcbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	types "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types3 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*types1.Event_Vcb](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponInputValue(m dsl.Module, handler func(txIds []string, txs []*types2.Transaction) error) {
	UponEvent[*types.Event_InputValue](m, func(ev *types.InputValue) error {
		return handler(ev.TxIds, ev.Txs)
	})
}

func UponDeliver(m dsl.Module, handler func(batchId string, signature tctypes.FullSig, srcModule types3.ModuleID) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		return handler(ev.BatchId, ev.Signature, ev.SrcModule)
	})
}

func UponQuorumDone(m dsl.Module, handler func(srcModule types3.ModuleID) error) {
	UponEvent[*types.Event_QuorumDone](m, func(ev *types.QuorumDone) error {
		return handler(ev.SrcModule)
	})
}

func UponAllDone(m dsl.Module, handler func(srcModule types3.ModuleID) error) {
	UponEvent[*types.Event_AllDone](m, func(ev *types.AllDone) error {
		return handler(ev.SrcModule)
	})
}
