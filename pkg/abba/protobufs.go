package abba

import (
	"github.com/filecoin-project/mir/pkg/abba/abbadsl"
	"github.com/filecoin-project/mir/pkg/pb/abbapb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func Message(moduleID t.ModuleID, msg *abbapb.Message) *messagepb.Message {
	return &messagepb.Message{
		DestModule: moduleID.Pb(),
		Type: &messagepb.Message_Abba{
			Abba: msg,
		},
	}
}

func FinishMessage(moduleID t.ModuleID, value bool) *messagepb.Message {
	return Message(moduleID, &abbapb.Message{
		Type: &abbapb.Message_FinishMessage{
			FinishMessage: &abbapb.FinishMessage{
				Value: value,
			},
		},
	})
}

func InitMessage(moduleID t.ModuleID, roundNumber uint64, estimate bool) *messagepb.Message {
	return Message(moduleID, &abbapb.Message{
		Type: &abbapb.Message_InitMessage{
			InitMessage: &abbapb.InitMessage{
				RoundNumber: roundNumber,
				Estimate:    estimate,
			},
		},
	})
}

func AuxMessage(moduleID t.ModuleID, roundNumber uint64, value bool) *messagepb.Message {
	return Message(moduleID, &abbapb.Message{
		Type: &abbapb.Message_AuxMessage{
			AuxMessage: &abbapb.AuxMessage{
				RoundNumber: roundNumber,
				Value:       value,
			},
		},
	})
}

func ConfMessage(moduleID t.ModuleID, roundNumber uint64, values abbadsl.ValueSet) *messagepb.Message {
	return Message(moduleID, &abbapb.Message{
		Type: &abbapb.Message_ConfMessage{
			ConfMessage: &abbapb.ConfMessage{
				RoundNumber: roundNumber,
				Values:      values.Pb(),
			},
		},
	})
}

func CoinMessage(moduleID t.ModuleID, roundNumber uint64, coinShare []byte) *messagepb.Message {
	return Message(moduleID, &abbapb.Message{
		Type: &abbapb.Message_CoinMessage{
			CoinMessage: &abbapb.CoinMessage{
				RoundNumber: roundNumber,
				CoinShare:   coinShare,
			},
		},
	})
}
