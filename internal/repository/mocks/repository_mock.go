// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	model "avito-internship/internal/model"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTransaction is a mock of Transaction interface.
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction.
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance.
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// SendCoin mocks base method.
func (m *MockTransaction) SendCoin(ctx context.Context, fromUser, toUser, amount int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoin", ctx, fromUser, toUser, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoin indicates an expected call of SendCoin.
func (mr *MockTransactionMockRecorder) SendCoin(ctx, fromUser, toUser, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoin", reflect.TypeOf((*MockTransaction)(nil).SendCoin), ctx, fromUser, toUser, amount)
}

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUser) CreateUser(ctx context.Context, user model.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUser)(nil).CreateUser), ctx, user)
}

// GetUserByID mocks base method.
func (m *MockUser) GetUserByID(ctx context.Context, userID int) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserMockRecorder) GetUserByID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUser)(nil).GetUserByID), ctx, userID)
}

// GetUserByUsername mocks base method.
func (m *MockUser) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", ctx, username)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUserMockRecorder) GetUserByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUser)(nil).GetUserByUsername), ctx, username)
}

// GetUsersInfo mocks base method.
func (m *MockUser) GetUsersInfo(ctx context.Context, userID int) (int, []model.Purchases, []model.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersInfo", ctx, userID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]model.Purchases)
	ret2, _ := ret[2].([]model.Transaction)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetUsersInfo indicates an expected call of GetUsersInfo.
func (mr *MockUserMockRecorder) GetUsersInfo(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersInfo", reflect.TypeOf((*MockUser)(nil).GetUsersInfo), ctx, userID)
}

// MockBuyer is a mock of Buyer interface.
type MockBuyer struct {
	ctrl     *gomock.Controller
	recorder *MockBuyerMockRecorder
}

// MockBuyerMockRecorder is the mock recorder for MockBuyer.
type MockBuyerMockRecorder struct {
	mock *MockBuyer
}

// NewMockBuyer creates a new mock instance.
func NewMockBuyer(ctrl *gomock.Controller) *MockBuyer {
	mock := &MockBuyer{ctrl: ctrl}
	mock.recorder = &MockBuyerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuyer) EXPECT() *MockBuyerMockRecorder {
	return m.recorder
}

// Buy mocks base method.
func (m *MockBuyer) Buy(ctx context.Context, userID int, item model.Merch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Buy", ctx, userID, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Buy indicates an expected call of Buy.
func (mr *MockBuyerMockRecorder) Buy(ctx, userID, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Buy", reflect.TypeOf((*MockBuyer)(nil).Buy), ctx, userID, item)
}

// GetItem mocks base method.
func (m *MockBuyer) GetItem(ctx context.Context, itemName string) (model.Merch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", ctx, itemName)
	ret0, _ := ret[0].(model.Merch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockBuyerMockRecorder) GetItem(ctx, itemName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockBuyer)(nil).GetItem), ctx, itemName)
}
