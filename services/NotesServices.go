package services

import (
	"context"
	"notes-server/interfaces"
	"notes-server/loggers"
	"notes-server/models"
	"notes-server/utils"
)

type notesService struct {
	repo   interfaces.INotesRepository
	logger *loggers.Logger
}

func NewNotesService(logger *loggers.Logger, repo interfaces.INotesRepository) interfaces.INotesService {
	return &notesService{
		repo:   repo,
		logger: logger,
	}
}

// GetNotes - retrieves all the notes of the user
func (s *notesService) GetNotes(ctx context.Context) ([]models.Note, error) {
	email := utils.GetEmailFromCtx(ctx)
	notes, err := s.repo.GetNotes(ctx, email)
	if err != nil {
		s.logger.Warn(ctx, "Error in notesService.GetNotes(), error from repo.GetNotes()")
		return []models.Note{}, err
	}
	return notes, nil
}

// AddNote - add a new note
func (s *notesService) AddNote(ctx context.Context, request models.AddNoteRequest) (models.AddNoteResponse, error) {
	email := utils.GetEmailFromCtx(ctx)
	request.Email = email
	id, err := s.repo.AddNote(ctx, request)
	if err != nil {
		s.logger.Warn(ctx, "Error in notesService.AddNote(), error from repo.AddNote()")
		return models.AddNoteResponse{}, err
	}
	return models.AddNoteResponse{Id: id}, nil
}

// DeleteNote - delete a note
func (s *notesService) DeleteNote(ctx context.Context, request models.DeleteNoteRequest) error {
	err := s.repo.DeleteNote(ctx, request.Id)
	if err != nil {
		s.logger.Warn(ctx, "Error in notesService.DeleteNote(), error from repo.DeleteNote()")
		return err
	}
	return nil
}
