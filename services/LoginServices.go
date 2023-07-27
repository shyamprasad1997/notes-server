package services

import (
	"context"
	"errors"
	"notes-server/constants"
	"notes-server/interfaces"
	"notes-server/loggers"
	"notes-server/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type loginService struct {
	repo   interfaces.ILoginRepository
	logger *loggers.Logger
}

func NewLoginService(logger *loggers.Logger, repo interfaces.ILoginRepository) interfaces.ILoginService {
	return &loginService{
		repo:   repo,
		logger: logger,
	}
}

// Login - service layer for POST /login route, validates creds and returns JWT token
func (s *loginService) Login(ctx context.Context, request models.LoginRequest) (models.LoginResponse, error) {
	s.logger.Info(ctx, "Entering LoginService.Login()")
	defer s.logger.Info(ctx, "Entering LoginService.Login()")
	response, err := s.repo.Login(ctx, request)
	if err != nil {
		s.logger.Warn(ctx, "Error in LoginService.Login(), error from s.repo.Login()")
		return models.LoginResponse{}, err
	}
	token, err := generateJWTToken(response.Email, response.Name)
	if err != nil {
		s.logger.Warn(ctx, "Error in LoginService.Login(), error from generateJWTToken()")
		return models.LoginResponse{}, err
	}
	return models.LoginResponse{
		SID: token,
	}, nil
}

// generateJWTToken - Creates a JWT token with HS256 signing method and signed with the secret
func generateJWTToken(email string, name string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		Email: email,
		Name:  name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString(constants.JwtSecretEnvKey)))
}

func (s *loginService) SignUp(ctx context.Context, request models.SignUpRequest) error {
	s.logger.Info(ctx, "Entering LoginService.SignUp()")
	defer s.logger.Info(ctx, "Entering LoginService.SignUp()")
	isUserExists, err := s.repo.CheckIfUserExists(ctx, request.Email)
	if err != nil {
		s.logger.Warn(ctx, "Error in LoginService.SignUp(), error from s.repo.CheckIfUserExists()")
		return err
	}
	if isUserExists {
		s.logger.Warn(ctx, "Error in LoginService.SignUp(), user already exists")
		return errors.New("user already exists")
	}
	err = s.repo.SignUp(ctx, request)
	if err != nil {
		s.logger.Warn(ctx, "Error in LoginService.SignUp(), error from s.repo.SignUp()")
		return err
	}
	return nil
}
