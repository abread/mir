package localclock

import (
	"time"

	"github.com/filecoin-project/mir/pkg/pb/eventpb"
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

func AttachTs(event *eventpb.Event) *eventpb.Event {
	event.LocalTs = Now().Nanoseconds()
	return event
}
