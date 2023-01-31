package abbatypes

import t "github.com/filecoin-project/mir/pkg/types"

type RoundNumber uint64

func (r RoundNumber) Pb() uint64 {
	return uint64(r)
}

func (r RoundNumber) Bytes() []byte {
	return t.Uint64ToBytes(uint64(r))
}
