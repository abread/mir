package batchdbpbdsl

import (
	trace "go.opentelemetry.io/otel/trace"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/types"
	requestpb "github.com/filecoin-project/mir/pkg/pb/requestpb"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func LookupBatch[C any](m dsl.Module, destModule types.ModuleID, batchId []uint8, context *C) {
	kind := trace.WithSpanKind(trace.SpanKindProducer)
	m.DslHandle().PushSpan("LookupBatch", kind)
	defer m.DslHandle().PopSpan()

	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.LookupBatchOrigin{
		Module: m.ModuleID(),
		Type:   &types1.LookupBatchOrigin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.LookupBatch(destModule, batchId, origin))
}

func LookupBatchResponse(m dsl.Module, destModule types.ModuleID, found bool, txs []*requestpb.Request, metadata []uint8, origin *types1.LookupBatchOrigin) {
	dsl.EmitMirEvent(m, events.LookupBatchResponse(destModule, found, txs, metadata, origin))
}

func StoreBatch[C any](m dsl.Module, destModule types.ModuleID, batchId types.BatchID, txIds []types.TxID, txs []*requestpb.Request, metadata []uint8, context *C) {
	kind := trace.WithSpanKind(trace.SpanKindProducer)
	m.DslHandle().PushSpan("StoreBatch", kind)
	defer m.DslHandle().PopSpan()

	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.StoreBatchOrigin{
		Module: m.ModuleID(),
		Type:   &types1.StoreBatchOrigin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.StoreBatch(destModule, batchId, txIds, txs, metadata, origin))
}

func BatchStored(m dsl.Module, destModule types.ModuleID, origin *types1.StoreBatchOrigin) {
	dsl.EmitMirEvent(m, events.BatchStored(destModule, origin))
}
