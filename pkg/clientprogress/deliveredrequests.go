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
	dr := NewDeliveredTXs(logger, len(pb.Delivered))
	dr.lowWM = tt.TxNo(pb.LowWm)
	for _, txNo := range pb.Delivered {
		dr.delivered[tt.TxNo(txNo)] = struct{}{}
	}
	return dr
}

func DeliveredTXsFromDslStruct(ds *trantorpbtypes.DeliveredTXs, logger logging.Logger) *DeliveredTXs {
	dr := NewDeliveredTXs(logger, len(ds.Delivered))
	dr.lowWM = ds.LowWm
	for _, txNo := range ds.Delivered {
		dr.delivered[txNo] = struct{}{}
	}
	return dr
}

// Add adds a transaction number that is considered delivered to the DeliveredTXs.
// Returns true if the transaction number has been added now (after not being previously present).
// Returns false if the transaction number has already been added before the call to Add.
func (dr *DeliveredTXs) Add(txNo tt.TxNo) bool {
	if dr.CanAdd(txNo) {
		dr.delivered[txNo] = struct{}{}
		return true
	}

	return false
}

// CanAdd checks if a transaction number is considered delivered to the DeliveredTXs.
// Returns true if the transaction number has not been added previously.
// Returns false if the transaction number has already been added before the call to CanAdd.
func (dr *DeliveredTXs) CanAdd(txNo tt.TxNo) bool {
	if txNo < dr.lowWM {
		dr.logger.Log(logging.LevelDebug, "Transaction number below client's watermark window.",
			"lowWM", dr.lowWM, "txNo", txNo)
		return false
	}

	_, alreadyPresent := dr.delivered[txNo]
	return !alreadyPresent
}

// GarbageCollect reduces the memory footprint of the DeliveredTXs
// by deleting a contiguous prefix of delivered transaction numbers
// and increasing the low watermark accordingly.
// Returns the new low watermark.
func (dr *DeliveredTXs) GarbageCollect() tt.TxNo {

	for _, ok := dr.delivered[dr.lowWM]; ok; _, ok = dr.delivered[dr.lowWM] {
		delete(dr.delivered, dr.lowWM)
		dr.lowWM++
	}

	return dr.lowWM
}

func (dr *DeliveredTXs) Pb() *trantorpb.DeliveredTXs {
	delivered := make([]uint64, len(dr.delivered))
	for i, txNo := range maputil.GetSortedKeys(dr.delivered) {
		delivered[i] = txNo.Pb()
	}

	return &trantorpb.DeliveredTXs{
		LowWm:     dr.lowWM.Pb(),
		Delivered: delivered,
	}
}

func (dr *DeliveredTXs) DslStruct() *trantorpbtypes.DeliveredTXs {
	return &trantorpbtypes.DeliveredTXs{
		LowWm:     dr.lowWM,
		Delivered: maputil.GetSortedKeys(dr.delivered),
	}
}
