// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mockfactory is a generated GoMock package.
package mockfactory

import (
	reflect "reflect"

	model "github.com/toshiykst/go-layerd-architecture/app/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockUserFactory is a mock of UserFactory interface.
type MockUserFactory struct {
	ctrl     *gomock.Controller
	recorder *MockUserFactoryMockRecorder
}

// MockUserFactoryMockRecorder is the mock recorder for MockUserFactory.
type MockUserFactoryMockRecorder struct {
	mock *MockUserFactory
}

// NewMockUserFactory creates a new mock instance.
func NewMockUserFactory(ctrl *gomock.Controller) *MockUserFactory {
	mock := &MockUserFactory{ctrl: ctrl}
	mock.recorder = &MockUserFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserFactory) EXPECT() *MockUserFactoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserFactory) Create(name, email string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name, email)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserFactoryMockRecorder) Create(name, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserFactory)(nil).Create), name, email)
}
