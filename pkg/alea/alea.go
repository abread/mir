package alea

import (
	"fmt"
	"time"

	"github.com/filecoin-project/mir/pkg/alea/agreement"
	"github.com/filecoin-project/mir/pkg/alea/broadcast"
	"github.com/filecoin-project/mir/pkg/alea/director"
	"github.com/filecoin-project/mir/pkg/checkpoint"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	commonpbtypes "github.com/filecoin-project/mir/pkg/pb/commonpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

// Config sets the module ids. All replicas are expected to use identical module configurations.
type Config struct {
	AleaDirector  t.ModuleID
	AleaBroadcast t.ModuleID
	AleaAgreement t.ModuleID
	Consumer      t.ModuleID
	BatchDB       t.ModuleID
	Hasher        t.ModuleID
	Mempool       t.ModuleID
	Net           t.ModuleID
	ReliableNet   t.ModuleID
	ThreshCrypto  t.ModuleID
	Timer         t.ModuleID
}

// Params sets the values for the parameters of an instance of the protocol.
type Params struct {
	// Unique identifier for this instance of Alea.
	// Must be the same across all replicas.
	InstanceUID []byte

	// The identities of all nodes that execute the protocol.
	// Must not be empty.
	Membership *commonpbtypes.Membership

	// Maximum number of concurrent VCB instances per queue
	// Must be at least 1
	MaxConcurrentVcbPerQueue int

	// Maximum number of unagreed batches that the broadcast component can have in this node's queue
	// Must be at least 1
	MaxOwnUnagreedBatchCount uint64

	// TODO
	MaxAbbaRoundLookahead int

	// TODO
	MaxAgRoundLookahead int

	// Pad broadcast duration estimate
	// Must be non-negative
	// TODO: 1/2*RTT + time to verify ?
	BcEstimateMargin time.Duration

	// Time to wait before resorting to FILL-GAP messages
	FillGapDelay time.Duration

	// Maximum time to stall agreement round waiting for broadcasts to complete
	MaxAgreementDelay time.Duration
}

// DefaultConfig returns a valid module config with default names for all modules.
func DefaultConfig(consumer t.ModuleID) *Config {
	return &Config{
		AleaDirector:  "alea_dir",
		AleaBroadcast: "alea_bc",
		AleaAgreement: "alea_ag",
		Consumer:      consumer,
		BatchDB:       "batchdb",
		Hasher:        "hasher",
		Mempool:       "mempool",
		ReliableNet:   "reliablenet",
		Net:           "net",
		ThreshCrypto:  "threshcrypto",
		Timer:         "timer",
	}
}

// DefaultParams returns the default configuration for a given membership.
// There is no guarantee that this configuration ensures good performance or security,
// but it will pass the CheckParams test.
// DefaultParams is intended for use during testing and hello-world examples.
// A proper deployment is expected to craft a custom configuration,
// for which DefaultParams can serve as a starting point.
func DefaultParams(membership *commonpbtypes.Membership) *Params {
	aproxRTT := 220 * time.Microsecond
	aproxBcDuration := 3*aproxRTT/2 + 20*time.Millisecond

	return &Params{
		InstanceUID:              []byte{42},
		Membership:               membership,
		MaxConcurrentVcbPerQueue: 32,
		MaxOwnUnagreedBatchCount: 1,
		MaxAbbaRoundLookahead:    4,
		MaxAgRoundLookahead:      32,
		BcEstimateMargin:         7*time.Millisecond + aproxRTT/2,
		FillGapDelay:             aproxBcDuration,
		MaxAgreementDelay:        aproxBcDuration,
	}
}

// New returns a new initialized instance of the base Alea protocol modules to be used when instantiating a mir.Node.
// Arguments:
//   - ownID:        the ID of the node being instantiated with Alea.
//   - moduleConfig: the IDs of the modules Alea interacts with.
//   - params:       Alea protocol-specific configuration (e.g. membership,  etc...).
//     see the documentation of the ModuleParams type for details.
//   - startingChkp: the stable checkpoint defining the initial state of the protocol.
//   - logger:       Logger the Alea implementation uses to output log messages.
func New(ownID t.NodeID, config *Config, params *Params, startingChkp *checkpoint.StableCheckpoint, logger logging.Logger) (modules.Modules, error) {
	if logger == nil {
		logger = logging.ConsoleErrorLogger
	}

	// Check whether the passed configuration is valid.
	if err := params.Check(); err != nil {
		return nil, fmt.Errorf("invalid Alea parameters: %w", err)
	}

	if startingChkp != nil {
		// TODO: checkpointing
		return nil, fmt.Errorf("alea checkpointing not implemented")
	}

	allNodes := params.AllNodes()

	aleaDir := director.NewModule(
		&director.ModuleConfig{
			Self:          config.AleaDirector,
			Consumer:      config.Consumer,
			AleaBroadcast: config.AleaBroadcast,
			AleaAgreement: config.AleaAgreement,
			BatchDB:       config.BatchDB,
			Mempool:       config.Mempool,
			Net:           config.Net,
			ReliableNet:   config.ReliableNet,
			Hasher:        config.Hasher,
			ThreshCrypto:  config.ThreshCrypto,
			Timer:         config.Timer,
		},
		&director.ModuleParams{
			InstanceUID: append(params.InstanceUID, 'd'),
			AllNodes:    allNodes,
		},
		&director.ModuleTunables{
			MaxConcurrentVcbPerQueue: params.MaxConcurrentVcbPerQueue,
			MaxOwnUnagreedBatchCount: params.MaxOwnUnagreedBatchCount,
			BcEstimateMargin:         params.BcEstimateMargin,
			FillGapDelay:             params.FillGapDelay,
			MaxAgreementDelay:        params.MaxAgreementDelay,
		},
		ownID,
		logging.Decorate(logger, "AleaDirector: "),
	)

	aleaBc, errAleaBc := broadcast.NewModule(
		&broadcast.ModuleConfig{
			Self:         config.AleaBroadcast,
			Consumer:     config.AleaDirector,
			BatchDB:      config.BatchDB,
			Mempool:      config.Mempool,
			ReliableNet:  config.ReliableNet,
			Hasher:       config.Hasher,
			ThreshCrypto: config.ThreshCrypto,
		},
		&broadcast.ModuleParams{
			InstanceUID: append(params.InstanceUID, 'b'),
			AllNodes:    allNodes,
		},
		&broadcast.ModuleTunables{
			MaxConcurrentVcbPerQueue: params.MaxConcurrentVcbPerQueue,
		},
		ownID,
		logging.Decorate(logger, "AleaBroadcast: "),
	)
	if errAleaBc != nil {
		return nil, fmt.Errorf("error creating alea broadcast: %w", errAleaBc)
	}

	aleaAg, errAleaAg := agreement.NewModule(
		&agreement.ModuleConfig{
			Self:         config.AleaAgreement,
			Consumer:     config.AleaDirector,
			Hasher:       config.Hasher,
			ReliableNet:  config.ReliableNet,
			Net:          config.Net,
			ThreshCrypto: config.ThreshCrypto,
		},
		&agreement.ModuleParams{
			InstanceUID: append(params.InstanceUID, 'a'),
			AllNodes:    allNodes,
		},
		&agreement.ModuleTunables{
			MaxRoundLookahead:     params.MaxAgRoundLookahead,
			MaxAbbaRoundLookahead: params.MaxAbbaRoundLookahead,
		},
		ownID,
		logging.Decorate(logger, "AleaAgreement: "),
	)
	if errAleaAg != nil {
		return nil, fmt.Errorf("error creating alea agreement: %w", errAleaAg)
	}

	moduleSet := modules.Modules{
		config.AleaDirector:  aleaDir,
		config.AleaBroadcast: aleaBc,
		config.AleaAgreement: aleaAg,
	}

	// Return the initialized protocol module.
	return moduleSet, nil
}

func (p *Params) Check() error {
	if len(p.Membership.Nodes) == 0 {
		return fmt.Errorf("cannot start Alea with an empty membership")
	}

	if p.MaxConcurrentVcbPerQueue < 1 {
		return fmt.Errorf("this implementation requires MaxConcurrentVcbPerQueue >= 1")
	}

	if p.MaxOwnUnagreedBatchCount < 1 {
		return fmt.Errorf("alea requires TargetOwnUnagreedBatchCount >= 1")
	}

	return nil
}

func (p *Params) AllNodes() []t.NodeID {
	return maputil.GetSortedKeys(p.Membership.Nodes)
}
