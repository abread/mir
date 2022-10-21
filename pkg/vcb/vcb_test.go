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
	"github.com/filecoin-project/mir/pkg/pb/dslpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	"github.com/filecoin-project/mir/pkg/pb/vcbpb"
	"github.com/filecoin-project/mir/pkg/testsim"
	"github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/vcb/vcbdsl"
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
		Transport: "sim",
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
	deployment, err := newDeployment(conf)
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

	// Check if all requests were delivered exactly once in all replicas.
	app0 := deployment.TestReplicas[0].Modules["app"].(*countingApp)
	for _, replica := range deployment.TestReplicas {
		app := replica.Modules["app"].(*countingApp)
		assert.Equal(t, 1, app.deliveredCount)
		// TODO: check request data
		assert.Equal(t, app0.firstBatchID, app.firstBatchID)
		assert.ElementsMatch(t, app0.firstSignature, app.firstSignature)
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

func newDeployment(conf *TestConfig) (*deploytest.Deployment, error) {
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

	for _, nodeID := range nodeIDs {
		//nodeLogger := logging.Decorate(logger, fmt.Sprintf("Node %d: ", i))

		transport, err := transportLayer.Link(nodeID)
		if err != nil {
			return nil, fmt.Errorf("error initializing Mir transport: %w", err)
		}

		// Use a simple mempool for incoming requests.
		mempool := simplemempool.NewModule(
			&simplemempool.ModuleConfig{
				Self:   "mempool",
				Hasher: "hasher",
			},
			&simplemempool.ModuleParams{
				MaxTransactionsInBatch: 10,
			},
		)

		vcb := NewModule(DefaultModuleConfig("app"), &ModuleParams{
			InstanceUID: []byte{0},
			AllNodes:    nodeIDs,
			Leader:      leader,
			Origin: &vcbpb.Origin{
				Module: "app",
				Type: &vcbpb.Origin_Dsl{
					Dsl: &dslpb.Origin{
						ContextID: 0,
					},
				},
			},
		}, nodeID)

		modulesWithDefaults := map[types.ModuleID]modules.Module{
			"app":          newCountingApp(nodeID == leader),
			"vcb":          vcb,
			"threshcrypto": threshCryptoSystem.Module(nodeID),
			"hasher":       mirCrypto.NewHasher(crypto.SHA256),
			"net":          transport,
			"mempool":      mempool,
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
	firstData      []*requestpb.Request
	firstBatchID   types.BatchID
	firstSignature []byte
}

func newCountingApp(isLeader bool) *countingApp {
	m := dsl.NewModule("app")
	m.DslHandle().StoreContext(&struct{}{}) // fill context id 0

	app := &countingApp{
		module: m,
	}

	if isLeader {
		dsl.UponInit(m, func() error {
			vcbdsl.Request(m, "vcb", []*requestpb.Request{
				{
					ClientId: "asd",
					ReqNo:    42,
					Type:     0,
					Data:     []byte{4, 2},
				},
			})

			return nil
		})
	}

	vcbdsl.UponDeliver(m, func(data []*requestpb.Request, batchID types.BatchID, signature []byte, context *struct{}) error {
		if app.deliveredCount == 0 {
			app.firstData = data
			app.firstBatchID = batchID
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
