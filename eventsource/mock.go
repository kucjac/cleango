package eventsource

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/kucjac/cleango/codec"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
	*aggBaseSetter
	setEventsCommitted bool
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller, eventCodec, snapCodec codec.Codec, idGen IdGenerator) *MockStore {
	mock := &MockStore{
		ctrl:          ctrl,
		aggBaseSetter: newAggregateBaseSetter(eventCodec, snapCodec, idGen),
	}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// NewDefaultMockStore creates new default mock store.
func NewDefaultMockStore(ctrl *gomock.Controller) *MockStore {
	c := codec.JSON()
	mock := &MockStore{
		ctrl:          ctrl,
		aggBaseSetter: newAggregateBaseSetter(c, c, UUIDGenerator{}),
	}
	mock.recorder = &MockStoreMockRecorder{mock: mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockStore) Commit(arg0 context.Context, arg1 Aggregate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", arg0, arg1)
	ret0, _ := ret[0].(error)

	if ret0 == nil && m.setEventsCommitted {
		b := arg1.AggBase()
		b.committedEvents, b.uncommittedEvents = b.uncommittedEvents, nil
	}
	m.setEventsCommitted = false
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockStoreMockRecorder) Commit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	mr.mock.setEventsCommitted = true
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockStore)(nil).Commit), arg0, arg1)
}

// LoadEventStream mocks base method.
func (m *MockStore) LoadEventStream(arg0 context.Context, arg1 Aggregate) error {
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
func (m *MockStore) LoadEventStreamWithSnapshot(arg0 context.Context, arg1 Aggregate) error {
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
func (m *MockStore) SaveSnapshot(arg0 context.Context, arg1 Aggregate) error {
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

// StreamAggregates mocks base method.
func (m *MockStore) StreamAggregates(arg0 context.Context, arg1 string, arg2 int64, arg3 AggregateFactory) (<-chan Aggregate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StreamAggregates", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(<-chan Aggregate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StreamAggregates indicates an expected call of StreamAggregates.
func (mr *MockStoreMockRecorder) StreamAggregates(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamAggregates", reflect.TypeOf((*MockStore)(nil).StreamAggregates), arg0, arg1, arg2, arg3)
}

// StreamProjections mocks base method.
func (m *MockStore) StreamProjections(arg0 context.Context, arg1 string, arg2 int64, arg3 ProjectionFactory) (<-chan Projection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StreamProjections", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(<-chan Projection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StreamProjections indicates an expected call of StreamProjections.
func (mr *MockStoreMockRecorder) StreamProjections(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamProjections", reflect.TypeOf((*MockStore)(nil).StreamProjections), arg0, arg1, arg2, arg3)
}
