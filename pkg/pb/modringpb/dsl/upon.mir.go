package modringpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_Modring](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponFreeSubmodule(m dsl.Module, handler func(id uint64, origin *types.FreeSubmoduleOrigin) error) {
	UponEvent[*types.Event_Free](m, func(ev *types.FreeSubmodule) error {
		return handler(ev.Id, ev.Origin)
	})
}

func UponFreedSubmodule[C any](m dsl.Module, handler func(context *C) error) {
	UponEvent[*types.Event_Freed](m, func(ev *types.FreedSubmodule) error {
		originWrapper, ok := ev.Origin.Type.(*types.FreeSubmoduleOrigin_Dsl)
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

func UponPastMessagesRecvd(m dsl.Module, handler func(messages []*types.PastMessage) error) {
	UponEvent[*types.Event_PastMessagesRecvd](m, func(ev *types.PastMessagesRecvd) error {
		return handler(ev.Messages)
	})
}
