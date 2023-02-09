package abba

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
	"github.com/filecoin-project/mir/pkg/modules"
	abbadsl "github.com/filecoin-project/mir/pkg/pb/abbapb/dsl"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	"github.com/filecoin-project/mir/pkg/testsim"
	"github.com/filecoin-project/mir/pkg/timer"
	"github.com/filecoin-project/mir/pkg/types"
)

const (
	failedTestDir = "failed-test-data"
)

type TestConfig struct {
	Info              string
	RandomSeed        int64
	N                 int
	F                 int
	Transport         string
	Duration          time.Duration
	Directory         string
	InputValues       map[types.NodeID]bool
	DefaultInputValue bool
	Logger            logging.Logger
}

func runHappyTest(t *testing.T, F int, decision bool, customInputs map[types.NodeID]bool) {
	config := TestConfig{
		N:                 3*F + 1,
		F:                 F,
		Transport:         "sim",
		DefaultInputValue: decision,
		InputValues:       customInputs,
		Duration:          1 * time.Second,
	}

	r, _, _ := runTest(t, &config)
	assert.Equal(t, r, decision)
}

func TestAbbaHappyUnanimousTrue(t *testing.T) {
	runHappyTest(t, 2, true, nil)
}

func TestAbbaHappyUnanimousFalse(t *testing.T) {
	runHappyTest(t, 2, false, nil)
}

func TestAbbaHappyMajorityTrue(t *testing.T) {
	decision := true
	F := 2

	customInputs := make(map[types.NodeID]bool, F)
	for i := 0; i < F; i++ {
		customInputs[types.NewNodeIDFromInt(i)] = !decision
	}

	runHappyTest(t, 2, decision, customInputs)
}

func TestAbbaHappyMajorityFalse(t *testing.T) {
	decision := false
	F := 2

	customInputs := make(map[types.NodeID]bool, F)
	for i := 0; i < F; i++ {
		customInputs[types.NewNodeIDFromInt(i)] = !decision
	}

	runHappyTest(t, 2, decision, customInputs)
}

func runTest(t *testing.T, conf *TestConfig) (result bool, heapObjects int64, heapAlloc int64) {
	t.Helper()

	if v := os.Getenv("RANDOM_SEED"); v != "" {
		var err error
		conf.RandomSeed, err = strconv.ParseInt(v, 10, 64)
		require.NoError(t, err)
	} else {
		conf.RandomSeed = time.Now().UnixNano()
	}
	t.Logf("Random seed = %d", conf.RandomSeed)

	// Create a directory for the deployment-generated files and set the test directory name.
	// The directory will be automatically removed when the outer test function exits.
	createDeploymentDir(t, conf)

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
				t.Log("under simulation")
				deployment.Simulation.RunFor(conf.Duration)
			} else {
				t.Log("using time.sleep")
				time.Sleep(conf.Duration)
			}
			t.Log("simulation time is up")
			cancel()
		}()
	}

	// Run deployment until it stops and returns final node errors.
	var nodeErrors []error
	nodeErrors, heapObjects, heapAlloc = deployment.Run(ctx)

	t.Logf("deployment.Run done")

	// Check whether all the test replicas exited correctly.
	assert.Len(t, nodeErrors, conf.N)
	for _, err := range nodeErrors {
		if err != nil {
			assert.Equal(t, mir.ErrStopped, err)
		}
	}

	app0 := deployment.TestReplicas[0].Modules["app"].(*countingApp)
	assert.Equal(t, app0.firstDeliveredOrigin, types.ModuleID("abba"))
	for _, replica := range deployment.TestReplicas {
		// Check if all requests were delivered exactly once in all replicas.
		app := replica.Modules["app"].(*countingApp)
		assert.Equal(t, 1, app.deliveredCount)
		assert.Equal(t, app0.firstDelivered, app.firstDelivered)
		assert.Equal(t, app0.firstDeliveredOrigin, app.firstDeliveredOrigin)

		// Check if all messages were ACKed
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

	return app0.firstDelivered, heapObjects, heapAlloc
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

	for i, nodeID := range nodeIDs {
		nodeLogger := logging.Decorate(logger, fmt.Sprintf("Node %d: ", i))

		transport, err := transportLayer.Link(nodeID)
		if err != nil {
			return nil, fmt.Errorf("error initializing Mir transport: %w", err)
		}

		abbaConfig := DefaultModuleConfig("app")
		abba := NewModule(abbaConfig, &ModuleParams{
			InstanceUID: []byte{0},
			AllNodes:    nodeIDs,
		}, nodeID, nodeLogger)

		inputValue, ok := conf.InputValues[nodeID]
		if !ok {
			inputValue = conf.DefaultInputValue
		}

		// Use a small retransmission delay to increase likelihood of duplicate messages
		rnetParams := reliablenet.DefaultModuleParams(nodeIDs)
		rnetParams.RetransmissionLoopInterval = 10 * time.Millisecond

		rnet, err := reliablenet.New(
			nodeID,
			&reliablenet.ModuleConfig{
				Self:  abbaConfig.ReliableNet,
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
			abbaConfig.Consumer:     newCountingApp(inputValue),
			abbaConfig.Self:         abba,
			abbaConfig.ThreshCrypto: threshCryptoSystem.Module(nodeID),
			abbaConfig.Hasher:       mirCrypto.NewHasher(crypto.SHA256),
			abbaConfig.ReliableNet:  rnet,
			"net":                   transport,
			"timer":                 timer.New(),
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

	deliveredCount       int
	firstDelivered       bool
	firstDeliveredOrigin types.ModuleID
}

func newCountingApp(inputValue bool) *countingApp {
	m := dsl.NewModule("app")
	m.DslHandle().StoreContext(&struct{}{}) // fill context id 0

	app := &countingApp{
		module: m,
	}

	dsl.UponInit(m, func() error {
		abbadsl.InputValue(m, "abba", inputValue)

		return nil
	})

	abbadsl.UponDeliver(m, func(result bool, from types.ModuleID) error {
		if app.deliveredCount == 0 {
			app.firstDelivered = result
			app.firstDeliveredOrigin = from
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