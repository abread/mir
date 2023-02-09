package broadcast

import "github.com/filecoin-project/mir/pkg/serializing"

func VCBInstanceUID(aleaInstanceUID []byte, queueIdx uint32, queueSlot uint64) []byte {
	uid := make([]byte, len(aleaInstanceUID)+1+4+8)
	uid = append(uid, aleaInstanceUID...)
	uid = append(uid, 'b')
	uid = append(uid, serializing.Uint32ToBytes(queueIdx)...)
	uid = append(uid, serializing.Uint64ToBytes(queueSlot)...)

	return uid
}
