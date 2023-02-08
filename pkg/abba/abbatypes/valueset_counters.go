package abbatypes

type ValueSetCounters struct {
	emptyCount   int
	zeroCount    int
	oneCount     int
	zeroOneCount int
}

func (m *ValueSetCounters) Increment(v ValueSet) {
	switch v {
	case VSetEmpty:
		m.emptyCount++
	case VSetZero:
		m.zeroCount++
	case VSetOne:
		m.oneCount++
	case VSetZeroAndOne:
		m.zeroOneCount++
	default:
		panic("invalid value set variant")
	}
}

func (m *ValueSetCounters) Get(v ValueSet) int {
	switch v {
	case VSetEmpty:
		return m.emptyCount
	case VSetZero:
		return m.zeroCount
	case VSetOne:
		return m.oneCount
	case VSetZeroAndOne:
		return m.zeroOneCount
	default:
		panic("invalid value set variant")
	}
}

func (m *ValueSetCounters) Reset() {
	m.emptyCount = 0
	m.zeroCount = 0
	m.oneCount = 0
	m.zeroOneCount = 0
}
