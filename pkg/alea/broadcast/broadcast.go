package broadcast

import (
	"strconv"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/alea/aleatypes"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/availability"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	"github.com/filecoin-project/mir/pkg/alea/broadcast/bcqueue"
	"github.com/filecoin-project/mir/pkg/checkpoint"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	transportpbtypes "github.com/filecoin-project/mir/pkg/pb/transportpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleConfig = bccommon.ModuleConfig
type ModuleParams = bccommon.ModuleParams
type ModuleTunables = bccommon.ModuleTunables

type bcMod struct {
	selfID  t.ModuleID
	epochNr tt.EpochNr

	queues       []modules.PassiveModule
	availability modules.PassiveModule
}

func NewModule(
	mc ModuleConfig,
	params ModuleParams,
	tunables ModuleTunables,
	startingChkp *checkpoint.StableCheckpoint,
	nodeID t.NodeID,
	logger logging.Logger,
) (modules.PassiveModule, error) {
	if startingChkp == nil {
		return nil, es.Errorf("missing initial checkpoint")
	}
	epochNr := startingChkp.Epoch()

	queues, err := createQueues(mc, params, tunables, nodeID, logger)
	if err != nil {
		return nil, err
	}

	return &bcMod{
		selfID:       mc.Self,
		epochNr:      epochNr,
		availability: availability.New(mc, params, tunables, nodeID, logger),
		queues:       queues,
	}, nil
}

func createQueues(bcMc ModuleConfig, bcParams bccommon.ModuleParams, bcTunables ModuleTunables, nodeID t.NodeID, logger logging.Logger) ([]modules.PassiveModule, error) {
	queues := make([]modules.PassiveModule, 0, len(bcParams.AllNodes))

	tunables := bcqueue.ModuleTunables{
		MaxConcurrentVcb: bcTunables.MaxConcurrentVcbPerQueue,
	}

	for idx := range bcParams.AllNodes {
		if int(aleatypes.QueueIdx(idx)) != idx {
			return nil, es.Errorf("queue idx %v is not representable", idx)
		}

		mc := bcqueue.ModuleConfig{
			Self:         bccommon.BcQueueModuleID(bcMc.Self, aleatypes.QueueIdx(idx)),
			Consumer:     bcMc.Self,
			BatchDB:      bcMc.BatchDB,
			Mempool:      bcMc.Mempool,
			Net:          bcMc.Net,
			ReliableNet:  bcMc.ReliableNet,
			ThreshCrypto: bcMc.ThreshCrypto,
		}

		params := bcqueue.ModuleParams{
			BcInstanceUID: bcParams.InstanceUID, // TODO: review
			AllNodes:      bcParams.AllNodes,

			QueueIdx:   aleatypes.QueueIdx(idx),
			QueueOwner: bcParams.AllNodes[idx],
		}

		mod, err := bcqueue.New(mc, params, tunables, nodeID, logging.Decorate(logger, "BcQueue: ", "queueIdx", params.QueueIdx, "queueOwner", params.QueueOwner))
		if err != nil {
			return nil, err
		}

		queues = append(queues, mod)
	}

	return queues, nil
}

func (bc *bcMod) ImplementsModule() {}
func (bc *bcMod) ApplyEvents(evsIn events.EventList) (events.EventList, error) {
	evsInByMod, err := bc.splitEvsIn(evsIn)
	if err != nil {
		return events.EmptyList(), err
	}

	evsOutChan := make(chan events.EventList)
	errOutChan := make(chan error)
	for mod, evsIn := range evsInByMod {
		go func(mod modules.PassiveModule, evsIn events.EventList) {
			evsOut, err := modules.ApplyAllSafely(mod, evsIn)
			if err == nil {
				evsOutChan <- evsOut
			} else {
				errOutChan <- err
			}
		}(mod, *evsIn)
	}

	evsOut := events.EmptyList()
	var firstError error
	for i := 0; i < len(evsInByMod); i++ {
		select {
		case subEvsOut := <-evsOutChan:
			evsOut.PushBackList(subEvsOut)
		case err := <-errOutChan:
			if firstError == nil {
				firstError = err
			}
		}
	}

	return evsOut, firstError
}

func (bc *bcMod) splitEvsIn(evsIn events.EventList) (map[modules.PassiveModule]*events.EventList, error) {
	res := make(map[modules.PassiveModule]*events.EventList)

	for _, ev := range evsIn.Slice() {
		var destID t.ModuleID
		if ev.DestModule == bc.selfID {
			destID = ""
		} else {
			destID = ev.DestModule.StripParent(bc.selfID)
		}

		var mod modules.PassiveModule
		if destID == "" {
			mod = bc.availability
		} else if n, err := strconv.ParseUint(string(destID.Top()), 10, 64); err == nil && n < uint64(len(bc.queues)) {
			mod = bc.queues[int(n)]
		}

		if mod == nil {
			if transportEv, ok := ev.Type.(*eventpbtypes.Event_Transport); ok {
				if _, ok := transportEv.Transport.Type.(*transportpbtypes.Event_MessageReceived); ok {
					// ignore event, other node is byz
					// TODO: signal byz node
					continue
				}
			}

			return nil, es.Errorf("failed to find destination module %s for bc event", ev.DestModule)
		}

		if _, ok := res[mod]; !ok {
			res[mod] = &events.EventList{}
		}
		res[mod].PushBack(ev)
	}

	return res, nil
}
