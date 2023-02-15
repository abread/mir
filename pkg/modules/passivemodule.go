package modules

import (
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
	it := evs.Iterator()

	evsInner := &events.EventList{}
	evsSub := &events.EventList{}

	for ev := it.Next(); ev != nil; ev = it.Next() {
		if ev.DestModule == m.rootID {
			evsInner.PushBack(ev)
		} else {
			evsSub.PushBack(ev)
		}
	}

	innerOutChan := make(chan *events.EventList)
	innerErrChan := make(chan error)

	go func() {
		out, err := m.root.ApplyEvents(evsInner)
		if err != nil {
			innerErrChan <- err
		}
		innerOutChan <- out
	}()

	evsOut, err := m.subRouter.ApplyEvents(evsSub)

	select {
	case innerOut := <-innerOutChan:
		evsOut.PushBackList(innerOut)
	case innerErr := <-innerErrChan:
		return nil, innerErr
	}

	// we only check this error here to ensure the concurrent m.inner.ApplyEvents ends
	if err != nil {
		return nil, err
	}

	return evsOut, nil
}
