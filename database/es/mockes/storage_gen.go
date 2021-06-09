// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kucjac/cleango/database/es (interfaces: Storage)

// Package mockes is a generated GoMock package.
package mockes

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cgerrors "github.com/kucjac/cleango/cgerrors"
	es "github.com/kucjac/cleango/database/es"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// As mocks base method.
func (m *MockStorage) As(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "As", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// As indicates an expected call of As.
func (mr *MockStorageMockRecorder) As(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "As", reflect.TypeOf((*MockStorage)(nil).As), arg0)
}

// ErrorCode mocks base method.
func (m *MockStorage) ErrorCode(arg0 error) cgerrors.ErrorCode {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ErrorCode", arg0)
	ret0, _ := ret[0].(cgerrors.ErrorCode)
	return ret0
}

// ErrorCode indicates an expected call of ErrorCode.
func (mr *MockStorageMockRecorder) ErrorCode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErrorCode", reflect.TypeOf((*MockStorage)(nil).ErrorCode), arg0)
}

// GetSnapshot mocks base method.
func (m *MockStorage) GetSnapshot(arg0 context.Context, arg1, arg2 string, arg3 int64) (*es.Snapshot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSnapshot", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*es.Snapshot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSnapshot indicates an expected call of GetSnapshot.
func (mr *MockStorageMockRecorder) GetSnapshot(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSnapshot", reflect.TypeOf((*MockStorage)(nil).GetSnapshot), arg0, arg1, arg2, arg3)
}

// ListEvents mocks base method.
func (m *MockStorage) ListEvents(arg0 context.Context, arg1, arg2 string) ([]*es.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEvents", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*es.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEvents indicates an expected call of ListEvents.
func (mr *MockStorageMockRecorder) ListEvents(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEvents", reflect.TypeOf((*MockStorage)(nil).ListEvents), arg0, arg1, arg2)
}

// ListEventsAfterRevision mocks base method.
func (m *MockStorage) ListEventsAfterRevision(arg0 context.Context, arg1, arg2 string, arg3 int64) ([]*es.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEventsAfterRevision", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*es.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEventsAfterRevision indicates an expected call of ListEventsAfterRevision.
func (mr *MockStorageMockRecorder) ListEventsAfterRevision(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEventsAfterRevision", reflect.TypeOf((*MockStorage)(nil).ListEventsAfterRevision), arg0, arg1, arg2, arg3)
}

// NewCursor mocks base method.
func (m *MockStorage) NewCursor(arg0 context.Context, arg1 string, arg2 int64) (es.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCursor", arg0, arg1, arg2)
	ret0, _ := ret[0].(es.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewCursor indicates an expected call of NewCursor.
func (mr *MockStorageMockRecorder) NewCursor(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCursor", reflect.TypeOf((*MockStorage)(nil).NewCursor), arg0, arg1, arg2)
}

// SaveEvents mocks base method.
func (m *MockStorage) SaveEvents(arg0 context.Context, arg1 []*es.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveEvents", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveEvents indicates an expected call of SaveEvents.
func (mr *MockStorageMockRecorder) SaveEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveEvents", reflect.TypeOf((*MockStorage)(nil).SaveEvents), arg0, arg1)
}

// SaveSnapshot mocks base method.
func (m *MockStorage) SaveSnapshot(arg0 context.Context, arg1 *es.Snapshot) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSnapshot", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSnapshot indicates an expected call of SaveSnapshot.
func (mr *MockStorageMockRecorder) SaveSnapshot(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSnapshot", reflect.TypeOf((*MockStorage)(nil).SaveSnapshot), arg0, arg1)
}

// StreamEvents mocks base method.
func (m *MockStorage) StreamEvents(arg0 context.Context, arg1 *es.StreamEventsRequest) (<-chan *es.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StreamEvents", arg0, arg1)
	ret0, _ := ret[0].(<-chan *es.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StreamEvents indicates an expected call of StreamEvents.
func (mr *MockStorageMockRecorder) StreamEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamEvents", reflect.TypeOf((*MockStorage)(nil).StreamEvents), arg0, arg1)
}