// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kucjac/cleango/database/es (interfaces: Storage,TxStorage)

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

// BeginTx mocks base method.
func (m *MockStorage) BeginTx(arg0 context.Context) (es.TxStorage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", arg0)
	ret0, _ := ret[0].(es.TxStorage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockStorageMockRecorder) BeginTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockStorage)(nil).BeginTx), arg0)
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

// MockTxStorage is a mock of TxStorage interface.
type MockTxStorage struct {
	ctrl     *gomock.Controller
	recorder *MockTxStorageMockRecorder
}

// MockTxStorageMockRecorder is the mock recorder for MockTxStorage.
type MockTxStorageMockRecorder struct {
	mock *MockTxStorage
}

// NewMockTxStorage creates a new mock instance.
func NewMockTxStorage(ctrl *gomock.Controller) *MockTxStorage {
	mock := &MockTxStorage{ctrl: ctrl}
	mock.recorder = &MockTxStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxStorage) EXPECT() *MockTxStorageMockRecorder {
	return m.recorder
}

// As mocks base method.
func (m *MockTxStorage) As(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "As", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// As indicates an expected call of As.
func (mr *MockTxStorageMockRecorder) As(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "As", reflect.TypeOf((*MockTxStorage)(nil).As), arg0)
}

// Commit mocks base method.
func (m *MockTxStorage) Commit(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockTxStorageMockRecorder) Commit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockTxStorage)(nil).Commit), arg0)
}

// ErrorCode mocks base method.
func (m *MockTxStorage) ErrorCode(arg0 error) cgerrors.ErrorCode {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ErrorCode", arg0)
	ret0, _ := ret[0].(cgerrors.ErrorCode)
	return ret0
}

// ErrorCode indicates an expected call of ErrorCode.
func (mr *MockTxStorageMockRecorder) ErrorCode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErrorCode", reflect.TypeOf((*MockTxStorage)(nil).ErrorCode), arg0)
}

// GetSnapshot mocks base method.
func (m *MockTxStorage) GetSnapshot(arg0 context.Context, arg1, arg2 string, arg3 int64) (*es.Snapshot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSnapshot", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*es.Snapshot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSnapshot indicates an expected call of GetSnapshot.
func (mr *MockTxStorageMockRecorder) GetSnapshot(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSnapshot", reflect.TypeOf((*MockTxStorage)(nil).GetSnapshot), arg0, arg1, arg2, arg3)
}

// ListEvents mocks base method.
func (m *MockTxStorage) ListEvents(arg0 context.Context, arg1, arg2 string) ([]*es.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEvents", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*es.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEvents indicates an expected call of ListEvents.
func (mr *MockTxStorageMockRecorder) ListEvents(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEvents", reflect.TypeOf((*MockTxStorage)(nil).ListEvents), arg0, arg1, arg2)
}

// ListEventsAfterRevision mocks base method.
func (m *MockTxStorage) ListEventsAfterRevision(arg0 context.Context, arg1, arg2 string, arg3 int64) ([]*es.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEventsAfterRevision", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*es.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEventsAfterRevision indicates an expected call of ListEventsAfterRevision.
func (mr *MockTxStorageMockRecorder) ListEventsAfterRevision(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEventsAfterRevision", reflect.TypeOf((*MockTxStorage)(nil).ListEventsAfterRevision), arg0, arg1, arg2, arg3)
}

// Rollback mocks base method.
func (m *MockTxStorage) Rollback(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockTxStorageMockRecorder) Rollback(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockTxStorage)(nil).Rollback), arg0)
}

// SaveEvents mocks base method.
func (m *MockTxStorage) SaveEvents(arg0 context.Context, arg1 []*es.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveEvents", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveEvents indicates an expected call of SaveEvents.
func (mr *MockTxStorageMockRecorder) SaveEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveEvents", reflect.TypeOf((*MockTxStorage)(nil).SaveEvents), arg0, arg1)
}

// SaveSnapshot mocks base method.
func (m *MockTxStorage) SaveSnapshot(arg0 context.Context, arg1 *es.Snapshot) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSnapshot", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSnapshot indicates an expected call of SaveSnapshot.
func (mr *MockTxStorageMockRecorder) SaveSnapshot(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSnapshot", reflect.TypeOf((*MockTxStorage)(nil).SaveSnapshot), arg0, arg1)
}

// StreamEvents mocks base method.
func (m *MockTxStorage) StreamEvents(arg0 context.Context, arg1 *es.StreamEventsRequest) (<-chan *es.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StreamEvents", arg0, arg1)
	ret0, _ := ret[0].(<-chan *es.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StreamEvents indicates an expected call of StreamEvents.
func (mr *MockTxStorageMockRecorder) StreamEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamEvents", reflect.TypeOf((*MockTxStorage)(nil).StreamEvents), arg0, arg1)
}
