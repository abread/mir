package agreementpbmsgs

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func FinishAbbaMessage(destModule types.ModuleID, round uint64, value bool) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_AleaAgreement{
			AleaAgreement: &types2.Message{
				Type: &types2.Message_FinishAbba{
					FinishAbba: &types2.FinishAbbaMessage{
						Round: round,
						Value: value,
					},
				},
			},
		},
	}
}
