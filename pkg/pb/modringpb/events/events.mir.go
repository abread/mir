package modringpbevents

import (
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func FreeSubmodule(destModule types.ModuleID, id uint64) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Modring{
			Modring: &types2.Event{
				Type: &types2.Event_Free{
					Free: &types2.FreeSubmodule{
						Id: id,
					},
				},
			},
		},
	}
}
