package abcevents

import (
	"github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	aleapb "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func StartBroadcast(destModule t.ModuleID, slot uint64, txs []*requestpb.Request) *eventpb.Event {
	return Event(destModule, &bcpb.Event{
		Type: &bcpb.Event_StartBroadcast{
			StartBroadcast: &bcpb.StartBroadcast{
				QueueSlot: slot,
				Txs:       txs,
			},
		},
	})
}

func Deliver(destModule t.ModuleID, slot *aleapb.Slot, txIDs []t.TxID, txs []*requestpb.Request, signature []byte) *eventpb.Event {
	return Event(destModule, &bcpb.Event{
		Type: &bcpb.Event_Deliver{
			Deliver: &bcpb.Deliver{
				Slot:      slot,
				TxIds:     t.TxIDSlicePb(txIDs),
				Txs:       txs,
				Signature: signature,
			},
		},
	})
}

func FreeSlot(destModule t.ModuleID, slot *aleapb.Slot) *eventpb.Event {
	return Event(destModule, &bcpb.Event{
		Type: &bcpb.Event_FreeSlot{
			FreeSlot: &bcpb.FreeSlot{
				Slot: slot,
			},
		},
	})
}

func Event(destModule t.ModuleID, ev *bcpb.Event) *eventpb.Event {
	return &eventpb.Event{
		Type: &eventpb.Event_AleaBroadcast{
			AleaBroadcast: ev,
		},
		DestModule: destModule.Pb(),
	}
}
