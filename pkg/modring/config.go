package modring

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleConfig struct {
	Self            t.ModuleID
	PastMsgConsumer t.ModuleID
}

type ModuleGenerator func(id t.ModuleID, idx uint64) (modules.PassiveModule, *events.EventList, error)

type ModuleParams struct {
	Generator ModuleGenerator
}

func DefaultParams(generator ModuleGenerator) ModuleParams {
	return ModuleParams{
		Generator: generator,
	}
}
