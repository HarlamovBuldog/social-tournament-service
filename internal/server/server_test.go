package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HarlamovBuldog/social-tournament-service/internal/storage"
)

func TestCreateNewUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := storage.NewMockService(ctrl)
	mock.EXPECT().AddUser(gomock.Any(), gomock.Eq("Gennadiy")).Times(1).Return("code_str", nil)
	s := NewServer(mock)
	enc, err := json.Marshal(userNameJSON{
		Name: "Gennadiy",
	})
	require.NoError(t, err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("GET", "/user", b)
	w := httptest.NewRecorder()
	s.createNewUser(w, req)

	actualCode := w.Result().StatusCode
	assert := assert.New(t)
	assert.Equal(http.StatusOK, actualCode, "The two http codes should be the same")

	var actualUserID userIDJSON
	err = json.NewDecoder(w.Result().Body).Decode(&actualUserID)
	require.NoError(t, err)
	assert.Equal(actualUserID, userIDJSON{ID: "code_str"}, "The two bodies shoud be the same")
}

func TestCreateNewUser_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().AddUser(gomock.Any(), gomock.Eq("Vasiliy")).Times(1).Return("", errors.New("insert doc to collection"))
	s := NewServer(mock)

	enc, err := json.Marshal(userNameJSON{
		Name: "Vasiliy",
	})
	require.NoError(t, err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("GET", "/user", b)
	w := httptest.NewRecorder()
	s.createNewUser(w, req)

	actualCode := w.Result().StatusCode
	assert := assert.New(t)
	assert.Equal(http.StatusInternalServerError, actualCode, "The two http codes should be the same")
}

func TestCreateNewUser_Bad_Req(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().AddUser(gomock.Any(), gomock.Any()).Times(0)
	s := NewServer(mock)

	req := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	s.createNewUser(w, req)

	actualCode := w.Result().StatusCode
	assert := assert.New(t)
	assert.Equal(http.StatusBadRequest, actualCode, "The two http codes should be the same")
}
