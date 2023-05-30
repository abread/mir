package availability

import (
	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
)

func Include(m dsl.Module, mc common.ModuleConfig, params common.ModuleParams, tunables common.ModuleTunables, logger logging.Logger) {
	IncludeBatchFetching(m, mc, params, tunables, logging.Decorate(logger, "Batch Feching: "))
	IncludeCreatingCertificates(m)
	IncludeVerificationOfCertificates(m)
}
