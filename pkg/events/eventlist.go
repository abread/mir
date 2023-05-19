/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package events

import (
	"golang.org/x/exp/slices"

	"github.com/filecoin-project/mir/pkg/pb/eventpb"
)

// EventList represents a list of Events, e.g. as produced by a module.
type EventList struct {
	evs []*eventpb.Event
}

// EmptyList returns an empty EventList.
// TODO: consider passing EventList by value here and everywhere else.
func EmptyList() *EventList {
	return &EventList{}
}

func EmptyListWithCapacity(cap int) *EventList {
	return &EventList{
		evs: make([]*eventpb.Event, 0, cap),
	}
}

// ListOf returns EventList containing the given elements.
func ListOf(events ...*eventpb.Event) *EventList {
	return &EventList{
		evs: events,
	}
}

// Len returns the number of events in the EventList.
func (el *EventList) Len() int {
	return len(el.evs)
}

// PushBack appends an event to the end of the list.
// Returns the EventList itself, for the convenience of chaining multiple calls to PushBack.
func (el *EventList) PushBack(event ...*eventpb.Event) *EventList {
	if len(el.evs) == cap(el.evs) {
		el.grow(len(event))
	}

	el.evs = append(el.evs, event...)

	return el
}

// PushBackSlice appends all events in newEvents to the end of the current EventList.
func (el *EventList) PushBackSlice(events []*eventpb.Event) *EventList {
	return el.PushBack(events...)
}

// PushBackList appends all events in newEvents to the end of the current EventList.
func (el *EventList) PushBackList(newEvents *EventList) *EventList {
	return el.PushBack(newEvents.evs...)
}

func (el *EventList) grow(szHint int) {
	el.evs = slices.Grow(el.evs, szHint)
}

// Head returns the first up to n events in the list as a new list.
// The original list is not modified.
func (el *EventList) Head(n int) *EventList {
	if n >= len(el.evs) {
		return &EventList{
			evs: el.evs,
		}
	}

	return &EventList{
		evs: el.evs[:n],
	}
}

// RemoveFront removes the first up to n events from the list.
// Returns the number of events actually removed.
func (el *EventList) RemoveFront(n int) int {
	if n > len(el.evs) {
		n = len(el.evs)
		el.evs = nil
		return n
	}

	// TODO: consider reallocating slice to really free up memory
	el.evs = el.evs[n:]
	return n
}

// Slice returns a slice representation of the current state of the list.
// The returned slice only contains pointers to the events in this list, no deep copying is performed.
// Any modifications performed on the events will affect the contents of both the EventList and the returned slice.
func (el *EventList) Slice() []*eventpb.Event {
	return el.evs
}

// StripFollowUps collects all follow-up Events of the Events in the list.
// It returns two lists:
// 1. An EventList containing the same events as this list, but with all follow-up events removed.
// 2. An EventList containing only those follow-up events.
func (el *EventList) StripFollowUps() (*EventList, *EventList) {
	// Create list of follow-up Events.
	followUps := EventList{}

	// Create a new EventList for events with follow-ups removed.
	plainEvents := EventList{}

	// Populate list by follow-up events
	for _, event := range el.evs {
		plainEvent, strippedEvents := Strip(event)
		plainEvents.PushBack(plainEvent)
		followUps.PushBackList(strippedEvents)
	}

	// Return populated list of follow-up events.
	return &plainEvents, &followUps
}

// Iterator returns a pointer to an EventListIterator object used to iterate over the events in this list,
// starting from the beginning of the list.
func (el *EventList) Iterator() *EventListIterator {
	return &EventListIterator{
		evSlice: el.evs,
	}
}

// EventListIterator is an object returned from EventList.Iterator
// used to iterate over the elements (Events) of an EventList using the iterator's Next method.
type EventListIterator struct {
	evSlice []*eventpb.Event
}

// Next will return the next Event until the end of the associated EventList is encountered.
// Thereafter, it will return nil.
func (eli *EventListIterator) Next() *eventpb.Event {
	// Return nil if list has been exhausted.
	if len(eli.evSlice) == 0 {
		return nil
	}

	// Obtain current element and move on to the next one.
	result := eli.evSlice[0]
	eli.evSlice = eli.evSlice[1:]

	// Return current element.
	return result
}
