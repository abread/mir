// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	gonet "net"
	"os"
	"runtime"
	"strconv"
	"time"

	es "github.com/go-errors/errors"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir"
	"github.com/filecoin-project/mir/cmd/bench/aleatracer"
	"github.com/filecoin-project/mir/cmd/bench/stats"
	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/deploytest"
	"github.com/filecoin-project/mir/pkg/eventlog"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/membership"
	"github.com/filecoin-project/mir/pkg/net"
	libp2p2 "github.com/filecoin-project/mir/pkg/net/libp2p"
	"github.com/filecoin-project/mir/pkg/pb/trantorpb"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	"github.com/filecoin-project/mir/pkg/transactionreceiver"
	"github.com/filecoin-project/mir/pkg/trantor"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/libp2p"
)

const (
	TxReceiverBasePort = 20000
)

var (
	protocol       string
	batchSize      int
	statFileName   string
	statPeriod     time.Duration
	traceFileName  string
	cryptoImplType string

	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "Start a Mir node",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runNode(ctx)
		},
	}
)

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.Flags().StringVarP(&protocol, "protocol", "p", "iss", "protocol to use")
	nodeCmd.Flags().IntVarP(&batchSize, "batchSize", "b", 1024, "maximum number of transactions in a batch (mempool module)")
	nodeCmd.Flags().StringVarP(&statFileName, "statFile", "o", "", "output file for statistics")
	nodeCmd.Flags().DurationVar(&statPeriod, "statPeriod", time.Second, "statistic record period")
	nodeCmd.Flags().StringVar(&traceFileName, "traceFile", "", "output file for alea tracing")
	nodeCmd.Flags().StringVarP(&cryptoImplType, "cryptoImplType", "c", "pseudo", "type of cryptography to use (acceptable values: pseudo or dummy)")
}

func issSMRFactory(app *App, ownID t.NodeID, transport net.Transport, initialMembership *trantorpbtypes.Membership, smrParams trantor.Params, logger logging.Logger) (*trantor.System, error) {
	localCS, err := deploytest.NewLocalCryptoSystem(cryptoImplType, membership.GetIDs(initialMembership), logger)
	if err != nil {
		return nil, es.Errorf("could not create a local crypto system: %w", err)
	}
	localCrypto, err := localCS.Crypto(ownID)
	if err != nil {
		return nil, es.Errorf("could not create a local crypto system: %w", err)
	}

	genesisCheckpoint, err := trantor.GenesisCheckpoint([]byte{}, smrParams)
	if err != nil {
		return nil, es.Errorf("could not create genesis checkpoint: %w", err)
	}

	return trantor.NewISS(
		ownID,
		transport,
		genesisCheckpoint,
		localCrypto,
		app,
		smrParams,
		logger,
	)
}

func aleaSMRFactory(app *App, ownID t.NodeID, transport net.Transport, initialMembership *trantorpbtypes.Membership, smrParams trantor.Params, logger logging.Logger) (*trantor.System, error) {
	F := (len(initialMembership.Nodes) - 1) / 3
	localCS := deploytest.NewLocalThreshCryptoSystem(cryptoImplType, membership.GetIDs(initialMembership), 2*F+1)
	localCrypto, err := localCS.ThreshCrypto(ownID)
	if err != nil {
		return nil, es.Errorf("could not create a local threshcrypto system: %w", err)
	}

	return trantor.NewAlea(
		ownID,
		transport,
		nil,
		localCrypto,
		app,
		smrParams,
		logger,
	)
}

type smrFactory func(app *App, ownID t.NodeID, transport net.Transport, initialMembership *trantorpbtypes.Membership, smrParams trantor.Params, logger logging.Logger) (*trantor.System, error)

var smrFactories = map[string]smrFactory{
	"iss":  issSMRFactory,
	"alea": aleaSMRFactory,
}

func runNode(ctx context.Context) error {
	var logger logging.Logger
	if verbose {
		logger = logging.ConsoleDebugLogger
	} else {
		logger = logging.ConsoleWarnLogger
	}

	// Load system membership.
	nodeAddrs, err := membership.FromFileName(membershipFile)
	if err != nil {
		return es.Errorf("could not load membership: %w", err)
	}
	initialMembership, err := membership.DummyMultiAddrs(nodeAddrs)
	if err != nil {
		return es.Errorf("could not create dummy multiaddrs: %w", err)
	}

	// Parse own ID.
	ownNumericID, err := strconv.Atoi(id)
	if err != nil {
		return es.Errorf("unable to convert node ID: %w", err)
	} else if ownNumericID < 0 || ownNumericID >= len(initialMembership.Nodes) {
		return es.Errorf("ID must be in [0, %d]", len(initialMembership.Nodes)-1)
	}
	ownID := t.NodeID(id)

	// Set Trantor parameters.
	smrParams := trantor.DefaultParams(initialMembership)
	smrParams.Mempool.MaxTransactionsInBatch = batchSize

	// ensure network messages can accommodate the chosen batch size
	batchAdjustedMaxMsgSize := batchSize * 512 * 105 / 100
	if smrParams.Net.MaxMessageSize < batchAdjustedMaxMsgSize {
		smrParams.Net.MaxMessageSize = batchAdjustedMaxMsgSize
	}

	// derived to match mean alea latency with 1tx/s load (directed at the next leader replica) - 2ms
	//smrParams.AdjustSpeed(time.Duration(10879+1467*len(initialMembership.Nodes)) * time.Microsecond)

	// Assemble listening address.
	// In this benchmark code, we always listen on tha address 0.0.0.0.
	portStr, err := getPortStr(initialMembership.Nodes[ownID].Addr)
	if err != nil {
		return es.Errorf("could not parse port from own address: %w", err)
	}
	h, err := libp2p.NewDummyHostWithPrivKey(
		portStr,
		libp2p.NewDummyHostKey(ownNumericID),
	)
	if err != nil {
		return es.Errorf("failed to create libp2p host: %w", err)
	}

	// Initialize the libp2p transport subsystem.
	transport := libp2p2.NewTransport(smrParams.Net, ownID, h, logger)

	app := &App{Logger: logger, Membership: initialMembership}

	benchSystem, err := smrFactories[protocol](app, ownID, transport, initialMembership, smrParams, logger)
	if err != nil {
		return es.Errorf("could not create bench app: %w", err)
	}

	/*recorder, err := eventlog.NewRecorder(
		ownID,
		"bench-output",
		logging.Decorate(logger, "EVTLOG: "),
		eventlog.EventFilterOpt(func(e *eventpb.Event) bool {
			switch e := e.Type.(type) {
			case *eventpb.Event_Mempool:
				switch e.Mempool.Type.(type) {
				case *mempoolpb.Event_NewTransactions:
					return true
				}
			case *eventpb.Event_BatchFetcher:
				switch e.BatchFetcher.Type.(type) {
				case *batchfetcherpb.Event_NewOrderedBatch:
					return true
				}
			}
			return false
		}),
	)
	defer recorder.Stop()
	if err != nil {
		return es.Errorf("cannot create event recorder: %w", err)
	}*/

	var tracer eventlog.Interceptor = eventlog.NilInterceptor
	if traceFileName != "" {
		traceFile, err := os.Create(traceFileName)
		if err != nil {
			return es.Errorf("error creating trace output file: %w", err)
		}
		defer func() {
			_ = traceFile.Close()
		}()

		ownQueueIdx := slices.Index(smrParams.Alea.AllNodes(), ownID)
		aleaTracer := aleatracer.NewAleaTracer(ctx, aleatypes.QueueIdx(ownQueueIdx), len(initialMembership.Nodes), traceFile)
		defer aleaTracer.Stop()
		tracer = aleaTracer
	}

	var statTracer eventlog.Interceptor = eventlog.NilInterceptor
	if statFileName != "/dev/null" {
		stat := stats.NewStats()
		statTracer = stats.NewStatInterceptor(stat, "app")

		var statFile *os.File
		if statFileName != "" {
			statFile, err = os.Create(statFileName)
			if err != nil {
				return es.Errorf("could not open output file for statistics: %w", err)
			}
		} else {
			statFile = os.Stdout
		}

		statCSV := csv.NewWriter(statFile)
		stat.WriteCSVHeader(statCSV)

		go func() {
			timestamp := time.Now()
			for {
				ticker := time.NewTicker(statPeriod)
				defer ticker.Stop()

				select {
				case <-ctx.Done():
					return
				case ts := <-ticker.C:
					d := ts.Sub(timestamp)
					stat.WriteCSVRecordAndReset(statCSV, d)
					statCSV.Flush()
					timestamp = ts
				}
			}
		}()
	}

	interceptor := eventlog.MultiInterceptor(
		//recorder,
		statTracer,
		tracer,
	)

	nodeConfig := mir.DefaultNodeConfig().WithLogger(logger)
	nodeModules := benchSystem.Modules().ConvertConcurrentEventAppliersToGoroutinePools(ctx, runtime.NumCPU())
	node, err := mir.NewNode(t.NodeID(id), nodeConfig, nodeModules, interceptor)
	if err != nil {
		return es.Errorf("could not create node: %w", err)
	}

	txReceiverListener, err := gonet.Listen("tcp", fmt.Sprintf(":%v", TxReceiverBasePort+ownNumericID))
	if err != nil {
		return es.Errorf("could not create tx receiver listener: %w", err)
	}

	txReceiver := transactionreceiver.NewTransactionReceiver(node, "mempool", logger)
	app.DeliverHook = func(txs []*trantorpbtypes.Transaction) error {
		txsPb := make([]*trantorpb.Transaction, len(txs))
		for i, tx := range txs {
			txsPb[i] = tx.Pb()
		}
		txReceiver.NotifyBatchDeliver(txsPb)
		return nil
	}
	txReceiver.Start(txReceiverListener)
	defer txReceiver.Stop()

	if err := benchSystem.Start(); err != nil {
		return es.Errorf("could not start bench app: %w", err)
	}
	defer benchSystem.Stop()

	defer node.Stop()
	nodeErr := node.Run(ctx)
	if nodeErr != nil && errors.Is(nodeErr, mir.ErrStopped) {
		nodeErr = nil
	}
	return nodeErr
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
