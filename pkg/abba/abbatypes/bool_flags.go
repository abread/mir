package abbatypes

type BoolFlags struct {
	falseVal bool
	trueVal  bool
}

func (m *BoolFlags) Set(v bool) {
	if v {
		m.trueVal = true
	} else {
		m.falseVal = true
	}
}

func (m *BoolFlags) Get(v bool) bool {
	if v {
		return m.trueVal
	}
	return m.falseVal
}

func (m *BoolFlags) Reset() {
	m.trueVal = false
	m.falseVal = false
}
