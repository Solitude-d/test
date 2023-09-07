// Code generated by MockGen. DO NOT EDIT.
// Source: webook/internal/repository/dao/user.go

// Package daomocks is a generated GoMock package.
package daomocks

import (
	context "context"
	reflect "reflect"
	dao "test/webook/internal/repository/dao"

	gomock "go.uber.org/mock/gomock"
)

// MockUserDao is a mock of UserDao interface.
type MockUserDao struct {
	ctrl     *gomock.Controller
	recorder *MockUserDaoMockRecorder
}

// MockUserDaoMockRecorder is the mock recorder for MockUserDao.
type MockUserDaoMockRecorder struct {
	mock *MockUserDao
}

// NewMockUserDao creates a new mock instance.
func NewMockUserDao(ctrl *gomock.Controller) *MockUserDao {
	mock := &MockUserDao{ctrl: ctrl}
	mock.recorder = &MockUserDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserDao) EXPECT() *MockUserDaoMockRecorder {
	return m.recorder
}

// FindByEmail mocks base method.
func (m *MockUserDao) FindByEmail(ctx context.Context, email string) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserDaoMockRecorder) FindByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserDao)(nil).FindByEmail), ctx, email)
}

// FindByID mocks base method.
func (m *MockUserDao) FindByID(ctx context.Context, id int64) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockUserDaoMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockUserDao)(nil).FindByID), ctx, id)
}

// FindByPhone mocks base method.
func (m *MockUserDao) FindByPhone(ctx context.Context, phone string) (dao.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPhone", ctx, phone)
	ret0, _ := ret[0].(dao.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPhone indicates an expected call of FindByPhone.
func (mr *MockUserDaoMockRecorder) FindByPhone(ctx, phone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPhone", reflect.TypeOf((*MockUserDao)(nil).FindByPhone), ctx, phone)
}

// Insert mocks base method.
func (m *MockUserDao) Insert(ctx context.Context, u dao.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockUserDaoMockRecorder) Insert(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserDao)(nil).Insert), ctx, u)
}

// Update mocks base method.
func (m *MockUserDao) Update(ctx context.Context, u dao.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserDaoMockRecorder) Update(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserDao)(nil).Update), ctx, u)
}