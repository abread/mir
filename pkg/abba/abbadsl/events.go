package abbadsl

import (
	abbaEvents "github.com/filecoin-project/mir/pkg/abba/abbaevents"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/abbapb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func InputValue(m dsl.Module, dest t.ModuleID, input bool) {
	dsl.EmitEvent(m, abbaEvents.InputValue(dest, input))
}

func Deliver(m dsl.Module, dest t.ModuleID, result bool) {
	dsl.EmitEvent(m, abbaEvents.Deliver(dest, result, m.ModuleID()))
}

// Module-specific dsl functions for processing events.

func UponEvent[EvWrapper abbapb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*eventpb.Event_Abba](m, func(ev *abbapb.Event) error {
		evWrapper, ok := ev.Type.(EvWrapper)
		if !ok {
			return nil
		}
		return handler(evWrapper.Unwrap())
	})
}

func UponInputValue(m dsl.Module, handler func(input bool) error) {
	UponEvent[*abbapb.Event_InputValue](m, func(ev *abbapb.InputValue) error {
		return handler(ev.Input)
	})
}

func UponDeliver(m dsl.Module, handler func(result bool, from t.ModuleID) error) {
	UponEvent[*abbapb.Event_Deliver](m, func(ev *abbapb.Deliver) error {
		return handler(ev.Result, t.ModuleID(ev.OriginModule))
	})
}

func UponAbbaMessageReceived(m dsl.Module, handler func(from t.NodeID, msg *abbapb.Message) error) {
	dsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepb.Message) error {
		cbMsgWrapper, ok := msg.Type.(*messagepb.Message_Abba)
		if !ok {
			return nil
		}

		return handler(from, cbMsgWrapper.Abba)
	})
}

func UponFinishMessageReceived(m dsl.Module, handler func(from t.NodeID, value bool) error) {
	UponAbbaMessageReceived(m, func(from t.NodeID, msg *abbapb.Message) error {
		finishMsgWrapper, ok := msg.Type.(*abbapb.Message_FinishMessage)
		if !ok {
			return nil
		}

		return handler(from, finishMsgWrapper.FinishMessage.Value)
	})
}

func UponInitMessageReceived(m dsl.Module, handler func(from t.NodeID, roundNumber uint64, estimate bool) error) {
	UponAbbaMessageReceived(m, func(from t.NodeID, msg *abbapb.Message) error {
		initMsgWrapper, ok := msg.Type.(*abbapb.Message_InitMessage)
		if !ok {
			return nil
		}

		initMsg := initMsgWrapper.InitMessage
		return handler(from, initMsg.RoundNumber, initMsg.Estimate)
	})
}

func UponAuxMessageReceived(m dsl.Module, handler func(from t.NodeID, roundNumber uint64, value bool) error) {
	UponAbbaMessageReceived(m, func(from t.NodeID, msg *abbapb.Message) error {
		auxMsgWrapper, ok := msg.Type.(*abbapb.Message_AuxMessage)
		if !ok {
			return nil
		}

		auxMsg := auxMsgWrapper.AuxMessage
		return handler(from, auxMsg.RoundNumber, auxMsg.Value)
	})
}

func UponConfMessageReceived(m dsl.Module, handler func(from t.NodeID, roundNumber uint64, values ValueSet) error) {
	UponAbbaMessageReceived(m, func(from t.NodeID, msg *abbapb.Message) error {
		confMsgWrapper, ok := msg.Type.(*abbapb.Message_ConfMessage)
		if !ok {
			return nil
		}

		confMsg := confMsgWrapper.ConfMessage
		return handler(from, confMsg.RoundNumber, ValueSet(confMsg.Values))
	})
}

func UponCoinMessageReceived(m dsl.Module, handler func(from t.NodeID, roundNumber uint64, coinShare []byte) error) {
	UponAbbaMessageReceived(m, func(from t.NodeID, msg *abbapb.Message) error {
		coinMsgWrapper, ok := msg.Type.(*abbapb.Message_CoinMessage)
		if !ok {
			return nil
		}

		coinMsg := coinMsgWrapper.CoinMessage
		return handler(from, coinMsg.RoundNumber, coinMsg.CoinShare)
	})
}
