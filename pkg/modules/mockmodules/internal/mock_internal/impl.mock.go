// Code generated by MockGen. DO NOT EDIT.
// Source: ./impl.go

// Package mock_internal is a generated GoMock package.
package mock_internal

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	events "github.com/filecoin-project/mir/pkg/events"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

// MockModuleImpl is a mock of ModuleImpl interface.
type MockModuleImpl struct {
	ctrl     *gomock.Controller
	recorder *MockModuleImplMockRecorder
}

// MockModuleImplMockRecorder is the mock recorder for MockModuleImpl.
type MockModuleImplMockRecorder struct {
	mock *MockModuleImpl
}

// NewMockModuleImpl creates a new mock instance.
func NewMockModuleImpl(ctrl *gomock.Controller) *MockModuleImpl {
	mock := &MockModuleImpl{ctrl: ctrl}
	mock.recorder = &MockModuleImplMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModuleImpl) EXPECT() *MockModuleImplMockRecorder {
	return m.recorder
}

// Event mocks base method.
func (m *MockModuleImpl) Event(ev *eventpbtypes.Event) (events.EventList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Event", ev)
	ret0, _ := ret[0].(events.EventList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Event indicates an expected call of Event.
func (mr *MockModuleImplMockRecorder) Event(ev interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Event", reflect.TypeOf((*MockModuleImpl)(nil).Event), ev)
}
