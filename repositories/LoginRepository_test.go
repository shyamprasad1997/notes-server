package repositories

import (
	"context"
	"errors"
	"notes-server/db"
	"notes-server/loggers"
	"notes-server/models"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func Test_loginRepository_Login(t *testing.T) {
	type args struct {
		ctx     context.Context
		request models.LoginRequest
	}
	tests := []struct {
		name    string
		given   func(*db.MockDB)
		args    args
		want    models.LoginRepoResponse
		wantErr bool
	}{
		{
			name: "success case",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(&models.User{
					Id:       123,
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password",
				}, nil)
				mockTxn.EXPECT().Commit()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx: context.Background(),
				request: models.LoginRequest{
					Email:    "test@gmail.com",
					Password: "password",
				},
			},
			want: models.LoginRepoResponse{
				Email: "test@gmail.com",
				Name:  "test",
			},
			wantErr: false,
		},
		{
			name: "failure case - error in txn.First()",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(&models.User{}, errors.New("db error"))
				mockTxn.EXPECT().Abort()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx: context.Background(),
				request: models.LoginRequest{
					Email:    "test@gmail.com",
					Password: "password",
				},
			},
			want:    models.LoginRepoResponse{},
			wantErr: true,
		},
		{
			name: "failure case - data dosent exist in db",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				mockTxn.EXPECT().Abort()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx: context.Background(),
				request: models.LoginRequest{
					Email:    "test@gmail.com",
					Password: "password",
				},
			},
			want:    models.LoginRepoResponse{},
			wantErr: true,
		},
		{
			name: "failure case - passwords do not match",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(&models.User{
					Id:       123,
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "passwordold",
				}, nil)
				mockTxn.EXPECT().Abort()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx: context.Background(),
				request: models.LoginRequest{
					Email:    "test@gmail.com",
					Password: "password",
				},
			},
			want:    models.LoginRepoResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := db.MockDB{}
			tt.given(&mockDB)
			r := &loginRepository{
				db:     &mockDB,
				logger: loggers.NewLogger(),
			}
			got, err := r.Login(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("loginRepository.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loginRepository.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loginRepository_SignUp(t *testing.T) {
	type args struct {
		ctx     context.Context
		request models.SignUpRequest
	}
	tests := []struct {
		name    string
		given   func(*db.MockDB)
		args    args
		wantErr bool
	}{
		{
			name: "success case",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().Insert(mock.Anything, mock.Anything).Return(nil)
				mockTxn.EXPECT().Commit()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx: context.Background(),
				request: models.SignUpRequest{
					Email:    "test@gmail.com",
					Password: "password",
					Name:     "test",
				},
			},
			wantErr: false,
		},
		{
			name: "failure case - error in txn.Insert()",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().Insert(mock.Anything, mock.Anything).Return(errors.New("db error"))
				mockTxn.EXPECT().Abort()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx: context.Background(),
				request: models.SignUpRequest{
					Email:    "test@gmail.com",
					Password: "password",
					Name:     "test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := db.MockDB{}
			tt.given(&mockDB)
			r := &loginRepository{
				db:     &mockDB,
				logger: loggers.NewLogger(),
			}
			err := r.SignUp(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("loginRepository.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_loginRepository_CheckIfUserExists(t *testing.T) {
	type args struct {
		ctx     context.Context
		request string
	}
	tests := []struct {
		name    string
		given   func(*db.MockDB)
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success case",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(&models.User{
					Id:       123,
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password",
				}, nil)
				mockTxn.EXPECT().Commit()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx:     context.Background(),
				request: "test@gmail.com",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "failure case - error in txn.First()",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(&models.User{}, errors.New("db error"))
				mockTxn.EXPECT().Abort()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx:     context.Background(),
				request: "test@gmail.com",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "failure case - data dosent exist in db",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				mockTxn.EXPECT().Commit()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx:     context.Background(),
				request: "test@gmail.com",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := db.MockDB{}
			tt.given(&mockDB)
			r := &loginRepository{
				db:     &mockDB,
				logger: loggers.NewLogger(),
			}
			got, err := r.CheckIfUserExists(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("loginRepository.CheckIfUserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loginRepository.CheckIfUserExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loginRepository_ValidateUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		request struct {
			email, name string
		}
	}
	tests := []struct {
		name    string
		given   func(*db.MockDB)
		args    args
		wantErr bool
	}{
		{
			name: "success case",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(&models.User{
					Id:       123,
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password",
				}, nil)
				mockTxn.EXPECT().Commit()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx: context.Background(),
				request: struct{ email, name string }{
					email: "test@gmail.com",
					name:  "test",
				},
			},
			wantErr: false,
		},
		{
			name: "failure case - error in txn.First()",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(&models.User{}, errors.New("db error"))
				mockTxn.EXPECT().Abort()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx: context.Background(),
				request: struct{ email, name string }{
					email: "test@gmail.com",
					name:  "test",
				},
			},
			wantErr: true,
		},
		{
			name: "failure case - data dosent exist in db",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
				mockTxn.EXPECT().Commit()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx: context.Background(),
				request: struct{ email, name string }{
					email: "test@gmail.com",
					name:  "test",
				},
			},
			wantErr: true,
		},
		{
			name: "failure case - invalid user",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().First(mock.Anything, mock.Anything, mock.Anything).Return(&models.User{
					Id:       123,
					Name:     "test2",
					Email:    "test@gmail.com",
					Password: "password",
				}, nil)
				mockTxn.EXPECT().Commit()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx: context.Background(),
				request: struct{ email, name string }{
					email: "test@gmail.com",
					name:  "test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := db.MockDB{}
			tt.given(&mockDB)
			r := &loginRepository{
				db:     &mockDB,
				logger: loggers.NewLogger(),
			}
			err := r.ValidateUser(tt.args.ctx, tt.args.request.email, tt.args.request.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("loginRepository.ValidateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
