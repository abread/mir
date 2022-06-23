package alea

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/pb/aleapb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type SlotId uint64

func (s SlotId) Pb() uint64 {
	return uint64(s)
}

// message id without protocol buf trash
type MsgId struct {
	ProposerId t.NodeID
	Slot       SlotId
}

func (m *MsgId) Pb() *aleapb.MsgId {
	return &aleapb.MsgId{
		QueueIdx: m.ProposerId.Pb(),
		Slot:     m.Slot.Pb(),
	}
}

func MsgIdFromDomain(pb *aleapb.MsgId) MsgId {
	return MsgId{
		ProposerId: t.NodeID(pb.QueueIdx),
		Slot:       SlotId(pb.Slot),
	}
}

const THRESH_SIG_SHARE_SIZE = 1 // TODO
type threshSigShare [THRESH_SIG_SHARE_SIZE]byte

func (s threshSigShare) Pb() []byte {
	a := [THRESH_SIG_SHARE_SIZE]byte(s)
	return a[:]
}

func threshSigShareFromDomain(slice []byte) (threshSigShare, error) {
	if len(slice) == THRESH_SIG_SHARE_SIZE {
		var s [THRESH_SIG_SHARE_SIZE]byte
		copy(s[:], slice)
		return threshSigShare(s), nil
	} else {
		return threshSigShare{}, fmt.Errorf("bad threshold threshSignature share size. expected %d got %d", THRESH_SIG_SHARE_SIZE, len(slice))
	}

}

const THRESH_SIG_FULL_SIZE = 1 // TODO
type threshSigFull [THRESH_SIG_FULL_SIZE]byte

func (s threshSigFull) Pb() []byte {
	a := [THRESH_SIG_FULL_SIZE]byte(s)
	return a[:]
}

func threshSigFullFromDomain(slice []byte) (threshSigFull, error) {
	if len(slice) == THRESH_SIG_FULL_SIZE {
		var s [THRESH_SIG_FULL_SIZE]byte
		copy(s[:], slice)
		return threshSigFull(s), nil
	} else {
		return threshSigFull{}, fmt.Errorf("bad threshold threshSignature share size. expected %d got %d", THRESH_SIG_FULL_SIZE, len(slice))
	}

}
