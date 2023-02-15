package bcqueuepbdsl

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
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

func UponInputValue(m dsl.Module, handler func(slot *types2.Slot, txs []*requestpb.Request) error) {
	UponEvent[*types.Event_InputValue](m, func(ev *types.InputValue) error {
		return handler(ev.Slot, ev.Txs)
	})
}

func UponDeliver(m dsl.Module, handler func(slot *types2.Slot) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		return handler(ev.Slot)
	})
}

func UponFreeSlot(m dsl.Module, handler func(queueSlot aleatypes.QueueSlot, origin *types.FreeSlotOrigin) error) {
	UponEvent[*types.Event_FreeSlot](m, func(ev *types.FreeSlot) error {
		return handler(ev.QueueSlot, ev.Origin)
	})
}

func UponSlotFreed[C any](m dsl.Module, handler func(context *C) error) {
	UponEvent[*types.Event_SlotFreed](m, func(ev *types.SlotFreed) error {
		originWrapper, ok := ev.Origin.Type.(*types.FreeSlotOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		return handler(context)
	})
}
