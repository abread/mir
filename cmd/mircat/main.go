package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"

	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
)

// mircat is a tool for reviewing Mir state machine recordings.
// It understands the format encoded via github.com/filecoin-project/mir/eventlog
// and is able to parse and filter these log files based on the events.

// arguments represents the parameters passed to mircat.
type arguments struct {

	// File containing the event log to read.
	src *string

	// Number of events at the start of the event log to be skipped (including the ones not selected).
	// If debugging, the skipped events will still be passed to the node before the interactive debugging starts.
	offset int

	// The number of events to display.
	// When debugging, this is the maximum number of events to pause before submitting to the node.
	// All remaining events will simply be flushed without user prompt.
	limit int

	// Events selected by the user for displaying.
	selectedEventTypes *evTypeTree

	// Events with specific destination modules selected by the user for displaying
	selectedEventDests map[string]struct{}

	// If set to true, start a Node in debug mode with the given event log.
	debug bool

	// When this is true, the selected events are all passed to the custom module.
	// Note that this is completely separate from normal debugging.
	dbgModule bool

	// The rest of the fields are only used in debug mode and are otherwise ignored.

	// The ID of the node being debugged.
	// It must correspond to the ID of the node that produced the event log being read.
	ownID t.NodeID

	// IDs, in order, of all nodes in the deployment that produced the event log.
	membership []t.NodeID

	// If set to true, the events produced by the node are printed to standard output.
	// Otherwise, they are simply dropped.
	showNodeEvents bool
}

func main() {

	// Parse command-line arguments
	kingpin.Version("0.0.1")
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		kingpin.Fatalf("Cannot parse given argument", err)
	}

	filenames := strings.Split(*args.src, " ")
	fmt.Println("filenames:", filenames)
	// Scan the event log and collect all occurring event types.
	fmt.Println("Scanning input file.")
	allEvents, allDests, totalEvents, err := getEventList(&filenames)
	if err != nil {
		kingpin.Errorf("Error parsing src file", err)
		fmt.Printf("\n\n!!!\nContinuing after error. Event list might be incomplete!\n!!!\n\n")
	}
	fmt.Printf("Total number of events found: %d\n", totalEvents)

	// If no event types have been selected through command-line arguments,
	// have the user interactively select the events to include in the output.
	if args.selectedEventTypes.IsEmpty() {
		// Select events
		selectedNames := checkboxes("Please select the event types", eventSelectionOptions(allEvents))
		args.selectedEventTypes = toEvTypeTree(toList(selectedNames))
	}

	// If no event destinations have been selected through command-line arguments,
	// have the user interactively select the event destinations' to include in the output.
	if len(args.selectedEventDests) == 0 {

		allDestsStr := maputil.Transform(allDests, func(k t.ModuleID, v struct{}) (string, struct{}) { return string(k), v })
		// Select top-level events
		args.selectedEventDests = checkboxes("Please select the event destinations", allDestsStr)

	}

	fmt.Println("Command-line arguments for selecting the chosen filters:\n" +
		selectionArgs(args.selectedEventTypes, args.selectedEventDests) + "\n")

	// Display selected events or enter debug mode.
	if args.debug {
		err = debug(args)
		if err != nil {
			kingpin.Errorf("Error debugging node", err)
		}
	} else {
		err = displayEvents(args)
		if err != nil {
			kingpin.Errorf("Error Processing Events", err)
		}
	}

	fmt.Println("Command-line arguments for selecting the chosen filters:\n" +
		selectionArgs(args.selectedEventTypes, args.selectedEventDests) + "\n")
}

// parse the command line arguments
func parseArgs(args []string) (*arguments, error) {
	if len(args) == 0 {
		return nil, errors.Errorf("required input \" --src <Src_File> \" not found")
	}

	app := kingpin.New("mircat", "Utility for processing Mir state event logs.")
	src := app.Flag("src", "The input file to read.").Required().String()
	events := app.Flag("event", "Event types to be displayed.").Short('e').Strings()
	eventDests := app.Flag("event-dest", "Event destination types to be displayed.").Short('r').Strings()
	offset := app.Flag("offset", "The first offset events will not be displayed.").Default("0").Int()
	limit := app.Flag("limit", "Maximum number of events to consider for display or debug").Default("0").Int()
	dbg := app.Flag("debug", "Start a Node in debug mode with the given event log.").Short('d').Bool()
	id := app.Flag("own-id", "ID of the node to use for debugging.").String()
	dbgModule := app.Flag("module", "Debug the custom module.").Bool()
	membership := app.Flag(
		"node-id",
		"ID of one membership node, specified once for each node (debugging only).",
	).Short('m').Strings()
	showNodeEvents := app.Flag("show-node-output", "Show events generated by the node (debugging only)").Bool()

	_, err := app.Parse(args)

	if err != nil {
		return nil, err
	}

	return &arguments{
		src:                src,
		debug:              *dbg,
		dbgModule:          *dbgModule,
		ownID:              t.NodeID(*id),
		membership:         t.NodeIDSlice(*membership),
		showNodeEvents:     *showNodeEvents,
		offset:             *offset,
		limit:              *limit,
		selectedEventTypes: toEvTypeTree(*events),
		selectedEventDests: toSet(*eventDests),
	}, nil
}

// Prompts users with a list of available Events to select from.
// Returns a set of selected Events.
func checkboxes(label string, opts map[string]struct{}) map[string]struct{} {

	// Use survey library to get a list selected event names.
	selected := make([]string, 0)
	prompt := &survey.MultiSelect{
		Message: label,
		Options: toList(opts),
	}
	if err := survey.AskOne(prompt, &selected); err != nil {
		fmt.Printf("Error selecting event types: %v", err)
	}

	return toSet(selected)
}

func eventSelectionOptions(allEvents *evTypeTree) map[string]struct{} {
	options := make(map[string]struct{})
	allEvents.Walk(func(path string, _ bool, hasChildren bool) IterControl {
		if path == "" {
			return IterControlContinue
		}

		if hasChildren {
			options[path+".*"] = struct{}{}
		} else {
			options[path] = struct{}{}
		}

		return IterControlContinue
	})

	return options
}

func selectionArgs(events *evTypeTree, dests map[string]struct{}) string {
	argStr := ""

	events.Walk(func(path string, allChildrenSelected, hasChildren bool) IterControl {
		if path == "" {
			return IterControlContinue
		}

		if allChildrenSelected {
			argStr += " --event " + path + ".\\*"
			return IterControlDontExpand
		}

		argStr += " --event " + path
		return IterControlContinue
	})

	for _, dest := range toList(dests) {
		argStr += " --event-dest " + dest
	}

	return argStr
}

func toEvTypeTree(eventTypeSpecs []string) *evTypeTree {
	tt := &evTypeTree{}

	for _, eventTypeSpec := range eventTypeSpecs {
		components := strings.Split(eventTypeSpec, ".")
		tree := tt
		for _, comp := range components {
			if comp == "*" || tree.allChildrenSelected {
				tree.allChildrenSelected = true
				break
			}

			if tree.leaves == nil {
				tree.leaves = make(map[string]*evTypeTree)
			}

			if _, ok := tree.leaves[comp]; !ok {
				tree.leaves[comp] = &evTypeTree{}
			}

			tree = tree.leaves[comp]
		}
	}

	return tt
}
