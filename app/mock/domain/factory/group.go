// Code generated by MockGen. DO NOT EDIT.
// Source: group.go

// Package mockfactory is a generated GoMock package.
package mockfactory

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

// MockGroupFactory is a mock of GroupFactory interface.
type MockGroupFactory struct {
	ctrl     *gomock.Controller
	recorder *MockGroupFactoryMockRecorder
}

// MockGroupFactoryMockRecorder is the mock recorder for MockGroupFactory.
type MockGroupFactoryMockRecorder struct {
	mock *MockGroupFactory
}

// NewMockGroupFactory creates a new mock instance.
func NewMockGroupFactory(ctrl *gomock.Controller) *MockGroupFactory {
	mock := &MockGroupFactory{ctrl: ctrl}
	mock.recorder = &MockGroupFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupFactory) EXPECT() *MockGroupFactoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockGroupFactory) Create(name string) (*model.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name)
	ret0, _ := ret[0].(*model.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockGroupFactoryMockRecorder) Create(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockGroupFactory)(nil).Create), name)
}