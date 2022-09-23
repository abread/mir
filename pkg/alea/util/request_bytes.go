package util

import (
	"bytes"
	"encoding/binary"

	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
)

// a unique string to ensure the signed message is really intended for Alea
const SLOT_MESSAGE_PAIR_TAG = "github.com/filecoin-project/mir/pkg/alea"

// SlotMessageData returns a [][]byte meant for signing/verifying the signature of a (slotId, payload) pair.
func SlotMessageData(slotId *aleapb.MsgId, payload *requestpb.Batch) [][]byte {
	result := make([][]byte, 0, 4+len(payload.Requests))

	result = append(result, []byte(SLOT_MESSAGE_PAIR_TAG))
	result = append(result, []byte(slotId.QueueIdx))
	result = append(result, ToBytes(slotId.Slot))
	result = append(result, ToBytes(len(payload.Requests)))

	for _, hashedReq := range payload.Requests {
		result = append(result, hashedReq.Digest)
	}

	return result
}

func ToBytes[T any](v T) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, binary.Size(v)))
	binary.Write(buf, binary.BigEndian, v)
	return buf.Bytes()
}
