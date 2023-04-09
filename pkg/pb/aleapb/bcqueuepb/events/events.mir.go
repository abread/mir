package bcqueuepbevents

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	types4 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/requestpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule types.ModuleID, slot *types1.Slot, txs []*types2.Request) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_AleaBcqueue{
			AleaBcqueue: &types4.Event{
				Type: &types4.Event_InputValue{
					InputValue: &types4.InputValue{
						Slot: slot,
						Txs:  txs,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, slot *types1.Slot) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_AleaBcqueue{
			AleaBcqueue: &types4.Event{
				Type: &types4.Event_Deliver{
					Deliver: &types4.Deliver{
						Slot: slot,
					},
				},
			},
		},
	}
}

func FreeSlot(destModule types.ModuleID, queueSlot aleatypes.QueueSlot) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_AleaBcqueue{
			AleaBcqueue: &types4.Event{
				Type: &types4.Event_FreeSlot{
					FreeSlot: &types4.FreeSlot{
						QueueSlot: queueSlot,
					},
				},
			},
		},
	}
}

func PastVcbFinal(destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txs []*types2.Request, signature tctypes.FullSig) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_AleaBcqueue{
			AleaBcqueue: &types4.Event{
				Type: &types4.Event_PastVcbFinal{
					PastVcbFinal: &types4.PastVcbFinal{
						QueueSlot: queueSlot,
						Txs:       txs,
						Signature: signature,
					},
				},
			},
		},
	}
}

func BcStarted(destModule types.ModuleID, slot *types1.Slot) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_AleaBcqueue{
			AleaBcqueue: &types4.Event{
				Type: &types4.Event_BcStarted{
					BcStarted: &types4.BcStarted{
						Slot: slot,
					},
				},
			},
		},
	}
}
