// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	gonet "net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	es "github.com/go-errors/errors"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir"
	"github.com/filecoin-project/mir/cmd/bench/aleatracer"
	"github.com/filecoin-project/mir/cmd/bench/localtxgenerator"
	"github.com/filecoin-project/mir/cmd/bench/stats"
	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/deploytest"
	"github.com/filecoin-project/mir/pkg/eventlog"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/membership"
	libp2p2 "github.com/filecoin-project/mir/pkg/net/libp2p"
	batchfetcherpbtypes "github.com/filecoin-project/mir/pkg/pb/batchfetcherpb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	"github.com/filecoin-project/mir/pkg/pb/trantorpb"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/rendezvous"
	"github.com/filecoin-project/mir/pkg/transactionreceiver"
	"github.com/filecoin-project/mir/pkg/trantor"
	"github.com/filecoin-project/mir/pkg/trantor/appmodule"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/libp2p"
	"github.com/filecoin-project/mir/pkg/util/maputil"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

const (
	syncLimit        = 30 * time.Second
	syncPollInterval = 100 * time.Millisecond
)

var (
	configFileName      string
	statSummaryFileName string
	liveStatsFileName   string
	clientStatsFileName string
	netStatsFileName    string
	statPeriod          time.Duration
	traceFileName       string
	readySyncFileName   string
	deliverSyncFileName string

	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "Start a Mir node",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if err := runNode(ctx); !es.Is(err, mir.ErrStopped) {
				return err
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(nodeCmd)

	// Required arguments
	nodeCmd.Flags().StringVarP(&configFileName, "config-file", "c", "", "configuration file")
	_ = nodeCmd.MarkFlagRequired("config-file")
	nodeCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "node ID")
	_ = nodeCmd.MarkPersistentFlagRequired("id")

	// Optional arguments
	nodeCmd.Flags().DurationVar(&statPeriod, "stat-period", 1*time.Second, "statistic record period")
	nodeCmd.Flags().StringVar(&clientStatsFileName, "client-stat-file", "", "live cumulative client statistics output file")
	nodeCmd.Flags().StringVar(&netStatsFileName, "net-stat-file", "", "live cumulative net statistics output file")
	nodeCmd.Flags().StringVar(&liveStatsFileName, "replica-stat-file", "", "output file for live statistics, default is standard output")
	nodeCmd.Flags().StringVar(&statSummaryFileName, "summary-stat-file", "", "output file for summarized statistics")
	nodeCmd.Flags().StringVar(&traceFileName, "traceFile", "", "output file for alea tracing")

	// Sync files
	nodeCmd.Flags().StringVar(&readySyncFileName, "ready-sync-file", "", "file to use for initial synchronization when ready to start the benchmark")
	nodeCmd.Flags().StringVar(&deliverSyncFileName, "deliver-sync-file", "", "file to use for synchronization when waiting to deliver all transactions")
}

func runNode(ctx context.Context) error {
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

	// Parse own ID.
	ownNumericID, err := strconv.Atoi(id)
	if err != nil {
		return es.Errorf("unable to convert node ID: %w", err)
	}

	// Check if own id is in the membership
	initialMembership := params.Trantor.Iss.InitialMembership
	if _, ok := initialMembership.Nodes[t.NodeID(id)]; !ok {
		return es.Errorf("own ID (%v) not found in membership (%v)", id, maputil.GetKeys(initialMembership.Nodes))
	}
	ownID := t.NodeID(id)

	// Assemble listening address.
	// In this benchmark code, we always listen on the address 0.0.0.0.
	portStr, err := getPortStr(initialMembership.Nodes[ownID].Addr)
	if err != nil {
		return es.Errorf("could not parse port from own address: %w", err)
	}

	// Create libp2p host (supporting all default protocols: TCP, QUIC, and WebTransport over IPv4 and IPv6).
	h, err := libp2p.NewDummyHostWithPrivKey(
		portStr,
		libp2p.NewDummyHostKey(ownNumericID),
	)
	if err != nil {
		return es.Errorf("failed to create libp2p host: %w", err)
	}

	// Initialize tracking of networking statistics.
	netStats := stats.NewNetStats(time.Second)

	// Initialize the libp2p transport subsystem.
	transport := libp2p2.NewTransport(params.Trantor.Net, ownID, h, logger, netStats)

	// Instantiate the crypto module.
	logger.Log(logging.LevelWarn, "generating crypto keys...")
	localCryptoSystem, err := deploytest.NewLocalCryptoSystem(params.CryptoImpl, params.CryptoSeed, membership.GetIDs(initialMembership), logger)
	if err != nil {
		return es.Errorf("could not create a local crypto system: %w", err)
	}
	localCrypto, err := localCryptoSystem.Crypto(ownID)
	if err != nil {
		return es.Errorf("could not create a local crypto module: %w", err)
	}

	// Instantiate the threshold crypto module.
	logger.Log(logging.LevelWarn, "generating threshcrypto keys...")
	F := (len(initialMembership.Nodes) - 1) / 3
	thresh := 2*F + 1
	localThreshCryptoSystem, err := deploytest.NewLocalThreshCryptoSystem(params.ThreshCryptoImpl, params.CryptoSeed, membership.GetIDs(initialMembership), thresh)
	if err != nil {
		return es.Errorf("could not create threshcrypto system: %w", err)
	}
	localThreshCrypto, err := localThreshCryptoSystem.ThreshCrypto(ownID)
	if err != nil {
		return es.Errorf("could not create a local crypto module: %w", err)
	}

	// Generate the initial checkpoint.
	genesisCheckpoint, err := trantor.GenesisCheckpoint([]byte{}, params.Trantor)
	if err != nil {
		return es.Errorf("could not create genesis checkpoint: %w", err)
	}

	// Create a local transaction generator.
	// It has, at the same time, the interface of a trantor App,
	// so it knows when transactions are delivered and can submit new ones accordingly.
	// If the client ID is not specified, use the local node's ID
	if params.TxGen.ClientID == "" {
		params.TxGen.ClientID = tt.ClientID(ownID)
	}
	txGen := localtxgenerator.New(localtxgenerator.DefaultModuleConfig(), params.TxGen)

	// Create a Trantor instance.
	logger.Log(logging.LevelWarn, "creating trantor instance")
	trantorInstance, err := trantor.New(
		ownID,
		transport,
		genesisCheckpoint,
		localCrypto,
		localThreshCrypto,
		appmodule.AppLogicFromStatic(txGen, initialMembership), // The transaction generator is also a static app.
		params.Trantor,
		logger,
	)
	if err != nil {
		return es.Errorf("could not create bench app: %w", err)
	}

	var tracer eventlog.Interceptor = eventlog.NilInterceptor
	if traceFileName != "" {
		traceFile, err := os.Create(traceFileName)
		if err != nil {
			return es.Errorf("error creating trace output file: %w", err)
		}
		defer func() {
			_ = traceFile.Close()
		}()

		ownQueueIdx := slices.Index(params.Trantor.Alea.AllNodes(), ownID)
		aleaTracer := aleatracer.NewAleaTracer(ctx, aleatypes.QueueIdx(ownQueueIdx), len(initialMembership.Nodes), traceFile)
		defer aleaTracer.Stop()
		tracer = aleaTracer
	}

	// Add transaction generator module to the setup.
	trantorInstance.WithModule("localtxgen", txGen)

	// Create trackers for gathering statistics about the performance.
	liveStats := stats.NewLiveStats(aleatypes.QueueIdx(slices.Index(params.Trantor.Alea.AllNodes(), ownID)))
	clientStats := stats.NewClientStats(time.Millisecond, 5*time.Second, 4*params.Trantor.Mempool.MaxTransactionsInBatch)
	txGen.TrackStats(liveStats)
	txGen.TrackStats(clientStats)

	txReceiverInterceptor := &txReceiverInterceptor{AppModuleID: trantor.DefaultModuleConfig().App}
	interceptor := eventlog.MultiInterceptor(
		tracer,
		stats.NewStatInterceptor(liveStats, trantor.DefaultModuleConfig().App, string(params.TxGen.ClientID)+"."),
		txReceiverInterceptor,
	)
	// Instantiate the Mir Node.
	nodeConfig := mir.DefaultNodeConfig().WithLogger(logger)
	nodeConfig.Stats.Period = 5 * time.Second
	nodeModules := trantorInstance.Modules().ConvertConcurrentEventAppliersToGoroutinePools(ctx, runtime.NumCPU())
	node, err := mir.NewNode(t.NodeID(id), nodeConfig, nodeModules, interceptor)
	if err != nil {
		return es.Errorf("could not create node: %w", err)
	}

	numericID, err := strconv.Atoi(id)
	if err != nil {
		return es.Errorf("only numeric node ids are supported")
	}
	txRecvListener, err := gonet.Listen("tcp", fmt.Sprintf(":%d", TxReceiverBasePort+numericID))
	if err != nil {
		return es.Errorf("could not create tx listener: %w")
	}

	txReceiver := transactionreceiver.NewTransactionReceiver(node, trantor.DefaultModuleConfig().Mempool, logger)
	txReceiverInterceptor.TxReceiver = txReceiver

	logger.Log(logging.LevelWarn, "starting trantor instance")
	if err := trantorInstance.Start(); err != nil {
		return es.Errorf("could not start bench app: %w", err)
	}

	logger.Log(logging.LevelWarn, "waiting for all nodes to connect")
	if err := transport.WaitFor(len(initialMembership.Nodes)); err != nil {
		return es.Errorf("failed waiting for network connections: %w", err)
	}
	time.Sleep(3 * time.Second) // wait a bit for the network to settle everywhere

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

	// Output the statistics.
	var statFile *os.File
	if liveStatsFileName != "" {
		statFile, err = os.Create(liveStatsFileName)
		if err != nil {
			return es.Errorf("could not open output file for statistics: %w", err)
		}
	} else {
		statFile = os.Stdout
	}

	trantorStopped := make(chan struct{})
	statsWg := &sync.WaitGroup{}
	statsCtx, stopStats := context.WithCancel(ctx)
	defer stopStats()

	replicaStatsCSV := csv.NewWriter(statFile)
	goDisplayLiveStats(statsCtx, statsWg, trantorStopped, liveStats, replicaStatsCSV)

	if clientStatsFileName != "" {
		clientStatFile, err := os.Create(clientStatsFileName)
		if err != nil {
			return es.Errorf("could not open output file for client statistics: %w", err)
		}
		clientStatsCSV := csv.NewWriter(clientStatFile)
		goDisplayLiveStats(statsCtx, statsWg, trantorStopped, clientStats, clientStatsCSV)
	}

	if netStatsFileName != "" {
		netStatFile, err := os.Create(netStatsFileName)
		if err != nil {
			return es.Errorf("could not open output file for net statistics: %w", err)
		}
		netStatsCSV := csv.NewWriter(netStatFile)
		goDisplayLiveStats(statsCtx, statsWg, trantorStopped, netStats, netStatsCSV)
	}

	// Stop outputting real-time stats and submitting transactions,
	// wait until everything is delivered, and stop node.
	shutDown := func() {

		stopStats()
		txGen.Stop()
		txReceiver.Stop()

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

		// Stop Mir node and Trantor instance.
		logger.Log(logging.LevelWarn, "Stopping Mir node.")
		node.Stop()
		logger.Log(logging.LevelWarn, "Mir node stopped.")
		logger.Log(logging.LevelWarn, "Stopping Trantor.")
		trantorInstance.Stop()
		logger.Log(logging.LevelWarn, "Trantor stopped.")

		close(trantorStopped)
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
		// Setup signal notification channel
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

		go func() {
			// Wait for a closing signal and shut down the node.
			select {
			case <-ctx.Done():
			case <-sigs:
			}
			shutDown()
			close(done)
		}()
	}

	if params.CrashAfter > 0 {
		go func() {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(params.CrashAfter)):
				logger.Log(logging.LevelWarn, "Simulating node crash")
				os.Exit(0)
			}
		}()
	}

	// Start accepting transactions from the outside world
	txReceiver.Start(txRecvListener)

	// Start generating the load and measuring performance.
	logger.Log(logging.LevelWarn, "applying load")
	clientStats.Start()
	netStats.Start()
	txGen.Start()

	nodeError := node.Run(ctx)
	<-done
	writeFinalStats(clientStats, netStats, statSummaryFileName, logger)
	return nodeError
}

func writeFinalStats(
	clientStats *stats.ClientStats,
	netStats *stats.NetStats,
	statFileName string,
	logger logging.Logger,
) {
	if statFileName == "" {
		return
	}

	data := make(map[string]any)
	data["Net"] = netStats
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

func goDisplayLiveStats(ctx context.Context, statsWg *sync.WaitGroup, trantorStopped chan struct{}, statsProducer stats.Stats, statCSVConsumer *csv.Writer) {
	statsWg.Add(1)
	statsProducer.WriteCSVHeader(statCSVConsumer)
	statCSVConsumer.Flush()

	go func() {
		defer statsWg.Done()

		timestamp := time.Now()
		ticker := time.NewTicker(statPeriod)
		defer ticker.Stop()

	StatOutputLoop:
		for {
			select {
			case <-ctx.Done():
				break StatOutputLoop
			case ts := <-ticker.C:
				d := ts.Sub(timestamp)
				statsProducer.WriteCSVRecord(statCSVConsumer, d)
				statCSVConsumer.Flush()
				timestamp = ts
			}
		}

		// wait for trantor to fully stop
		<-trantorStopped

		statsProducer.Fill()
		d := time.Since(timestamp)
		statsProducer.WriteCSVRecord(statCSVConsumer, d)
		statCSVConsumer.Flush()
	}()
}

func getPortStr(addressStr string) (string, error) {
	address, err := multiaddr.NewMultiaddr(addressStr)
	if err != nil {
		return "", err
	}

	_, addrStr, err := manet.DialArgs(address)
	if err != nil {
		return "", err
	}

	_, portStr, err := gonet.SplitHostPort(addrStr)
	if err != nil {
		return "", err
	}

	return portStr, nil
}

type txReceiverInterceptor struct {
	AppModuleID t.ModuleID
	TxReceiver  *transactionreceiver.TransactionReceiver
}

func (i *txReceiverInterceptor) Intercept(evs events.EventList) error {
	for _, ev := range evs.Slice() {
		if bfEvW, ok := ev.Type.(*eventpbtypes.Event_BatchFetcher); ok {
			if newBatchEvW, ok := bfEvW.BatchFetcher.Type.(*batchfetcherpbtypes.Event_NewOrderedBatch); ok {
				i.TxReceiver.NotifyBatchDeliver(
					sliceutil.Transform(
						newBatchEvW.NewOrderedBatch.Txs,
						func(_ int, tx *trantorpbtypes.Transaction) *trantorpb.Transaction { return tx.Pb() },
					),
				)
			}
		}
	}

	return nil
}
