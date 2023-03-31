package bcpbdsl

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_AleaBroadcast](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponStartBroadcast(m dsl.Module, handler func(queueSlot aleatypes.QueueSlot, txs []*requestpb.Request) error) {
	UponEvent[*types.Event_StartBroadcast](m, func(ev *types.StartBroadcast) error {
		return handler(ev.QueueSlot, ev.Txs)
	})
}

func UponDeliver(m dsl.Module, handler func(slot *types2.Slot) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		return handler(ev.Slot)
	})
}

func UponFreeSlot(m dsl.Module, handler func(slot *types2.Slot) error) {
	UponEvent[*types.Event_FreeSlot](m, func(ev *types.FreeSlot) error {
		return handler(ev.Slot)
	})
}

func UponDoFillGap(m dsl.Module, handler func(slot *types2.Slot) error) {
	UponEvent[*types.Event_FillGap](m, func(ev *types.DoFillGap) error {
		return handler(ev.Slot)
	})
}

func UponBcStarted(m dsl.Module, handler func(slot *types2.Slot) error) {
	UponEvent[*types.Event_BcStarted](m, func(ev *types.BcStarted) error {
		return handler(ev.Slot)
	})
}
