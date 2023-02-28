package batchfetcher

import (
	"go.opentelemetry.io/otel/trace"

	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
)

type outputItem struct {
	event *eventpb.Event
	f     func(e *eventpb.Event)
	span  trace.Span
}

type outputQueue struct {
	items []*outputItem
}

func (oq *outputQueue) Enqueue(item *outputItem) {
	oq.items = append(oq.items, item)
}

func (oq *outputQueue) Flush(m dsl.Module) {
	m.DslHandle().PushSpan("flush outputQueue")
	defer m.DslHandle().PopSpan()

	for len(oq.items) > 0 && oq.items[0].event != nil {
		// Convenience variable.
		item := oq.items[0]

		// Execute event output hook.
		if item.f != nil {
			item.f(item.event)
		}

		item.span.AddEvent("delivering during flush")
		item.span.End()

		// Emit queued event.
		dsl.EmitEvent(m, item.event)

		// Remove item from queue.
		oq.items = oq.items[1:]
	}
}
