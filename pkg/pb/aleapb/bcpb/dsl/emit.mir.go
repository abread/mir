// Code generated by Mir codegen. DO NOT EDIT.

package bcpbdsl

import (
	"time"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func RequestCert(m dsl.Module, destModule types.ModuleID) {
	dsl.EmitEvent(m, events.RequestCert(destModule))
}

func DeliverCert(m dsl.Module, destModule types.ModuleID, cert *types1.Cert) {
	dsl.EmitEvent(m, events.DeliverCert(destModule, cert))
}

func BcStarted(m dsl.Module, destModule types.ModuleID, slot *types1.Slot) {
	dsl.EmitEvent(m, events.BcStarted(destModule, slot))
}

func FreeSlot(m dsl.Module, destModule types.ModuleID, slot *types1.Slot) {
	dsl.EmitEvent(m, events.FreeSlot(destModule, slot))
}

func EstimateUpdate(m dsl.Module, destModule types.ModuleID, maxOwnBcDuration time.Duration, maxOwnBcLocalDuration time.Duration, maxExtBcDuration time.Duration) {
	dsl.EmitEvent(m, events.EstimateUpdate(destModule, maxOwnBcDuration, maxOwnBcLocalDuration, maxExtBcDuration))
}

func DoFillGap(m dsl.Module, destModule types.ModuleID, slot *types1.Slot) {
	dsl.EmitEvent(m, events.DoFillGap(destModule, slot))
}

func MarkStableProposal(m dsl.Module, destModule types.ModuleID, slot *types1.Slot) {
	dsl.EmitEvent(m, events.MarkStableProposal(destModule, slot))
}
