package alea

import (
	"github.com/filecoin-project/mir/pkg/alea/internal/common"
	"github.com/filecoin-project/mir/pkg/alea/internal/parts/general"
	"github.com/filecoin-project/mir/pkg/alea/util"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleConfig = common.ModuleConfig
type ModuleParams = common.ModuleParams
type ModuleTunables = common.ModuleTunables

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig(consumer t.ModuleID) *ModuleConfig {
	return &ModuleConfig{
		Self:          "alea",
		AleaBroadcast: "alea_bc",
		AleaAgreement: "alea_ag",
		BatchDB:       "batchdb",
		Mempool:       "mempool",
		Net:           "net",
		ThreshCrypto:  "threshcrypto",
	}
}

type aleaState struct {
	ownBcQueueWinSzCtrl *util.WindowSizeController
}

func NewModule(mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) dsl.Module {
	m := dsl.NewModule(mc.Self)

	general.Include(m, mc, params, tunables, nodeID, logger)

	return m
}
