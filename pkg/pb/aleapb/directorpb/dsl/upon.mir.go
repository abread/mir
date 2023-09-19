// Code generated by Mir codegen. DO NOT EDIT.

package directorpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/trantor/types"
	types3 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*types1.Event_AleaDirector](m, func(ev *types.Event) error {
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

func UponNewEpoch(m dsl.Module, handler func(epoch types2.EpochNr) error) {
	UponEvent[*types.Event_NewEpoch](m, func(ev *types.NewEpoch) error {
		return handler(ev.Epoch)
	})
}

func UponGCEpochs(m dsl.Module, handler func(minEpoch types2.EpochNr) error) {
	UponEvent[*types.Event_GcEpochs](m, func(ev *types.GCEpochs) error {
		return handler(ev.MinEpoch)
	})
}

func UponHelpNode(m dsl.Module, handler func(nodeId types3.NodeID) error) {
	UponEvent[*types.Event_HelpNode](m, func(ev *types.HelpNode) error {
		return handler(ev.NodeId)
	})
}
