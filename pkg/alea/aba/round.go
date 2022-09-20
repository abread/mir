package aba

import (
	abaDsl "github.com/filecoin-project/mir/pkg/alea/aba/dsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type initialRoundState struct {
	initTracker MsgSupportTrackerBool
}

func (rs *initialRoundState) Initialize(m dsl.Module, moduleState *ModuleState) {
	rs.initTracker = *newMsgSupportTrackerBool(&moduleState.config)

	abaDsl.BroadcastInit(m, &moduleState.config, moduleState.roundNumber, moduleState.currentEstimate)
}

func (rs *initialRoundState) UponInitMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, estimate bool) error {
	return nil // noop
}

func (rs *initialRoundState) UponAuxMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, value bool) error {
	return nil // noop
}

func (rs *initialRoundState) UponConfMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, values aleapb.CobaltValueSet) error {
	return nil // noop
}

func (rs *initialRoundState) UponCoinTossShareMsg(m dsl.Module, moduleState *ModuleState, from t.NodeID, coinShare []byte) error {
	return nil // noop
}

func (rs *initialRoundState) UponCondition(m dsl.Module, moduleState *ModuleState) error {
	return nil // noop
}
