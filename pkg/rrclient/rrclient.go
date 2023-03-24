package rrclient

import (
	"context"
	"crypto"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/requestreceiver"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

const (
	// Maximum size of a gRPC message
	maxMessageSize = 1073741824
)

// TODO: Update the comments around crypto, hasher, and request signing.

type RoundRobinClient struct {
	ownID     t.ClientID
	hasher    crypto.Hash
	nextReqNo t.ReqNo
	conns     map[t.NodeID]*grpc.ClientConn
	clients   map[t.NodeID]requestreceiver.RequestReceiver_ListenClient
	clientIDs []t.NodeID
	logger    logging.Logger

	nextClientIdx int
}

func NewRoundRobinClient(
	clientID t.ClientID,
	hasher crypto.Hash,
	l logging.Logger,
) *RoundRobinClient {

	// If no logger was given, only write errors to the console.
	if l == nil {
		l = logging.ConsoleErrorLogger
	}

	return &RoundRobinClient{
		ownID:     clientID,
		hasher:    hasher,
		nextReqNo: 0,
		clients:   make(map[t.NodeID]requestreceiver.RequestReceiver_ListenClient),
		conns:     make(map[t.NodeID]*grpc.ClientConn),
		logger:    l,
	}
}

// Connect establishes (in parallel) network connections to all nodes in the system.
// The nodes' RequestReceivers must be running.
// Only after Connect() returns, sending requests through this RoundRobinClient is possible.
// TODO: Deal with errors, e.g. when the connection times out (make sure the RPC call in connectToNode() has a timeout).
func (rrc *RoundRobinClient) Connect(ctx context.Context, membership map[t.NodeID]string) {
	// Initialize wait group used by the connecting goroutines
	wg := sync.WaitGroup{}
	wg.Add(len(membership))

	// Synchronizes concurrent access to connections.
	lock := sync.Mutex{}

	// For each node in the membership
	for nodeID, nodeAddr := range membership {

		// Launch a goroutine that connects to the node.
		go func(id t.NodeID, addr string) {
			defer wg.Done()

			// Create and store connection
			conn, sink, err := rrc.connectToNode(ctx, addr) // May take long time, execute before acquiring the lock.
			lock.Lock()
			rrc.conns[id] = conn
			rrc.clients[id] = sink
			lock.Unlock()

			// Print debug info.
			if err != nil {
				rrc.logger.Log(logging.LevelWarn, "Failed to connect to node.", "id", id, "addr", addr, "err", err)
			} else {
				rrc.logger.Log(logging.LevelDebug, "Node connected.", "id", id, "addr", addr)
			}

		}(nodeID, nodeAddr)
	}

	// Wait for connecting goroutines to finish.
	wg.Wait()

	rrc.clientIDs = maputil.GetSortedKeys(rrc.clients)
}

// SubmitRequest submits a request by sending it to one node, chosen in a round-robin fashion from
// the full node list (as configured when creating the RoundRobinClient).
// It automatically appends meta-info like client ID and request number.
// SubmitRequest must not be called concurrently.
// If an error occurs, SubmitRequest returns immediately.
func (rrc *RoundRobinClient) SubmitRequest(data []byte) error {

	// Create new request message.
	reqMsg := events.ClientRequest(rrc.ownID, rrc.nextReqNo, data)
	rrc.nextReqNo++

	nID := rrc.clientIDs[rrc.nextClientIdx]
	client := rrc.clients[nID]
	rrc.nextClientIdx++

	if err := client.Send(reqMsg); err != nil {
		return fmt.Errorf("failed sending request to node (%v): %w", nID, err)
	}

	return nil
}

// Disconnect closes all open connections to Mir nodes.
func (rrc *RoundRobinClient) Disconnect() {
	// Close connections to all nodes.
	for id, client := range rrc.clients {
		if client == nil {
			rrc.logger.Log(logging.LevelWarn, fmt.Sprintf("No gRPC client to close to node %v", id))
		} else if _, err := client.CloseAndRecv(); err != nil {
			rrc.logger.Log(logging.LevelWarn, fmt.Sprintf("Could not close gRPC client %v", id))
		}
	}

	for id, conn := range rrc.conns {
		if conn == nil {
			rrc.logger.Log(logging.LevelWarn, fmt.Sprintf("No connection to close to node %v", id))
		} else if err := conn.Close(); err != nil {
			rrc.logger.Log(logging.LevelWarn, fmt.Sprintf("Could not close connection to node %v", id))
		}
	}
}

// Establishes a connection to a single node at address addrString.
func (rrc *RoundRobinClient) connectToNode(ctx context.Context, addrString string) (*grpc.ClientConn, requestreceiver.RequestReceiver_ListenClient, error) {

	rrc.logger.Log(logging.LevelDebug, fmt.Sprintf("Connecting to node: %s", addrString))

	// Set general gRPC dial options.
	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize), grpc.MaxCallSendMsgSize(maxMessageSize)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Set up a gRPC connection.
	conn, err := grpc.DialContext(ctx, addrString, dialOpts...)
	if err != nil {
		return nil, nil, err
	}

	// Register client stub.
	client := requestreceiver.NewRequestReceiverClient(conn)

	// Remotely invoke the Listen function on the other node's gRPC server.
	// As this is "stream of requests"-type RPC, it returns a message sink.
	msgSink, err := client.Listen(context.Background())
	if err != nil {
		if cerr := conn.Close(); cerr != nil {
			rrc.logger.Log(logging.LevelWarn, fmt.Sprintf("Failed to close connection: %v", cerr))
		}
		return nil, nil, err
	}

	// Return the message sink connected to the node.
	return conn, msgSink, nil
}
