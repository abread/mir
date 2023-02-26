package abbapbmsgs

import (
	abbatypes "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	types2 "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func FinishMessage(destModule types.ModuleID, value bool) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Abba{
			Abba: &types2.Message{
				Type: &types2.Message_Finish{
					Finish: &types2.FinishMessage{
						Value: value,
					},
				},
			},
		},
	}
}

func RoundInitMessage(destModule types.ModuleID, estimate bool) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Abba{
			Abba: &types2.Message{
				Type: &types2.Message_Round{
					Round: &types2.RoundMessage{
						Type: &types2.RoundMessage_Init{
							Init: &types2.RoundInitMessage{
								Estimate: estimate,
							},
						},
					},
				},
			},
		},
	}
}

func RoundAuxMessage(destModule types.ModuleID, value bool) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Abba{
			Abba: &types2.Message{
				Type: &types2.Message_Round{
					Round: &types2.RoundMessage{
						Type: &types2.RoundMessage_Aux{
							Aux: &types2.RoundAuxMessage{
								Value: value,
							},
						},
					},
				},
			},
		},
	}
}

func RoundConfMessage(destModule types.ModuleID, values abbatypes.ValueSet) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Abba{
			Abba: &types2.Message{
				Type: &types2.Message_Round{
					Round: &types2.RoundMessage{
						Type: &types2.RoundMessage_Conf{
							Conf: &types2.RoundConfMessage{
								Values: values,
							},
						},
					},
				},
			},
		},
	}
}

func RoundCoinMessage(destModule types.ModuleID, coinShare tctypes.SigShare) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Abba{
			Abba: &types2.Message{
				Type: &types2.Message_Round{
					Round: &types2.RoundMessage{
						Type: &types2.RoundMessage_Coin{
							Coin: &types2.RoundCoinMessage{
								CoinShare: coinShare,
							},
						},
					},
				},
			},
		},
	}
}