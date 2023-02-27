package batchdbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types2 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_BatchDb](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponLookupBatch(m dsl.Module, handler func(batchId []uint8, origin *types.LookupBatchOrigin) error) {
	UponEvent[*types.Event_Lookup](m, func(ev *types.LookupBatch) error {
		originWrapper, ok := ev.Origin.Type.(*types.LookupBatchOrigin_Dsl)
		if ok {
			m.DslHandle().ImportTraceContextFromMap(originWrapper.Dsl.TraceContext)
		}

		m.DslHandle().PushSpan("LookupBatch")
		defer m.DslHandle().PopSpan()

		return handler(ev.BatchId, ev.Origin)
	})
}

func UponLookupBatchResponse[C any](m dsl.Module, handler func(found bool, txs []*requestpb.Request, metadata []uint8, context *C) error) {
	UponEvent[*types.Event_LookupResponse](m, func(ev *types.LookupBatchResponse) error {
		originWrapper, ok := ev.Origin.Type.(*types.LookupBatchOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		m.DslHandle().ImportTraceContextFromMap(originWrapper.Dsl.TraceContext)

		return handler(ev.Found, ev.Txs, ev.Metadata, context)
	})
}

func UponStoreBatch(m dsl.Module, handler func(batchId types2.BatchID, txIds []types2.TxID, txs []*requestpb.Request, metadata []uint8, origin *types.StoreBatchOrigin) error) {
	UponEvent[*types.Event_Store](m, func(ev *types.StoreBatch) error {
		originWrapper, ok := ev.Origin.Type.(*types.StoreBatchOrigin_Dsl)
		if ok {
			m.DslHandle().ImportTraceContextFromMap(originWrapper.Dsl.TraceContext)
		}

		m.DslHandle().PushSpan("StoreBatch")
		defer m.DslHandle().PopSpan()

		return handler(ev.BatchId, ev.TxIds, ev.Txs, ev.Metadata, ev.Origin)
	})
}

func UponBatchStored[C any](m dsl.Module, handler func(context *C) error) {
	UponEvent[*types.Event_Stored](m, func(ev *types.BatchStored) error {
		originWrapper, ok := ev.Origin.Type.(*types.StoreBatchOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		m.DslHandle().ImportTraceContextFromMap(originWrapper.Dsl.TraceContext)

		return handler(context)
	})
}
