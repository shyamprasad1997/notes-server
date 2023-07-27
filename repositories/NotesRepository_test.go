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

func Test_notesRepository_GetNotes(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		given   func(*db.MockDB)
		args    args
		want    []models.Note
		wantErr bool
	}{
		{
			name: "success case",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				t := mockResultIterator{
					NextResp: &models.Note{
						Id:   123,
						Note: "test note",
					},
				}
				mockTxn.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything).Return(&t, nil)
				mockTxn.EXPECT().Commit()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx:   context.Background(),
				email: "test@gmail.com",
			},
			want: []models.Note{{
				Id:   123,
				Note: "test note",
			}},
			wantErr: false,
		},
		{
			name: "failure case - error in txn.Get()",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db errpr"))
				mockTxn.EXPECT().Abort()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx:   context.Background(),
				email: "test@gmail.com",
			},
			want:    []models.Note{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := db.MockDB{}
			tt.given(&mockDb)
			r := &notesRepository{
				db:     &mockDb,
				logger: loggers.NewLogger(),
			}
			got, err := r.GetNotes(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("notesRepository.GetNotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("notesRepository.GetNotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notesRepository_AddNote(t *testing.T) {
	type args struct {
		ctx     context.Context
		request models.AddNoteRequest
	}
	tests := []struct {
		name    string
		given   func(*db.MockDB)
		args    args
		want    int32
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
				request: models.AddNoteRequest{
					Email: "test@gmail.com",
					Note:  "test note",
				},
			},
			want:    123,
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
				request: models.AddNoteRequest{
					Email: "test@gmail.com",
					Note:  "test note",
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := db.MockDB{}
			tt.given(&mockDb)
			r := &notesRepository{
				db:     &mockDb,
				logger: loggers.NewLogger(),
			}
			_, err := r.AddNote(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("notesRepository.AddNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_notesRepository_DeleteNote(t *testing.T) {
	type args struct {
		ctx     context.Context
		request int32
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
				mockTxn.EXPECT().Delete(mock.Anything, mock.Anything).Return(nil)
				mockTxn.EXPECT().Commit()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx:     context.Background(),
				request: 123,
			},
			wantErr: false,
		},
		{
			name: "failure case - error in txn.Delete()",
			given: func(dab *db.MockDB) {
				mockTxn := db.MockMemDbTxn{}
				mockTxn.EXPECT().Delete(mock.Anything, mock.Anything).Return(errors.New("db error"))
				mockTxn.EXPECT().Abort()
				dab.EXPECT().Txn(mock.Anything, mock.Anything).Return(&mockTxn)
			},
			args: args{
				ctx:     context.Background(),
				request: 123,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := db.MockDB{}
			tt.given(&mockDb)
			r := &notesRepository{
				db:     &mockDb,
				logger: loggers.NewLogger(),
			}
			err := r.DeleteNote(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("notesRepository.DeleteNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type mockResultIterator struct {
	WatchChResp chan struct{}
	NextResp    interface{}
}

func (m mockResultIterator) WatchCh() <-chan struct{} {
	t := make(chan struct{})
	return t
}

func (m *mockResultIterator) Next() interface{} {
	value := m.NextResp
	m.NextResp = nil
	return value
}
