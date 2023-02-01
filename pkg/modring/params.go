package modring

import (
	"github.com/filecoin-project/mir/pkg/modules"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleGenerator func(id t.ModuleID) (modules.PassiveModule, error)

type ModuleParams struct {
	Generator ModuleGenerator
}

func DefaultParams(generator ModuleGenerator) ModuleParams {
	return ModuleParams{
		Generator: generator,
	}
}
