package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
func (db *DB) AddTournament(ctx context.Context, name string, deposit float64) (string, error) {
	return "", nil
}
func (db *DB) GetTournament(ctx context.Context, id string) (*sts.Tournament, error) {
	return &sts.Tournament{}, nil
}
func (db *DB) DeleteTournament(ctx context.Context, id string) error {
	return nil
}
func (db *DB) CalculateTournamentPrize(ctx context.Context, id string) error {
	return nil
}
func (db *DB) SetTournamentWinner(ctx context.Context, tournamentID, userID string) error {
	return nil
}
func (db *DB) AddUserToTournamentList(ctx context.Context, tournamentID, userID string) error {
	return nil
}
