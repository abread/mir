package common

import (
	"time"

	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self          t.ModuleID // id of this module
	Consumer      t.ModuleID
	BcQueuePrefix string
	AleaAgreement t.ModuleID
	BatchDB       t.ModuleID
	Mempool       t.ModuleID
	Net           t.ModuleID
	ReliableNet   t.ModuleID
	Hasher        t.ModuleID
	ThreshCrypto  t.ModuleID
	Timer         t.ModuleID
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // unique identifier for this instance of Alea, must be the same in Broadcast and Agreement components
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
}

// ModuleTunables sets the values of protocol tunables that need not be the same across all nodes.
type ModuleTunables struct {
	// Maximum number of concurrent VCB instances per queue
	// Must be at least 1
	MaxConcurrentVcbPerQueue int

	// Maximum number of unagreed batches that the broadcast component can have in this node's queue
	// Must be at least 1
	MaxOwnUnagreedBatchCount uint64

	// Pad broadcast duration estimate
	InitialBcEstimateMargin time.Duration

	// Time to wait before resorting to FILL-GAP messages
	FillGapDelay time.Duration

	// Maximum time to stall agreement round waiting for broadcasts to complete
	MaxAgreementDelay time.Duration
}
