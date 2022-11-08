package common

import t "github.com/filecoin-project/mir/pkg/types"

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self          t.ModuleID // id of this module
	Consumer      t.ModuleID
	AleaBroadcast t.ModuleID
	AleaAgreement t.ModuleID
	BatchDB       t.ModuleID
	Mempool       t.ModuleID
	Net           t.ModuleID
	ThreshCrypto  t.ModuleID
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters, apart from the Origin.
type ModuleParams struct {
	InstanceUID []byte     // unique identifier for this instance of Alea
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
}

// ModuleTunables sets the values of protocol tunables that need not be the same across all nodes.
type ModuleTunables struct {
	MaxConcurrentVcbInstancesPerQueue int
	TargetOwnUnagreedBatchCount       int
}
