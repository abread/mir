package alea

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type vcbcReceiverInst struct {
	payload   *requestpb.Batch
	delivered bool
	signature threshSigFull
}

func (alea *Alea) applyVCBCSendMessage(idPb *aleapb.MsgId, msg *aleapb.VCBCSend, source t.NodeID) *events.EventList {
	id := MsgIdFromDomain(idPb)

	if _, ok := alea.vcbcReceiverInstances[id]; ok {
		// broadcast instance already created
		// TODO: we can detect byzantine behavior from here (SEND message with different messages for the same slot)

		// Safety: we will never sign a different message for the same queue/slot
		return &events.EventList{}
	}

	inst := vcbcReceiverInst{
		payload:   msg.Payload,
		delivered: false,
	}
	alea.vcbcReceiverInstances[id] = inst

	// TODO: dispatch share-sign event for (id, inst.payload)
	return &events.EventList{}
}

func (alea *Alea) applyThreshShareSignResult(id MsgId, sig threshSigShare) *events.EventList {
	if inst, instOk := alea.vcbcReceiverInstances[id]; instOk {
		// TODO: does the FINAL message really need a payload?
		finalMsg := VCBCEcho(id, inst.payload, sig)
		return alea.sendMessage(finalMsg, id.ProposerId)
	}

	// nothing to do, broadcast instance was discarded (too old)
	return &events.EventList{}
}

func (alea *Alea) applyVCBCFinalMessage(idPb *aleapb.MsgId, msg *aleapb.VCBCFinal, source t.NodeID) *events.EventList {
	id := MsgIdFromDomain(idPb)
	// TODO: ignore if id is too old and was already discarded

	inst, ok := alea.vcbcReceiverInstances[id]
	if !ok {
		// TODO: consider removing payload from final message
		inst = vcbcReceiverInst{
			payload:   msg.Payload,
			delivered: false,
		}
		alea.vcbcReceiverInstances[id] = inst
	}

	if inst.delivered {
		// avoid working for nothing
		return &events.EventList{}
	}

	// TODO: dispatch signature verification for (id, inst.Payload, msg.Signature)
	return &events.EventList{}
}

func (alea *Alea) applyThreshValidateSignatureResult(id MsgId, sig threshSigFull, sigOk bool) *events.EventList {
	if !sigOk {
		// TODO: log
		// nothing to do, bad signature
		return &events.EventList{}
	}

	if inst, instOk := alea.vcbcReceiverInstances[id]; instOk {
		inst.delivered = true
		inst.signature = sig
		alea.vcbcReceiverInstances[id] = inst

		return (&events.EventList{}).PushBack(VCBCDeliver(id, inst.payload))
	}

	// nothing to do, broadcast instance was discarded (too old)
	return &events.EventList{}
}
