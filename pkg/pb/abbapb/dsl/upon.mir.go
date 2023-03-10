package abbapbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/types"
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

func UponInputValue(m dsl.Module, handler func(input bool) error) {
	UponEvent[*types.Event_InputValue](m, func(ev *types.InputValue) error {
		return handler(ev.Input)
	})
}

func UponDeliver(m dsl.Module, handler func(result bool, srcModule types2.ModuleID) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		return handler(ev.Result, ev.SrcModule)
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
