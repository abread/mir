// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"crypto"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	es "github.com/go-errors/errors"
	"github.com/spf13/cobra"
	rateLimiter "golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/filecoin-project/mir/pkg/dummyclient"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/membership"
	"github.com/filecoin-project/mir/pkg/rrclient"
	"github.com/filecoin-project/mir/pkg/transactionreceiver"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

var (
	txSize     int
	rate       float64
	burst      int
	duration   time.Duration
	clientType string
	//statFileName string
	//statPeriod   time.Duration

	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Generate and submit transactions to a Mir cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runClient(ctx)
		},
	}
)

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().IntVarP(&txSize, "txSize", "s", 256, "size of each transaction in bytes")
	clientCmd.Flags().Float64VarP(&rate, "rate", "r", 1000, "average number of transactions per second")
	clientCmd.Flags().IntVarP(&burst, "burst", "b", 1, "maximum number of transactions in a burst")
	clientCmd.Flags().DurationVarP(&duration, "duration", "T", 10*time.Second, "benchmarking duration")
	clientCmd.Flags().StringVarP(&clientType, "type", "t", "dummy", "client type (one of: dummy, rr)")
	clientCmd.Flags().StringVarP(&statFileName, "statFile", "o", "", "output file for statistics")
	clientCmd.Flags().DurationVar(&statPeriod, "statPeriod", time.Second, "statistic record period")
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

func runClient(ctx context.Context) error {
	var logger logging.Logger
	if verbose {
		logger = logging.ConsoleDebugLogger
	} else {
		logger = logging.ConsoleWarnLogger
	}

	initialMembership, err := membership.FromFileName(membershipFile)
	if err != nil {
		return es.Errorf("could not load membership: %w", err)
	}
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

	ctxStats, stopStats := context.WithCancel(ctx)
	defer stopStats()
	ctx, stop := context.WithCancel(ctx)

	client := clientFactories[clientType](
		tt.ClientID(id),
		crypto.SHA256,
		logger,
	)
	client.Connect(ctx, txReceiverAddrs)
	defer client.Disconnect()

	clock := time.Now()
	txs := make(map[tt.TxNo]time.Duration, int(rate*statPeriod.Seconds()))
	txsMutex := &sync.Mutex{}
	nextTxNo := tt.TxNo(0)
	if statFileName != "" {
		go clientStats(ctxStats, txReceiverAddrs, clock, txs, txsMutex)
	}

	go func() {
		time.Sleep(duration)
		stop()
	}()

	limiter := rateLimiter.NewLimiter(rateLimiter.Limit(rate), 1)
	txBytes := make([]byte, txSize)
SendLoop:
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			break SendLoop
		default:
		}

		if err := limiter.Wait(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				err = nil
			}
			return err
		}

		rand.Read(txBytes) //nolint:gosec
		logger.Log(logging.LevelDebug, fmt.Sprintf("Submitting transaction #%d", i))
		if statFileName != "" {
			txsMutex.Lock()
			txs[nextTxNo] = time.Since(clock)
			txsMutex.Unlock()
			nextTxNo++
		}
		if err := client.SubmitTransaction(txBytes); err != nil {
			return err
		}
	}

	// sleep a bit before returning
	time.Sleep(10 * time.Second)
	return nil
}

func clientStats(ctx context.Context, txReceiverAddrs map[t.NodeID]string, clock time.Time, txs map[tt.TxNo]time.Duration, txsMutex *sync.Mutex) {
	confirmations := make(chan time.Duration, int(rate*statPeriod.Seconds()))

	for _, addr := range txReceiverAddrs {
		go populateClientStats(ctx, addr, clock, txs, txsMutex, confirmations)
	}

	ticker := time.NewTicker(statPeriod)
	defer ticker.Stop()

	latencySum := time.Duration(0)
	txCount := int64(0)

	var outputFile *os.File
	if statFileName == "" {
		outputFile = os.Stdout
	} else {
		var err error
		outputFile, err = os.Create(statFileName)
		if err != nil {
			panic(es.Errorf("could not open output file for statistics: %w", err))
		}
		defer outputFile.Close()
	}

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write([]string{"ts", "nrDelivered", "tps", "avgLatency"}); err != nil {
		panic(err)
	}
	writer.Flush()

	defer func() {
		txsMutex.Lock()
		defer txsMutex.Unlock()

		if len(txs) == 0 {
			return
		}

		now := time.Since(clock)
		ts := fmt.Sprintf("%.4f", now.Seconds())

		if err := writer.Write([]string{ts, fmt.Sprintf("%v", len(txs)), "0", "inf"}); err != nil {
			panic(err)
		}
	}()

	lastWrite := time.Since(clock)
	writeStats := func() {
		now := time.Since(clock)

		ts := fmt.Sprintf("%.4f", now.Seconds())
		nrDelivered := fmt.Sprintf("%v", txCount)
		tps := fmt.Sprintf("%.5f", float64(txCount)/(float64(now-lastWrite)/float64(time.Second)))
		avgLatency := fmt.Sprintf("%.5f", float64(latencySum)/float64(txCount)/float64(time.Second))

		if err := writer.Write([]string{ts, nrDelivered, tps, avgLatency}); err != nil {
			panic(err)
		}
		writer.Flush()

		lastWrite = now
		txCount = 0
		latencySum = 0
	}

	for {
		select {
		case <-ctx.Done():
			writeStats()
			return
		case <-ticker.C:
			writeStats()
		case ts := <-confirmations:
			latencySum += ts
			txCount++
		}
	}
}

func populateClientStats(ctx context.Context, txReceiverAddr string, clock time.Time, txs map[tt.TxNo]time.Duration, txsMutex *sync.Mutex, confirmations chan<- time.Duration) {
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

		txsMutex.Lock()
		for _, tx := range batch.Txs {
			if tx.ClientId != id {
				continue
			}

			if ts, ok := txs[tt.TxNo(tx.TxNo)]; ok {
				confirmations <- time.Since(clock) - ts
				delete(txs, tt.TxNo(tx.TxNo))
			}
		}
		txsMutex.Unlock()
	}

}
