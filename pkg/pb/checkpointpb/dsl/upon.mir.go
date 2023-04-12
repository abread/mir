package checkpointpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/checkpointpb/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/commonpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_Checkpoint](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponStableCheckpoint(m dsl.Module, handler func(sn types2.SeqNr, snapshot *types3.StateSnapshot, cert map[types2.NodeID][]uint8) error) {
	UponEvent[*types.Event_StableCheckpoint](m, func(ev *types.StableCheckpoint) error {
		return handler(ev.Sn, ev.Snapshot, ev.Cert)
	})
}

func UponEpochProgress(m dsl.Module, handler func(nodeId types2.NodeID, epoch types2.EpochNr) error) {
	UponEvent[*types.Event_EpochProgress](m, func(ev *types.EpochProgress) error {
		return handler(ev.NodeId, ev.Epoch)
	})
}
