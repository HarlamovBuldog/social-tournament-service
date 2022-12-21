package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HarlamovBuldog/social-tournament-service/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateNewTournament_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTournamentID := primitive.NewObjectID().Hex()
	expectedTournamentName := "Tournament_1"
	expectedTournamentDeposit := 1500.0

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().AddTournament(gomock.Any(), gomock.Eq(expectedTournamentName),
		gomock.Eq(expectedTournamentDeposit)).Times(1).Return(expectedTournamentID, nil)

	enc, err := json.Marshal(tournament{
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", "/tournament", b)

	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.createNewTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusOK, actualCode, "The two http codes should be the same")

	var actualtournamentID tournamentID
	err = json.NewDecoder(w.Result().Body).Decode(&actualtournamentID)
	require.NoError(err)
	require.Equal(tournamentID{ID: expectedTournamentID}, actualtournamentID, "The two bodies shoud be the same")
}

func TestCreateNewTournament_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTournamentName := "Tournament_1"
	expectedTournamentDeposit := 1500.0
	expectedError := errors.New("add doc to collection")

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().AddTournament(gomock.Any(), gomock.Eq(expectedTournamentName),
		gomock.Eq(expectedTournamentDeposit)).Times(1).Return("", expectedError)

	enc, err := json.Marshal(tournament{
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", "/tournament", b)

	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.createNewTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusInternalServerError, actualCode, "The two http codes should be the same")
}

func TestCreateNewTournament_Bad_Req(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().AddTournament(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	req := httptest.NewRequest("POST", "/tournament", nil)
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.createNewTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestGetTournamentInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	expectedTournamentID := primitive.NewObjectID()
	expectedTournamentName := "Tournament_1"
	expectedTournamentDeposit := 1500.0
	expectedTournament := &storage.Tournament{
		ID:      expectedTournamentID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit}
	mock.EXPECT().GetTournament(gomock.Any(), gomock.Eq(expectedTournamentID.Hex())).
		Times(1).Return(expectedTournament, nil)

	expectedURLPath := fmt.Sprintf("/tournament/%s", expectedTournamentID.Hex())
	req := httptest.NewRequest("GET", expectedURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{"id": expectedTournamentID.Hex()})

	w := httptest.NewRecorder()
	s := NewServer(mock)
	s.getTournamentInfo(w, req)

	actualCode := w.Result().StatusCode
	require := require.New(t)
	require.Equal(http.StatusOK, actualCode, "The two http codes should be the same")

	var actualTournament storage.Tournament
	err := json.NewDecoder(w.Result().Body).Decode(&actualTournament)
	require.NoError(err)
	require.Equal(*expectedTournament, actualTournament, "The two objects shoud be the same")
}

func TestGetTournamentInfo_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	expectedTournamentID := primitive.NewObjectID()
	expectedError := errors.New("get doc from collection")
	mock.EXPECT().GetTournament(gomock.Any(), gomock.Eq(expectedTournamentID.Hex())).
		Times(1).Return(nil, expectedError)

	expectedURLPath := fmt.Sprintf("/tournament/%s", expectedTournamentID.Hex())
	req := httptest.NewRequest("GET", expectedURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{"id": expectedTournamentID.Hex()})

	w := httptest.NewRecorder()
	s := NewServer(mock)
	s.getTournamentInfo(w, req)

	actualCode := w.Result().StatusCode
	require := require.New(t)
	require.Equal(http.StatusInternalServerError, actualCode, "The two http codes should be the same")

	var actualTournament storage.Tournament
	err := json.NewDecoder(w.Result().Body).Decode(&actualTournament)
	require.Error(err)
}

func TestGetTournamentInfo_Bad_Req(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().GetTournament(gomock.Any(), gomock.Any()).Times(0)

	badURLPath := fmt.Sprint("/tournament/")
	req := httptest.NewRequest("GET", badURLPath, nil)

	w := httptest.NewRecorder()
	s := NewServer(mock)
	s.getTournamentInfo(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestJoinTournament_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	expectedTournamentID := primitive.NewObjectID().Hex()
	expectedUserID := primitive.NewObjectID().Hex()
	mock.EXPECT().JoinTournament(gomock.Any(), gomock.Eq(expectedTournamentID), gomock.Eq(expectedUserID)).
		Times(1).Return(nil)

	expectedURLPath := fmt.Sprintf("/tournament/%s/join", expectedTournamentID)
	enc, err := json.Marshal(userID{
		ID: expectedUserID,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", expectedURLPath, b)
	req = mux.SetURLVars(req, map[string]string{"id": expectedTournamentID})

	w := httptest.NewRecorder()
	s := NewServer(mock)
	s.joinTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusOK, actualCode, "The two http codes should be the same")
}

func TestJoinTournament_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	expectedTournamentID := primitive.NewObjectID().Hex()
	expectedUserID := primitive.NewObjectID().Hex()
	expectedError := errors.New("any error cause it's transaction")
	mock.EXPECT().JoinTournament(gomock.Any(), gomock.Eq(expectedTournamentID), gomock.Eq(expectedUserID)).
		Times(1).Return(expectedError)

	expectedURLPath := fmt.Sprintf("/tournament/%s/join", expectedTournamentID)
	enc, err := json.Marshal(userID{
		ID: expectedUserID,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", expectedURLPath, b)
	req = mux.SetURLVars(req, map[string]string{"id": expectedTournamentID})

	w := httptest.NewRecorder()
	s := NewServer(mock)
	s.joinTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusInternalServerError, actualCode, "The two http codes should be the same")
}

func TestJoinTournament_Bad_Req(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().JoinTournament(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	badURLPath := fmt.Sprint("/tournament/")
	expectedUserID := primitive.NewObjectID().Hex()
	enc, err := json.Marshal(userID{
		ID: expectedUserID,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", badURLPath, b)

	w := httptest.NewRecorder()
	s := NewServer(mock)
	s.joinTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestFinishTournament_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	tournamentID := primitive.NewObjectID().Hex()
	winnerUsrID := primitive.NewObjectID().Hex()
	mock.EXPECT().FinishTournament(gomock.Any(), gomock.Eq(tournamentID), gomock.Eq(winnerUsrID)).
		Times(1).Return(nil)

	expectedURLPath := fmt.Sprintf("/tournament/%s/finish", tournamentID)
	enc, err := json.Marshal(winnerUserID{
		ID: winnerUsrID,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", expectedURLPath, b)
	req = mux.SetURLVars(req, map[string]string{"id": tournamentID})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.finishTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusOK, actualCode, "The two http codes should be the same")
}

func TestFinishTournament_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	tournamentID := primitive.NewObjectID().Hex()
	winnerUsrID := primitive.NewObjectID().Hex()
	expectedError := errors.New("any error cause it's transaction")
	mock.EXPECT().FinishTournament(gomock.Any(), gomock.Eq(tournamentID), gomock.Eq(winnerUsrID)).
		Times(1).Return(expectedError)

	expectedURLPath := fmt.Sprintf("/tournament/%s/finish", tournamentID)
	enc, err := json.Marshal(winnerUserID{
		ID: winnerUsrID,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", expectedURLPath, b)
	req = mux.SetURLVars(req, map[string]string{"id": tournamentID})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.finishTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusInternalServerError, actualCode, "The two http codes should be the same")
}

func TestFinishTournament_Bad_Req(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().FinishTournament(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	badURLPath := "/tournament/finish"
	winnerUsrID := primitive.NewObjectID().Hex()
	enc, err := json.Marshal(winnerUserID{
		ID: winnerUsrID,
	})
	require := require.New(t)
	require.NoError(err)

	b := bytes.NewBuffer(enc)
	req := httptest.NewRequest("POST", badURLPath, b)
	req = mux.SetURLVars(req, map[string]string{})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.finishTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(http.StatusBadRequest, actualCode, "The two http codes should be the same")
}

func TestCancelTournament_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	tournamentID := primitive.NewObjectID().Hex()
	mock.EXPECT().DeleteTournament(gomock.Any(), gomock.Eq(tournamentID)).Times(1).Return(nil)

	expectedURLPath := fmt.Sprintf("/tournament/%s", tournamentID)
	req := httptest.NewRequest("DELETE", expectedURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{"id": tournamentID})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.cancelTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusOK, actualCode, "The two http codes should be the same")
}

func TestCancelTournament_DB_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	tournamentID := primitive.NewObjectID().Hex()
	expectedError := errors.New("any error cause it's transaction")
	mock.EXPECT().DeleteTournament(gomock.Any(), gomock.Eq(tournamentID)).Times(1).Return(expectedError)

	expectedURLPath := fmt.Sprintf("/tournament/%s", tournamentID)
	req := httptest.NewRequest("DELETE", expectedURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{"id": tournamentID})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.cancelTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusInternalServerError, actualCode, "The two http codes should be the same")
}

func TestCancelTournament_Bad_Req(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := storage.NewMockService(ctrl)
	mock.EXPECT().DeleteTournament(gomock.Any(), gomock.Any()).Times(0)

	badURLPath := "/tournament"
	req := httptest.NewRequest("DELETE", badURLPath, nil)
	req = mux.SetURLVars(req, map[string]string{})
	w := httptest.NewRecorder()

	s := NewServer(mock)
	s.cancelTournament(w, req)

	actualCode := w.Result().StatusCode
	require.Equal(t, http.StatusBadRequest, actualCode, "The two http codes should be the same")
}
