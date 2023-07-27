package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"notes-server/interfaces"
	"notes-server/loggers"
	"notes-server/models"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestNotesController_GetNotes(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name  string
		given func(*interfaces.MockINotesService)
		args  args
		want  int
	}{
		{
			name: "success case",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(``),
			},
			given: func(s *interfaces.MockINotesService) {
				s.EXPECT().GetNotes(mock.Anything).Return([]models.Note{{
					Id:        1,
					Note:      "test note",
					CreatedBy: "test@gmail.com",
				}}, nil)
			},
			want: http.StatusOK,
		},
		{
			name: "failure case - error in service.GetNotes()",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(``),
			},
			given: func(s *interfaces.MockINotesService) {
				s.EXPECT().GetNotes(mock.Anything).Return([]models.Note{}, errors.New("db error"))
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := interfaces.MockINotesService{}
			tt.given(&mockService)
			c := &NotesController{
				service: &mockService,
				logger:  loggers.NewLogger(),
			}
			c.GetNotes(tt.args.w, tt.args.r)
			if tt.args.w.Result().StatusCode != tt.want {
				t.Errorf("expected status code %d, got %d", tt.want, tt.args.w.Result().StatusCode)
			}
		})
	}
}

func TestNotesController_AddNote(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name  string
		given func(*interfaces.MockINotesService)
		args  args
		want  int
	}{
		{
			name: "success case",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"note":"test note"}`),
			},
			given: func(s *interfaces.MockINotesService) {
				s.EXPECT().AddNote(mock.Anything, mock.Anything).Return(models.AddNoteResponse{
					Id: 123,
				}, nil)
			},
			want: http.StatusCreated,
		},
		{
			name: "failure case - invalid request",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"note": test note"}`),
			},
			given: func(s *interfaces.MockINotesService) {
			},
			want: http.StatusBadRequest,
		},
		{
			name: "failure case - note missing",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"note": ""}`),
			},
			given: func(s *interfaces.MockINotesService) {
			},
			want: http.StatusBadRequest,
		},
		{
			name: "failure case - error in service.AddNote()",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"note":"test note"}`),
			},
			given: func(s *interfaces.MockINotesService) {
				s.EXPECT().AddNote(mock.Anything, mock.Anything).Return(models.AddNoteResponse{
					Id: 123,
				}, errors.New("db error"))
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := interfaces.MockINotesService{}
			tt.given(&mockService)
			c := &NotesController{
				service: &mockService,
				logger:  loggers.NewLogger(),
			}
			c.AddNote(tt.args.w, tt.args.r)
			if tt.args.w.Result().StatusCode != tt.want {
				t.Errorf("expected status code %d, got %d", tt.want, tt.args.w.Result().StatusCode)
			}
		})
	}
}

func TestNotesController_DeleteNote(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name  string
		given func(*interfaces.MockINotesService)
		args  args
		want  int
	}{
		{
			name: "success case",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"id":123}`),
			},
			given: func(s *interfaces.MockINotesService) {
				s.EXPECT().DeleteNote(mock.Anything, mock.Anything).Return(nil)
			},
			want: http.StatusOK,
		},
		{
			name: "failure case - invalid request",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"id":"123"}`),
			},
			given: func(s *interfaces.MockINotesService) {
			},
			want: http.StatusBadRequest,
		},
		{
			name: "failure case - error in service.DeleteNote()",
			args: args{
				w: httptest.NewRecorder(),
				r: CreateReq(`{"id":123}`),
			},
			given: func(s *interfaces.MockINotesService) {
				s.EXPECT().DeleteNote(mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := interfaces.MockINotesService{}
			tt.given(&mockService)
			c := &NotesController{
				service: &mockService,
				logger:  loggers.NewLogger(),
			}
			c.DeleteNote(tt.args.w, tt.args.r)
			if tt.args.w.Result().StatusCode != tt.want {
				t.Errorf("expected status code %d, got %d", tt.want, tt.args.w.Result().StatusCode)
			}
		})
	}
}
