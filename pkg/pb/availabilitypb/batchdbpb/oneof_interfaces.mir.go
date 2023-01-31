package batchdbpb

import (
	bcpb "github.com/filecoin-project/mir/pkg/pb/aleapb/bcpb"
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
)

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_Lookup) Unwrap() *LookupBatch {
	return w.Lookup
}

func (w *Event_LookupResponse) Unwrap() *LookupBatchResponse {
	return w.LookupResponse
}

func (w *Event_Store) Unwrap() *StoreBatch {
	return w.Store
}

func (w *Event_Stored) Unwrap() *BatchStored {
	return w.Stored
}

type LookupBatchOrigin_Type = isLookupBatchOrigin_Type

type LookupBatchOrigin_TypeWrapper[T any] interface {
	LookupBatchOrigin_Type
	Unwrap() *T
}

func (w *LookupBatchOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *LookupBatchOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}

type StoreBatchOrigin_Type = isStoreBatchOrigin_Type

type StoreBatchOrigin_TypeWrapper[T any] interface {
	StoreBatchOrigin_Type
	Unwrap() *T
}

func (w *StoreBatchOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *StoreBatchOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}

func (w *StoreBatchOrigin_AleaBroadcast) Unwrap() *bcpb.StoreBatchOrigin {
	return w.AleaBroadcast
}
