// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kucjac/cleango/eventsource (interfaces: Store)

// Package storemock is a generated GoMock package.
package storemock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	eventsource "github.com/kucjac/cleango/eventsource"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockStore) Commit(arg0 context.Context, arg1 eventsource.Aggregate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockStoreMockRecorder) Commit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockStore)(nil).Commit), arg0, arg1)
}

// LoadEventStream mocks base method.
func (m *MockStore) LoadEventStream(arg0 context.Context, arg1 eventsource.Aggregate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadEventStream", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadEventStream indicates an expected call of LoadEventStream.
func (mr *MockStoreMockRecorder) LoadEventStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadEventStream", reflect.TypeOf((*MockStore)(nil).LoadEventStream), arg0, arg1)
}

// LoadEventStreamWithSnapshot mocks base method.
func (m *MockStore) LoadEventStreamWithSnapshot(arg0 context.Context, arg1 eventsource.Aggregate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadEventStreamWithSnapshot", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadEventStreamWithSnapshot indicates an expected call of LoadEventStreamWithSnapshot.
func (mr *MockStoreMockRecorder) LoadEventStreamWithSnapshot(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadEventStreamWithSnapshot", reflect.TypeOf((*MockStore)(nil).LoadEventStreamWithSnapshot), arg0, arg1)
}

// SaveSnapshot mocks base method.
func (m *MockStore) SaveSnapshot(arg0 context.Context, arg1 eventsource.Aggregate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSnapshot", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSnapshot indicates an expected call of SaveSnapshot.
func (mr *MockStoreMockRecorder) SaveSnapshot(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSnapshot", reflect.TypeOf((*MockStore)(nil).SaveSnapshot), arg0, arg1)
}
