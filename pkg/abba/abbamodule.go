package abba

import (
	"time"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/abba/abbaround"
	abbat "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modring"
	"github.com/filecoin-project/mir/pkg/modules"
	abbadsl "github.com/filecoin-project/mir/pkg/pb/abbapb/dsl"
	abbapbmsgs "github.com/filecoin-project/mir/pkg/pb/abbapb/msgs"
	rnetdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	reliablenetpbevents "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/events"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID
	ReliableNet  t.ModuleID
	ThreshCrypto t.ModuleID
	Hasher       t.ModuleID
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // unique identifier for this instance of VCB
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
}

// ModuleTunables sets the values of protocol tunables that need not be the same across all nodes.
type ModuleTunables struct {
	// Maximum number of ABBA rounds for which we process messages
	// Must be at least 1
	MaxRoundLookahead int
}

type abbaPhase uint8

const (
	phaseAwaitingInput abbaPhase = iota
	phaseRunning
	phaseDelivered
	phaseDone
)

type state struct {
	phase        abbaPhase
	currentRound uint64

	finishRecvd            abbat.RecvTracker
	finishRecvdValueCounts abbat.BoolCounters

	finishSent bool
}

const ModringSubName t.ModuleID = t.ModuleID("r")

func NewModule(mc ModuleConfig, params ModuleParams, tunables ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	if tunables.MaxRoundLookahead <= 0 {
		return nil, es.Errorf("MaxRoundLookahead must be at least 1")
	}

	// rounds use a submodule namespace to allow us to mark all round messages as received at once
	rounds := modring.New(mc.Self.Then(ModringSubName), tunables.MaxRoundLookahead, modring.ModuleParams{
		Generator: newRoundGenerator(mc, params, nodeID, logger),
		CleanupHandler: func(roundID uint64) (events.EventList, error) {
			return events.ListOf(
				reliablenetpbevents.MarkModuleMsgsRecvd(mc.ReliableNet, mc.Self.Then(ModringSubName).Then(t.NewModuleIDFromInt(roundID)), params.AllNodes),
			), nil
		},
	}, logging.Decorate(logger, "Modring controller: "))

	controller := newController(mc, params, logger, rounds)

	return modules.RoutedModule(mc.Self, controller, rounds), nil
}

func newController(mc ModuleConfig, params ModuleParams, logger logging.Logger, rounds *modring.Module) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)
	_ = logger // silence warnings about unused variables

	state := &state{
		phase:       phaseAwaitingInput,
		finishRecvd: make(abbat.RecvTracker, params.GetN()),
		finishSent:  false,
	}

	abbadsl.UponContinueExecution(m, func() error {
		if state.phase == phaseAwaitingInput {
			return es.Errorf("abba ordered to make progress without input")
		} else if state.phase == phaseDone {
			// nothing to do
			return nil
		}

		// tell the first round to make progress
		abbadsl.RoundContinue(m, mc.RoundID(0))
		return nil
	})

	abbadsl.UponRoundFinishAll(m, func(decision bool, unanimous bool) error {
		if state.phase < phaseRunning {
			return es.Errorf("abba round cannot deliver before input")
		} else if state.phase > phaseRunning {
			return nil // no need
		}

		if !state.finishSent {
			rnetdsl.SendMessage(m, mc.ReliableNet, FinishMsgID(), abbapbmsgs.FinishMessage(mc.Self, decision), params.AllNodes)
			state.finishSent = true
		}

		if unanimous {
			logger.Log(logging.LevelDebug, "delivering with unanimous fast path", "decision", decision)
			abbadsl.Deliver(m, mc.Consumer, decision, mc.Self)
			state.phase = phaseDelivered
		}

		return nil
	})

	abbadsl.UponFinishMessageReceived(m, func(from t.NodeID, value bool) error {
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, FinishMsgID(), from)

		if !state.finishRecvd.Register(from) {
			// logger.Log(logging.LevelWarn, "duplicate FINISH(v)", "v", value)
			return nil // duplicate message
		}

		state.finishRecvdValueCounts.Increment(value)

		return nil
	})

	dsl.UponStateUpdates(m, func() error {
		if state.phase == phaseDone {
			return nil
		}

		for _, value := range []bool{true, false} {
			// 1. upon receiving weak support for FINISH(v), broadcast FINISH(v)
			if !state.finishSent && state.finishRecvdValueCounts.Get(value) >= params.weakSupportThresh() {
				// logger.Log(logging.LevelDebug, "received weak support for FINISH(v)", "v", value)
				rnetdsl.SendMessage(m, mc.ReliableNet,
					FinishMsgID(),
					abbapbmsgs.FinishMessage(mc.Self, value),
					params.AllNodes,
				)
				state.finishSent = true
			}

			// 2. upon receiving strong support for FINISH(v), output v and terminate
			if state.finishRecvdValueCounts.Get(value) >= params.strongSupportThresh() {
				logger.Log(logging.LevelDebug, "received strong support for FINISH(v)", "v", value)

				if state.phase != phaseDelivered {
					abbadsl.Deliver(m, mc.Consumer, value, mc.Self)
				}
				abbadsl.Done(m, mc.Consumer, mc.Self)
				state.phase = phaseDone

				// only care about finish messages from now on
				// eventually instances that are out-of-date will receive them and be happy
				evsOut, err := rounds.MarkPastAndFreeAll()
				if err != nil {
					return err
				}
				dsl.EmitEvents(m, evsOut)
			}
		}
		return nil
	})

	abbadsl.UponInputValue(m, func(input bool) error {
		if state.phase == phaseDone {
			// this may be a duplicate input, but we're ignoring that sanity check.
			return nil
		} else if state.phase != phaseAwaitingInput {
			return es.Errorf("input value provided twice to abba")
		}

		abbadsl.RoundInputValue(m, mc.RoundID(0), input)
		state.phase = phaseRunning
		return nil
	})

	abbadsl.UponRoundDeliver(m, func(nextEstimate bool, prevRoundNum uint64, _ time.Duration) error {
		if state.phase == phaseAwaitingInput {
			return es.Errorf("round delivered before input")
		} else if state.phase == phaseDone {
			return nil // already done
		}
		if prevRoundNum != state.currentRound {
			return es.Errorf("round number mismatch: expected %v, but round %v delivered", state.currentRound, prevRoundNum)
		}

		// clean up old round
		evsOut, err := rounds.MarkSubmodulePast(state.currentRound)
		if err != nil {
			return es.Errorf("failed to clean up previous round: %w", err)
		}
		dsl.EmitEvents(m, evsOut)

		state.currentRound++
		abbadsl.RoundInputValue(m, mc.RoundID(state.currentRound), nextEstimate)

		return nil
	})

	return m
}

func newRoundGenerator(controllerMc ModuleConfig, controllerParams ModuleParams, nodeID t.NodeID, logger logging.Logger) func(id t.ModuleID, idx uint64) (modules.PassiveModule, events.EventList, error) {

	return func(id t.ModuleID, idx uint64) (modules.PassiveModule, events.EventList, error) {
		mc := abbaround.ModuleConfig{
			Self:         id,
			Consumer:     controllerMc.Self,
			ReliableNet:  controllerMc.ReliableNet,
			ThreshCrypto: controllerMc.ThreshCrypto,
			Hasher:       controllerMc.Hasher,
		}

		params := abbaround.ModuleParams{
			InstanceUID: controllerParams.InstanceUID, // TODO: review
			AllNodes:    controllerParams.AllNodes,
			RoundNumber: idx,
		}

		mod := abbaround.New(mc, params, nodeID, logging.Decorate(logger, "Abba Round: ", "abbaRound", idx))

		return mod, events.EmptyList(), nil
	}
}

// GetN returns the total number of nodes.
func (params *ModuleParams) GetN() int {
	return len(params.AllNodes)
}

// GetF returns the maximum tolerated number of faulty nodes.
func (params *ModuleParams) GetF() int {
	return (params.GetN() - 1) / 3
}

func (params *ModuleParams) weakSupportThresh() int {
	return params.GetF() + 1
}

func (params *ModuleParams) strongSupportThresh() int {
	return params.GetN() - params.GetF()
}

func (mc *ModuleConfig) RoundID(roundNumber uint64) t.ModuleID {
	return mc.Self.Then(ModringSubName).Then(t.NewModuleIDFromInt(roundNumber))
}
