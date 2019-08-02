package storage

import (
	"context"

	"github.com/HarlamovBuldog/social-tournament-service/pkg/sts"
)

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
