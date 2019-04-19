// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/db/db.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	db "gitlab.com/jskswamy/nightfury/pkg/db"
	reflect "reflect"
)

// MockModel is a mock of Model interface
type MockModel struct {
	ctrl     *gomock.Controller
	recorder *MockModelMockRecorder
}

// MockModelMockRecorder is the mock recorder for MockModel
type MockModelMockRecorder struct {
	mock *MockModel
}

// NewMockModel creates a new mock instance
func NewMockModel(ctrl *gomock.Controller) *MockModel {
	mock := &MockModel{ctrl: ctrl}
	mock.recorder = &MockModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockModel) EXPECT() *MockModelMockRecorder {
	return m.recorder
}

// ID mocks base method
func (m *MockModel) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockModelMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockModel)(nil).ID))
}

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Save mocks base method
func (m *MockRepository) Save(bucketName string, model db.Model) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", bucketName, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockRepositoryMockRecorder) Save(bucketName, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepository)(nil).Save), bucketName, model)
}

// Delete mocks base method
func (m *MockRepository) Delete(bucketName string, model db.Model) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", bucketName, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepositoryMockRecorder) Delete(bucketName, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), bucketName, model)
}

// Fetch mocks base method
func (m *MockRepository) Fetch(bucketName, name string, model db.Model) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", bucketName, name, model)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch
func (mr *MockRepositoryMockRecorder) Fetch(bucketName, name, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockRepository)(nil).Fetch), bucketName, name, model)
}

// FetchAll mocks base method
func (m *MockRepository) FetchAll(bucketName string, modelFn func([]byte) (db.Model, error)) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchAll", bucketName, modelFn)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchAll indicates an expected call of FetchAll
func (mr *MockRepositoryMockRecorder) FetchAll(bucketName, modelFn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchAll", reflect.TypeOf((*MockRepository)(nil).FetchAll), bucketName, modelFn)
}

// Close mocks base method
func (m *MockRepository) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockRepositoryMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepository)(nil).Close))
}
