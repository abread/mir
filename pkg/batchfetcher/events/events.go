package events

import (
	batchfetcherpbtypes "github.com/filecoin-project/mir/pkg/pb/batchfetcherpb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

// These are specific events that are not automatically generated but are necessary for the batchfetcher
// to queue them and treat them in order.

// Event creates an eventpb.Event out of an batchfetcherpb.Event.
func Event(dest t.ModuleID, ev *batchfetcherpbtypes.Event) *eventpbtypes.Event {
	return &eventpbtypes.Event{
		DestModule: dest,
		Type: &eventpbtypes.Event_BatchFetcher{
			BatchFetcher: ev,
		},
	}
}

func ClientProgress(dest t.ModuleID, progress *trantorpbtypes.ClientProgress) *eventpbtypes.Event {
	return Event(dest, &batchfetcherpbtypes.Event{
		Type: &batchfetcherpbtypes.Event_ClientProgress{
			ClientProgress: progress,
		},
	})
}
