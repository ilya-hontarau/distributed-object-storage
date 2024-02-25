// Code generated by MockGen. DO NOT EDIT.
// Source: internal/gateway/gateway.go
//
// Generated by this command:
//
//	mockgen -source internal/gateway/gateway.go -destination internal/mock/gateway.go -package mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	io "io"
	fs "io/fs"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStorageNode is a mock of StorageNode interface.
type MockStorageNode struct {
	ctrl     *gomock.Controller
	recorder *MockStorageNodeMockRecorder
}

// MockStorageNodeMockRecorder is the mock recorder for MockStorageNode.
type MockStorageNodeMockRecorder struct {
	mock *MockStorageNode
}

// NewMockStorageNode creates a new mock instance.
func NewMockStorageNode(ctrl *gomock.Controller) *MockStorageNode {
	mock := &MockStorageNode{ctrl: ctrl}
	mock.recorder = &MockStorageNodeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageNode) EXPECT() *MockStorageNodeMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockStorageNode) Download(ctx context.Context, id string) (io.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", ctx, id)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockStorageNodeMockRecorder) Download(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockStorageNode)(nil).Download), ctx, id)
}

// Upload mocks base method.
func (m *MockStorageNode) Upload(ctx context.Context, id string, file fs.File) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", ctx, id, file)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upload indicates an expected call of Upload.
func (mr *MockStorageNodeMockRecorder) Upload(ctx, id, file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockStorageNode)(nil).Upload), ctx, id, file)
}
