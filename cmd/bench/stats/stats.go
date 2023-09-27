package stats

import (
	"encoding/csv"
	"time"

	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
)

type Tracker interface {
	Submit(tx *trantorpbtypes.Transaction)
	Deliver(tx *trantorpbtypes.Transaction)
	AssumeDelivered(tx *trantorpbtypes.Transaction)
}

type Stats interface {
	Fill()
	WriteCSVHeader(w *csv.Writer) error
	WriteCSVRecord(w *csv.Writer, d time.Duration) error
}
