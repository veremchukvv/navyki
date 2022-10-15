// Code generated by MockGen. DO NOT EDIT.
// Source: order.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	storage "funny_test/storage"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOrderStorage is a mock of OrderStorage interface.
type MockOrderStorage struct {
	ctrl     *gomock.Controller
	recorder *MockOrderStorageMockRecorder
}

// MockOrderStorageMockRecorder is the mock recorder for MockOrderStorage.
type MockOrderStorageMockRecorder struct {
	mock *MockOrderStorage
}

// NewMockOrderStorage creates a new mock instance.
func NewMockOrderStorage(ctrl *gomock.Controller) *MockOrderStorage {
	mock := &MockOrderStorage{ctrl: ctrl}
	mock.recorder = &MockOrderStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderStorage) EXPECT() *MockOrderStorageMockRecorder {
	return m.recorder
}

// Order mocks base method.
func (m *MockOrderStorage) Order(ctx context.Context, ID int) *storage.Order {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Order", ctx, ID)
	ret0, _ := ret[0].(*storage.Order)
	return ret0
}

// Order indicates an expected call of Order.
func (mr *MockOrderStorageMockRecorder) Order(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Order", reflect.TypeOf((*MockOrderStorage)(nil).Order), ctx, ID)
}
