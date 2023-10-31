// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"crypto"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	es "github.com/go-errors/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/filecoin-project/mir/cmd/bench/localtxgenerator"
	"github.com/filecoin-project/mir/cmd/bench/stats"
	"github.com/filecoin-project/mir/pkg/dummyclient"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/membership"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	mempoolpbtypes "github.com/filecoin-project/mir/pkg/pb/mempoolpb/types"
	"github.com/filecoin-project/mir/pkg/pb/trantorpb"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/rendezvous"
	"github.com/filecoin-project/mir/pkg/rrclient"
	"github.com/filecoin-project/mir/pkg/transactionreceiver"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

const (
	TxReceiverBasePort = 20000
)

var (
	clientType string

	latClientCmd = &cobra.Command{
		Use:   "latclient",
		Short: "Generate and submit transactions to a Mir cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runLatClient(ctx)
		},
	}
)

func init() {
	rootCmd.AddCommand(latClientCmd)

	// Required arguments
	latClientCmd.Flags().StringVarP(&configFileName, "config-file", "c", "", "configuration file")
	_ = latClientCmd.MarkFlagRequired("config-file")
	latClientCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "client ID")
	_ = latClientCmd.MarkPersistentFlagRequired("id")

	// Optional arguments
	latClientCmd.Flags().StringVarP(&clientType, "type", "t", "dummy", "client type (one of: dummy, rr)")

	latClientCmd.Flags().DurationVar(&statPeriod, "stat-period", 5*time.Second, "statistic record period")
	latClientCmd.Flags().StringVar(&clientStatsFileName, "live-stat-file", "", "output file for live statistics, default is standard output")
	latClientCmd.Flags().StringVar(&statSummaryFileName, "summary-stat-file", "", "output file for summarized statistics")

	// Sync files
	latClientCmd.Flags().StringVar(&readySyncFileName, "ready-sync-file", "", "file to use for initial synchronization when ready to start the benchmark")
	latClientCmd.Flags().StringVar(&deliverSyncFileName, "deliver-sync-file", "", "file to use for synchronization when waiting to deliver all transactions")
}

type mirClient interface {
	Connect(ctx context.Context, membership map[t.NodeID]string)
	SubmitTransaction(data []byte) error
	Disconnect()
}

type mirClientFactory func(clientID tt.ClientID, hasher crypto.Hash, logger logging.Logger) mirClient

var clientFactories = map[string]mirClientFactory{
	"dummy": func(clientID tt.ClientID, hasher crypto.Hash, logger logging.Logger) mirClient {
		return dummyclient.NewDummyClient(clientID, hasher, logger)
	},
	"rr": func(clientID tt.ClientID, hasher crypto.Hash, logger logging.Logger) mirClient {
		return rrclient.NewRoundRobinClient(clientID, hasher, logger)
	},
}

func runLatClient(ctx context.Context) error {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	var logger logging.Logger
	if verbose {
		logger = logging.ConsoleDebugLogger
	} else {
		logger = logging.ConsoleWarnLogger
	}

	// Load configuration parameters
	var params BenchParams
	if err := loadFromFile(configFileName, &params); err != nil {
		return es.Errorf("could not load parameters from file '%s': %w", configFileName, err)
	}

	// Check if own id is in the membership
	initialMembership := params.Trantor.Iss.InitialMembership
	if _, ok := initialMembership.Nodes[t.NodeID(id)]; !ok {
		return es.Errorf("own ID (%v) not found in membership (%v)", id, maputil.GetKeys(initialMembership.Nodes))
	}
	ownID := t.NodeID(id)

	addresses, err := membership.GetIPs(initialMembership)
	if err != nil {
		return es.Errorf("could not load node IPs: %w", err)
	}

	// Generate addresses and ports for transaction receivers.
	// Each node uses different ports for receiving protocol messages and transactions.
	// These addresses will be used by the client code to know where to send its transactions.
	txReceiverAddrs := make(map[t.NodeID]string)
	for nodeID, nodeIP := range addresses {
		numericID, err := strconv.Atoi(string(nodeID))
		if err != nil {
			return es.Errorf("node IDs must be numeric in the sample app: %w", err)
		}
		txReceiverAddrs[nodeID] = net.JoinHostPort(nodeIP, fmt.Sprintf("%d", TxReceiverBasePort+numericID))

		// The break statement causes the client to send its transactions to only one single node.
		// Remove it for the client to send its transactions to all nodes.
		// TODO: Make this properly configurable and remove the hack.
		// break // sending to all nodes is more reproducible
	}
	logger.Log(logging.LevelInfo, "processed node list", "txReceiverAddrs", txReceiverAddrs)

	client := clientFactories[clientType](
		tt.ClientID(id+".0"),
		crypto.SHA256,
		logger,
	)

	// Synchronize with other nodes if necessary.
	// If invoked, this code blocks until all the nodes have connected to each other.
	// (The file created by Ready must be deleted by some external code (or manually) after all nodes have created it.)
	if readySyncFileName != "" {
		syncCtx, cancelFunc := context.WithTimeout(ctx, syncLimit)
		err = rendezvous.NewFileSyncer(readySyncFileName, syncPollInterval).Ready(syncCtx)
		cancelFunc()
		if err != nil {
			return fmt.Errorf("error synchronizing nodes: %w", err)
		}
	}

	client.Connect(ctx, txReceiverAddrs)
	defer client.Disconnect()
	logger.Log(logging.LevelInfo, "client connected")
	time.Sleep(3)

	// Create a local transaction generator.
	// It has, at the same time, the interface of a trantor App,
	// so it knows when transactions are delivered and can submit new ones accordingly.
	// We will abuse it to send transactions over the network and track statistics
	params.TxGen.ClientID = tt.ClientID(ownID)
	params.TxGen.NumClients = 1
	txGen := localtxgenerator.New(localtxgenerator.DefaultModuleConfig(), params.TxGen)

	// Create trackers for gathering statistics about the performance.
	clientStats := stats.NewClientStats(time.Millisecond, 5*time.Second, 4*params.Trantor.Mempool.MaxTransactionsInBatch)
	txGen.TrackStats(clientStats)

	// Output the statistics.
	var statFile *os.File
	if clientStatsFileName != "" {
		statFile, err = os.Create(clientStatsFileName)
		if err != nil {
			return es.Errorf("could not open output file for statistics: %w", err)
		}
	} else {
		statFile = os.Stdout
	}

	clientStopped := make(chan struct{})
	statsWg := &sync.WaitGroup{}
	statsCtx, stopStats := context.WithCancel(ctx)
	defer stopStats()

	liveStatsCSV := csv.NewWriter(statFile)
	goDisplayLiveStats(statsCtx, statsWg, clientStopped, clientStats, liveStatsCSV)

	clientAdapterWg := &sync.WaitGroup{}
	clientAdapterStop := make(chan struct{})

	// Stop outputting real-time stats and submitting transactions,
	// wait until everything is delivered, and stop node.
	shutDown := func() {
		stopStats()
		txGen.Stop()

		// Wait for other nodes to deliver their transactions.
		if deliverSyncFileName != "" {
			syncerCtx, stopWaiting := context.WithTimeout(ctx, syncLimit)
			err := rendezvous.NewFileSyncer(deliverSyncFileName, syncPollInterval).Ready(syncerCtx)
			stopWaiting()
			if err != nil {
				logger.Log(logging.LevelError, "Aborting waiting for other nodes transaction delivery.", "error", err)
			} else {
				logger.Log(logging.LevelWarn, "All nodes successfully delivered all transactions they submitted.",
					"error", err)
			}
		}

		close(clientAdapterStop)
		clientAdapterWg.Wait()

		close(clientStopped)
		statsWg.Wait()
	}

	done := make(chan struct{})
	if params.Duration > 0 {
		go func() {
			// Wait until the end of the benchmark and shut down the node.
			select {
			case <-ctx.Done():
			case <-time.After(time.Duration(params.Duration)):
			}
			shutDown()
			close(done)
		}()
	} else {
		// TODO: This is not right. Only have this branch to quit on node error.
		//   Set up signal handlers so that the nodes stops and cleans up after itself upon SIGINT and / or SIGTERM.
		close(done)
	}

	// submit transactions through client
	clientAdapterWg.Add(1)
	go func() {
		defer clientAdapterWg.Done()

		for {
			select {
			case <-clientAdapterStop:
				return
			case evs := <-txGen.EventsOut():
				for _, ev := range evs.Slice() {
					txs := ev.Type.(*eventpbtypes.Event_Mempool).Mempool.Type.(*mempoolpbtypes.Event_NewTransactions).NewTransactions.Transactions

					for _, tx := range txs {
						if err := client.SubmitTransaction(tx.Data); err != nil {
							if errors.Is(err, io.EOF) {
								select {
								case <-clientAdapterStop:
									return
								default:
									panic(err)
								}

							}
							panic(err)
						}
					}
				}
			}
		}
	}()

	// assume last node is healthy and can tell us which txs were properly received
	go registerDeliveredTxs(ctx, txReceiverAddrs[t.NewNodeIDFromInt(len(txReceiverAddrs)-1)], txGen)

	// Start generating the load and measuring performance.
	logger.Log(logging.LevelWarn, "applying load")
	clientStats.Start()
	txGen.Start()

	<-done
	writeClientFinalStats(clientStats, statSummaryFileName, logger)
	return nil
}

func writeClientFinalStats(
	clientStats *stats.ClientStats,
	statFileName string,
	logger logging.Logger,
) {
	if statFileName == "" {
		return
	}

	data := make(map[string]any)
	data["Client"] = clientStats

	statsData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logger.Log(logging.LevelError, "Could not marshal benchmark output", "error", err)
	}
	if err = os.WriteFile(statFileName, []byte(fmt.Sprintf("%s\n", string(statsData))), 0644); err != nil {
		logger.Log(logging.LevelError, "Could not write benchmark output to file",
			"file", statFileName, "error", err)
	}
}

func registerDeliveredTxs(ctx context.Context, txReceiverAddr string, txGen *localtxgenerator.LocalTXGen) {
	maxMessageSize := 1073741824
	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize), grpc.MaxCallSendMsgSize(maxMessageSize)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Set up a gRPC connection.
	conn, err := grpc.DialContext(ctx, txReceiverAddr, dialOpts...)
	if err != nil {
		panic(err)
	}

	// Register client stub.
	client := transactionreceiver.NewTransactionReceiverClient(conn)

	outputStream, err := client.Output(ctx, &transactionreceiver.Empty{})
	if err != nil {
		panic(err)
	}

	for {
		batch, err := outputStream.Recv()
		if errors.Is(err, io.EOF) {
			return
		} else if err != nil {
			fmt.Printf("error reading replica: %s", es.Wrap(err, 0).ErrorStack())
			return
		}

		if err := txGen.ApplyTXs(sliceutil.Transform(batch.Txs, func(_ int, tx *trantorpb.Transaction) *trantorpbtypes.Transaction {
			return trantorpbtypes.TransactionFromPb(tx)
		})); err != nil {
			panic(err)
		}
	}

}
