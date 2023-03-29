package reliablenetpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	events "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/events"
	rntypes "github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func SendMessage(m dsl.Module, destModule types.ModuleID, msgId rntypes.MsgID, msg *types1.Message, destinations []types.NodeID) {
	dsl.EmitMirEvent(m, events.SendMessage(destModule, msgId, msg, destinations))
}

func Ack(m dsl.Module, destModule types.ModuleID, destModule0 types.ModuleID, msgId rntypes.MsgID, source types.NodeID) {
	dsl.EmitMirEvent(m, events.Ack(destModule, destModule0, msgId, source))
}

func MarkRecvd(m dsl.Module, destModule types.ModuleID, destModule0 types.ModuleID, msgId rntypes.MsgID, destinations []types.NodeID) {
	dsl.EmitMirEvent(m, events.MarkRecvd(destModule, destModule0, msgId, destinations))
}

func MarkModuleMsgsRecvd(m dsl.Module, destModule types.ModuleID, destModule0 types.ModuleID, destinations []types.NodeID) {
	dsl.EmitMirEvent(m, events.MarkModuleMsgsRecvd(destModule, destModule0, destinations))
}

func RetransmitAll(m dsl.Module, destModule types.ModuleID) {
	dsl.EmitMirEvent(m, events.RetransmitAll(destModule))
}