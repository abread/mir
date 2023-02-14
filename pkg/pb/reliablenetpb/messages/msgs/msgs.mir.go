package messagesmsgs

import (
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types2 "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/messages/types"
	rntypes "github.com/filecoin-project/mir/pkg/reliablenet/rntypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

func AckMessage(destModule types.ModuleID, msgDestModule types.ModuleID, msgId rntypes.MsgID) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_ReliableNet{
			ReliableNet: &types2.Message{
				Type: &types2.Message_Ack{
					Ack: &types2.AckMessage{
						MsgDestModule: msgDestModule,
						MsgId:         msgId,
					},
				},
			},
		},
	}
}
