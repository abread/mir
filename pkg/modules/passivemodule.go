package modules

import (
	"runtime/debug"

	es "github.com/go-errors/errors"

	"github.com/filecoin-project/mir/pkg/events"
	t "github.com/filecoin-project/mir/pkg/types"
)

type PassiveModule interface {
	Module

	// ApplyEvents applies a list of input events to the module, making it advance its state
	// and returns a (potentially empty) list of output events that the application of the input events results in.
	ApplyEvents(events events.EventList) (events.EventList, error)
}

func RoutedModule(rootID t.ModuleID, root PassiveModule, subRouter PassiveModule) PassiveModule {
	return &routedModule{
		rootID:    rootID,
		root:      root,
		subRouter: subRouter,
	}
}

type routedModule struct {
	rootID    t.ModuleID
	root      PassiveModule
	subRouter PassiveModule
}

func (m *routedModule) ImplementsModule() {}
func (m *routedModule) ApplyEvents(evs events.EventList) (events.EventList, error) {
	rootEvsIn := events.EmptyList()
	subRouterEvsIn := events.EmptyList()

	it := evs.Iterator()
	for ev := it.Next(); ev != nil; ev = it.Next() {
		if ev.DestModule == m.rootID {
			rootEvsIn.PushBack(ev)
		} else {
			subRouterEvsIn.PushBack(ev)
		}
	}

	rootEvsOut, rootErr := ApplyAllSafely(m.root, rootEvsIn)
	subRouterEvsOut, subRouterErr := ApplyAllSafely(m.subRouter, subRouterEvsIn)

	if subRouterErr != nil {
		return events.EmptyList(), subRouterErr
	} else if rootErr != nil {
		return events.EmptyList(), rootErr
	}

	rootEvsOut.PushBackList(subRouterEvsOut)

	return rootEvsOut, nil
}

func MultiApplyModule(subs []PassiveModule) PassiveModule {
	return &multiApplyModule{
		subs: subs,
	}
}

type multiApplyModule struct {
	subs []PassiveModule
}

func (m *multiApplyModule) ImplementsModule() {}
func (m *multiApplyModule) ApplyEvents(evs events.EventList) (events.EventList, error) {
	allEvsOut := events.EmptyList()

	for _, sub := range m.subs {
		evsOut, err := ApplyAllSafely(sub, evs)
		if err != nil {
			return events.EmptyList(), err
		}

		allEvsOut.PushBackList(evsOut)
	}

	return allEvsOut, nil
}

func ApplyAllSafely(m PassiveModule, evs events.EventList) (result events.EventList, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rErr, ok := r.(error); ok {
				err = es.Errorf("event application panicked: %w\nStack trace:\n%s", rErr, string(debug.Stack()))
			} else {
				err = es.Errorf("event application panicked: %v\nStack trace:\n%s", r, string(debug.Stack()))
			}
		}
	}()

	return m.ApplyEvents(evs)
}
