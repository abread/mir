package threshcryptopbdsl

import (
	trace "go.opentelemetry.io/otel/trace"

	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/events"
	types1 "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/types"
	tctypes "github.com/filecoin-project/mir/pkg/threshcrypto/tctypes"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func SignShare[C any](m dsl.Module, destModule types.ModuleID, data [][]uint8, context *C) {
	kind := trace.WithSpanKind(trace.SpanKindProducer)
	m.DslHandle().PushSpan("SignShare", kind)
	defer m.DslHandle().PopSpan()

	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.SignShareOrigin{
		Module: m.ModuleID(),
		Type:   &types1.SignShareOrigin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.SignShare(destModule, data, origin))
}

func SignShareResult(m dsl.Module, destModule types.ModuleID, signatureShare tctypes.SigShare, origin *types1.SignShareOrigin) {
	dsl.EmitMirEvent(m, events.SignShareResult(destModule, signatureShare, origin))
}

func VerifyShare[C any](m dsl.Module, destModule types.ModuleID, data [][]uint8, signatureShare tctypes.SigShare, nodeId types.NodeID, context *C) {
	kind := trace.WithSpanKind(trace.SpanKindProducer)
	m.DslHandle().PushSpan("VerifyShare", kind)
	defer m.DslHandle().PopSpan()

	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.VerifyShareOrigin{
		Module: m.ModuleID(),
		Type:   &types1.VerifyShareOrigin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.VerifyShare(destModule, data, signatureShare, nodeId, origin))
}

func VerifyShareResult(m dsl.Module, destModule types.ModuleID, ok bool, error string, origin *types1.VerifyShareOrigin) {
	dsl.EmitMirEvent(m, events.VerifyShareResult(destModule, ok, error, origin))
}

func VerifyFull[C any](m dsl.Module, destModule types.ModuleID, data [][]uint8, fullSignature tctypes.FullSig, context *C) {
	kind := trace.WithSpanKind(trace.SpanKindProducer)
	m.DslHandle().PushSpan("VerifyFull", kind)
	defer m.DslHandle().PopSpan()

	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.VerifyFullOrigin{
		Module: m.ModuleID(),
		Type:   &types1.VerifyFullOrigin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.VerifyFull(destModule, data, fullSignature, origin))
}

func VerifyFullResult(m dsl.Module, destModule types.ModuleID, ok bool, error string, origin *types1.VerifyFullOrigin) {
	dsl.EmitMirEvent(m, events.VerifyFullResult(destModule, ok, error, origin))
}

func Recover[C any](m dsl.Module, destModule types.ModuleID, data [][]uint8, signatureShares []tctypes.SigShare, context *C) {
	kind := trace.WithSpanKind(trace.SpanKindProducer)
	m.DslHandle().PushSpan("Recover", kind)
	defer m.DslHandle().PopSpan()

	contextID := m.DslHandle().StoreContext(context)
	traceCtx := m.DslHandle().TraceContextAsMap()

	origin := &types1.RecoverOrigin{
		Module: m.ModuleID(),
		Type:   &types1.RecoverOrigin_Dsl{Dsl: dsl.MirOrigin(contextID, traceCtx)},
	}

	dsl.EmitMirEvent(m, events.Recover(destModule, data, signatureShares, origin))
}

func RecoverResult(m dsl.Module, destModule types.ModuleID, fullSignature tctypes.FullSig, ok bool, error string, origin *types1.RecoverOrigin) {
	dsl.EmitMirEvent(m, events.RecoverResult(destModule, fullSignature, ok, error, origin))
}
