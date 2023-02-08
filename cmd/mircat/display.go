// Handles the processing, display and retrieval of events from a given eventlog file
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

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

	// Keep track of the position of events in the recorded log.
	index := 0

	filenames := strings.Split(*args.src, " ")

	// Create a reader for the input event log file.
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			kingpin.Errorf("Error opening src file", filename, ": ", err)
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
			kingpin.Errorf("Error opening new reader for src file", filename, ": ", err)
			fmt.Printf("\n\n!!!\nContinuing after error with next file\n!!!\n\n")
			continue
		}

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
				event = customTransform(event) // Apply custom transformation to event.

				if validEvent && validDest && event != nil && index >= args.offset && (args.limit == 0 || index < args.offset+args.limit) {
					// If event type has been selected for displaying
					displayEvent(event, metadata)
				}

				index++
			}
		}

		if !errors.Is(err, io.EOF) {
			kingpin.Errorf("Error reading event log for src file", filename, ": ", err)
			fmt.Printf("\n\n!!!\nContinuing after error with next file\n!!!\n\n")
			continue
		}
	}
	fmt.Println("End of trace.")

	return nil //TODO Wrap non-blocking errors and return here
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
