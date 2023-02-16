package modringpbevents

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func FreeSubmodule(destModule types.ModuleID, id uint64, origin *types1.FreeSubmoduleOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Modring{
			Modring: &types1.Event{
				Type: &types1.Event_Free{
					Free: &types1.FreeSubmodule{
						Id:     id,
						Origin: origin,
					},
				},
			},
		},
	}
}

func FreedSubmodule(destModule types.ModuleID, origin *types1.FreeSubmoduleOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Modring{
			Modring: &types1.Event{
				Type: &types1.Event_Freed{
					Freed: &types1.FreedSubmodule{
						Origin: origin,
					},
				},
			},
		},
	}
}

func PastMessagesRecvd(destModule types.ModuleID, messages []*types1.PastMessage) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Modring{
			Modring: &types1.Event{
				Type: &types1.Event_PastMessagesRecvd{
					PastMessagesRecvd: &types1.PastMessagesRecvd{
						Messages: messages,
					},
				},
			},
		},
	}
}
