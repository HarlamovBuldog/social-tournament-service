//go:generate mockgen -source=storage.go -destination=storage_mock.go -package=storage

package storage

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DB is struct that holds database object
type DB struct {
	conn *mongo.Database
}

const (
	usersCollectionName       = "users"
	tournamentsCollectionName = "tournaments"
)

// CreateNew is constructor for db
func CreateNew(db *mongo.Database) *DB {
	return &DB{
		conn: db,
	}
}

func (db *DB) JoinTournament(ctx context.Context, tournamentID, userID string) error {
	session, err := db.conn.Client().StartSession()
	if err != nil {
		return errors.Wrap(err, "error start mongoDB session")
	}
	if err = session.StartTransaction(); err != nil {
		return errors.Wrap(err, "error start transaction")
	}
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := db.AddUserToTournamentList(sc, tournamentID, userID); err != nil {
			return errors.Wrap(err, "AddUserToTournamentList")
		}

		tournament, err := db.GetTournament(sc, tournamentID)
		if err != nil {
			return errors.Wrap(err, "GetTournament")
		}

		if err := db.IncreaseTournamentPrize(sc, tournamentID, tournament.Deposit); err != nil {
			return errors.Wrap(err, "IncreaseTournamentPrize")
		}

		if err = session.CommitTransaction(sc); err != nil {
			return errors.Wrap(err, "commit transaction")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "error processing transaction")
	}

	session.EndSession(ctx)
	return nil
}

func (db *DB) FinishTournament(ctx context.Context, tournamentID string) error {
	session, err := db.conn.Client().StartSession()
	if err != nil {
		return errors.Wrap(err, "error start mongoDB session")
	}
	if err = session.StartTransaction(); err != nil {
		return errors.Wrap(err, "error start transaction")
	}
	status := "finished"
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := db.SetTournamentStatus(sc, tournamentID, status); err != nil {
			return errors.Wrap(err, "SetTournamentStatus")
		}
		userID := primitive.NewObjectID()
		err := db.SetTournamentWinner(sc, tournamentID, userID.Hex())
		if err != nil {
			return errors.Wrap(err, "SetTournamentWinner")
		}

		tournament, err := db.GetTournament(sc, tournamentID)
		if err != nil {
			return errors.Wrap(err, "GetTournament")
		}

		err = db.FundUserBalance(sc, userID.Hex(), tournament.Prize)
		if err != nil {
			return errors.Wrap(err, "FundUserBalance")
		}

		if err = session.CommitTransaction(sc); err != nil {
			return errors.Wrap(err, "commit transaction")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "error processing transaction")
	}

	session.EndSession(ctx)
	return nil
}

// Service is the wrapper for all methods working with db.
type Service interface {
	// AddUser adds user to db with given name, auto-increments id
	// and sets user balance to zero. Returns userID if succeed.
	AddUser(ctx context.Context, name string) (string, error)

	// GetUser returns *User that contains all information about user with help of provided id
	GetUser(ctx context.Context, id string) (*User, error)

	DeleteUser(ctx context.Context, id string) error

	// TakeUserBalance finds user with provided id and deducts from his balance provided points
	TakeUserBalance(ctx context.Context, id string, points float64) error

	// FundUserBalance finds user with provided id and adds to his balance provided points
	FundUserBalance(ctx context.Context, id string, points float64) error

	AddTournament(ctx context.Context, name string, deposit float64) (string, error)
	GetTournament(ctx context.Context, id string) (*Tournament, error)
	DeleteTournament(ctx context.Context, id string) error
	IncreaseTournamentPrize(ctx context.Context, id string, amount float64) error
	DecreaseTournamentPrize(ctx context.Context, id string, amount float64) error
	SetTournamentWinner(ctx context.Context, tournamentID, userID string) error
	SetTournamentStatus(ctx context.Context, tournamentID, status string) error
	AddUserToTournamentList(ctx context.Context, tournamentID, userID string) error

	JoinTournament(ctx context.Context, tournamentID, userID string) error
	FinishTournament(ctx context.Context, tournamentID string) error
}
