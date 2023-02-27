package abbapbdsl

import (
	trace "go.opentelemetry.io/otel/trace"

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

		kind := trace.WithSpanKind(trace.SpanKindConsumer)
		m.DslHandle().PushSpan("InputValue", kind)
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

		kind := trace.WithSpanKind(trace.SpanKindConsumer)
		m.DslHandle().PushSpan("Deliver", kind)
		defer m.DslHandle().PopSpan()

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

func UponRoundInputValue(m dsl.Module, handler func(input bool) error) {
	UponRoundEvent[*types.RoundEvent_InputValue](m, func(ev *types.RoundInputValue) error {
		return handler(ev.Input)
	})
}

func UponRoundDeliver(m dsl.Module, handler func(nextEstimate bool, roundNumber uint64) error) {
	UponRoundEvent[*types.RoundEvent_Deliver](m, func(ev *types.RoundDeliver) error {
		return handler(ev.NextEstimate, ev.RoundNumber)
	})
}

func UponRoundFinishAll(m dsl.Module, handler func(decision bool) error) {
	UponRoundEvent[*types.RoundEvent_Finish](m, func(ev *types.RoundFinishAll) error {
		return handler(ev.Decision)
	})
}
