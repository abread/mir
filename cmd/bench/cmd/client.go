// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	rateLimiter "golang.org/x/time/rate"

	"github.com/filecoin-project/mir/pkg/dummyclient"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/membership"
	"github.com/filecoin-project/mir/pkg/rrclient"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

var (
	txSize     int
	rate       float64
	burst      int
	duration   time.Duration
	clientType string

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
		return fmt.Errorf("could not load membership: %w", err)
	}
	addresses, err := membership.GetIPs(initialMembership)
	if err != nil {
		return fmt.Errorf("could not load node IPs: %w", err)
	}

	// Generate addresses and ports for transaction receivers.
	// Each node uses different ports for receiving protocol messages and transactions.
	// These addresses will be used by the client code to know where to send its transactions.
	txReceiverAddrs := make(map[t.NodeID]string)
	for nodeID, nodeIP := range addresses {
		numericID, err := strconv.Atoi(string(nodeID))
		if err != nil {
			return fmt.Errorf("node IDs must be numeric in the sample app: %w", err)
		}
		txReceiverAddrs[nodeID] = net.JoinHostPort(nodeIP, fmt.Sprintf("%d", TxReceiverBasePort+numericID))

		// The break statement causes the client to send its transactions to only one single node.
		// Remove it for the client to send its transactions to all nodes.
		// TODO: Make this properly configurable and remove the hack.
		// break // sending to all nodes is more reproducible
	}

	ctx, stop := context.WithCancel(ctx)

	client := clientFactories[clientType](
		tt.ClientID(id),
		crypto.SHA256,
		logger,
	)
	client.Connect(ctx, txReceiverAddrs)
	defer client.Disconnect()

	go func() {
		time.Sleep(duration)
		stop()
	}()

	limiter := rateLimiter.NewLimiter(rateLimiter.Limit(rate), 1)
	txBytes := make([]byte, txSize)
	for i := 0; ; i++ {
		if err := limiter.Wait(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				err = nil
			}
			return err
		}
		rand.Read(txBytes) //nolint:gosec
		logger.Log(logging.LevelDebug, fmt.Sprintf("Submitting transaction #%d", i))
		if err := client.SubmitTransaction(txBytes); err != nil {
			return err
		}
	}
}
