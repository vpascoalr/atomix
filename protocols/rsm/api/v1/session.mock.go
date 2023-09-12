// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/vpascoalr/atomix/protocols/rsm/api/v1 (interfaces: SessionServer)

// Package v1 is a generated GoMock package.
package v1

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSessionServer is a mock of SessionServer interface.
type MockSessionServer struct {
	ctrl     *gomock.Controller
	recorder *MockSessionServerMockRecorder
}

// MockSessionServerMockRecorder is the mock recorder for MockSessionServer.
type MockSessionServerMockRecorder struct {
	mock *MockSessionServer
}

// NewMockSessionServer creates a new mock instance.
func NewMockSessionServer(ctrl *gomock.Controller) *MockSessionServer {
	mock := &MockSessionServer{ctrl: ctrl}
	mock.recorder = &MockSessionServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionServer) EXPECT() *MockSessionServerMockRecorder {
	return m.recorder
}

// ClosePrimitive mocks base method.
func (m *MockSessionServer) ClosePrimitive(arg0 context.Context, arg1 *ClosePrimitiveRequest) (*ClosePrimitiveResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClosePrimitive", arg0, arg1)
	ret0, _ := ret[0].(*ClosePrimitiveResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClosePrimitive indicates an expected call of ClosePrimitive.
func (mr *MockSessionServerMockRecorder) ClosePrimitive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClosePrimitive", reflect.TypeOf((*MockSessionServer)(nil).ClosePrimitive), arg0, arg1)
}

// CreatePrimitive mocks base method.
func (m *MockSessionServer) CreatePrimitive(arg0 context.Context, arg1 *CreatePrimitiveRequest) (*CreatePrimitiveResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePrimitive", arg0, arg1)
	ret0, _ := ret[0].(*CreatePrimitiveResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePrimitive indicates an expected call of CreatePrimitive.
func (mr *MockSessionServerMockRecorder) CreatePrimitive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePrimitive", reflect.TypeOf((*MockSessionServer)(nil).CreatePrimitive), arg0, arg1)
}