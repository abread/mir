package bcqueue

import (
	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID
	BatchDB      t.ModuleID
	Mempool      t.ModuleID
	ReliableNet  t.ModuleID
	ThreshCrypto t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(self t.ModuleID, consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:         self,
		Consumer:     consumer,
		BatchDB:      "batchdb",
		Mempool:      "mempool",
		ReliableNet:  "reliablenet",
		ThreshCrypto: "threshcrypto",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	BcInstanceUID []byte
	AllNodes      []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module

	QueueIdx   aleatypes.QueueIdx
	QueueOwner t.NodeID
}

// ModuleTunables sets the values of protocol tunables that need not be the same across all nodes.
type ModuleTunables struct {
	// Maximum number of concurrent VCB instances per queue
	// Must match the equally named tunable in the main Alea module
	// Must be at least 1
	MaxConcurrentVcb int
}
