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

	es "github.com/go-errors/errors"

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
	abbapbtypes "github.com/filecoin-project/mir/pkg/pb/abbapb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	"github.com/filecoin-project/mir/pkg/reliablenet"
	"github.com/filecoin-project/mir/pkg/testsim"
	"github.com/filecoin-project/mir/pkg/timer"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	"github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

const (
	failedTestDir = "failed-test-data"
)

type TestConfig struct {
	Info              string
	RandomSeed        int64
	NodeIDsWeight     map[types.NodeID]tt.VoteWeight
	F                 int
	Transport         string
	Duration          time.Duration
	Directory         string
	InputValues       map[types.NodeID]bool
	DefaultInputValue bool

	TransportFilter func(msg *messagepbtypes.Message, from types.NodeID, to types.NodeID) bool
	Logger          logging.Logger
}

func runHappyTest(t *testing.T, F int, decision bool, customInputs map[types.NodeID]bool, duration time.Duration) {
	config := TestConfig{
		NodeIDsWeight:     deploytest.NewNodeIDsDefaultWeights(3*F + 1),
		F:                 F,
		Transport:         "sim",
		DefaultInputValue: decision,
		InputValues:       customInputs,
		Duration:          duration,
	}

	r, _, _ := runTest(t, &config)
	assert.Equal(t, r, decision)
}

func TestAbbaHappyUnanimousTrueOneNode(t *testing.T) {
	runHappyTest(t, 0, true, nil, 5*time.Second)
}

func TestAbbaHappyUnanimousTrue(t *testing.T) {
	runHappyTest(t, 2, true, nil, 5*time.Second)
}

func TestAbbaHappyUnanimousFalse(t *testing.T) {
	runHappyTest(t, 2, false, nil, 5*time.Second)
}

func TestAbbaHappyMajorityTrue(t *testing.T) {
	decision := true
	F := 2

	customInputs := make(map[types.NodeID]bool, F)
	for i := 0; i < F; i++ {
		customInputs[types.NewNodeIDFromInt(i)] = !decision
	}

	runHappyTest(t, 2, decision, customInputs, 20*time.Second)
}

func TestAbbaHappyMajorityFalse(t *testing.T) {
	decision := false
	F := 2

	customInputs := make(map[types.NodeID]bool, F)
	for i := 0; i < F; i++ {
		customInputs[types.NewNodeIDFromInt(i)] = !decision
	}

	runHappyTest(t, 2, decision, customInputs, 20*time.Second)
}

func runUnanimousOptimizationTest(t *testing.T, F int, decision bool, duration time.Duration) {
	config := TestConfig{
		NodeIDsWeight:     deploytest.NewNodeIDsDefaultWeights(3*F + 1),
		F:                 F,
		Transport:         "sim",
		DefaultInputValue: decision,
		Duration:          duration,
		TransportFilter: func(msg *messagepbtypes.Message, from, to types.NodeID) bool {
			if abbaMsg, ok := msg.Type.(*messagepbtypes.Message_Abba); ok {
				if abbaRoundMsg, ok := abbaMsg.Abba.Type.(*abbapbtypes.Message_Round); ok {
					switch abbaRoundMsg.Round.Type.(type) {
					case *abbapbtypes.RoundMessage_Input:
						return true
					default:
						// filter out all non-INPUT ABBA messages to force the result to come only from the unanimity optimization
						return false
					}
				}
			}

			return true
		},
	}

	r, _, _ := runTest(t, &config)
	assert.Equal(t, r, decision)
}

func TestAbbaUnanimousOptimizationTrue(t *testing.T) {
	runUnanimousOptimizationTest(t, 2, true, 5*time.Second)
}

func TestAbbaUnanimousOptimizationFalse(t *testing.T) {
	runUnanimousOptimizationTest(t, 2, false, 5*time.Second)
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
	assert.Len(t, nodeErrors, len(conf.NodeIDsWeight))
	for _, err := range nodeErrors {
		if err != nil {
			assert.Equal(t, mir.ErrStopped, err)
		}
	}

	app0 := deployment.TestReplicas[0].Modules[abbaConfig.Consumer].(*countingApp)
	for _, replica := range deployment.TestReplicas {
		// Check if all requests were delivered exactly once in all replicas.
		app := replica.Modules[abbaConfig.Consumer].(*countingApp)
		assert.Equal(t, 1, app.deliveredCount)
		assert.Equal(t, abbaConfig.Self, app.firstSrcModule)
		assert.Equal(t, app0.firstDelivered, app.firstDelivered)

		// Check if ABBA reported protocol termination
		assert.Equal(t, 1, app.doneCount)

		// Check if all messages were ACKed
		rnet := replica.Modules[abbaConfig.ReliableNet].(*reliablenet.Module)
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

var abbaConfig = ModuleConfig{
	Self:         "abba",
	Consumer:     "app",
	ReliableNet:  "rnet",
	ThreshCrypto: "tc",
	Hasher:       "hasher",
}
var rnetConfig = reliablenet.ModuleConfig{
	Self:  abbaConfig.ReliableNet,
	Net:   "net",
	Timer: "timer",
}

func newDeployment(conf *TestConfig) (*deploytest.Deployment, error) {
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
		return nil, es.Errorf("could not create transport: %w", err)
	}

	threshCryptoSystem, err := deploytest.NewLocalThreshCryptoSystem("pseudo", nodeIDs, 2*conf.F+1)
	if err != nil {
		return nil, es.Errorf("could not create threshcrypto system: %w", err)
	}

	nodeModules := make(map[types.NodeID]modules.Modules)

	for i, nodeID := range nodeIDs {
		nodeLogger := logging.Decorate(logger, fmt.Sprintf("Node %d: ", i))

		transport, err := transportLayer.Link(nodeID)
		if err != nil {
			return nil, es.Errorf("error initializing Mir transport: %w", err)
		}
		transport = deploytest.NewFilteredTransport(transport, nodeID, conf.TransportFilter)

		abba, err := NewModule(
			abbaConfig,
			ModuleParams{
				InstanceUID: []byte{0},
				AllNodes:    nodeIDs,
			},
			ModuleTunables{
				MaxRoundLookahead: 1,
			},
			nodeID,
			nodeLogger,
		)
		if err != nil {
			return nil, err
		}

		inputValue, ok := conf.InputValues[nodeID]
		if !ok {
			inputValue = conf.DefaultInputValue
		}

		// Use a small retransmission delay to increase likelihood of duplicate messages
		rnetParams := reliablenet.DefaultModuleParams(nodeIDs)
		rnetParams.RetransmissionLoopInterval = 50 * time.Millisecond

		rnet, err := reliablenet.New(nodeID, rnetConfig, rnetParams, logging.Decorate(nodeLogger, "ReliableNet: "))
		if err != nil {
			return nil, es.Errorf("error creating reliablenet module: %w", err)
		}

		tc, err := threshCryptoSystem.Module(nodeID)
		if err != nil {
			return nil, es.Errorf("failed to create threshcrypto module: %w", err)
		}

		modulesWithDefaults := map[types.ModuleID]modules.Module{
			abbaConfig.Consumer:     newCountingApp(abbaConfig, inputValue),
			abbaConfig.Self:         abba,
			abbaConfig.ThreshCrypto: tc,
			abbaConfig.Hasher:       mirCrypto.NewHasher(crypto.SHA256),
			abbaConfig.ReliableNet:  rnet,
			rnetConfig.Net:          transport,
			rnetConfig.Timer:        timer.New(),
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
	doneCount      int
	firstSrcModule types.ModuleID
	firstDelivered bool
}

func newCountingApp(abbaMc ModuleConfig, inputValue bool) *countingApp {
	m := dsl.NewModule(abbaMc.Consumer)

	app := &countingApp{
		module: m,
	}

	dsl.UponInit(m, func() error {
		abbadsl.InputValue(m, abbaMc.Self, inputValue)
		abbadsl.ContinueExecution(m, abbaMc.Self)

		return nil
	})

	abbadsl.UponDeliver(m, func(result bool, srcModule types.ModuleID) error {
		if app.deliveredCount == 0 {
			app.firstSrcModule = srcModule
			app.firstDelivered = result
		}

		app.deliveredCount++

		return nil
	})

	abbadsl.UponDone(m, func(srcModule types.ModuleID) error {
		app.doneCount++
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
