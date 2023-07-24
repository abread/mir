// Code generated by Mir codegen. DO NOT EDIT.

package pprepvalidatorpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/pb/ordererpb/pprepvalidatorpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/pbftpb/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_PprepValiadtor](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponValidatePreprepare(m dsl.Module, handler func(preprepare *types2.Preprepare, origin *types.ValidatePreprepareOrigin) error) {
	UponEvent[*types.Event_ValidatePreprepare](m, func(ev *types.ValidatePreprepare) error {
		return handler(ev.Preprepare, ev.Origin)
	})
}

func UponPreprepareValidated[C any](m dsl.Module, handler func(error error, context *C) error) {
	UponEvent[*types.Event_PreprepareValidated](m, func(ev *types.PreprepareValidated) error {
		originWrapper, ok := ev.Origin.Type.(*types.ValidatePreprepareOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		return handler(ev.Error, context)
	})
}
