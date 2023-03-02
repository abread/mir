package modules

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/types"
)

func RoutedModule(ctx context.Context, rootID types.ModuleID, root PassiveModule, subRouter ActiveModule, outChan chan *events.EventList) ActiveModule {
	mod := &routedModule{
		rootID:    string(rootID),
		root:      root,
		subRouter: subRouter,
		outChan:   outChan,
	}

	return mod
}

type routedModule struct {
	rootID    string
	root      PassiveModule
	subRouter ActiveModule

	outChan chan *events.EventList
}

func (m *routedModule) ImplementsModule() {}

func (m *routedModule) EventsOut() <-chan *events.EventList {
	return m.outChan
}

func (m *routedModule) ApplyEvents(ctx context.Context, evs *events.EventList) error {
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

	rootEvsOut, rootErr := passiveApplySafely(m.root, rootEvsIn)
	if rootErr != nil {
		return rootErr
	}
	m.outChan <- rootEvsOut

	subRouterErr := activeApplySafely(ctx, m.subRouter, subRouterEvsIn)
	if subRouterErr != nil {
		return subRouterErr
	}

	return nil
}

func passiveApplySafely(m PassiveModule, evs *events.EventList) (result *events.EventList, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rErr, ok := r.(error); ok {
				err = fmt.Errorf("module panicked: %w\nStack trace:\n%s", rErr, string(debug.Stack()))
			} else {
				err = fmt.Errorf("module panicked: %v\nStack trace:\n%s", r, string(debug.Stack()))
			}
		}
	}()

	return m.ApplyEvents(evs)
}

func activeApplySafely(ctx context.Context, m ActiveModule, evs *events.EventList) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if rErr, ok := r.(error); ok {
				err = fmt.Errorf("module panicked: %w\nStack trace:\n%s", rErr, string(debug.Stack()))
			} else {
				err = fmt.Errorf("module panicked: %v\nStack trace:\n%s", r, string(debug.Stack()))
			}
		}
	}()

	return m.ApplyEvents(ctx, evs)
}
