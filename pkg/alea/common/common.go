package common

import (
	"fmt"

	aleapbCommon "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	t "github.com/filecoin-project/mir/pkg/types"
)

func FormatAleaBatchID(slot *aleapbCommon.Slot) t.BatchID {
	return t.BatchID(fmt.Sprintf("alea-%d:%d", slot.QueueIdx, slot.QueueSlot))
}
