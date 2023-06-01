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
	"sync/atomic"
	"time"

	es "github.com/go-errors/errors"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/membership"
	"github.com/filecoin-project/mir/pkg/transactionreceiver"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

var (
	//txSize     int
	latclientTxCount int
	cooldown         time.Duration
	//clientType string
	//statFileName string
	//statPeriod   time.Duration

	latclientCmd = &cobra.Command{
		Use:   "latclient",
		Short: "Measure E2E transaction latency in a Mir cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runLatClient(ctx)
		},
	}
)

func init() {
	rootCmd.AddCommand(latclientCmd)
	latclientCmd.Flags().IntVarP(&txSize, "txSize", "s", 256, "size of each transaction in bytes")
	latclientCmd.Flags().IntVarP(&latclientTxCount, "count", "c", 500, "number of txs to send and wait for reply")
	latclientCmd.Flags().DurationVarP(&cooldown, "cooldown", "C", 200*time.Millisecond, "cooldown between transactions")
	latclientCmd.Flags().StringVarP(&clientType, "type", "t", "dummy", "client type (one of: dummy, rr)")
	latclientCmd.Flags().StringVarP(&statFileName, "statFile", "o", "", "output file for statistics")
	latclientCmd.Flags().DurationVar(&statPeriod, "statPeriod", time.Second, "statistic record period")
}

func runLatClient(ctx context.Context) error {
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

	ctx, stop := context.WithCancel(ctx)
	defer func() {
		stop()
		time.Sleep(5 * time.Second)
	}()

	client := clientFactories[clientType](
		tt.ClientID(id),
		crypto.SHA256,
		logger,
	)
	client.Connect(ctx, txReceiverAddrs)
	defer client.Disconnect()

	clock := time.Now()

	currentTxNo := uint64(0)
	confirmations := make(chan time.Duration, len(txReceiverAddrs))
	for _, addr := range txReceiverAddrs {
		go collectConfirmations(ctx, addr, clock, confirmations, &currentTxNo)
	}

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

	if err := writer.Write([]string{"firstLatency", "fp1Latency", "fp2Latency", "fullPropLatency"}); err != nil {
		panic(err)
	}
	writer.Flush()

	F := (len(txReceiverAddrs) - 1) / 3

	txBytes := make([]byte, txSize)
	for i := 0; i < latclientTxCount; i++ {
		rand.Read(txBytes) //nolint:gosec

		submitTs := time.Since(clock) // nolint:stylecheck

		logger.Log(logging.LevelDebug, fmt.Sprintf("Submitting transaction #%d", i))
		if err := client.SubmitTransaction(txBytes); err != nil {
			return err
		}

		firstLat := <-confirmations - submitTs
		fp1LatSum := firstLat
		fp2LatSum := firstLat
		fullPropLatSum := firstLat

		for i := 0; i < F; i++ {
			l := <-confirmations - submitTs
			fp1LatSum += l
			fp2LatSum += l
			fullPropLatSum += l
		}
		for i := 0; i < F; i++ {
			l := <-confirmations - submitTs
			fp2LatSum += l
			fullPropLatSum += l
		}
		for i := 0; i < F; i++ {
			l := <-confirmations - submitTs
			fullPropLatSum += l
		}

		if err := writer.Write([]string{
			fmt.Sprintf("%.6f", firstLat.Seconds()),
			fmt.Sprintf("%.6f", fp1LatSum.Seconds()/float64(F+1)),
			fmt.Sprintf("%.6f", fp2LatSum.Seconds()/float64(2*F+1)),
			fmt.Sprintf("%.6f", fullPropLatSum.Seconds()/float64(len(txReceiverAddrs))),
		}); err != nil {
			panic(err)
		}
		writer.Flush()

		time.Sleep(cooldown)
		atomic.AddUint64(&currentTxNo, 1)
	}

	return nil
}

func collectConfirmations(ctx context.Context, txReceiverAddr string, clock time.Time, confirmations chan time.Duration, currentTxNo *uint64) {
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
		ts := time.Since(clock)

		if errors.Is(err, io.EOF) {
			return
		} else if err != nil {
			fmt.Printf("error reading replica: %v", err)
			return
		}

		for _, tx := range batch.Txs {
			if tx.ClientId != id {
				continue
			} else if tx.TxNo != atomic.LoadUint64(currentTxNo) {
				continue // error?
			}

			confirmations <- ts
		}
	}
}
