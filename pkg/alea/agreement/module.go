package agreement

import (
	"fmt"
	"strconv"

	"github.com/filecoin-project/mir/pkg/abba"
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
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
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

func NewModule(mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	if tunables.MaxRoundLookahead <= 0 {
		return nil, fmt.Errorf("MaxRoundLookahead must be at least 1")
	} else if tunables.MaxAbbaRoundLookahead <= 0 {
		return nil, fmt.Errorf("MaxAbbaRoundLookahead must be at least 1")
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

	controller := newAgController(mc, params, nodeID, logger, agRounds)

	return modules.RoutedModule(mc.Self, controller, agRounds), nil
}

func (mc *ModuleConfig) agRoundModuleID(id uint64) t.ModuleID {
	return mc.Self.Then(t.NewModuleIDFromInt(id))
}

type state struct {
	roundDecisionHistory AgRoundHistory

	currentRound uint64
}

func newAgController(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, logger logging.Logger, agRounds *modring.Module) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	state := state{}

	agreementpbdsl.UponInputValue(m, func(round uint64, input bool) error {
		if round != state.currentRound {
			return fmt.Errorf("round input provided out-of-order. expected round %v got %v", state.currentRound, round)
		}

		// logger.Log(logging.LevelDebug, "inputting value to agreement round", "agRound", round, "value", input)
		dsl.EmitMirEvent(m, abbapbevents.InputValue(
			mc.agRoundModuleID(round),
			input,
		))
		return nil
	})

	abbapbdsl.UponEvent[*abbapbtypes.Event_Deliver](m, func(ev *abbapbtypes.Deliver) error {
		roundStr := ev.SrcModule.StripParent(mc.Self).Top()
		round, err := strconv.ParseUint(string(roundStr), 10, 64)
		if err != nil {
			return fmt.Errorf("deliver event for invalid round: %w", err)
		}

		decision := ev.Result

		if round != state.currentRound {
			return fmt.Errorf("rounds delivered out-of-order. expected round %v got %v", state.currentRound, round)
		}

		// logger.Log(logging.LevelDebug, "delivering round", "agRound", round, "decision", decision)

		agreementpbdsl.Deliver(m, mc.Consumer, round, decision)

		state.roundDecisionHistory.Push(decision)
		state.currentRound++

		if err := agRounds.MarkSubmodulePast(round); err != nil {
			return fmt.Errorf("failed to clean up finished agreement round: %w", err)
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

func newPastMessageHandler(mc *ModuleConfig) func(pastMessages []*modringpbtypes.PastMessage) (*events.EventList, error) {
	return func(pastMessages []*modringpbtypes.PastMessage) (*events.EventList, error) {
		return events.ListOf(
			agreementpbevents.StaleMsgsRecvd(mc.Self, pastMessages).Pb(),
		), nil
	}
}

func newAbbaGenerator(agMc *ModuleConfig, agParams *ModuleParams, agTunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
	params := &abba.ModuleParams{
		InstanceUID: agParams.InstanceUID, // TODO: review
		AllNodes:    agParams.AllNodes,
	}
	tunables := &abba.ModuleTunables{
		MaxRoundLookahead: agTunables.MaxAbbaRoundLookahead,
	}

	return func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error) {
		mc := &abba.ModuleConfig{
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
