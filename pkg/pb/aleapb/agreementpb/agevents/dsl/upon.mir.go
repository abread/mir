// Code generated by Mir codegen. DO NOT EDIT.

package ageventsdsl

import (
	"time"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
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

func UponInputValue(m dsl.Module, handler func(round uint64, input bool) error) {
	UponEvent[*types.Event_InputValue](m, func(ev *types.InputValue) error {
		return handler(ev.Round, ev.Input)
	})
}

func UponDeliver(m dsl.Module, handler func(round uint64, decision bool, posQuorumWait time.Duration, posTotalWait time.Duration) error) {
	UponEvent[*types.Event_Deliver](m, func(ev *types.Deliver) error {
		return handler(ev.Round, ev.Decision, ev.PosQuorumWait, ev.PosTotalWait)
	})
}

func UponStaleMsgsRecvd(m dsl.Module, handler func(messages []*types2.PastMessage) error) {
	UponEvent[*types.Event_StaleMsgsRevcd](m, func(ev *types.StaleMsgsRecvd) error {
		return handler(ev.Messages)
	})
}

func UponInnerAbbaRoundTime(m dsl.Module, handler func(duration time.Duration) error) {
	UponEvent[*types.Event_InnerAbbaRoundTime](m, func(ev *types.InnerAbbaRoundTime) error {
		return handler(ev.Duration)
	})
}
