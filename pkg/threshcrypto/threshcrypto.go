// Package threshcrypto provides an implementation of the MirModule module.
// It supports TBLS signatures.
package threshcrypto

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/dslpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/threshcryptopb"
	tcEvents "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/events"
	threshcryptopbtypes "github.com/filecoin-project/mir/pkg/pb/threshcryptopb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

var name = "github.com/filecoin-project/mir/pkg/threshcrypto"

type MirModule struct {
	ctx          context.Context
	threshCrypto ThreshCrypto
}

func New(ctx context.Context, threshCrypto ThreshCrypto) *MirModule {
	return &MirModule{ctx: ctx, threshCrypto: threshCrypto}
}

func (c *MirModule) ApplyEvents(eventsIn *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsConcurrently(eventsIn, c.ApplyEvent)
}

func (c *MirModule) ApplyEvent(event *eventpb.Event) (*events.EventList, error) {
	switch e := event.Type.(type) {
	case *eventpb.Event_Init:
		// no actions on init
		return events.EmptyList(), nil
	case *eventpb.Event_ThreshCrypto:
		return c.applyTCEvent(e.ThreshCrypto)
	default:
		// Complain about all other incoming event types.
		return nil, fmt.Errorf("unexpected type of event in threshcrypto MirModule: %T", event.Type)
	}
}

func startSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	tracer := otel.Tracer(name)
	ctx, span := tracer.Start(ctx, name, opts...)
	return ctx, span
}

func propagateCtxFromOrigin(ctx context.Context, origin *dslpb.Origin) context.Context {
	return propagation.TraceContext{}.Extract(ctx, propagation.MapCarrier(origin.TraceContext))
}

// apply a thresholdcryptopb.Event
func (c *MirModule) applyTCEvent(event *threshcryptopb.Event) (*events.EventList, error) {
	ctx := c.ctx

	switch e := event.Type.(type) {
	case *threshcryptopb.Event_SignShare:
		// Compute signature share
		if origin, ok := e.SignShare.Origin.Type.(*threshcryptopb.SignShareOrigin_Dsl); ok {
			ctx = propagateCtxFromOrigin(ctx, origin.Dsl)
		}
		ctx, span := startSpan(ctx, "SignShare", trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		sigShare, err := c.threshCrypto.SignShare(e.SignShare.Data)
		if err != nil {
			return nil, err
		}

		_, spanResp := startSpan(ctx, "SignShareResult", trace.WithSpanKind(trace.SpanKindProducer))
		defer spanResp.End()

		origin := threshcryptopbtypes.SignShareOriginFromPb(e.SignShare.Origin)
		return events.ListOf(
			tcEvents.SignShareResult(origin.Module, sigShare, origin).Pb(),
		), nil

	case *threshcryptopb.Event_VerifyShare:
		// Verify signature share
		if origin, ok := e.VerifyShare.Origin.Type.(*threshcryptopb.VerifyShareOrigin_Dsl); ok {
			ctx = propagateCtxFromOrigin(ctx, origin.Dsl)
		}
		ctx, span := startSpan(ctx, "VerifyShare", trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		err := c.threshCrypto.VerifyShare(e.VerifyShare.Data, e.VerifyShare.SignatureShare, t.NodeID(e.VerifyShare.NodeId))

		_, spanResp := startSpan(ctx, "VerifyShareResult", trace.WithSpanKind(trace.SpanKindProducer))
		defer spanResp.End()

		ok := err == nil
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		origin := threshcryptopbtypes.VerifyShareOriginFromPb(e.VerifyShare.Origin)
		return events.ListOf(
			tcEvents.VerifyShareResult(origin.Module, ok, errStr, origin).Pb(),
		), nil

	case *threshcryptopb.Event_VerifyFull:
		// Verify full signature
		if origin, ok := e.VerifyFull.Origin.Type.(*threshcryptopb.VerifyFullOrigin_Dsl); ok {
			ctx = propagateCtxFromOrigin(ctx, origin.Dsl)
		}
		ctx, span := startSpan(ctx, "VerifyFull", trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		err := c.threshCrypto.VerifyFull(e.VerifyFull.Data, e.VerifyFull.FullSignature)

		_, spanResp := startSpan(ctx, "VerifyFullResult", trace.WithSpanKind(trace.SpanKindProducer))
		defer spanResp.End()

		ok := err == nil
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		origin := threshcryptopbtypes.VerifyFullOriginFromPb(e.VerifyFull.Origin)
		return events.ListOf(
			tcEvents.VerifyFullResult(origin.Module, ok, errStr, origin).Pb(),
		), nil

	case *threshcryptopb.Event_Recover:
		// Recover full signature from shares
		if origin, ok := e.Recover.Origin.Type.(*threshcryptopb.RecoverOrigin_Dsl); ok {
			ctx = propagateCtxFromOrigin(ctx, origin.Dsl)
		}
		ctx, span := startSpan(ctx, "Recover", trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		fullSig, err := c.threshCrypto.Recover(e.Recover.Data, e.Recover.SignatureShares)

		_, spanResp := startSpan(ctx, "RecoverResult", trace.WithSpanKind(trace.SpanKindProducer))
		defer spanResp.End()

		ok := err == nil
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		origin := threshcryptopbtypes.RecoverOriginFromPb(e.Recover.Origin)
		return events.ListOf(
			tcEvents.RecoverResult(origin.Module, fullSig, ok, errStr, origin).Pb(),
		), nil
	default:
		// Complain about all other incoming event types.
		return nil, fmt.Errorf("unexpected type of threshcrypto event in threshcrypto MirModule: %T", event.Type)
	}
}

// The ImplementsModule method only serves the purpose of indicating that this is a Module and must not be called.
func (c *MirModule) ImplementsModule() {}
