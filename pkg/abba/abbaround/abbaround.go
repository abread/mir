package abbaround

import (
	"runtime"
	"time"

	es "github.com/go-errors/errors"

	abbat "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	abbadsl "github.com/filecoin-project/mir/pkg/pb/abbapb/dsl"
	abbapbevents "github.com/filecoin-project/mir/pkg/pb/abbapb/events"
	abbapbmsgs "github.com/filecoin-project/mir/pkg/pb/abbapb/msgs"
	hasherpbdsl "github.com/filecoin-project/mir/pkg/pb/hasherpb/dsl"
	hasherpbtypes "github.com/filecoin-project/mir/pkg/pb/hasherpb/types"
	rnetdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	threshDsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tsagg"
	t "github.com/filecoin-project/mir/pkg/types"
)

type roundPhase uint8

const (
	phaseAwaitingInput roundPhase = iota // immediately before protocol step 4

	// nodes in the first round will stall in this phase until
	// they receive RoundContinue, or skip to the next phase
	// straight away if they are not in the first round
	phaseAwaitingContinue // same as before, but with input already provided

	phaseAwaitingNiceAux  // until protocol step 7 completes (until CONF is broadcast)
	phaseAwaitingNiceConf // until protocol step 8 completes (until we receive enough CONF)
	phaseTossingCoin      // protocol step 9, while preparing coin message
	phaseDone             // protocol step 10 (coin was tossed, next estimate was delivered)
)

var timeRef = time.Now()

type state struct {
	phase           roundPhase
	estimate        bool
	values          abbat.ValueSet
	requestedFinish bool

	initRecvd               abbat.BoolRecvTrackers
	initRecvdEstimateCounts abbat.BoolCounters

	auxRecvd            abbat.RecvTracker
	auxRecvdValueCounts abbat.BoolCounters

	confRecvd               abbat.RecvTracker
	confRecvdValueSetCounts abbat.ValueSetCounters

	ownCoinShare tctypes.SigShare
	coinSig      *tsagg.ThreshSigAggregator
	hashingCoin  bool

	initWeakSupportReachedForValue abbat.BoolFlags
	auxSent                        bool

	relStartTime            time.Duration
	relCoinRecoverStartTime time.Duration
}

// nolint: gocognit,gocyclo
func New(mc ModuleConfig, params ModuleParams, nodeID t.NodeID, logger logging.Logger) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	coinData := genCoinData(mc.Self, params.InstanceUID)
	state := state{
		phase: phaseAwaitingInput,

		initRecvd: abbat.NewBoolRecvTrackers(len(params.AllNodes)),
		auxRecvd:  make(abbat.RecvTracker, len(params.AllNodes)),
		confRecvd: make(abbat.RecvTracker, len(params.AllNodes)),
		coinSig: tsagg.New(m, &tsagg.Params{
			TCModuleID:              mc.ThreshCrypto,
			Threshold:               params.strongSupportThresh(),
			MaxVerifyShareBatchSize: runtime.NumCPU(),
			SigData: func() [][]byte {
				return coinData
			},
			InitialNodeCount: len(params.AllNodes),
		}, logging.Decorate(logger, "ThresholdSigAggregator: ")),

		relStartTime: time.Since(timeRef),
	}

	abbadsl.UponRoundInputValue(m, func(input bool) error {
		// 10. est^r+1_i = v OR 3. set est^r_i = v_in
		state.estimate = input

		// 4. Broadcast INIT(r, est_r_i)
		// changed to INPUT(r, est_r_i) to allow the unaminimity optimization to work
		rnetdsl.SendMessage(m, mc.ReliableNet, InputMsgID(), abbapbmsgs.RoundInputMessage(
			mc.Self,
			state.estimate,
		), params.AllNodes)

		state.phase = phaseAwaitingContinue

		return nil
	})

	if params.RoundNumber == 0 {
		abbadsl.UponRoundContinue(m, func() error {
			// no more time, we must go on with ABBA, with unanimity or not
			state.phase = phaseAwaitingNiceAux

			// precompute coin sig share
			threshDsl.SignShare[struct{}](m, mc.ThreshCrypto, coinData, nil)
			return nil
		})
	} else {
		abbadsl.UponRoundInputValue(m, func(input bool) error {
			// skip ahead to sending AUX straight away
			state.phase = phaseAwaitingNiceAux

			// precompute coin sig share
			threshDsl.SignShare[struct{}](m, mc.ThreshCrypto, coinData, nil)
			return nil
		})
	}

	threshDsl.UponSignShareResult(m, func(sigShare tctypes.SigShare, _ctx *struct{}) error {
		state.ownCoinShare = sigShare
		state.coinSig.Add(sigShare, nodeID)
		return nil
	})

	registerInitVal := func(from t.NodeID, est bool) {
		if state.initRecvd.Register(est, from) {
			state.initRecvdEstimateCounts.Increment(est)
		}
	}
	abbadsl.UponRoundInputMessageReceived(m, func(from t.NodeID, est bool) error {
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, InputMsgID(), from)
		registerInitVal(from, est)
		return nil
	})
	abbadsl.UponRoundInitMessageReceived(m, func(from t.NodeID, est bool) error {
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, InitMsgID(est), from)
		registerInitVal(from, est)
		return nil
	})

	dsl.UponStateUpdates(m, func() error {
		if state.phase <= phaseAwaitingContinue || state.phase == phaseDone {
			return nil // did not perform the initial INIT broadcast or has already terminated
		}

		for _, est := range []bool{true, false} {
			// 5. upon receiving weak support for INIT(r, v), add v to values and broadcast INIT(r, v)
			if !state.initWeakSupportReachedForValue.Get(est) && state.initRecvdEstimateCounts.Get(est) >= params.weakSupportThresh() {
				// logger.Log(logging.LevelDebug, "received weak support for INIT(est)", "est", est)

				state.values.Add(est)

				// if our round.estimate is r, then we already sent this message
				if state.estimate != est {
					rnetdsl.SendMessage(m, mc.ReliableNet,
						InitMsgID(est),
						abbapbmsgs.RoundInitMessage(mc.Self, est),
						params.AllNodes,
					)
				}

				state.initWeakSupportReachedForValue.Set(est)
			}

			// 6. upon receiving strong support for INIT(r, v), broadcast AUX(r, v) if we have not already broadcast AUX(r, _)
			if !state.auxSent && state.initRecvdEstimateCounts.Get(est) >= params.strongSupportThresh() {
				// logger.Log(logging.LevelDebug, "received strong support for INIT(est)", "est", est)
				rnetdsl.SendMessage(m, mc.ReliableNet,
					AuxMsgID(),
					abbapbmsgs.RoundAuxMessage(mc.Self, est),
					params.AllNodes,
				)
				state.auxSent = true
			}
		}

		return nil
	})

	abbadsl.UponRoundAuxMessageReceived(m, func(from t.NodeID, value bool) error {
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, AuxMsgID(), from)

		if !state.auxRecvd.Register(from) {
			// logger.Log(logging.LevelWarn, "duplicate AUX(_)", "from", from)
			return nil // duplicate message
		}
		state.auxRecvdValueCounts.Increment(value)

		// logger.Log(logging.LevelDebug, "recvd AUX(v)", "v", value)

		return nil
	})

	// 7. wait until there exists a subset of nodes with size >= q_S(= N-F), from which we have received AUX(v', r) with any v' in round.values, then broadcast CONF(values, r)
	dsl.UponStateUpdates(m, func() error {
		if state.phase != phaseAwaitingNiceAux {
			return nil
		}

		if state.isNiceAuxValueCount(&params) {
			// logger.Log(logging.LevelDebug, "received enough support for AUX(v in values)", "values", state.values)
			rnetdsl.SendMessage(m, mc.ReliableNet,
				ConfMsgID(),
				abbapbmsgs.RoundConfMessage(mc.Self, state.values),
				params.AllNodes,
			)

			state.phase = phaseAwaitingNiceConf
		}

		return nil
	})

	abbadsl.UponRoundConfMessageReceived(m, func(from t.NodeID, values abbat.ValueSet) error {
		values = values.Sanitized()
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, ConfMsgID(), from)

		if !state.confRecvd.Register(from) {
			// logger.Log(logging.LevelWarn, "duplicate CONF(_)", "from", from)
			return nil // duplicate message
		}

		// logger.Log(logging.LevelDebug, "recvd CONF(C)", "C", values)
		state.confRecvdValueSetCounts.Increment(values)

		return nil
	})

	// 8. wait until there exists a subset of nodes with size >= q_S(= N-F), from which we have received CONF(vs', r) with any vs' subset of round.values
	dsl.UponStateUpdates(m, func() error {
		if state.phase != phaseAwaitingNiceConf {
			return nil
		}

		if state.isNiceConfValuesCount(&params) && state.relCoinRecoverStartTime == 0 {
			state.relCoinRecoverStartTime = time.Since(timeRef)
		}

		if state.isNiceConfValuesCount(&params) && state.ownCoinShare != nil {
			// logger.Log(logging.LevelDebug, "received enough support for CONF(C subset of values)", "values", state.values)

			// 9. sample coin
			state.phase = phaseTossingCoin

			// logger.Log(logging.LevelDebug, "tossing coin", "ownShare", state.ownCoinShare)
			rnetdsl.SendMessage(m, mc.ReliableNet,
				CoinMsgID(),
				abbapbmsgs.RoundCoinMessage(mc.Self, state.ownCoinShare),
				params.AllNodes,
			)
		}

		return nil
	})

	// working in advance for 9. sample coin
	abbadsl.UponRoundCoinMessageReceived(m, func(from t.NodeID, coinShare tctypes.SigShare) error {
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, CoinMsgID(), from)

		if state.phase == phaseDone {
			// logger.Log(logging.LevelWarn, "already terminated (recvd COIN)")
			return nil // already terminated
		}

		state.coinSig.Add(coinShare, from)
		// logger.Log(logging.LevelDebug, "recvd COIN(share)", "from", from)
		return nil
	})

	// still in 9. sample coin
	dsl.UponStateUpdates(m, func() error {
		if state.phase != phaseTossingCoin {
			return nil
		}

		if state.coinSig.FullSig() != nil && !state.hashingCoin {
			hasherpbdsl.RequestOne[struct{}](m, mc.Hasher, &hasherpbtypes.HashData{
				Data: [][]byte{state.coinSig.FullSig()},
			}, nil)
			state.hashingCoin = true
		}

		return nil
	})

	// still in 9. sample coin
	hasherpbdsl.UponResultOne(m, func(hash []byte, _ctx *struct{}) error {
		if state.phase != phaseTossingCoin {
			return es.Errorf("impossible state reached: coin being tossed but round is already over")
		}

		// finishing step 9
		sR := (hash[0] & 1) == 1 // TODO: this is ok, right?

		// 10.
		if state.values.Len() == 2 {
			state.estimate = sR
		} else if state.values.Len() == 1 { // values = {v}
			v := state.values.Has(true) // if values contains true, v=true, otherwise v=false
			state.estimate = v

			// If in fact values = { s_r }, broadcast FINISH(s_r) if we haven't broadcast FINISH(_) already
			// request ABBA controller to broadcast FINISH (if not done already)
			if v == sR && !state.requestedFinish {
				abbadsl.RoundFinishAll(m, mc.Consumer, sR, false)
				state.requestedFinish = true
			}
		}

		durationNoCoinRecover := state.relCoinRecoverStartTime - state.relStartTime

		// (still 10.) Set r = r + 1, and return to step 4
		dsl.EmitEvent(m, abbapbevents.RoundDeliver(
			mc.Consumer,
			state.estimate,
			params.RoundNumber,
			durationNoCoinRecover,
		))

		state.phase = phaseDone

		return nil
	})

	// unanimity optimization
	// when *all nodes input the same value* to abba, it is guaranteed that they will output that value
	// inform abba module of that condition
	// NOTE: we consider the first INIT(v) message received from a node to be its input. This is ok,
	// because it is impossible for a correct node to broadcast INIT(0) and INIT(1) when correct nodes
	// are in unanimity (because there are only up to F INIT(other value) from byz nodes received).
	if params.RoundNumber == 0 {
		inputRecvTracker := make(abbat.RecvTracker, len(params.AllNodes))
		trueCount := 0
		falseCount := 0

		abbadsl.UponRoundInputMessageReceived(m, func(from t.NodeID, estimate bool) error {
			if !inputRecvTracker.Register(from) {
				return nil
			}

			if estimate {
				trueCount++
			} else {
				falseCount++
			}

			return nil
		})

		dsl.UponStateUpdates(m, func() error {
			if state.requestedFinish {
				return nil
			}

			if trueCount == len(params.AllNodes) {
				abbadsl.RoundFinishAll(m, mc.Consumer, true, true)
				state.requestedFinish = true
			} else if falseCount == len(params.AllNodes) {
				abbadsl.RoundFinishAll(m, mc.Consumer, false, true)
				state.requestedFinish = true
			}

			return nil
		})
	}

	return m
}

// Check if there exists a subset of nodes with size >= q_S(= N-F/strong support),
// from which we have received AUX(v', r) with any v' in round.values, then broadcast CONF(values, r).
// Meant for step 7.
func (rs *state) isNiceAuxValueCount(params *ModuleParams) bool {
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
func (rs *state) isNiceConfValuesCount(params *ModuleParams) bool {
	total := 0

	for _, set := range []abbat.ValueSet{abbat.VSetEmpty, abbat.VSetZero, abbat.VSetOne, abbat.VSetZeroAndOne} {
		if set.SubsetOf(rs.values) {
			total += rs.confRecvdValueSetCounts.Get(set)
		}
	}

	return total >= params.strongSupportThresh()
}

func genCoinData(ownModID t.ModuleID, instanceUID []byte) [][]byte {
	return [][]byte{
		instanceUID,
		[]byte(ownModID),
		[]byte("coin"),
	}
}
