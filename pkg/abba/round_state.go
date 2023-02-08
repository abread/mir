package abba

import (
	abbat "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	"github.com/filecoin-project/mir/pkg/serializing"
)

type abbaRoundState struct {
	number   uint64
	estimate bool
	values   abbat.ValueSet

	initRecvd               abbat.BoolRecvTrackers
	initRecvdEstimateCounts abbat.BoolCounters

	auxRecvd            abbat.RecvTracker
	auxRecvdValueCounts abbat.BoolCounters

	confRecvd               abbat.RecvTracker
	confRecvdValueSetCounts abbat.ValueSetCounters

	coinRecvd         abbat.RecvTracker
	coinRecvdOkShares [][]byte

	initWeakSupportReachedForValue abbat.BoolFlags
	auxSent                        bool
	coinRecoverInProgress          bool
	coinRecoverMinShareCount       int
}

// resets round state, apart from the round number and estimate
func (rs *abbaRoundState) resetState(params *ModuleParams) {
	rs.values = abbat.VSetEmpty

	rs.initRecvd.Reset(params.GetN())
	rs.initRecvdEstimateCounts.Reset()

	rs.auxRecvd = make(abbat.RecvTracker, params.GetN())
	rs.auxRecvdValueCounts.Reset()

	rs.confRecvd = make(abbat.RecvTracker, params.GetN())
	rs.confRecvdValueSetCounts.Reset()

	rs.coinRecvd = make(abbat.RecvTracker, params.GetN())
	rs.coinRecvdOkShares = make([][]byte, 0, params.GetN())

	rs.initWeakSupportReachedForValue.Reset()
	rs.auxSent = false
	rs.coinRecoverInProgress = false
	rs.coinRecoverMinShareCount = params.strongSupportThresh() - 1 // TODO: review
}

// Check if there exists a subset of nodes with size >= q_S(= N-F/strong support),
// from which we have received AUX(v', r) with any v' in round.values, then broadcast CONF(values, r).
// Meant for step 7.
func (rs *abbaRoundState) isNiceAuxValueCount(params *ModuleParams) bool {
	total := 0

	for _, val := range []bool{true, false} {
		if rs.values.Has(val) {
			total += rs.auxRecvdValueCounts.Get(val)
		}
	}

	return total >= params.strongSupportThresh()
}

// Check if there exists a subset of nodes with size >= q_S(= N-F/strong support),
// from which we have received CONF(vs', r) with any vs' subset of round.values.
// Meant for step 8.
func (rs *abbaRoundState) isNiceConfValuesCount(params *ModuleParams) bool {
	total := 0

	for _, set := range []abbat.ValueSet{abbat.VSetEmpty, abbat.VSetZero, abbat.VSetOne, abbat.VSetZeroAndOne} {
		if set.SubsetOf(rs.values) {
			total += rs.confRecvdValueSetCounts.Get(set)
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
