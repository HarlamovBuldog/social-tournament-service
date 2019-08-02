package sts

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a player with id, name
// and certain amount of points as a balance
type User struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Balance float64            `json:"balance" bson:"balance"`
}

// Tournament represents a competition between players
// with deposit to enter and prize as a product of number
// of all players by deposit for winner.
type Tournament struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Deposit float64            `json:"deposit" bson:"deposit"`
	Status  string             `json:"status" bson:"status"`
	Prize   float64            `json:"prize" bson:"prize"`
	Users   []*User            `json:"users" bson:"users"`
	Winner  *User              `json:"winner" bson:"winner"`
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
	CalculateTournamentPrize(ctx context.Context, id string) error
	SetTournamentWinner(ctx context.Context, tournamentID, userID string) error
	AddUserToTournamentList(ctx context.Context, tournamentID, userID string) error
}
