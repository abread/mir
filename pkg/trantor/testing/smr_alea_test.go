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
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir"
	"github.com/filecoin-project/mir/pkg/alea"
	"github.com/filecoin-project/mir/pkg/checkpoint"
	"github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/deploytest"
	"github.com/filecoin-project/mir/pkg/eventmangler"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/net/libp2p"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	"github.com/filecoin-project/mir/pkg/reliablenet" // nolint: typecheck
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

func testIntegrationWithAlea(t *testing.T) { // nolint: gocognit,gocyclo
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
				Duration:      12 * time.Second,
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
				Duration:      12 * time.Second,
			}},

		// TODO: abstract away Go's time to make tests deterministic under simulation (and allow them to pass).
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
		11: {"Submit 10 fake requests with 4 nodes in simulation",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(4),
				Transport:     "sim",
				NumFakeTXs:    10,
				Directory:     "mirbft-deployment-test",
				Duration:      10 * time.Second,
			}},
		12: {"Submit 100 fake requests with 1 node in simulation",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(1),
				NumClients:    0,
				Transport:     "sim",
				NumFakeTXs:    100,
				Duration:      10 * time.Second,
			}},
		13: {"Submit 100 fake requests with 4 nodes in simulation",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(1),
				NumClients:    0,
				Transport:     "sim",
				NumFakeTXs:    100,
				Duration:      15 * time.Second,
				ParamsModifier: func(params *trantor.Params) {
					params.Alea.MaxConcurrentVcbPerQueue = 2
					params.Alea.MaxOwnUnagreedBatchCount = 2
					params.Alea.MaxAbbaRoundLookahead = 1
					params.Alea.MaxAgRoundLookahead = 1
					params.Alea.MaxAgRoundAdvanceInput = 0

					// Use small retransmission delay to stress-test the reliablenet module.
					params.ReliableNet.RetransmissionLoopInterval = 250 * time.Millisecond
				},
			}},*/
		14: {"Submit 100 requests with 4 nodes and libp2p networking with optimizations",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(1),
				Transport:     "libp2p",
				NumClients:    4,
				NumNetTXs:     100,
				Duration:      15 * time.Second,
				ParamsModifier: func(params *trantor.Params) {
					// disable retransmissions
					params.ReliableNet.RetransmissionLoopInterval = math.MaxInt64

					// bring all abba instances into view
					params.Alea.MaxAbbaRoundLookahead = 10
					params.Alea.MaxAgRoundLookahead = 10
					params.Alea.MaxAgRoundEagerInput = 5
				},
			}},

		/*100: {"Submit 10 requests with 4 nodes and libp2p networking, with 1 replica not receiving broadcasts",
			&TestConfig{
				Info:          "libp2p 10 requests and 4 nodes, force FILL-GAP/FILLER",
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(4),
				Transport:     "sim",
				NumFakeTXs:    10,
				Duration:      15 * time.Second,
				TransportFilter: func(msg *messagepbtypes.Message, from, to types.NodeID) bool {
					node0 := types.NewNodeIDFromInt(0)
					_, isVcb := msg.Type.(*messagepbtypes.Message_Vcb)

					// drop all broadcast messages involving node 0
					// node 0 will not receive broadcasts, but should be able to deliver
					if isVcb && to == node0 && from != to && !msg.DestModule.IsSubOf("availability/0") {
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
					params.Alea.MaxAgRoundAdvanceInput = 4
				},
			}},

		// TODO: fix test - check app state, not delivered tx events
		101: {"Test checkpoint recovery with 4 nodes in simulation",
			&TestConfig{
				NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(4),
				Transport:     "sim",
				NumFakeTXs:    128,
				Duration:      120 * time.Second,
				TransportFilter: func(msg *messagepbtypes.Message, from, to types.NodeID) bool {
					node0 := types.NewNodeIDFromInt(0)
					_, isVcb := msg.Type.(*messagepbtypes.Message_Vcb)

					// drop all broadcast messages involving node 0
					// node 0 will not receive broadcasts, but should be able to initiate them
					if isVcb && to == node0 && from != to && msg.DestModule.Sub().Top() != types.ModuleID("0") {
						return false
					}

					// drop nearly all agreement messages involving node 0
					_, isAbba := msg.Type.(*messagepbtypes.Message_Abba)
					if isAbba && (to == node0 || from == node0) && from != to {
						agRoundStr := msg.DestModule.Sub().Top()
						agRound, _ := strconv.ParseUint(string(agRoundStr), 10, 64)
						return agRound > 32-8-3
					}
					agMsgW, isAg := msg.Type.(*messagepbtypes.Message_AleaAgreement)
					if isAg && (to == node0 || from == node0) && from != to {
						finishMsgW := agMsgW.AleaAgreement.Type.(*agreementpbtypes.Message_FinishAbba)
						return finishMsgW.FinishAbba.Round > 32-8-2
					}

					// drop all checkpoint-building messages involving node 0 up to checkpoint 32/8
					_, isChkpBuild := msg.Type.(*messagepbtypes.Message_Threshcheckpoint)
					if isChkpBuild && to == node0 && from != to {
						chkpNum, _ := strconv.ParseUint(string(msg.DestModule.Sub().Top()), 10, 64)
						return chkpNum > 32/8
					}

					return true
				},
				ParamsModifier: func(params *trantor.Params) {
					params.Alea.Adjust(8, 2)
					params.Alea.RetainEpochs = 2
					params.Alea.MaxAgStall = 250 * time.Millisecond
					params.ReliableNet.MaxRetransmissionBurst = 128
					params.ReliableNet.RetransmissionLoopInterval = 500 * time.Millisecond
				},
				CheckFunc: func(tb testing.TB, deployment *deploytest.Deployment, conf *TestConfig) {
					require.Error(tb, conf.ErrorExpected)
					for _, replica := range deployment.TestReplicas {
						app := deployment.TestConfig.FakeApps[replica.ID]
						require.Equal(tb, conf.NumNetTXs+conf.NumFakeTXs, int(app.TransactionsProcessed))

						// Check if there are no un-acked messages
						rnet := replica.Modules["reliablenet"].(*reliablenet.Module)
						pendingMsgs := rnet.GetPendingMessages()
						assert.Emptyf(tb, pendingMsgs, "replica %v has pending messages", replica.ID)
						if len(pendingMsgs) > 0 {
							for _, msg := range pendingMsgs {
								msgJSON, _ := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(msg.Msg.Pb())
								tb.Logf("pending message for %v: %v", msg.Destinations, string(msgJSON))
							}
						}
					}
				},
			}},*/
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
		// Check if all transactions were delivered.
		app := deployment.TestConfig.FakeApps[replica.ID]
		assert.Equal(tb, conf.NumNetTXs*conf.NumClients+conf.NumFakeTXs, int(app.TransactionsProcessed))

		// Check if there are no un-acked messages
		rnet := replica.Modules["reliablenet"].(*reliablenet.Module)
		pendingMsgs := rnet.GetPendingMessages()
		assert.Emptyf(tb, pendingMsgs, "replica %v has pending messages", replica.ID)
		if len(pendingMsgs) > 0 {
			for _, msg := range pendingMsgs {
				msgJSON, _ := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(msg.Msg.Pb())
				tb.Logf("pending message for %v: %v", msg.Destinations, string(msgJSON))
			}
		}
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
		eventDelayFn := func(e *eventpbtypes.Event) time.Duration {
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
	threshCryptoSystem, err := deploytest.NewLocalThreshCryptoSystem("dummy", crypto.DefaultPseudoSeed, nodeIDs, 2*F+1)
	if err != nil {
		return nil, es.Errorf("could not create threshcrypto system: %w", err)
	}

	cryptoSystem, err := deploytest.NewLocalCryptoSystem("pseudo", crypto.DefaultPseudoSeed, nodeIDs, everythingLogger)
	if err != nil {
		return nil, es.Errorf("could not create a local crypto system: %w", err)
	}

	nodeModules := make(map[types.NodeID]modules.Modules)
	fakeApps := make(map[types.NodeID]*deploytest.FakeApp)

	for i, nodeID := range nodeIDs {
		nodeLogger := logging.NewMultiLogger(append(
			[]logging.Logger{nodeFileLoggers[i]},
			commonLoggers...,
		))
		nodeLogger = logging.Decorate(nodeLogger, fmt.Sprintf("Node %s: ", nodeID))
		fakeApp := deploytest.NewFakeApp(logging.Decorate(nodeLogger, "FakeApp: "))
		initialSnapshot, err := fakeApp.Snapshot()
		if err != nil {
			return nil, err
		}

		tConf := trantor.DefaultParams(transportLayer.Membership())
		// Use Alea
		tConf.Protocol = trantor.Alea
		// Use small batches so even a few transactions keep being proposed even after epoch transitions.
		tConf.Mempool.MaxTransactionsInBatch = 10
		// Keep retransmission bursts low to avoid overloading the test system.
		tConf.ReliableNet.MaxRetransmissionBurst = 4
		tConf.Net = libp2p.Params{} // should be unused anyway

		// make epochs progress fast
		tConf.Alea.Adjust(4, 1)
		tConf.Alea.RetainEpochs = 1
		tConf.Alea.MaxAgStall = 150 * time.Millisecond

		if conf.ParamsModifier != nil {
			conf.ParamsModifier(&tConf)
		}

		transport, err := transportLayer.Link(nodeID)
		if err != nil {
			return nil, es.Errorf("error initializing Mir transport: %w", err)
		}
		transport = deploytest.NewFilteredTransport(transport, nodeID, conf.TransportFilter)

		localCrypto, err := cryptoSystem.Crypto(nodeID)
		if err != nil {
			return nil, es.Errorf("error creating local crypto system for node %v: %w", nodeID, err)
		}
		localTCrypto, err := threshCryptoSystem.ThreshCrypto(nodeID)
		if err != nil {
			return nil, es.Errorf("error creating local threshcrypto system for node %v: %w", nodeID, err)
		}

		stateSnapshot, err := alea.InitialStateSnapshot(initialSnapshot, tConf.Alea)
		if err != nil {
			return nil, es.Errorf("error initializing Mir state snapshot: %w", err)
		}
		system, err := trantor.New(
			nodeID,
			transport,
			checkpoint.Genesis(stateSnapshot),
			localCrypto,
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
