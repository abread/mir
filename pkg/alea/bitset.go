package alea

import "github.com/filecoin-project/mir/pkg/pb/aleapb"

type bit bool
type BitSet uint8

const (
	BITS_NONE BitSet = 0b00
	BITS_ZERO BitSet = 0b01
	BITS_ONE  BitSet = 0b10
	BITS_BOTH BitSet = 0b11
)

func bitToBitset(b bit) BitSet {
	if b {
		return BITS_ONE
	} else {
		return BITS_ZERO
	}
}

func (bs BitSet) Union(other BitSet) BitSet {
	return bs | other
}

func (bs BitSet) Intersection(other BitSet) BitSet {
	return bs & other
}

func (bs BitSet) IsEmpty() bool {
	return bs == BITS_NONE
}

func (bs BitSet) WithBit(b bit) BitSet {
	return bs.Union(bitToBitset(b))
}

func (bs BitSet) Cleared() BitSet {
	return BITS_NONE
}

func (bs BitSet) Contains(b bit) bool {
	return !(bs.Intersection(bitToBitset(b))).IsEmpty()
}

var POWERSET_BITS_NONE = []BitSet{BITS_NONE}
var POWERSET_BITS_ZERO = []BitSet{BITS_NONE, BITS_ZERO}
var POWERSET_BITS_ONE = []BitSet{BITS_NONE, BITS_ONE}
var POWERSET_BITS_BOTH = []BitSet{BITS_NONE, BITS_ZERO, BITS_ONE, BITS_BOTH}

func (bs BitSet) PowerSet() []BitSet {
	switch bs {
	case BITS_NONE:
		return POWERSET_BITS_NONE[:]
	case BITS_ZERO:
		return POWERSET_BITS_ZERO[:]
	case BITS_ONE:
		return POWERSET_BITS_ONE[:]
	case BITS_BOTH:
		return POWERSET_BITS_BOTH[:]
	default:
		panic("invalid bitset")
	}
}

func BitSetFromCobaltValueSetPb(pb aleapb.CobaltValueSet) BitSet {
	switch pb {
	case aleapb.CobaltValueSet_EMPTY:
		return BITS_NONE
	case aleapb.CobaltValueSet_ZERO:
		return BITS_ZERO
	case aleapb.CobaltValueSet_ONE:
		return BITS_ONE
	case aleapb.CobaltValueSet_ZERO_AND_ONE:
		return BITS_BOTH
	default:
		panic("invalid cobalt value set (should be a set of bits)")
	}
}

func (bs BitSet) Pb() aleapb.CobaltValueSet {
	switch bs {
	case BITS_NONE:
		return aleapb.CobaltValueSet_EMPTY
	case BITS_ZERO:
		return aleapb.CobaltValueSet_ZERO
	case BITS_ONE:
		return aleapb.CobaltValueSet_ONE
	case BITS_BOTH:
		return aleapb.CobaltValueSet_ZERO_AND_ONE
	default:
		panic("invalid bitset")
	}
}

func (b bit) Pb() bool {
	return bool(b)
}
