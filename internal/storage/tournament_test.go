package storage

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/require"
)

func TestAddTournament(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	expectedTournamentObjID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTournament := Tournament{
		ID:      expectedTournamentObjID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users:   []primitive.ObjectID{},
	}

	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	cleanUp(t)
}

func TestGetTournament(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	expectedTournamentObjID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTournament := Tournament{
		ID:      expectedTournamentObjID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users:   []primitive.ObjectID{},
	}

	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	badTournamentID := "bad_t_id"
	actualTournament, err = db.GetTournament(context.TODO(), badTournamentID)
	require.Nil(actualTournament, "the tournament object should be nil")
	require.EqualError(err, "convert string value to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'")

	notExistTournamentID := primitive.NewObjectID().Hex()
	actualTournament, err = db.GetTournament(context.TODO(), notExistTournamentID)
	require.Nil(actualTournament, "the tournament object should be nil")
	require.EqualError(err, "decode returned doc: mongo: no documents in result")

	cleanUp(t)
}

func TestDeleteTournament(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	err = db.DeleteTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	badTournamentID := "bad_t_id"
	err = db.DeleteTournament(context.TODO(), badTournamentID)
	require.EqualError(err, "convert string value to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'")

	notExistTournamentID := primitive.NewObjectID().Hex()
	err = db.DeleteTournament(context.TODO(), notExistTournamentID)
	require.EqualError(err, "delete doc from collection: DeletedCount != 1")

	cleanUp(t)
}

func TestAddUserToTournamentList(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	userID := primitive.NewObjectID()
	err = db.AddUserToTournamentList(context.TODO(), expectedTournamentID, userID.Hex())
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTournamentObjID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)

	expectedTournament := Tournament{
		ID:      expectedTournamentObjID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users:   []primitive.ObjectID{userID},
	}

	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	badTournamentID := "bad_t_id"
	err = db.AddUserToTournamentList(context.TODO(), badTournamentID, userID.Hex())
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badTournamentID))

	badUserID := "bad_user_id"
	err = db.AddUserToTournamentList(context.TODO(), expectedTournamentID, badUserID)
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badUserID))

	notExistTournamentID := primitive.NewObjectID().Hex()
	err = db.AddUserToTournamentList(context.TODO(), notExistTournamentID, userID.Hex())
	require.EqualError(err, "update doc in collection: ModifiedCount != 1")

	cleanUp(t)
}

func TestSetTournamentWinner(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	userWinnerID := primitive.NewObjectID()
	err = db.SetTournamentWinner(context.TODO(), expectedTournamentID, userWinnerID.Hex())
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTournamentObjID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)

	expectedTournament := Tournament{
		ID:      expectedTournamentObjID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users:   []primitive.ObjectID{},
		Winner:  userWinnerID,
	}

	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	badTournamentID := "bad_t_id"
	err = db.SetTournamentWinner(context.TODO(), badTournamentID, userWinnerID.Hex())
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badTournamentID))

	badUserID := "bad_user_id"
	err = db.SetTournamentWinner(context.TODO(), expectedTournamentID, badUserID)
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badUserID))

	notExistTournamentID := primitive.NewObjectID().Hex()
	err = db.SetTournamentWinner(context.TODO(), notExistTournamentID, userWinnerID.Hex())
	require.EqualError(err, "update doc in collection: ModifiedCount != 1")

	cleanUp(t)
}

func TestIncreaseTournamentPrize(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	incAmount := 1000.0
	err = db.IncreaseTournamentPrize(context.TODO(), expectedTournamentID, incAmount)
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTournamentObjID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)

	expectedTournamentPrize := 1000.0
	expectedTournament := Tournament{
		ID:      expectedTournamentObjID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users:   []primitive.ObjectID{},
		Prize:   expectedTournamentPrize,
	}

	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	badTournamentID := "bad_t_id"
	err = db.IncreaseTournamentPrize(context.TODO(), badTournamentID, incAmount)
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badTournamentID))

	notExistTournamentID := primitive.NewObjectID().Hex()
	err = db.IncreaseTournamentPrize(context.TODO(), notExistTournamentID, incAmount)
	require.EqualError(err, "update doc in collection: ModifiedCount != 1")

	cleanUp(t)
}

func TestDecreaseTournamentPrize(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	incAmount := 1000.0
	err = db.IncreaseTournamentPrize(context.TODO(), expectedTournamentID, incAmount)
	require.NoError(err)

	decAmount := 250.0
	err = db.DecreaseTournamentPrize(context.TODO(), expectedTournamentID, decAmount)
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTournamentObjID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)

	expectedTournamentPrize := 750.0
	expectedTournament := Tournament{
		ID:      expectedTournamentObjID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users:   []primitive.ObjectID{},
		Prize:   expectedTournamentPrize,
	}

	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	badTournamentID := "bad_t_id"
	err = db.DecreaseTournamentPrize(context.TODO(), badTournamentID, decAmount)
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badTournamentID))

	notExistTournamentID := primitive.NewObjectID().Hex()
	err = db.DecreaseTournamentPrize(context.TODO(), notExistTournamentID, decAmount)
	require.EqualError(err, "update doc in collection: ModifiedCount != 1")

	cleanUp(t)
}

func TestSetTournamentStatus(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	expectedStatus := StatusFinished
	err = db.SetTournamentStatus(context.TODO(), expectedTournamentID, expectedStatus)
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTournamentObjID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)

	expectedTournament := Tournament{
		ID:      expectedTournamentObjID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users:   []primitive.ObjectID{},
		Status:  expectedStatus,
	}

	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	badTournamentID := "bad_t_id"
	err = db.SetTournamentStatus(context.TODO(), badTournamentID, expectedStatus)
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badTournamentID))

	notExistTournamentID := primitive.NewObjectID().Hex()
	err = db.SetTournamentStatus(context.TODO(), notExistTournamentID, expectedStatus)
	require.EqualError(err, "update doc in collection: ModifiedCount != 1")

	cleanUp(t)
}

func TestJoinTournament(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	userJoinTorneyID := primitive.NewObjectID()
	err = db.JoinTournament(context.TODO(), expectedTournamentID, userJoinTorneyID.Hex())
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTournamentObjID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)

	expectedTournamentPrize := expectedTournamentDeposit
	expectedTournament := Tournament{
		ID:      expectedTournamentObjID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users:   []primitive.ObjectID{userJoinTorneyID},
		Prize:   expectedTournamentPrize,
	}
	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	badTournamentID := "bad_t_id"
	actualErr := db.JoinTournament(context.TODO(), badTournamentID, userJoinTorneyID.Hex())
	expectedErr := fmt.Sprintf("error processing transaction: AddUserToTournamentList: convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badTournamentID)
	require.EqualError(actualErr, expectedErr, "The two errors should be the same")

	badUserID := "bad_user_id"
	actualErr = db.JoinTournament(context.TODO(), expectedTournamentID, badUserID)
	expectedErr = fmt.Sprintf("error processing transaction: AddUserToTournamentList: convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badUserID)
	require.EqualError(actualErr, expectedErr, "The two errors should be the same")

	notExistTournamentID := primitive.NewObjectID().Hex()
	actualErr = db.JoinTournament(context.TODO(), notExistTournamentID, userJoinTorneyID.Hex())
	expectedErr = "error processing transaction: AddUserToTournamentList: update doc in collection: ModifiedCount != 1"
	require.EqualError(actualErr, expectedErr, "The two errors should be the same")

	cleanUp(t)
}

func TestFinishTournament(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	expectedUserName := "Vasya"
	expectedUserID, err := db.AddUser(context.TODO(), expectedUserName)
	require.NoError(err)

	err = db.JoinTournament(context.TODO(), expectedTournamentID, expectedUserID)
	require.NoError(err)

	err = db.FinishTournament(context.TODO(), expectedTournamentID, expectedUserID)
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTourneyID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)
	expectedUsrID, err := primitive.ObjectIDFromHex(expectedUserID)
	require.NoError(err)
	expectedTournament := Tournament{
		ID:      expectedTourneyID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users: []primitive.ObjectID{
			expectedUsrID,
		},
		Status: StatusFinished,
		Winner: expectedUsrID,
		Prize:  expectedTournamentDeposit,
	}
	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	actualUser, err := db.GetUser(context.TODO(), expectedUserID)
	require.NoError(err)

	expectedUser := User{
		ID:      expectedUsrID,
		Name:    expectedUserName,
		Balance: expectedTournamentDeposit,
	}
	require.Equal(expectedUser, *actualUser, "The two user objects should be the same")

	badTournamentID := "bad_t_id"
	actualErr := db.FinishTournament(context.TODO(), badTournamentID, expectedUserID)
	expectedErr := fmt.Sprintf("error processing transaction: SetTournamentStatus: convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badTournamentID)
	require.EqualError(actualErr, expectedErr, "The two errors should be the same")

	err = db.SetTournamentStatus(context.TODO(), expectedTournamentID, StatusStarted)
	require.NoError(err)

	badUserID := "bad_user_id"
	actualErr = db.FinishTournament(context.TODO(), expectedTournamentID, badUserID)
	expectedErr = fmt.Sprintf("error processing transaction: SetTournamentWinner: convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badUserID)
	require.EqualError(actualErr, expectedErr, "The two errors should be the same")

	notExistTournamentID := primitive.NewObjectID().Hex()
	actualErr = db.FinishTournament(context.TODO(), notExistTournamentID, expectedUserID)
	expectedErr = "error processing transaction: SetTournamentStatus: update doc in collection: ModifiedCount != 1"
	require.EqualError(actualErr, expectedErr, "The two errors should be the same")

	cleanUp(t)
}
