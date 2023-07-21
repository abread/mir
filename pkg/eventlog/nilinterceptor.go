package eventlog

import "github.com/filecoin-project/mir/pkg/events"

var NilInterceptor = &nilInterceptor{}

type nilInterceptor struct{}

func (i *nilInterceptor) Intercept(_ events.EventList) error {
	return nil
}
