package deploytest

import (
	"context"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/net"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	messagepbtypes "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	transportpbtypes "github.com/filecoin-project/mir/pkg/pb/transportpb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

type TransportFilter = func(msg *messagepbtypes.Message, source t.NodeID, dest t.NodeID) bool

type FilteredLink struct {
	link   net.Transport
	ownID  t.NodeID
	filter TransportFilter
}

func NewFilteredTransport(inner net.Transport, ownID t.NodeID, filter TransportFilter) net.Transport {
	if filter == nil {
		return inner
	}

	return &FilteredLink{
		link:   inner,
		ownID:  ownID,
		filter: filter,
	}
}

func getSendMessageEv(ev *eventpbtypes.Event) (*transportpbtypes.SendMessage, bool) {
	if te, ok := ev.Type.(*eventpbtypes.Event_Transport); ok {
		if sme, ok := te.Transport.Type.(*transportpbtypes.Event_SendMessage); ok {
			return sme.SendMessage, true
		}
	}

	return nil, false
}

func getForceSendMessageEv(ev *eventpbtypes.Event) (*transportpbtypes.ForceSendMessage, bool) {
	if te, ok := ev.Type.(*eventpbtypes.Event_Transport); ok {
		if sme, ok := te.Transport.Type.(*transportpbtypes.Event_ForceSendMessage); ok {
			return sme.ForceSendMessage, true
		}
	}

	return nil, false
}

func (fl *FilteredLink) applyFilterToEventSlice(events []*eventpbtypes.Event) []*eventpbtypes.Event {
	filtered := make([]*eventpbtypes.Event, 0, len(events))

	for _, ev := range events {
		ev = fl.filterEvent(ev)
		if ev != nil {
			filtered = append(filtered, ev)
		}
	}

	return filtered
}

func (fl *FilteredLink) filterEvent(event *eventpbtypes.Event) *eventpbtypes.Event {
	if ev, ok := getSendMessageEv(event); ok {
		// don't mutate the original event to avoid races
		newEv := *ev
		ev = &newEv

		ev.Destinations = sliceutil.Filter(ev.Destinations, func(_ int, dest t.NodeID) bool {
			return fl.filter(ev.Msg, fl.ownID, dest)
		})

		if len(ev.Destinations) > 0 {
			return &eventpbtypes.Event{
				DestModule: event.DestModule,
				Type: &eventpbtypes.Event_Transport{
					Transport: &transportpbtypes.Event{
						Type: &transportpbtypes.Event_SendMessage{
							SendMessage: ev,
						},
					},
				},
				Next: fl.applyFilterToEventSlice(event.Next),
			}
		}

		return nil
	} else if ev, ok := getForceSendMessageEv(event); ok {
		// don't mutate the original event to avoid races
		newEv := *ev
		ev = &newEv

		ev.Destinations = sliceutil.Filter(ev.Destinations, func(_ int, dest t.NodeID) bool {
			return fl.filter(ev.Msg, fl.ownID, dest)
		})

		if len(ev.Destinations) > 0 {
			return &eventpbtypes.Event{
				DestModule: event.DestModule,
				Type: &eventpbtypes.Event_Transport{
					Transport: &transportpbtypes.Event{
						Type: &transportpbtypes.Event_ForceSendMessage{
							ForceSendMessage: ev,
						},
					},
				},
				Next: fl.applyFilterToEventSlice(event.Next),
			}
		}

		return nil
	}

	return event
}

func (fl *FilteredLink) ApplyEvents(
	ctx context.Context,
	eventList events.EventList,
) error {
	filtered := events.ListOf(fl.applyFilterToEventSlice(eventList.Slice())...)
	return fl.link.ApplyEvents(ctx, filtered)
}

// The ImplementsModule method only serves the purpose of indicating that this is a Module and must not be called.
func (fl *FilteredLink) ImplementsModule() {}

func (fl *FilteredLink) Send(dest t.NodeID, msg *messagepbtypes.Message) error {
	if fl.filter(msg, fl.ownID, dest) {
		return fl.link.Send(dest, msg)
	}
	return nil
}

func (fl *FilteredLink) EventsOut() <-chan events.EventList {
	return fl.link.EventsOut()
}

func (fl *FilteredLink) CloseOldConnections(m *trantorpbtypes.Membership) {
	fl.link.CloseOldConnections(m)
}

func (fl *FilteredLink) Start() error {
	return fl.link.Start()
}

func (fl *FilteredLink) Connect(m *trantorpbtypes.Membership) {
	fl.link.Connect(m)
}

func (fl *FilteredLink) WaitFor(x int) error {
	return fl.link.WaitFor(x)
}

func (fl *FilteredLink) Stop() {
	fl.link.Stop()
}
