package availability

import (
	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
	t "github.com/filecoin-project/mir/pkg/types"
)

func Include(m dsl.Module, mc common.ModuleConfig, params common.ModuleParams, tunables common.ModuleTunables, nodeID t.NodeID, logger logging.Logger) {
	IncludeBatchFetching(m, mc, params, tunables, nodeID, logging.Decorate(logger, "Batch Feching: "))
	IncludeCreatingCertificates(m, mc, params, nodeID)
	IncludeVerificationOfCertificates(m, mc, params, nodeID)
}
