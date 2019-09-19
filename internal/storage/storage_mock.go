// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

package storage

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockService) EXPECT() *MockServiceMockRecorder {
	return _m.recorder
}

// AddUser mocks base method
func (_m *MockService) AddUser(ctx context.Context, name string) (string, error) {
	ret := _m.ctrl.Call(_m, "AddUser", ctx, name)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser
func (_mr *MockServiceMockRecorder) AddUser(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "AddUser", reflect.TypeOf((*MockService)(nil).AddUser), arg0, arg1)
}

// GetUser mocks base method
func (_m *MockService) GetUser(ctx context.Context, id string) (*User, error) {
	ret := _m.ctrl.Call(_m, "GetUser", ctx, id)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (_mr *MockServiceMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetUser", reflect.TypeOf((*MockService)(nil).GetUser), arg0, arg1)
}

// DeleteUser mocks base method
func (_m *MockService) DeleteUser(ctx context.Context, id string) error {
	ret := _m.ctrl.Call(_m, "DeleteUser", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (_mr *MockServiceMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "DeleteUser", reflect.TypeOf((*MockService)(nil).DeleteUser), arg0, arg1)
}

// TakeUserBalance mocks base method
func (_m *MockService) TakeUserBalance(ctx context.Context, id string, points float64) error {
	ret := _m.ctrl.Call(_m, "TakeUserBalance", ctx, id, points)
	ret0, _ := ret[0].(error)
	return ret0
}

// TakeUserBalance indicates an expected call of TakeUserBalance
func (_mr *MockServiceMockRecorder) TakeUserBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "TakeUserBalance", reflect.TypeOf((*MockService)(nil).TakeUserBalance), arg0, arg1, arg2)
}

// FundUserBalance mocks base method
func (_m *MockService) FundUserBalance(ctx context.Context, id string, points float64) error {
	ret := _m.ctrl.Call(_m, "FundUserBalance", ctx, id, points)
	ret0, _ := ret[0].(error)
	return ret0
}

// FundUserBalance indicates an expected call of FundUserBalance
func (_mr *MockServiceMockRecorder) FundUserBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "FundUserBalance", reflect.TypeOf((*MockService)(nil).FundUserBalance), arg0, arg1, arg2)
}

// AddTournament mocks base method
func (_m *MockService) AddTournament(ctx context.Context, name string, deposit float64) (string, error) {
	ret := _m.ctrl.Call(_m, "AddTournament", ctx, name, deposit)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTournament indicates an expected call of AddTournament
func (_mr *MockServiceMockRecorder) AddTournament(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "AddTournament", reflect.TypeOf((*MockService)(nil).AddTournament), arg0, arg1, arg2)
}

// GetTournament mocks base method
func (_m *MockService) GetTournament(ctx context.Context, id string) (*Tournament, error) {
	ret := _m.ctrl.Call(_m, "GetTournament", ctx, id)
	ret0, _ := ret[0].(*Tournament)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTournament indicates an expected call of GetTournament
func (_mr *MockServiceMockRecorder) GetTournament(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetTournament", reflect.TypeOf((*MockService)(nil).GetTournament), arg0, arg1)
}

// DeleteTournament mocks base method
func (_m *MockService) DeleteTournament(ctx context.Context, id string) error {
	ret := _m.ctrl.Call(_m, "DeleteTournament", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTournament indicates an expected call of DeleteTournament
func (_mr *MockServiceMockRecorder) DeleteTournament(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "DeleteTournament", reflect.TypeOf((*MockService)(nil).DeleteTournament), arg0, arg1)
}

// IncreaseTournamentPrize mocks base method
func (_m *MockService) IncreaseTournamentPrize(ctx context.Context, id string, amount float64) error {
	ret := _m.ctrl.Call(_m, "IncreaseTournamentPrize", ctx, id, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncreaseTournamentPrize indicates an expected call of IncreaseTournamentPrize
func (_mr *MockServiceMockRecorder) IncreaseTournamentPrize(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "IncreaseTournamentPrize", reflect.TypeOf((*MockService)(nil).IncreaseTournamentPrize), arg0, arg1, arg2)
}

// DecreaseTournamentPrize mocks base method
func (_m *MockService) DecreaseTournamentPrize(ctx context.Context, id string, amount float64) error {
	ret := _m.ctrl.Call(_m, "DecreaseTournamentPrize", ctx, id, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecreaseTournamentPrize indicates an expected call of DecreaseTournamentPrize
func (_mr *MockServiceMockRecorder) DecreaseTournamentPrize(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "DecreaseTournamentPrize", reflect.TypeOf((*MockService)(nil).DecreaseTournamentPrize), arg0, arg1, arg2)
}

// SetTournamentWinner mocks base method
func (_m *MockService) SetTournamentWinner(ctx context.Context, tournamentID string, userID string) error {
	ret := _m.ctrl.Call(_m, "SetTournamentWinner", ctx, tournamentID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTournamentWinner indicates an expected call of SetTournamentWinner
func (_mr *MockServiceMockRecorder) SetTournamentWinner(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "SetTournamentWinner", reflect.TypeOf((*MockService)(nil).SetTournamentWinner), arg0, arg1, arg2)
}

// SetTournamentStatus mocks base method
func (_m *MockService) SetTournamentStatus(ctx context.Context, tournamentID string, status TournamentStatus) error {
	ret := _m.ctrl.Call(_m, "SetTournamentStatus", ctx, tournamentID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTournamentStatus indicates an expected call of SetTournamentStatus
func (_mr *MockServiceMockRecorder) SetTournamentStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "SetTournamentStatus", reflect.TypeOf((*MockService)(nil).SetTournamentStatus), arg0, arg1, arg2)
}

// AddUserToTournamentList mocks base method
func (_m *MockService) AddUserToTournamentList(ctx context.Context, tournamentID string, userID string) error {
	ret := _m.ctrl.Call(_m, "AddUserToTournamentList", ctx, tournamentID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToTournamentList indicates an expected call of AddUserToTournamentList
func (_mr *MockServiceMockRecorder) AddUserToTournamentList(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "AddUserToTournamentList", reflect.TypeOf((*MockService)(nil).AddUserToTournamentList), arg0, arg1, arg2)
}

// JoinTournament mocks base method
func (_m *MockService) JoinTournament(ctx context.Context, tournamentID string, userID string) error {
	ret := _m.ctrl.Call(_m, "JoinTournament", ctx, tournamentID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// JoinTournament indicates an expected call of JoinTournament
func (_mr *MockServiceMockRecorder) JoinTournament(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "JoinTournament", reflect.TypeOf((*MockService)(nil).JoinTournament), arg0, arg1, arg2)
}

// FinishTournament mocks base method
func (_m *MockService) FinishTournament(ctx context.Context, tournamentID string, winnerUserID string) error {
	ret := _m.ctrl.Call(_m, "FinishTournament", ctx, tournamentID, winnerUserID)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinishTournament indicates an expected call of FinishTournament
func (_mr *MockServiceMockRecorder) FinishTournament(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "FinishTournament", reflect.TypeOf((*MockService)(nil).FinishTournament), arg0, arg1, arg2)
}
