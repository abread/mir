package abba

import (
	abbat "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	"github.com/filecoin-project/mir/pkg/serializing"
)

type abbaRoundState struct {
	number   uint64
	estimate bool
	values   abbat.ValueSet

	initRecvd          map[bool]recvTracker
	initRecvdEstimates map[bool]int

	auxRecvd       recvTracker
	auxRecvdValues map[bool]int

	confRecvd       recvTracker
	confRecvdValues map[abbat.ValueSet]int

	coinRecvd         recvTracker
	coinRecvdOkShares [][]byte

	initWeakSupportReached   map[bool]bool
	auxSent                  bool
	coinRecoverInProgress    bool
	coinRecoverMinShareCount int
}

// resets round state, apart from the round number and estimate
func (rs *abbaRoundState) resetState(params *ModuleParams) {
	rs.values = abbat.VSetEmpty

	rs.initRecvd = make(map[bool]recvTracker, 2)
	for _, v := range []bool{false, true} {
		rs.initRecvd[v] = make(recvTracker, params.GetN())
	}

	rs.initRecvdEstimates = makeBoolCounterMap()

	rs.auxRecvd = make(recvTracker, params.GetN())
	rs.auxRecvdValues = makeBoolCounterMap()

	rs.confRecvd = make(recvTracker, params.GetN())
	rs.confRecvdValues = makeValueSetCounterMap()

	rs.coinRecvd = make(recvTracker, params.strongSupportThresh())
	rs.coinRecvdOkShares = make([][]byte, 0, params.GetN())

	rs.initWeakSupportReached = makeBoolBoolMap()
	rs.auxSent = false
	rs.coinRecoverInProgress = false
	rs.coinRecoverMinShareCount = params.strongSupportThresh() - 1
}

// Check if there exists a subset of nodes with size >= q_S(= N-F/strong support),
// from which we have received AUX(v', r) with any v' in round.values, then broadcast CONF(values, r).
// Meant for step 7.
func (rs *abbaRoundState) isNiceAuxValueCount(params *ModuleParams) bool {
	total := 0

	for val, count := range rs.auxRecvdValues {
		if rs.values.Has(val) {
			total += count
		}
	}

	return total >= params.strongSupportThresh()
}

// Check if there exists a subset of nodes with size >= q_S(= N-F/strong support),
// from which we have received CONF(vs', r) with any vs' subset of round.values.
// Meant for step 8.
func (rs *abbaRoundState) isNiceConfValuesCount(params *ModuleParams) bool {
	total := 0

	for set, count := range rs.confRecvdValues {
		if set.SubsetOf(rs.values) {
			total += count
		}
	}

	return total >= params.strongSupportThresh()
}

const CoinSignDataPrefix = "github.com/filecoin-project/mir/pkg/alea/aba"

func (rs *abbaRoundState) coinData(params *ModuleParams) [][]byte {
	return [][]byte{
		[]byte(CoinSignDataPrefix),
		params.InstanceUID,
		serializing.Uint64ToBytes(rs.number),
	}
}
