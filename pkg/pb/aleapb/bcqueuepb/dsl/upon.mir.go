// Code generated by Mir codegen. DO NOT EDIT.

package bcqueuepbdsl

import (
	"time"

	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types3 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*types1.Event_AleaBcqueue](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponInputValue(m dsl.Module, handler func(queueSlot aleatypes.QueueSlot, txs []*types2.Transaction) error) {
	UponEvent[*types.Event_InputValue](m, func(ev *types.InputValue) error {
		return handler(ev.QueueSlot, ev.Txs)
	})
}

func UponDeliver(m dsl.Module, handler func(cert *types3.Cert) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		return handler(ev.Cert)
	})
}

func UponFreeSlot(m dsl.Module, handler func(queueSlot aleatypes.QueueSlot) error) {
	UponEvent[*types.Event_FreeSlot](m, func(ev *types.FreeSlot) error {
		return handler(ev.QueueSlot)
	})
}

func UponPastVcbFinal(m dsl.Module, handler func(queueSlot aleatypes.QueueSlot, txs []*types2.Transaction, signature tctypes.FullSig) error) {
	UponEvent[*types.Event_PastVcbFinal](m, func(ev *types.PastVcbFinal) error {
		return handler(ev.QueueSlot, ev.Txs, ev.Signature)
	})
}

func UponBcStarted(m dsl.Module, handler func(slot *types3.Slot) error) {
	UponEvent[*types.Event_BcStarted](m, func(ev *types.BcStarted) error {
		return handler(ev.Slot)
	})
}

func UponBcQuorumDone(m dsl.Module, handler func(slot *types3.Slot, deliverDelta time.Duration) error) {
	UponEvent[*types.Event_BcQuorumDone](m, func(ev *types.BcQuorumDone) error {
		return handler(ev.Slot, ev.DeliverDelta)
	})
}

func UponBcAllDone(m dsl.Module, handler func(slot *types3.Slot, quorumDoneDelta time.Duration) error) {
	UponEvent[*types.Event_BcAllDone](m, func(ev *types.BcAllDone) error {
		return handler(ev.Slot, ev.QuorumDoneDelta)
	})
}
