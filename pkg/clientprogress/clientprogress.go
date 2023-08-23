package clientprogress

import (
	"github.com/filecoin-project/mir/pkg/pb/trantorpb"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
)

// ClientProgress tracks watermarks for all the clients.
type ClientProgress struct {
	ClientTrackers map[tt.ClientID]*DeliveredTXs
}

func NewClientProgress() *ClientProgress {
	return &ClientProgress{
		ClientTrackers: make(map[tt.ClientID]*DeliveredTXs),
	}
}

func (cp *ClientProgress) Contains(clID tt.ClientID, txNo tt.TxNo) bool {
	if _, ok := cp.ClientTrackers[clID]; !ok {
		return false
	}
	return cp.ClientTrackers[clID].Contains(txNo)
}

func (cp *ClientProgress) Add(clID tt.ClientID, txNo tt.TxNo) bool {
	if _, ok := cp.ClientTrackers[clID]; !ok {
		cp.ClientTrackers[clID] = EmptyDeliveredTXs()
	}
	return cp.ClientTrackers[clID].Add(txNo)
}

func (cp *ClientProgress) IsBelowWatermarkWindow(clID tt.ClientID, txNo tt.TxNo) bool {
	if ct, ok := cp.ClientTrackers[clID]; ok {
		return ct.IsBelowWatermarkWindow(txNo)
	}
	return false
}

func (cp *ClientProgress) GarbageCollectGetWatermarks() map[tt.ClientID]tt.TxNo {
	lowWMs := make(map[tt.ClientID]tt.TxNo)
	for clientID, cwmt := range cp.ClientTrackers {
		lowWMs[clientID] = cwmt.GarbageCollect()
	}
	return lowWMs
}

func (cp *ClientProgress) GarbageCollect() {
	for _, cwmt := range cp.ClientTrackers {
		cwmt.GarbageCollect()
	}
}

func (cp *ClientProgress) DslStruct() *trantorpbtypes.ClientProgress {
	ds := make(map[tt.ClientID]*trantorpbtypes.DeliveredTXs)
	for clientID, clientTracker := range cp.ClientTrackers {
		ds[clientID] = clientTracker.DslStruct()
	}
	return &trantorpbtypes.ClientProgress{Progress: ds}
}

func (cp *ClientProgress) LoadDslStruct(ds *trantorpbtypes.ClientProgress) {
	cp.ClientTrackers = make(map[tt.ClientID]*DeliveredTXs)
	for clientID, deliveredReqs := range ds.Progress {
		cp.ClientTrackers[clientID] = DeliveredTXsFromDslStruct(deliveredReqs)
	}
}

func (cp *ClientProgress) Pb() *trantorpb.ClientProgress {
	pb := make(map[string]*trantorpb.DeliveredTXs)
	for clientID, clientTracker := range cp.ClientTrackers {
		pb[clientID.Pb()] = clientTracker.Pb()
	}
	return &trantorpb.ClientProgress{Progress: pb}
}

func (cp *ClientProgress) LoadPb(pb *trantorpb.ClientProgress) {
	cp.ClientTrackers = make(map[tt.ClientID]*DeliveredTXs)
	for clientID, deliveredTXs := range pb.Progress {
		cp.ClientTrackers[tt.ClientID(clientID)] = DeliveredTXsFromPb(deliveredTXs)
	}
}

func FromPb(pb *trantorpb.ClientProgress) *ClientProgress {
	cp := NewClientProgress()
	cp.LoadPb(pb)
	return cp
}
