package agreement

import (
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

	state := &state{
		rounds: make([]round, tunables.MaxRoundLookahead),

		// stalled inputs count can never go beyond the number of concurrent agreement rounds
		pendingInput: make(map[uint64]bool, tunables.MaxRoundLookahead),
	}
	for i := range state.rounds {
		state.rounds[i].participants = make(abbatypes.RecvTracker, (len(params.AllNodes)-1)/3*2+1)
	}

	controller := newAgController(mc, params, logger, state, agRounds)

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
	participants abbatypes.RecvTracker // TODO: move recvtracker to another package?
	startTime    time.Time

	delivered bool
	decision  bool
	duration  time.Duration
}

type state struct {
	roundDecisionHistory AgRoundHistory

	currentRound uint64

	rounds []round

	pendingInput map[uint64]bool
}

func (s *state) storeRoundDecision(rounds *modring.Module, roundNum uint64, decision bool) error {
	if !rounds.IsInView(roundNum) {
		return es.Errorf("round out of view: %v", rounds)
	}

	idx := int(roundNum % uint64(len(s.rounds)))
	s.rounds[idx].decision = decision
	s.rounds[idx].delivered = true

	startTime := s.rounds[idx].startTime
	s.rounds[idx].duration = time.Since(startTime)

	return nil
}

func (s *state) loadRound(rounds *modring.Module, roundNum uint64) *round {
	if !rounds.IsInView(roundNum) {
		return nil
	}

	idx := int(roundNum % uint64(len(s.rounds)))
	return &s.rounds[idx]
}

func (s *state) clearRoundData(rounds *modring.Module, roundNum uint64, n int) error {
	if !rounds.IsInView(roundNum) {
		return es.Errorf("round out of view: %v", rounds)
	}

	idx := int(roundNum % uint64(len(s.rounds)))
	s.rounds[idx] = round{
		participants: make(abbatypes.RecvTracker, (n-1)/3*2+1),
	}

	return nil
}

func newAgController(mc ModuleConfig, params ModuleParams, logger logging.Logger, state *state, agRounds *modring.Module) modules.PassiveModule {
	_ = logger // silence warnings

	m := dsl.NewModule(mc.Self)

	agreementpbdsl.UponInputValue(m, func(roundNum uint64, input bool) error {
		// queue input for later
		state.pendingInput[roundNum] = input
		logger.Log(logging.LevelDebug, "queued input to agreement round", "agRound", roundNum, "value", input)
		return nil
	})

	abbapbdsl.UponDeliver(m, func(decision bool, srcModule t.ModuleID) error {
		round, err := mc.abbaRoundNumber(srcModule)
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

			agreementpbdsl.Deliver(m, mc.Consumer, state.currentRound, currentRound.decision, currentRound.duration)
			state.roundDecisionHistory.Push(currentRound.decision)
			if err := state.clearRoundData(agRounds, state.currentRound, len(params.AllNodes)); err != nil {
				return err
			}

			state.currentRound++

			currentRound = state.loadRound(agRounds, state.currentRound)
		}

		return nil
	})

	abbapbdsl.UponDone(m, func(srcModule t.ModuleID) error {
		// when abba terminates, clean it up
		round, err := mc.abbaRoundNumber(srcModule)
		if err != nil {
			return err
		}

		logger.Log(logging.LevelDebug, "garbage-collecting round", "agRound", round)
		if err := agRounds.MarkSubmodulePast(round); err != nil {
			return es.Errorf("failed to clean up finished agreement round: %w", err)
		}

		return nil
	})

	// give input to new rounds
	dsl.UponStateUpdates(m, func() error {
		for round, input := range state.pendingInput {
			if agRounds.IsInView(round) {
				logger.Log(logging.LevelDebug, "inputting value to agreement round", "agRound", round, "value", input)
				dsl.EmitMirEvent(m, abbapbevents.InputValue(
					mc.agRoundModuleID(round),
					input,
				))

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
	m := dsl.NewModule("__invalid__") // this module shall not emit events

	N := len(params.AllNodes)
	thresh := (N-1)/3*2 + 1

	transportpbdsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepbtypes.Message) error {
		roundNum, err := mc.abbaRoundNumber(msg.DestModule)
		if err != nil {
			return err
		}

		if !agRounds.IsInView(roundNum) {
			return nil // round not being processed yet
		}

		idx := int(roundNum % uint64(len(state.rounds)))
		var zeroTime time.Time
		if state.rounds[idx].startTime == zeroTime {
			state.rounds[idx].participants.Register(from)

			if state.rounds[idx].participants.Len() >= thresh {
				state.rounds[idx].startTime = time.Now()
			}
		}

		return nil
	})

	dsl.UponOtherMirEvent(m, func(_ *eventpbtypes.Event) error {
		return nil
	})

	return m
}

func (mc *ModuleConfig) abbaRoundNumber(id t.ModuleID) (uint64, error) {
	roundStr := id.StripParent(mc.Self).Top()
	round, err := strconv.ParseUint(string(roundStr), 10, 64)
	if err != nil {
		return 0, es.Errorf("deliver event for invalid round: %w", err)
	}

	return round, nil
}
