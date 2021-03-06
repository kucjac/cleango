// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kucjac/cleango/xblob (interfaces: Bucket)

// Package mockblob is a generated GoMock package.
package mockblob

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	xblob "github.com/kucjac/cleango/xblob"
	blob "gocloud.dev/blob"
)

// MockBucket is a mock of Bucket interface.
type MockBucket struct {
	ctrl     *gomock.Controller
	recorder *MockBucketMockRecorder
}

// MockBucketMockRecorder is the mock recorder for MockBucket.
type MockBucketMockRecorder struct {
	mock *MockBucket
}

// NewMockBucket creates a new mock instance.
func NewMockBucket(ctrl *gomock.Controller) *MockBucket {
	mock := &MockBucket{ctrl: ctrl}
	mock.recorder = &MockBucketMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBucket) EXPECT() *MockBucketMockRecorder {
	return m.recorder
}

// As mocks base method.
func (m *MockBucket) As(arg0 interface{}) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "As", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// As indicates an expected call of As.
func (mr *MockBucketMockRecorder) As(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "As", reflect.TypeOf((*MockBucket)(nil).As), arg0)
}

// Attributes mocks base method.
func (m *MockBucket) Attributes(arg0 context.Context, arg1 string) (*blob.Attributes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attributes", arg0, arg1)
	ret0, _ := ret[0].(*blob.Attributes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Attributes indicates an expected call of Attributes.
func (mr *MockBucketMockRecorder) Attributes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attributes", reflect.TypeOf((*MockBucket)(nil).Attributes), arg0, arg1)
}

// Bucket mocks base method.
func (m *MockBucket) Bucket() *blob.Bucket {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bucket")
	ret0, _ := ret[0].(*blob.Bucket)
	return ret0
}

// Bucket indicates an expected call of Bucket.
func (mr *MockBucketMockRecorder) Bucket() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bucket", reflect.TypeOf((*MockBucket)(nil).Bucket))
}

// Close mocks base method.
func (m *MockBucket) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockBucketMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBucket)(nil).Close))
}

// Copy mocks base method.
func (m *MockBucket) Copy(arg0 context.Context, arg1, arg2 string, arg3 *blob.CopyOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Copy indicates an expected call of Copy.
func (mr *MockBucketMockRecorder) Copy(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockBucket)(nil).Copy), arg0, arg1, arg2, arg3)
}

// Delete mocks base method.
func (m *MockBucket) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockBucketMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBucket)(nil).Delete), arg0, arg1)
}

// ErrorAs mocks base method.
func (m *MockBucket) ErrorAs(arg0 error, arg1 interface{}) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ErrorAs", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ErrorAs indicates an expected call of ErrorAs.
func (mr *MockBucketMockRecorder) ErrorAs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErrorAs", reflect.TypeOf((*MockBucket)(nil).ErrorAs), arg0, arg1)
}

// Exists mocks base method.
func (m *MockBucket) Exists(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockBucketMockRecorder) Exists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockBucket)(nil).Exists), arg0, arg1)
}

// IsAccessible mocks base method.
func (m *MockBucket) IsAccessible(arg0 context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAccessible", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsAccessible indicates an expected call of IsAccessible.
func (mr *MockBucketMockRecorder) IsAccessible(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAccessible", reflect.TypeOf((*MockBucket)(nil).IsAccessible), arg0)
}

// List mocks base method.
func (m *MockBucket) List(arg0 *blob.ListOptions) *blob.ListIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].(*blob.ListIterator)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockBucketMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockBucket)(nil).List), arg0)
}

// ListPage mocks base method.
func (m *MockBucket) ListPage(arg0 context.Context, arg1 []byte, arg2 int, arg3 *blob.ListOptions) ([]*blob.ListObject, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPage", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*blob.ListObject)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListPage indicates an expected call of ListPage.
func (mr *MockBucketMockRecorder) ListPage(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPage", reflect.TypeOf((*MockBucket)(nil).ListPage), arg0, arg1, arg2, arg3)
}

// NewRangeReader mocks base method.
func (m *MockBucket) NewRangeReader(arg0 context.Context, arg1 string, arg2, arg3 int64, arg4 *blob.ReaderOptions) (xblob.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRangeReader", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(xblob.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRangeReader indicates an expected call of NewRangeReader.
func (mr *MockBucketMockRecorder) NewRangeReader(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRangeReader", reflect.TypeOf((*MockBucket)(nil).NewRangeReader), arg0, arg1, arg2, arg3, arg4)
}

// NewReader mocks base method.
func (m *MockBucket) NewReader(arg0 context.Context, arg1 string, arg2 *blob.ReaderOptions) (xblob.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewReader", arg0, arg1, arg2)
	ret0, _ := ret[0].(xblob.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewReader indicates an expected call of NewReader.
func (mr *MockBucketMockRecorder) NewReader(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewReader", reflect.TypeOf((*MockBucket)(nil).NewReader), arg0, arg1, arg2)
}

// NewWriter mocks base method.
func (m *MockBucket) NewWriter(arg0 context.Context, arg1 string, arg2 *blob.WriterOptions) (xblob.Writer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewWriter", arg0, arg1, arg2)
	ret0, _ := ret[0].(xblob.Writer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewWriter indicates an expected call of NewWriter.
func (mr *MockBucketMockRecorder) NewWriter(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewWriter", reflect.TypeOf((*MockBucket)(nil).NewWriter), arg0, arg1, arg2)
}

// ReadAll mocks base method.
func (m *MockBucket) ReadAll(arg0 context.Context, arg1 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAll", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAll indicates an expected call of ReadAll.
func (mr *MockBucketMockRecorder) ReadAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAll", reflect.TypeOf((*MockBucket)(nil).ReadAll), arg0, arg1)
}

// SignedURL mocks base method.
func (m *MockBucket) SignedURL(arg0 context.Context, arg1 string, arg2 *blob.SignedURLOptions) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignedURL", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignedURL indicates an expected call of SignedURL.
func (mr *MockBucketMockRecorder) SignedURL(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignedURL", reflect.TypeOf((*MockBucket)(nil).SignedURL), arg0, arg1, arg2)
}

// WriteAll mocks base method.
func (m *MockBucket) WriteAll(arg0 context.Context, arg1 string, arg2 []byte, arg3 *blob.WriterOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteAll", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteAll indicates an expected call of WriteAll.
func (mr *MockBucketMockRecorder) WriteAll(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAll", reflect.TypeOf((*MockBucket)(nil).WriteAll), arg0, arg1, arg2, arg3)
}
