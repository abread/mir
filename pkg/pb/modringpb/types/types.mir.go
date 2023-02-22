package modringpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	modringpb "github.com/filecoin-project/mir/pkg/pb/modringpb"
	types "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type PastMessage struct {
	DestId  uint64
	From    types.NodeID
	Message *types1.Message
}

func PastMessageFromPb(pb *modringpb.PastMessage) *PastMessage {
	return &PastMessage{
		DestId:  pb.DestId,
		From:    (types.NodeID)(pb.From),
		Message: types1.MessageFromPb(pb.Message),
	}
}

func (m *PastMessage) Pb() *modringpb.PastMessage {
	return &modringpb.PastMessage{
		DestId:  m.DestId,
		From:    (string)(m.From),
		Message: (m.Message).Pb(),
	}
}

func (*PastMessage) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*modringpb.PastMessage]()}
}
