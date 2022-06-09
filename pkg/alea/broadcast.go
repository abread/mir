package alea

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type void struct{}

type broadcastInst struct {
	initiatorId uint64
	priority    uint64
	message     *requestpb.Batch
	sigShares   map[sigShare]void
}

func (alea *Alea) newBroadcastInstance(priority uint64, message *requestpb.Batch) *events.EventList {
	// TODO: save inst somewhere
	/*inst := broadcastInst{
		message:   message,
		sigShares: make(map[sigShare]void),
	}*/

	sendMsg := &aleapb.VCBC{
		Type: &aleapb.VCBC_Send{},
	}
	msg := broadcastMessage(sendMsg)

	return (&events.EventList{}).PushBack(events.SendMessage(msg, alea.config.Membership))
}

func (alea *Alea) applyBroadcastMessage(broadcastMsg *aleapb.VCBC, source t.NodeID) *events.EventList {
	switch msg := broadcastMsg.Type.(type) {
	case *aleapb.VCBC_Send:

	case *aleapb.VCBC_Echo:
		panic("TODO")
	case *aleapb.VCBC_Final:
		panic("TODO")
	default:
		panic(fmt.Errorf("unknown alea broadcast message type %T", msg))
	}
}

func (inst *broadcastInst) broadcastMessageSend(message *requestpb.Batch) *messagepb.Message {
	return broadcastMessage(
		&aleapb.VCBC{
			Type: &aleapb.VCBC_Send{
				Send: &aleapb.VCBCSend{
					Payload: message,
				},
			},
			Instance: inst.identifierMsg(),
		},
	)
}

func (inst *broadcastInst) broadcastMessageEcho(message *requestpb.Batch, share *sigShare) *messagepb.Message {
	return broadcastMessage(
		&aleapb.VCBC{
			Type: &aleapb.VCBC_Echo{
				Echo: &aleapb.VCBCEcho{
					Payload:        message,
					SignatureShare: share[:],
				},
			},
			Instance: inst.identifierMsg(),
		},
	)
}

func (inst *broadcastInst) broadcastMessageFinal(message *requestpb.Batch, signature *sigFull) *messagepb.Message {
	return broadcastMessage(
		&aleapb.VCBC{
			Type: &aleapb.VCBC_Final{
				Final: &aleapb.VCBCFinal{
					Payload:   message,
					Signature: signature[:],
				},
			},
			Instance: inst.identifierMsg(),
		},
	)
}

func (inst *broadcastInst) identifierMsg() *aleapb.VCBC_VCBCId {
	return &aleapb.VCBC_VCBCId{
		InitiatorId: inst.initiatorId,
		Priority:    inst.priority,
	}
}

func (inst *broadcastInst) broadcastMessage(innerMsg *aleapb.VCBC) *events.EventList {
	msg := &messagepb.Message{
		Type: &messagepb.Message_Alea{
			Alea: &aleapb.AleaMessage{
				Type: &aleapb.AleaMessage_Broadcast{
					Broadcast: innerMsg,
				},
			},
		},
	}

	return (&events.EventList{}).PushBack(events.SendMessage(msg, alea.config.Membership))
}
