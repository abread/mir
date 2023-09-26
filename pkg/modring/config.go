package modring

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleConfig struct {
	Self t.ModuleID
}

type ModuleGenerator func(id t.ModuleID, idx uint64) (modules.PassiveModule, events.EventList, error)
type CleanupHandler func(idx uint64) (events.EventList, error)
type PastMessagesHandler func(pastMessages []*modringpbtypes.PastMessage) (events.EventList, error)

type ModuleParams struct {
	Generator      ModuleGenerator
	CleanupHandler CleanupHandler
	PastMsgHandler PastMessagesHandler
}

func DefaultParams(generator ModuleGenerator) ModuleParams {
	return ModuleParams{
		Generator:      generator,
		CleanupHandler: nilCleanup,
		PastMsgHandler: nilPastMsgHandler,
	}
}

func nilCleanup(_ uint64) (events.EventList, error) {
	return events.EmptyList(), nil
}

func nilPastMsgHandler(_ []*modringpbtypes.PastMessage) (events.EventList, error) {
	return events.EmptyList(), nil
}
