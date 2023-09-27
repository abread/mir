package clientprogress

import (
	"github.com/filecoin-project/mir/pkg/pb/trantorpb"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

// DeliveredTxs tracks the watermarks of delivered transactions for a single client.
type DeliveredTXs struct {

	// LowWm is the lowest watermark of the client.
	lowWm tt.TxNo

	// Set of transaction numbers indicating whether the corresponding transaction has been delivered.
	// In addition, all transaction numbers strictly smaller than LowWM are considered delivered.
	delivered map[tt.TxNo]struct{}
}

// EmptyDeliveredReqs allocates and returns a new DeliveredReqs.
func EmptyDeliveredTXs() *DeliveredTXs {
	return NewDeliveredTXs(1)
}

// NewDeliveredReqs allocates (with the given capacity for delivered requests) and returns a new DeliveredReqs.
func NewDeliveredTXs(capacity int) *DeliveredTXs {
	return &DeliveredTXs{
		lowWm:     0,
		delivered: make(map[tt.TxNo]struct{}, capacity),
	}
}

func DeliveredTXsFromPb(pb *trantorpb.DeliveredTXs) *DeliveredTXs {
	dt := NewDeliveredTXs(len(pb.Delivered))
	dt.lowWm = tt.TxNo(pb.LowWm)
	for _, txNo := range pb.Delivered {
		dt.delivered[tt.TxNo(txNo)] = struct{}{}
	}
	return dt
}

func DeliveredTXsFromDslStruct(ds *trantorpbtypes.DeliveredTXs) *DeliveredTXs {
	dt := NewDeliveredTXs(len(ds.Delivered))
	dt.lowWm = ds.LowWm
	for _, txNo := range ds.Delivered {
		dt.delivered[txNo] = struct{}{}
	}
	return dt
}

// Contains returns true if the given txNo has already been added.
func (dt *DeliveredTXs) Contains(txNo tt.TxNo) bool {

	if txNo < dt.lowWm {
		return true
	}

	_, alreadyPresent := dt.delivered[txNo]
	return alreadyPresent
}

// IsBelowWatermarkWindow returns true if the given txNo is lower than the client's watermark window.
func (dt *DeliveredTXs) IsBelowWatermarkWindow(txNo tt.TxNo) bool {
	return txNo < dt.lowWm
}

// LowWm returns the lowest transaction number that can be inside the clien't window.
func (dt *DeliveredTXs) LowWm() tt.TxNo {
	return dt.lowWm
}

// DeliveredExactlyUpTo checks if this client's transactions were delivered exactly up to the provided
// tx number, no more, no less.
func (dt *DeliveredTXs) DeliveredExactlyUpTo(txNo tt.TxNo) bool {
	dt.GarbageCollect()
	return txNo+1 == dt.lowWm && len(dt.delivered) == 0
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

func (dt *DeliveredTXs) MergeIn(other *DeliveredTXs) {
	if other.lowWm > dt.lowWm {
		dt.lowWm = other.lowWm
	}

	for txNo := range other.delivered {
		if txNo < dt.lowWm {
			continue
		}
		dt.delivered[txNo] = struct{}{}
	}

	dt.GarbageCollect()
}

// GarbageCollect reduces the memory footprint of the DeliveredTXs
// by deleting a contiguous prefix of delivered transaction numbers
// and increasing the low watermark accordingly.
// Returns the new low watermark.
func (dt *DeliveredTXs) GarbageCollect() tt.TxNo {

	for _, ok := dt.delivered[dt.lowWm]; ok; _, ok = dt.delivered[dt.lowWm] {
		delete(dt.delivered, dt.lowWm)
		dt.lowWm++
	}

	return dt.lowWm
}

// Clone returns a copy of the DeliveredTXs object.
func (dt *DeliveredTXs) Clone() *DeliveredTXs {
	cloned := &DeliveredTXs{
		lowWm:     dt.lowWm,
		delivered: make(map[tt.TxNo]struct{}, len(dt.delivered)),
	}

	for txNo := range dt.delivered {
		cloned.delivered[txNo] = struct{}{}
	}

	return cloned
}

func (dt *DeliveredTXs) Pb() *trantorpb.DeliveredTXs {
	delivered := make([]uint64, len(dt.delivered))
	for i, txNo := range maputil.GetSortedKeys(dt.delivered) {
		delivered[i] = txNo.Pb()
	}

	return &trantorpb.DeliveredTXs{
		LowWm:     dt.lowWm.Pb(),
		Delivered: delivered,
	}
}

func (dt *DeliveredTXs) DslStruct() *trantorpbtypes.DeliveredTXs {
	return &trantorpbtypes.DeliveredTXs{
		LowWm:     dt.lowWm,
		Delivered: maputil.GetSortedKeys(dt.delivered),
	}
}
