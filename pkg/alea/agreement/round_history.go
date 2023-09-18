package agreement

type AgRoundHistory struct {
	epochSize int
	minEpoch  uint64 // first stored epoch
	history   [][]bool
}

func NewAgRoundHistory(epochSize int) AgRoundHistory {
	history := make([][]bool, 1)
	history[0] = make([]bool, 0, epochSize)

	return AgRoundHistory{
		epochSize: epochSize,
		minEpoch:  0,
		history:   history,
	}
}

func (h *AgRoundHistory) Store(sn uint64, v bool) {
	epoch := sn / uint64(h.epochSize)
	idxInEpoch := int(sn % uint64(h.epochSize))

	if epoch == h.minEpoch+uint64(len(h.history)) && idxInEpoch == 0 && len(h.history[len(h.history)-1]) == h.epochSize {
		// first of new epoch
		h.history = append(h.history, make([]bool, 0, h.epochSize))
		h.history[len(h.history)-1] = append(h.history[len(h.history)-1], v)
	} else if epoch == h.minEpoch+uint64(len(h.history))-1 && idxInEpoch == len(h.history[len(h.history)-1]) {
		// next of current epoch
		h.history[len(h.history)-1] = append(h.history[len(h.history)-1], v)
	} else {
		panic("agroundhistory.store: unexpected sn - in the future or in the past")
	}
}

func (h AgRoundHistory) Get(sn uint64) (bool, bool) {
	epoch := sn / uint64(h.epochSize)
	idxInEpoch := int(sn % uint64(h.epochSize))

	if epoch < h.minEpoch || epoch > h.minEpoch+uint64(len(h.history)) {
		// epoch is out of range
		return false, false
	} else if idxInEpoch >= len(h.history[epoch-h.minEpoch]) {
		// idxInEpoch is out of range
		return false, false
	} else {
		return h.history[epoch-h.minEpoch][idxInEpoch], true
	}
}

func (h AgRoundHistory) Len() int {
	if len(h.history) > 0 {
		pastEpochs := (len(h.history) - 1) * h.epochSize
		return pastEpochs + len(h.history[len(h.history)-1])
	}

	return 0
}

func (h *AgRoundHistory) FreeEpochs(newMinEpoch uint64) {
	if newMinEpoch > h.minEpoch {
		h.history = h.history[newMinEpoch-h.minEpoch:]
		h.minEpoch = newMinEpoch

		if len(h.history) == 0 {
			// allocate slice for new epoch
			h.history = append(h.history, make([]bool, 0, h.epochSize))
		}
	}
}
