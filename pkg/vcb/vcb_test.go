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
	mpdsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	"github.com/filecoin-project/mir/pkg/mempool/simplemempool"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	vcbdsl "github.com/filecoin-project/mir/pkg/pb/vcbpb/dsl"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	"github.com/filecoin-project/mir/pkg/testsim"
	"github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	"github.com/filecoin-project/mir/pkg/timer"
	"github.com/filecoin-project/mir/pkg/types"
)

const (
	failedTestDir = "failed-test-data"
)

type TestConfig struct {
	Info       string
	RandomSeed int64
	N          int
	F          int
	Transport  string
	Duration   time.Duration
	Directory  string
	Logger     logging.Logger
}

func TestVcb(t *testing.T) {
	config := TestConfig{
		N:         7,
		F:         2,
		Transport: "libp2p", // TODO: fix sim for goroutine pool active modules (threshcrypto breaks it)
		Duration:  5 * time.Second,
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
	assert.Len(t, nodeErrors, conf.N)
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
	nodeIDs := deploytest.NewNodeIDs(conf.N)
	logger := deploytest.NewLogger(conf.Logger)

	var simulation *deploytest.Simulation
	if conf.Transport == "sim" {
		r := rand.New(rand.NewSource(conf.RandomSeed)) // nolint: gosec
		eventDelayFn := func(e *eventpb.Event) time.Duration {
			// TODO: Make min and max event processing delay configurable
			return testsim.RandDuration(r, 0, time.Microsecond)
		}
		simulation = deploytest.NewSimulation(r, nodeIDs, eventDelayFn)
	}
	transportLayer := deploytest.NewLocalTransportLayer(simulation, conf.Transport, nodeIDs, logger)
	threshCryptoSystem := deploytest.NewLocalThreshCryptoSystem("pseudo", nodeIDs, conf.F+1, logger)

	nodeModules := make(map[types.NodeID]modules.Modules)

	leader := nodeIDs[rand.New(rand.NewSource(conf.RandomSeed)).Intn(len(nodeIDs))] // nolint: gosec

	for i, nodeID := range nodeIDs {
		nodeLogger := logging.Decorate(logger, fmt.Sprintf("Node %d: ", i))

		transport, err := transportLayer.Link(nodeID)
		if err != nil {
			return nil, fmt.Errorf("error initializing Mir transport: %w", err)
		}

		// Use a simple mempool for incoming requests.
		mempool := simplemempool.NewModule(
			ctx,
			&simplemempool.ModuleConfig{
				Self:   "mempool",
				Hasher: "hasher",
			},
			&simplemempool.ModuleParams{
				MaxTransactionsInBatch: 10,
			},
		)

		vcbConfig := DefaultModuleConfig("app")
		vcb := NewModule(ctx, vcbConfig, &ModuleParams{
			InstanceUID: []byte{0},
			AllNodes:    nodeIDs,
			Leader:      leader,
		}, nodeID, logging.Decorate(nodeLogger, "Vcb: "))

		// Use a small retransmission delay to increase likelihood of duplicate messages
		rnetParams := reliablenet.DefaultModuleParams(nodeIDs)
		rnetParams.RetransmissionLoopInterval = 10 * time.Millisecond

		rnet, err := reliablenet.New(
			nodeID,
			&reliablenet.ModuleConfig{
				Self:  vcbConfig.ReliableNet,
				Net:   "net",
				Timer: "timer",
			},
			rnetParams,
			logging.Decorate(nodeLogger, "ReliableNet: "),
		)
		if err != nil {
			return nil, fmt.Errorf("error creating reliablenet module: %w", err)
		}

		modulesWithDefaults := map[types.ModuleID]modules.Module{
			vcbConfig.Consumer:     newCountingApp(ctx, nodeID == leader),
			vcbConfig.Self:         vcb,
			vcbConfig.ThreshCrypto: threshCryptoSystem.Module(ctx, nodeID),
			vcbConfig.Mempool:      mempool,
			"hasher":               mirCrypto.NewHasher(ctx, mirCrypto.DefaultHasherModuleParams(), crypto.SHA256),
			vcbConfig.ReliableNet:  rnet,
			"net":                  transport,
			"timer":                timer.New(),
		}

		nodeModules[nodeID] = modulesWithDefaults
	}

	deployConf := &deploytest.TestConfig{
		Info:           conf.Info,
		Simulation:     simulation,
		TransportLayer: transportLayer,
		NodeIDs:        nodeIDs,
		Nodes:          transportLayer.Nodes(),
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
	firstData      []*requestpb.Request
	firstTxIDs     []types.TxID
	firstSignature tctypes.FullSig
}

func newCountingApp(ctx context.Context, isLeader bool) *countingApp {
	m := dsl.NewModule(ctx, "app")

	app := &countingApp{
		module: m,
	}

	if isLeader {
		txs := []*requestpb.Request{
			{
				ClientId: "asd",
				ReqNo:    42,
				Type:     0,
				Data:     []byte{4, 2},
			},
			{
				ClientId: "asd",
				ReqNo:    4242,
				Type:     0,
				Data:     []byte{2, 4, 2, 4},
			},
		}

		dsl.UponInit(m, func() error {
			mpdsl.RequestTransactionIDs[struct{}](m, "mempool", txs, nil)
			return nil
		})

		mpdsl.UponTransactionIDsResponse(m, func(txIDs []types.TxID, _ctx *struct{}) error {
			vcbdsl.InputValue(m, "vcb", txs)
			return nil
		})
	}

	vcbdsl.UponDeliver(m, func(data []*requestpb.Request, txIDs []types.TxID, signature tctypes.FullSig, srcModule types.ModuleID) error {
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

func (app *countingApp) ApplyEvents(evs *events.EventList) (*events.EventList, error) {
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
