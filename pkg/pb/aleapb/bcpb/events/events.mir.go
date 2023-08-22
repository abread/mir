// Code generated by Mir codegen. DO NOT EDIT.

package bcpbevents

import (
	"time"

	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func RequestCert(destModule types.ModuleID) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaBc{
			AleaBc: &types2.Event{
				Type: &types2.Event_RequestCert{
					RequestCert: &types2.RequestCert{},
				},
			},
		},
	}
}

func DeliverCert(destModule types.ModuleID, cert *types2.Cert) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaBc{
			AleaBc: &types2.Event{
				Type: &types2.Event_DeliverCert{
					DeliverCert: &types2.DeliverCert{
						Cert: cert,
					},
				},
			},
		},
	}
}

func BcStarted(destModule types.ModuleID, slot *types2.Slot) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaBc{
			AleaBc: &types2.Event{
				Type: &types2.Event_BcStarted{
					BcStarted: &types2.BcStarted{
						Slot: slot,
					},
				},
			},
		},
	}
}

func FreeSlot(destModule types.ModuleID, slot *types2.Slot) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaBc{
			AleaBc: &types2.Event{
				Type: &types2.Event_FreeSlot{
					FreeSlot: &types2.FreeSlot{
						Slot: slot,
					},
				},
			},
		},
	}
}

func EstimateUpdate(destModule types.ModuleID, maxOwnBcDuration time.Duration, maxOwnBcLocalDuration time.Duration, maxExtBcDuration time.Duration) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaBc{
			AleaBc: &types2.Event{
				Type: &types2.Event_EstimateUpdate{
					EstimateUpdate: &types2.EstimateUpdate{
						MaxOwnBcDuration:      maxOwnBcDuration,
						MaxOwnBcLocalDuration: maxOwnBcLocalDuration,
						MaxExtBcDuration:      maxExtBcDuration,
					},
				},
			},
		},
	}
}

func DoFillGap(destModule types.ModuleID, slot *types2.Slot) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaBc{
			AleaBc: &types2.Event{
				Type: &types2.Event_FillGap{
					FillGap: &types2.DoFillGap{
						Slot: slot,
					},
				},
			},
		},
	}
}
