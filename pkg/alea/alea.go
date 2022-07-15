package alea

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type Alea struct {
	ownID  t.NodeID
	logger logging.Logger
	config *Config

	vcbcSenderInstances   map[SlotId]vcbcSenderInst
	vcbcReceiverInstances map[MsgId]vcbcReceiverInst
	abbaInstances map[MsgId]cobaltAbbaInst
}

var (
	aleaModuleName t.ModuleID = "alea"
	netModuleName  t.ModuleID = "net"
	//appModuleName    t.ModuleID = "app"
	//walModuleName    t.ModuleID = "wal"
	//hasherModuleName t.ModuleID = "hasher"
	//cryptoModuleName t.ModuleID = "crypto"
	//timerModuleName  t.ModuleID = "timer"
)

func New(ownID t.NodeID, config *Config, logger *logging.Logger) (*Alea, error) {
	if err := CheckConfig(config); err != nil {
		return nil, fmt.Errorf("invalid Alea configuration: %w", err)
	}

	alea := &Alea{
		ownID:                 ownID,
		logger:                *logger,
		config:                config,
		vcbcSenderInstances:   make(map[SlotId]vcbcSenderInst),
		vcbcReceiverInstances: make(map[MsgId]vcbcReceiverInst),
	}

	// TODO: WAL recovery(?)

	return alea, nil
}

func (alea *Alea) ApplyEvent(event *eventpb.Event) *events.EventList {
	switch e := event.Type.(type) {
	case *eventpb.Event_MessageReceived:
		return alea.applyMessageReceived(e.MessageReceived)
	case *eventpb.Event_Init:
		// TODO: really ensure nothing is required
		return &events.EventList{}
	case *eventpb.Event_Tick:
		return alea.applyTick(e.Tick)
	case *eventpb.Event_HashResult:
		return alea.applyHashResult(e.HashResult)
	case *eventpb.Event_RequestReady:
		return alea.applyRequestReady(e.RequestReady)
	case *eventpb.Event_AppSnapshot:
		panic("TODO: implement snapshotting")
	default:
		panic(fmt.Sprintf("unknown protocol (Alea) event type: %T", event.Type))
	}
}

func (alea *Alea) applyMessageReceived(messageReceived *eventpb.MessageReceived) *events.EventList {
	message := messageReceived.Msg
	from := t.NodeID(messageReceived.From)

	switch msg := message.Type.(*messagepb.Message_Alea).Alea.Type.(type) {
	case *aleapb.AleaMessage_Agreement:
		return alea.applyAgreementMessage(msg.Agreement, from)
	case *aleapb.AleaMessage_Vcbc:
		return alea.applyVcbcMessage(msg.Vcbc, from)
	default:
		panic(fmt.Errorf("unknown Alea message type: %T", msg))
	}
}

func (alea *Alea) applyVcbcMessage(message *aleapb.VCBC, src t.NodeID) *events.EventList {
	id := MsgIdFromDomain(message.Instance)

	switch msg := message.Type.(type) {
	case *aleapb.VCBC_Send:
		return alea.applyVCBCSendMessage(id, msg.Send, src)
	case *aleapb.VCBC_Echo:
		return alea.applyVCBCEchoMessage(id, msg.Echo, src)
	case *aleapb.VCBC_Final:
		return alea.applyVCBCFinalMessage(id, msg.Final, src)
	default:
		panic(fmt.Errorf("unknown Alea VCBC message type: %T", msg))
	}
}

func (alea *Alea) applyRequestReady(requestReady *eventpb.RequestReady) *events.EventList {
	// TODO: plumbing, batching
	return &events.EventList{}
}

func (alea *Alea) applyTick(tick *eventpb.Tick) *events.EventList {
	eventsOut := &events.EventList{}

	// TODO: relay tick to protocol instances to retransmit stuff

	return eventsOut
}

func (alea *Alea) applyHashResult(result *eventpb.HashResult) *events.EventList {
	panic("TODO: route hashresult to correct protocol instance or whatever")
}
