package libp2p

import (
	"sync/atomic"

	es "github.com/go-errors/errors"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/protobuf/proto"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	transportpbevents "github.com/filecoin-project/mir/pkg/pb/transportpb/events"
	t "github.com/filecoin-project/mir/pkg/types"
)

// selfConnection represents a connection of a node to itself.
// It bypasses the network completely and feeds sent messages directly into the code that handles message delivery.
type selfConnection struct {
	ownID       t.NodeID
	peerID      peer.ID
	msgBuffer   chan []byte
	deliverChan chan<- events.EventList

	// We receive all messages in selfConnection, so this probably isn't strictly necessary
	// but it's here for complete consistency with remoteConnections
	forceSendMsg        atomic.Value
	forceSendMsgPresent chan struct{}

	stop chan struct{}
	done chan struct{}
}

// newSelfConnection returns a connection to self.
// Addr is the own address and deliverChan is the channel to which delivered messages need to be written.
// Messages sent to this connection will be buffered and eventually written to deliverChan in form of MessageDelivered
// events (unless the buffer fills up, in which case sent messages will be dropped.)
func newSelfConnection(params Params, ownID t.NodeID, ownAddr t.NodeAddress, deliverChan chan<- events.EventList) (*selfConnection, error) {
	addrInfo, err := peer.AddrInfoFromP2pAddr(ownAddr)
	if err != nil {
		return nil, es.Errorf("failed to parse address: %w", err)
	}

	conn := &selfConnection{
		ownID:               ownID,
		peerID:              addrInfo.ID,
		msgBuffer:           make(chan []byte, params.ConnectionBufferSize),
		deliverChan:         deliverChan,
		forceSendMsgPresent: make(chan struct{}),
		stop:                make(chan struct{}),
		done:                make(chan struct{}),
	}

	go conn.process()

	return conn, nil
}

// PeerID returns the peer ID of the node itself, as it is itself on the other side of the connection.
func (conn *selfConnection) PeerID() peer.ID {
	return conn.peerID
}

// Send feeds the given message directly to the sink of delivered messages.
// Send is non-blocking and if the buffer for delivered messages is full, the message is dropped.
func (conn *selfConnection) Send(msg []byte, _ t.ModuleID) error {

	select {
	case conn.msgBuffer <- msg:
		return nil
	default:
		go func() {
			conn.msgBuffer <- msg
		}()
		return es.Errorf("send buffer full")
	}
}

// ForceSend places a message at the top of the send queue in a special slot.
// If a message was previously placed in this slot, it is replaced by the new one.
// This operation is non-blocking.
func (conn *selfConnection) ForceSend(msg []byte, _ t.ModuleID) error {
	conn.forceSendMsg.Store(msg)

	select {
	case conn.forceSendMsgPresent <- struct{}{}:
		return nil
	default:
		// We could not mark the forceSendMsgPresent channel as present, but this is fine.
		// It means the worker goroutine hasn't processed the force message yet and will consider
		// the new message.
		return nil
	}
}

// Close makes the connection stop processing messages. No messages will be delivered by it after Close returns.
func (conn *selfConnection) Close() {

	// Do nothing if connection already has been closed.
	select {
	case <-conn.stop:
		return
	default:
	}

	// Stop processing and wait until it finishes.
	close(conn.stop)
	<-conn.done
}

// Wait returns a channel and a function.
// Since this is a connection to self, there is no underlying network stream to wait for.
// If the connection is closed, the channel will contain a single error value written to it and will never be closed.
// Otherwise, the channel will be closed without any values.
// The returned function has no effect.
func (conn *selfConnection) Wait() (chan error, func()) {
	result := make(chan error, 1)
	select {
	case <-conn.stop:
		result <- es.Errorf("connection closed")
	default:
		close(result)
	}
	return result, func() {}
}

// process shovels messages from the message buffer to the deliver channel.
func (conn *selfConnection) process() {
	// When done, make the Close method return.
	defer close(conn.done)

	recvMsg := func(encodedMsg []byte) {
		if len(encodedMsg) == 0 {
			return
		}

		var msgPb messagepb.Message
		err := proto.Unmarshal(encodedMsg, &msgPb)
		if err != nil {
			panic(es.Errorf("failed to unmarshal message from self: %w", err))
		}

		msg := messagepbtypes.MessageFromPb(&msgPb)

		select {
		case <-conn.stop:
			return
		case conn.deliverChan <- events.ListOf(transportpbevents.MessageReceived(
			msg.DestModule,
			conn.ownID,
			msg,
		)):
			// Nothing to do in this case, message has been delivered.
		}
	}

	for {
		select {
		case <-conn.stop:
			return
		default:
			// Nothing to do in this case, continue.
		}

		// Give priority to force send
		select {
		case <-conn.stop:
			return
		case <-conn.forceSendMsgPresent:
			recvMsg(conn.forceSendMsg.Swap([]byte{}).([]byte))
		default:
			// Nothing to do in this case, continue.
		}

		select {
		case <-conn.stop:
			return
		// still consider force send here
		case <-conn.forceSendMsgPresent:
			recvMsg(conn.forceSendMsg.Swap([]byte{}).([]byte))
		case encodedMsg := <-conn.msgBuffer:
			recvMsg(encodedMsg)
		}
	}
}
