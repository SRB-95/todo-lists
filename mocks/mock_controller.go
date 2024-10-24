// Code generated by MockGen. DO NOT EDIT.
// Source: todo-lists/controllers (interfaces: IController)

// Package controller is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockIController is a mock of IController interface.
type MockIController struct {
	ctrl     *gomock.Controller
	recorder *MockIControllerMockRecorder
}

// MockIControllerMockRecorder is the mock recorder for MockIController.
type MockIControllerMockRecorder struct {
	mock *MockIController
}

// NewMockIController creates a new mock instance.
func NewMockIController(ctrl *gomock.Controller) *MockIController {
	mock := &MockIController{ctrl: ctrl}
	mock.recorder = &MockIControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIController) EXPECT() *MockIControllerMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockIController) CreateTask(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateTask", arg0)
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockIControllerMockRecorder) CreateTask(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockIController)(nil).CreateTask), arg0)
}

// DeleteTask mocks base method.
func (m *MockIController) DeleteTask(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteTask", arg0)
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockIControllerMockRecorder) DeleteTask(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockIController)(nil).DeleteTask), arg0)
}

// FilterTasksByDeadline mocks base method.
func (m *MockIController) FilterTasksByDeadline(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FilterTasksByDeadline", arg0)
}

// FilterTasksByDeadline indicates an expected call of FilterTasksByDeadline.
func (mr *MockIControllerMockRecorder) FilterTasksByDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterTasksByDeadline", reflect.TypeOf((*MockIController)(nil).FilterTasksByDeadline), arg0)
}

// GetTaskById mocks base method.
func (m *MockIController) GetTaskById(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetTaskById", arg0)
}

// GetTaskById indicates an expected call of GetTaskById.
func (mr *MockIControllerMockRecorder) GetTaskById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskById", reflect.TypeOf((*MockIController)(nil).GetTaskById), arg0)
}

// GetTaskByTag mocks base method.
func (m *MockIController) GetTaskByTag(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetTaskByTag", arg0)
}

// GetTaskByTag indicates an expected call of GetTaskByTag.
func (mr *MockIControllerMockRecorder) GetTaskByTag(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskByTag", reflect.TypeOf((*MockIController)(nil).GetTaskByTag), arg0)
}

// GetTasks mocks base method.
func (m *MockIController) GetTasks(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetTasks", arg0)
}

// GetTasks indicates an expected call of GetTasks.
func (mr *MockIControllerMockRecorder) GetTasks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasks", reflect.TypeOf((*MockIController)(nil).GetTasks), arg0)
}

// SearchTasks mocks base method.
func (m *MockIController) SearchTasks(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SearchTasks", arg0)
}

// SearchTasks indicates an expected call of SearchTasks.
func (mr *MockIControllerMockRecorder) SearchTasks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchTasks", reflect.TypeOf((*MockIController)(nil).SearchTasks), arg0)
}

// UpdateTask mocks base method.
func (m *MockIController) UpdateTask(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateTask", arg0)
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockIControllerMockRecorder) UpdateTask(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockIController)(nil).UpdateTask), arg0)
}
