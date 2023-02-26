package commontypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	aleatypes "github.com/filecoin-project/mir/pkg/alea/aleatypes"
	common "github.com/filecoin-project/mir/pkg/pb/aleapb/common"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Slot struct {
	QueueIdx  aleatypes.QueueIdx
	QueueSlot aleatypes.QueueSlot
}

func SlotFromPb(pb *common.Slot) *Slot {
	return &Slot{
		QueueIdx:  (aleatypes.QueueIdx)(pb.QueueIdx),
		QueueSlot: (aleatypes.QueueSlot)(pb.QueueSlot),
	}
}

func (m *Slot) Pb() *common.Slot {
	return &common.Slot{
		QueueIdx:  (uint32)(m.QueueIdx),
		QueueSlot: (uint64)(m.QueueSlot),
	}
}

func (*Slot) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*common.Slot]()}
}