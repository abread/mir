package abba

import (
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
	modringpbdsl "github.com/filecoin-project/mir/pkg/pb/modringpb/dsl"
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

const ModringSubName t.ModuleID = t.ModuleID("round")

func NewModule(mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	if tunables.MaxRoundLookahead <= 0 {
		return nil, fmt.Errorf("MaxRoundLookahead must be at least 1")
	}

	state := &state{
		phase:       phaseAwaitingInput,
		finishRecvd: make(abbat.RecvTracker, params.GetN()),
		finishSent:  false,
	}
	controller := newController(mc, params, tunables, nodeID, logger, state)

	slots := modring.New(&modring.ModuleConfig{Self: mc.Self.Then(ModringSubName)}, tunables.MaxRoundLookahead, modring.ModuleParams{
		Generator: newRoundGenerator(mc, params, nodeID, logger),
	}, logging.Decorate(logger, "Modring controller: "))

	return modules.RoutedModule(mc.Self, controller, wrapModuleStopAfterDelivered(slots, state)), nil
}

type roundAdvanceCtx struct {
	prevRoundNumber uint64
	nextEstimate    bool
}

func newController(mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger, state *state) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

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
			logger.Log(logging.LevelDebug, "duplicate FINISH(v)", "v", value)
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

			// only care about finish messages from now on
			// eventually instances that are out-of-date will receive them and be happy
			rnetdsl.MarkModuleMsgsRecvd(m, mc.ReliableNet, mc.Self.Then(ModringSubName), params.AllNodes)
		}

		return nil
	})

	abbadsl.UponInputValue(m, func(input bool, origin *abbapbtypes.Origin) error {
		if state.phase != phaseAwaitingInput {
			return fmt.Errorf("input value provided twice to abba")
		}

		state.origin = origin

		roundNum := uint64(0)
		abbadsl.RoundInputValue(m, mc.RoundID(0), input, &roundNum)

		state.phase = phaseRunning

		return nil
	})

	abbadsl.UponRoundDeliver(m, func(nextEstimate bool, prevRoundNum *uint64) error {
		if state.phase != phaseRunning {
			return nil // already done
		}
		if *prevRoundNum != state.currentRound {
			return fmt.Errorf("round number mismatch: expected %v, but round %v delivered", state.currentRound, *prevRoundNum)
		}

		// clean up old round
		context := &roundAdvanceCtx{
			prevRoundNumber: state.currentRound,
			nextEstimate:    nextEstimate,
		}
		modringpbdsl.FreeSubmodule(m, mc.Self.Then(ModringSubName), state.currentRound, context)

		return nil
	})
	modringpbdsl.UponFreedSubmodule(m, func(context *roundAdvanceCtx) error {
		if context.prevRoundNumber != state.currentRound {
			return fmt.Errorf("round number mismatch: expected %v, but round %v was freed", state.currentRound, context.prevRoundNumber)
		}

		state.currentRound++
		roundNum := state.currentRound // copy of round number, to make our previous sanity check work
		abbadsl.RoundInputValue(m, mc.RoundID(state.currentRound), context.nextEstimate, &roundNum)

		return nil
	})

	return m
}

func newRoundGenerator(controllerMc *ModuleConfig, controllerParams *ModuleParams, nodeID t.NodeID, logger logging.Logger) func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
	params := &abbaround.ModuleParams{
		InstanceUID: controllerParams.InstanceUID, // TODO: review
		AllNodes:    controllerParams.AllNodes,
	}

	return func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
		mc := &abbaround.ModuleConfig{
			Self:         id,
			Consumer:     controllerMc.Self,
			ReliableNet:  controllerMc.ReliableNet,
			ThreshCrypto: controllerMc.ThreshCrypto,
			Hasher:       controllerMc.Hasher,
		}

		mod := abbaround.New(mc, params, nodeID, logging.Decorate(logger, "Abba Round: ", "abbaRound", idx))

		initialEvs := &events.EventList{}
		return mod, initialEvs, nil
	}
}

type stopsAfterDelivered struct {
	inner modules.PassiveModule
	state *state
}

func (m *stopsAfterDelivered) ImplementsModule() {}
func (m *stopsAfterDelivered) ApplyEvents(evs *events.EventList) (*events.EventList, error) {
	if m.state.phase == phaseDelivered {
		return &events.EventList{}, nil
	}

	return m.inner.ApplyEvents(evs)
}

func wrapModuleStopAfterDelivered(module modules.PassiveModule, state *state) modules.PassiveModule {
	return &stopsAfterDelivered{
		inner: module,
		state: state,
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
