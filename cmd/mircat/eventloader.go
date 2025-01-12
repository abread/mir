package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/filecoin-project/mir/pkg/eventlog"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	"github.com/filecoin-project/mir/pkg/pb/recordingpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type eventMetadata struct {
	time   int64
	nodeID t.NodeID
	index  uint64
}

type evTypeTree struct {
	allChildrenSelected bool
	leaves              map[string]*evTypeTree
}

// Returns the list of event names and destinations present in the given eventlog file,
// along with the total number of events present in the file.
func getEventList(filenames *[]string) (*evTypeTree, map[t.ModuleID]struct{}, int, error) {
	events := &evTypeTree{
		leaves: make(map[string]*evTypeTree),
	}
	eventDests := make(map[t.ModuleID]struct{})

	totalCount := 0
	for _, filename := range *filenames {
		// Open the file
		file, err := os.Open(filename)
		if err != nil {
			kingpin.Errorf("Error opening src file", filename, ": ", err)
			fmt.Printf("\n\n!!!\nContinuing after error with new file. Event list might be incomplete!\n!!!\n\n")
			continue
		}

		defer func(file *os.File, offset int64, whence int) {
			_, _ = file.Seek(offset, whence) // resets the file offset for successive reading
			if err = file.Close(); err != nil {
				kingpin.Errorf("Error closing src file", filename, ": ", err)
			}
		}(file, 0, 0)

		reader, err := eventlog.NewReader(file)
		if err != nil {
			kingpin.Errorf("Error creating new reader of src file", filename, ": ", err)
			fmt.Printf("\n\n!!!\nContinuing after error with new file. Event list might be incomplete!\n!!!\n\n")
			continue
		}

		cnt := 0 // Counts the total number of events in the event log.
		var entry *recordingpb.Entry
		for entry, err = reader.ReadEntry(); err == nil; entry, err = reader.ReadEntry() {
			// For each entry of the event log

			for _, eventPb := range entry.Events {
				event := eventpbtypes.EventFromPb(eventPb)

				// For each Event in the entry
				cnt++

				// Add the Event type to the set of known Events.
				tree := events
				walkEventTypeName(event, func(nameComponent string) bool {
					if _, ok := tree.leaves[nameComponent]; !ok {
						tree.leaves[nameComponent] = &evTypeTree{
							leaves: make(map[string]*evTypeTree),
						}
					}
					tree = tree.leaves[nameComponent]

					return true
				})

				eventDests[event.DestModule] = struct{}{}
			}
		}

		if !errors.Is(err, io.EOF) {
			kingpin.Errorf("Error reading src file", filename, ": ", err)
			fmt.Printf("\n\n!!!\nContinuing after error. Event list might be incomplete!\n!!!\n\n")
		}
		totalCount += cnt
	}

	return events, eventDests, totalCount, nil // TODO wrap non-blocking errors and return them here
}

var EventPrefix = regexp.MustCompile("^[^_]+_")

func walkEventTypeName(event *eventpbtypes.Event, f func(nameComponent string) bool) {
	evType := reflect.ValueOf(event.Type)
	for evType.IsValid() && evType.Kind() == reflect.Pointer {
		name := evType.Elem().Type().Name()
		name = EventPrefix.ReplaceAllString(name, "")
		if !f(name) {
			break
		}

		inner := derefGetFieldByName(evType, name)

		evType = derefGetFieldByName(inner, "Type")
		if !evType.IsValid() {
			// try to parse it as a message event
			msg := derefGetFieldByName(inner, "Msg")
			evType = derefGetFieldByName(msg, "Type")
		}

		if evType.Kind() == reflect.Interface {
			evType = evType.Elem()
		}
	}
}

func autoderef(val reflect.Value) reflect.Value {
	for val.IsValid() && (val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface) {
		val = val.Elem()
	}

	return val
}

func derefGetFieldByName(val reflect.Value, fieldName string) reflect.Value {
	val = autoderef(val)

	if val.Kind() != reflect.Struct {
		var zero reflect.Value
		return zero
	}

	return val.FieldByName(fieldName)
}

func eventName(event *eventpbtypes.Event) string {
	name := make([]string, 0, 1)
	walkEventTypeName(event, func(nameComponent string) bool {
		name = append(name, nameComponent)
		return true
	})

	return strings.Join(name, ".")
}

// selected returns true if the given event has been selected by the user according to the given criteria.
func (tt *evTypeTree) IsEventSelected(event *eventpbtypes.Event) bool {
	isSelected := true

	tree := tt
	walkEventTypeName(event, func(nameComponent string) bool {
		if tree.allChildrenSelected {
			return false
		}

		if _, ok := tree.leaves[nameComponent]; !ok {
			isSelected = false
			return false
		}

		tree = tree.leaves[nameComponent]
		return true
	})

	return isSelected
}

type IterControl uint8

const (
	IterControlStop = iota
	IterControlContinue
	IterControlDontExpand
)

func (tt *evTypeTree) Walk(f func(path string, allChildrenSelected bool, hasChildren bool) IterControl) {
	tt.walk("", f)
}

func (tt *evTypeTree) walk(pathPrefix string, f func(path string, allChildrenSelected bool, hasChildren bool) IterControl) IterControl {
	switch f(pathPrefix, tt.allChildrenSelected, len(tt.leaves) != 0) {
	case IterControlStop:
		return IterControlStop
	case IterControlDontExpand:
		return IterControlContinue
	}

	if tt.leaves == nil {
		return IterControlContinue
	}

	if pathPrefix != "" {
		pathPrefix += "."
	}

	for subname, subtree := range tt.leaves {
		switch subtree.walk(pathPrefix+subname, f) {
		case IterControlStop:
			return IterControlStop
		case IterControlDontExpand:
			panic("cannot stop expansion of expanded node")
		}
	}

	return IterControlContinue
}

func (tt *evTypeTree) IsEmpty() bool {
	return !tt.allChildrenSelected && len(tt.leaves) == 0
}

// Converts a set of strings (represented as a map) to a list.
// Returns a slice containing all the keys present in the given set.
// toList is used to convert sets to a format used by the survey library.
func toList(set map[string]struct{}) []string {
	list := make([]string, 0, len(set))
	for item := range set {
		list = append(list, item)
	}
	sort.Strings(list)
	return list
}

// Converts a list of strings to a set (represented as a map).
// Returns a map of empty structs with one entry for each unique item of the given list (the item being the map key).
// toSet is used to convert lists produced by the survey library to sets for easier lookup.
func toSet(list []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, item := range list {
		set[item] = struct{}{}
	}
	return set
}
