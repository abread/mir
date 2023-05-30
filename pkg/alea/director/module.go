package director

import (
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

func NewModule(mc ModuleConfig, params ModuleParams, tunables ModuleTunables, nodeID t.NodeID, logger logging.Logger) dsl.Module {
	m := dsl.NewModule(mc.Self)

	general.Include(m, mc, params, tunables, nodeID, logger)
	availability.Include(m, mc, params, tunables, logger)

	return m
}
