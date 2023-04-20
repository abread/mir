package directorpbevents

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func Heartbeat(destModule types.ModuleID) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaDirector{
			AleaDirector: &types2.Event{
				Type: &types2.Event_Heartbeat{
					Heartbeat: &types2.Heartbeat{},
				},
			},
		},
	}
}