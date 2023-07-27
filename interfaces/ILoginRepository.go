package interfaces

import (
	"context"
	"notes-server/models"
)

type ILoginRepository interface {
	Login(ctx context.Context, request models.LoginRequest) (models.LoginRepoResponse, error)
	CheckIfUserExists(ctx context.Context, email string) (bool, error)
	SignUp(ctx context.Context, request models.SignUpRequest) error
	ValidateUser(ctx context.Context, email, name string) error
}
