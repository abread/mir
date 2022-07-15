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
		Type: &aleapb.AleaMessage_Vcbc{
			Vcbc: msg,
		},
	})
}

func CobaltAbbaFinish(id MsgId, v bit) *messagepb.Message {
	return prepareAbaMessage(id, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_Finish{
			Finish: &aleapb.CobaltFinish{
				Value: v.Pb(),
			},
		},
	})
}

func CobaltAbbaInit(id MsgId, r uint64, v bit) *messagepb.Message {
	return prepareAbaMessage(id, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_Init{
			Init: &aleapb.CobaltInit{
				Round: r,
				Value: v.Pb(),
			},
		},
	})
}

func CobaltAbbaAux(id MsgId, r uint64, v bit) *messagepb.Message {
	return prepareAbaMessage(id, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_Aux{
			Aux: &aleapb.CobaltAux{
				Round: r,
				Value: v.Pb(),
			},
		},
	})
}

func CobaltAbbaConf(id MsgId, r uint64, c BitSet) *messagepb.Message {
	return prepareAbaMessage(id, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_Conf{
			Conf: &aleapb.CobaltConf{
				Round:  r,
				Config: c.Pb(),
			},
		},
	})
}

func CobaltAbbaCoinShare(id MsgId, r uint64, share threshSigShare) *messagepb.Message {
	return prepareAbaMessage(id, &aleapb.CobaltABBA{
		Type: &aleapb.CobaltABBA_CoinShare{
			CoinShare: &aleapb.CobaltCoinShare{
				Round:     r,
				CoinShare: share.Pb(),
			},
		},
	})
}

func prepareAbaMessage(id MsgId, msg *aleapb.CobaltABBA) *messagepb.Message {
	return prepareAleaMessage(&aleapb.AleaMessage{
		Type: &aleapb.AleaMessage_Agreement{
			Agreement: &aleapb.Agreement{
				Instance: id.Pb(),
				Msg:      msg,
			},
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

func AbaDeliver(id MsgId, result bit) *eventpb.Event {
	return prepareAleaEvent(&aleapb.AleaEvent{
		Type: &aleapb.AleaEvent_AbaDeliver{
			AbaDeliver: &aleapb.CobaltABBADeliver{
				Id:     id.Pb(),
				Result: result.Pb(),
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
