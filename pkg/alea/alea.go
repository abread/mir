package alea

import (
	"time"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/agreement"
	"github.com/filecoin-project/mir/pkg/alea/broadcast"
	"github.com/filecoin-project/mir/pkg/alea/director"
	"github.com/filecoin-project/mir/pkg/checkpoint"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
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
	Null          t.ModuleID
}

// Params sets the values for the parameters of an instance of the protocol.
type Params struct {
	// Unique identifier for this instance of Alea.
	// Must be the same across all replicas.
	InstanceUID []byte

	// The identities of all nodes that execute the protocol.
	// Must not be empty.
	Membership *trantorpbtypes.Membership

	// Maximum number of concurrent VCB instances per queue
	// Must be at least 1
	MaxConcurrentVcbPerQueue int

	// Maximum number of unagreed batches that the broadcast component can have in this node's queue
	// Must be at least 1
	MaxOwnUnagreedBatchCount int

	// Maximum number of concurrent ABBA rounds for which we process messages
	// Must be at least 1
	MaxAbbaRoundLookahead int

	// Maximum number of concurrent agreement rounds for which we process messages
	// Must be at least 1
	MaxAgRoundLookahead int

	// Maximum number of agreement rounds for which we send input before
	// allowing them to progress in the normal path.
	// Must be at least 0, must be less that MaxRoundLookahead
	// Setting this parameter too high will lead to costly retransmissions!
	// Should likely be less than MaxRoundLookahead/2 - 1.
	MaxAgRoundAdvanceInput int

	// Subprotocol duration estimates window size
	EstimateWindowSize int

	// Maximum time to stall agreement round waiting for broadcasts to complete
	MaxAgreementDelay time.Duration

	// How slower can the F slowest nodes be compared to the rest
	// A factor of 2, means we allow F nodes to take double the time we expect the rest to take.
	// Must be >=1
	MaxExtSlowdownFactor float64
}

// DefaultParams returns the default configuration for a given membership.
// There is no guarantee that this configuration ensures good performance or security,
// but it will pass the CheckParams test.
// DefaultParams is intended for use during testing and hello-world examples.
// A proper deployment is expected to craft a custom configuration,
// for which DefaultParams can serve as a starting point.
func DefaultParams(membership *trantorpbtypes.Membership) Params {
	return Params{
		InstanceUID:              []byte{42},
		Membership:               membership,
		MaxConcurrentVcbPerQueue: 32,
		MaxOwnUnagreedBatchCount: 8,
		MaxAbbaRoundLookahead:    4,
		MaxAgRoundLookahead:      32,
		MaxAgRoundAdvanceInput:   15, // AgRoundLookahead / 2 - 1
		EstimateWindowSize:       32,
		MaxAgreementDelay:        time.Second,
		MaxExtSlowdownFactor:     1.25,
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
func New(ownID t.NodeID, config Config, params Params, startingChkp *checkpoint.StableCheckpoint, logger logging.Logger) (modules.Modules, error) {
	if logger == nil {
		logger = logging.ConsoleErrorLogger
	}

	// Check whether the passed configuration is valid.
	if err := params.Check(); err != nil {
		return nil, es.Errorf("invalid Alea parameters: %w", err)
	}

	if startingChkp != nil {
		// TODO: checkpointing
		return nil, es.Errorf("alea checkpointing not implemented")
	}

	// TODO: support weighted voting
	// easiest solution is distributing <node weight> keys to each node, and sign with all of them in
	// each operation this is probably best to encode in the threshcrypto module, since weighted
	// threshold crypto schemes seem to exist (see "An Efficient and Secure Weighted Threshold
	// Signcryption Scheme", Chien-Hua Tsai, Journal of Internet Technology, 2019)
	w := params.Membership.Nodes[ownID].Weight
	for _, identity := range params.Membership.Nodes {
		if identity.Weight != w {
			return nil, es.Errorf("alea does not support weighted voting (yet): node %v cannot have weight %v != %v (own weight)", identity.Id, identity.Weight, w)
		}
	}

	allNodes := params.AllNodes()

	aleaDir := director.NewModule(
		director.ModuleConfig{
			Self:     config.AleaDirector,
			Consumer: config.Consumer,
			// TODO: implement proper epoch change/checkpointing
			AleaBroadcast: config.AleaBroadcast.Then(t.NewModuleIDFromInt(0)),
			AleaAgreement: config.AleaAgreement,
			BatchDB:       config.BatchDB,
			Mempool:       config.Mempool,
			Net:           config.Net,
			ReliableNet:   config.ReliableNet,
			Hasher:        config.Hasher,
			ThreshCrypto:  config.ThreshCrypto,
			Timer:         config.Timer,
			Null:          config.Null,
		},
		director.ModuleParams{
			InstanceUID: append(params.InstanceUID, 'd'),
			AllNodes:    allNodes,
		},
		director.ModuleTunables{
			MaxConcurrentVcbPerQueue: params.MaxConcurrentVcbPerQueue,
			MaxOwnUnagreedBatchCount: params.MaxOwnUnagreedBatchCount,
			EstimateWindowSize:       params.EstimateWindowSize,
			MaxAgreementDelay:        params.MaxAgreementDelay,
			MaxExtSlowdownFactor:     params.MaxExtSlowdownFactor,
		},
		ownID,
		logging.Decorate(logger, "AleaDirector: "),
	)

	aleaBc, errAleaBc := broadcast.NewMulti(
		broadcast.ModuleConfig{
			Self:         config.AleaBroadcast,
			Consumer:     config.AleaDirector,
			BatchDB:      config.BatchDB,
			Mempool:      config.Mempool,
			Net:          config.Net,
			ReliableNet:  config.ReliableNet,
			ThreshCrypto: config.ThreshCrypto,
			Timer:        config.Timer,
		},
		broadcast.ModuleParams{
			InstanceUID: append(params.InstanceUID, 'b'),
			AllNodes:    allNodes,
		},
		broadcast.ModuleTunables{
			MaxConcurrentVcbPerQueue: params.MaxConcurrentVcbPerQueue,
			EstimateWindowSize:       params.EstimateWindowSize,
			MaxExtSlowdownFactor:     params.MaxExtSlowdownFactor,
		},
		ownID,
		logging.Decorate(logger, "AleaBroadcast: "),
	)
	if errAleaBc != nil {
		return nil, es.Errorf("error creating alea broadcast: %w", errAleaBc)
	}

	aleaAg, errAleaAg := agreement.NewModule(
		agreement.ModuleConfig{
			Self:         config.AleaAgreement,
			Consumer:     config.AleaDirector,
			Hasher:       config.Hasher,
			ReliableNet:  config.ReliableNet,
			Net:          config.Net,
			ThreshCrypto: config.ThreshCrypto,
		},
		agreement.ModuleParams{
			InstanceUID: append(params.InstanceUID, 'a'),
			AllNodes:    allNodes,
		},
		agreement.ModuleTunables{
			MaxRoundLookahead:     params.MaxAgRoundLookahead,
			MaxAbbaRoundLookahead: params.MaxAbbaRoundLookahead,
			MaxRoundAdvanceInput:  params.MaxAgRoundAdvanceInput,
		},
		ownID,
		logging.Decorate(logger, "AleaAgreement: "),
	)
	if errAleaAg != nil {
		return nil, es.Errorf("error creating alea agreement: %w", errAleaAg)
	}

	moduleSet := make(modules.Modules, 3)
	moduleSet[config.AleaDirector] = aleaDir
	moduleSet[config.AleaAgreement] = aleaAg
	moduleSet[config.AleaBroadcast] = aleaBc

	// Return the initialized protocol module.
	return moduleSet, nil
}

func (p *Params) Check() error {
	if len(p.Membership.Nodes) == 0 {
		return es.Errorf("cannot start Alea with an empty membership")
	}

	if p.MaxConcurrentVcbPerQueue < 1 {
		return es.Errorf("this implementation requires MaxConcurrentVcbPerQueue >= 1")
	}

	if p.MaxOwnUnagreedBatchCount < 1 {
		return es.Errorf("alea requires TargetOwnUnagreedBatchCount >= 1")
	}

	return nil
}

func (p *Params) AllNodes() []t.NodeID {
	return maputil.GetSortedKeys(p.Membership.Nodes)
}
