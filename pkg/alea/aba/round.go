package aba

import (
	abaDsl "github.com/filecoin-project/mir/pkg/alea/aba/dsl"
	"github.com/filecoin-project/mir/pkg/alea/util"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	threshDsl "github.com/filecoin-project/mir/pkg/threshcrypto/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

// TODO: stop processing messages when strong support thresh is reached

type roundSetup struct {
	number   uint64
	estimate bool
	values   valueSet

	initRecvd            map[bool]Set[t.NodeID]
	auxSent              bool
	auxFromValuesRecvd   Set[t.NodeID]
	confValueSubsetRecvd Set[t.NodeID]
	coinSharesRecvd      Set[t.NodeID]
	coinSharesOk         [][]byte

	coinSignData [][]byte
}

func (rs *roundSetup) Initialize(m dsl.Module, moduleState *ModuleState) {
	rs.values = newValueSet()

	rs.initRecvd = make(map[bool]Set[t.NodeID], 2)
	for _, v := range []bool{true, false} {
		rs.initRecvd[v] = make(Set[t.NodeID], moduleState.strongSupportThresh())
	}

	rs.auxSent = false

	rs.auxFromValuesRecvd = make(Set[t.NodeID], moduleState.strongSupportThresh())
	rs.confValueSubsetRecvd = make(Set[t.NodeID], moduleState.strongSupportThresh())

	rs.coinSharesRecvd = make(Set[t.NodeID], len(moduleState.config.Members))
	rs.coinSharesOk = make([][]byte, moduleState.strongSupportThresh())

	rs.coinSignData = genCoinSignData(moduleState.config.InstanceId, rs.number)

	abaDsl.BroadcastInit(m, &moduleState.config, moduleState.round.number, moduleState.round.estimate)
}

func (rs *roundSetup) UponInitMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, estimate bool) error {
	if rs.initRecvd[estimate].Has(from) {
		return nil // already processed
	}
	rs.initRecvd[estimate].Add(from)

	if rs.initRecvd[estimate].Len() == moduleState.weakSupportThresh() {
		abaDsl.BroadcastInit(m, &moduleState.config, moduleState.round.number, estimate)
	}

	if rs.initRecvd[estimate].Len() == moduleState.strongSupportThresh() {
		rs.values.Add(estimate)

		if !rs.auxSent {
			abaDsl.BroadcastAux(m, &moduleState.config, moduleState.round.number, estimate)
			rs.auxSent = true
		}
	}

	return nil
}

func (rs *roundSetup) UponAuxMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, value bool) error {
	if !rs.values.Has(value) || rs.auxFromValuesRecvd.Has(from) {
		return nil // already processed or not to be processed
	}
	rs.auxFromValuesRecvd.Add(from)

	if rs.auxFromValuesRecvd.Len() == moduleState.strongSupportThresh() {
		abaDsl.BroadcastConf(m, &moduleState.config, moduleState.round.number, rs.values.Pb())
	}

	return nil
}

func (rs *roundSetup) UponConfMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, values aleapb.CobaltValueSet) error {
	if !valueSet(values).SubsetOf(rs.values) || rs.confValueSubsetRecvd.Has(from) {
		return nil // already processed or not to be processed
	}
	rs.confValueSubsetRecvd.Add(from)

	if rs.confValueSubsetRecvd.Len() == moduleState.strongSupportThresh() {
		threshDsl.SignShare(m, moduleState.config.SelfModuleID, rs.coinSignData, rs)
	}

	return nil
}

func (rs *roundSetup) UponSignShareResult(m dsl.Module, moduleState *ModuleState, sigShare []byte, context *signShareCtx) error {
	if context.roundNumber != moduleState.round.number {
		return nil // stale
	}

	coinShare := sigShare
	abaDsl.BroadcastCoinTossShare(m, &moduleState.config, moduleState.round.number, coinShare)

	return nil
}

func (rs *roundSetup) UponCoinTossShareMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, coinShare []byte) error {
	if rs.coinSharesRecvd.Has(from) {
		// TODO: detect/notify byz behavior
		return nil // already processed
	}
	rs.coinSharesRecvd.Add(from)

	context := &verifyShareCtx{
		roundNumber: rs.number,
		nodeID:      from,
		sigShare:    coinShare,
	}
	threshDsl.VerifyShare(m, moduleState.config.SelfModuleID, rs.coinSignData, coinShare, from, context)

	return nil
}

func (rs *roundSetup) UponVerifyShareResult(m dsl.Module, moduleState *ModuleState, ok bool, err string, context *verifyShareCtx) error {
	if context.roundNumber != moduleState.round.number {
		return nil // stale
	}

	if !ok {
		// TODO: notify byz behavior
		return nil
	}

	rs.coinSharesOk = append(rs.coinSharesOk, context.sigShare)

	if len(rs.coinSharesRecvd) >= moduleState.strongSupportThresh() {
		// TODO: avoid calling Recover on every verified share after 2f+1 are received

		context := &recoverCtx{
			roundNumber: rs.number,
		}
		threshDsl.Recover(m, moduleState.config.SelfModuleID, rs.coinSignData, rs.coinSharesOk, context)
	}

	return nil
}

func (rs *roundSetup) UponRecoverResult(m dsl.Module, moduleState *ModuleState, ok bool, fullSig []byte, err string, context *recoverCtx) error {
	if context.roundNumber != moduleState.round.number {
		return nil // stale
	}

	if !ok {
		return nil
	}

	s_r := sigToCoin(fullSig)

	if rs.values.Len() == 2 {
		rs.estimate = s_r
	} else if rs.values.Len() == 1 {
		if rs.values.Has(true) {
			rs.estimate = true
		} else {
			// rs.values == {false}
			rs.estimate = false
		}

		if rs.values.Has(s_r) {
			moduleState.finishIfNotAlready(m, s_r)
		}
	}

	moduleState.nextRound(m)
	return nil
}

const COIN_SIGN_DATA_PREFIX = "github.com/filecoin-project/mir/pkg/alea/aba"

func genCoinSignData(instanceId uint64, roundNumber uint64) [][]byte {
	return [][]byte{
		[]byte(COIN_SIGN_DATA_PREFIX),
		util.ToBytes(instanceId),
		util.ToBytes(roundNumber),
	}
}

func sigToCoin(sig []byte) bool {
	// TODO: is this safe? or do i need a proper hash?

	b := byte(0)
	for _, bSig := range sig {
		b = b ^ bSig
	}

	finalB := byte(0)
	for i := 0; i < 8; i++ {
		finalB ^= (b >> i) & 1
	}

	return finalB == 0
}
