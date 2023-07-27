package services

import (
	"context"
	"errors"
	"notes-server/interfaces"
	"notes-server/loggers"
	"notes-server/models"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func Test_notesService_GetNotes(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		given   func(*interfaces.MockINotesRepository)
		args    args
		want    []models.Note
		wantErr bool
	}{
		{
			name: "success case",
			given: func(r *interfaces.MockINotesRepository) {
				r.EXPECT().GetNotes(mock.Anything, mock.Anything).Return([]models.Note{{
					Id:        1,
					Note:      "test",
					CreatedBy: "test@gmail.com",
				}}, nil)
			},
			args: args{
				ctx: context.Background(),
			},
			want: []models.Note{{
				Id:        1,
				Note:      "test",
				CreatedBy: "test@gmail.com",
			}},
			wantErr: false,
		},
		{
			name: "failure case - error in repo.GetNotes()",
			given: func(r *interfaces.MockINotesRepository) {
				r.EXPECT().GetNotes(mock.Anything, mock.Anything).Return([]models.Note{}, errors.New("db error"))
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []models.Note{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := interfaces.MockINotesRepository{}
			tt.given(&mockRepo)
			s := &notesService{
				repo:   &mockRepo,
				logger: loggers.NewLogger(),
			}
			got, err := s.GetNotes(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("notesService.GetNotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("notesService.GetNotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notesService_AddNote(t *testing.T) {
	type args struct {
		ctx     context.Context
		request models.AddNoteRequest
	}
	tests := []struct {
		name    string
		given   func(*interfaces.MockINotesRepository)
		args    args
		want    models.AddNoteResponse
		wantErr bool
	}{
		{
			name: "success case",
			given: func(r *interfaces.MockINotesRepository) {
				r.EXPECT().AddNote(mock.Anything, mock.Anything).Return(123, nil)
			},
			args: args{
				ctx: context.Background(),
				request: models.AddNoteRequest{
					Email: "test@gmail.com",
					Note:  "test note",
				},
			},
			want: models.AddNoteResponse{
				Id: 123,
			},
			wantErr: false,
		},
		{
			name: "failure case - error in repo.AddNote()",
			given: func(r *interfaces.MockINotesRepository) {
				r.EXPECT().AddNote(mock.Anything, mock.Anything).Return(0, errors.New("db error"))
			},
			args: args{
				ctx: context.Background(),
				request: models.AddNoteRequest{
					Email: "test@gmail.com",
					Note:  "test note",
				},
			},
			want:    models.AddNoteResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := interfaces.MockINotesRepository{}
			tt.given(&mockRepo)
			s := &notesService{
				repo:   &mockRepo,
				logger: loggers.NewLogger(),
			}
			got, err := s.AddNote(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("notesService.AddNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("notesService.AddNote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notesService_DeleteNote(t *testing.T) {
	type args struct {
		ctx     context.Context
		request models.DeleteNoteRequest
	}
	tests := []struct {
		name    string
		given   func(*interfaces.MockINotesRepository)
		args    args
		wantErr bool
	}{
		{
			name: "success case",
			given: func(r *interfaces.MockINotesRepository) {
				r.EXPECT().DeleteNote(mock.Anything, mock.Anything).Return(nil)
			},
			args: args{
				ctx: context.Background(),
				request: models.DeleteNoteRequest{
					Id: 123,
				},
			},
			wantErr: false,
		},
		{
			name: "failure case - error in repo.DeleteNote()",
			given: func(r *interfaces.MockINotesRepository) {
				r.EXPECT().DeleteNote(mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			args: args{
				ctx: context.Background(),
				request: models.DeleteNoteRequest{
					Id: 123,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := interfaces.MockINotesRepository{}
			tt.given(&mockRepo)
			s := &notesService{
				repo:   &mockRepo,
				logger: loggers.NewLogger(),
			}
			err := s.DeleteNote(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("notesService.DeleteNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
