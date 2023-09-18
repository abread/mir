// Code generated by Mir codegen. DO NOT EDIT.

package bcpbdsl

import (
	"time"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*types1.Event_AleaBc](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponRequestCert(m dsl.Module, handler func() error) {
	UponEvent[*types.Event_RequestCert](m, func(ev *types.RequestCert) error {
		return handler()
	})
}

func UponDeliverCert(m dsl.Module, handler func(cert *types.Cert) error) {
	UponEvent[*types.Event_DeliverCert](m, func(ev *types.DeliverCert) error {
		return handler(ev.Cert)
	})
}

func UponBcStarted(m dsl.Module, handler func(slot *types.Slot) error) {
	UponEvent[*types.Event_BcStarted](m, func(ev *types.BcStarted) error {
		return handler(ev.Slot)
	})
}

func UponFreeSlot(m dsl.Module, handler func(slot *types.Slot) error) {
	UponEvent[*types.Event_FreeSlot](m, func(ev *types.FreeSlot) error {
		return handler(ev.Slot)
	})
}

func UponEstimateUpdate(m dsl.Module, handler func(maxOwnBcDuration time.Duration, maxOwnBcLocalDuration time.Duration, maxExtBcDuration time.Duration) error) {
	UponEvent[*types.Event_EstimateUpdate](m, func(ev *types.EstimateUpdate) error {
		return handler(ev.MaxOwnBcDuration, ev.MaxOwnBcLocalDuration, ev.MaxExtBcDuration)
	})
}

func UponDoFillGap(m dsl.Module, handler func(slot *types.Slot) error) {
	UponEvent[*types.Event_FillGap](m, func(ev *types.DoFillGap) error {
		return handler(ev.Slot)
	})
}

func UponMarkStableProposal(m dsl.Module, handler func(slot *types.Slot) error) {
	UponEvent[*types.Event_MarkStable](m, func(ev *types.MarkStableProposal) error {
		return handler(ev.Slot)
	})
}
