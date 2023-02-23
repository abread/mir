package agreement

type AgRoundHistory [][]bool

const subSliceSize = 1024

func (h *AgRoundHistory) Push(v bool) {
	if !h.hasFreeSpace() {
		h.grow()
	}

	subslice := h.tail()
	*subslice = append(*subslice, v)
}

func (h AgRoundHistory) Get(num uint64) (bool, bool) {
	if num >= h.Len() {
		return false, false
	}

	subslice := h[num/subSliceSize]
	return subslice[num%subSliceSize], true
}

func (h AgRoundHistory) Len() uint64 {
	if len(h) == 0 {
		return 0
	}

	fullSubslices := uint64(len(h)-1) * uint64(subSliceSize)
	return fullSubslices + uint64(len(*h.tail()))
}

func (h AgRoundHistory) hasFreeSpace() bool {
	if len(h) == 0 {
		return false
	}
	tail := *h.tail()
	return cap(tail) > len(tail)
}

func (h AgRoundHistory) tail() *[]bool {
	return &h[len(h)-1]
}

func (h *AgRoundHistory) grow() {
	*h = append(*h, make([]bool, 0, subSliceSize))
}
