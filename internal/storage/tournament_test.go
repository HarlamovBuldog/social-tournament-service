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
		Users:   &[]string{},
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
		Users:   &[]string{},
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

	userID := primitive.NewObjectID().Hex()
	err = db.AddUserToTournamentList(context.TODO(), expectedTournamentID, userID)
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTournamentObjID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)

	expectedTournament := Tournament{
		ID:      expectedTournamentObjID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users:   &[]string{userID},
	}

	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	badTournamentID := "bad_t_id"
	err = db.AddUserToTournamentList(context.TODO(), badTournamentID, userID)
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badTournamentID))

	badUserID := "bad_user_id"
	err = db.AddUserToTournamentList(context.TODO(), expectedTournamentID, badUserID)
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badUserID))

	notExistTournamentID := primitive.NewObjectID().Hex()
	err = db.AddUserToTournamentList(context.TODO(), notExistTournamentID, userID)
	require.EqualError(err, "update doc in collection: ModifiedCount != 1")

	cleanUp(t)
}

func TestSetTournamentWinner(t *testing.T) {
	expectedTournamentName := "tournament-1"
	expectedTournamentDeposit := 1000.0
	expectedTournamentID, err := db.AddTournament(context.TODO(), expectedTournamentName, expectedTournamentDeposit)
	require := require.New(t)
	require.NoError(err)

	userWinnerID := primitive.NewObjectID().Hex()
	err = db.SetTournamentWinner(context.TODO(), expectedTournamentID, userWinnerID)
	require.NoError(err)

	actualTournament, err := db.GetTournament(context.TODO(), expectedTournamentID)
	require.NoError(err)

	expectedTournamentObjID, err := primitive.ObjectIDFromHex(expectedTournamentID)
	require.NoError(err)

	expectedTournament := Tournament{
		ID:      expectedTournamentObjID,
		Name:    expectedTournamentName,
		Deposit: expectedTournamentDeposit,
		Users:   &[]string{},
		Winner:  userWinnerID,
	}

	require.Equal(expectedTournament, *actualTournament, "The two tournament objects should be the same")

	badTournamentID := "bad_t_id"
	err = db.SetTournamentWinner(context.TODO(), badTournamentID, userWinnerID)
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badTournamentID))

	badUserID := "bad_user_id"
	err = db.SetTournamentWinner(context.TODO(), expectedTournamentID, badUserID)
	require.EqualError(err, fmt.Sprintf("convert string %s to primitive.ObjectID type: encoding/hex: invalid byte: U+005F '_'", badUserID))

	notExistTournamentID := primitive.NewObjectID().Hex()
	err = db.SetTournamentWinner(context.TODO(), notExistTournamentID, userWinnerID)
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
		Users:   &[]string{},
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
		Users:   &[]string{},
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
