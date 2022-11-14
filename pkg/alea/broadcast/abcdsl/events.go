package abcdsl

import (
	"github.com/filecoin-project/mir/pkg/alea/broadcast/abcevents"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	aleapb "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func StartBroadcast(m dsl.Module, destModule t.ModuleID, slot uint64, txIDs []t.TxID, txs []*requestpb.Request) {
	dsl.EmitEvent(m, abcevents.StartBroadcast(destModule, slot, txIDs, txs))
}

func Deliver(m dsl.Module, destModule t.ModuleID, slot *aleapb.Slot, txIDs []t.TxID, txs []*requestpb.Request, signature []byte) {
	dsl.EmitEvent(m, abcevents.Deliver(destModule, slot, txIDs, txs, signature))
}

func FreeSlot(m dsl.Module, destModule t.ModuleID, slot *aleapb.Slot) {
	dsl.EmitEvent(m, abcevents.FreeSlot(destModule, slot))
}

func UponStartBroadcast(m dsl.Module, handler func(slot uint64, txIDs []t.TxID, txs []*requestpb.Request) error) {
	UponEvent[*bcpb.Event_StartBroadcast](m, func(ev *bcpb.StartBroadcast) error {
		return handler(ev.QueueSlot, t.TxIDSlice(ev.TxIds), ev.Txs)
	})
}

func UponDeliver(m dsl.Module, handler func(slot *aleapb.Slot, txIDs []t.TxID, txs []*requestpb.Request, signature []byte) error) {
	UponEvent[*bcpb.Event_Deliver](m, func(ev *bcpb.Deliver) error {
		return handler(ev.Slot, t.TxIDSlice(ev.TxIds), ev.Txs, ev.Signature)
	})
}

func UponFreeSlot(m dsl.Module, handler func(slot *aleapb.Slot) error) {
	UponEvent[*bcpb.Event_FreeSlot](m, func(ev *bcpb.FreeSlot) error {
		return handler(ev.Slot)
	})
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
