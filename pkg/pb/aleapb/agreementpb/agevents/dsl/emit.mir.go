package ageventsdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/aleapb/agreementpb/agevents/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue(m dsl.Module, destModule types.ModuleID, round uint64, input bool) {
	dsl.EmitMirEvent(m, events.InputValue(destModule, round, input))
}

func Deliver(m dsl.Module, destModule types.ModuleID, round uint64, decision bool) {
	dsl.EmitMirEvent(m, events.Deliver(destModule, round, decision))
}

func StaleMsgsRecvd(m dsl.Module, destModule types.ModuleID, messages []*types1.PastMessage) {
	dsl.EmitMirEvent(m, events.StaleMsgsRecvd(destModule, messages))
}
