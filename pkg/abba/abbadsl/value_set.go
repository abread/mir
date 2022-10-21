package abbadsl

import "github.com/filecoin-project/mir/pkg/pb/abbapb"

type ValueSet abbapb.ValueSet

func EmptyValueSet() ValueSet {
	return ValueSet(abbapb.ValueSet_EMPTY)
}

func (vs ValueSet) Pb() abbapb.ValueSet {
	return abbapb.ValueSet(vs)
}

func cobaltValueSetFromBool(v bool) abbapb.ValueSet {
	if v {
		return abbapb.ValueSet_ONE
	}

	return abbapb.ValueSet_ZERO
}

func valueSetFromBool(v bool) ValueSet {
	return ValueSet(cobaltValueSetFromBool(v))
}

func (vs *ValueSet) Add(v bool) {
	if vs.Has(v) {
		return // already in set
	}

	if *vs == ValueSet(abbapb.ValueSet_EMPTY) {
		*vs = valueSetFromBool(v)
	} else {
		// vs only contains !v

		*vs = ValueSet(abbapb.ValueSet_ZERO_AND_ONE)
	}
}

func (vs ValueSet) Has(v bool) bool {
	return vs == ValueSet(abbapb.ValueSet_ZERO_AND_ONE) || vs == valueSetFromBool(v)
}

func (vs ValueSet) Len() int {
	switch abbapb.ValueSet(vs) {
	case abbapb.ValueSet_EMPTY:
		return 0
	case abbapb.ValueSet_ZERO_AND_ONE:
		return 2
	default:
		return 1
	}
}

func (vs ValueSet) SubsetOf(other ValueSet) bool {
	return vs == other || vs == ValueSet(abbapb.ValueSet_EMPTY) || other == ValueSet(abbapb.ValueSet_ZERO_AND_ONE)
}
