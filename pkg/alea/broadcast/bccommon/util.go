package bccommon

import (
	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/serializing"
	t "github.com/filecoin-project/mir/pkg/types"
)

func VCBInstanceUID(aleaInstanceUID []byte, queueIdx aleatypes.QueueIdx, queueSlot aleatypes.QueueSlot) []byte {
	uid := make([]byte, 0, len(aleaInstanceUID)+1+4+8)
	uid = append(uid, aleaInstanceUID...)
	uid = append(uid, byte('b'))
	uid = append(uid, serializing.Uint32ToBytes(uint32(queueIdx))...)
	uid = append(uid, serializing.Uint64ToBytes(uint64(queueSlot))...)

	return uid
}

func BcQueueModuleID(aleaBc t.ModuleID, idx aleatypes.QueueIdx) t.ModuleID {
	return aleaBc.Then(t.NewModuleIDFromInt(idx))
}
