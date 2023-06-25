package agreement

import (
	"math"
	"strconv"
	"time"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/abba"
	"github.com/filecoin-project/mir/pkg/abba/abbatypes"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modring"
	"github.com/filecoin-project/mir/pkg/modules"
	abbapbdsl "github.com/filecoin-project/mir/pkg/pb/abbapb/dsl"
	abbapbevents "github.com/filecoin-project/mir/pkg/pb/abbapb/events"
	abbapbmsgs "github.com/filecoin-project/mir/pkg/pb/abbapb/msgs"
	abbapbtypes "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	agreementpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/dsl"
	agreementpbevents "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/events"
	agreementpbmsgdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/dsl"
	agreementpbmsgs "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/msgs"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	reliablenetpbdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	transportpbdsl "github.com/filecoin-project/mir/pkg/pb/transportpb/dsl"
	transportpbevents "github.com/filecoin-project/mir/pkg/pb/transportpb/events"
	t "github.com/filecoin-project/mir/pkg/types"
)

var timeRef = time.Now()

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID
	Hasher       t.ModuleID
	ReliableNet  t.ModuleID
	Net          t.ModuleID
	ThreshCrypto t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self:         "alea_ag",
		Consumer:     "alea_dir",
		Hasher:       "hasher",
		ReliableNet:  "reliablenet",
		Net:          "net",
		ThreshCrypto: "threshcrypto",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // must be the alea ID followed by 'a'
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
}

// ModuleTunables sets the values of protocol tunables that need not be the same across all nodes.
type ModuleTunables struct {
	// Maximum number of agreement rounds for which we process messages
	// Must be at least 1
	MaxRoundLookahead int

	// Maximum number of ABBA rounds for which we process messages
	// Must be at least 1
	MaxAbbaRoundLookahead int
}

func NewModule(mc ModuleConfig, params ModuleParams, tunables ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	if tunables.MaxRoundLookahead <= 0 {
		return nil, es.Errorf("MaxRoundLookahead must be at least 1")
	} else if tunables.MaxAbbaRoundLookahead <= 0 {
		return nil, es.Errorf("MaxAbbaRoundLookahead must be at least 1")
	}

	agRounds := modring.New(
		mc.Self,
		tunables.MaxRoundLookahead,
		modring.ModuleParams{
			Generator:      newAbbaGenerator(mc, params, tunables, nodeID, logger),
			PastMsgHandler: newPastMessageHandler(mc),
		},
		logging.Decorate(logger, "Modring controller: "),
	)

	N := len(params.AllNodes)
	F := (N - 1) / 3

	state := &state{
		rounds: make(map[uint64]*round, tunables.MaxRoundLookahead),

		// stalled inputs count can never go beyond the number of concurrent agreement rounds
		pendingInput: make(map[uint64]pendingInput, tunables.MaxRoundLookahead),

		N:         N,
		strongMaj: 2*F + 1,
	}

	controller := newAgController(mc, logger, state, agRounds)

	perfSniffer := newAgRoundSniffer(mc, params, agRounds, state)

	agRounds2 := modules.MultiApplyModule([]modules.PassiveModule{
		perfSniffer,
		agRounds,
	})

	return modules.RoutedModule(mc.Self, controller, agRounds2), nil
}

func (mc *ModuleConfig) agRoundModuleID(id uint64) t.ModuleID {
	return mc.Self.Then(t.NewModuleIDFromInt(id))
}

type round struct {
	nodesWithInput         abbatypes.RecvTracker // TODO: move recvtracker to another package?
	nodesWithPosInputCount int
	relInputTime           time.Duration
	relPosQuorumTime       time.Duration // instant where 2F+1 nodes provided input=1
	relPosTotalTime        time.Duration // instant where all nodes provided input=1

	abbaRoundNumber       uint64
	relAbbaRoundStartTime time.Duration

	delivered         bool
	decision          bool
	posQuorumDuration time.Duration
	posTotalDuration  time.Duration
}

type pendingInput struct {
	input   bool
	relTime time.Duration
}

type state struct {
	roundDecisionHistory AgRoundHistory

	currentRound uint64

	rounds map[uint64]*round

	pendingInput map[uint64]pendingInput

	N         int
	strongMaj int
}

func (s *state) storeRoundDecision(rounds *modring.Module, roundNum uint64, decision bool) error {
	if !rounds.IsInView(roundNum) {
		return es.Errorf("round out of view: %v", rounds)
	}

	r := s.rounds[roundNum]

	r.decision = decision
	r.delivered = true

	if r.relInputTime == 0 && r.relPosQuorumTime != 0 {
		r.posQuorumDuration = 0
	} else if r.relInputTime != 0 && r.relPosQuorumTime == 0 {
		r.posQuorumDuration = time.Duration(math.MaxInt64)
	} else {
		r.posQuorumDuration = r.relPosQuorumTime - r.relInputTime
	}

	if r.relPosTotalTime != 0 {
		r.posTotalDuration = r.relPosTotalTime - r.relPosQuorumTime
	} else {
		r.posTotalDuration = time.Duration(math.MaxInt64)
	}

	return nil
}

func (s *state) ensureRoundInitialized(roundNum uint64) *round {
	if _, ok := s.rounds[roundNum]; !ok {
		s.rounds[roundNum] = &round{
			nodesWithInput: make(abbatypes.RecvTracker, s.N),
		}
	}

	return s.rounds[roundNum]
}

func (s *state) loadRound(rounds *modring.Module, roundNum uint64) *round {
	if !rounds.IsInView(roundNum) {
		return nil
	}

	if _, ok := s.rounds[roundNum]; !ok {
		return nil
	}

	return s.rounds[roundNum]
}

func (s *state) clearRoundData(rounds *modring.Module, roundNum uint64) error {
	if !rounds.IsInView(roundNum) {
		return es.Errorf("round out of view: %v", rounds)
	}

	delete(s.rounds, roundNum)
	return nil
}

func newAgController(mc ModuleConfig, logger logging.Logger, state *state, agRounds *modring.Module) modules.PassiveModule {
	_ = logger // silence warnings

	m := dsl.NewModule(mc.Self)

	agreementpbdsl.UponInputValue(m, func(roundNum uint64, input bool) error {
		// queue input for later
		state.pendingInput[roundNum] = pendingInput{
			input:   input,
			relTime: time.Since(timeRef),
		}
		logger.Log(logging.LevelDebug, "queued input to agreement round", "agRound", roundNum, "value", input)
		return nil
	})

	abbapbdsl.UponDeliver(m, func(decision bool, srcModule t.ModuleID) error {
		round, err := mc.agRoundNumber(srcModule)
		if err != nil {
			return err
		}

		// queue result for delivery (when ready)
		logger.Log(logging.LevelDebug, "round delivered", "agRound", round, "decision", decision)
		return state.storeRoundDecision(agRounds, round, decision)
	})

	// deliver undelivered rounds
	// TODO: currently clashes with round clearing in UponDone. change them to allow batching state updates
	dsl.UponStateUpdate(m, func() error {
		currentRound := state.loadRound(agRounds, state.currentRound)
		if currentRound == nil || !currentRound.delivered {
			return nil
		}

		for currentRound != nil && currentRound.delivered {
			logger.Log(logging.LevelDebug, "delivering round", "agRound", state.currentRound, "decision", currentRound.decision)

			agreementpbdsl.Deliver(m, mc.Consumer, state.currentRound, currentRound.decision, currentRound.posQuorumDuration, currentRound.posTotalDuration)
			state.roundDecisionHistory.Push(currentRound.decision)
			if err := state.clearRoundData(agRounds, state.currentRound); err != nil {
				return err
			}

			state.currentRound++
			currentRound = state.loadRound(agRounds, state.currentRound)
		}

		return nil
	})

	abbapbdsl.UponDone(m, func(srcModule t.ModuleID) error {
		// when abba terminates, clean it up
		roundNumber, err := mc.agRoundNumber(srcModule)
		if err != nil {
			return err
		}

		logger.Log(logging.LevelDebug, "garbage-collecting round", "agRound", roundNumber)
		if err := agRounds.MarkSubmodulePast(roundNumber); err != nil {
			return es.Errorf("failed to clean up finished agreement round: %w", err)
		}

		delete(state.rounds, roundNumber)
		return nil
	})

	// give input to new rounds
	dsl.UponStateUpdates(m, func() error {
		for round, input := range state.pendingInput {
			if agRounds.IsInView(round) {
				logger.Log(logging.LevelDebug, "inputting value to agreement round", "agRound", round, "value", input)
				dsl.EmitMirEvent(m, abbapbevents.InputValue(
					mc.agRoundModuleID(round),
					input.input,
				))

				state.ensureRoundInitialized(round).relInputTime = input.relTime
				delete(state.pendingInput, round)
			}

			// else: not enough rounds have terminated/been freed yet
		}

		return nil
	})

	// TODO: find a way to avoid this weird double-package split. it harms both sight and soul.
	// the problem is that go and protoc's errors do the same.
	agreementpbmsgdsl.UponFinishAbbaMessageReceived(m, func(from t.NodeID, round uint64, value bool) error {
		destModule := mc.agRoundModuleID(round)

		// receiving this is also an implicit acknowdledgement of all protocol messages
		// the ABBA instance has terminated in the sender and does not require any further message
		// Note: this must always be done regardless of the current round to ensure there are no leftover
		// messages in reliablenet queues for old ABBA instances
		reliablenetpbdsl.MarkModuleMsgsRecvd(m, mc.ReliableNet, destModule, []t.NodeID{from})

		if state.currentRound <= round {
			// we're behind, so we can make use of this information
			// Note: this must not be done after the relevant round delivers, or it will lead to an
			// infinite loop of FinishAbbaMessages, as it will trigger the PastMessagesRecvd handler.

			dsl.EmitMirEvent(m, transportpbevents.MessageReceived(
				destModule,
				from,
				abbapbmsgs.FinishMessage(destModule, value),
			))
		}

		return nil
	})

	agreementpbdsl.UponStaleMsgsRecvd(m, func(messages []*modringpbtypes.PastMessage) error {
		for _, msg := range messages {
			decision, ok := state.roundDecisionHistory.Get(msg.DestId)
			if !ok {
				continue // we can't help
			}

			if _, ok := msg.Message.Type.(*messagepbtypes.Message_Abba); ok {
				transportpbdsl.SendMessage(m, mc.Net,
					agreementpbmsgs.FinishAbbaMessage(
						mc.Self,
						msg.DestId,
						decision,
					),
					[]t.NodeID{msg.From},
				)
			}
		}

		return nil
	})

	return m
}

func newPastMessageHandler(mc ModuleConfig) func(pastMessages []*modringpbtypes.PastMessage) (*events.EventList, error) {
	return func(pastMessages []*modringpbtypes.PastMessage) (*events.EventList, error) {
		return events.ListOf(
			agreementpbevents.StaleMsgsRecvd(mc.Self, pastMessages).Pb(),
		), nil
	}
}

func newAbbaGenerator(agMc ModuleConfig, agParams ModuleParams, agTunables ModuleTunables, nodeID t.NodeID, logger logging.Logger) func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
	params := abba.ModuleParams{
		InstanceUID: agParams.InstanceUID, // TODO: review
		AllNodes:    agParams.AllNodes,
	}
	tunables := abba.ModuleTunables{
		MaxRoundLookahead: agTunables.MaxAbbaRoundLookahead,
	}

	return func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
		mc := abba.ModuleConfig{
			Self:         id,
			Consumer:     agMc.Self,
			ReliableNet:  agMc.ReliableNet,
			ThreshCrypto: agMc.ThreshCrypto,
			Hasher:       agMc.Hasher,
		}

		mod, err := abba.NewModule(mc, params, tunables, nodeID, logging.Decorate(logger, "Abba: ", "agRound", idx))
		if err != nil {
			return nil, nil, err
		}

		initialEvs := &events.EventList{}
		return mod, initialEvs, nil
	}
}

func newAgRoundSniffer(mc ModuleConfig, params ModuleParams, agRounds *modring.Module, state *state) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	N := len(params.AllNodes)
	strongQuorumThresh := (N-1)/3*2 + 1

	trackInput := func(r *round, from t.NodeID, msg *messagepbtypes.Message) {
		abbaRoundNum, err := mc.abbaRoundNumber(msg.DestModule)
		if err != nil || abbaRoundNum != 0 {
			return
		}

		abbaMsg, ok := msg.Type.(*messagepbtypes.Message_Abba)
		if !ok {
			return
		}

		abbaRoundMsg, ok := abbaMsg.Abba.Type.(*abbapbtypes.Message_Round)
		if !ok {
			return
		}

		abbaRoundInitMsg, ok := abbaRoundMsg.Round.Type.(*abbapbtypes.RoundMessage_Init)
		if !ok || !abbaRoundInitMsg.Init.IsInput {
			return
		}

		if !r.nodesWithInput.Register(from) {
			return
		}

		if abbaRoundInitMsg.Init.Estimate {
			r.nodesWithPosInputCount++

			if r.nodesWithPosInputCount == strongQuorumThresh {
				r.relPosQuorumTime = time.Since(timeRef)
			} else if r.nodesWithPosInputCount == N {
				r.relPosTotalTime = time.Since(timeRef)
			}
		}
	}

	transportpbdsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepbtypes.Message) error {
		roundNum, err := mc.agRoundNumber(msg.DestModule)
		if err != nil {
			return nil // invalid message
		}

		if !agRounds.IsInView(roundNum) {
			return nil // round not being processed yet
		}

		// track inputs
		r := state.ensureRoundInitialized(roundNum)
		trackInput(r, from, msg)

		return nil
	})

	dsl.UponOtherMirEvent(m, func(ev *eventpbtypes.Event) error {
		switch e := ev.Type.(type) {
		case *eventpbtypes.Event_Abba:
			switch abbaEv := e.Abba.Type.(type) {
			case *abbapbtypes.Event_Round:
				switch abbaRoundEv := abbaEv.Round.Type.(type) {
				case *abbapbtypes.RoundEvent_InputValue:
					agRoundNum, abbaRoundNum, err := parseAgAbbaRoundNumbersFromModule(ev.DestModule, &mc)
					if err != nil {
						return err
					}

					if abbaRoundNum == 0 {
						return nil // skip first round
					}

					r := state.loadRound(agRounds, agRoundNum)
					if r == nil {
						return nil
					}

					if r.abbaRoundNumber+1 != abbaRoundNum {
						return es.Errorf("abba rounds are not being processed sequentially! (expected input for %d, got %d)", r.abbaRoundNumber+1, abbaRoundNum)
					}
					r.abbaRoundNumber = abbaRoundNum
					r.relAbbaRoundStartTime = time.Since(timeRef)
				case *abbapbtypes.RoundEvent_Deliver:
					agRoundNum, err := parseAgRoundNumberFromModule(ev.DestModule, &mc)
					if err != nil {
						return err
					}
					abbaRoundNum := abbaRoundEv.Deliver.RoundNumber

					if abbaRoundNum == 0 {
						return nil // skip first round
					}

					r := state.loadRound(agRounds, agRoundNum)
					if r == nil {
						return nil
					}

					if r.abbaRoundNumber != abbaRoundNum {
						return es.Errorf("abba rounds are not being processed sequentially! (expected deliver from %d, got %d)", r.abbaRoundNumber, abbaRoundNum)
					}

					now := time.Since(timeRef)
					duration := now - r.relAbbaRoundStartTime
					agreementpbdsl.InnerAbbaRoundTime(m, mc.Consumer, duration)

					r.relAbbaRoundStartTime = 0
				}

			}
		}
		return nil
	})

	return m
}

func parseAgRoundNumberFromModule(destModule t.ModuleID, mc *ModuleConfig) (uint64, error) {
	agRoundNumStr := destModule.StripParent(mc.Self).Top()
	agRoundNum, err := strconv.ParseUint(string(agRoundNumStr), 10, 64)
	if err != nil {
		return 0, es.Errorf("failed to parse ag round number from %s (%s): %w", agRoundNumStr, destModule, err)
	}

	return agRoundNum, nil
}

func parseAgAbbaRoundNumbersFromModule(destModule t.ModuleID, mc *ModuleConfig) (uint64, uint64, error) {
	agRoundNum, err := parseAgRoundNumberFromModule(destModule, mc)
	if err != nil {
		return 0, 0, err
	}

	abbaRoundNumStr := destModule.StripParent(mc.Self).Sub().Sub().Top()
	abbaRoundNum, err := strconv.ParseUint(string(abbaRoundNumStr), 10, 64)
	if err != nil {
		return 0, 0, es.Errorf("failed to parse abba round number from %s (%s): %w", abbaRoundNumStr, destModule, err)
	}

	return agRoundNum, abbaRoundNum, nil
}

func (mc *ModuleConfig) agRoundNumber(id t.ModuleID) (uint64, error) {
	roundStr := id.StripParent(mc.Self).Top()
	round, err := strconv.ParseUint(string(roundStr), 10, 64)
	if err != nil {
		return 0, es.Errorf("invalid ag round id: %w", err)
	}

	return round, nil
}

func (mc *ModuleConfig) abbaRoundNumber(id t.ModuleID) (uint64, error) {
	if _, err := mc.agRoundNumber(id); err != nil {
		return 0, err
	}

	if id.StripParent(mc.Self).Sub().Top() != abba.ModringSubName {
		return 0, es.Errorf("invalid abba round id: %v does not match expected structure (missing modring sub name)", id)
	}

	roundStr := id.StripParent(mc.Self).Sub().Sub().Top()
	round, err := strconv.ParseUint(string(roundStr), 10, 64)
	if err != nil {
		return 0, es.Errorf("invalid abba round id: %w", err)
	}

	return round, nil
}
