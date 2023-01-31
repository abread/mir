package abbapbmsgs

import (
	abbatypes "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	types2 "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func FinishMessage(destModule types.ModuleID, value bool) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Abba{
			Abba: &types2.Message{
				Type: &types2.Message_FinishMessage{
					FinishMessage: &types2.FinishMessage{
						Value: value,
					},
				},
			},
		},
	}
}

func InitMessage(destModule types.ModuleID, roundNumber uint64, estimate bool) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Abba{
			Abba: &types2.Message{
				Type: &types2.Message_InitMessage{
					InitMessage: &types2.InitMessage{
						RoundNumber: roundNumber,
						Estimate:    estimate,
					},
				},
			},
		},
	}
}

func AuxMessage(destModule types.ModuleID, roundNumber uint64, value bool) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Abba{
			Abba: &types2.Message{
				Type: &types2.Message_AuxMessage{
					AuxMessage: &types2.AuxMessage{
						RoundNumber: roundNumber,
						Value:       value,
					},
				},
			},
		},
	}
}

func ConfMessage(destModule types.ModuleID, roundNumber uint64, values abbatypes.ValueSet) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Abba{
			Abba: &types2.Message{
				Type: &types2.Message_ConfMessage{
					ConfMessage: &types2.ConfMessage{
						RoundNumber: roundNumber,
						Values:      values,
					},
				},
			},
		},
	}
}

func CoinMessage(destModule types.ModuleID, roundNumber uint64, coinShare []uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Abba{
			Abba: &types2.Message{
				Type: &types2.Message_CoinMessage{
					CoinMessage: &types2.CoinMessage{
						RoundNumber: roundNumber,
						CoinShare:   coinShare,
					},
				},
			},
		},
	}
}
