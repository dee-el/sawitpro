// Code generated by MockGen. DO NOT EDIT.
// Source: service/user.interface.go

// Package mock_service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for Mock
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetMe mocks base method.
func (m *MockService) GetMe(ctx context.Context, token string) (GetMeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMe", ctx, token)
	ret0, _ := ret[0].(GetMeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMe indicates an expected call of GetMe.
func (mr *MockServiceMockRecorder) GetMe(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMe", reflect.TypeOf((*MockService)(nil).GetMe), ctx, token)
}

// Login mocks base method.
func (m *MockService) Login(ctx context.Context, params LoginRequest) (LoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, params)
	ret0, _ := ret[0].(LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockServiceMockRecorder) Login(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockService)(nil).Login), ctx, params)
}

// Register mocks base method.
func (m *MockService) Register(ctx context.Context, params RegisterRequest) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, params)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockServiceMockRecorder) Register(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockService)(nil).Register), ctx, params)
}

// UpdateMe mocks base method.
func (m *MockService) UpdateMe(ctx context.Context, token string, params UpdateRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMe", ctx, token, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMe indicates an expected call of UpdateMe.
func (mr *MockServiceMockRecorder) UpdateMe(ctx, token, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMe", reflect.TypeOf((*MockService)(nil).UpdateMe), ctx, token, params)
}
