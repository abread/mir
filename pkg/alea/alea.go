package alea

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/alea/agreement"
	"github.com/filecoin-project/mir/pkg/alea/broadcast"
	"github.com/filecoin-project/mir/pkg/alea/director"
	"github.com/filecoin-project/mir/pkg/checkpoint"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
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
	ThreshCrypto  t.ModuleID
}

// Params sets the values for the parameters of an instance of the protocol.
type Params struct {
	// Unique identifier for this instance of Alea.
	// Must be the same across all replicas.
	InstanceUID []byte

	// The list of participating nodes, which must be the same as the set of nodes in the threshcrypto module.
	// Must not be empty, and must be the same across all replicas.
	AllNodes []t.NodeID

	// Maximum number of concurrent VCB instances per queue
	// Must be at least 1
	MaxConcurrentVcbPerQueue int

	// Number of batches that the broadcast component tries to have broadcast at all times in own queue
	// Must be at least 1
	TargetOwnUnagreedBatchCount int
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
		Net:           "net",
		ThreshCrypto:  "threshcrypto",
	}
}

// DefaultParams returns the default configuration for a given membership.
// There is no guarantee that this configuration ensures good performance or security,
// but it will pass the CheckParams test.
// DefaultParams is intended for use during testing and hello-world examples.
// A proper deployment is expected to craft a custom configuration,
// for which DefaultParams can serve as a starting point.
func DefaultParams(initialMembership map[t.NodeID]t.NodeAddress) *Params {
	return &Params{
		InstanceUID:                 []byte{42},
		AllNodes:                    maputil.GetKeys(initialMembership),
		MaxConcurrentVcbPerQueue:    10,
		TargetOwnUnagreedBatchCount: 1,
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

	aleaDir := director.NewModule(
		&director.ModuleConfig{
			Self:          config.AleaDirector,
			Consumer:      config.Consumer,
			AleaBroadcast: config.AleaBroadcast,
			AleaAgreement: config.AleaAgreement,
			BatchDB:       config.BatchDB,
			Mempool:       config.Mempool,
			Net:           config.Net,
			ThreshCrypto:  config.ThreshCrypto,
		},
		&director.ModuleParams{
			InstanceUID: params.InstanceUID,
			AllNodes:    params.AllNodes,
		},
		&director.ModuleTunables{
			MaxConcurrentVcbPerQueue:    params.MaxConcurrentVcbPerQueue,
			TargetOwnUnagreedBatchCount: params.TargetOwnUnagreedBatchCount,
		},
		ownID,
		logging.Decorate(logger, "AleaDirector: "),
	)

	aleaBc, errAleaBc := broadcast.NewModule(
		&broadcast.ModuleConfig{
			Self:         config.AleaBroadcast,
			Consumer:     config.AleaDirector,
			Mempool:      config.Mempool,
			Net:          config.Net,
			ThreshCrypto: config.ThreshCrypto,
		},
		&broadcast.ModuleParams{
			InstanceUID: params.InstanceUID,
			AllNodes:    params.AllNodes,
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

	aleaAg := agreement.NewModule(
		&agreement.ModuleConfig{
			Self:         config.AleaAgreement,
			Consumer:     config.AleaDirector,
			Hasher:       config.Hasher,
			Net:          config.Net,
			ThreshCrypto: config.ThreshCrypto,
		},
		&agreement.ModuleParams{
			InstanceUID: params.InstanceUID,
			AllNodes:    params.AllNodes,
		},
		ownID,
		logging.Decorate(logger, "AleaAgreement: "),
	)

	moduleSet := modules.Modules{
		config.AleaDirector:  aleaDir,
		config.AleaBroadcast: aleaBc,
		config.AleaAgreement: aleaAg,
	}

	// Return the initialized protocol module.
	return moduleSet, nil
}

func (p *Params) Check() error {
	if len(p.AllNodes) == 0 {
		return fmt.Errorf("cannot start Alea with an empty membership")
	}

	if p.MaxConcurrentVcbPerQueue < 1 {
		return fmt.Errorf("this implementation requires MaxConcurrentVcbPerQueue >= 1")
	}

	if p.TargetOwnUnagreedBatchCount < 1 {
		return fmt.Errorf("alea requires TargetOwnUnagreedBatchCount >= 1")
	}

	return nil
}
