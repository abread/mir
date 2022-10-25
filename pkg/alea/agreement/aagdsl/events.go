package aagdsl

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func RequestInput(m dsl.Module, destModule t.ModuleID, round uint64) {
	EmitEvent(m, destModule, &agreementpb.Event{
		Type: &agreementpb.Event_RequestInput{
			RequestInput: &agreementpb.RequestInput{
				Round: round,
			},
		},
	})
}

func InputValue(m dsl.Module, destModule t.ModuleID, round uint64, input bool) {
	EmitEvent(m, destModule, &agreementpb.Event{
		Type: &agreementpb.Event_InputValue{
			InputValue: &agreementpb.InputValue{
				Round: round,
				Input: input,
			},
		},
	})
}

func Deliver(m dsl.Module, destModule t.ModuleID, round uint64, decision bool) {
	EmitEvent(m, destModule, &agreementpb.Event{
		Type: &agreementpb.Event_Deliver{
			Deliver: &agreementpb.Deliver{
				Round:    round,
				Decision: decision,
			},
		},
	})
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

func EmitEvent(m dsl.Module, destModule t.ModuleID, ev *agreementpb.Event) {
	evWrapped := &eventpb.Event{
		DestModule: destModule.Pb(),
		Type: &eventpb.Event_AleaAgreement{
			AleaAgreement: ev,
		},
	}

	dsl.EmitEvent(m, evWrapped)
}
