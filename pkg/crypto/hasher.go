package crypto

import (
	"context"
	"fmt"
	"hash"
	"runtime"

	"go.opentelemetry.io/otel/trace"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type HashImpl interface {
	New() hash.Hash
}

type HasherModuleParams struct {
	InputBufferSize int
	NumWorkers      int
}

func DefaultHasherModuleParams() *HasherModuleParams {
	return &HasherModuleParams{
		InputBufferSize: 64,
		NumWorkers:      runtime.NumCPU() - 1,
	}
}

func NewHasher(ctx context.Context, params *HasherModuleParams, hashImpl HashImpl) modules.ActiveModule {
	return modules.NewGoRoutinePoolModule(ctx, &hasherEventProc{ctx, hashImpl}, params.InputBufferSize, params.NumWorkers)
}

type hasherEventProc struct {
	ctx      context.Context
	hashImpl HashImpl
}

func (hasher *hasherEventProc) ApplyEvent(ctx context.Context, event *eventpb.Event) *events.EventList {
	switch e := event.Type.(type) {
	case *eventpb.Event_Init:
		// no actions on init
		return events.EmptyList()
	case *eventpb.Event_HashRequest:
		// HashRequest is the only event understood by the hasher module.
		if origin, ok := e.HashRequest.Origin.Type.(*eventpb.HashOrigin_Dsl); ok {
			ctx = propagateCtxFromOrigin(ctx, origin.Dsl)
		}
		ctx, span := startSpan(ctx, "HashRequest", trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		// Create a slice for the resulting digests containing one element for each data item to be hashed.
		digests := make([][]byte, len(e.HashRequest.Data))

		// Hash each data item contained in the event
		for i, data := range e.HashRequest.Data {

			// One data item consists of potentially multiple byte slices.
			// Add each of them to the hash function.
			h := hasher.hashImpl.New()
			for _, d := range data.Data {
				h.Write(d)
			}

			// Save resulting digest in the result slice
			digests[i] = h.Sum(nil)
		}

		_, spanRes := startSpan(ctx, "HashResult", trace.WithSpanKind(trace.SpanKindProducer))
		defer spanRes.End()

		// Return all computed digests in one common event.
		return events.ListOf(
			events.HashResult(t.ModuleID(e.HashRequest.Origin.Module), digests, e.HashRequest.Origin),
		)
	default:
		// Complain about all other incoming event types.
		panic(fmt.Errorf("unexpected type of Hash event: %T", event.Type))
	}
}
