package availability

import (
	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/alea/director/internal/parts/estimators"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/logging"
)

func Include(m dsl.Module, mc common.ModuleConfig, params common.ModuleParams, logger logging.Logger, est *estimators.Estimators) {
	IncludeBatchFetching(m, mc, params, logging.Decorate(logger, "Batch Feching: "), est)
	IncludeCreatingCertificates(m)
	IncludeVerificationOfCertificates(m)
}
