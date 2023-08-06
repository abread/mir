package clientprogress

import (
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/pb/trantorpb"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

// DeliveredTxs tracks the watermarks of delivered transactions for a single client.
type DeliveredTXs struct {

	// LowWM is the lowest watermark of the client.
	lowWM tt.TxNo

	// Set of transaction numbers indicating whether the corresponding transaction has been delivered.
	// In addition, all transaction numbers strictly smaller than LowWM are considered delivered.
	delivered map[tt.TxNo]struct{}

	// Logger to use for logging output.
	logger logging.Logger
}

// EmptyDeliveredReqs allocates and returns a new DeliveredReqs.
func EmptyDeliveredTXs(logger logging.Logger) *DeliveredTXs {
	return NewDeliveredTXs(logger, 1)
}

// NewDeliveredReqs allocates (with the given capacity for delivered requests) and returns a new DeliveredReqs.
func NewDeliveredTXs(logger logging.Logger, capacity int) *DeliveredTXs {
	return &DeliveredTXs{
		lowWM:     0,
		delivered: make(map[tt.TxNo]struct{}, capacity),
		logger:    logger,
	}
}

func DeliveredTXsFromPb(pb *trantorpb.DeliveredTXs, logger logging.Logger) *DeliveredTXs {
	dt := NewDeliveredTXs(logger, len(pb.Delivered))
	dt.lowWM = tt.TxNo(pb.LowWm)
	for _, txNo := range pb.Delivered {
		dt.delivered[tt.TxNo(txNo)] = struct{}{}
	}
	return dt
}

func DeliveredTXsFromDslStruct(ds *trantorpbtypes.DeliveredTXs, logger logging.Logger) *DeliveredTXs {
	dt := NewDeliveredTXs(logger, len(ds.Delivered))
	dt.lowWM = ds.LowWm
	for _, txNo := range ds.Delivered {
		dt.delivered[txNo] = struct{}{}
	}
	return dt
}

// Contains returns true if the given txNo has already been added.
func (dt *DeliveredTXs) Contains(txNo tt.TxNo) bool {

	if txNo < dt.lowWM {
		return true
	}

	_, alreadyPresent := dt.delivered[txNo]
	return alreadyPresent
}

// Add adds a transaction number that is considered delivered to the DeliveredTXs.
// Returns true if the transaction number has been added now (after not being previously present).
// Returns false if the transaction number has already been added before the call to Add.
func (dt *DeliveredTXs) Add(txNo tt.TxNo) bool {
	if !dt.Contains(txNo) {
		dt.delivered[txNo] = struct{}{}
		return true
	}

	return false
}

// GarbageCollect reduces the memory footprint of the DeliveredTXs
// by deleting a contiguous prefix of delivered transaction numbers
// and increasing the low watermark accordingly.
// Returns the new low watermark.
func (dt *DeliveredTXs) GarbageCollect() tt.TxNo {

	for _, ok := dt.delivered[dt.lowWM]; ok; _, ok = dt.delivered[dt.lowWM] {
		delete(dt.delivered, dt.lowWM)
		dt.lowWM++
	}

	return dt.lowWM
}

func (dt *DeliveredTXs) Pb() *trantorpb.DeliveredTXs {
	delivered := make([]uint64, len(dt.delivered))
	for i, txNo := range maputil.GetSortedKeys(dt.delivered) {
		delivered[i] = txNo.Pb()
	}

	return &trantorpb.DeliveredTXs{
		LowWm:     dt.lowWM.Pb(),
		Delivered: delivered,
	}
}

func (dt *DeliveredTXs) DslStruct() *trantorpbtypes.DeliveredTXs {
	return &trantorpbtypes.DeliveredTXs{
		LowWm:     dt.lowWM,
		Delivered: maputil.GetSortedKeys(dt.delivered),
	}
}
