package modules

import (
	"context"
	"runtime/debug"
	"sync"
	"sync/atomic"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/events"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

// ApplyEventsSequentially takes a list of events and applies the given applyEvent function to each event in the list.
// Processing is performed sequentially, one event at a time, in the order of the input list.
// The EventLists returned by applyEvent are aggregated in a single EventList (in order of creation)
// and returned by ApplyEventsSequentially.
func ApplyEventsSequentially(
	eventsIn events.EventList,
	applyEvent func(*eventpbtypes.Event) (events.EventList, error),
) (events.EventList, error) {

	eventsOut := events.EmptyList()

	iter := eventsIn.Iterator()
	for event := iter.Next(); event != nil; event = iter.Next() {
		evts, err := applyEvent(event)
		if err != nil {
			return events.EmptyList(), err
		}
		eventsOut.PushBackList(evts)
	}

	return eventsOut, nil
}

// ApplyEventsConcurrently takes a list of events and applies the given applyEvent function to each event in the list.
// Processing is performed concurrently on all events in the input list.
// Thus, the provided applyEvent function must be thread-safe.
// The EventLists returned by applyEvent are aggregated in a single EventList and returned by ApplyEventsConcurrently.
// Despite being executed concurrently,
// the order of the returned results preserves the order of the corresponding input events.
// Thus, if applyEvent is deterministic, the output of ApplyEventsConcurrently is also deterministic.
// If one or more errors occur during processing, ApplyEventsConcurrently returns the first of them,
// along with an empty EventList.
func ApplyEventsConcurrently(
	eventsIn events.EventList,
	applyEvent func(*eventpbtypes.Event) (events.EventList, error),
) (events.EventList, error) {

	// Initialize channels into which the results of each invocation of applyEvent will be written.
	results := make([]chan events.EventList, eventsIn.Len())
	errors := make([]chan error, eventsIn.Len())
	for i := 0; i < eventsIn.Len(); i++ {
		results[i] = make(chan events.EventList)
		errors[i] = make(chan error)
	}

	// Start one concurrent worker for each event in the input list.
	// Note that the processing starts concurrently for all events, and only the writing of the results is synchronized.
	iter := eventsIn.Iterator()
	i := 0
	for event := iter.Next(); event != nil; event = iter.Next() {

		go func(e *eventpbtypes.Event, j int) {

			// Apply the input event, catching potential panics.
			res, err := applySafely(e, applyEvent)

			// Write processing results to the output channels.
			// Attention: Those (unbuffered) channels must be read by the aggregator in the same order
			//            as they are being written here, otherwise the system gets stuck.
			results[j] <- res
			errors[j] <- err

		}(event, i)
		i++

	}

	// The event processing results will be aggregated here.
	var firstError error
	eventsOut := events.EmptyList()

	// For each input event, read the processing result from the common channels and aggregate it with the rest.
	for i := 0; i < eventsIn.Len(); i++ {

		// Attention: The (unbuffered) errors and results channels must be read in the same order
		//            as they are being written by the worker goroutines, otherwise the system gets stuck.

		// Read results from common channel and add it to the accumulator.
		evList := <-results[i]
		eventsOut.PushBackList(evList)

		// Read error from common channel.
		// We only consider the first error, as ApplyEventsConcurrently only returns a single error.
		// TODO: Explore possibilities of aggregating multiple errors in one.
		if err := <-errors[i]; err != nil && firstError == nil {
			firstError = err
		}

	}

	// Return the resulting events or an error.
	if firstError != nil {
		return events.EmptyList(), firstError
	}
	return eventsOut, nil
}

type EventProcessor interface {
	// TODO: add if needed
	// CheckEvent(ev *eventpbtypes.Event) error

	ApplyEvent(ev *eventpbtypes.Event) events.EventList
}

type SimpleEventApplier struct {
	EventProcessor
}

func (sea SimpleEventApplier) ImplementsModule() {}
func (sea SimpleEventApplier) ApplyEvents(evs events.EventList) (events.EventList, error) {
	return ApplyEventsSequentially(evs, func(e *eventpbtypes.Event) (events.EventList, error) {
		// if err := sea.CheckEvent(e); err != nil { return err }
		return sea.ApplyEvent(e), nil
	})
}

func (sea SimpleEventApplier) IntoGoroutinePool(ctx context.Context, workers int) ActiveModule {
	return NewGoRoutinePoolModule(ctx, sea.EventProcessor, workers)
}

type poolModule struct {
	processor  EventProcessor
	inputChan  chan *eventpbtypes.Event
	outputChan chan events.EventList

	// ensures events are scheduled for processing in order of arrival
	// ApplyEvents will spawn a goroutine with some seqno. This goroutine will wait for nextInputSeqNo
	// to be its seqno, then push all its input to inputChan, and aftewards increment nextInputSeqNo
	// and signal other routines to try to make progress.
	// Sequence numbers are assigned from lastFutureInputSeqNo.
	nextInputReady       *sync.Cond // protects nextInputSeqNo
	nextInputSeqNo       uint64
	lastFutureInputSeqNo uint64 // must be manipulated with atomics
}

func NewGoRoutinePoolModule(ctx context.Context, processor EventProcessor, workers int) ActiveModule {
	mod := &poolModule{
		processor:  processor,
		inputChan:  make(chan *eventpbtypes.Event, workers),
		outputChan: make(chan events.EventList, workers*2),

		nextInputReady: sync.NewCond(&sync.Mutex{}),
		nextInputSeqNo: 1, // 0 will never be used
	}
	doneChan := ctx.Done()

	for i := 0; i < workers; i++ {
		go func() {
		Loop:
			for {
				select {
				case <-doneChan:
					break Loop
				case ev := <-mod.inputChan:
					select {
					case mod.outputChan <- processor.ApplyEvent(ev):
						continue
					case <-doneChan:
						break Loop
					}
				}
			}
		}()
	}

	return mod
}

func (m *poolModule) ImplementsModule() {}
func (m *poolModule) EventsOut() <-chan events.EventList {
	return m.outputChan
}
func (m *poolModule) ApplyEvents(ctx context.Context, events events.EventList) error {
	seqno := atomic.AddUint64(&m.lastFutureInputSeqNo, 1)

	select {
	case _, channelOpen := <-ctx.Done():
		if !channelOpen {
			// we're done
			return nil
		}
	default:
	}

	go func() {
		m.nextInputReady.L.Lock()
		defer m.nextInputReady.L.Unlock()

		for m.nextInputSeqNo != seqno {
			m.nextInputReady.Wait()
		}

		it := events.Iterator()
	InputLoop:
		for ev := it.Next(); ev != nil; ev = it.Next() {
			select {
			case <-ctx.Done():
				break InputLoop
			case m.inputChan <- ev:
			}
		}

		m.nextInputSeqNo++
		m.nextInputReady.Broadcast()
	}()
	return nil
}

// applySafely is a wrapper around an event processing function that catches its panic and returns it as an error.
func applySafely(
	event *eventpbtypes.Event,
	processingFunc func(*eventpbtypes.Event) (events.EventList, error),
) (result events.EventList, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rErr, ok := r.(error); ok {
				err = es.Errorf("event application panicked: %w\nStack trace:\n%s", rErr, string(debug.Stack()))
			} else {
				err = es.Errorf("event application panicked: %v\nStack trace:\n%s", r, string(debug.Stack()))
			}
		}
	}()

	return processingFunc(event)
}
