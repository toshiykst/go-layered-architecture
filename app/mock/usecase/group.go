// Code generated by MockGen. DO NOT EDIT.
// Source: group.go

// Package mockusecase is a generated GoMock package.
package mockusecase

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/toshiykst/go-layerd-architecture/app/usecase/dto"
)

// MockGroupUsecase is a mock of GroupUsecase interface.
type MockGroupUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockGroupUsecaseMockRecorder
}

// MockGroupUsecaseMockRecorder is the mock recorder for MockGroupUsecase.
type MockGroupUsecaseMockRecorder struct {
	mock *MockGroupUsecase
}

// NewMockGroupUsecase creates a new mock instance.
func NewMockGroupUsecase(ctrl *gomock.Controller) *MockGroupUsecase {
	mock := &MockGroupUsecase{ctrl: ctrl}
	mock.recorder = &MockGroupUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupUsecase) EXPECT() *MockGroupUsecaseMockRecorder {
	return m.recorder
}

// CreateGroup mocks base method.
func (m *MockGroupUsecase) CreateGroup(in *dto.CreateGroupInput) (*dto.CreateGroupOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", in)
	ret0, _ := ret[0].(*dto.CreateGroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockGroupUsecaseMockRecorder) CreateGroup(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockGroupUsecase)(nil).CreateGroup), in)
}

// DeleteGroup mocks base method.
func (m *MockGroupUsecase) DeleteGroup(in *dto.DeleteGroupInput) (*dto.DeleteGroupOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroup", in)
	ret0, _ := ret[0].(*dto.DeleteGroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteGroup indicates an expected call of DeleteGroup.
func (mr *MockGroupUsecaseMockRecorder) DeleteGroup(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockGroupUsecase)(nil).DeleteGroup), in)
}

// GetGroup mocks base method.
func (m *MockGroupUsecase) GetGroup(in *dto.GetGroupInput) (*dto.GetGroupOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", in)
	ret0, _ := ret[0].(*dto.GetGroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup.
func (mr *MockGroupUsecaseMockRecorder) GetGroup(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockGroupUsecase)(nil).GetGroup), in)
}

// GetGroups mocks base method.
func (m *MockGroupUsecase) GetGroups(in *dto.GetGroupsInput) (*dto.GetGroupsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroups", in)
	ret0, _ := ret[0].(*dto.GetGroupsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroups indicates an expected call of GetGroups.
func (mr *MockGroupUsecaseMockRecorder) GetGroups(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroups", reflect.TypeOf((*MockGroupUsecase)(nil).GetGroups), in)
}

// UpdateGroup mocks base method.
func (m *MockGroupUsecase) UpdateGroup(in *dto.UpdateGroupInput) (*dto.UpdateGroupOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroup", in)
	ret0, _ := ret[0].(*dto.UpdateGroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateGroup indicates an expected call of UpdateGroup.
func (mr *MockGroupUsecaseMockRecorder) UpdateGroup(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroup", reflect.TypeOf((*MockGroupUsecase)(nil).UpdateGroup), in)
}