package common

import (
	"fmt"

	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

func FormatAleaBatchID(slot *commontypes.Slot) t.BatchID {
	return t.BatchID(fmt.Sprintf("alea-%d:%d", slot.QueueIdx, slot.QueueSlot))
}
