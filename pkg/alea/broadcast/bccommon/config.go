package bccommon

import (
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // prefix for queue module IDs
	AleaDirector t.ModuleID
	BatchDB      t.ModuleID
	Mempool      t.ModuleID
	Net          t.ModuleID
	ReliableNet  t.ModuleID
	ThreshCrypto t.ModuleID
	Timer        t.ModuleID
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // must be the alea instance uid followed by 'b'
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module

	EpochLength uint64
}

// ModuleTunables sets the values of protocol tunables that need not be the same across all nodes.
type ModuleTunables struct {
	// Maximum number of concurrent VCB instances per queue
	// Must be at least 1
	MaxConcurrentVcbPerQueue int

	// Subprotocol duration estimates window size
	EstimateWindowSize int

	// Maximum factor that the F slowest nodes are allowed to delay completion vs the rest.
	// A factor of K will let the system wait up to <time for local completion> + K * <time for 2F+1 quorum done>
	// Must be >= 1
	MaxExtSlowdownFactor float64
}
