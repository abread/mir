package vcbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types "github.com/filecoin-project/mir/pkg/pb/vcbpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types2 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_Vcb](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponInputValue(m dsl.Module, handler func(txs []*requestpb.Request) error) {
	UponEvent[*types.Event_InputValue](m, func(ev *types.InputValue) error {
		return handler(ev.Txs)
	})
}

func UponDeliver(m dsl.Module, handler func(txs []*requestpb.Request, txIds []types2.TxID, signature tctypes.FullSig, srcModule types2.ModuleID) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		return handler(ev.Txs, ev.TxIds, ev.Signature, ev.SrcModule)
	})
}
