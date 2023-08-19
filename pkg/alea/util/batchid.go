package util

import (
	"fmt"

	msct "github.com/filecoin-project/mir/pkg/availability/multisigcollector/types"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
)

func FormatAleaBatchID(slot *bcpbtypes.Slot) msct.BatchID {
	return fmt.Sprintf("alea-%d:%d", slot.QueueIdx, slot.QueueSlot)
}
