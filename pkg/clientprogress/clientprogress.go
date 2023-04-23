package clientprogress

import (
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/pb/commonpb"
	commonpbtypes "github.com/filecoin-project/mir/pkg/pb/commonpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
)

// ClientProgress tracks watermarks for all the clients.
type ClientProgress struct {
	ClientTrackers map[tt.ClientID]*DeliveredReqs
	logger         logging.Logger
}

func NewClientProgress(logger logging.Logger) *ClientProgress {
	return &ClientProgress{
		ClientTrackers: make(map[tt.ClientID]*DeliveredReqs),
		logger:         logger,
	}
}

func (cp *ClientProgress) Add(clID tt.ClientID, reqNo tt.ReqNo) bool {
	if _, ok := cp.ClientTrackers[clID]; !ok {
		cp.ClientTrackers[clID] = EmptyDeliveredReqs(logging.Decorate(cp.logger, "", "clID", clID))
	}
	return cp.ClientTrackers[clID].Add(reqNo)
}

func (cp *ClientProgress) CanAdd(clID tt.ClientID, reqNo tt.ReqNo) bool {
	if _, ok := cp.ClientTrackers[clID]; !ok {
		cp.ClientTrackers[clID] = EmptyDeliveredReqs(logging.Decorate(cp.logger, "", "clID", clID))
	}
	return cp.ClientTrackers[clID].CanAdd(reqNo)
}

func (cp *ClientProgress) GarbageCollect() map[tt.ClientID]tt.ReqNo {
	lowWMs := make(map[tt.ClientID]tt.ReqNo)
	for clientID, cwmt := range cp.ClientTrackers {
		lowWMs[clientID] = cwmt.GarbageCollect()
	}
	return lowWMs
}

func (cp *ClientProgress) DslStruct() *commonpbtypes.ClientProgress {
	ds := make(map[string]*commonpbtypes.DeliveredReqs)
	for clientID, clientTracker := range cp.ClientTrackers {
		ds[clientID.Pb()] = clientTracker.DslStruct()
	}
	return &commonpbtypes.ClientProgress{Progress: ds}
}

func (cp *ClientProgress) LoadDslStruct(ds *commonpbtypes.ClientProgress) {
	cp.ClientTrackers = make(map[tt.ClientID]*DeliveredReqs)
	for clientID, deliveredReqs := range ds.Progress {
		cp.ClientTrackers[tt.ClientID(clientID)] = DeliveredReqsFromDslStruct(
			deliveredReqs,
			logging.Decorate(cp.logger, "", "clID", clientID),
		)
	}
}

func (cp *ClientProgress) Pb() *commonpb.ClientProgress {
	pb := make(map[string]*commonpb.DeliveredReqs)
	for clientID, clientTracker := range cp.ClientTrackers {
		pb[clientID.Pb()] = clientTracker.Pb()
	}
	return &commonpb.ClientProgress{Progress: pb}
}

func (cp *ClientProgress) LoadPb(pb *commonpb.ClientProgress) {
	cp.ClientTrackers = make(map[tt.ClientID]*DeliveredReqs)
	for clientID, deliveredReqs := range pb.Progress {
		cp.ClientTrackers[tt.ClientID(clientID)] = DeliveredReqsFromPb(
			deliveredReqs,
			logging.Decorate(cp.logger, "", "clID", clientID),
		)
	}
}

func FromPb(pb *commonpb.ClientProgress, logger logging.Logger) *ClientProgress {
	cp := NewClientProgress(logger)
	cp.LoadPb(pb)
	return cp
}
