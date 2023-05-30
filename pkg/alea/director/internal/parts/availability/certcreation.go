package availability

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/dsl"
	adsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/dsl"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
)

// IncludeCreatingCertificates registers event handlers for processing availabilitypb.RequestCert events.
func IncludeCreatingCertificates(
	m dsl.Module,
) {
	// When a batch is requested by the consensus layer, request a batch of transactions from the mempool.
	adsl.UponRequestCert(m, func(origin *availabilitypbtypes.RequestCertOrigin) error {
		return fmt.Errorf("operation not supported")
	})
}
