package modring

import (
	"context"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	modringpbtypes "github.com/filecoin-project/mir/pkg/pb/modringpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleConfig struct {
	Self t.ModuleID
}

type ModuleGenerator func(ctx context.Context, id t.ModuleID, idx uint64) (modules.Module, *events.EventList, error)
type PastMessagesHandler func(pastMessages []*modringpbtypes.PastMessage) (*events.EventList, error)

type ModuleParams struct {
	Generator      ModuleGenerator
	PastMsgHandler PastMessagesHandler
	InputQueueSize int
}

func DefaultParams(generator ModuleGenerator) ModuleParams {
	return ModuleParams{
		Generator:      generator,
		InputQueueSize: 8,
	}
}
