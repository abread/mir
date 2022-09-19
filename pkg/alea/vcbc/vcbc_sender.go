package vcbc

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	threshDsl "github.com/filecoin-project/mir/pkg/threshcrypto/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

type vcbcSenderInit struct {
	pendingVerification map[t.NodeID]void
	validShares         map[t.NodeID][]byte
}

func (v *vcbcSenderInit) UponInit(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl]) {
	bcVCBCSend(m, moduleState, moduleState.payload)
}

func (v *vcbcSenderInit) UponVCBCEcho(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCEcho) error {
	_, pendingVerification := v.pendingVerification[from]
	if pendingVerification {
		// TODO: notify byz behavior (?)
		return nil
	}

	v.pendingVerification[from] = void{}

	threshDsl.VerifyShare(m, moduleState.config.SelfModuleID, moduleState.dataToSign(), msg.SignatureShare, from, &verifyShareCtx{
		node:     from,
		sigShare: msg.SignatureShare,
	})
	return nil
}

func (v *vcbcSenderInit) UponVerifyShareResult(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl], opCtx *verifyShareCtx, ok bool, err string) error {
	if ok {
		v.validShares[opCtx.node] = opCtx.sigShare
		return nil
	}

	// TODO: return err? go to err state? emit event for instance to be cleaned up?
	// this **is** byz behavior
	return nil
}

func (v *vcbcSenderInit) HasReachedStrongSupport(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl]) bool {
	return len(v.validShares) >= moduleState.config.Threshold
}

func (v *vcbcSenderInit) UponReachedStrongSupport(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl]) error {
	shares := make([][]byte, 0, len(v.validShares))
	for _, share := range v.validShares {
		shares = append(shares, share)
	}

	threshDsl.Recover(m, moduleState.config.SelfModuleID, moduleState.dataToSign(), shares, &void{})
	return nil
}

func (v *vcbcSenderInit) UponRecoverResult(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl], opCtx *verifyShareCtx, ok bool, fullSig []byte, err string) error {
	if ok {
		moduleState.signature = fullSig

		bcVCBCFinal(m, moduleState, moduleState.payload, fullSig)
		// TODO: emit delivery event

		moduleState.protocolState = &vcbcSenderDone{}
		return nil
	}
	// will retry after getting more signature shares, can't tell who's being byzantine

	return nil
}

type vcbcSenderDone struct{}

func (v *vcbcSenderDone) UponVCBCEcho(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCEcho) error {
	return nil // noop
}

func (v *vcbcSenderDone) UponVerifyShareResult(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl], opCtx *verifyShareCtx, ok bool, err string) error {
	return nil // noop
}

func (v *vcbcSenderDone) HasReachedStrongSupport(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl]) bool {
	return false // noop
}

func (v *vcbcSenderDone) UponReachedStrongSupport(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl]) error {
	return nil // noop
}

func (v *vcbcSenderDone) UponRecoverResult(m dsl.Module, moduleState *VCBCModuleState[vcbcSenderProtocolStateImpl], opCtx *verifyShareCtx, ok bool, fullSig []byte, err string) error {
	return nil // noop
}
