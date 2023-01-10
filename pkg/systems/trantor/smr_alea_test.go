//nolint:dupl
package trantor

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir"
	"github.com/filecoin-project/mir/pkg/alea"
	"github.com/filecoin-project/mir/pkg/deploytest"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/net/libp2p"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	"github.com/filecoin-project/mir/pkg/testsim"
	t "github.com/filecoin-project/mir/pkg/types"
)

// TODO: try to unify with smr_test

func TestIntegrationAlea(t *testing.T) {
	defer goleak.VerifyNone(t,
		goleak.IgnoreTopFunction("github.com/filecoin-project/mir/pkg/deploytest.newSimModule"),
		goleak.IgnoreTopFunction("github.com/filecoin-project/mir/pkg/testsim.(*Chan).recv"),

		// If an observable is not exhausted when checking an event trace...
		goleak.IgnoreTopFunction("github.com/reactivex/rxgo/v2.Item.SendContext"),
	)
	t.Run("Alea", testIntegrationWithAlea)
}

func BenchmarkIntegrationAlea(b *testing.B) {
	b.Run("Alea", benchmarkIntegrationWithAlea)
}

func testIntegrationWithAlea(t *testing.T) {
	tests := map[int]struct {
		Desc   string // test description
		Config *TestConfig
	}{
		0: {"Do nothing with 1 node",
			&TestConfig{
				NumReplicas: 1,
				Transport:   "fake",
				Duration:    4 * time.Second,
			}},
		2: {"Submit 10 fake requests with 1 node",
			&TestConfig{
				NumReplicas:     1,
				Transport:       "fake",
				NumFakeRequests: 10,
				Directory:       "mirbft-deployment-test",
				Duration:        4 * time.Second,
			}},
		5: {"Submit 10 fake requests with 4 nodes and gRPC networking",
			&TestConfig{
				NumReplicas:     4,
				NumClients:      1,
				Transport:       "grpc",
				NumFakeRequests: 10,
				Duration:        4 * time.Second,
			}},
		7: {"Submit 10 requests with 4 nodes and gRPC networking",
			&TestConfig{
				Info:           "grpc 10 requests and 4 nodes",
				NumReplicas:    4,
				NumClients:     1,
				Transport:      "grpc",
				NumNetRequests: 10,
				Duration:       4 * time.Second,
			}},
		8: {"Submit 10 fake requests with 4 nodes and libp2p networking",
			&TestConfig{
				NumReplicas:     4,
				NumClients:      1,
				Transport:       "libp2p",
				NumFakeRequests: 10,
				Duration:        15 * time.Second,
			}},
		9: {"Submit 10 requests with 1 node and libp2p networking",
			&TestConfig{
				NumReplicas:    1,
				NumClients:     1,
				Transport:      "libp2p",
				NumNetRequests: 10,
				Duration:       10 * time.Second,
			}},
		10: {"Submit 10 requests with 4 nodes and libp2p networking",
			&TestConfig{
				Info:           "libp2p 10 requests and 4 nodes",
				NumReplicas:    4,
				NumClients:     1,
				Transport:      "libp2p",
				NumNetRequests: 10,
				Duration:       15 * time.Second,
			}},
		11: {"Do nothing with 1 node in simulation",
			&TestConfig{
				NumReplicas: 1,
				Transport:   "sim",
				Duration:    4 * time.Second,
			}},

		13: {"Submit 10 fake requests with 1 node in simulation",
			&TestConfig{
				NumReplicas:     1,
				Transport:       "sim",
				NumFakeRequests: 10,
				Directory:       "mirbft-deployment-test",
				Duration:        4 * time.Second,
			}},
		15: {"Submit 100 fake requests with 1 node in simulation",
			&TestConfig{
				NumReplicas:     1,
				NumClients:      0,
				Transport:       "sim",
				NumFakeRequests: 100,
				Duration:        20 * time.Second,
			}},
	}

	for i, test := range tests {
		i, test := i, test

		// Create a directory for the deployment-generated files and set the test directory name.
		// The directory will be automatically removed when the outer test function exits.
		createDeploymentDir(t, test.Config)

		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			simMode := (test.Config.Transport == "sim") //nolint:goconst
			if testing.Short() && !simMode {
				t.SkipNow()
			}

			if simMode {
				if v := os.Getenv("RANDOM_SEED"); v != "" {
					var err error
					test.Config.RandomSeed, err = strconv.ParseInt(v, 10, 64)
					require.NoError(t, err)
				} else {
					test.Config.RandomSeed = time.Now().UnixNano()
				}
				t.Logf("Random seed = %d", test.Config.RandomSeed)
			}

			runIntegrationWithAleaConfig(t, test.Config)

			if t.Failed() {
				t.Logf("Test #%03d (%s) failed", i, test.Desc)
				if simMode {
					t.Logf("Reproduce with RANDOM_SEED=%d", test.Config.RandomSeed)
				}
			}
		})
	}
}

func benchmarkIntegrationWithAlea(b *testing.B) {
	benchmarks := []struct {
		Desc   string // test description
		Config *TestConfig
	}{
		0: {"Runs for 10s with 4 nodes",
			&TestConfig{
				NumReplicas: 4,
				NumClients:  1,
				Transport:   "fake",
				Duration:    10 * time.Second,
				Logger:      logging.ConsoleErrorLogger,
			}},
		1: {"Runs for 100s with 4 nodes",
			&TestConfig{
				NumReplicas: 4,
				NumClients:  1,
				Transport:   "fake",
				Duration:    100 * time.Second,
				Logger:      logging.ConsoleErrorLogger,
			}},
	}

	for i, bench := range benchmarks {
		i, bench := i, bench
		b.Run(fmt.Sprintf("%03d", i), func(b *testing.B) {
			b.ReportAllocs()

			var totalHeapObjects, totalHeapAlloc float64
			for i := 0; i < b.N; i++ {
				createDeploymentDir(b, bench.Config)
				heapObjects, heapAlloc := runIntegrationWithAleaConfig(b, bench.Config)
				totalHeapObjects += float64(heapObjects)
				totalHeapAlloc += float64(heapAlloc)
			}
			b.ReportMetric(totalHeapObjects/float64(b.N), "heapObjects/op")
			b.ReportMetric(totalHeapAlloc/float64(b.N), "heapAlloc/op")

			if b.Failed() {
				b.Logf("Benchmark #%03d (%s) failed", i, bench.Desc)
			} else {
				b.Logf("Benchmark #%03d (%s) done", i, bench.Desc)
			}
		})
	}
}

func runIntegrationWithAleaConfig(tb testing.TB, conf *TestConfig) (heapObjects int64, heapAlloc int64) {
	tb.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create new test deployment.
	deployment, err := newDeploymentAlea(conf)
	require.NoError(tb, err)

	defer deployment.TestConfig.TransportLayer.Close()

	// Schedule shutdown of test deployment
	if conf.Duration > 0 {
		go func() {
			if deployment.Simulation != nil {
				deployment.Simulation.RunFor(conf.Duration)
			} else {
				time.Sleep(conf.Duration)
			}
			cancel()
		}()
	}

	// Run deployment until it stops and returns final node errors.
	var nodeErrors []error
	nodeErrors, heapObjects, heapAlloc = deployment.Run(ctx)

	// Check whether all the test replicas exited correctly.
	assert.Len(tb, nodeErrors, conf.NumReplicas)
	for _, err := range nodeErrors {
		if err != nil {
			assert.Equal(tb, mir.ErrStopped, err)
		}
	}

	// Check event logs
	require.NoError(tb, checkEventTraces(deployment.EventLogFiles(), conf.NumNetRequests+conf.NumFakeRequests))

	for _, replica := range deployment.TestReplicas {
		// Check if all requests were delivered.
		app := deployment.TestConfig.FakeApps[replica.ID]
		assert.Equalf(tb, conf.NumNetRequests+conf.NumFakeRequests, int(app.RequestsProcessed), "replica %v missed requests", replica.ID)

		// Check if there are no un-acked messages
		rnet := replica.Modules["reliablenet"].(*reliablenet.Module)
		pendingMsgs := rnet.GetPendingMessages()
		assert.Emptyf(tb, pendingMsgs, "replica %v has pending messages", replica.ID)
	}

	// If the test failed, keep the generated data.
	if tb.Failed() {

		// Save the test data.
		testRelDir, err := filepath.Rel(os.TempDir(), conf.Directory)
		require.NoError(tb, err)
		retainedDir := filepath.Join(failedTestDir, testRelDir)

		tb.Logf("Test failed. Saving deployment data to: %s\n", retainedDir)
		err = copy.Copy(conf.Directory, retainedDir)
		require.NoError(tb, err)
	}

	return heapObjects, heapAlloc
}

func fileLogger(conf *TestConfig, name string) (logging.Logger, error) {
	logFile, err := os.OpenFile(path.Join(conf.Directory, name), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return logging.NewStreamLogger(logging.LevelDebug, logFile), nil
}

func newDeploymentAlea(conf *TestConfig) (*deploytest.Deployment, error) {
	nodeIDs := deploytest.NewNodeIDs(conf.NumReplicas)
	testCommonLogger := deploytest.NewLogger(conf.Logger)

	commonFileLogger, err := fileLogger(conf, "common.log")
	if err != nil {
		return nil, err
	}

	nodeFileLoggers := make([]logging.Logger, len(nodeIDs))
	for i, nodeID := range nodeIDs {
		nodeFileLoggers[i], err = fileLogger(conf, fmt.Sprintf("node-%s.log", nodeID))
		if err != nil {
			return nil, err
		}
	}

	commonLoggers := []logging.Logger{
		testCommonLogger,
		commonFileLogger,
	}
	allLoggers := append(slices.Clone(commonLoggers), nodeFileLoggers...)
	everythingLogger := logging.NewMultiLogger(allLoggers)

	var simulation *deploytest.Simulation
	if conf.Transport == "sim" {
		r := rand.New(rand.NewSource(conf.RandomSeed)) // nolint: gosec
		eventDelayFn := func(e *eventpb.Event) time.Duration {
			// TODO: Make min and max event processing delay configurable
			return testsim.RandDuration(r, 0, time.Microsecond)
		}
		simulation = deploytest.NewSimulation(r, nodeIDs, eventDelayFn)
	}
	transportLayer := deploytest.NewLocalTransportLayer(simulation, conf.Transport, nodeIDs, logging.Decorate(everythingLogger, "LocalTransport: "))

	F := (conf.NumReplicas - 1) / 3
	cryptoSystem := deploytest.NewLocalThreshCryptoSystem("pseudo", nodeIDs, 2*F+1, logging.Decorate(everythingLogger, "ThreshCrypto: "))

	nodeModules := make(map[t.NodeID]modules.Modules)
	fakeApps := make(map[t.NodeID]*deploytest.FakeApp)

	for i, nodeID := range nodeIDs {
		// Alea configuration
		aleaConfig := alea.DefaultParams(transportLayer.Nodes())
		aleaConfig.BatchCutFailRetryDelay = t.TimeDuration(100 * time.Millisecond)

		nodeLogger := logging.NewMultiLogger(append(
			[]logging.Logger{nodeFileLoggers[i]},
			commonLoggers...,
		))
		nodeLogger = logging.Decorate(nodeLogger, fmt.Sprintf("Node %s: ", nodeID))
		fakeApp := deploytest.NewFakeApp()

		transport, err := transportLayer.Link(nodeID)
		if err != nil {
			return nil, fmt.Errorf("error initializing Mir transport: %w", err)
		}

		system, err := NewAlea(
			nodeID,
			transport,
			nil,
			cryptoSystem.ThreshCrypto(nodeID),
			AppLogicFromStatic(fakeApp, transportLayer.Nodes()),
			Params{
				Mempool: &simplemempool.ModuleParams{
					MaxTransactionsInBatch: 10,
				},
				Alea: aleaConfig,
				Net:  libp2p.Params{},
			},
			nodeLogger,
		)

		nodeModules[nodeID] = system.Modules()
		fakeApps[nodeID] = fakeApp
	}

	deployConf := &deploytest.TestConfig{
		Info:                   conf.Info,
		Simulation:             simulation,
		TransportLayer:         transportLayer,
		NodeIDs:                nodeIDs,
		Nodes:                  transportLayer.Nodes(),
		NodeModules:            nodeModules,
		NumClients:             conf.NumClients,
		NumFakeRequests:        conf.NumFakeRequests,
		NumNetRequests:         conf.NumNetRequests,
		FakeRequestsDestModule: t.ModuleID("mempool"),
		Directory:              conf.Directory,
		Logger:                 everythingLogger,
		FakeApps:               fakeApps,
	}

	return deploytest.NewDeployment(deployConf)
}
