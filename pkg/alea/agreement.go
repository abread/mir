package alea

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type agreementImpl struct {
	currentRound uint64
	nextSlotId   map[t.NodeID]uint64

	pendingAgreementRounds InfVec[bool]
}

func (alea *Alea) applyAgreementMessage(msg *aleapb.AgreementMsg, from t.NodeID) *events.EventList {
	if msg.InstanceId != alea.agreement.currentRound {
		// not yet/already done
		return &events.EventList{}
	}

	destModule := abaModuleId(alea, alea.agreement.currentRound)

	return events.ListOf(AbaMessageRecvd(destModule, alea.agreement.currentRound, from, msg.Message))
}

func (alea *Alea) uponAbaDelivery(instanceId uint64, result bool) *events.EventList {
	if instanceId != alea.agreement.currentRound {
		panic("aba delivery for unknown agreement round")
	}

	outEvents := &events.EventList{}

	hasMsg := alea.hasMsgId(alea.agreement.nextSlotId[alea.queueSelection(instanceId)])
	hasPendingDeliveries := alea.agreement.pendingAgreementRounds.Len() > 0
	if hasMsg && !hasPendingDeliveries {
		// TODO: deliver
	} else {
		if !hasMsg {
			// TODO: send FILL-GAP message
		}

		// can't deliver to the application before other instances/VCBC completes
		alea.agreement.pendingAgreementRounds.TryPut(instanceId, result)

		if instanceId < alea.agreement.pendingAgreementRounds.UpperBound() {
			alea.agreement.currentRound += 1
		}
	}

	return outEvents
}

func (alea *Alea) uponFillerMsg() {
	// TODO
}

func abaModuleId(alea *Alea, round uint64) t.ModuleID {
	return alea.config.ABAModuleFactoryName.Then(t.NewModuleIDFromInt(round))
}
