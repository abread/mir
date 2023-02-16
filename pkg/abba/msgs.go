package abba

import (
	"github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
)

const (
	MsgTypeFinish = "f"
)

func FinishMsgID() rntypes.MsgID {
	return MsgTypeFinish
}
