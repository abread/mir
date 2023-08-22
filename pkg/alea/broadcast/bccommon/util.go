package bccommon

import (
	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/serializing"
	t "github.com/filecoin-project/mir/pkg/types"
)

func VCBInstanceUID(aleaBcInstanceUID []byte, queueIdx aleatypes.QueueIdx, queueSlot aleatypes.QueueSlot) []byte {
	uid := make([]byte, len(aleaBcInstanceUID)+1+4+8)
	uid = append(uid, aleaBcInstanceUID...)
	uid = append(uid, serializing.Uint32ToBytes(uint32(queueIdx))...)
	uid = append(uid, serializing.Uint64ToBytes(uint64(queueSlot))...)

	return uid
}

func BcQueueModuleID(aleaBc t.ModuleID, idx aleatypes.QueueIdx) t.ModuleID {
	return aleaBc.Then(t.NewModuleIDFromInt(idx))
}
