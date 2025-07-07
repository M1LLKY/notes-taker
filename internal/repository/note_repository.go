package repository

import (
	"context"
	"notes-taker/internal/models"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, title, content string, userID int) (int, error)
	GetAllNotes(ctx context.Context, userID int) ([]models.NoteDB, error)
	GetNoteByID(ctx context.Context, noteID int) (*models.NoteDB, error)
	UpdateNoteByID(ctx context.Context, title, content string, noteID int) error
	DeleteNoteByID(ctx context.Context, noteID int) error
}
