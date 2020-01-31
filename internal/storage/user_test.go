package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestAddUser(t *testing.T) {
	userIDExpected, err := db.AddUser(context.TODO(), "gennadiy")
	require.NoError(t, err, "AddUser func should return nil error")

	actualUser, err := db.GetUser(context.TODO(), userIDExpected)
	require.NoError(t, err, "GetUser func should return nil error")

	assert := assert.New(t)
	assert.Equal(userIDExpected, actualUser.ID.Hex(), "The two IDs should be the same.")

	userIDExpected2, err := primitive.ObjectIDFromHex(userIDExpected)
	require.NoError(t, err, "ObjectIDFromHex func should return nil error")

	assert.Equal(&User{ID: userIDExpected2, Name: "gennadiy"}, actualUser, "The two users should be the same.")

	cleanUp(t)
}

func TestGetUser(t *testing.T) {
	userIDExpected, err := db.AddUser(context.TODO(), "Vasya")
	require.NoError(t, err, "AddUser func should return nil error")

	actualUser, err := db.GetUser(context.TODO(), userIDExpected)
	require.NoError(t, err, "GetUser func should return nil error")

	assert := assert.New(t)
	assert.Equal(userIDExpected, actualUser.ID.Hex(), "The two IDs should be the same.")

	userIDExpected2, err := primitive.ObjectIDFromHex(userIDExpected)
	require.NoError(t, err, "ObjectIDFromHex func should return nil error")

	assert.Equal(&User{ID: userIDExpected2, Name: "Vasya"}, actualUser, "The two users should be the same.")

	badUserID := "safasf2412"
	_, err = db.GetUser(context.TODO(), badUserID)
	assert.EqualError(err,
		"convert string value to primitive.ObjectID type: encoding/hex: invalid byte: U+0073 's'",
		"The error should contain text")

	notExistUserID := primitive.NewObjectID()
	actualUser, err = db.GetUser(context.TODO(), notExistUserID.Hex())
	assert.EqualError(err, "decode returned doc: "+mongo.ErrNoDocuments.Error(), "The two errors should be the same")

	cleanUp(t)
}

func TestDeleteUser(t *testing.T) {
	userIDExpected, err := db.AddUser(context.TODO(), "Vasya")
	require.NoError(t, err, "AddUser func should return nil error")

	err = db.DeleteUser(context.TODO(), userIDExpected)
	require.NoError(t, err, "DeleteUser func should return nil error")

	badUserID := "safasf2412"
	err = db.DeleteUser(context.TODO(), badUserID)
	assert := assert.New(t)
	assert.EqualError(err,
		"convert string value to primitive.ObjectID type: encoding/hex: invalid byte: U+0073 's'",
		"The error should contain text")

	err = db.DeleteUser(context.TODO(), userIDExpected)
	assert.EqualError(err, "delete doc from collection: DeletedCount != 1", "The two errors should be the same")

	cleanUp(t)
}

func TestTakeUserBalance(t *testing.T) {
	generatedUserID := primitive.NewObjectID()
	amount := 100.0

	err := db.TakeUserBalance(context.TODO(), generatedUserID.Hex(), amount)
	assert := assert.New(t)
	assert.EqualError(err,
		"update doc in collection: ModifiedCount != 1",
		"The error should contain text")

	badUserID := "safasf2412"
	err = db.TakeUserBalance(context.TODO(), badUserID, amount)
	assert.EqualError(err,
		"convert string value to primitive.ObjectID type: encoding/hex: invalid byte: U+0073 's'",
		"The error should contain text")

	addedUserID, err := db.AddUser(context.TODO(), "Vasya")
	require.NoError(t, err, "AddUser func should return nil error")

	err = db.TakeUserBalance(context.TODO(), addedUserID, amount)
	require.NoError(t, err, "TakeUserBalance func should return nil error")

	addedUser, err := db.GetUser(context.TODO(), addedUserID)
	require.NoError(t, err, "GetUser func should return nil error")

	addedUserObjectID, err := primitive.ObjectIDFromHex(addedUserID)
	require.NoError(t, err, "ObjectIDFromHex func should return nil error")

	assert.Equal(&User{ID: addedUserObjectID, Name: "Vasya", Balance: -100.0}, addedUser,
		"The two users should be the same.")

	cleanUp(t)
}

func TestFundUserBalance(t *testing.T) {
	generatedUserID := primitive.NewObjectID()
	amount := 100.0

	err := db.FundUserBalance(context.TODO(), generatedUserID.Hex(), amount)
	assert := assert.New(t)
	assert.EqualError(err,
		"update doc in collection: ModifiedCount != 1",
		"The error should contain text")

	badUserID := "safasf2412"
	err = db.FundUserBalance(context.TODO(), badUserID, amount)
	assert.EqualError(err,
		"convert string value to primitive.ObjectID type: encoding/hex: invalid byte: U+0073 's'",
		"The error should contain text")

	addedUserID, err := db.AddUser(context.TODO(), "Vasya")
	require.NoError(t, err, "AddUser func should return nil error")

	err = db.FundUserBalance(context.TODO(), addedUserID, amount)
	require.NoError(t, err, "FundUserBalance func should return nil error")

	addedUser, err := db.GetUser(context.TODO(), addedUserID)
	require.NoError(t, err, "GetUser func should return nil error")

	addedUserObjectID, err := primitive.ObjectIDFromHex(addedUserID)
	require.NoError(t, err, "ObjectIDFromHex func should return nil error")

	assert.Equal(&User{ID: addedUserObjectID, Name: "Vasya", Balance: 100.0}, addedUser,
		"The two users should be the same.")

	cleanUp(t)
}
