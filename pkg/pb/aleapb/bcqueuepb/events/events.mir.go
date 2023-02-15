package bcqueuepbevents

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	types3 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule types.ModuleID, slot *types1.Slot, txs []*requestpb.Request) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBcqueue{
			AleaBcqueue: &types3.Event{
				Type: &types3.Event_InputValue{
					InputValue: &types3.InputValue{
						Slot: slot,
						Txs:  txs,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, slot *types1.Slot) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBcqueue{
			AleaBcqueue: &types3.Event{
				Type: &types3.Event_Deliver{
					Deliver: &types3.Deliver{
						Slot: slot,
					},
				},
			},
		},
	}
}

func FreeSlot(destModule types.ModuleID, queueSlot aleatypes.QueueSlot, origin *types3.FreeSlotOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBcqueue{
			AleaBcqueue: &types3.Event{
				Type: &types3.Event_FreeSlot{
					FreeSlot: &types3.FreeSlot{
						QueueSlot: queueSlot,
						Origin:    origin,
					},
				},
			},
		},
	}
}

func SlotFreed(destModule types.ModuleID, origin *types3.FreeSlotOrigin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBcqueue{
			AleaBcqueue: &types3.Event{
				Type: &types3.Event_SlotFreed{
					SlotFreed: &types3.SlotFreed{
						Origin: origin,
					},
				},
			},
		},
	}
}
