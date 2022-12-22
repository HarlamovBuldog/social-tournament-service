package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	storage2 "github.com/HarlamovBuldog/social-tournament-service/internal/pkg/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateNewUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randomUserID := primitive.NewObjectID().Hex()
	mock := storage2.NewMockService(ctrl)
	mock.EXPECT().AddUser(gomock.Any(), gomock.Eq("Gennadiy")).Times(1).Return(randomUserID, nil)

	enc, err := json.Marshal(userName{
		Name: "Gennadiy",
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", "/user", b)

	actualURLPath := req.URL.Path
	require.Equal("/user", actualURLPath, "The two URL pathes should be the same")

	actualReqMethod := req.Method
	require.Equal("POST", actualReqMethod, "The two request methods should be the same")

	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.createNewUser(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusOK, actualCode, "The two http codes should be the same")

	var actualUserID userID
	err = json.NewDecoder(w.Result().Body).Decode(&actualUserID)
	require.NoError(err)
	require.Equal(userID{ID: randomUserID}, actualUserID, "The two bodies shoud be the same")
}

func TestCreateNewUser_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	mock.EXPECT().AddUser(gomock.Any(), gomock.Eq("Vasiliy")).Times(1).Return("", errors.New("insert doc to collection"))

	enc, err := json.Marshal(userName{
		Name: "Vasiliy",
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", "/user", b)
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.createNewUser(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusInternalServerError, actualCode, "The two http codes should be the same")
}

func TestCreateNewUser_Bad_Req(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	mock.EXPECT().AddUser(gomock.Any(), gomock.Any()).Times(0)

	req := httptest.NewRequest("POST", "/user", nil)
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.createNewUser(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestGetUserInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	expectedUserID := primitive.NewObjectID()
	expectedUser := &storage2.User{ID: expectedUserID, Name: "Gennadiy", Balance: 0}
	mock.EXPECT().GetUser(gomock.Any(), gomock.Eq(expectedUserID.Hex())).Times(1).Return(expectedUser, nil)

	expectedURLPath := fmt.Sprintf("/user/%s", expectedUserID.Hex())
	req := httptest.NewRequest("GET", expectedURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{"id": expectedUserID.Hex()})

	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.getUserInfo(w, req)

	actualCode := w.Result().StatusCode
	require := require.New(t)
	require.Equal(http.StatusOK, actualCode, "The two http codes should be the same")

	var actualUser storage2.User
	err := json.NewDecoder(w.Result().Body).Decode(&actualUser)
	require.NoError(err)
	require.Equal(*expectedUser, actualUser, "The two objects shoud be the same")
}

func TestGetUserInfo_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	expectedUserID := primitive.NewObjectID()
	mock.EXPECT().GetUser(gomock.Any(), gomock.Eq(expectedUserID.Hex())).
		Times(1).Return(nil, errors.New("get doc from collection"))

	expectedURLPath := fmt.Sprintf("/user/%s", expectedUserID.Hex())
	req := httptest.NewRequest("GET", expectedURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{"id": expectedUserID.Hex()})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.getUserInfo(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusInternalServerError, actualCode, "The two http codes should be the same")
}

func TestGetUserInfo_Bad_Req(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	mock.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)

	expectedURLPath := fmt.Sprintf("/user/%s", "garbage")
	req := httptest.NewRequest("GET", expectedURLPath, nil)
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.getUserInfo(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestRemoveUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	userID := primitive.NewObjectID()
	mock.EXPECT().DeleteUser(gomock.Any(), gomock.Eq(userID.Hex())).Times(1).Return(nil)

	expectedURLPath := fmt.Sprintf("/user/%s", userID.Hex())
	req := httptest.NewRequest("DELETE", expectedURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{"id": userID.Hex()})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.removeUser(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusOK, actualCode, "The two http codes should be the same")
}

func TestRemoveUser_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	userID := primitive.NewObjectID()
	mock.EXPECT().DeleteUser(gomock.Any(), gomock.Eq(userID.Hex())).Times(1).Return(errors.New("delete doc from collection"))

	expectedURLPath := fmt.Sprintf("/user/%s", userID.Hex())
	req := httptest.NewRequest("DELETE", expectedURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{"id": userID.Hex()})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.removeUser(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusInternalServerError, actualCode, "The two http codes should be the same")
}

func TestRemoveUser_Bad_Req(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	mock.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Times(0)

	expectedURLPath := fmt.Sprintf("/user/%s", "garbage")
	req := httptest.NewRequest("DELETE", expectedURLPath, nil)
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.removeUser(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestTakeUserBonusPoints_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	userID := primitive.NewObjectID()
	userPointsToTake := 200.0
	mock.EXPECT().TakeUserBalance(gomock.Any(), gomock.Eq(userID.Hex()), gomock.Eq(userPointsToTake)).
		Times(1).Return(nil)

	enc, err := json.Marshal(userPoints{
		Points: userPointsToTake,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	expectedURLPath := fmt.Sprintf("/user/%s/take", userID.Hex())
	req := httptest.NewRequest("POST", expectedURLPath, b)
	req = mux.SetURLVars(req, map[string]string{"id": userID.Hex()})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.takeUserBonusPoints(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusOK, actualCode, "The two http codes should be the same")
}

func TestTakeUserBonusPoints_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	userID := primitive.NewObjectID()
	userPointsToTake := 200.0
	mock.EXPECT().TakeUserBalance(gomock.Any(), gomock.Eq(userID.Hex()), gomock.Eq(userPointsToTake)).
		Times(1).Return(errors.New("update doc in collection"))

	enc, err := json.Marshal(userPoints{
		Points: userPointsToTake,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	expectedURLPath := fmt.Sprintf("/user/%s/take", userID.Hex())
	req := httptest.NewRequest("POST", expectedURLPath, b)
	req = mux.SetURLVars(req, map[string]string{"id": userID.Hex()})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.takeUserBonusPoints(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusInternalServerError, actualCode, "The two http codes should be the same")
}

func TestTakeUserBonusPoints_Bad_Req_Body(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	mock.EXPECT().TakeUserBalance(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	userID := primitive.NewObjectID()
	expectedURLPath := fmt.Sprintf("/user/%s/take", userID.Hex())
	req := httptest.NewRequest("POST", expectedURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{"id": userID.Hex()})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.takeUserBonusPoints(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestTakeUserBonusPoints_Bad_Req_URL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	mock.EXPECT().TakeUserBalance(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	userPointsToTake := 200.0
	enc, err := json.Marshal(userPoints{
		Points: userPointsToTake,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	expectedURLPath := fmt.Sprintf("/user/%s/take", "garbage")
	req := httptest.NewRequest("POST", expectedURLPath, b)
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.takeUserBonusPoints(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestAddUserBonusPoints_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	userID := primitive.NewObjectID()
	userPointsToAdd := 200.0
	mock.EXPECT().FundUserBalance(gomock.Any(), gomock.Eq(userID.Hex()), gomock.Eq(userPointsToAdd)).
		Times(1).Return(nil)

	enc, err := json.Marshal(userPoints{
		Points: userPointsToAdd,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	expectedURLPath := fmt.Sprintf("/user/%s/fund", userID.Hex())
	req := httptest.NewRequest("POST", expectedURLPath, b)
	req = mux.SetURLVars(req, map[string]string{"id": userID.Hex()})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.addUserBonusPoints(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusOK, actualCode, "The two http codes should be the same")
}

func TestAddUserBonusPoints_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	userID := primitive.NewObjectID()
	userPointsToAdd := 200.0
	mock.EXPECT().FundUserBalance(gomock.Any(), gomock.Eq(userID.Hex()), gomock.Eq(userPointsToAdd)).
		Times(1).Return(errors.New("update doc in collection"))

	enc, err := json.Marshal(userPoints{
		Points: userPointsToAdd,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	expectedURLPath := fmt.Sprintf("/user/%s/fund", userID.Hex())
	req := httptest.NewRequest("POST", expectedURLPath, b)
	req = mux.SetURLVars(req, map[string]string{"id": userID.Hex()})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.addUserBonusPoints(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusInternalServerError, actualCode, "The two http codes should be the same")
}

func TestAddUserBonusPoints_Bad_Req_Body(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	mock.EXPECT().FundUserBalance(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	userID := primitive.NewObjectID()
	expectedURLPath := fmt.Sprintf("/user/%s/fund", userID.Hex())
	req := httptest.NewRequest("POST", expectedURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{"id": userID.Hex()})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.addUserBonusPoints(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestAddUserBonusPoints_Bad_Req_URL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage2.NewMockService(ctrl)
	mock.EXPECT().FundUserBalance(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	userPointsToAdd := 200.0
	enc, err := json.Marshal(userPoints{
		Points: userPointsToAdd,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	expectedURLPath := fmt.Sprintf("/user/%s/fund", "garbage")
	req := httptest.NewRequest("POST", expectedURLPath, b)
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.addUserBonusPoints(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusBadRequest, actualCode, "The two http codes should be the same")
}
