package alea

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type vcbcSenderInst struct {
	payload   *requestpb.Batch
	sigShares map[t.NodeID]threshSigShare
}

func (alea *Alea) newVcbcSenderInstance(slot SlotId, batch *requestpb.Batch) *events.EventList {
	inst := vcbcSenderInst{
		payload:   batch,
		sigShares: make(map[t.NodeID]threshSigShare),
	}

	id := MsgId{
		ProposerId: alea.ownID,
		Slot:       slot,
	}

	if _, ok := alea.vcbcSenderInstances[slot]; ok {
		panic(fmt.Errorf("duplicate broadcast instance for (%s, %ld)", id.ProposerId, id.Slot))
	}

	alea.vcbcSenderInstances[slot] = inst

	sendMsg := VCBCSend(id, batch)
	return alea.broadcastMessage(sendMsg)
}

func (alea *Alea) applyVCBCEchoMessage(id MsgId, msg *aleapb.VCBCEcho, source t.NodeID) *events.EventList {
	if id.ProposerId != alea.ownID {
		// message is not for us
		return &events.EventList{}
	}

	inst, ok := alea.vcbcSenderInstances[id.Slot]
	if !ok {
		// broadcast instance already destroyed/never created
		return &events.EventList{}
	}

	if _, ok := inst.sigShares[source]; ok {
		// already processed
		return &events.EventList{}
	}

	if share, err := threshSigShareFromDomain(msg.SignatureShare); err == nil {
		inst.sigShares[source] = share
	} else {
		// we can't do anything else
		return &events.EventList{}
	}

	// we need 2f+1
	if len(inst.sigShares) >= 2*alea.config.F+1 {
		// TODO: dispatch share combine + validate operation for (id, payload, sigShares)
	}

	return &events.EventList{}
}

func (alea *Alea) applyThreshCombineValidateResult(slot SlotId, sig threshSigFull, sigOk bool) *events.EventList {
	if sigOk {
		if inst, instOk := alea.vcbcSenderInstances[slot]; instOk {
			// TODO: does the FINAL message really need a payload?
			id := MsgId{
				ProposerId: alea.ownID,
				Slot:       slot,
			}
			finalMsg := VCBCFinal(id, inst.payload, sig)
			delete(alea.vcbcSenderInstances, slot)
			return alea.broadcastMessage(finalMsg)
		}
	}

	// we can't do anything
	return &events.EventList{}
}
