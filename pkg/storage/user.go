package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/HarlamovBuldog/social-tournament-service/pkg/sts"
)

// User represents a player with id, name
// and certain amount of points as a balance
type User struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Balance float64            `json:"balance" bson:"balance"`
}
func (db *DB) AddUser(ctx context.Context, name string) (string, error) {
	insertResult, err := db.conn.Collection(usersCollectionName).InsertOne(ctx, sts.User{
		Name:    name,
		Balance: 0,
	})
	if err != nil {
		return "", fmt.Errorf("error insert doc to collection: %v", err)
	}
	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("error converting inserted id to primitive.ObjectID")
	}

	return insertedID.Hex(), nil
}

func (db *DB) GetUser(ctx context.Context, id string) (*sts.User, error) {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("error converting string value to primitive id: %v", err)
	}
	docReturned := db.conn.Collection(usersCollectionName).FindOne(ctx, bson.M{"_id": primID})
	if err := docReturned.Err(); err != nil {
		return nil, fmt.Errorf("error getting doc from collection: %v", err)
	}
	var user sts.User
	err = docReturned.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("error decoding returned doc: %v", err)
	}

	return &user, nil
}

func (db *DB) DeleteUser(ctx context.Context, id string) error {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error converting string value to primitive id: %v", err)
	}
	deleteResult, err := db.conn.Collection(usersCollectionName).DeleteOne(ctx, bson.M{"_id": primID})
	if err != nil {
		log.Printf("Error on deleting user: %v\n", err)
		return fmt.Errorf("Error on deleting user: %v", err)
	}
	if deleteResult.DeletedCount < 1 {
		log.Println("Error on deleting user: DeletedCount < 1")
		return errors.New("Error on deleting user: DeletedCount < 1")
	}
	return nil
}

func (db *DB) TakeUserBalance(ctx context.Context, id string, points float64) error {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error converting string value to primitive id: %v", err)
	}
	update := bson.D{
		{"$inc", bson.D{
			{"balance", -points},
		}},
	}
	updateResult, err := db.conn.Collection(usersCollectionName).UpdateOne(ctx, bson.M{"_id": primID}, update)
	if err != nil {
		log.Printf("Error on taking user balance: %v\n", err)
		return fmt.Errorf("Error on taking user balance: %v", err)
	}
	if updateResult.ModifiedCount < 1 {
		log.Println("Error on taking user balance: ModifiedCount < 1")
		return errors.New("Error on taking user balance: ModifiedCount < 1")
	}
	return nil
}

func (db *DB) FundUserBalance(ctx context.Context, id string, points float64) error {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error converting string value to primitive id: %v", err)
	}
	update := bson.D{
		{"$inc", bson.D{
			{"balance", points},
		}},
	}
	updateResult, err := db.conn.Collection(usersCollectionName).UpdateOne(ctx, bson.M{"_id": primID}, update)
	if err != nil {
		log.Printf("Error on funding user balance: %v\n", err)
		return fmt.Errorf("Error on funding user balance: %v", err)
	}
	if updateResult.ModifiedCount < 1 {
		log.Println("Error on funding user balance: ModifiedCount < 1")
		return errors.New("Error on funding user balance: ModifiedCount < 1")
	}
	return nil
}
