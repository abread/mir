package modules

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/types"
)

func RoutedModule(ctx context.Context, rootID types.ModuleID, root PassiveModule, subRouter ActiveModule) ActiveModule {
	outChan := make(chan *events.EventList)

	mod := &routedModule{
		rootID:    string(rootID),
		root:      root,
		subRouter: subRouter,
		doneC:     ctx.Done(),
		outChan:   outChan,
	}

	mod.runningChildren.Add(1)
	go importEvents(ctx, &mod.runningChildren, outChan, subRouter)

	// close output when done
	go func() {
		<-mod.doneC
		mod.runningChildren.Wait()

		close(mod.outChan)
	}()

	return mod
}

type routedModule struct {
	rootID    string
	root      PassiveModule
	subRouter ActiveModule

	outChan         chan *events.EventList
	doneC           <-chan struct{}
	runningChildren sync.WaitGroup
}

func (m *routedModule) ImplementsModule() {}

func (m *routedModule) EventsOut() <-chan *events.EventList {
	return m.outChan
}

func (m *routedModule) ApplyEvents(ctx context.Context, evs *events.EventList) error {
	// the cleanup goroutine will close the output channel, so we must block it with this
	m.runningChildren.Add(1)
	defer m.runningChildren.Done()

	select {
	case <-m.doneC:
		// we are done here
		// perhaps the output channel was already closed by the cleanup goroutine
		return nil
	default:
		// not done yet, keep going
	}

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

func importEvents(ctx context.Context, wg *sync.WaitGroup, outChan chan *events.EventList, mod ActiveModule) {
	defer wg.Done()

	input := mod.EventsOut()
	doneC := ctx.Done()

	for {
		select {
		case <-doneC:
			return
		case evs := <-input:
			outChan <- evs
		}
	}
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
