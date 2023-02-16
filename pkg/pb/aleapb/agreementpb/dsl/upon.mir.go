package agreementpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_AleaAgreement](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponRequestInput(m dsl.Module, handler func(round uint64) error) {
	UponEvent[*types.Event_RequestInput](m, func(ev *types.RequestInput) error {
		return handler(ev.Round)
	})
}

func UponInputValue(m dsl.Module, handler func(round uint64, input bool) error) {
	UponEvent[*types.Event_InputValue](m, func(ev *types.InputValue) error {
		return handler(ev.Round, ev.Input)
	})
}

func UponDeliver(m dsl.Module, handler func(round uint64, decision bool) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		return handler(ev.Round, ev.Decision)
	})
}