package abbapbevents

import (
	types1 "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule types.ModuleID, input bool, origin *types1.Origin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Abba{
			Abba: &types1.Event{
				Type: &types1.Event_InputValue{
					InputValue: &types1.InputValue{
						Input:  input,
						Origin: origin,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, result bool, origin *types1.Origin) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Abba{
			Abba: &types1.Event{
				Type: &types1.Event_Deliver{
					Deliver: &types1.Deliver{
						Result: result,
						Origin: origin,
					},
				},
			},
		},
	}
}

func RoundInputValue(destModule types.ModuleID, input bool) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Abba{
			Abba: &types1.Event{
				Type: &types1.Event_Round{
					Round: &types1.RoundEvent{
						Type: &types1.RoundEvent_InputValue{
							InputValue: &types1.RoundInputValue{
								Input: input,
							},
						},
					},
				},
			},
		},
	}
}

func RoundDeliver(destModule types.ModuleID, nextEstimate bool, roundNumber uint64) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Abba{
			Abba: &types1.Event{
				Type: &types1.Event_Round{
					Round: &types1.RoundEvent{
						Type: &types1.RoundEvent_Deliver{
							Deliver: &types1.RoundDeliver{
								NextEstimate: nextEstimate,
								RoundNumber:  roundNumber,
							},
						},
					},
				},
			},
		},
	}
}

func RoundFinishAll(destModule types.ModuleID, decision bool) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_Abba{
			Abba: &types1.Event{
				Type: &types1.Event_Round{
					Round: &types1.RoundEvent{
						Type: &types1.RoundEvent_Finish{
							Finish: &types1.RoundFinishAll{
								Decision: decision,
							},
						},
					},
				},
			},
		},
	}
}
