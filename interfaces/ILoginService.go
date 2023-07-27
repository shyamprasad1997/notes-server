package interfaces

import (
	"context"
	"notes-server/models"
)

type ILoginService interface {
	Login(ctx context.Context, request models.LoginRequest) (models.LoginResponse, error)
	SignUp(ctx context.Context, request models.SignUpRequest) error
}
