package abbapbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_Abba](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponInputValue(m dsl.Module, handler func(input bool, origin *types.Origin) error) {
	UponEvent[*types.Event_InputValue](m, func(ev *types.InputValue) error {
		originWrapper, ok := ev.Origin.Type.(*types.Origin_Dsl)
		if ok {
			m.DslHandle().ImportTraceContextFromMap(originWrapper.Dsl.TraceContext)
		}

		m.DslHandle().PushSpan("InputValue")
		defer m.DslHandle().PopSpan()

		return handler(ev.Input, ev.Origin)
	})
}

func UponDeliver[C any](m dsl.Module, handler func(result bool, context *C) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		originWrapper, ok := ev.Origin.Type.(*types.Origin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		m.DslHandle().ImportTraceContextFromMap(originWrapper.Dsl.TraceContext)

		return handler(ev.Result, context)
	})
}

func UponRoundEvent[W types.RoundEvent_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	UponEvent[*types.Event_Round](m, func(ev *types.RoundEvent) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponRoundInputValue(m dsl.Module, handler func(input bool, origin *types.RoundOrigin) error) {
	UponRoundEvent[*types.RoundEvent_InputValue](m, func(ev *types.RoundInputValue) error {
		originWrapper, ok := ev.Origin.Type.(*types.RoundOrigin_Dsl)
		if ok {
			m.DslHandle().ImportTraceContextFromMap(originWrapper.Dsl.TraceContext)
		}

		m.DslHandle().PushSpan("RoundInputValue")
		defer m.DslHandle().PopSpan()

		return handler(ev.Input, ev.Origin)
	})
}

func UponRoundDeliver[C any](m dsl.Module, handler func(nextEstimate bool, context *C) error) {
	UponRoundEvent[*types.RoundEvent_Deliver](m, func(ev *types.RoundDeliver) error {
		originWrapper, ok := ev.Origin.Type.(*types.RoundOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		m.DslHandle().ImportTraceContextFromMap(originWrapper.Dsl.TraceContext)

		return handler(ev.NextEstimate, context)
	})
}

func UponRoundFinishAll(m dsl.Module, handler func(decision bool) error) {
	UponRoundEvent[*types.RoundEvent_Finish](m, func(ev *types.RoundFinishAll) error {
		return handler(ev.Decision)
	})
}
