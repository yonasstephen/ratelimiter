// Code generated by MockGen. DO NOT EDIT.
// Source: repository/repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// IncrementByKey mocks base method.
func (m *MockRepository) IncrementByKey(key string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncrementByKey", key)
}

// IncrementByKey indicates an expected call of IncrementByKey.
func (mr *MockRepositoryMockRecorder) IncrementByKey(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementByKey", reflect.TypeOf((*MockRepository)(nil).IncrementByKey), key)
}
