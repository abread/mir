package vcb

import (
	"context"
	"crypto"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/otiai10/copy"

	"github.com/filecoin-project/mir"
	mirCrypto "github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/deploytest"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool"
	"github.com/filecoin-project/mir/pkg/modules"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	mpdsl "github.com/filecoin-project/mir/pkg/pb/mempoolpb/dsl"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	vcbpbdsl "github.com/filecoin-project/mir/pkg/pb/vcbpb/dsl"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	"github.com/filecoin-project/mir/pkg/testsim"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	"github.com/filecoin-project/mir/pkg/timer"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	"github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"

	es "github.com/go-errors/errors"
)

const (
	failedTestDir = "failed-test-data"
)

type TestConfig struct {
	Info          string
	RandomSeed    int64
	NodeIDsWeight map[types.NodeID]tt.VoteWeight
	F             int
	Transport     string
	Duration      time.Duration
	Directory     string
	Logger        logging.Logger
}

func TestVcb(t *testing.T) {
	F := 2
	config := TestConfig{
		NodeIDsWeight: deploytest.NewNodeIDsDefaultWeights(3*F + 1),
		F:             F,
		Transport:     "libp2p", // TODO: fix sim for goroutine pool active modules (threshcrypto breaks it)
		Duration:      5 * time.Second,
	}

	if v := os.Getenv("RANDOM_SEED"); v != "" {
		var err error
		config.RandomSeed, err = strconv.ParseInt(v, 10, 64)
		require.NoError(t, err)
	} else {
		config.RandomSeed = time.Now().UnixNano()
	}
	t.Logf("Random seed = %d", config.RandomSeed)

	// Create a directory for the deployment-generated files and set the test directory name.
	// The directory will be automatically removed when the outer test function exits.
	createDeploymentDir(t, &config)

	runTest(t, &config)
}

func runTest(t *testing.T, conf *TestConfig) (heapObjects int64, heapAlloc int64) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create new test deployment.
	deployment, err := newDeployment(ctx, conf)
	require.NoError(t, err)

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
	assert.Len(t, nodeErrors, len(conf.NodeIDsWeight))
	for _, err := range nodeErrors {
		if err != nil {
			assert.Equal(t, mir.ErrStopped, err)
		}
	}

	app0 := deployment.TestReplicas[0].Modules["app"].(*countingApp)
	for _, replica := range deployment.TestReplicas {
		// Check if all requests were delivered exactly once in all replicas.
		app := replica.Modules["app"].(*countingApp)
		assert.Equal(t, 1, app.deliveredCount)
		assert.Equal(t, types.ModuleID("vcb"), app.firstSrcModule)
		// TODO: check request data
		assert.ElementsMatch(t, app0.firstTxIDs, app.firstTxIDs)
		assert.ElementsMatch(t, app0.firstSignature, app.firstSignature)

		// Check that all messages were properly ACKed
		rnet := replica.Modules["reliablenet"].(*reliablenet.Module)
		assert.Empty(t, rnet.GetPendingMessages())
	}

	// If the test failed, keep the generated data.
	if t.Failed() {

		// Save the test data.
		testRelDir, err := filepath.Rel(os.TempDir(), conf.Directory)
		require.NoError(t, err)
		retainedDir := filepath.Join(failedTestDir, testRelDir)

		t.Logf("Test failed. Saving deployment data to: %s\n", retainedDir)
		err = copy.Copy(conf.Directory, retainedDir)
		require.NoError(t, err)
	}

	return heapObjects, heapAlloc
}

func newDeployment(ctx context.Context, conf *TestConfig) (*deploytest.Deployment, error) {
	nodeIDs := maputil.GetSortedKeys(conf.NodeIDsWeight)
	logger := deploytest.NewLogger(conf.Logger)

	var simulation *deploytest.Simulation
	if conf.Transport == "sim" {
		r := rand.New(rand.NewSource(conf.RandomSeed)) // nolint: gosec
		eventDelayFn := func(e *eventpbtypes.Event) time.Duration {
			// TODO: Make min and max event processing delay configurable
			return testsim.RandDuration(r, 0, time.Microsecond)
		}
		simulation = deploytest.NewSimulation(r, nodeIDs, eventDelayFn)
	}
	transportLayer, err := deploytest.NewLocalTransportLayer(simulation, conf.Transport, conf.NodeIDsWeight, logger)
	if err != nil {
		return nil, es.Errorf("failed to create local transport: %w", err)
	}

	threshCryptoSystem := deploytest.NewLocalThreshCryptoSystem("pseudo", nodeIDs, 2*conf.F+1)

	nodeModules := make(map[types.NodeID]modules.Modules)

	leader := nodeIDs[rand.New(rand.NewSource(conf.RandomSeed)).Intn(len(nodeIDs))] // nolint: gosec

	for i, nodeID := range nodeIDs {
		nodeLogger := logging.Decorate(logger, fmt.Sprintf("Node %d: ", i))

		transport, err := transportLayer.Link(nodeID)
		if err != nil {
			return nil, es.Errorf("error initializing Mir transport: %w", err)
		}
		vcbConfig := ModuleConfig{
			Self:         "vcb",
			Consumer:     "app",
			Net:          "net",
			ReliableNet:  "reliablenet",
			Hasher:       "hasher",
			ThreshCrypto: "threshcrypto",
			Mempool:      "mempool",
		}

		// Use a simple mempool for incoming requests.
		mempool := simplemempool.NewModule(
			simplemempool.ModuleConfig{
				Self:   vcbConfig.Mempool,
				Hasher: vcbConfig.Hasher,
			},
			&simplemempool.ModuleParams{
				MaxTransactionsInBatch: 10,
			},
			logging.Decorate(nodeLogger, "Mempool: "),
		)

		vcb := NewModule(vcbConfig, &ModuleParams{
			InstanceUID: []byte{0},
			AllNodes:    nodeIDs,
			Leader:      leader,
		}, nodeID, logging.Decorate(nodeLogger, "Vcb: "))

		// Use a small retransmission delay to increase likelihood of duplicate messages
		rnetParams := reliablenet.DefaultModuleParams(nodeIDs)
		rnetParams.RetransmissionLoopInterval = 10 * time.Millisecond

		rnet, err := reliablenet.New(
			nodeID,
			reliablenet.ModuleConfig{
				Self:  vcbConfig.ReliableNet,
				Net:   vcbConfig.Net,
				Timer: "timer",
			},
			rnetParams,
			logging.Decorate(nodeLogger, "ReliableNet: "),
		)
		if err != nil {
			return nil, es.Errorf("error creating reliablenet module: %w", err)
		}

		tc, err := threshCryptoSystem.Module(ctx, nodeID)
		if err != nil {
			return nil, es.Errorf("failed to build threshcrypto: %w", err)
		}

		modulesWithDefaults := map[types.ModuleID]modules.Module{
			vcbConfig.Consumer:     newCountingApp(nodeID == leader),
			vcbConfig.Self:         vcb,
			vcbConfig.ThreshCrypto: tc,
			vcbConfig.Mempool:      mempool,
			vcbConfig.Hasher:       mirCrypto.NewHasher(ctx, mirCrypto.DefaultHasherModuleParams(), crypto.SHA256),
			vcbConfig.ReliableNet:  rnet,
			vcbConfig.Net:          transport,
			"timer":                timer.New(),
		}

		nodeModules[nodeID] = modulesWithDefaults
	}

	deployConf := &deploytest.TestConfig{
		Info:           conf.Info,
		Simulation:     simulation,
		TransportLayer: transportLayer,
		NodeIDs:        nodeIDs,
		Membership:     transportLayer.Membership(),
		NodeModules:    nodeModules,
		Directory:      conf.Directory,
		Logger:         logger,
	}

	return deploytest.NewDeployment(deployConf)
}

type countingApp struct {
	module dsl.Module

	deliveredCount int
	firstSrcModule types.ModuleID
	firstData      []*trantorpbtypes.Transaction
	firstTxIDs     []tt.TxID
	firstSignature tctypes.FullSig
}

func newCountingApp(isLeader bool) *countingApp {
	m := dsl.NewModule("app")

	app := &countingApp{
		module: m,
	}

	if isLeader {
		txs := []*trantorpbtypes.Transaction{
			{
				ClientId: "asd",
				TxNo:     42,
				Type:     0,
				Data:     []byte{4, 2},
			},
			{
				ClientId: "asd",
				TxNo:     4242,
				Type:     0,
				Data:     []byte{2, 4, 2, 4},
			},
		}

		dsl.UponInit(m, func() error {
			mpdsl.RequestTransactionIDs[struct{}](m, "mempool", txs, nil)
			return nil
		})

		mpdsl.UponTransactionIDsResponse(m, func(txIDs []tt.TxID, _ctx *struct{}) error {
			vcbpbdsl.InputValue(m, "vcb", txs)
			return nil
		})
	}

	vcbpbdsl.UponDeliver(m, func(data []*trantorpbtypes.Transaction, txIDs []tt.TxID, signature tctypes.FullSig, srcModule types.ModuleID) error {
		if app.deliveredCount == 0 {
			app.firstSrcModule = srcModule
			app.firstData = data
			app.firstTxIDs = txIDs
			app.firstSignature = signature
		}

		app.deliveredCount++

		return nil
	})

	return app
}

func (app *countingApp) ApplyEvents(evs events.EventList) (events.EventList, error) {
	return app.module.ApplyEvents(evs)
}

func (app *countingApp) ImplementsModule() {}

// If conf.Directory is not empty, creates a directory with that path if it does not yet exist.
// If conf.Directory is empty, creates a directory in the OS-default temporary path
// and sets conf.Directory to that path.
func createDeploymentDir(tb testing.TB, conf *TestConfig) {
	tb.Helper()
	if conf == nil {
		return
	}

	if conf.Directory != "" {
		conf.Directory = filepath.Join(os.TempDir(), conf.Directory)
		tb.Logf("Using deployment dir: %s\n", conf.Directory)
		err := os.MkdirAll(conf.Directory, 0777)
		require.NoError(tb, err)
		tb.Cleanup(func() { os.RemoveAll(conf.Directory) })
	} else {
		// If no directory is configured, create a temporary directory in the OS-default location.
		conf.Directory = tb.TempDir()
		tb.Logf("Created temp dir: %s\n", conf.Directory)
	}
}
