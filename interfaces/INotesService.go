package interfaces

import (
	"context"
	"notes-server/models"
)

type INotesService interface {
	GetNotes(ctx context.Context) ([]models.Note, error)
	AddNote(ctx context.Context, request models.AddNoteRequest) (models.AddNoteRsponse, error)
	DeleteNote(ctx context.Context, request models.DeleteNoteRequest) error
}
