package repositories

import (
	"context"
	"notes-server/db"
	"notes-server/interfaces"
	"notes-server/loggers"
	"notes-server/models"
	"notes-server/utils"
)

type notesRepository struct {
	db     db.DB
	logger *loggers.Logger
}

func NewNotesRepository(db db.DB, logger *loggers.Logger) interfaces.INotesRepository {
	return &notesRepository{db: db, logger: logger}
}

func (r *notesRepository) GetNotes(ctx context.Context, email string) ([]models.Note, error) {
	r.logger.Info(ctx, "Entering notesRepository.GetNotes()")
	defer r.logger.Info(ctx, "Exiting notesRepository.GetNotes()")
	notes := make([]models.Note, 0)
	txn := r.db.Txn(ctx, false)
	rows, err := txn.Get("notes", "created_by", email)
	if err != nil {
		txn.Abort()
		r.logger.Warn(ctx, "error in notesRepository.GetNotes(), error from txn.Get()", err)
		return []models.Note{}, err
	}
	txn.Commit()
	for obj := rows.Next(); obj != nil; obj = rows.Next() {
		note := obj.(*models.Note)
		notes = append(notes, models.Note{
			Id:   note.Id,
			Note: note.Note,
		})
	}
	return notes, nil
}

func (r *notesRepository) AddNote(ctx context.Context, request models.AddNoteRequest) (int32, error) {
	r.logger.Info(ctx, "Entering notesRepository.AddNote()")
	defer r.logger.Info(ctx, "Exiting notesRepository.AddNote()")
	txn := r.db.Txn(ctx, true)
	note := models.Note{
		Note:      request.Note,
		CreatedBy: request.Email,
		Id:        utils.NewID(),
	}
	err := txn.Insert("notes", &note)
	if err != nil {
		txn.Abort()
		r.logger.Warn(ctx, "error in notesRepository.AddNote(), error from txn.Insert()", err)
		return 0, err
	}
	txn.Commit()
	return note.Id, nil
}

func (r *notesRepository) DeleteNote(ctx context.Context, noteID int32) error {
	r.logger.Info(ctx, "Entering notesRepository.DeleteNote()")
	defer r.logger.Info(ctx, "Exiting notesRepository.DeleteNote()")
	txn := r.db.Txn(ctx, true)
	err := txn.Delete("notes", models.Note{Id: noteID})
	if err != nil {
		txn.Abort()
		r.logger.Warn(ctx, "error in notesRepository.DeleteNote(), error from txn.Delete()", err)
		return err
	}
	txn.Commit()
	return nil
}
