package alea

import (
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
)

// ============================
// Messages
// ============================

func VCBCSend(id MsgId, batch *requestpb.Batch) *messagepb.Message {
	return prepareBroadcastMessage(&aleapb.VCBC{
		Instance: &aleapb.MsgId{
			QueueIdx: id.ProposerId.Pb(),
			Slot:     id.Slot.Pb(),
		},
		Type: &aleapb.VCBC_Send{
			Send: &aleapb.VCBCSend{
				Payload: batch,
			},
		},
	})
}

func VCBCEcho(id MsgId, batch *requestpb.Batch, sigShare threshSigShare) *messagepb.Message {
	return prepareBroadcastMessage(&aleapb.VCBC{
		Instance: &aleapb.MsgId{
			QueueIdx: id.ProposerId.Pb(),
			Slot:     id.Slot.Pb(),
		},
		Type: &aleapb.VCBC_Echo{
			Echo: &aleapb.VCBCEcho{
				Payload:        batch,
				SignatureShare: sigShare.Pb(),
			},
		},
	})
}

func VCBCFinal(id MsgId, batch *requestpb.Batch, sig threshSigFull) *messagepb.Message {
	return prepareBroadcastMessage(&aleapb.VCBC{
		Instance: &aleapb.MsgId{
			QueueIdx: id.ProposerId.Pb(),
			Slot:     id.Slot.Pb(),
		},
		Type: &aleapb.VCBC_Final{
			Final: &aleapb.VCBCFinal{
				Payload:   batch,
				Signature: sig.Pb(),
			},
		},
	})
}

func prepareBroadcastMessage(msg *aleapb.VCBC) *messagepb.Message {
	return prepareAleaMessage(&aleapb.AleaMessage{
		Type: &aleapb.AleaMessage_Broadcast{
			Broadcast: msg,
		},
	})
}

func prepareAleaMessage(msg *aleapb.AleaMessage) *messagepb.Message {
	return &messagepb.Message{
		Type: &messagepb.Message_Alea{
			Alea: msg,
		},
	}
}

// ============================
// Events
// ============================

func VCBCDeliver(id MsgId, batch *requestpb.Batch) *eventpb.Event {
	return prepareAleaEvent(&aleapb.AleaEvent{
		Type: &aleapb.AleaEvent_VcbcDeliver{
			VcbcDeliver: &aleapb.VCBCDeliver{
				Id:      id.Pb(),
				Payload: batch,
			},
		},
	})
}

func prepareAleaEvent(ev *aleapb.AleaEvent) *eventpb.Event {
	return &eventpb.Event{
		Type: &eventpb.Event_Alea{
			Alea: ev,
		},
	}
}
