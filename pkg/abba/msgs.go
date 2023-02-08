package abba

import (
	abbat "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	abbapbmsgs "github.com/filecoin-project/mir/pkg/pb/abbapb/msgs"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
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
	MsgTypeFinish uint8 = iota
	MsgTypeInit
	MsgTypeAux
	MsgTypeConf
	MsgTypeCoin
)

func FinishMsgID() []byte {
	return []byte{MsgTypeFinish}
}

func InitMsgID(v bool) []byte {
	return []byte{MsgTypeInit, boolToNum(v)}
}

func AuxMsgID() []byte {
	return []byte{MsgTypeAux}
}

func ConfMsgID() []byte {
	return []byte{MsgTypeConf}
}

func CoinMsgID() []byte {
	return []byte{MsgTypeCoin}
}
