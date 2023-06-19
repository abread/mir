// Code generated by Mir codegen. DO NOT EDIT.

package directorpbdsl

import (
	"time"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_AleaDirector](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponHeartbeat(m dsl.Module, handler func() error) {
	UponEvent[*types.Event_Heartbeat](m, func(ev *types.Heartbeat) error {
		return handler()
	})
}

func UponDoFillGap(m dsl.Module, handler func(slot *types2.Slot) error) {
	UponEvent[*types.Event_FillGap](m, func(ev *types.DoFillGap) error {
		return handler(ev.Slot)
	})
}

func UponStats(m dsl.Module, handler func(slotsWaitingDelivery uint64, minAgDurationEst time.Duration, maxAgDurationEst time.Duration, minBcDurationEst time.Duration, maxBcDurationEst time.Duration, minOwnBcDurationEst time.Duration, maxOwnBcDurationEst time.Duration, bcEstMargin time.Duration) error) {
	UponEvent[*types.Event_Stats](m, func(ev *types.Stats) error {
		return handler(ev.SlotsWaitingDelivery, ev.MinAgDurationEst, ev.MaxAgDurationEst, ev.MinBcDurationEst, ev.MaxBcDurationEst, ev.MinOwnBcDurationEst, ev.MaxOwnBcDurationEst, ev.BcEstMargin)
	})
}
