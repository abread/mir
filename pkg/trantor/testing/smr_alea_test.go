//nolint:dupl
package testing

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	es "github.com/go-errors/errors"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir"
	"github.com/filecoin-project/mir/pkg/deploytest"
	"github.com/filecoin-project/mir/pkg/eventmangler"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	"github.com/filecoin-project/mir/pkg/testsim"
	"github.com/filecoin-project/mir/pkg/trantor"
	"github.com/filecoin-project/mir/pkg/trantor/appmodule"
	"github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

// TODO: try to unify with smr_test

func TestIntegrationAlea(t *testing.T) {
	defer goleak.VerifyNone(t,
		goleak.IgnoreTopFunction("github.com/filecoin-project/mir/pkg/deploytest.newSimModule"),
		goleak.IgnoreTopFunction("github.com/filecoin-project/mir/pkg/testsim.(*Chan).recv"),

		// Problems with this started occurring after an update to a new version of the quic implementation.
		// Assuming it has nothing to do with Mir or Trantor.
		goleak.IgnoreTopFunction("github.com/libp2p/go-libp2p/p2p/transport/quicreuse.(*reuse).gc"),

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
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(1),
				Transport:     "fake",
				Duration:      2 * time.Second,
			}},
		2: {"Submit 10 fake requests with 1 node",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(1),
				Transport:     "fake",
				NumFakeTXs:    10,
				Directory:     "mirbft-deployment-test",
				Duration:      2 * time.Second,
			}},
		5: {"Submit 10 fake requests with 4 nodes and libp2p networking",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(4),
				Transport:     "libp2p",
				NumFakeTXs:    10,
				Duration:      5 * time.Second,
			}},
		6: {"Submit 10 requests with 1 node and libp2p networking",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(1),
				NumClients:    1,
				Transport:     "libp2p",
				NumNetTXs:     10,
				Duration:      5 * time.Second,
			}},
		7: {"Submit 10 requests with 4 nodes and libp2p networking",
			&TestConfig{
				Info:          "libp2p 10 requests and 4 nodes",
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(4),
				NumClients:    1,
				Transport:     "libp2p",
				NumNetTXs:     10,
				Duration:      5 * time.Second,
			}},

		// TODO: fix sim transport with non-transport active modules (threshcrypto breaks it)
		/*8: {"Do nothing with 1 node in simulation",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(1),
				Transport:     "sim",
				Duration:      4 * time.Second,
			}},
		10: {"Submit 10 fake requests with 1 node in simulation",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(1),
				Transport:     "sim",
				NumFakeTXs:    10,
				Directory:     "mirbft-deployment-test",
				Duration:      4 * time.Second,
			}},
		12: {"Submit 100 fake requests with 1 node in simulation",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(1),
				NumClients:    0,
				Transport:     "sim",
				NumFakeTXs:    100,
				Duration:      10 * time.Second,
			}},*/

		100: {"Submit 10 requests with 4 nodes and libp2p networking, with 1 replica not receiving broadcasts",
			&TestConfig{
				Info:          "libp2p 10 requests and 4 nodes, force FILL-GAP/FILLER",
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(4),
				Transport:     "libp2p",
				NumNetTXs:     10,
				NumClients:    1,
				Duration:      10 * time.Second,
				TransportFilter: func(msg *messagepb.Message, from, to types.NodeID) bool {
					node0 := types.NewNodeIDFromInt(0)
					_, isVcb := msg.Type.(*messagepb.Message_Vcb)

					// drop all broadcast messages involving node 0
					// node 0 will not receive broadcasts, but should be able to deliver
					if isVcb && to == node0 && from != to && !types.ModuleID(msg.DestModule).IsSubOf("abc-0") {
						return false
					}
					return true
				},
				ParamsModifier: func(params *trantor.Params) {
					// disable retransmissions
					params.ReliableNet.RetransmissionLoopInterval = math.MaxInt64

					// bring all abba instances into view
					params.Alea.MaxAbbaRoundLookahead = 10
					params.Alea.MaxAgRoundLookahead = 10
				},
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
		0: {"Runs for 20s/400txs with 4 nodes",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(4),
				NumClients:    1,
				Transport:     "sim",
				NumFakeTXs:    100,
				Duration:      20 * time.Second,
				Logger:        logging.ConsoleErrorLogger,
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
	deployment, err := newDeploymentAlea(ctx, conf)
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
	assert.Len(tb, nodeErrors, len(conf.NodeIDsWeight))
	for _, err := range nodeErrors {
		if err != nil {
			assert.Equal(tb, mir.ErrStopped, err)
		}
	}

	// Check event logs
	if conf.CheckFunc != nil {
		conf.CheckFunc(tb, deployment, conf)
		return heapObjects, heapAlloc
	}

	if conf.ErrorExpected != nil {
		require.Error(tb, conf.ErrorExpected)
		return heapObjects, heapAlloc
	}

	assert.NoError(tb, checkEventTraces(deployment.EventLogFiles(), conf.NumNetTXs*conf.NumClients+conf.NumFakeTXs))

	for _, replica := range deployment.TestReplicas {
		// Check if all requests were delivered.
		app := deployment.TestConfig.FakeApps[replica.ID]
		assert.Equal(tb, conf.NumNetTXs*conf.NumClients+conf.NumFakeTXs, int(app.TransactionsProcessed))

		// Check if there are no un-acked messages
		rnet := replica.Modules["rnet"].(*reliablenet.Module)
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

func newDeploymentAlea(ctx context.Context, conf *TestConfig) (*deploytest.Deployment, error) {
	nodeIDs := maputil.GetSortedKeys(conf.NodeIDsWeight)
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
	transportLayer, err := deploytest.NewLocalTransportLayer(simulation, conf.Transport, conf.NodeIDsWeight, logging.Decorate(everythingLogger, "LocalTransport: "))
	if err != nil {
		return nil, es.Errorf("error creating local transport system: %w", err)
	}

	// TODO: fix for weighted stuff
	F := (len(conf.NodeIDsWeight) - 1) / 3
	cryptoSystem := deploytest.NewLocalThreshCryptoSystem("dummy", nodeIDs, 2*F+1)

	nodeModules := make(map[types.NodeID]modules.Modules)
	fakeApps := make(map[types.NodeID]*deploytest.FakeApp)

	for i, nodeID := range nodeIDs {
		nodeLogger := logging.NewMultiLogger(append(
			[]logging.Logger{nodeFileLoggers[i]},
			commonLoggers...,
		))
		nodeLogger = logging.Decorate(nodeLogger, fmt.Sprintf("Node %s: ", nodeID))
		fakeApp := deploytest.NewFakeApp(logging.Decorate(nodeLogger, "FakeApp: "))

		tConf := trantor.DefaultParams(transportLayer.Membership())
		tConf.Alea.MaxConcurrentVcbPerQueue = 2
		tConf.Alea.MaxOwnUnagreedBatchCount = 2
		tConf.Alea.MaxAbbaRoundLookahead = 1
		tConf.Alea.MaxAgRoundLookahead = 1

		// Use small batches so even a few transactions keep being proposed even after epoch transitions.
		tConf.Mempool.MaxTransactionsInBatch = 10
		// Use small retransmission delay to stress-test the reliablenet module.
		tConf.ReliableNet.RetransmissionLoopInterval = 250 * time.Millisecond
		// Keep retransmission bursts low to avoid overloading the test system.
		tConf.ReliableNet.MaxRetransmissionBurst = 4

		if conf.ParamsModifier != nil {
			conf.ParamsModifier(&tConf)
		}

		transport, err := transportLayer.Link(nodeID)
		if err != nil {
			return nil, es.Errorf("error initializing Mir transport: %w", err)
		}
		transport = deploytest.NewFilteredTransport(transport, nodeID, conf.TransportFilter)

		localTCrypto, err := cryptoSystem.ThreshCrypto(nodeID)
		if err != nil {
			return nil, es.Errorf("error creating local threshcrypto system for node %v: %w", nodeID, err)
		}

		system, err := trantor.NewAlea(
			ctx,
			nodeID,
			transport,
			nil,
			localTCrypto,
			appmodule.AppLogicFromStatic(fakeApp, transportLayer.Membership()),
			tConf,
			nodeLogger,
		)
		if err != nil {
			return nil, es.Errorf("error initializing Alea: %w", err)
		}

		if conf.CrashedReplicas[i] {
			err := trantor.PerturbMessages(&eventmangler.ModuleParams{
				DropRate: 1,
			}, trantor.DefaultModuleConfig().Net, system)
			if err != nil {
				return nil, err
			}
		}

		nodeModules[nodeID] = system.Modules()
		fakeApps[nodeID] = fakeApp
	}

	deployConf := &deploytest.TestConfig{
		Info:             conf.Info,
		Simulation:       simulation,
		TransportLayer:   transportLayer,
		NodeIDs:          nodeIDs,
		Membership:       transportLayer.Membership(),
		NodeModules:      nodeModules,
		NumClients:       conf.NumClients,
		NumFakeTXs:       conf.NumFakeTXs,
		NumNetTXs:        conf.NumNetTXs,
		FakeTXDestModule: types.ModuleID("mempool"),
		Directory:        conf.Directory,
		Logger:           everythingLogger,
		FakeApps:         fakeApps,
	}

	return deploytest.NewDeployment(deployConf)
}
