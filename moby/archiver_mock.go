// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/quay/claircore/moby (interfaces: Archiver)

// Package moby is a generated GoMock package.
package moby

import (
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
)

// MockArchiver is a mock of Archiver interface
type MockArchiver struct {
	ctrl     *gomock.Controller
	recorder *MockArchiverMockRecorder
}

// MockArchiverMockRecorder is the mock recorder for MockArchiver
type MockArchiverMockRecorder struct {
	mock *MockArchiver
}

// NewMockArchiver creates a new mock instance
func NewMockArchiver(ctrl *gomock.Controller) *MockArchiver {
	mock := &MockArchiver{ctrl: ctrl}
	mock.recorder = &MockArchiverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockArchiver) EXPECT() *MockArchiverMockRecorder {
	return m.recorder
}

// DecompressStream mocks base method
func (m *MockArchiver) DecompressStream(arg0 io.Reader) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecompressStream", arg0)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecompressStream indicates an expected call of DecompressStream
func (mr *MockArchiverMockRecorder) DecompressStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecompressStream", reflect.TypeOf((*MockArchiver)(nil).DecompressStream), arg0)
}
