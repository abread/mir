// Handles the processing, display and retrieval of events from a given eventlog file
package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/ttacon/chalk"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/filecoin-project/mir/pkg/eventlog"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/recordingpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// extracts events from eventlog entries and
// forwards them for display
func displayEvents(args *arguments) error {

	// new reader
	reader, err := eventlog.NewReader(args.srcFile)

	if err != nil {
		return err
	}

	// Keep track of the position of events in the recorded log.
	index := 0

	var entry *recordingpb.Entry
	for entry, err = reader.ReadEntry(); err == nil; entry, err = reader.ReadEntry() {
		metadata := eventMetadata{
			nodeID: t.NodeID(entry.NodeId),
			time:   entry.Time,
		}
		// getting events from entry
		for _, event := range entry.Events {
			metadata.index = uint64(index)

			validEvent := args.selectedEventTypes.IsEventSelected(event)
			_, validDest := args.selectedEventDests[event.DestModule]

			if validEvent && validDest && index >= args.offset && (args.limit == 0 || index < args.offset+args.limit) {
				// If event type has been selected for displaying
				displayEvent(event, metadata)
			}

			index++
		}
	}

	if errors.Is(err, io.EOF) {
		return fmt.Errorf("error reading event log: %w", err)
	}

	fmt.Println("End of trace.")

	return nil
}

// Displays one event according to its type.
func displayEvent(event *eventpb.Event, metadata eventMetadata) {
	display(eventName(event), protojson.Format(event), metadata)
}

// Creates and returns a prefix tag for event display using event metadata
func getMetaTag(eventType string, metadata eventMetadata) string {
	boldGreen := chalk.Green.NewStyle().WithBackground(chalk.ResetColor).WithTextStyle(chalk.Bold) // setting font color and style
	boldCyan := chalk.Cyan.NewStyle().WithBackground(chalk.ResetColor).WithTextStyle(chalk.Bold)
	return fmt.Sprintf("%s %s",
		boldGreen.Style(fmt.Sprintf("[ Event_%s ]", eventType)),
		boldCyan.Style(fmt.Sprintf("[ Node #%v ] [ Time _%s ] [ Index #%s ]",
			metadata.nodeID,
			strconv.FormatInt(metadata.time, 10),
			strconv.FormatUint(metadata.index, 10))),
	)
}

// displays the event
func display(eventType string, event string, metadata eventMetadata) {
	whiteText := chalk.Bold.NewStyle().WithBackground(chalk.ResetColor).WithForeground(chalk.ResetColor)
	metaTag := getMetaTag(eventType, metadata)
	fmt.Printf("%s\n%s \n", metaTag, whiteText.Style(event))
}
