package availability

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/alea/director/internal/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	adsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/dsl"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

// IncludeVerificationOfCertificates registers event handlers for processing availabilitypb.VerifyCert events.
func IncludeVerificationOfCertificates(
	m dsl.Module,
	mc common.ModuleConfig,
	params common.ModuleParams,
	nodeID t.NodeID,
) {
	// When receive a request to verify a certificate, check that it is structurally correct and verify the signatures.
	adsl.UponVerifyCert(m, func(genericCert *availabilitypbtypes.Cert, origin *availabilitypbtypes.VerifyCertOrigin) error {
		return fmt.Errorf("operation not supported")
	})
}
