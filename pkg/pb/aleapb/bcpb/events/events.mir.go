package bcpbevents

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	types3 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	types4 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func StartBroadcast(destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txs []*types1.Transaction) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBroadcast{
			AleaBroadcast: &types3.Event{
				Type: &types3.Event_StartBroadcast{
					StartBroadcast: &types3.StartBroadcast{
						QueueSlot: queueSlot,
						Txs:       txs,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, slot *types4.Slot) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBroadcast{
			AleaBroadcast: &types3.Event{
				Type: &types3.Event_Deliver{
					Deliver: &types3.Deliver{
						Slot: slot,
					},
				},
			},
		},
	}
}

func FreeSlot(destModule types.ModuleID, slot *types4.Slot) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBroadcast{
			AleaBroadcast: &types3.Event{
				Type: &types3.Event_FreeSlot{
					FreeSlot: &types3.FreeSlot{
						Slot: slot,
					},
				},
			},
		},
	}
}

func DoFillGap(destModule types.ModuleID, slot *types4.Slot) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBroadcast{
			AleaBroadcast: &types3.Event{
				Type: &types3.Event_FillGap{
					FillGap: &types3.DoFillGap{
						Slot: slot,
					},
				},
			},
		},
	}
}

func BcStarted(destModule types.ModuleID, slot *types4.Slot) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBroadcast{
			AleaBroadcast: &types3.Event{
				Type: &types3.Event_BcStarted{
					BcStarted: &types3.BcStarted{
						Slot: slot,
					},
				},
			},
		},
	}
}
