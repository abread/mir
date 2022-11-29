package abba

import (
	"fmt"
	"math"

	"github.com/filecoin-project/mir/pkg/abba/abbadsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/factorymodule"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/abbapb"
	"github.com/filecoin-project/mir/pkg/pb/factorymodulepb"
	"github.com/filecoin-project/mir/pkg/reliablenet/rnetdsl"
	"github.com/filecoin-project/mir/pkg/serializing"
	threshDsl "github.com/filecoin-project/mir/pkg/threshcrypto/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID // id of the module to send the "Deliver" event to
	ReliableNet  t.ModuleID
	ThreshCrypto t.ModuleID
	Hasher       t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         "abba",
		Consumer:     consumer,
		ReliableNet:  "reliablenet",
		ThreshCrypto: "threshcrypto",
		Hasher:       "hasher",
	}
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

	initRecvd          map[bool]recvTracker
	initRecvdEstimates map[bool]int

	auxRecvd       recvTracker
	auxRecvdValues map[bool]int

	confRecvd       recvTracker
	confRecvdValues map[abbadsl.ValueSet]int

	coinRecvd         recvTracker
	coinRecvdOkShares [][]byte

	initWeakSupportReached   map[bool]bool
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

const MaxStep uint8 = 10

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
					logger,
				)
				return inst, nil
			},
		),
		logger,
	)
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, logger logging.Logger) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)
	state := &abbaModuleState{
		step: 0,

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
			logger.Log(logging.LevelDebug, "duplicate FINISH(v)", "v", value)
			return nil // duplicate message
		}
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, FinishMsgID(), from)
		state.finishRecvd[from] = struct{}{}
		state.finishRecvdValues[value]++

		// 1. upon receiving weak support for FINISH(v), broadcast FINISH(v)
		if !state.finishSent && state.finishRecvdValues[value] >= params.weakSupportThresh() {
			logger.Log(logging.LevelDebug, "received weak support for FINISH(v)", "v", value)
			rnetdsl.SendMessage(m, mc.ReliableNet,
				FinishMsgID(),
				FinishMessage(mc.Self, value),
				params.AllNodes,
			)
			state.finishSent = true
		}

		// 2. upon receiving strong support for FINISH(v), output v and terminate
		if state.step <= MaxStep && state.finishRecvdValues[value] >= params.strongSupportThresh() {
			logger.Log(logging.LevelDebug, "received strong support for FINISH(v)", "v", value)
			abbadsl.Deliver(m, mc.Consumer, value)
			state.updateStep(math.MaxUint8) // no more progress can be made

			// only care about finish messages from now on
			// eventually instances that are out-of-date will receive them and be happy
			// TODO: do this more efficiently (by splitting the rounds into another module maybe)
			for i := uint64(0); i <= state.round.number; i++ {
				for _, msgID := range [][]byte{
					InitMsgID(i, false),
					InitMsgID(i, true),
					AuxMsgID(i),
					ConfMsgID(i),
					CoinMsgID(i),
				} {
					rnetdsl.MarkRecvd(m, mc.ReliableNet, mc.Self, msgID, params.AllNodes)
				}
			}
		}

		return nil
	})

	state.updateStep(3)

	// 3. upon P_i providing input value v_in, set est^r_i=v_in, r=0
	abbadsl.UponInputValue(m, func(input bool) error {
		if state.step > 3 {
			logger.Log(logging.LevelDebug, "input already provided to this ABBA instance (or it terminated before we could have input a value to it)")
			return nil
		}

		state.round.number = 0
		state.round.estimate = input

		// 4. broadcast INIT(r, est^r_i)
		msg := InitMessage(mc.Self, state.round.number, state.round.estimate)
		rnetdsl.SendMessage(m, mc.ReliableNet, InitMsgID(state.round.number, state.round.estimate), msg, params.AllNodes)

		state.updateStep(5)
		return nil
	})

	registerRoundEvents(m, state, mc, params, nodeID, logger)

	return m
}

func registerRoundEvents(m dsl.Module, state *abbaModuleState, mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, logger logging.Logger) { // nolint: gocognit, gocyclo
	// TODO: isolate coin sampling to different module to reduce complexity/noise

	abbadsl.UponInitMessageReceived(m, func(from t.NodeID, r uint64, est bool) error {
		if r != state.round.number {
			logger.Log(logging.LevelDebug, "wrong round for INIT(r, est)", "current", state.round.number, "got", r)
			return nil // not ready yet or wrong round
		}
		if _, present := state.round.initRecvd[est][from]; present {
			logger.Log(logging.LevelDebug, "duplicate INIT(r, _)", "r", r)
			return nil // duplicate message
		}
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, InitMsgID(r, est), from)
		state.round.initRecvd[est][from] = struct{}{}
		state.round.initRecvdEstimates[est]++

		return nil
	})

	dsl.UponCondition(m, func() error {
		if state.step < 5 || state.step > MaxStep {
			return nil // did not perform the initial INIT broadcast or has already terminated
		}

		r := state.round.number
		for _, est := range []bool{false, true} {
			// 5. upon receiving weak support for INIT(r, v), add v to values and broadcast INIT(r, v)
			if !state.round.initWeakSupportReached[est] && state.round.initRecvdEstimates[est] >= params.weakSupportThresh() {
				logger.Log(logging.LevelDebug, "received weak support for INIT(r, v)", "r", r, "v", est)

				state.round.values.Add(est)

				// if we're in round 0, and our round.estimate is r, then we already sent this message
				if !(state.round.number == 0 && state.round.estimate == est) {
					rnetdsl.SendMessage(m, mc.ReliableNet, InitMsgID(r, est), InitMessage(mc.Self, r, est), params.AllNodes)
				}

				state.round.initWeakSupportReached[est] = true

				state.updateStep(6)
			}

			// 6. upon receiving strong support for INIT(r, v), broadcast AUX(r, v) if we have not already broadcast AUX(r, _)
			if !state.round.auxSent && state.round.initRecvdEstimates[est] >= params.strongSupportThresh() {
				logger.Log(logging.LevelDebug, "received strong support for INIT(r, v)", "r", r, "v", est)
				rnetdsl.SendMessage(m, mc.ReliableNet, AuxMsgID(r), AuxMessage(mc.Self, r, est), params.AllNodes)
				state.round.auxSent = true
				state.updateStep(7)
			}
		}

		return nil
	})

	abbadsl.UponAuxMessageReceived(m, func(from t.NodeID, r uint64, value bool) error {
		if r != state.round.number {
			logger.Log(logging.LevelDebug, "wrong round for AUX(r, v)", "current", state.round.number, "got", r)
			return nil // not processing this round
		}
		if _, present := state.round.auxRecvd[from]; present {
			logger.Log(logging.LevelDebug, "duplicate AUX(r, _)", "r", r)
			return nil // duplicate message
		}
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, AuxMsgID(r), from)
		state.round.auxRecvd[from] = struct{}{}
		state.round.auxRecvdValues[value]++

		logger.Log(logging.LevelDebug, "recvd AUX(r, v)", "r", r, "v", value)

		return nil
	})

	// 7. wait until there exists a subset of nodes with size >= q_S(= N-F), from which we have received AUX(v', r) with any v' in round.values, then broadcast CONF(values, r)
	dsl.UponCondition(m, func() error {
		if state.step != 7 {
			return nil
		}

		if state.round.isNiceAuxValueCount(params) {
			r := state.round.number
			logger.Log(logging.LevelDebug, "received enough support for AUX(r, v in values)", "r", r, "values", state.round.values)
			rnetdsl.SendMessage(m, mc.ReliableNet, ConfMsgID(r), ConfMessage(mc.Self, r, state.round.values), params.AllNodes)
			state.updateStep(8)
		}

		return nil
	})

	abbadsl.UponConfMessageReceived(m, func(from t.NodeID, r uint64, values abbadsl.ValueSet) error {
		if r != state.round.number {
			logger.Log(logging.LevelDebug, "wrong round for CONF(r, c)", "current", state.round.number, "got", r)
			return nil // wrong round
		}
		if _, present := state.round.confRecvd[from]; present {
			logger.Log(logging.LevelDebug, "duplicate CONF(r, _)", "r", r)
			return nil // duplicate message
		}
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, ConfMsgID(r), from)
		state.round.confRecvd[from] = struct{}{}
		state.round.confRecvdValues[values]++

		return nil
	})

	// 8. wait until there exists a subset of nodes with size >= q_S(= N-F), from which we have received CONF(vs', r) with any vs' subset of round.values
	dsl.UponCondition(m, func() error {
		if state.step != 8 {
			return nil
		}

		if state.round.isNiceConfValuesCount(params) {
			r := state.round.number
			logger.Log(logging.LevelDebug, "received enough support for CONF(r, C subset of values)", "r", r, "values", state.round.values)

			// 9. sample coin
			state.updateStep(9)
			threshDsl.SignShare(m, mc.ThreshCrypto, state.round.coinData(params), &signCoinShareCtx{
				roundNumber: r,
			})
		}

		return nil
	})

	// still in 9. sample coin
	threshDsl.UponSignShareResult(m, func(sigShare []byte, context *signCoinShareCtx) error {
		rnetdsl.SendMessage(m, mc.ReliableNet, CoinMsgID(context.roundNumber), CoinMessage(mc.Self, context.roundNumber, sigShare), params.AllNodes)

		return nil
	})

	// working in advance for 9. sample coin
	abbadsl.UponCoinMessageReceived(m, func(from t.NodeID, r uint64, coinShare []byte) error {
		if r != state.round.number {
			logger.Log(logging.LevelDebug, "wrong round for COIN(r, s)", "current", state.round.number, "got", r)
			return nil // wrong round or already terminated
		}
		if state.step > MaxStep {
			logger.Log(logging.LevelDebug, "already terminated (recvd COIN)")
			return nil // already terminated
		}
		if _, present := state.round.coinRecvd[from]; present {
			logger.Log(logging.LevelDebug, "duplicate COIN(r, _)", "r", r)
			return nil // duplicate message
		}

		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, CoinMsgID(r), from)
		logger.Log(logging.LevelDebug, "recvd COIN(r, share)", "r", r)
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
			logger.Log(logging.LevelDebug, "wrong round for verifyshares", "current", state.round.number, "got", context.roundNumber)
			return nil
		}

		if ok {
			state.round.coinRecvdOkShares = append(state.round.coinRecvdOkShares, context.sigShare)
		}

		return nil
	})

	// still in 9. sample coin
	dsl.UponCondition(m, func() error {
		if state.step != 9 {
			return nil
		}

		if len(state.round.coinRecvdOkShares) > state.round.coinRecoverMinShareCount && !state.round.coinRecoverInProgress {
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
		if context.roundNumber != state.round.number {
			return fmt.Errorf("impossible condition: changed round without coin toss")
		}
		if state.step < 9 {
			return fmt.Errorf("impossible condition: changed round without coin toss")
		} else if state.step > 9 {
			return nil // stale result
		}

		if ok {
			// we have a signature, all we need to do is hash it and
			dsl.HashOneMessage(m, mc.Hasher, [][]byte{fullSig}, context)
		} else {
			logger.Log(logging.LevelDebug, "could not recover coin ...YET")
			// will attempt to recover when more signature shares arrive
			state.round.coinRecoverInProgress = false
		}

		return nil
	})

	dsl.UponOneHashResult(m, func(hash []byte, context *recoverCoinCtx) error {
		if state.step < 9 {
			return fmt.Errorf("impossible condition: changed round without coin toss")
		} else if state.step > 9 {
			return nil // stale result
		}

		if state.round.number != context.roundNumber {
			return fmt.Errorf("impossible condition: changed round without coin toss")
		}

		// finishing step 9
		sR := (hash[0] & 1) == 1 // TODO: this is ok, right?
		state.updateStep(10)

		// 10.
		if state.round.values.Len() == 2 {
			state.round.estimate = sR
		} else if state.round.values.Len() == 1 { // values = {v}
			v := state.round.values.Has(true) // if values contains true, v=true, otherwise v=false
			state.round.estimate = v

			// If in fact values = { s_r }, broadcast FINISH(s_r) if we haven't broadcast FINISH(_) already
			if v == sR && !state.finishSent {
				rnetdsl.SendMessage(m, mc.ReliableNet, FinishMsgID(), FinishMessage(mc.Self, sR), params.AllNodes)
				state.finishSent = true
			}
		}

		// (still 10.) Set r = r + 1, and return to step 4
		state.round.number++
		state.round.resetState(params)
		state.step = 4 // can't use updateStep, must go backwards

		// 4. broadcast INIT(r, est)
		rnetdsl.SendMessage(m, mc.ReliableNet, InitMsgID(state.round.number, state.round.estimate), InitMessage(mc.Self, state.round.number, state.round.estimate), params.AllNodes)
		state.updateStep(5)

		return nil
	})
}

// resets round state, apart from the round number and estimate
func (rs *abbaRoundState) resetState(params *ModuleParams) {
	rs.values = abbadsl.EmptyValueSet()

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
	rs.coinRecvdOkShares = make([][]byte, 0, params.GetN()-params.GetF())

	rs.initWeakSupportReached = makeBoolBoolMap()
	rs.auxSent = false
	rs.coinRecoverInProgress = false
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

func (s *abbaModuleState) updateStep(newStep uint8) {
	if newStep > s.step {
		s.step = newStep
	}
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

const (
	MsgTypeFinish uint8 = iota
	MsgTypeInit
	MsgTypeAux
	MsgTypeConf
	MsgTypeCoin
)

func FinishMsgID() []byte {
	return []byte{MsgTypeFinish}
}

func InitMsgID(r uint64, v bool) []byte {
	s := make([]byte, 0, 1+8+1)
	s = append(s, MsgTypeInit)
	s = append(s, serializing.Uint64ToBytes(r)...)
	s = append(s, boolToNum(v))
	return s
}

func AuxMsgID(r uint64) []byte {
	s := make([]byte, 0, 1+8)
	s = append(s, MsgTypeAux)
	s = append(s, serializing.Uint64ToBytes(r)...)
	return s
}

func ConfMsgID(r uint64) []byte {
	s := make([]byte, 0, 1+8)
	s = append(s, MsgTypeConf)
	s = append(s, serializing.Uint64ToBytes(r)...)
	return s
}

func CoinMsgID(r uint64) []byte {
	s := make([]byte, 0, 1+8)
	s = append(s, MsgTypeCoin)
	s = append(s, serializing.Uint64ToBytes(r)...)
	return s
}

func boolToNum(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}
