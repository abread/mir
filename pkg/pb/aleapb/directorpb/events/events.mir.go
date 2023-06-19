// Code generated by Mir codegen. DO NOT EDIT.

package directorpbevents

import (
	"time"

	types3 "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
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

func DoFillGap(destModule types.ModuleID, slot *types3.Slot) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaDirector{
			AleaDirector: &types2.Event{
				Type: &types2.Event_FillGap{
					FillGap: &types2.DoFillGap{
						Slot: slot,
					},
				},
			},
		},
	}
}

func Stats(destModule types.ModuleID, slotsWaitingDelivery uint64, minAgDurationEst time.Duration, maxAgDurationEst time.Duration, minBcDurationEst time.Duration, maxBcDurationEst time.Duration, minOwnBcDurationEst time.Duration, maxOwnBcDurationEst time.Duration, bcEstMargin time.Duration) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaDirector{
			AleaDirector: &types2.Event{
				Type: &types2.Event_Stats{
					Stats: &types2.Stats{
						SlotsWaitingDelivery: slotsWaitingDelivery,
						MinAgDurationEst:     minAgDurationEst,
						MaxAgDurationEst:     maxAgDurationEst,
						MinBcDurationEst:     minBcDurationEst,
						MaxBcDurationEst:     maxBcDurationEst,
						MinOwnBcDurationEst:  minOwnBcDurationEst,
						MaxOwnBcDurationEst:  maxOwnBcDurationEst,
						BcEstMargin:          bcEstMargin,
					},
				},
			},
		},
	}
}
