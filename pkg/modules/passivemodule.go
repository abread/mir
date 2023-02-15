package modules

import (
	"fmt"
	"runtime/debug"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/types"
)

type PassiveModule interface {
	Module

	// ApplyEvents applies a list of input events to the module, making it advance its state
	// and returns a (potentially empty) list of output events that the application of the input events results in.
	ApplyEvents(events *events.EventList) (*events.EventList, error)
}

func RoutedModule(rootID types.ModuleID, root PassiveModule, subRouter PassiveModule) PassiveModule {
	return &routedModule{
		rootID:    string(rootID),
		root:      root,
		subRouter: subRouter,
	}
}

type routedModule struct {
	rootID    string
	root      PassiveModule
	subRouter PassiveModule
}

func (m *routedModule) ImplementsModule() {}
func (m *routedModule) ApplyEvents(evs *events.EventList) (*events.EventList, error) {
	rootEvsIn := &events.EventList{}
	subRouterEvsIn := &events.EventList{}

	it := evs.Iterator()
	for ev := it.Next(); ev != nil; ev = it.Next() {
		if ev.DestModule == m.rootID {
			rootEvsIn.PushBack(ev)
		} else {
			subRouterEvsIn.PushBack(ev)
		}
	}

	rootOutChan := make(chan *events.EventList)
	rootErrChan := make(chan error)

	go func(rootModule PassiveModule, evs *events.EventList) {
		out, err := applyAllSafely(rootModule, evs)
		rootOutChan <- out
		rootErrChan <- err
	}(m.root, rootEvsIn)

	subRouterEvsOut, subRouterErr := applyAllSafely(m.subRouter, subRouterEvsIn)

	// Attention: Those (unbuffered) channels must be read by the aggregator in the same order
	//            as they are being written here, otherwise the system gets stuck.
	rootEvsOut := <-rootOutChan
	rootErr := <-rootErrChan

	if subRouterErr != nil {
		return nil, subRouterErr
	} else if rootErr != nil {
		return nil, rootErr
	}

	return rootEvsOut.PushBackList(subRouterEvsOut), nil
}

func applyAllSafely(m PassiveModule, evs *events.EventList) (result *events.EventList, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rErr, ok := r.(error); ok {
				err = fmt.Errorf("event application panicked: %w\nStack trace:\n%s", rErr, string(debug.Stack()))
			} else {
				err = fmt.Errorf("event application panicked: %v\nStack trace:\n%s", r, string(debug.Stack()))
			}
		}
	}()

	return m.ApplyEvents(evs)
}
