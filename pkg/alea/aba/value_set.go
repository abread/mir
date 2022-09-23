package aba

import "github.com/filecoin-project/mir/pkg/pb/aleapb"

type valueSet aleapb.CobaltValueSet

func newValueSet() valueSet {
	return valueSet(aleapb.CobaltValueSet_EMPTY)
}

func (vs valueSet) Pb() aleapb.CobaltValueSet {
	return aleapb.CobaltValueSet(vs)
}

func cobaltValueSetFromBool(v bool) aleapb.CobaltValueSet {
	if v {
		return aleapb.CobaltValueSet_ONE
	} else {
		return aleapb.CobaltValueSet_ZERO
	}
}

func valueSetFromBool(v bool) valueSet {
	return valueSet(cobaltValueSetFromBool(v))
}

func (vs *valueSet) Add(v bool) {
	if vs.Has(v) {
		return // already in set
	}

	if *vs == valueSet(aleapb.CobaltValueSet_EMPTY) {
		*vs = valueSetFromBool(v)
	} else {
		// vs only contains !v

		*vs = valueSet(aleapb.CobaltValueSet_ZERO_AND_ONE)
	}
}

func (vs valueSet) Has(v bool) bool {
	return vs == valueSet(aleapb.CobaltValueSet_ZERO_AND_ONE) || vs == valueSetFromBool(v)
}

func (vs valueSet) Len() int {
	switch aleapb.CobaltValueSet(vs) {
	case aleapb.CobaltValueSet_EMPTY:
		return 0
	case aleapb.CobaltValueSet_ZERO_AND_ONE:
		return 2
	default:
		return 1
	}
}

func (vs valueSet) SubsetOf(other valueSet) bool {
	return vs == other || vs == valueSet(aleapb.CobaltValueSet_EMPTY) || other == valueSet(aleapb.CobaltValueSet_ZERO_AND_ONE)
}
