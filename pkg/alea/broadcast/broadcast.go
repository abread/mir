package broadcast

import (
	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bcqueue"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bcutil"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ConfigTemplate sets the module ids. All replicas are expected to use identical module configurations.
type ConfigTemplate struct {
	SelfPrefix   string // prefix for queue module IDs
	Consumer     t.ModuleID
	BatchDB      t.ModuleID
	Mempool      t.ModuleID
	ReliableNet  t.ModuleID
	Hasher       t.ModuleID
	ThreshCrypto t.ModuleID
}

// ParamsTemplate sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ParamsTemplate struct {
	InstanceUID []byte     // must be the alea instance uid followed by 'b'
	AllNodes    []t.NodeID // the list of participating nodes, which must be the same as the set of nodes in the threshcrypto module
}

func CreateQueues(mcTemplate ConfigTemplate, paramsTemplate ParamsTemplate, tunables bcqueue.ModuleTunables, nodeID t.NodeID, logger logging.Logger) (modules.Modules, error) {
	queues := make(map[t.ModuleID]modules.Module, len(paramsTemplate.AllNodes))

	for idx := range paramsTemplate.AllNodes {
		if int(aleatypes.QueueIdx(idx)) != idx {
			return nil, es.Errorf("queue idx %v is not representable", idx)
		}

		mc := bcqueue.ModuleConfig{
			Self:         bcutil.BcQueueModuleID(mcTemplate.SelfPrefix, aleatypes.QueueIdx(idx)),
			Consumer:     mcTemplate.Consumer,
			BatchDB:      mcTemplate.BatchDB,
			Mempool:      mcTemplate.Mempool,
			ReliableNet:  mcTemplate.ReliableNet,
			Hasher:       mcTemplate.Hasher,
			ThreshCrypto: mcTemplate.ThreshCrypto,
		}

		params := bcqueue.ModuleParams{
			BcInstanceUID: paramsTemplate.InstanceUID, // TODO: review
			AllNodes:      paramsTemplate.AllNodes,

			QueueIdx:   aleatypes.QueueIdx(idx),
			QueueOwner: paramsTemplate.AllNodes[idx],
		}

		mod, err := bcqueue.New(mc, params, tunables, nodeID, logging.Decorate(logger, "BcQueue: ", "queueIdx", idx))
		if err != nil {
			return nil, err
		}

		queues[mc.Self] = mod
	}

	return queues, nil
}
