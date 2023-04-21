package abbatypes

import (
	"github.com/filecoin-project/mir/pkg/serializing"
)

type RoundNumber uint64

func (r RoundNumber) Pb() uint64 {
	return uint64(r)
}

func (r RoundNumber) Bytes() []byte {
	return serializing.Uint64ToBytes(uint64(r))
}
