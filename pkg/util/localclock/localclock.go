package localclock

import (
	"time"

	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

var clock time.Time

func init() {
	clock = time.Now()
}

func Now() time.Duration {
	return time.Since(clock)
}

func RefTime() time.Time {
	return clock
}

func AttachTS(event *eventpbtypes.Event) *eventpbtypes.Event {
	event.LocalTs = Now().Nanoseconds()
	return event
}
