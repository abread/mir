package agreement

type AgRoundHistory [][]bool

const subSliceSize = 1024

func (h *AgRoundHistory) Push(v bool) {
	if h.cap() == 0 {
		h.grow()
	}

	subslice := h.tail()
	*subslice = append(*subslice, v)
}

func (h *AgRoundHistory) Get(num uint64) (bool, bool) {
	if num >= h.Len() {
		return false, false
	}

	return (*h)[num/subSliceSize][num%subSliceSize], true
}

func (h *AgRoundHistory) Len() uint64 {
	fullSubslices := uint64(len(*h)) * uint64(subSliceSize)
	if fullSubslices != 0 {
		fullSubslices += uint64(len(*h.tail()))
	}

	return fullSubslices
}

func (h *AgRoundHistory) cap() int {
	if len(*h) == 0 {
		return 0
	}
	return cap((*h)[len(*h)-1])
}

func (h *AgRoundHistory) tail() *[]bool {
	return &(*h)[len(*h)-1]
}

func (h *AgRoundHistory) grow() {
	*h = append(*h, make([]bool, 0, subSliceSize))
}
