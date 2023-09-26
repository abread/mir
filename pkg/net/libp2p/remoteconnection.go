package libp2p

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fxamacker/cbor/v2"
	es "github.com/go-errors/errors"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-yamux/v4"

	"github.com/filecoin-project/mir/pkg/logging"
	t "github.com/filecoin-project/mir/pkg/types"
)

type remoteConnection struct {
	params              Params
	ownID               t.NodeID
	addrInfo            *peer.AddrInfo
	logger              logging.Logger
	host                host.Host
	encMode             cbor.EncMode
	streamWriter        *cbor.Encoder
	stream              network.Stream
	msgBuffer           chan PendingMsg
	forceSendMsg        atomic.Value
	forceSendMsgPresent chan struct{}
	stop                chan struct{}
	done                chan struct{}
	connectedCond       *sync.Cond

	currentMsgStatsLabel string
	stats                Stats
}

type PendingMsg struct {
	Payload    []byte
	DestModule t.ModuleID
}

func newRemoteConnection(
	params Params,
	encMode cbor.EncMode,
	ownID t.NodeID,
	addr t.NodeAddress,
	h host.Host,
	logger logging.Logger,
	stats Stats,
) (*remoteConnection, error) {
	addrInfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return nil, es.Errorf("failed to parse address: %w", err)
	}
	conn := &remoteConnection{
		params:              params,
		ownID:               ownID,
		addrInfo:            addrInfo,
		logger:              logger,
		host:                h,
		encMode:             encMode,
		streamWriter:        nil,
		stream:              nil,
		msgBuffer:           make(chan PendingMsg, params.ConnectionBufferSize),
		forceSendMsgPresent: make(chan struct{}),
		stop:                make(chan struct{}),
		done:                make(chan struct{}),
		connectedCond:       sync.NewCond(&sync.Mutex{}),
		stats:               stats,
	}
	go conn.process()
	return conn, nil
}

// PeerID returns the libp2p peer ID of the other side of this connection.
func (conn *remoteConnection) PeerID() peer.ID {
	return conn.addrInfo.ID
}

// Send makes a non-blocking attempt to send a message to this connection.
// Send might use internal buffering. Thus, even if it returns nil,
// the message might not have yet been sent to the network.
func (conn *remoteConnection) Send(msg []byte, destModule t.ModuleID) error {

	select {
	case conn.msgBuffer <- PendingMsg{msg, destModule}:
		return nil
	default:
		return es.Errorf("send buffer full (" + conn.addrInfo.String() + ")")
	}
}

// ForceSend places a message at the top of the send queue in a special slot.
// If a message was previously placed in this slot, it is replaced by the new one.
// This operation is non-blocking.
func (conn *remoteConnection) ForceSend(msg []byte, destModule t.ModuleID) error {
	conn.forceSendMsg.Store(PendingMsg{msg, destModule})

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

// Close closes the connection. No data will be sent to the underlying network stream after Close returns.
func (conn *remoteConnection) Close() {

	// Do nothing if connection already has been closed.
	select {
	case <-conn.stop:
		return
	default:
	}

	// Stop processing and wait until it finishes.
	conn.connectedCond.L.Lock()
	close(conn.stop)
	conn.connectedCond.Broadcast()
	conn.connectedCond.L.Unlock()
	<-conn.done
}

// Wait returns an error channel and a cancel function.
// The channel will be closed without any value being written to it
// when the underlying network stream has been established.
// Waiting is aborted when the cancel function is called or when the connection is closed.
// In both cases, an error is written in the returned channel.
func (conn *remoteConnection) Wait() (chan error, func()) {

	// The channel to be returned.
	result := make(chan error, 1)

	// This flag is set by the returned abort function.
	// It makes the goroutine waiting for the connection return an error.
	abort := false

	go func() {
		conn.connectedCond.L.Lock()
		defer conn.connectedCond.L.Unlock()

		// Wait while
		for conn.stream == nil && !abort {
			// the connection has not yet been established and the caller of Wait has not called the abort function
			select {
			case <-conn.stop:
				// and the connection is not closing.
				result <- es.Errorf("connection closed")
				return
			default:
				conn.connectedCond.Wait()
			}
		}

		// If the waiting was aborted, output an error, otherwise exit successfully (by closing the channel).
		if abort {
			result <- es.Errorf("waiting aborted")
		} else {
			close(result)
		}
	}()

	// Return the result channel and the cancel function.
	return result, func() {
		conn.connectedCond.L.Lock()
		abort = true
		conn.connectedCond.Broadcast()
		conn.connectedCond.L.Unlock()
	}
}

// connect establishes the underlying network stream.
// In case of a network failure, connect keeps retrying until it succeeds (returning nil)
// or until the connection is closed (returning a non-nil error).
// connect blocks until the stream is created or Close is called.
func (conn *remoteConnection) connect() error {

	// Create a context that will be canceled when the connection is closed.
	// It will be used for aborting opening the stream.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // This is necessary for stopping the goroutine below in case of a successful connection.
	go func() {
		select {
		case <-ctx.Done():
			// Garbage-collect this goroutine when the context it is supposed to cancel is canceled by someone else.
		case <-conn.stop:
			// Cancel the new stream creation if connection is closed.
			cancel()
		}
	}()

	// Retry connecting until we succeed.
	for conn.stream == nil {

		select {
		case <-conn.stop:
			// Stop connection attempts if connection is closing.
			return es.Errorf("context canceled")
		default:

			// Try connecting to the peer.
			if err := conn.tryConnecting(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

// tryConnecting makes a single attempt to open the underlying network stream.
// It blocks until it succeeds or, if it fails, at least for ReconnectionPeriod (unless the connection is closing).
// tryConnecting returns nil on success or if it is meaningful to retry connecting (i.e., call tryConnecting again).
// It only returns a non-nil error if it is not meaningful to retry.
func (conn *remoteConnection) tryConnecting(ctx context.Context) error {

	// Try creating a network stream.
	conn.host.Peerstore().AddAddrs(conn.addrInfo.ID, conn.addrInfo.Addrs, conn.params.ConnectionTTL)
	stream, err := conn.host.NewStream(ctx, conn.addrInfo.ID, conn.params.ProtocolID)

	// If connecting failed, wait a moment before returning.
	// (tryConnecting is likely to be called again immediately if it fails.)
	after := time.NewTimer(conn.params.ReconnectionPeriod)
	defer after.Stop()

	if err != nil {
		select {
		case <-conn.stop:
			return es.Errorf("context canceled")
		case <-after.C:
			conn.logger.Log(logging.LevelWarn, "Failed connecting.", "err", err)
			return nil
		}
	}

	streamWriter := conn.encMode.NewEncoder(&disconnectAwareOutputStream{conn, stream})

	// If connecting succeeded, save the new stream
	// and notify any goroutines waiting for the connection establishment (Wait method).
	conn.connectedCond.L.Lock()
	conn.streamWriter = streamWriter
	conn.stream = stream
	conn.connectedCond.Broadcast()
	conn.connectedCond.L.Unlock()

	conn.logger.Log(logging.LevelWarn, "Connected to peer.", "peer", conn.addrInfo.ID.String())
	return nil
}

// process is the main processing loop.
// It keeps reading the input data buffer and writes its contents to the network.
// It automatically creates the underlying network stream and re-establishes it as needed, in case it is dropped.
func (conn *remoteConnection) process() {
	// In case there is a panic in the main processing loop, log an error message.
	// (Otherwise, since this function is run as a goroutine, panicking would be completely silent.)
	defer func() {
		if r := recover(); r != nil {
			err := es.New(r)
			conn.logger.Log(logging.LevelError, "Remote connection panicked.", "cause", r, "stack", err.ErrorStack())
		}
	}()

	// When processing finishes, close the underlying stream and signal to the Stop method that it can return.
	// Note that the defer order is thus inverted.
	defer close(conn.done)
	defer conn.closeStream()

	// Data to be sent to the connection.
	// If nil, a new message from conn.msgBuffer will be read, encoded, and stored here.
	var msgData PendingMsg

	for {
		// The processing loop runs indefinitely (until interrupted by explicitly returning).
		// One iteration corresponds to one attempt of sending a message.

		// Check if connection is being closed.
		// This is necessary for not getting stuck trying to send an unsent message
		// (failing all the time and retrying forever).
		select {
		case <-conn.stop:
			return
		default:
		}

		// Create a network connection if there is none.
		if conn.streamWriter == nil {
			if err := conn.connect(); err != nil {
				// Unless the connection is closing, connect() will keep retrying to connect indefinitely.
				// Thus, if it returns an error, it means that there is no point in continuing the processing.
				conn.logger.Log(logging.LevelWarn, "Gave up connecting", "err", err)
				return
			}
		}

		// Get the next message if there is no pending unsent message.
		if msgData.Payload == nil {
			// prioritize force-send message
			select {
			case <-conn.stop:
				return
			case <-conn.forceSendMsgPresent:
				msgData = conn.forceSendMsg.Swap(PendingMsg{}).(PendingMsg)
			default:
				// Nothing to do in this case, continue.
			}
		}
		for len(msgData.Payload) == 0 {
			select {
			case <-conn.stop:
				return
			// still consider the force send message here
			case <-conn.forceSendMsgPresent:
				msgData = conn.forceSendMsg.Swap(PendingMsg{}).(PendingMsg)
			case msgData = <-conn.msgBuffer:
			}
		}
		conn.currentMsgStatsLabel = string(msgData.DestModule.Top())

		// Write the encoded data to the network stream.
		if err := conn.streamWriter.Encode(TransportMessage{msgData.Payload}); err != nil {
			// If writing fails, close the stream, such that a new one will be re-established in the next iteration.
			conn.logger.Log(logging.LevelWarn, "Failed sending data.", "err", err)
			conn.closeStream()
		} else {
			// On success, clear the pending message (that has just been sent)
			// so a new one can be read from the msbBuffer on the next iteration.
			msgData = PendingMsg{}
		}
	}
}

// closeStream closes the underlying network stream if it is open.
func (conn *remoteConnection) closeStream() {
	conn.connectedCond.L.Lock()
	stream := conn.stream
	// conn.stream == nil is used as a condition in the Wait method and thus needs to be guarded by the lock.
	conn.streamWriter = nil
	conn.stream = nil
	conn.connectedCond.L.Unlock()

	if stream != nil {
		if err := stream.Close(); err != nil {
			conn.logger.Log(logging.LevelWarn, "Failed closing stream.", "err", err)
		}
	}
}

// disconnectAwareOutputStream writes data to the underlying network stream.
// It blocks until all data is written, the connection closes, or an error occurs.
// In the first case, Write returns nil. Otherwise, it returns the corresponding error.
type disconnectAwareOutputStream struct {
	*remoteConnection
	stream network.Stream
}

func (s *disconnectAwareOutputStream) Write(data []byte) (int, error) {
	n := 0

	// Retry sending data until:
	// - all data is sent, or
	// - the connection closes, or
	// - an error occurs.
	for {

		// Set a timeout for the data to be written, so the conn.stream.Write call does not block forever.
		// This is required so that we can periodically check the conn.stop channel.
		if err := s.stream.SetWriteDeadline(time.Now().Add(s.params.StreamWriteTimeout)); err != nil {
			return n, es.Errorf("could not set stream write deadline")
		}

		// Try writing a chunk of data to the underlying network stream.
		var bytesWritten int
		var err error
		if len(data) > s.params.MaxDataPerWrite {
			bytesWritten, err = s.stream.Write(data[:s.params.MaxDataPerWrite])
		} else {
			bytesWritten, err = s.stream.Write(data)
		}
		data = data[bytesWritten:]
		n += bytesWritten

		// Gather statistics if applicable.
		if bytesWritten > 0 && s.stats != nil {
			s.stats.Sent(bytesWritten, s.currentMsgStatsLabel)
		}

		if err == nil && len(data) == 0 {
			// If all data was successfully written, return.
			return n, nil
		} else if errors.Is(err, yamux.ErrTimeout) {
			// If a timeout occurred, check if the connection has not been closed in the meantime.
			// If the connection is still open, retry sending the rest of the data in the next iteration.

			select {
			case <-s.stop:
				return n, es.Errorf("connection closing")
			default:
			}

		} else if err != nil {
			// If any other error occurred, just return it.
			return n, es.Errorf("failed sending data: %w", err)
		}
	}
}
