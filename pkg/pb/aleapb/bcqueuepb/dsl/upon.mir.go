package bcqueuepbdsl

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_AleaBcqueue](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponInputValue(m dsl.Module, handler func(slot *types2.Slot, txs []*types3.Transaction) error) {
	UponEvent[*types.Event_InputValue](m, func(ev *types.InputValue) error {
		return handler(ev.Slot, ev.Txs)
	})
}

func UponDeliver(m dsl.Module, handler func(slot *types2.Slot) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		return handler(ev.Slot)
	})
}

func UponFreeSlot(m dsl.Module, handler func(queueSlot aleatypes.QueueSlot) error) {
	UponEvent[*types.Event_FreeSlot](m, func(ev *types.FreeSlot) error {
		return handler(ev.QueueSlot)
	})
}

func UponPastVcbFinal(m dsl.Module, handler func(queueSlot aleatypes.QueueSlot, txs []*types3.Transaction, signature tctypes.FullSig) error) {
	UponEvent[*types.Event_PastVcbFinal](m, func(ev *types.PastVcbFinal) error {
		return handler(ev.QueueSlot, ev.Txs, ev.Signature)
	})
}

func UponBcStarted(m dsl.Module, handler func(slot *types2.Slot) error) {
	UponEvent[*types.Event_BcStarted](m, func(ev *types.BcStarted) error {
		return handler(ev.Slot)
	})
}
