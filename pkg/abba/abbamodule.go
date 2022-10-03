package abba

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/filecoin-project/mir/pkg/abba/abbadsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/factorymodule"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/abbapb"
	"github.com/filecoin-project/mir/pkg/pb/factorymodulepb"
	threshDsl "github.com/filecoin-project/mir/pkg/threshcrypto/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID // id of the module to send the "Deliver" event to
	Net          t.ModuleID
	ThreshCrypto t.ModuleID
	Hasher       t.ModuleID
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // unique identifier for this instance of VCB
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
}

// GetN returns the total number of nodes.
func (params *ModuleParams) GetN() int {
	return len(params.AllNodes)
}

// GetF returns the maximum tolerated number of faulty nodes.
func (params *ModuleParams) GetF() int {
	return (params.GetN() - 1) / 3
}

type abbaModuleState struct {
	round abbaRoundState

	step uint8 // next/in-progress protocol step

	finishRecvd       recvTracker
	finishRecvdValues map[bool]int

	finishSent bool
}

type abbaRoundState struct {
	number   uint64
	estimate bool
	values   abbadsl.ValueSet

	initRecvd          recvTracker
	initRecvdEstimates map[bool]int

	auxRecvd       recvTracker
	auxRecvdValues map[bool]int

	confRecvd       recvTracker
	confRecvdValues map[abbadsl.ValueSet]int

	coinRecvd         recvTracker
	coinRecvdOkShares [][]byte

	initSent                 map[bool]bool
	auxSent                  bool
	coinRecoverInProgress    bool
	coinRecoverMinShareCount int
}

type signCoinShareCtx struct {
	roundNumber uint64
}

type verifyCoinShareCtx struct {
	roundNumber uint64
	sigShare    []byte
}

type recoverCoinCtx struct {
	roundNumber uint64
}

const MAX_STEP uint8 = 10

func NewReconfigurableModule(mc *ModuleConfig, nodeID t.NodeID, logger logging.Logger) modules.PassiveModule {
	return factorymodule.New(
		mc.Self,
		factorymodule.DefaultParams(

			// This function will be called whenever the factory module
			// is asked to create a new instance of the vcb protocol.
			func(abbaID t.ModuleID, genericParams *factorymodulepb.GeneratorParams) (modules.PassiveModule, error) {
				params := genericParams.Type.(*factorymodulepb.GeneratorParams_Abba).Abba

				// Extract the IDs of the nodes in the membership associated with this instance
				allNodes := maputil.GetSortedKeys(t.Membership(params.Membership))

				// Crate a copy of basic module config with an adapted ID for the submodule.
				submc := *mc
				submc.Self = abbaID

				// Create a new instance of the vcb protocol.
				inst := NewModule(
					&submc,
					&ModuleParams{
						// TODO: Use InstanceUIDs properly.
						//       (E.g., concatenate this with the instantiating protocol's InstanceUID when introduced.)
						InstanceUID: []byte(abbaID),
						AllNodes:    allNodes,
					},
					nodeID,
				)
				return inst, nil
			},
		),
		logger,
	)
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)
	state := &abbaModuleState{
		round: abbaRoundState{
			number: 0,
		},

		finishRecvd:       make(recvTracker, params.GetN()),
		finishRecvdValues: makeBoolCounterMap(),

		finishSent: false,
	}

	state.round.resetState(params)

	abbadsl.UponFinishMessageReceived(m, func(from t.NodeID, value bool) error {
		if _, present := state.finishRecvd[from]; present {
			return nil // duplicate message
		}
		state.finishRecvd[from] = struct{}{}
		state.finishRecvdValues[value] += 1

		// 1. upon receiving weak support for FINISH(v), broadcast FINISH(v)
		if !state.finishSent && state.finishRecvdValues[value] >= params.weakSupportThresh() {
			dsl.SendMessage(m, mc.Net, FinishMessage(mc.Self, value), params.AllNodes)
			state.finishSent = true
		}

		// 2. upon receiving strong support for FINISH(v), output v and terminate
		if state.step <= 10 && state.finishRecvdValues[value] >= params.strongSupportThresh() {
			abbadsl.Deliver(m, mc.Consumer, value)
			state.step = math.MaxUint8 // no more progress can be made
		}

		return nil
	})

	state.step = 3

	// 3. upon P_i providing input value v_in, set est^r_i=v_in, r=0
	abbadsl.UponInputValue(m, func(input bool) error {
		if state.step > 3 {
			return fmt.Errorf("Input value already provided for ABBA instance")
		}

		state.round.number = 0
		state.round.estimate = input

		// 4. broadcast INIT(r, est^r_i)
		msg := InitMessage(mc.Self, state.round.number, state.round.estimate)
		dsl.SendMessage(m, mc.Net, msg, params.AllNodes)

		state.step = 5
		return nil
	})

	abbadsl.UponInitMessageReceived(m, func(from t.NodeID, r uint64, est bool) error {
		if r != state.round.number {
			return nil // not ready yet or wrong round
		}
		if _, present := state.round.initRecvd[from]; present {
			return nil // duplicate message
		}
		state.round.initRecvd[from] = struct{}{}
		state.round.initRecvdEstimates[est] += 1

		if state.step < 5 || state.step > MAX_STEP {
			return nil // not ready yet for the next steps (or already delivered)
		}

		// 5. upon receiving weak support for INIT(r, v), broadcast INIT(r, v)
		if !state.round.initSent[est] && state.round.initRecvdEstimates[est] == params.weakSupportThresh() {
			dsl.SendMessage(m, mc.Net, InitMessage(mc.Self, r, est), params.AllNodes)
			state.step = 6
		}

		// 6. upon receiving strong support for INIT(r, v), broadcast AUX(r, v)
		if !state.round.auxSent && state.round.initRecvdEstimates[est] == params.strongSupportThresh() {
			dsl.SendMessage(m, mc.Net, AuxMessage(mc.Self, r, est), params.AllNodes)
			state.step = 7
		}

		return nil
	})

	abbadsl.UponAuxMessageReceived(m, func(from t.NodeID, r uint64, value bool) error {
		if r != state.round.number {
			return nil // not processing this round
		}
		if _, present := state.round.auxRecvd[from]; present {
			return nil // duplicate message
		}
		state.round.auxRecvd[from] = struct{}{}
		state.round.auxRecvdValues[value] += 1

		// 7. wait until there exists a subset of nodes with size >= q_S(= N-F), from which we have received AUX(v', r) with any v' in round.values, then broadcast CONF(values, r)
		if state.step == 7 && state.round.isNiceAuxValueCount(params) {
			dsl.SendMessage(m, mc.Net, ConfMessage(mc.Self, r, state.round.values), params.AllNodes)
			state.step = 8
		}

		return nil
	})

	abbadsl.UponConfMessageReceived(m, func(from t.NodeID, r uint64, values abbadsl.ValueSet) error {
		if r != state.round.number {
			return nil // wrong round
		}
		if _, present := state.round.confRecvd[from]; present {
			return nil // duplicate message
		}
		state.round.confRecvd[from] = struct{}{}
		state.round.confRecvdValues[values] += 1

		// 8. wait until there exists a subset of nodes with size >= q_S(= N-F), from which we have received CONF(vs', r) with any vs' subset of round.values
		if state.step == 8 && state.round.isNiceConfValuesCount(params) {

			// 9. sample coin
			state.step = 9
			threshDsl.SignShare(m, mc.ThreshCrypto, state.round.coinData(params), &signCoinShareCtx{
				roundNumber: state.round.number,
			})
		}

		return nil
	})

	// still in 9. sample coin
	threshDsl.UponSignShareResult(m, func(sigShare []byte, context *signCoinShareCtx) error {
		dsl.SendMessage(m, mc.ThreshCrypto, CoinMessage(mc.Self, context.roundNumber, sigShare), params.AllNodes)

		return nil
	})

	// working in advance for 9. sample coin
	abbadsl.UponCoinMessageReceived(m, func(from t.NodeID, r uint64, coinShare []byte) error {
		if r != state.round.number || state.step > MAX_STEP {
			return nil // wrong round or already terminated
		}
		if _, present := state.round.coinRecvd[from]; present {
			return nil // duplicate message
		}
		state.round.coinRecvd[from] = struct{}{}

		context := &verifyCoinShareCtx{
			roundNumber: r,
			sigShare:    coinShare,
		}
		threshDsl.VerifyShare(m, mc.ThreshCrypto, state.round.coinData(params), coinShare, from, context)

		return nil
	})

	// working in advance for 9. sample coin
	threshDsl.UponVerifyShareResult(m, func(ok bool, err string, context *verifyCoinShareCtx) error {
		if context.roundNumber != state.round.number {
			return nil
		}

		if ok {
			state.round.coinRecvdOkShares = append(state.round.coinRecvdOkShares, context.sigShare)
		}

		return nil
	})

	dsl.UponCondition(m, func() error {
		// still in 9. sample coin
		if state.step == 9 && len(state.round.coinRecvdOkShares) > state.round.coinRecoverMinShareCount && !state.round.coinRecoverInProgress {
			context := &recoverCoinCtx{
				roundNumber: state.round.number,
			}
			threshDsl.Recover(m, mc.ThreshCrypto, state.round.coinData(params), state.round.coinRecvdOkShares, context)

			state.round.coinRecoverInProgress = true
			state.round.coinRecoverMinShareCount = len(state.round.coinRecvdOkShares)
		}

		return nil
	})

	// still in 9. sample coin
	threshDsl.UponRecoverResult(m, func(ok bool, fullSig []byte, err string, context *recoverCoinCtx) error {
		if context.roundNumber != state.round.number || state.step != 9 {
			// TODO: is this branch even possible?
			return nil // stale result
		}

		if ok {
			// we have a signature, all we need to do is hash it and
			dsl.HashOneMessage(m, mc.Hasher, [][]byte{fullSig}, context)
		} else {
			// will attempt to recover when more signature shares arrive
			state.round.coinRecoverInProgress = false
		}

		return nil
	})

	dsl.UponOneHashResult(m, func(hash []byte, context *recoverCoinCtx) error {
		if state.round.number != context.roundNumber || state.step != 9 {
			// TODO: is this branch even possible?
			return nil // stale result
		}

		// finishing step 9
		s_r := (hash[0] & 1) == 1 // TODO: this is ok, right?
		state.step = 10

		// 10.
		if state.round.values.Len() == 2 {
			state.round.estimate = s_r
		} else if state.round.values.Len() == 1 { // values = {v}
			v := state.round.values.Has(true) // if values contains true, v=true, otherwise v=false
			state.round.estimate = v

			// If in fact values = { s_r }, broadcast FINISH(s_r) if we haven't broadcast FINISH(_) already
			if v == s_r && !state.finishSent {
				dsl.SendMessage(m, mc.Net, FinishMessage(mc.Self, s_r), params.AllNodes)
				state.finishSent = true
			}
		}

		// (still 10.) Set r = r + 1, and return to step 4
		state.round.number += 1
		state.round.resetState(params)
		state.step = 4

		return nil
	})

	return m
}

// resets round state, apart from the round number and estimate
func (rs *abbaRoundState) resetState(params *ModuleParams) {
	rs.values = abbadsl.EmptyValueSet()

	rs.initRecvd = make(recvTracker, params.GetN())
	rs.initRecvdEstimates = makeBoolCounterMap()

	rs.auxRecvd = make(recvTracker, params.GetN())
	rs.auxRecvdValues = makeBoolCounterMap()

	rs.confRecvd = make(recvTracker, params.GetN())
	rs.confRecvdValues = makeValueSetCounterMap()

	rs.coinRecvd = make(recvTracker, params.strongSupportThresh())
	rs.coinRecvdOkShares = make([][]byte, 0, params.GetN()-params.GetF())

	rs.initSent = makeBoolBoolMap()
	rs.auxSent = false
	rs.coinRecoverMinShareCount = params.strongSupportThresh() - 1
}

func makeBoolCounterMap() map[bool]int {
	m := make(map[bool]int, 2)
	m[false] = 0
	m[true] = 0

	return m
}

func makeValueSetCounterMap() map[abbadsl.ValueSet]int {
	m := make(map[abbadsl.ValueSet]int, 4)

	for _, v := range []abbapb.ValueSet{abbapb.ValueSet_EMPTY, abbapb.ValueSet_ONE, abbapb.ValueSet_ZERO, abbapb.ValueSet_ZERO_AND_ONE} {
		m[abbadsl.ValueSet(v)] = 0
	}

	return m
}

func makeBoolBoolMap() map[bool]bool {
	m := make(map[bool]bool, 2)
	m[false] = false
	m[true] = false

	return m
}

func (params *ModuleParams) weakSupportThresh() int {
	return params.GetF() + 1
}

func (params *ModuleParams) strongSupportThresh() int {
	return params.GetN() - params.GetF()
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

const COIN_SIGN_DATA_PREFIX = "github.com/filecoin-project/mir/pkg/alea/aba"

func (rs *abbaRoundState) coinData(params *ModuleParams) [][]byte {
	return [][]byte{
		[]byte(COIN_SIGN_DATA_PREFIX),
		params.InstanceUID,
		toBytes(rs.number),
	}
}

func toBytes[T any](v T) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, binary.Size(v)))
	binary.Write(buf, binary.BigEndian, v)
	return buf.Bytes()
}
