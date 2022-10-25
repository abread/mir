package abcdsl

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	aleapb "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func StartBroadcast(m dsl.Module, destModule t.ModuleID, slot uint64, txs []*requestpb.Request) {
	EmitEvent(m, destModule, &bcpb.Event{
		Type: &bcpb.Event_StartBroadcast{
			StartBroadcast: &bcpb.StartBroadcast{
				QueueSlot: slot,
				Txs:       txs,
			},
		},
	})
}

func Deliver(m dsl.Module, destModule t.ModuleID, slot *aleapb.Slot, txs []*requestpb.Request, signature []byte) {
	EmitEvent(m, destModule, &bcpb.Event{
		Type: &bcpb.Event_Deliver{
			Deliver: &bcpb.Deliver{
				Slot: slot,
				Txs:  txs,
			},
		},
	})
}

func FreeSlot(m dsl.Module, destModule t.ModuleID, slot *aleapb.Slot) {
	EmitEvent(m, destModule, &bcpb.Event{
		Type: &bcpb.Event_FreeSlot{
			FreeSlot: &bcpb.FreeSlot{
				Slot: slot,
			},
		},
	})
}

func UponStartBroadcast(m dsl.Module, handler func(slot uint64, txs []*requestpb.Request) error) {
	UponEvent[*bcpb.Event_StartBroadcast](m, func(ev *bcpb.StartBroadcast) error {
		return handler(ev.QueueSlot, ev.Txs)
	})
}

func UponDeliver(m dsl.Module, handler func(slot *aleapb.Slot, txs []*requestpb.Request) error) {
	UponEvent[*bcpb.Event_Deliver](m, func(ev *bcpb.Deliver) error {
		return handler(ev.Slot, ev.Txs)
	})
}

func UponFreeSlot(m dsl.Module, handler func(slot *aleapb.Slot) error) {
	UponEvent[*bcpb.Event_FreeSlot](m, func(ev *bcpb.FreeSlot) error {
		return handler(ev.Slot)
	})
}

func EmitEvent(m dsl.Module, destModule t.ModuleID, ev *bcpb.Event) {
	evWrapped := &eventpb.Event{
		Type: &eventpb.Event_AleaBroadcast{
			AleaBroadcast: ev,
		},
		DestModule: destModule.Pb(),
	}

	dsl.EmitEvent(m, evWrapped)
}

func UponEvent[EvWrapper bcpb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*eventpb.Event_AleaBroadcast](m, func(ev *bcpb.Event) error {
		evWrapper, ok := ev.Type.(EvWrapper)
		if !ok {
			return nil
		}
		return handler(evWrapper.Unwrap())
	})
}
