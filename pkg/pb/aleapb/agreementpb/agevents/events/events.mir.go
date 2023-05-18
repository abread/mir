// Code generated by Mir codegen. DO NOT EDIT.

package ageventsevents

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func InputValue(destModule types.ModuleID, round uint64, input bool) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaAgreement{
			AleaAgreement: &types2.Event{
				Type: &types2.Event_InputValue{
					InputValue: &types2.InputValue{
						Round: round,
						Input: input,
					},
				},
			},
		},
	}
}

func Deliver(destModule types.ModuleID, round uint64, decision bool) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaAgreement{
			AleaAgreement: &types2.Event{
				Type: &types2.Event_Deliver{
					Deliver: &types2.Deliver{
						Round:    round,
						Decision: decision,
					},
				},
			},
		},
	}
}

func StaleMsgsRecvd(destModule types.ModuleID, messages []*types3.PastMessage) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_AleaAgreement{
			AleaAgreement: &types2.Event{
				Type: &types2.Event_StaleMsgsRevcd{
					StaleMsgsRevcd: &types2.StaleMsgsRecvd{
						Messages: messages,
					},
				},
			},
		},
	}
}
