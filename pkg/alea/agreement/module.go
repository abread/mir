package agreement

import (
	"fmt"

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
	agreementpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/dsl"
	agreementpbmsgs "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/msgs"
	agreementpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/types"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	eventpbevents "github.com/filecoin-project/mir/pkg/pb/eventpb/events"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpbdsl "github.com/filecoin-project/mir/pkg/pb/modringpb/dsl"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	reliablenetpbdsl "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/dsl"
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

const modringSubName = t.ModuleID("r")

func NewModule(mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.PassiveModule, error) {
	if tunables.MaxRoundLookahead <= 0 {
		return nil, fmt.Errorf("MaxRoundLookahead must be at least 1")
	} else if tunables.MaxAbbaRoundLookahead <= 0 {
		return nil, fmt.Errorf("MaxAbbaRoundLookahead must be at least 1")
	}

	controller := newAgController(mc, params, nodeID, logger)

	slots := modring.New(
		&modring.ModuleConfig{
			Self:            mc.Self.Then(modringSubName),
			PastMsgConsumer: mc.Self,
		},
		tunables.MaxRoundLookahead,
		modring.ModuleParams{
			Generator: newAbbaGenerator(mc, params, tunables, nodeID, logger),
		},
		logging.Decorate(logger, "Modring controller: "),
	)

	return modules.RoutedModule(mc.Self, controller, slots), nil
}

func (mc *ModuleConfig) agRoundModuleID(id uint64) t.ModuleID {
	return mc.Self.Then(modringSubName).Then(t.NewModuleIDFromInt(id))
}

func (mc *ModuleConfig) agRoundOrigin(roundID uint64) *abbapbtypes.Origin {
	return &abbapbtypes.Origin{
		Module: mc.Self,
		Type: &abbapbtypes.Origin_AleaAg{
			AleaAg: &agreementpbtypes.AbbaOrigin{
				AgRound: roundID,
			},
		},
	}
}

type state struct {
	roundDecisionHistory AgRoundHistory

	currentRound uint64
}

type roundDeliverCtx struct {
	round    uint64
	decision bool
}

func newAgController(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID, logger logging.Logger) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	state := state{}

	dsl.UponOtherEvent(m, func(ev *eventpb.Event) error {
		logger.Log(logging.LevelDebug, "unknown event", "ev", ev)
		return nil
	})

	agreementpbdsl.UponInputValue(m, func(round uint64, input bool) error {
		if round != state.currentRound {
			return fmt.Errorf("round input provided out-of-order. expected round %v got %v", state.currentRound, round)
		}

		logger.Log(logging.LevelDebug, "inputting value to agreement round", "agRound", round, "value", input)
		dsl.EmitMirEvent(m, abbapbevents.InputValue(
			mc.agRoundModuleID(round),
			input,
			mc.agRoundOrigin(round),
		))
		return nil
	})

	abbapbdsl.UponEvent[*abbapbtypes.Event_Deliver](m, func(ev *abbapbtypes.Deliver) error {
		round := ev.Origin.Type.(*abbapbtypes.Origin_AleaAg).AleaAg.AgRound

		if round != state.currentRound {
			return fmt.Errorf("rounds delivered out-of-order. expected round %v got %v", state.currentRound, round)
		}

		context := &roundDeliverCtx{
			round:    round,
			decision: ev.Result,
		}
		modringpbdsl.FreeSubmodule(m, mc.Self.Then(modringSubName), round, context)

		return nil
	})
	modringpbdsl.UponFreedSubmodule(m, func(context *roundDeliverCtx) error {
		if context.round != state.currentRound {
			return fmt.Errorf("race in round cleanup. expected round %v got %v", state.currentRound, context.round)
		}

		logger.Log(logging.LevelDebug, "delivering round", "round", "agRound", context.round, "decision", context.decision)
		agreementpbdsl.Deliver(m, mc.Consumer, context.round, context.decision)

		state.roundDecisionHistory.Push(context.decision)

		state.currentRound++
		return nil
	})

	agreementpbdsl.UponFinishAbbaMessageReceived(m, func(from t.NodeID, round uint64, value bool) error {
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

			dsl.EmitMirEvent(m, eventpbevents.MessageReceived(
				destModule,
				from,
				abbapbmsgs.FinishMessage(destModule, value),
			))
		}

		return nil
	})

	modringpbdsl.UponPastMessagesRecvd(m, func(messages []*modringpbtypes.PastMessage) error {
		for _, msg := range messages {
			decision, ok := state.roundDecisionHistory.Get(msg.DestId)
			if !ok {
				continue // we can't help
			}

			if _, ok := msg.Message.Type.(*messagepbtypes.Message_Abba); ok {
				dsl.SendMessage(m, mc.Net,
					agreementpbmsgs.FinishAbbaMessage(
						mc.Self,
						msg.DestId,
						decision,
					).Pb(),
					[]t.NodeID{msg.From},
				)
			}
		}
		return nil
	})

	return m
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
