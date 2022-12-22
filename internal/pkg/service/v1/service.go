package v1

import (
	"context"

	v1 "github.com/HarlamovBuldog/social-tournament-service/internal/pkg/api/v1"
	"github.com/HarlamovBuldog/social-tournament-service/internal/pkg/storage"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// TournamentService is implementation of v1.Tournament proto interface.
type TournamentService struct {
	v1.UnimplementedTournamentServer
	db *storage.DB
}

// NewToDoServiceServer creates ToDo service
func NewToDoServiceServer(db *storage.DB) v1.TournamentServer {
	return &TournamentService{db: db}
}

func (t TournamentService) CreateUser(ctx context.Context, r *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	return &v1.CreateUserResponse{Id: "newID"}, nil
}

func (t TournamentService) UserList(ctx context.Context, r *v1.GetUserListRequest) (*v1.GetUserListResponse, error) {
	return &v1.GetUserListResponse{Users: getUserList()}, nil
}

func getUserList() []*v1.User {
	return []*v1.User{
		{
			Name: "Stas",
			Age:  12,
		},
		{
			Name: "Vlad",
			Age:  25,
		},
	}
}
