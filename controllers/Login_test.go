package controllers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"notes-server/interfaces"
	"notes-server/loggers"
	"notes-server/models"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestLoginController_Login(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name  string
		given func(*interfaces.MockILoginService)
		args  args
		want  int
	}{
		{
			name: "success case",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"test@gmail.com", "password":"testpassword"}`),
			},
			given: func(s *interfaces.MockILoginService) {
				s.EXPECT().Login(mock.Anything, mock.Anything).Return(models.LoginResponse{
					SID: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGFjY3Vrbm94LmNvbSIsIm5hbWUiOiJBZG1pbiIsImV4cCI6MTY5MDQ2MDY2MX0.i-SrhNfMcVjHH7rLlDkczkt5TyQ6Oa1McmUk5cnKAKQ",
				}, nil)
			},
			want: http.StatusOK,
		},
		{
			name: "failure case - invalid request",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"test@gmail.com", password":"testpassword"}`),
			},
			given: func(s *interfaces.MockILoginService) {
			},
			want: http.StatusBadRequest,
		},
		{
			name: "failure case - missing email",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"", "password":"testpassword"}`),
			},
			given: func(s *interfaces.MockILoginService) {
			},
			want: http.StatusBadRequest,
		},
		{
			name: "failure case - missing password",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"test@gmail.com", "password":""}`),
			},
			given: func(s *interfaces.MockILoginService) {
			},
			want: http.StatusBadRequest,
		},
		{
			name: "failure case - error in service.Login()",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"test@gmail.com", "password":"testpassword"}`),
			},
			given: func(s *interfaces.MockILoginService) {
				s.EXPECT().Login(mock.Anything, mock.Anything).Return(models.LoginResponse{}, errors.New("db errpr"))
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := interfaces.MockILoginService{}
			tt.given(&mockService)
			c := &LoginController{
				service: &mockService,
				logger:  loggers.NewLogger(),
			}
			c.Login(tt.args.w, tt.args.r)
			if tt.args.w.Result().StatusCode != tt.want {
				t.Errorf("expected status code %d, got %d", tt.want, tt.args.w.Result().StatusCode)
			}
		})
	}
}

func CreateReq(req string) *http.Request {
	return &http.Request{
		Method: http.MethodPost,
		Body:   io.NopCloser(bytes.NewBuffer([]byte(req))),
	}
}

func TestLoginController_SignUp(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name  string
		given func(*interfaces.MockILoginService)
		args  args
		want  int
	}{
		{
			name: "success case",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"test@gmail.com", "password":"testpassword", "name":"test"}`),
			},
			given: func(s *interfaces.MockILoginService) {
				s.EXPECT().SignUp(mock.Anything, mock.Anything).Return(nil)
			},
			want: http.StatusCreated,
		},
		{
			name: "failure case - invalid request",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"test@gmail.com", password":"testpassword", "name":"test"}`),
			},
			given: func(s *interfaces.MockILoginService) {
				s.EXPECT().SignUp(mock.Anything, mock.Anything).Return(nil)
			},
			want: http.StatusBadRequest,
		},
		{
			name: "failure case - missing email",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"", "password":"testpassword", "name":"test"}`),
			},
			given: func(s *interfaces.MockILoginService) {
				s.EXPECT().SignUp(mock.Anything, mock.Anything).Return(nil)
			},
			want: http.StatusBadRequest,
		},
		{
			name: "failure case - missing password",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"test@gmail.com", "password":"", "name":"test"}`),
			},
			given: func(s *interfaces.MockILoginService) {
				s.EXPECT().SignUp(mock.Anything, mock.Anything).Return(nil)
			},
			want: http.StatusBadRequest,
		},
		{
			name: "failure case - missing name",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"test@gmail.com", "password":"testpassword", "name":""}`),
			},
			given: func(s *interfaces.MockILoginService) {
				s.EXPECT().SignUp(mock.Anything, mock.Anything).Return(nil)
			},
			want: http.StatusBadRequest,
		},
		{
			name: "failure case - error in service.SignUp()",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"email":"test@gmail.com", "password":"testpassword", "name":"test"}`),
			},
			given: func(s *interfaces.MockILoginService) {
				s.EXPECT().SignUp(mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := interfaces.MockILoginService{}
			tt.given(&mockService)
			c := &LoginController{
				service: &mockService,
				logger:  loggers.NewLogger(),
			}
			c.SignUp(tt.args.w, tt.args.r)
			if tt.args.w.Result().StatusCode != tt.want {
				t.Errorf("expected status code %d, got %d", tt.want, tt.args.w.Result().StatusCode)
			}
		})
	}
}
