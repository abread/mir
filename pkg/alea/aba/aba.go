package aba

import (
	"fmt"

	abaConfig "github.com/filecoin-project/mir/pkg/alea/aba/config"
	abaDsl "github.com/filecoin-project/mir/pkg/alea/aba/dsl"
	"github.com/filecoin-project/mir/pkg/alea/util"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	threshDsl "github.com/filecoin-project/mir/pkg/threshcrypto/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleState struct {
	config abaConfig.Config

	round roundSetup

	finishRecvd map[bool]util.Set[t.NodeID]
	finishSent  bool
}

type signShareCtx struct {
	roundNumber uint64
}

type verifyShareCtx struct {
	roundNumber uint64
	nodeID      t.NodeID
	sigShare    []byte
}

type recoverCtx struct {
	roundNumber uint64
}

func New(config abaConfig.Config, initialValue bool) modules.PassiveModule {
	m := dsl.NewModule(config.SelfModuleID)
	moduleState := &ModuleState{
		config: config,

		round: roundSetup{
			number:   0,
			estimate: initialValue,
		},

		finishRecvd: make(map[bool]util.Set[t.NodeID], 2),
		finishSent:  false,
	}

	for _, v := range []bool{true, false} {
		moduleState.finishRecvd[v] = make(util.Set[t.NodeID], len(config.Members))
	}

	abaDsl.UponRoundMessage(m, &moduleState.config, func(from t.NodeID, msg *aleapb.CobaltRoundMessage) error {
		if msg.RoundNumber != moduleState.round.number {
			// not in that round
			// TODO: distinguish between lower and higher received round numbers?
			return nil
		}

		switch innerMsg := msg.Type.(type) {
		case *aleapb.CobaltRoundMessage_Init:
			return moduleState.round.UponInitMsg(m, moduleState, from, innerMsg.Init.Estimate)
		case *aleapb.CobaltRoundMessage_Aux:
			return moduleState.round.UponAuxMsg(m, moduleState, from, innerMsg.Aux.Value)
		case *aleapb.CobaltRoundMessage_Conf:
			return moduleState.round.UponConfMsg(m, moduleState, from, innerMsg.Conf.Values)
		case *aleapb.CobaltRoundMessage_CoinTossShare:
			return moduleState.round.UponCoinTossShareMsg(m, moduleState, from, innerMsg.CoinTossShare.CoinShare)
		default:
			panic(fmt.Sprintf("unknown ABA protocol message: %T", msg.Type))
		}
	})

	abaDsl.UponFinishMessage(m, &moduleState.config, func(from t.NodeID, value bool) error {
		if moduleState.finishRecvd[value].Has(from) || moduleState.finishRecvd[!value].Has(from) {
			// already received
			// TODO: detect/notify byz behavior (finish with !value)
			return nil
		}
		moduleState.finishRecvd[value].Add(from)

		if moduleState.finishRecvd[value].Len() == moduleState.weakSupportThresh() {
			moduleState.finishIfNotAlready(m, value)
		}

		if moduleState.finishRecvd[value].Len() == moduleState.strongSupportThresh() {
			// TODO: deliver value
		}

		return nil
	})

	threshDsl.UponSignShareResult(m, func(sigShare []byte, context *signShareCtx) error {
		return moduleState.round.UponSignShareResult(m, moduleState, sigShare, context)
	})

	threshDsl.UponVerifyShareResult(m, func(ok bool, err string, context *verifyShareCtx) error {

		return moduleState.round.UponVerifyShareResult(m, moduleState, ok, err, context)
	})

	threshDsl.UponRecoverResult(m, func(ok bool, fullSig []byte, err string, context *recoverCtx) error {
		if context.roundNumber != moduleState.round.number {
			return nil // stale
		}

		return moduleState.round.UponRecoverResult(m, moduleState, ok, fullSig, err, context)
	})

	moduleState.round.Initialize(m, moduleState)

	return m
}

func (moduleState *ModuleState) finishIfNotAlready(m dsl.Module, value bool) {
	if !moduleState.finishSent {
		abaDsl.BroadcastFinish(m, &moduleState.config, value)
		moduleState.finishSent = true
	}
}

func (moduleState *ModuleState) nextRound(m dsl.Module) {
	moduleState.round.number += 1
	moduleState.round.Initialize(m, moduleState)
}

func (moduleState *ModuleState) weakSupportThresh() int {
	return moduleState.config.F + 1
}

func (moduleState *ModuleState) strongSupportThresh() int {
	// TODO: review (2f+1 vs N-f)
	return 2*moduleState.config.F + 1
}
