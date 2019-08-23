package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tournament represents a competition between players
// with deposit to enter and prize as a product of number
// of all players by deposit for winner.
type Tournament struct {
	ID      primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name    string               `json:"name" bson:"name"`
	Deposit float64              `json:"deposit" bson:"deposit"`
	Status  string               `json:"status" bson:"status"`
	Prize   float64              `json:"prize" bson:"prize"`
	Users   []primitive.ObjectID `json:"users" bson:"users"`
	Winner  primitive.ObjectID   `json:"winner" bson:"winner"`
}

// AddTournament func fills tournament info with provided name, provided deposit
// and with automatically generated id, then adds generated tournament info to database.
// It returns added tournamentID in string format if succeed and null string and err if smth wrong.
func (db *DB) AddTournament(ctx context.Context, name string, deposit float64) (string, error) {
	insertResult, err := db.conn.Collection(tournamentsCollectionName).InsertOne(ctx, Tournament{
		Name:    name,
		Deposit: deposit,
		Users:   []primitive.ObjectID{},
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

// GetTournament func tries to find tournament with provided id string.
// If succeed it returns *Tournament and nil error. If smth wrong it
// returns nil *Tournament and corresponding error.
func (db *DB) GetTournament(ctx context.Context, id string) (*Tournament, error) {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "convert string value to primitive.ObjectID type")
	}

	docReturned := db.conn.Collection(tournamentsCollectionName).FindOne(ctx, bson.M{"_id": primID})
	if err = docReturned.Err(); err != nil {
		return nil, errors.Wrap(err, "get doc from collection")
	}

	var tournament Tournament
	if err = docReturned.Decode(&tournament); err != nil {
		return nil, errors.Wrap(err, "decode returned doc")
	}

	return &tournament, nil
}

// DeleteTournament func tries to delete tournament with provided id string.
// If smth wrong it returns corresponding error, and nil error otherwise.
func (db *DB) DeleteTournament(ctx context.Context, id string) error {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(err, "convert string value to primitive.ObjectID type")
	}

	docDeleted, err := db.conn.Collection(tournamentsCollectionName).DeleteOne(ctx, bson.M{"_id": primID})
	if err != nil {
		return errors.Wrap(err, "delete doc from collection")
	}

	if docDeleted.DeletedCount != 1 {
		return errors.New("delete doc from collection: DeletedCount != 1")
	}

	return nil
}

// AddUserToTournamentList func ..
func (db *DB) AddUserToTournamentList(ctx context.Context, tournamentID, userID string) error {
	primTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return errors.Wrapf(err, "convert string %s to primitive.ObjectID type", tournamentID)
	}

	primUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.Wrapf(err, "convert string %s to primitive.ObjectID type", userID)
	}

	update := bson.D{
		{"$addToSet", bson.D{
			{"users", primUserID},
		}},
	}
	updateResult, err := db.conn.Collection(tournamentsCollectionName).UpdateOne(ctx,
		bson.M{"_id": primTournamentID}, update)
	if err != nil {
		return errors.Wrap(err, "update doc in collection")
	}

	if updateResult.ModifiedCount != 1 {
		return errors.New("update doc in collection: ModifiedCount != 1")
	}

	return nil
}

// SetTournamentWinner func ..
func (db *DB) SetTournamentWinner(ctx context.Context, tournamentID, userID string) error {
	primTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return errors.Wrapf(err, "convert string %s to primitive.ObjectID type", tournamentID)
	}

	primUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.Wrapf(err, "convert string %s to primitive.ObjectID type", userID)
	}

	update := bson.D{
		{"$set", bson.D{
			{"winner", primUserID},
		}},
	}
	updateResult, err := db.conn.Collection(tournamentsCollectionName).UpdateOne(ctx,
		bson.M{"_id": primTournamentID}, update)
	if err != nil {
		return errors.Wrap(err, "update doc in collection")
	}

	if updateResult.ModifiedCount != 1 {
		return errors.New("update doc in collection: ModifiedCount != 1")
	}

	return nil
}

// IncreaseTournamentPrize func ...
func (db *DB) IncreaseTournamentPrize(ctx context.Context, id string, amount float64) error {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrapf(err, "convert string %s to primitive.ObjectID type", id)
	}

	update := bson.D{
		{"$inc", bson.D{
			{"prize", amount},
		}},
	}
	updateResult, err := db.conn.Collection(tournamentsCollectionName).UpdateOne(ctx, bson.M{"_id": primID}, update)
	if err != nil {
		return errors.Wrap(err, "update doc in collection")
	}

	if updateResult.ModifiedCount != 1 {
		return errors.New("update doc in collection: ModifiedCount != 1")
	}

	return nil
}

// DecreaseTournamentPrize func ...
func (db *DB) DecreaseTournamentPrize(ctx context.Context, id string, amount float64) error {
	primID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrapf(err, "convert string %s to primitive.ObjectID type", id)
	}

	update := bson.D{
		{"$inc", bson.D{
			{"prize", -amount},
		}},
	}
	updateResult, err := db.conn.Collection(tournamentsCollectionName).UpdateOne(ctx, bson.M{"_id": primID}, update)
	if err != nil {
		return errors.Wrap(err, "update doc in collection")
	}

	if updateResult.ModifiedCount != 1 {
		return errors.New("update doc in collection: ModifiedCount != 1")
	}

	return nil
}

// SetTournamentStatus func ...
func (db *DB) SetTournamentStatus(ctx context.Context, tournamentID, status string) error {
	primTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return errors.Wrapf(err, "convert string %s to primitive.ObjectID type", tournamentID)
	}

	update := bson.D{
		{"$set", bson.D{
			{"status", status},
		}},
	}
	updateResult, err := db.conn.Collection(tournamentsCollectionName).UpdateOne(ctx,
		bson.M{"_id": primTournamentID}, update)
	if err != nil {
		return errors.Wrap(err, "update doc in collection")
	}

	if updateResult.ModifiedCount != 1 {
		return errors.New("update doc in collection: ModifiedCount != 1")
	}

	return nil
}
