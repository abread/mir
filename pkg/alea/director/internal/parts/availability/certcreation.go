package availability

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	adsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/dsl"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

// IncludeCreatingCertificates registers event handlers for processing availabilitypb.RequestCert events.
func IncludeCreatingCertificates(
	m dsl.Module,
	mc common.ModuleConfig,
	params common.ModuleParams,
	nodeID t.NodeID,
) {
	// When a batch is requested by the consensus layer, request a batch of transactions from the mempool.
	adsl.UponRequestCert(m, func(origin *availabilitypbtypes.RequestCertOrigin) error {
		return fmt.Errorf("operation not supported")
	})
}
