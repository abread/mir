package aleatypes

type QueueIdx uint32
type QueueSlot uint64

func (qi QueueIdx) Pb() uint32 {
	return uint32(qi)
}

func (qs QueueSlot) Pb() uint64 {
	return uint64(qs)
}
