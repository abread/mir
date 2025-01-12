package availability

import (
	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	bcpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/dsl"
	bcpbtypes "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb/types"
	bcqueuepbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/bcqueuepb/dsl"
	directorpbdsl "github.com/filecoin-project/mir/pkg/pb/aleapb/directorpb/dsl"
	apppbdsl "github.com/filecoin-project/mir/pkg/pb/apppb/dsl"
	checkpointpbtypes "github.com/filecoin-project/mir/pkg/pb/checkpointpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
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
		bcpbdsl.BcStarted(m, mc.AleaDirector, slot)
		return nil
	})

	// propagate epoch changes to queues
	directorpbdsl.UponNewEpoch(m, func(epoch tt.EpochNr) error {
		for i := aleatypes.QueueIdx(0); i < aleatypes.QueueIdx(len(params.AllNodes)); i++ {
			directorpbdsl.NewEpoch(m, bccommon.BcQueueModuleID(mc.Self, i), epoch)
		}
		return nil
	})
	directorpbdsl.UponGCEpochs(m, func(minEpoch tt.EpochNr) error {
		for i := aleatypes.QueueIdx(0); i < aleatypes.QueueIdx(len(params.AllNodes)); i++ {
			directorpbdsl.GCEpochs(m, bccommon.BcQueueModuleID(mc.Self, i), minEpoch)
		}
		return nil
	})
	apppbdsl.UponRestoreState(m, func(checkpoint *checkpointpbtypes.StableCheckpoint) error {
		for i := aleatypes.QueueIdx(0); i < aleatypes.QueueIdx(len(params.AllNodes)); i++ {
			apppbdsl.RestoreState(m, bccommon.BcQueueModuleID(mc.Self, i), checkpoint)
		}
		return nil
	})

	return m
}
