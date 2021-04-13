// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kucjac/cleango/messages/pubsub (interfaces: Publisher,Subscriber)

// Package mockps is a generated GoMock package.
package mockps

import (
	context "context"
	reflect "reflect"

	message "github.com/ThreeDotsLabs/watermill/message"
	gomock "github.com/golang/mock/gomock"
	pubsub "github.com/kucjac/cleango/messages/pubsub"
)

// MockPublisher is a mock of Publisher interface.
type MockPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockPublisherMockRecorder
}

// MockPublisherMockRecorder is the mock recorder for MockPublisher.
type MockPublisherMockRecorder struct {
	mock *MockPublisher
}

// NewMockPublisher creates a new mock instance.
func NewMockPublisher(ctrl *gomock.Controller) *MockPublisher {
	mock := &MockPublisher{ctrl: ctrl}
	mock.recorder = &MockPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPublisher) EXPECT() *MockPublisherMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockPublisher) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockPublisherMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPublisher)(nil).Close))
}

// Publish mocks base method.
func (m *MockPublisher) Publish(arg0 string, arg1 ...*message.Message) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Publish", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockPublisherMockRecorder) Publish(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockPublisher)(nil).Publish), varargs...)
}

// PublishMessage mocks base method.
func (m *MockPublisher) PublishMessage(arg0 string, arg1 interface{}, arg2 ...pubsub.PublishOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PublishMessage", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishMessage indicates an expected call of PublishMessage.
func (mr *MockPublisherMockRecorder) PublishMessage(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishMessage", reflect.TypeOf((*MockPublisher)(nil).PublishMessage), varargs...)
}

// MockSubscriber is a mock of Subscriber interface.
type MockSubscriber struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriberMockRecorder
}

// MockSubscriberMockRecorder is the mock recorder for MockSubscriber.
type MockSubscriberMockRecorder struct {
	mock *MockSubscriber
}

// NewMockSubscriber creates a new mock instance.
func NewMockSubscriber(ctrl *gomock.Controller) *MockSubscriber {
	mock := &MockSubscriber{ctrl: ctrl}
	mock.recorder = &MockSubscriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriber) EXPECT() *MockSubscriberMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockSubscriber) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockSubscriberMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSubscriber)(nil).Close))
}

// Subscribe mocks base method.
func (m *MockSubscriber) Subscribe(arg0 context.Context, arg1 string) (<-chan *message.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0, arg1)
	ret0, _ := ret[0].(<-chan *message.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockSubscriberMockRecorder) Subscribe(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockSubscriber)(nil).Subscribe), arg0, arg1)
}