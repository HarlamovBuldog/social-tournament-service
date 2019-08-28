package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/HarlamovBuldog/social-tournament-service/internal/storage"
	gomock "github.com/golang/mock/gomock"
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

	enc, err := json.Marshal(tournamentInitJSON{
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

	var actualtournamentID tournamentIDJSON
	err = json.NewDecoder(w.Result().Body).Decode(&actualtournamentID)
	require.NoError(err)
	require.Equal(tournamentIDJSON{ID: expectedTournamentID}, actualtournamentID, "The two bodies shoud be the same")
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

	enc, err := json.Marshal(tournamentInitJSON{
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
