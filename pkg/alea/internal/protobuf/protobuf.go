package protobuf

import (
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	aleapbCommon "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func Message(moduleID t.ModuleID, msg *aleapb.Message) *messagepb.Message {
	return &messagepb.Message{
		DestModule: moduleID.Pb(),
		Type: &messagepb.Message_Alea{
			Alea: msg,
		},
	}
}

func FillerMessage(moduleID t.ModuleID, slot *aleapbCommon.Slot, txs []*requestpb.Request, signature []byte) *messagepb.Message {
	return Message(moduleID, &aleapb.Message{
		Type: &aleapb.Message_FillerMessage{
			FillerMessage: &aleapb.FillerMessage{
				Slot:      slot,
				Txs:       txs,
				Signature: signature,
			},
		},
	})
}

func FillGapMessage(moduleID t.ModuleID, slot *aleapbCommon.Slot) *messagepb.Message {
	return Message(moduleID, &aleapb.Message{
		Type: &aleapb.Message_FillGapMessage{
			FillGapMessage: &aleapb.FillGapMessage{
				Slot: slot,
			},
		},
	})
}

func AleaCert(slot *aleapbCommon.Slot) *aleapb.Cert {
	return &aleapb.Cert{
		Slot: slot,
	}
}

func AvailabilityCert(slot *aleapbCommon.Slot) *apb.Cert {
	return &apb.Cert{
		Type: &apb.Cert_Alea{
			Alea: AleaCert(slot),
		},
	}
}
