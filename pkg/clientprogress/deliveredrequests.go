package clientprogress

import (
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/pb/commonpb"
	commonpbtypes "github.com/filecoin-project/mir/pkg/pb/commonpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

// DeliveredReqs tracks the watermarks of delivered transactions for a single client.
type DeliveredReqs struct {

	// LowWM is the lowest watermark of the client.
	lowWM tt.ReqNo

	// Set of request numbers indicating whether the corresponding transaction has been delivered.
	// In addition, all request numbers strictly smaller than LowWM are considered delivered.
	delivered map[tt.ReqNo]struct{}

	// Logger to use for logging output.
	logger logging.Logger
}

// EmptyDeliveredReqs allocates and returns a new DeliveredReqs.
func EmptyDeliveredReqs(logger logging.Logger) *DeliveredReqs {
	return NewDeliveredReqs(logger, 1)
}

// NewDeliveredReqs allocates (with the given capacity for delivered requests) and returns a new DeliveredReqs.
func NewDeliveredReqs(logger logging.Logger, capacity int) *DeliveredReqs {
	return &DeliveredReqs{
		lowWM:     0,
		delivered: make(map[tt.ReqNo]struct{}, capacity),
		logger:    logger,
	}
}

func DeliveredReqsFromPb(pb *commonpb.DeliveredReqs, logger logging.Logger) *DeliveredReqs {
	dr := NewDeliveredReqs(logger, len(pb.Delivered))
	dr.lowWM = tt.ReqNo(pb.LowWm)
	for _, reqNo := range pb.Delivered {
		dr.delivered[tt.ReqNo(reqNo)] = struct{}{}
	}
	return dr
}

func DeliveredReqsFromDslStruct(ds *commonpbtypes.DeliveredReqs, logger logging.Logger) *DeliveredReqs {
	dr := NewDeliveredReqs(logger, len(ds.Delivered))
	dr.lowWM = ds.LowWm
	for _, reqNo := range ds.Delivered {
		dr.delivered[reqNo] = struct{}{}
	}
	return dr
}

// Add adds a request number that is considered delivered to the DeliveredReqs.
// Returns true if the request number has been added now (after not being previously present).
// Returns false if the request number has already been added before the call to Add.
func (dr *DeliveredReqs) Add(reqNo tt.ReqNo) bool {
	if dr.CanAdd(reqNo) {
		dr.delivered[reqNo] = struct{}{}
		return true
	}

	return false
}

// CanAdd checks if a request number is considered delivered in the DeliveredReqs.
// Returns true if the request number has not been added previously.
// Returns false if the request number has already been added before the call to CanAdd.
func (dr *DeliveredReqs) CanAdd(reqNo tt.ReqNo) bool {
	if reqNo < dr.lowWM {
		dr.logger.Log(logging.LevelDebug, "Request sequence number below client's watermark window.",
			"lowWM", dr.lowWM, "reqNo", reqNo)
		return false
	}

	_, alreadyPresent := dr.delivered[reqNo]
	return !alreadyPresent
}

// GarbageCollect reduces the memory footprint of the DeliveredReqs
// by deleting a contiguous prefix of delivered request numbers
// and increasing the low watermark accordingly.
// Returns the new low watermark.
func (dr *DeliveredReqs) GarbageCollect() tt.ReqNo {

	for _, ok := dr.delivered[dr.lowWM]; ok; _, ok = dr.delivered[dr.lowWM] {
		delete(dr.delivered, dr.lowWM)
		dr.lowWM++
	}

	return dr.lowWM
}

func (dr *DeliveredReqs) Pb() *commonpb.DeliveredReqs {
	delivered := make([]uint64, len(dr.delivered))
	for i, reqNo := range maputil.GetSortedKeys(dr.delivered) {
		delivered[i] = reqNo.Pb()
	}

	return &commonpb.DeliveredReqs{
		LowWm:     dr.lowWM.Pb(),
		Delivered: delivered,
	}
}

func (dr *DeliveredReqs) DslStruct() *commonpbtypes.DeliveredReqs {
	return &commonpbtypes.DeliveredReqs{
		LowWm:     dr.lowWM,
		Delivered: maputil.GetSortedKeys(dr.delivered),
	}
}
