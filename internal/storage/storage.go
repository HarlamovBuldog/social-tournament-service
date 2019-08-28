//go:generate mockgen -source=storage.go -destination=storage_mock.go -package=storage

package storage

import (
	"context"

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
}
