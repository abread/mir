package abbapbevents

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule types.ModuleID, input bool) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Abba{
			Abba: &types2.Event{
				Type: &types2.Event_InputValue{
					InputValue: &types2.InputValue{
						Input: input,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, result bool, originModule types.ModuleID) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Abba{
			Abba: &types2.Event{
				Type: &types2.Event_Deliver{
					Deliver: &types2.Deliver{
						Result:       result,
						OriginModule: originModule,
					},
				},
			},
		},
	}
}
