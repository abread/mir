package availability

import (
	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/dsl"

	"github.com/filecoin-project/mir/pkg/alea/broadcast/bccommon"
	availabilitypbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/dsl"
	availabilitypbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	threshcryptopbdsl "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/dsl"
	"github.com/filecoin-project/mir/pkg/vcb"
)

func includeCertVerification(m dsl.Module, mc ModuleConfig, params ModuleParams) {
	availabilitypbdsl.UponVerifyCert(m, func(cert *availabilitypbtypes.Cert, origin *availabilitypbtypes.VerifyCertOrigin) error {
		aleaCertWrapped, ok := cert.Type.(*availabilitypbtypes.Cert_Alea)
		if !ok {
			return es.Errorf("alea can only validate alea certs")
		}
		aleaCert := aleaCertWrapped.Alea

		sigData := vcb.SigData(
			bccommon.VCBInstanceUID(params.AleaInstanceUID, aleaCert.Slot.QueueIdx, aleaCert.Slot.QueueSlot),
			aleaCert.BatchId,
		)
		threshcryptopbdsl.VerifyFull(m, mc.ThreshCrypto, sigData, aleaCert.Signature, origin)
		return nil
	})
	threshcryptopbdsl.UponVerifyFullResult(m, func(ok bool, err string, origin *availabilitypbtypes.VerifyCertOrigin) error {
		availabilitypbdsl.CertVerified(m, origin.Module, ok, err, origin)
		return nil
	})
}
