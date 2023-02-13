package reliablenetpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/types"
	types3 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_ReliableNet](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponSendMessage(m dsl.Module, handler func(msg *types2.Message, destinations []types3.NodeID, msgId []uint8) error) {
	UponEvent[*types.Event_SendMessage](m, func(ev *types.SendMessage) error {
		return handler(ev.Msg, ev.Destinations, ev.MsgId)
	})
}

func UponAck(m dsl.Module, handler func(destModule0 types3.ModuleID, msgId []uint8, source types3.NodeID) error) {
	UponEvent[*types.Event_Ack](m, func(ev *types.Ack) error {
		return handler(ev.DestModule, ev.MsgId, ev.Source)
	})
}

func UponMarkRecvd(m dsl.Module, handler func(destModule0 types3.ModuleID, msgId []uint8, destinations []types3.NodeID) error) {
	UponEvent[*types.Event_MarkRecvd](m, func(ev *types.MarkRecvd) error {
		return handler(ev.DestModule, ev.MsgId, ev.Destinations)
	})
}

func UponMarkModuleMsgsRecvd(m dsl.Module, handler func(destModule0 types3.ModuleID, destinations []types3.NodeID) error) {
	UponEvent[*types.Event_MarkModuleMsgsRecvd](m, func(ev *types.MarkModuleMsgsRecvd) error {
		return handler(ev.DestModule, ev.Destinations)
	})
}

func UponRetransmitAll(m dsl.Module, handler func() error) {
	UponEvent[*types.Event_RetransmitAll](m, func(ev *types.RetransmitAll) error {
		return handler()
	})
}
