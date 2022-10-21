package alea

import (
	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type SlotID uint64

func (s SlotID) Pb() uint64 {
	return uint64(s)
}

// message id without protocol buf trash
type MsgID struct {
	ProposerID t.NodeID
	Slot       SlotID
}

func (m *MsgID) Pb() *aleapb.MsgId {
	return &aleapb.MsgId{
		QueueIdx: m.ProposerID.Pb(),
		Slot:     m.Slot.Pb(),
	}
}

func MsgIDFromDomain(pb *aleapb.MsgId) MsgID {
	return MsgID{
		ProposerID: t.NodeID(pb.QueueIdx),
		Slot:       SlotID(pb.Slot),
	}
}
