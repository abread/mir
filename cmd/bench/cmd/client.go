// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	rateLimiter "golang.org/x/time/rate"

	es "github.com/go-errors/errors"
	"github.com/spf13/cobra"

	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/membership"
	"github.com/filecoin-project/mir/pkg/rendezvous"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

var (
	//clientType string
	rate float64

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

	// Required arguments
	clientCmd.Flags().StringVarP(&configFileName, "config-file", "c", "", "configuration file")
	_ = clientCmd.MarkFlagRequired("config-file")
	clientCmd.Flags().Float64VarP(&rate, "rate", "r", 128, "rate of requests to submit per second")
	_ = clientCmd.MarkFlagRequired("rate")
	clientCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "client ID")
	_ = clientCmd.MarkPersistentFlagRequired("id")

	// Optional arguments
	clientCmd.Flags().StringVarP(&clientType, "type", "t", "dummy", "client type (one of: dummy, rr)")

	clientCmd.Flags().DurationVar(&statPeriod, "stat-period", 5*time.Second, "statistic record period")
	clientCmd.Flags().StringVar(&clientStatsFileName, "live-stat-file", "", "output file for live statistics, default is standard output")
	clientCmd.Flags().StringVar(&statSummaryFileName, "summary-stat-file", "", "output file for summarized statistics")

	// Sync files
	clientCmd.Flags().StringVar(&readySyncFileName, "ready-sync-file", "", "file to use for initial synchronization when ready to start the benchmark")
	clientCmd.Flags().StringVar(&deliverSyncFileName, "deliver-sync-file", "", "file to use for synchronization when waiting to deliver all transactions")
}

func runClient(ctx context.Context) error {
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
	time.Sleep(3 * time.Second)

	clientCtx, stop := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}

	// Start generating the load and measuring performance.
	logger.Log(logging.LevelWarn, "applying load")
	wg.Add(1)
	go func() {
		defer wg.Done()

		limiter := rateLimiter.NewLimiter(rateLimiter.Limit(rate), 1)
		txBytes := make([]byte, params.TxGen.PayloadSize)
		for {
			select {
			case <-clientCtx.Done():
				return
			default:
			}

			if err := limiter.Wait(clientCtx); err != nil {
				if !errors.Is(err, context.Canceled) {
					panic(err)
				}
			}

			rand.Read(txBytes) //nolint:gosec
			if err := client.SubmitTransaction(txBytes); err != nil {
				if errors.Is(err, io.EOF) {
					select {
					case <-clientCtx.Done():
						return
					default:
						panic(es.Errorf("server closed: %w", err))
					}
				}
				panic(err)
			}
		}
	}()

	done := make(chan struct{})
	if params.Duration > 0 {
		go func() {
			// Wait until the end of the benchmark and shut down the node.
			select {
			case <-ctx.Done():
			case <-time.After(time.Duration(params.Duration)):
			}
			stop()
			wg.Wait()
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
			stop()
			wg.Wait()
			close(done)
		}()
	}

	<-done
	return nil
}
