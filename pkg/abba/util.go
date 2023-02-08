package abba

import abbat "github.com/filecoin-project/mir/pkg/abba/abbatypes"

func makeBoolCounterMap() map[bool]int {
	m := make(map[bool]int, 2)
	m[false] = 0
	m[true] = 0

	return m
}

func makeValueSetCounterMap() map[abbat.ValueSet]int {
	m := make(map[abbat.ValueSet]int, 4)

	for _, v := range []abbat.ValueSet{abbat.VSetEmpty, abbat.VSetZero, abbat.VSetOne, abbat.VSetZeroAndOne} {
		m[v] = 0
	}

	return m
}

func makeBoolBoolMap() map[bool]bool {
	m := make(map[bool]bool, 2)
	m[false] = false
	m[true] = false

	return m
}

func boolToNum(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}
