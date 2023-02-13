package broadcast

import (
	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/serializing"
)

func VCBInstanceUID(aleaInstanceUID []byte, queueIdx aleatypes.QueueIdx, queueSlot aleatypes.QueueSlot) []byte {
	uid := make([]byte, len(aleaInstanceUID)+1+4+8)
	uid = append(uid, aleaInstanceUID...)
	uid = append(uid, 'b')
	uid = append(uid, serializing.Uint32ToBytes(uint32(queueIdx))...)
	uid = append(uid, serializing.Uint64ToBytes(uint64(queueSlot))...)

	return uid
}
