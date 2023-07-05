package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	bcbpbtypes "github.com/filecoin-project/mir/pkg/pb/bcbpb/types"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

type controlModule struct {
	eventsOut chan *events.EventList
	isLeader  bool
}

func newControlModule(isLeader bool) modules.ActiveModule {
	return &controlModule{
		eventsOut: make(chan *events.EventList),
		isLeader:  isLeader,
	}
}

func (m *controlModule) ImplementsModule() {}

func (m *controlModule) ApplyEvents(_ context.Context, events *events.EventList) error {
	iter := events.Iterator()
	for event := iter.Next(); event != nil; event = iter.Next() {
		switch event.Type.(type) {

		case *eventpbtypes.Event_Init:
			if m.isLeader {
				go func() {
					err := m.readMessageFromConsole()
					if err != nil {
						panic(err)
					}
				}()
			} else {
				fmt.Println("Waiting for the message...")
			}

		case *eventpbtypes.Event_Bcb:
			bcbEvent := event.Type.(*eventpbtypes.Event_Bcb).Bcb
			switch bcbEvent.Type.(type) {

			case *bcbpbtypes.Event_Deliver:
				deliverEvent := bcbEvent.Type.(*bcbpbtypes.Event_Deliver).Deliver
				fmt.Println("Leader says: ", string(deliverEvent.Data))

			default:
				return es.Errorf("unknown bcb event type: %T", bcbEvent.Type)
			}

		default:
			return es.Errorf("unknown event type: %T", event.Type)
		}
	}

	return nil
}

func (m *controlModule) EventsOut() <-chan *events.EventList {
	return m.eventsOut
}

func (m *controlModule) readMessageFromConsole() error {
	// Read the user input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Type in a message and press Enter: ")
	scanner.Scan()
	if scanner.Err() != nil {
		return es.Errorf("error reading from console: %w", scanner.Err())
	}

	m.eventsOut <- events.ListOf(&eventpbtypes.Event{
		DestModule: "bcb",
		Type: &eventpbtypes.Event_Bcb{
			Bcb: &bcbpbtypes.Event{
				Type: &bcbpbtypes.Event_Request{
					Request: &bcbpbtypes.BroadcastRequest{
						Data: []byte(scanner.Text()),
					},
				},
			},
		},
	})

	return nil
}
