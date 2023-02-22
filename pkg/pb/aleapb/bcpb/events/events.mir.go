package bcpbevents

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types "github.com/filecoin-project/mir/pkg/types"
)

func StartBroadcast(destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txs []*requestpb.Request) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaBroadcast{
			AleaBroadcast: &types2.Event{
				Type: &types2.Event_StartBroadcast{
					StartBroadcast: &types2.StartBroadcast{
						QueueSlot: queueSlot,
						Txs:       txs,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, slot *types3.Slot) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaBroadcast{
			AleaBroadcast: &types2.Event{
				Type: &types2.Event_Deliver{
					Deliver: &types2.Deliver{
						Slot: slot,
					},
				},
			},
		},
	}
}

func FreeSlot(destModule types.ModuleID, slot *types3.Slot) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaBroadcast{
			AleaBroadcast: &types2.Event{
				Type: &types2.Event_FreeSlot{
					FreeSlot: &types2.FreeSlot{
						Slot: slot,
					},
				},
			},
		},
	}
}

func DoFillGap(destModule types.ModuleID, slot *types3.Slot) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaBroadcast{
			AleaBroadcast: &types2.Event{
				Type: &types2.Event_FillGap{
					FillGap: &types2.DoFillGap{
						Slot: slot,
					},
				},
			},
		},
	}
}
