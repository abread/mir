package mir

import (
	"context"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/modules/mock_modules"
	"github.com/filecoin-project/mir/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNode_Run(t *testing.T) {
	testCases := map[string]func(t *testing.T) (m modules.Modules, done <-chan struct{}){
		"InitEvents": func(t *testing.T) (modules.Modules, <-chan struct{}) {
			ctrl := gomock.NewController(t)
			mockModule1 := mock_modules.NewMockPassiveModule(ctrl)
			mockModule2 := mock_modules.NewMockPassiveModule(ctrl)

			var wg sync.WaitGroup
			wg.Add(2)

			mockModule1.EXPECT().ApplyEvents(events.ListOf(events.Init("mock1"))).
				Do(func(_ *events.EventList) { wg.Done() }).
				Return(events.EmptyList(), nil)
			mockModule2.EXPECT().ApplyEvents(events.ListOf(events.Init("mock2"))).
				Do(func(_ *events.EventList) { wg.Done() }).
				Return(events.EmptyList(), nil)

			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			m := map[types.ModuleID]modules.Module{
				"mock1": mockModule1,
				"mock2": mockModule2,
			}
			return m, done
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			m, tcDone := tc(t)

			logger := logging.ConsoleWarnLogger
			n, err := NewNode(
				"testnode",
				&NodeConfig{Logger: logger},
				m,
				nil,
			)

			assert.Nil(t, err)
			ctx, stopNode := context.WithCancel(context.Background())

			nodeStopped := make(chan struct{})
			go func() {
				err := n.Run(ctx)
				assert.Equal(t, ErrStopped, err)
				close(nodeStopped)
			}()

			// Wait until either the test case is done or a 2 seconds deadline
			select {
			case <-tcDone:
			case <-time.After(2 * time.Second):
			}

			stopNode()
			<-nodeStopped
		})
	}
}