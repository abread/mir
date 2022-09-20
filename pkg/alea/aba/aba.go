package aba

import (
	"fmt"

	abaConfig "github.com/filecoin-project/mir/pkg/alea/aba/config"
	abaDsl "github.com/filecoin-project/mir/pkg/alea/aba/dsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleState struct {
	config      abaConfig.Config
	roundNumber uint64
	roundState  roundStateImpl

	currentEstimate bool

	finishTracker *MsgSupportTrackerBool
}

func New(config abaConfig.Config) modules.PassiveModule {
	m := dsl.NewModule(config.SelfModuleID)
	moduleState := &ModuleState{
		config:      config,
		roundNumber: 0,
		roundState:  nil,

		finishTracker: newMsgSupportTrackerBool(&config),
	}

	abaDsl.UponRoundMessage(m, &moduleState.config, func(from t.NodeID, msg *aleapb.CobaltRoundMessage) error {
		if msg.RoundNumber != moduleState.roundNumber {
			// not in that round
			// TODO: distinguish between lower and higher received round numbers?
			return nil
		}

		switch innerMsg := msg.Type.(type) {
		case *aleapb.CobaltRoundMessage_Init:
			return moduleState.roundState.UponInitMsg(m, moduleState, from, innerMsg.Init.Estimate)
		case *aleapb.CobaltRoundMessage_Aux:
			return moduleState.roundState.UponAuxMsg(m, moduleState, from, innerMsg.Aux.Value)
		case *aleapb.CobaltRoundMessage_Conf:
			return moduleState.roundState.UponConfMsg(m, moduleState, from, innerMsg.Conf.Values)
		case *aleapb.CobaltRoundMessage_CoinTossShare:
			return moduleState.roundState.UponCoinTossShareMsg(m, moduleState, from, innerMsg.CoinTossShare.CoinShare)
		default:
			panic(fmt.Sprintf("unknown ABA protocol message: %T", msg.Type))
		}
	})

	abaDsl.UponFinishMessage(m, &moduleState.config, func(from t.NodeID, value bool) error {
		if !moduleState.finishTracker.RegisterMsg(from, value) {
			return nil // was already received before
		}

		if !moduleState.finishTracker.SelfAlreadySentAnything() {
			// check for weak support
			if moduleState.finishTracker.ValueCountFor(false) >= moduleState.weakSupportThresh() {
				moduleState.finishIfNotAlready(m, false)
			} else if moduleState.finishTracker.ValueCountFor(true) >= moduleState.weakSupportThresh() {
				moduleState.finishIfNotAlready(m, true)
			}
		}

		// check for strong support
		if moduleState.finishTracker.ValueCountFor(false) >= moduleState.strongSupportThresh() {
			// TODO: deliver false
		} else if moduleState.finishTracker.ValueCountFor(true) >= moduleState.strongSupportThresh() {
			// TODO: deliver true
		}

		return nil
	})

	// let rounds track conditions
	dsl.UponCondition(m, func() error {
		return moduleState.roundState.UponCondition(m, moduleState)
	})

	moduleState.initializeRound(m)

	return m
}

func (moduleState *ModuleState) finishIfNotAlready(m dsl.Module, value bool) {
	if !moduleState.finishTracker.SelfAlreadySentAnything() {
		abaDsl.BroadcastFinish(m, &moduleState.config, value)
		moduleState.finishTracker.RegisterMsg(moduleState.config.SelfNodeID(), value)
	}
}

func (moduleState *ModuleState) nextRound(m dsl.Module) {
	moduleState.roundNumber += 1
	moduleState.initializeRound(m)
}

func (moduleState *ModuleState) initializeRound(m dsl.Module) {
	moduleState.roundState = &initialRoundState{}
	moduleState.roundState.Initialize(m, moduleState)
}

func (moduleState *ModuleState) weakSupportThresh() int {
	return moduleState.config.F + 1
}

func (moduleState *ModuleState) strongSupportThresh() int {
	// TODO: review (2f+1 vs N-f)
	return 2*moduleState.config.F + 1
}

type void struct{}

type roundStateImpl interface {
	Initialize(m dsl.Module, moduleState *ModuleState)

	UponInitMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, estimate bool) error
	UponAuxMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, value bool) error
	UponConfMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, values aleapb.CobaltValueSet) error
	UponCoinTossShareMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, coinShare []byte) error

	UponCondition(m dsl.Module, moduleState *ModuleState) error
}
