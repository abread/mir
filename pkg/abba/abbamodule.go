package abba

import (
	"context"
	"fmt"

	"github.com/filecoin-project/mir/pkg/abba/abbaround"
	abbat "github.com/filecoin-project/mir/pkg/abba/abbatypes"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modring"
	"github.com/filecoin-project/mir/pkg/modules"
	abbadsl "github.com/filecoin-project/mir/pkg/pb/abbapb/dsl"
	abbapbmsgs "github.com/filecoin-project/mir/pkg/pb/abbapb/msgs"
	abbapbtypes "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	rnetdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	ReliableNet  t.ModuleID
	ThreshCrypto t.ModuleID
	Hasher       t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         "abba",
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
)

type state struct {
	phase        abbaPhase
	currentRound uint64

	finishRecvd            abbat.RecvTracker
	finishRecvdValueCounts abbat.BoolCounters

	finishSent bool

	origin *abbapbtypes.Origin
}

const modringSubName t.ModuleID = t.ModuleID("r")

func NewModule(ctx context.Context, mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger, outChan chan *events.EventList) (modules.ActiveModule, error) {
	if tunables.MaxRoundLookahead <= 0 {
		return nil, fmt.Errorf("MaxRoundLookahead must be at least 1")
	}

	// rounds use a submodule namespace to allow us to mark all round messages as received at once
	rounds := modring.New(ctx, mc.Self.Then(modringSubName), tunables.MaxRoundLookahead, &modring.ModuleParams{
		Generator:      newRoundGenerator(mc, params, nodeID, logger),
		InputQueueSize: 8, // TODO: make configurable
	}, logging.Decorate(logger, "Modring controller: "), outChan)

	controller := newController(ctx, mc, params, tunables, nodeID, logger, rounds)

	return modules.RoutedModule(ctx, mc.Self, controller, rounds, outChan), nil
}

func newController(ctx context.Context, mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger, rounds *modring.Module) modules.PassiveModule {
	m := dsl.NewModule(ctx, mc.Self)

	state := &state{
		phase:       phaseAwaitingInput,
		finishRecvd: make(abbat.RecvTracker, params.GetN()),
		finishSent:  false,
	}

	abbadsl.UponRoundFinishAll(m, func(decision bool) error {
		if !state.finishSent && state.phase == phaseRunning {
			rnetdsl.SendMessage(m, mc.ReliableNet, FinishMsgID(), abbapbmsgs.FinishMessage(mc.Self, decision), params.AllNodes)
			state.finishSent = true
		}

		return nil
	})

	abbadsl.UponFinishMessageReceived(m, func(from t.NodeID, value bool) error {
		rnetdsl.Ack(m, mc.ReliableNet, mc.Self, FinishMsgID(), from)

		if !state.finishRecvd.Register(from) {
			logger.Log(logging.LevelWarn, "duplicate FINISH(v)", "v", value)
			return nil // duplicate message
		}

		state.finishRecvdValueCounts.Increment(value)

		if state.phase != phaseRunning {
			return nil
		}

		// 1. upon receiving weak support for FINISH(v), broadcast FINISH(v)
		if !state.finishSent && state.finishRecvdValueCounts.Get(value) >= params.weakSupportThresh() {
			logger.Log(logging.LevelDebug, "received weak support for FINISH(v)", "v", value)
			rnetdsl.SendMessage(m, mc.ReliableNet,
				FinishMsgID(),
				abbapbmsgs.FinishMessage(mc.Self, value),
				params.AllNodes,
			)
			state.finishSent = true
		}

		// 2. upon receiving strong support for FINISH(v), output v and terminate
		if state.origin != nil && state.finishRecvdValueCounts.Get(value) >= params.strongSupportThresh() {
			logger.Log(logging.LevelDebug, "received strong support for FINISH(v)", "v", value)

			abbadsl.Deliver(m, state.origin.Module, value, state.origin)
			state.phase = phaseDelivered

			if err := rounds.MarkAllPast(); err != nil {
				return fmt.Errorf("could not stop abba rounds: %w", err)
			}

			// only care about finish messages from now on
			// eventually instances that are out-of-date will receive them and be happy
			rnetdsl.MarkModuleMsgsRecvd(m, mc.ReliableNet, mc.Self.Then(modringSubName), params.AllNodes)
		}

		return nil
	})

	abbadsl.UponInputValue(m, func(input bool, origin *abbapbtypes.Origin) error {
		if state.phase != phaseAwaitingInput {
			return fmt.Errorf("input value provided twice to abba")
		}

		state.origin = origin

		abbadsl.RoundInputValue(m, mc.RoundID(0), input)

		state.phase = phaseRunning

		return nil
	})

	abbadsl.UponRoundDeliver(m, func(nextEstimate bool, prevRoundNum uint64) error {
		if state.phase != phaseRunning {
			return nil // already done
		}
		if prevRoundNum != state.currentRound {
			return fmt.Errorf("round number mismatch: expected %v, but round %v delivered", state.currentRound, prevRoundNum)
		}

		// clean up old round
		if err := rounds.MarkSubmodulePast(state.currentRound); err != nil {
			return fmt.Errorf("failed to clean up previous round: %w", err)
		}

		state.currentRound++
		abbadsl.RoundInputValue(m, mc.RoundID(state.currentRound), nextEstimate)

		return nil
	})

	return m
}

func newRoundGenerator(controllerMc *ModuleConfig, controllerParams *ModuleParams, nodeID t.NodeID, logger logging.Logger) func(ctx context.Context, id t.ModuleID, idx uint64, outChan chan *events.EventList) (modules.Module, *events.EventList, error) {

	return func(ctx context.Context, id t.ModuleID, idx uint64, outChan chan *events.EventList) (modules.Module, *events.EventList, error) {
		mc := &abbaround.ModuleConfig{
			Self:         id,
			Consumer:     controllerMc.Self,
			ReliableNet:  controllerMc.ReliableNet,
			ThreshCrypto: controllerMc.ThreshCrypto,
			Hasher:       controllerMc.Hasher,
		}

		params := &abbaround.ModuleParams{
			InstanceUID: controllerParams.InstanceUID, // TODO: review
			AllNodes:    controllerParams.AllNodes,
			RoundNumber: idx,
		}

		mod := abbaround.New(ctx, mc, params, nodeID, logging.Decorate(logger, "Abba Round: ", "abbaRound", idx))

		initialEvs := &events.EventList{}
		return mod, initialEvs, nil
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
	return mc.Self.Then(modringSubName).Then(t.NewModuleIDFromInt(roundNumber))
}
