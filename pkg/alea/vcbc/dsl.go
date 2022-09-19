package vcbc

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func bcVCBCSend[PS any](m dsl.Module, state *VCBCModuleState[PS], batch *requestpb.Batch) {
	broadcastMessage(m, state, &aleapb.VCBC{
		Type: &aleapb.VCBC_Send{
			Send: &aleapb.VCBCSend{
				Payload: batch,
			},
		},
	})
}

func sendVCBCEcho[PS any](m dsl.Module, state *VCBCModuleState[PS], dest t.NodeID, batch *requestpb.Batch, sigShare []byte) {
	unicastMessage(m, state, dest, &aleapb.VCBC{
		Type: &aleapb.VCBC_Echo{
			Echo: &aleapb.VCBCEcho{
				Payload:        batch,
				SignatureShare: sigShare,
			},
		},
	})
}

func bcVCBCFinal[PS any](m dsl.Module, state *VCBCModuleState[PS], batch *requestpb.Batch, sigFull []byte) {
	broadcastMessage(m, state, &aleapb.VCBC{
		Type: &aleapb.VCBC_Final{
			Final: &aleapb.VCBCFinal{
				Payload:   batch,
				Signature: sigFull,
			},
		},
	})
}

func broadcastMessage[PS any](m dsl.Module, state *VCBCModuleState[PS], msg *aleapb.VCBC) {
	msgPrepared := prepareMessage(state, msg)
	// TODO: use better broadcast
	dsl.SendMessage(m, state.config.NetModuleID, msgPrepared, state.config.Members)
}

func unicastMessage[PS any](m dsl.Module, state *VCBCModuleState[PS], dest t.NodeID, msg *aleapb.VCBC) {
	msgPrepared := prepareMessage(state, msg)
	dsl.SendMessage(m, state.config.NetModuleID, msgPrepared, []t.NodeID{dest})
}

func prepareMessage[PS any](state *VCBCModuleState[PS], msg *aleapb.VCBC) *messagepb.Message {
	return &messagepb.Message{
		Type: &messagepb.Message_Alea{
			Alea: &aleapb.AleaMessage{
				Type: &aleapb.AleaMessage_Broadcast{
					Broadcast: &aleapb.BroadcastMsg{
						InstanceId: state.config.Id,
						Message:    msg,
					},
				},
			},
		},
	}
}
