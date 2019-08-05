package storage

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

// User represents a player with id, name
// and certain amount of points as a balance
type User struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Balance float64            `json:"balance" bson:"balance"`
}

// AddUser func fills user info with provided name, zero balance by default
// and with automatically generated id, then adds generated user info to database.
// It returns added userID in string format if succeed and null string and err if smth wrong.
func (db *DB) AddUser(ctx context.Context, name string) (string, error) {
	insertResult, err := db.conn.Collection(usersCollectionName).InsertOne(ctx, User{
		Name: name,
	})
	if err != nil {
		return "", errors.Wrap(err, "insert doc to collection")
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("convert inserted id to primitive.ObjectID")
	}

	return insertedID.Hex(), nil
}

// GetUser func tries to find user with provided id string.
// If succeed it returns *User and nil error. If smth wrong it
// returns nil *User and corresponding error.
func (db *DB) GetUser(ctx context.Context, id string) (*User, error) {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "convert string value to primitive.ObjectID type")
	}

	docReturned := db.conn.Collection(usersCollectionName).FindOne(ctx, bson.M{"_id": primID})
	if err = docReturned.Err(); err != nil {
		return nil, errors.Wrap(err, "get doc from collection")
	}

	var user User
	if err = docReturned.Decode(&user); err != nil {
		return nil, errors.Wrap(err, "decode returned doc")
	}

	return &user, nil
}

// DeleteUser func tries to delete user with provided id string.
// If smth wrong it returns corresponding error, and nil error otherwise.
func (db *DB) DeleteUser(ctx context.Context, id string) error {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "convert string value to primitive.ObjectID type")
	}

	deleteResult, err := db.conn.Collection(usersCollectionName).DeleteOne(ctx, bson.M{"_id": primID})
	if err != nil {
		return errors.Wrap(err, "delete doc from collection")
	}

	if deleteResult.DeletedCount != 1 {
		return errors.New("delete doc from collection: DeletedCount != 1")
	}

	return nil
}

// TakeUserBalance func tries to decrease user balance with provided id string.
// If smth wrong it returns corresponding error, and nil error otherwise.
func (db *DB) TakeUserBalance(ctx context.Context, id string, points float64) error {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "convert string value to primitive.ObjectID type")
	}

	update := bson.D{
		{"$inc", bson.D{
			{"balance", -points},
		}},
	}
	updateResult, err := db.conn.Collection(usersCollectionName).UpdateOne(ctx, bson.M{"_id": primID}, update)
	if err != nil {
		return errors.Wrap(err, "update doc from collection")
	}

	if updateResult.ModifiedCount != 1 {
		return errors.New("update doc from collection: ModifiedCount != 1")
	}

	return nil
}

// FundUserBalance func tries to increase user balance with provided id string.
// If smth wrong it returns corresponding error, and nil error otherwise
func (db *DB) FundUserBalance(ctx context.Context, id string, points float64) error {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "convert string value to primitive.ObjectID type")
	}

	update := bson.D{
		{"$inc", bson.D{
			{"balance", points},
		}},
	}
	updateResult, err := db.conn.Collection(usersCollectionName).UpdateOne(ctx, bson.M{"_id": primID}, update)
	if err != nil {
		return errors.Wrap(err, "update doc from collection")
	}

	if updateResult.ModifiedCount != 1 {
		return errors.New("update doc from collection: ModifiedCount != 1")
	}

	return nil
}
