package availability

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/alea/internal/common"
	adsl "github.com/filecoin-project/mir/pkg/availability/dsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// IncludeVerificationOfCertificates registers event handlers for processing availabilitypb.VerifyCert events.
func IncludeVerificationOfCertificates(
	m dsl.Module,
	mc *common.ModuleConfig,
	params *common.ModuleParams,
	nodeID t.NodeID,
) {
	// When receive a request to verify a certificate, check that it is structurally correct and verify the signatures.
	adsl.UponVerifyCert(m, func(genericCert *apb.Cert, origin *apb.VerifyCertOrigin) error {
		return fmt.Errorf("operation not supported")
	})
}
