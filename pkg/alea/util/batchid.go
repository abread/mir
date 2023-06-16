package util

import (
	"fmt"

	msct "github.com/filecoin-project/mir/pkg/availability/multisigcollector/types"
	commontypes "github.com/filecoin-project/mir/pkg/pb/aleapb/common/types"
)

func FormatAleaBatchID(slot *commontypes.Slot) msct.BatchID {
	return fmt.Sprintf("alea-%d:%d", slot.QueueIdx, slot.QueueSlot)
}
