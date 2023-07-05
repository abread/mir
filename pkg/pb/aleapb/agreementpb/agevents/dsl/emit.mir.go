// Code generated by Mir codegen. DO NOT EDIT.

package ageventsdsl

import (
	"time"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue(m dsl.Module, destModule types.ModuleID, round uint64, input bool) {
	dsl.EmitEvent(m, events.InputValue(destModule, round, input))
}

func Deliver(m dsl.Module, destModule types.ModuleID, round uint64, decision bool, posQuorumWait time.Duration, posTotalWait time.Duration) {
	dsl.EmitEvent(m, events.Deliver(destModule, round, decision, posQuorumWait, posTotalWait))
}

func StaleMsgsRecvd(m dsl.Module, destModule types.ModuleID, messages []*types1.PastMessage) {
	dsl.EmitEvent(m, events.StaleMsgsRecvd(destModule, messages))
}

func InnerAbbaRoundTime(m dsl.Module, destModule types.ModuleID, durationNoCoin time.Duration) {
	dsl.EmitEvent(m, events.InnerAbbaRoundTime(destModule, durationNoCoin))
}
