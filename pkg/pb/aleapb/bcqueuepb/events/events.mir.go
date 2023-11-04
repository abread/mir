// Code generated by Mir codegen. DO NOT EDIT.

package bcqueuepbevents

import (
	"time"

	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	types5 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	types4 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types1 "github.com/filecoin-project/mir/pkg/trantor/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txIds []types1.TxID, txs []*types2.Transaction) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_AleaBcqueue{
			AleaBcqueue: &types4.Event{
				Type: &types4.Event_InputValue{
					InputValue: &types4.InputValue{
						QueueSlot: queueSlot,
						TxIds:     txIds,
						Txs:       txs,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, cert *types5.Cert) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_AleaBcqueue{
			AleaBcqueue: &types4.Event{
				Type: &types4.Event_Deliver{
					Deliver: &types4.Deliver{
						Cert: cert,
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

func PastVcbFinal(destModule types.ModuleID, queueSlot aleatypes.QueueSlot, txs []*types2.Transaction, signature tctypes.FullSig) *types3.Event {
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

func BcStarted(destModule types.ModuleID, slot *types5.Slot) *types3.Event {
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

func BcQuorumDone(destModule types.ModuleID, slot *types5.Slot, deliverDelta time.Duration) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_AleaBcqueue{
			AleaBcqueue: &types4.Event{
				Type: &types4.Event_BcQuorumDone{
					BcQuorumDone: &types4.BcQuorumDone{
						Slot:         slot,
						DeliverDelta: deliverDelta,
					},
				},
			},
		},
	}
}

func BcAllDone(destModule types.ModuleID, slot *types5.Slot, quorumDoneDelta time.Duration) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_AleaBcqueue{
			AleaBcqueue: &types4.Event{
				Type: &types4.Event_BcAllDone{
					BcAllDone: &types4.BcAllDone{
						Slot:            slot,
						QuorumDoneDelta: quorumDoneDelta,
					},
				},
			},
		},
	}
}

func NetLatencyEstimate(destModule types.ModuleID, minEstimate time.Duration) *types3.Event {
	return &types3.Event{
		DestModule: destModule,
		Type: &types3.Event_AleaBcqueue{
			AleaBcqueue: &types4.Event{
				Type: &types4.Event_NetLatEst{
					NetLatEst: &types4.NetLatencyEstimate{
						MinEstimate: minEstimate,
					},
				},
			},
		},
	}
}
