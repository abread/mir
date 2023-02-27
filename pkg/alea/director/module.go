package director

import (
	"context"

	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/parts/availability"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/parts/general"
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
		Self:          "alea_dir",
		AleaBroadcast: "alea_bc",
		AleaAgreement: "alea_ag",
		BatchDB:       "batchdb",
		Mempool:       "mempool",
		ReliableNet:   "reliablenet",
		ThreshCrypto:  "threshcrypto",
		Timer:         "timer",
	}
}

func NewModule(ctx context.Context, mc *ModuleConfig, params *ModuleParams, tunables *ModuleTunables, nodeID t.NodeID, logger logging.Logger) dsl.Module {
	m := dsl.NewModule(ctx, mc.Self)

	general.Include(m, mc, params, tunables, nodeID, logger)
	availability.Include(m, mc, params, tunables, nodeID, logger)

	return m
}
