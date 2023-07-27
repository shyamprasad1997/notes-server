package interfaces

import (
	"context"
	"notes-server/models"
)

type INotesRepository interface {
	GetNotes(ctx context.Context, email string) ([]models.Note, error)
	AddNote(ctx context.Context, request models.AddNoteRequest) (int32, error)
	DeleteNote(ctx context.Context, noteID int32) error
}
