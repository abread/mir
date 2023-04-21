package serializing

import (
	"encoding/binary"
)

func Uint64ToBytes(n uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, n)
	return buf
}

func Uint64FromBytes(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

// Uint32ToBytes returns a 4-byte slice encoding an unsigned 32-bit integer.
func Uint32ToBytes(n uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, n)
	return buf
}
