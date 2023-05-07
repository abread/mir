package bcqueuepbevents

import (
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	types3 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	types4 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txs []*types1.Transaction) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBcqueue{
			AleaBcqueue: &types3.Event{
				Type: &types3.Event_InputValue{
					InputValue: &types3.InputValue{
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

func FreeSlot(destModule types.ModuleID, queueSlot aleatypes.QueueSlot) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBcqueue{
			AleaBcqueue: &types3.Event{
				Type: &types3.Event_FreeSlot{
					FreeSlot: &types3.FreeSlot{
						QueueSlot: queueSlot,
					},
				},
			},
		},
	}
}

func PastVcbFinal(destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txs []*types1.Transaction, signature tctypes.FullSig) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBcqueue{
			AleaBcqueue: &types3.Event{
				Type: &types3.Event_PastVcbFinal{
					PastVcbFinal: &types3.PastVcbFinal{
						QueueSlot: queueSlot,
						Txs:       txs,
						Signature: signature,
					},
				},
			},
		},
	}
}

func BcStarted(destModule types.ModuleID, slot *types4.Slot) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_AleaBcqueue{
			AleaBcqueue: &types3.Event{
				Type: &types3.Event_BcStarted{
					BcStarted: &types3.BcStarted{
						Slot: slot,
					},
				},
			},
		},
	}
}
