package abba

import (
	"fmt"

	abbat "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	abbapbmsgs "github.com/filecoin-project/mir/pkg/pb/abbapb/msgs"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	t "github.com/filecoin-project/mir/pkg/types"
)

const (
	GlobalMsgsNs = t.ModuleID("__global")
	RoundMsgsNs  = t.ModuleID("__round")
)

func FinishMessage(moduleID t.ModuleID, value bool) *messagepbtypes.Message {
	return abbapbmsgs.FinishMessage(moduleID.Then(GlobalMsgsNs), value)
}

func InitMessage(moduleID t.ModuleID, roundNumber uint64, estimate bool) *messagepbtypes.Message {
	return abbapbmsgs.InitMessage(subidForRoundMsg(moduleID, roundNumber), roundNumber, estimate)
}

func AuxMessage(moduleID t.ModuleID, roundNumber uint64, value bool) *messagepbtypes.Message {
	return abbapbmsgs.AuxMessage(subidForRoundMsg(moduleID, roundNumber), roundNumber, value)
}

func ConfMessage(moduleID t.ModuleID, roundNumber uint64, values abbat.ValueSet) *messagepbtypes.Message {
	return abbapbmsgs.ConfMessage(subidForRoundMsg(moduleID, roundNumber), roundNumber, values)
}

func CoinMessage(moduleID t.ModuleID, roundNumber uint64, coinShare []byte) *messagepbtypes.Message {
	return abbapbmsgs.CoinMessage(subidForRoundMsg(moduleID, roundNumber), roundNumber, coinShare)
}

func subidForRoundMsg(moduleID t.ModuleID, roundNumber uint64) t.ModuleID {
	return moduleID.Then(RoundMsgsNs).Then(t.NewModuleIDFromInt(roundNumber))
}

const (
	MsgTypeFinish = "f"
	MsgTypeInit   = "i"
	MsgTypeAux    = "a"
	MsgTypeConf   = "c"
	MsgTypeCoin   = "r"
)

func FinishMsgID() rntypes.MsgID {
	return MsgTypeFinish
}

func InitMsgID(v bool) rntypes.MsgID {
	return rntypes.MsgID(fmt.Sprintf("%s.%d", MsgTypeInit, boolToNum(v)))
}

func AuxMsgID() rntypes.MsgID {
	return MsgTypeAux
}

func ConfMsgID() rntypes.MsgID {
	return MsgTypeConf
}

func CoinMsgID() rntypes.MsgID {
	return MsgTypeCoin
}

func boolToNum(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}
