package aagdsl

import (
	aagEvents "github.com/filecoin-project/mir/pkg/alea/agreement/aagevents"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func RequestInput(m dsl.Module, dest t.ModuleID, round uint64) {
	dsl.EmitEvent(m, aagEvents.RequestInput(dest, round))
}

func InputValue(m dsl.Module, dest t.ModuleID, round uint64, input bool) {
	dsl.EmitEvent(m, aagEvents.InputValue(dest, round, input))
}

func Deliver(m dsl.Module, dest t.ModuleID, round uint64, decision bool) {
	dsl.EmitEvent(m, aagEvents.Deliver(dest, round, decision))
}

func UponRequestInput(m dsl.Module, handler func(round uint64) error) {
	UponEvent[*agreementpb.Event_RequestInput](m, func(ev *agreementpb.RequestInput) error {
		return handler(ev.Round)
	})
}

func UponInputValue(m dsl.Module, handler func(round uint64, input bool) error) {
	UponEvent[*agreementpb.Event_InputValue](m, func(ev *agreementpb.InputValue) error {
		return handler(ev.Round, ev.Input)
	})
}

func UponDeliver(m dsl.Module, handler func(round uint64, decision bool) error) {
	UponEvent[*agreementpb.Event_Deliver](m, func(ev *agreementpb.Deliver) error {
		return handler(ev.Round, ev.Decision)
	})
}

func UponEvent[EvWrapper agreementpb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*eventpb.Event_AleaAgreement](m, func(ev *agreementpb.Event) error {
		evWrapper, ok := ev.Type.(EvWrapper)
		if !ok {
			return nil
		}
		return handler(evWrapper.Unwrap())
	})
}
