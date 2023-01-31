package abbatypes

type ValueSet uint32

const (
	VSetEmpty      ValueSet = 0b00
	VSetZero       ValueSet = 0b01
	VSetOne        ValueSet = 0b10
	VSetZeroAndOne ValueSet = 0b11
)

func (vs ValueSet) Pb() uint32 {
	return uint32(vs)
}

func boolToValueSet(v bool) ValueSet {
	if v {
		return VSetOne
	} else {
		return VSetZero
	}
}

func (vs *ValueSet) Add(v bool) {
	*vs = *vs | boolToValueSet(v)
}

func (vs ValueSet) Has(v bool) bool {
	return (vs & boolToValueSet(v)) != VSetEmpty
}

func (vs ValueSet) Len() int {
	switch vs.Sanitized() {
	case VSetEmpty:
		return 0
	case VSetZeroAndOne:
		return 2
	default:
		return 1
	}
}

func (vs ValueSet) SubsetOf(other ValueSet) bool {
	return vs.Sanitized() == other || vs.Sanitized() == VSetEmpty || other.Sanitized() == VSetZeroAndOne
}

// ValueSet.Sanitized returns a copy of the value set with all padding bits set to zero.
//
// When receiving a ValueSet from the network, we must ensure the padding bits are set to zero
// before comparing sets.
func (vs ValueSet) Sanitized() ValueSet {
	return vs & VSetZeroAndOne
}
