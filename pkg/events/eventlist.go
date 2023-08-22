/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package events

import (
	"golang.org/x/exp/slices"

	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

// EventList represents a list of Events, e.g. as produced by a module.
type EventList struct {
	evs []*eventpbtypes.Event
}

// EmptyList returns an empty EventList.
func EmptyList() EventList {
	return EventList{}
}

func EmptyListWithCapacity(cap int) EventList {
	return EventList{
		evs: make([]*eventpbtypes.Event, 0, cap),
	}
}

// ListOf returns EventList containing the given elements.
func ListOf(events ...*eventpbtypes.Event) EventList {
	return EventList{
		evs: events,
	}
}

// Len returns the number of events in the EventList.
func (el EventList) Len() int {
	return len(el.evs)
}

// PushBack appends an event to the end of the list.
func (el *EventList) PushBack(event *eventpbtypes.Event) {
	// Note: a previous attempt to ammortize the cost of growing the slice by growing it by a factor of 2 was not successful.

	el.evs = append(el.evs, event)
}

// PushBackSlice appends all events in newEvents to the end of the current EventList.
func (el *EventList) PushBackSlice(events []*eventpbtypes.Event) {
	// Note: a previous attempt to ammortize the cost of growing the slice by growing it by a factor of 2 was not successful.

	el.evs = append(el.evs, events...)
}

// PushBackList appends all events in newEvents to the end of the current EventList.
func (el *EventList) PushBackList(newEvents EventList) {
	el.PushBackSlice(newEvents.evs)
}

func (el *EventList) ReserveExtraSpace(sz int) {
	el.evs = slices.Grow(el.evs, sz)
}

// Head returns the first up to n events in the list as a new list.
// The original list is not modified.
func (el EventList) Head(n int) EventList {
	if n >= len(el.evs) {
		return EventList{
			evs: el.evs,
		}
	}

	return EventList{
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

func (el EventList) Transform(evTransformer func(ev *eventpbtypes.Event) *eventpbtypes.Event) EventList {
	newList := EmptyListWithCapacity(el.Len())

	for _, ev := range el.evs {
		newEv := evTransformer(ev)
		newList.PushBack(newEv)
	}

	return newList
}

// Slice returns a slice representation of the current state of the list.
// The returned slice only contains pointers to the events in this list, no deep copying is performed.
// Any modifications performed on the events will affect the contents of both the EventList and the returned slice.
func (el EventList) Slice() []*eventpbtypes.Event {
	return el.evs
}

// StripFollowUps collects all follow-up Events of the Events in a new EventList.
// It returns a copy of the original EventList with follow-ups removed, and a new EventList containing the collected follow-up events.
func (el EventList) StripFollowUps() (EventList, EventList) {
	// Create list of follow-up Events.
	followUps := EventList{}

	// Populate list by follow-up events, removing them from the original
	newEl := el.Transform(func(ev *eventpbtypes.Event) *eventpbtypes.Event {
		if len(ev.Next) > 0 {
			followUps.PushBackSlice(ev.Next)

			// replace event with shallow copy missing follow-ups
			newEv := *ev
			ev = &newEv
			ev.Next = nil
		}

		return ev
	})

	// Return populated list of follow-up events.
	return newEl, followUps
}

// Iterator returns a pointer to an EventListIterator object used to iterate over the events in this list,
// starting from the beginning of the list.
func (el EventList) Iterator() EventListIterator {
	return EventListIterator{
		evSlice: el.evs,
	}
}

// EventListIterator is an object returned from EventList.Iterator
// used to iterate over the elements (Events) of an EventList using the iterator's Next method.
type EventListIterator struct {
	evSlice []*eventpbtypes.Event
}

// Next will return the next Event until the end of the associated EventList is encountered.
// Thereafter, it will return nil.
func (eli *EventListIterator) Next() *eventpbtypes.Event {
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
