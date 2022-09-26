package vcbc

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	threshDsl "github.com/filecoin-project/mir/pkg/threshcrypto/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

type vcbcReceiverInit struct{}

func (v *vcbcReceiverInit) UponVCBCSend(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCSend) error {
	if ctx.config.Id.QueueIdx != from.Pb() {
		// TODO: signal byz behavior
		// TODO: consider moving this check to VCBCModuleState
		return nil
	}

	// stop advancing state while waiting for thresh signature
	ctx.protocolState = &vcbcReceiverAwaitingSignShare{}

	ctx.payload = msg.Payload
	threshDsl.SignShare(m, ctx.config.ThreshCryptoModuleID, ctx.dataToSign(), &void{})

	return nil
}

func (v *vcbcReceiverInit) UponVCBCFinal(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCFinal) error {
	return nil // noop
}

func (v *vcbcReceiverInit) UponSignShareResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *void, sigShare []byte) error {
	return nil // noop
}

func (v *vcbcReceiverInit) UponVerifyFullResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *fullSig, ok bool, err string) error {
	return nil // noop
}

type vcbcReceiverAwaitingSignShare struct{}

func (v *vcbcReceiverAwaitingSignShare) UponVCBCSend(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCSend) error {
	return nil // noop
}

func (v *vcbcReceiverAwaitingSignShare) UponVCBCFinal(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCFinal) error {
	return nil
}

func (v *vcbcReceiverAwaitingSignShare) UponSignShareResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *void, sigShare []byte) error {
	sendVCBCEcho(m, ctx, t.NodeID(ctx.config.Id.QueueIdx), sigShare)

	ctx.protocolState = &vcbcReceiverEchoed{}

	return nil
}

func (v *vcbcReceiverAwaitingSignShare) UponVerifyFullResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *fullSig, ok bool, err string) error {
	return nil // noop
}

type vcbcReceiverEchoed struct{}

func (v *vcbcReceiverEchoed) UponVCBCSend(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCSend) error {
	return nil // noop
}

func (v *vcbcReceiverEchoed) UponVCBCFinal(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCFinal) error {
	opCtx := fullSig(msg.Signature)
	threshDsl.VerifyFull(m, ctx.config.ThreshCryptoModuleID, ctx.dataToSign(), msg.Signature, &opCtx)

	ctx.protocolState = &vcbcReceiverAwaitingVerifyFull{}

	return nil
}

func (v *vcbcReceiverEchoed) UponSignShareResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *void, sigShare []byte) error {
	return nil // noop
}

func (v *vcbcReceiverEchoed) UponVerifyFullResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *fullSig, ok bool, err string) error {
	return nil // noop
}

type vcbcReceiverAwaitingVerifyFull struct{}

func (v *vcbcReceiverAwaitingVerifyFull) UponVCBCSend(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCSend) error {
	return nil // noop
}

func (v *vcbcReceiverAwaitingVerifyFull) UponVCBCFinal(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCFinal) error {
	opCtx := fullSig(msg.Signature)
	threshDsl.VerifyFull(m, ctx.config.ThreshCryptoModuleID, ctx.dataToSign(), msg.Signature, &opCtx)

	ctx.protocolState = &vcbcReceiverAwaitingVerifyFull{}

	return nil
}

func (v *vcbcReceiverAwaitingVerifyFull) UponSignShareResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *void, sigShare []byte) error {
	return nil // noop
}

func (v *vcbcReceiverAwaitingVerifyFull) UponVerifyFullResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *fullSig, ok bool, err string) error {
	if ok {
		fullSig := []byte(*opCtx)

		deliverBatch(m, ctx, ctx.payload, fullSig)

		ctx.protocolState = &vcbcReceiverDone{}
		return nil
	}

	// TODO: return err? emit event for instance to be cleaned up?
	// this **is** byz behavior
	return nil
}

type vcbcReceiverDone struct{}

func (v *vcbcReceiverDone) UponVCBCSend(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCSend) error {
	return nil // noop
}

func (v *vcbcReceiverDone) UponVCBCFinal(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCFinal) error {
	return nil // noop
}

func (v *vcbcReceiverDone) UponSignShareResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *void, sigShare []byte) error {
	return nil // noop
}

func (v *vcbcReceiverDone) UponVerifyFullResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *fullSig, ok bool, err string) error {
	return nil // noop
}
