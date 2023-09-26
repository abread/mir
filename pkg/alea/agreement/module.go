package agreement

import (
	"math"
	"strconv"
	"time"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/abba"
	"github.com/filecoin-project/mir/pkg/abba/abbatypes"
	"github.com/filecoin-project/mir/pkg/checkpoint"
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
	directorpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/dsl"
	apppbdsl "github.com/filecoin-project/mir/pkg/pb/apppb/dsl"
	checkpointpbtypes "github.com/filecoin-project/mir/pkg/pb/checkpointpb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	reliablenetpbdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
	reliablenetpbevents "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/events"
	transportpbdsl "github.com/filecoin-project/mir/pkg/pb/transportpb/dsl"
	transportpbevents "github.com/filecoin-project/mir/pkg/pb/transportpb/events"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

var timeRef = time.Now()

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	AleaDirector t.ModuleID
	Hasher       t.ModuleID
	ReliableNet  t.ModuleID
	Net          t.ModuleID
	ThreshCrypto t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self:         "alea_ag",
		AleaDirector: "alea_dir",
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
	EpochLength uint64
}

// ModuleTunables sets the values of protocol tunables that need not be the same across all nodes.
type ModuleTunables struct {
	// Maximum number of concurrent agreement rounds for which we process messages
	// Must be at least 1
	MaxRoundLookahead int

	// Maximum number of concurrent ABBA rounds for which we process messages
	// Must be at least 1
	MaxAbbaRoundLookahead int

	// Maximum number of agreement rounds for which we send input before
	// allowing them to progress in the normal path.
	// Must be at least 0, must be less that MaxRoundLookahead
	// Setting this parameter too high will lead to costly retransmissions!
	// Should likely be less than MaxRoundLookahead/2 - 1.
	MaxRoundAdvanceInput int
}

func NewModule(
	mc ModuleConfig,
	params ModuleParams,
	tunables ModuleTunables,
	startingChkp *checkpoint.StableCheckpoint,
	nodeID t.NodeID,
	logger logging.Logger,
) (modules.PassiveModule, error) {
	if startingChkp == nil {
		return nil, es.Errorf("missing initial checkpoint")
	}
	startingSn := startingChkp.SeqNr()

	if tunables.MaxRoundLookahead <= 0 {
		return nil, es.Errorf("MaxRoundLookahead must be at least 1")
	} else if tunables.MaxAbbaRoundLookahead <= 0 {
		return nil, es.Errorf("MaxAbbaRoundLookahead must be at least 1")
	} else if tunables.MaxRoundAdvanceInput < 0 {
		return nil, es.Errorf("MaxRoundAdvanceInput must be at least 0")
	} else if tunables.MaxRoundAdvanceInput > tunables.MaxRoundLookahead {
		return nil, es.Errorf("MaxRoundAdvanceInput must be less than MaxRoundLookahead")
	} else if startingSn%tt.SeqNr(params.EpochLength) != 0 {
		return nil, es.Errorf("startingSn must be a multiple of EpochLength (the start of an epoch)")
	}

	agRounds := modring.New(
		mc.Self,
		tunables.MaxRoundLookahead,
		modring.ModuleParams{
			Generator: newAbbaGenerator(mc, params, tunables, nodeID, logger),
			CleanupHandler: func(id uint64) (events.EventList, error) {
				// free all remaining pending messages
				return events.ListOf(
					reliablenetpbevents.MarkModuleMsgsRecvd(mc.ReliableNet, mc.agRoundModuleID(id), params.AllNodes),
				), nil
			},
			PastMsgHandler: newPastMessageHandler(mc),
		},
		logging.Decorate(logger, "Modring controller: "),
	)
	if err := agRounds.AdvanceViewToAtLeastSubmodule(uint64(startingSn)); err != nil {
		return nil, err
	}

	N := len(params.AllNodes)
	F := (N - 1) / 3

	state := &state{
		roundDecisionHistory: NewAgRoundHistory(int(params.EpochLength)),

		currentRound: uint64(startingSn),

		rounds: make(map[uint64]*round, tunables.MaxRoundLookahead),

		// stalled inputs count can never go beyond the number of nodes
		pendingInput: make(map[uint64]bool, len(params.AllNodes)),

		N:         N,
		strongMaj: 2*F + 1,
	}
	state.roundDecisionHistory.FreeEpochs(uint64(startingSn) / params.EpochLength)

	controller := newAgController(mc, params, tunables, logger, state, agRounds)

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
	relInputKnownTime      time.Duration
	relPosQuorumTime       time.Duration // instant where 2F+1 nodes provided input=1
	relPosTotalTime        time.Duration // instant where all nodes provided input=1

	continued         bool
	delivered         bool
	decision          bool
	posQuorumDuration time.Duration
	posTotalDuration  time.Duration
}

type state struct {
	roundDecisionHistory AgRoundHistory

	currentRound uint64

	rounds map[uint64]*round

	pendingInput map[uint64]bool

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

	if r.relInputKnownTime == 0 && r.relPosQuorumTime != 0 {
		r.posQuorumDuration = 0
	} else if r.relInputKnownTime != 0 && r.relPosQuorumTime == 0 {
		r.posQuorumDuration = time.Duration(math.MaxInt64)
	} else {
		r.posQuorumDuration = r.relPosQuorumTime - r.relInputKnownTime
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

func newAgController(mc ModuleConfig, params ModuleParams, tunables ModuleTunables, logger logging.Logger, state *state, agRounds *modring.Module) modules.PassiveModule {
	_ = logger // silence warnings

	m := dsl.NewModule(mc.Self)

	agreementpbdsl.UponInputValue(m, func(roundNum uint64, input bool) error {
		// ensure we are not sending duplicate inputs
		_, pendingPresent := state.pendingInput[roundNum]
		round := state.ensureRoundInitialized(roundNum)
		if pendingPresent || round.relInputKnownTime != 0 {
			return es.Errorf("duplicate input for round %v", roundNum)
		}

		round.relInputKnownTime = time.Since(timeRef)

		// queue input for later
		state.pendingInput[roundNum] = input
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
	dsl.UponStateUpdates(m, func() error {
		currentRound, currentRoundRunning := state.rounds[state.currentRound]
		if !currentRoundRunning || !currentRound.delivered {
			return nil
		}

		for currentRoundRunning && currentRound.delivered {
			logger.Log(logging.LevelDebug, "delivering round", "agRound", state.currentRound, "decision", currentRound.decision)

			agreementpbdsl.Deliver(m, mc.AleaDirector, state.currentRound, currentRound.decision, currentRound.posQuorumDuration, currentRound.posTotalDuration)
			state.roundDecisionHistory.Store(state.currentRound, currentRound.decision)

			// clean up metadata
			delete(state.rounds, state.currentRound)

			state.currentRound++
			currentRound, currentRoundRunning = state.rounds[state.currentRound]
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
		evsOut, err := agRounds.MarkSubmodulePast(roundNumber)
		if err != nil {
			return es.Errorf("failed to clean up finished agreement round: %w", err)
		}
		dsl.EmitEvents(m, evsOut)

		return nil
	})

	// give input to new rounds
	dsl.UponStateUpdates(m, func() error {
		for round, input := range state.pendingInput {
			if agRounds.IsInView(round) && round <= state.currentRound+uint64(tunables.MaxRoundAdvanceInput) {
				logger.Log(logging.LevelDebug, "inputting value to agreement round", "agRound", round, "value", input)
				roundModID := mc.agRoundModuleID(round)
				dsl.EmitEvent(m, abbapbevents.InputValue(
					roundModID,
					input,
				))

				delete(state.pendingInput, round)
			}
			// else: not enough rounds have terminated/been freed yet
		}

		return nil
	})

	// force full ABA when current round has input but hasn't delivered yet
	dsl.UponStateUpdates(m, func() error {
		currentRound, ok := state.rounds[state.currentRound]
		if !ok {
			return nil // current round hasn't started at all yet
		}

		// if the input is known, this is the current round, and it is in view, then we must have emitted
		// an ABBA InputValue Event already
		roundHasInput := currentRound.relInputKnownTime != 0 && agRounds.IsInView(state.currentRound)
		if !roundHasInput {
			return nil // current round hasn't started yet (from our end)
		}

		// current round must make progress if unanimity wasn't reached yet
		if !currentRound.delivered && !currentRound.continued {
			logger.Log(logging.LevelDebug, "continuing normal execution", "agRound", state.currentRound)
			abbapbdsl.ContinueExecution(m, mc.agRoundModuleID(state.currentRound))

			currentRound.continued = true
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

			dsl.EmitEvent(m, transportpbevents.MessageReceived(
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
				if msg.DestId < state.currentRound {
					// remote node is behind, ask director to help (send checkpoint)
					directorpbdsl.HelpNode(m, mc.AleaDirector, msg.From)
				}

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

	directorpbdsl.UponNewEpoch(m, func(_ tt.EpochNr) error {
		// business as usual, we don't need to perform any special action
		return nil
	})
	directorpbdsl.UponGCEpochs(m, func(epoch tt.EpochNr) error {
		// we can discard older rounds
		state.roundDecisionHistory.FreeEpochs(uint64(epoch))
		return nil
	})
	apppbdsl.UponRestoreState(m, func(checkpoint *checkpointpbtypes.StableCheckpoint) error {
		// restore the state from the checkpoint
		state.currentRound = uint64(checkpoint.Sn)

		// its simpler to destroy data from all older epochs, because we cannot easily reconstruct it
		state.roundDecisionHistory.FreeEpochs(uint64(checkpoint.Snapshot.EpochData.EpochConfig.EpochNr))
		for round := range state.pendingInput {
			if round < uint64(checkpoint.Sn) {
				delete(state.pendingInput, round)
			}
		}
		for round := range state.rounds {
			if round < uint64(checkpoint.Sn) {
				delete(state.rounds, round)
			}
		}
		if err := agRounds.AdvanceViewToAtLeastSubmodule(uint64(checkpoint.Sn)); err != nil {
			return err
		}
		evsOut, err := agRounds.FreePast()
		if err != nil {
			return err
		}
		dsl.EmitEvents(m, evsOut)

		return nil
	})

	return m
}

func newPastMessageHandler(mc ModuleConfig) func(pastMessages []*modringpbtypes.PastMessage) (events.EventList, error) {
	return func(pastMessages []*modringpbtypes.PastMessage) (events.EventList, error) {
		return events.ListOf(
			agreementpbevents.StaleMsgsRecvd(mc.Self, pastMessages),
		), nil
	}
}

func newAbbaGenerator(agMc ModuleConfig, agParams ModuleParams, agTunables ModuleTunables, nodeID t.NodeID, logger logging.Logger) func(id t.ModuleID, idx uint64) (modules.PassiveModule, events.EventList, error) {
	params := abba.ModuleParams{
		InstanceUID: agParams.InstanceUID, // TODO: review
		AllNodes:    agParams.AllNodes,
	}
	tunables := abba.ModuleTunables{
		MaxRoundLookahead: agTunables.MaxAbbaRoundLookahead,
	}

	return func(id t.ModuleID, idx uint64) (modules.PassiveModule, events.EventList, error) {
		mc := abba.ModuleConfig{
			Self:         id,
			Consumer:     agMc.Self,
			ReliableNet:  agMc.ReliableNet,
			ThreshCrypto: agMc.ThreshCrypto,
			Hasher:       agMc.Hasher,
		}

		mod, err := abba.NewModule(mc, params, tunables, nodeID, logging.Decorate(logger, "Abba: ", "agRound", idx))
		if err != nil {
			return nil, events.EmptyList(), err
		}

		return mod, events.EmptyList(), nil
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

		abbaRoundInputMsg, ok := abbaRoundMsg.Round.Type.(*abbapbtypes.RoundMessage_Input)
		if !ok {
			return
		}

		if !r.nodesWithInput.Register(from) {
			return
		}

		if abbaRoundInputMsg.Input.Estimate {
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

	dsl.UponOtherEvent(m, func(ev *eventpbtypes.Event) error {
		switch e := ev.Type.(type) {
		case *eventpbtypes.Event_Abba:
			switch abbaEv := e.Abba.Type.(type) {
			case *abbapbtypes.Event_Round:
				switch abbaRoundEv := abbaEv.Round.Type.(type) {
				case *abbapbtypes.RoundEvent_Deliver:
					abbaRoundNum := abbaRoundEv.Deliver.RoundNumber
					if abbaRoundNum == 0 {
						return nil // skip first round
					}

					durationNoCoin := abbaRoundEv.Deliver.DurationNoCoin
					agreementpbdsl.InnerAbbaRoundTime(m, mc.AleaDirector, durationNoCoin)
				}

			}
		}
		return nil
	})

	return m
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
