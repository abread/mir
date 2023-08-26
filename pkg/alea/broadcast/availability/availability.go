package availability

import (
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	bcpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleConfig = bccommon.ModuleConfig
type ModuleTunables = bccommon.ModuleTunables
type ModuleParams = bccommon.ModuleParams

func New(mc ModuleConfig, params ModuleParams, tunables ModuleTunables, nodeID t.NodeID, logger logging.Logger) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	certDB := make(map[bcpbtypes.Slot]*bcpbtypes.Cert)

	est := newBcEstimators(m, mc, params, tunables, nodeID)
	includeCertVerification(m, mc, params)
	includeCertCreation(m, mc, params, nodeID, certDB)
	includeBatchFetching(m, mc, params, logger, certDB, est)

	bcpbdsl.UponFreeSlot(m, func(slot *bcpbtypes.Slot) error {
		// propagate event to queue
		destModule := bccommon.BcQueueModuleID(mc.Self, slot.QueueIdx)
		bcqueuepbdsl.FreeSlot(m, destModule, slot.QueueSlot)
		return nil
	})

	bcqueuepbdsl.UponBcStarted(m, func(slot *bcpbtypes.Slot) error {
		// propagate event to consumer
		bcpbdsl.BcStarted(m, mc.Consumer, slot)
		return nil
	})

	return m
}
