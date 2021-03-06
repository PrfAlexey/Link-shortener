// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_pkg is a generated GoMock package.
package mock_pkg

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
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

// GetURL mocks base method.
func (m *MockService) GetURL(link string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", link)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockServiceMockRecorder) GetURL(link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockService)(nil).GetURL), link)
}

// SaveURL mocks base method.
func (m *MockService) SaveURL(URL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveURL", URL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveURL indicates an expected call of SaveURL.
func (mr *MockServiceMockRecorder) SaveURL(URL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveURL", reflect.TypeOf((*MockService)(nil).SaveURL), URL)
}
