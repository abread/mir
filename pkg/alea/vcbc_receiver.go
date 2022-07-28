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

func (alea *Alea) applyVCBCSendMessage(id MsgId, msg *aleapb.VCBCSend, source t.NodeID) *events.EventList {
	inst := vcbcReceiverInst{
		payload:   msg.Payload,
		delivered: false,
	}

	if !alea.vcbcReceiverInstances[id.ProposerId].TryPut(uint64(id.Slot), inst) {
		// broadcast instance already created (possibly already ended)
		// TODO: we can maybe detect byzantine behavior from here (SEND message with different messages for the same slot)

		// Safety: we will never sign a different message for the same queue/slot
		return &events.EventList{}
	}

	// TODO: dispatch share-sign event for (id, inst.payload)
	return &events.EventList{}
}

func (alea *Alea) applyThreshShareSignResult(id MsgId, sig threshSigShare) *events.EventList {
	if inst, instOk := alea.vcbcReceiverInstances[id.ProposerId].Get(uint64(id.Slot)); instOk {
		// TODO: does the FINAL message really need a payload?
		finalMsg := VCBCEcho(id, inst.payload, sig)
		return alea.sendMessage(finalMsg, id.ProposerId)
	}

	// nothing to do, broadcast instance was discarded (too old)
	return &events.EventList{}
}

func (alea *Alea) applyVCBCFinalMessage(id MsgId, msg *aleapb.VCBCFinal, source t.NodeID) *events.EventList {
	alea.vcbcReceiverInstances[id.ProposerId].TryPut(uint64(id.Slot), vcbcReceiverInst{
		payload:   msg.Payload,
		delivered: false,
	})

	if inst, ok := alea.vcbcReceiverInstances[id.ProposerId].Get(uint64(id.Slot)); ok {
		// TODO: dispatch signature verification for (id, inst.Payload, msg.Signature)
		inst.payload.String() // XXX: remove - this is just to keep the compiler happy
		return &events.EventList{}
	}

	// instance doesn't exist: too old or was fully-processed already
	return &events.EventList{}
}

func (alea *Alea) applyThreshValidateSignatureResult(id MsgId, sig threshSigFull, sigOk bool) *events.EventList {
	if !sigOk {
		// TODO: log
		// nothing to do, bad signature
		return &events.EventList{}
	}

	if inst, instOk := alea.vcbcReceiverInstances[id.ProposerId].Get(uint64(id.Slot)); instOk {
		inst.delivered = true
		inst.signature = sig

		return (&events.EventList{}).PushBack(VCBCDeliver(id))
	}

	// nothing to do, broadcast instance was discarded (too old/already processed)
	return &events.EventList{}
}
