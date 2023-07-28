package modules

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

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

var ref = time.Now()

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

	tStart := time.Since(ref)
	nIn := rootEvsIn.Len()
	rootEvsOut, rootErr := applyAllSafely(m.root, rootEvsIn)
	fmt.Fprintf(os.Stderr, "tr,%s,root,%d,%d\n", m.rootID, nIn, time.Since(ref)-tStart)

	tStart = time.Since(ref)
	nIn = subRouterEvsIn.Len()
	subRouterEvsOut, subRouterErr := applyAllSafely(m.subRouter, subRouterEvsIn)
	fmt.Fprintf(os.Stderr, "tr,%s,sub,%d,%d\n", m.rootID, nIn, time.Since(ref)-tStart)

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

	nIn := evs.Len()
	for i, sub := range m.subs {
		tStart := time.Since(ref)
		evsOut, err := applyAllSafely(sub, evs)
		fmt.Fprintf(os.Stderr, "tm,%d,%d,%d\n", i, nIn, time.Since(ref)-tStart)

		if err != nil {
			return events.EmptyList(), err
		}

		allEvsOut.PushBackList(evsOut)
	}

	return allEvsOut, nil
}

func applyAllSafely(m PassiveModule, evs events.EventList) (result events.EventList, err error) {
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
