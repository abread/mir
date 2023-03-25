// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	gonet "net"
	"os"
	"strconv"
	"time"

	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
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

	//"github.com/filecoin-project/mir/pkg/pb/batchfetcherpb"
	//"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/requestreceiver"
	"github.com/filecoin-project/mir/pkg/systems/trantor"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/libp2p"
)

const (
	ReqReceiverBasePort = 20000
)

var (
	protocol      string
	batchSize     int
	statFileName  string
	statPeriod    time.Duration
	enableOTLP    bool
	traceFileName string

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
	nodeCmd.Flags().StringVar(&traceFileName, "traceFile", "", "output file for *alea* traces")
	nodeCmd.Flags().BoolVar(&enableOTLP, "enableOTLP", false, "enable OTLP exporter")
}

func issSMRFactory(ctx context.Context, ownID t.NodeID, transport net.Transport, initialMembership map[t.NodeID]t.NodeAddress, smrParams trantor.Params, logger logging.Logger) (*trantor.System, error) {
	localCrypto := deploytest.NewLocalCryptoSystem("pseudo", membership.GetIDs(initialMembership), logger)

	genesisCheckpoint, err := trantor.GenesisCheckpoint([]byte{}, smrParams)
	if err != nil {
		return nil, fmt.Errorf("could not create genesis checkpoint: %w", err)
	}

	return trantor.NewISS(
		ctx,
		ownID,
		transport,
		genesisCheckpoint,
		localCrypto.Crypto(ownID),
		&App{Logger: logger, Membership: initialMembership},
		smrParams,
		logger,
	)
}

func aleaSMRFactory(ctx context.Context, ownID t.NodeID, transport net.Transport, initialMembership map[t.NodeID]t.NodeAddress, smrParams trantor.Params, logger logging.Logger) (*trantor.System, error) {
	F := (len(initialMembership) - 1) / 3
	localCrypto := deploytest.NewLocalThreshCryptoSystem("pseudo", membership.GetIDs(initialMembership), 2*F+1, logger)

	return trantor.NewAlea(
		ctx,
		ownID,
		transport,
		nil,
		localCrypto.ThreshCrypto(ownID),
		&App{Logger: logger, Membership: initialMembership},
		smrParams,
		logger,
	)
}

type smrFactory func(ctx context.Context, ownID t.NodeID, transport net.Transport, initialMembership map[t.NodeID]t.NodeAddress, smrParams trantor.Params, logger logging.Logger) (*trantor.System, error)

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

	if enableOTLP {
		otlpExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(otlptracehttp.WithInsecure()))
		if err != nil {
			return fmt.Errorf("error creating otlp exporter: %w", err)
		}
		defer func() {
			_ = otlpExporter.Shutdown(context.Background())
		}()

		tp := trace.NewTracerProvider(
			trace.WithBatcher(otlpExporter),
			trace.WithSampler(trace.AlwaysSample()),
		)
		defer func() {
			_ = tp.Shutdown(context.Background())
		}()
		otel.SetTracerProvider(tp)
	}

	// Load system membership.
	nodeAddrs, err := membership.FromFileName(membershipFile)
	if err != nil {
		return fmt.Errorf("could not load membership: %w", err)
	}
	initialMembership, err := membership.DummyMultiAddrs(nodeAddrs)
	if err != nil {
		return fmt.Errorf("could not create dummy multiaddrs: %w", err)
	}

	// Parse own ID.
	ownNumericID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("unable to convert node ID: %w", err)
	} else if ownNumericID < 0 || ownNumericID >= len(initialMembership) {
		return fmt.Errorf("ID must be in [0, %d]", len(initialMembership)-1)
	}
	ownID := t.NodeID(id)

	// Set Trantor parameters.
	smrParams := trantor.DefaultParams(initialMembership)
	smrParams.Mempool.MaxTransactionsInBatch = batchSize

	// derived to match mean alea latency with 1tx/s load (directed at the next leader replica)
	smrParams.AdjustSpeed(time.Duration(13483+7826*math.Log(float64(len(initialMembership)))) * time.Nanosecond)

	// Assemble listening address.
	// In this benchmark code, we always listen on tha address 0.0.0.0.
	portStr, err := getPortStr(initialMembership[ownID])
	if err != nil {
		return fmt.Errorf("could not parse port from own address: %w", err)
	}
	addrStr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", portStr)
	listenAddr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		return fmt.Errorf("could not create listen address: %w", err)
	}
	h, err := libp2p.NewDummyHostWithPrivKey(
		t.NodeAddress(libp2p.NewDummyMultiaddr(ownNumericID, listenAddr)),
		libp2p.NewDummyHostKey(ownNumericID),
	)
	if err != nil {
		return fmt.Errorf("failed to create libp2p host: %w", err)
	}

	// Initialize the libp2p transport subsystem.
	transport := libp2p2.NewTransport(smrParams.Net, ownID, h, logger)

	benchApp, err := smrFactories[protocol](ctx, ownID, transport, initialMembership, smrParams, logger)
	if err != nil {
		return fmt.Errorf("could not create bench app: %w", err)
	}

	/*recorder, err := eventlog.NewRecorder(
		ownID,
		"bench-output",
		logging.Decorate(logger, "EVTLOG: "),
		eventlog.EventFilterOpt(func(e *eventpb.Event) bool {
			switch e := e.Type.(type) {
			case *eventpb.Event_NewRequests:
				return true
			case *eventpb.Event_BatchFetcher:
				switch e.BatchFetcher.Type.(type) {
				case *batchfetcherpb.Event_NewOrderedBatch:
					return true
				}
			}
			return false
		}),
	)
	if err != nil {
		return fmt.Errorf("cannot create event recorder: %w", err)
	}*/

	var tracer eventlog.Interceptor = eventlog.NilInterceptor
	if traceFileName != "" {
		traceFile, err := os.Create(traceFileName)
		if err != nil {
			return fmt.Errorf("error creating trace output file: %w", err)
		}
		defer func() {
			_ = traceFile.Close()
		}()

		ownQueueIdx := slices.Index(smrParams.Alea.AllNodes(), ownID)
		aleaTracer := aleatracer.NewAleaTracer(ctx, aleatypes.QueueIdx(ownQueueIdx), len(initialMembership), traceFile)
		defer aleaTracer.Stop()
		tracer = aleaTracer
	}

	stat := stats.NewStats()
	interceptor := eventlog.MultiInterceptor(
		stats.NewStatInterceptor(stat, "app"),
		//recorder,
		tracer,
	)

	nodeConfig := mir.DefaultNodeConfig().WithLogger(logger)
	node, err := mir.NewNode(t.NodeID(id), nodeConfig, benchApp.Modules(), nil, interceptor)
	if err != nil {
		return fmt.Errorf("could not create node: %w", err)
	}

	reqReceiverListener, err := gonet.Listen("tcp", fmt.Sprintf(":%v", ReqReceiverBasePort+ownNumericID))
	if err != nil {
		return fmt.Errorf("could not create request receiver listener: %w", err)
	}

	reqReceiver := requestreceiver.NewRequestReceiver(node, "mempool", logger)
	reqReceiver.Start(reqReceiverListener)
	defer reqReceiver.Stop()

	if err := benchApp.Start(); err != nil {
		return fmt.Errorf("could not start bench app: %w", err)
	}
	defer benchApp.Stop()

	var statFile *os.File
	if statFileName != "" {
		statFile, err = os.Create(statFileName)
		if err != nil {
			return fmt.Errorf("could not open output file for statistics: %w", err)
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

	defer node.Stop()
	nodeErr := node.Run(ctx)
	if nodeErr != nil && errors.Is(nodeErr, mir.ErrStopped) {
		nodeErr = nil
	}
	return nodeErr
}

func getPortStr(address t.NodeAddress) (string, error) {
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
