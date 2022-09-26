package vcbc

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/alea/util"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type VCBCConfig struct {
	Id        *aleapb.MsgId
	Members   []t.NodeID
	Threshold int

	SelfModuleID         t.ModuleID
	AleaModuleID         t.ModuleID
	NetModuleID          t.ModuleID
	ThreshCryptoModuleID t.ModuleID
}

type VCBCModuleState[PS any] struct {
	config        VCBCConfig
	protocolState PS
	payload       *requestpb.Batch
}

type void struct{}
type verifyShareCtx struct {
	node     t.NodeID
	sigShare []byte
}
type fullSig []byte

type vcbcReceiverProtocolStateImpl interface {
	UponVCBCSend(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCSend) error
	UponVCBCFinal(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCFinal) error

	UponSignShareResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *void, sigShare []byte) error
	UponVerifyFullResult(m dsl.Module, ctx *VCBCModuleState[vcbcReceiverProtocolStateImpl], opCtx *fullSig, ok bool, err string) error
}

type vcbcSenderProtocolStateImpl interface {
	UponVCBCEcho(m dsl.Module, ctx *VCBCModuleState[vcbcSenderProtocolStateImpl], from t.NodeID, msg *aleapb.VCBCEcho) error

	HasReachedStrongSupport(m dsl.Module, ctx *VCBCModuleState[vcbcSenderProtocolStateImpl]) bool
	UponReachedStrongSupport(m dsl.Module, ctx *VCBCModuleState[vcbcSenderProtocolStateImpl]) error

	UponVerifyShareResult(m dsl.Module, ctx *VCBCModuleState[vcbcSenderProtocolStateImpl], opCtx *verifyShareCtx, ok bool, err string) error
	UponRecoverResult(m dsl.Module, ctx *VCBCModuleState[vcbcSenderProtocolStateImpl], opCtx *verifyShareCtx, ok bool, fullSig []byte, err string) error
}

func NewVCBCReceiver(config VCBCConfig) modules.PassiveModule {
	m := dsl.NewModule(config.SelfModuleID)
	moduleState := &VCBCModuleState[vcbcReceiverProtocolStateImpl]{
		config:        config,
		protocolState: &vcbcReceiverInit{},
		payload:       nil,
	}

	uponVcbcMessage(m, config.Id, func(from t.NodeID, msgWrapped *aleapb.VCBC) error {
		switch msg := msgWrapped.Type.(type) {
		case *aleapb.VCBC_Send:
			return moduleState.protocolState.UponVCBCSend(m, moduleState, from, msg.Send)
		case *aleapb.VCBC_Echo:
			// TODO: is this the right thing?
			return fmt.Errorf("unexpected VCBC protocol message: %T", msgWrapped.Type)
		case *aleapb.VCBC_Final:
			return moduleState.protocolState.UponVCBCFinal(m, moduleState, from, msg.Final)
		default:
			panic(fmt.Sprintf("unknown VCBC protocol message: %T", msgWrapped.Type))
		}
	})

	// TODO: missing handlers

	return m
}

func NewVCBCSender(config VCBCConfig, payload *requestpb.Batch) modules.PassiveModule {
	m := dsl.NewModule(config.SelfModuleID)
	protocolState := &vcbcSenderInit{}
	moduleState := &VCBCModuleState[vcbcSenderProtocolStateImpl]{
		config:        config,
		protocolState: protocolState,
		payload:       payload,
	}

	uponVcbcMessage(m, config.Id, func(from t.NodeID, msgWrapped *aleapb.VCBC) error {
		switch msg := msgWrapped.Type.(type) {
		case *aleapb.VCBC_Send:
			// TODO: is this the right thing?
			return fmt.Errorf("unexpected VCBC protocol message: %T", msgWrapped.Type)
		case *aleapb.VCBC_Echo:
			return moduleState.protocolState.UponVCBCEcho(m, moduleState, from, msg.Echo)
		case *aleapb.VCBC_Final:
			// TODO: is this the right thing?
			return fmt.Errorf("unexpected VCBC protocol message: %T", msgWrapped.Type)
		default:
			panic(fmt.Sprintf("unknown VCBC protocol message: %T", msgWrapped.Type))
		}
	})

	// TODO: missing handlers

	dsl.UponCondition(m, func() error {
		if moduleState.protocolState.HasReachedStrongSupport(m, moduleState) {
			return moduleState.protocolState.UponReachedStrongSupport(m, moduleState)
		}

		return nil
	})

	protocolState.UponInit(m, moduleState)

	return m
}

func uponVcbcMessage(m dsl.Module, id *aleapb.MsgId, handler func(from t.NodeID, msg *aleapb.VCBC) error) {
	dsl.UponEvent[*eventpb.Event_Alea](m, func(ev *aleapb.Event) error {
		evWrapper, ok := ev.Type.(aleapb.Event_TypeWrapper[aleapb.VCBCMessageRecvd])
		if !ok {
			return nil
		}

		unwrapped := evWrapper.Unwrap()
		if unwrapped.InstanceId != id {
			panic("VCBC message routed to wrong VCBC instance")
		}

		return handler(t.NodeID(unwrapped.From), unwrapped.Message)
	})
}

func (vcbc *VCBCModuleState[any]) dataToSign() [][]byte {
	return util.SlotMessageData(vcbc.config.Id, vcbc.payload)
}
