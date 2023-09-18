package director

import (
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self          t.ModuleID // id of this module
	App           t.ModuleID
	Checkpoint    t.ModuleID
	ChkpValidator t.ModuleID
	AleaBroadcast t.ModuleID
	AleaAgreement t.ModuleID
	BatchDB       t.ModuleID
	Mempool       t.ModuleID
	Net           t.ModuleID
	ReliableNet   t.ModuleID
	Hasher        t.ModuleID
	ThreshCrypto  t.ModuleID
	Timer         t.ModuleID
	Null          t.ModuleID
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte                     // unique identifier for this instance of Alea, must be the same in Broadcast and Agreement components
	Membership  *trantorpbtypes.Membership // identities of all nodes that execute the protocol, must not be empty

	// length of an epoch in ag rounds
	EpochLength uint64

	RetainEpochs uint64
}

// ModuleTunables sets the values of protocol tunables that need not be the same across all nodes.
type ModuleTunables struct {
	// Maximum number of concurrent VCB instances per queue
	// Must be at least 1
	MaxConcurrentVcbPerQueue int

	// Maximum number of unagreed batches that the broadcast component can have in this node's queue
	// Must be at least 1
	MaxOwnUnagreedBatchCount int

	// Subprotocol duration estimates window size
	EstimateWindowSize int

	// Maximum factor that the F slowest nodes are allowed to delay bc completion vs the rest.
	// A factor of K will let the system wait up to <time for local deliver> + K * <time for 2F+1 VCB quorum done>
	// Must be >= 1
	MaxExtSlowdownFactor float64
}
