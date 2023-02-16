package abbaround

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
)

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
