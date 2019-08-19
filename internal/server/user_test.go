package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/HarlamovBuldog/social-tournament-service/internal/storage"
)

func TestCreateNewUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := storage.NewMockService(ctrl)
	mock.EXPECT().AddUser(gomock.Any(), gomock.Eq("Gennadiy")).Times(1).Return("code_str", nil)

	enc, err := json.Marshal(userNameJSON{
		Name: "Gennadiy",
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", "/user", b)
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.createNewUser(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusOK, actualCode, "The two http codes should be the same")

	var actualUserID userIDJSON
	err = json.NewDecoder(w.Result().Body).Decode(&actualUserID)
	require.NoError(err)
	require.Equal(userIDJSON{ID: "code_str"}, actualUserID, "The two bodies shoud be the same")
}

func TestCreateNewUser_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().AddUser(gomock.Any(), gomock.Eq("Vasiliy")).Times(1).Return("", errors.New("insert doc to collection"))

	enc, err := json.Marshal(userNameJSON{
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

	mock := storage.NewMockService(ctrl)
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

	mock := storage.NewMockService(ctrl)
	expectedUserID := primitive.NewObjectID()
	expectedUser := &storage.User{ID: expectedUserID, Name: "Gennadiy", Balance: 0}
	mock.EXPECT().GetUser(gomock.Any(), gomock.Eq(expectedUserID.Hex())).Times(1).Return(expectedUser, nil)

	req := httptest.NewRequest("GET", "/user/", nil)
	req = mux.SetURLVars(req, map[string]string{"id": expectedUserID.Hex()})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.getUserInfo(w, req)

	actualCode := w.Result().StatusCode
	require := require.New(t)
	require.Equal(http.StatusOK, actualCode, "The two http codes should be the same")

	var actualUser storage.User
	err := json.NewDecoder(w.Result().Body).Decode(&actualUser)
	require.NoError(err)
	require.Equal(*expectedUser, actualUser, "The two objects shoud be the same")
}

func TestGetUserInfo_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	expectedUserID := primitive.NewObjectID()
	mock.EXPECT().GetUser(gomock.Any(), gomock.Eq(expectedUserID.Hex())).
		Times(1).Return(nil, errors.New("get doc from collection"))

	req := httptest.NewRequest("GET", "/user", nil)
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

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)

	req := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.getUserInfo(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestRemoveUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	userID := primitive.NewObjectID()
	mock.EXPECT().DeleteUser(gomock.Any(), gomock.Eq(userID.Hex())).Times(1).Return(nil)

	req := httptest.NewRequest("DELETE", "/user/", nil)
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

	mock := storage.NewMockService(ctrl)
	userID := primitive.NewObjectID()
	mock.EXPECT().DeleteUser(gomock.Any(), gomock.Eq(userID.Hex())).Times(1).Return(errors.New("delete doc from collection"))

	req := httptest.NewRequest("DELETE", "/user/", nil)
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

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Times(0)

	req := httptest.NewRequest("DELETE", "/user", nil)
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.removeUser(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusBadRequest, actualCode, "The two http codes should be the same")
}
