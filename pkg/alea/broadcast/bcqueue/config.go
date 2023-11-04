package bcqueue

import (
	"strconv"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self         t.ModuleID // id of this module
	Consumer     t.ModuleID
	BatchDB      t.ModuleID
	Mempool      t.ModuleID
	Net          t.ModuleID
	ReliableNet  t.ModuleID
	ThreshCrypto t.ModuleID
}

func (mc *ModuleConfig) ParseQueueSlotFromModuleID(id t.ModuleID) (aleatypes.QueueSlot, error) {
	if !id.IsSubOf(mc.Self) {
		return 0, es.Errorf("could not parse queue slot: %v is not a submodule of %v", id, mc.Self)
	}

	queueSlotStr := id.StripParent(mc.Self).Top()
	queueSlot, err := strconv.ParseUint(string(queueSlotStr), 10, 64)
	if err != nil {
		return 0, es.Errorf("deliver event for invalid round: %w", err)
	}

	return aleatypes.QueueSlot(queueSlot), nil
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	AleaInstanceUID []byte
	AllNodes        []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module

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
