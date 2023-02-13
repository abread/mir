package reliablenetpbevents

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types3 "github.com/filecoin-project/mir/pkg/pb/reliablenetpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func SendMessage(destModule types.ModuleID, msgId []uint8, msg *types1.Message, destinations []types.NodeID) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ReliableNet{
			ReliableNet: &types3.Event{
				Type: &types3.Event_SendMessage{
					SendMessage: &types3.SendMessage{
						MsgId:        msgId,
						Msg:          msg,
						Destinations: destinations,
					},
				},
			},
		},
	}
}

func Ack(destModule types.ModuleID, destModule0 types.ModuleID, msgId []uint8, source types.NodeID) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ReliableNet{
			ReliableNet: &types3.Event{
				Type: &types3.Event_Ack{
					Ack: &types3.Ack{
						DestModule: destModule0,
						MsgId:      msgId,
						Source:     source,
					},
				},
			},
		},
	}
}

func MarkRecvd(destModule types.ModuleID, destModule0 types.ModuleID, msgId []uint8, destinations []types.NodeID) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ReliableNet{
			ReliableNet: &types3.Event{
				Type: &types3.Event_MarkRecvd{
					MarkRecvd: &types3.MarkRecvd{
						DestModule:   destModule0,
						MsgId:        msgId,
						Destinations: destinations,
					},
				},
			},
		},
	}
}

func MarkModuleMsgsRecvd(destModule types.ModuleID, destModule0 types.ModuleID, destinations []types.NodeID) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ReliableNet{
			ReliableNet: &types3.Event{
				Type: &types3.Event_MarkModuleMsgsRecvd{
					MarkModuleMsgsRecvd: &types3.MarkModuleMsgsRecvd{
						DestModule:   destModule0,
						Destinations: destinations,
					},
				},
			},
		},
	}
}

func RetransmitAll(destModule types.ModuleID) *types2.Event {
	return &types2.Event{
		DestModule: destModule,
		Type: &types2.Event_ReliableNet{
			ReliableNet: &types3.Event{
				Type: &types3.Event_RetransmitAll{
					RetransmitAll: &types3.RetransmitAll{},
				},
			},
		},
	}
}
