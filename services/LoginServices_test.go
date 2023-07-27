package services

import (
	"context"
	"errors"
	"notes-server/interfaces"
	"notes-server/loggers"
	"notes-server/models"
	"testing"

	"github.com/stretchr/testify/mock"
)

func Test_loginService_Login(t *testing.T) {
	type args struct {
		ctx     context.Context
		request models.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		given   func(*interfaces.MockILoginRepository)
		wantErr bool
	}{
		{
			name: "success case",
			args: args{
				ctx: context.Background(),
				request: models.LoginRequest{
					Email:    "test@gmail.com",
					Password: "testpassword",
				},
			},
			given: func(r *interfaces.MockILoginRepository) {
				r.EXPECT().Login(mock.Anything, mock.Anything).Return(models.LoginRepoResponse{
					Email: "test@gmail.com",
					Name:  "test",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "failure case - error in repo.Login()",
			args: args{
				ctx: context.Background(),
				request: models.LoginRequest{
					Email:    "test@gmail.com",
					Password: "testpassword",
				},
			},
			given: func(r *interfaces.MockILoginRepository) {
				r.EXPECT().Login(mock.Anything, mock.Anything).Return(models.LoginRepoResponse{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := interfaces.MockILoginRepository{}
			tt.given(&mockRepo)
			s := &loginService{
				repo:   &mockRepo,
				logger: loggers.NewLogger(),
			}
			_, err := s.Login(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("loginService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_loginService_SignUp(t *testing.T) {
	type args struct {
		ctx     context.Context
		request models.SignUpRequest
	}
	tests := []struct {
		name    string
		args    args
		given   func(*interfaces.MockILoginRepository)
		wantErr bool
	}{
		{
			name: "success case",
			args: args{
				ctx: context.Background(),
				request: models.SignUpRequest{
					Email:    "test@gmail.com",
					Password: "testpassword",
					Name:     "test",
				},
			},
			given: func(r *interfaces.MockILoginRepository) {
				r.EXPECT().SignUp(mock.Anything, mock.Anything).Return(nil)
				r.EXPECT().CheckIfUserExists(mock.Anything, mock.Anything).Return(false, nil)
			},
			wantErr: false,
		},
		{
			name: "failure case - error in repo.CheckIfUserExists",
			args: args{
				ctx: context.Background(),
				request: models.SignUpRequest{
					Email:    "test@gmail.com",
					Password: "testpassword",
					Name:     "test",
				},
			},
			given: func(r *interfaces.MockILoginRepository) {
				r.EXPECT().CheckIfUserExists(mock.Anything, mock.Anything).Return(false, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "failure case - user already exists",
			args: args{
				ctx: context.Background(),
				request: models.SignUpRequest{
					Email:    "test@gmail.com",
					Password: "testpassword",
					Name:     "test",
				},
			},
			given: func(r *interfaces.MockILoginRepository) {
				r.EXPECT().SignUp(mock.Anything, mock.Anything).Return(nil)
				r.EXPECT().CheckIfUserExists(mock.Anything, mock.Anything).Return(true, nil)
			},
			wantErr: true,
		},
		{
			name: "failure case - error in repo.SignUp()",
			args: args{
				ctx: context.Background(),
				request: models.SignUpRequest{
					Email:    "test@gmail.com",
					Password: "testpassword",
					Name:     "test",
				},
			},
			given: func(r *interfaces.MockILoginRepository) {
				r.EXPECT().CheckIfUserExists(mock.Anything, mock.Anything).Return(false, nil)
				r.EXPECT().SignUp(mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := interfaces.MockILoginRepository{}
			tt.given(&mockRepo)
			s := &loginService{
				repo:   &mockRepo,
				logger: loggers.NewLogger(),
			}
			err := s.SignUp(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("loginService.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
