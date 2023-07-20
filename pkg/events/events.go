/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package events

import (
	"google.golang.org/protobuf/types/known/wrapperspb"

	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

func Redirect(event *eventpbtypes.Event, destination t.ModuleID) *eventpbtypes.Event {
	return &eventpbtypes.Event{
		Type:       event.Type,
		Next:       event.Next,
		DestModule: destination,
	}
}

// ============================================================
// Event Constructors
// ============================================================

func TestingString(dest t.ModuleID, s string) *eventpbtypes.Event {
	return &eventpbtypes.Event{
		DestModule: dest,
		Type: &eventpbtypes.Event_TestingString{
			TestingString: wrapperspb.String(s),
		},
	}
}

func TestingUint(dest t.ModuleID, u uint64) *eventpbtypes.Event {
	return &eventpbtypes.Event{
		DestModule: dest,
		Type: &eventpbtypes.Event_TestingUint{
			TestingUint: wrapperspb.UInt64(u),
		},
	}
}

// Init returns an event instructing a module to initialize.
// This event is the first to be applied to a module.
func Init(destModule t.ModuleID) *eventpbtypes.Event {
	return &eventpbtypes.Event{DestModule: destModule, Type: &eventpbtypes.Event_Init{Init: &eventpbtypes.Init{}}}
}
