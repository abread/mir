package batchfetcherpb

import trantorpb "github.com/filecoin-project/mir/pkg/pb/trantorpb"

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_NewOrderedBatch) Unwrap() *NewOrderedBatch {
	return w.NewOrderedBatch
}

func (w *Event_ClientProgress) Unwrap() *trantorpb.ClientProgress {
	return w.ClientProgress
}
