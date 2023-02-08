package abbatypes

type BoolCounters struct {
	falseCount int
	trueCount  int
}

func (m *BoolCounters) Increment(v bool) {
	if v {
		m.trueCount++
	} else {
		m.falseCount++
	}
}

func (m *BoolCounters) Get(v bool) int {
	if v {
		return m.trueCount
	}
	return m.falseCount
}

func (m *BoolCounters) Reset() {
	m.trueCount = 0
	m.falseCount = 0
}
