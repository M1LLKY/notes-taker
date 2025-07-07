package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"notes-taker/internal/models"
	"notes-taker/internal/repository"
	"time"
)

type NoteRepository struct {
	Repository
}

func NewNoteRepository(connector *DBConnector) repository.NoteRepository {
	return &NoteRepository{Repository{pool: connector.Pool}}
}

func scanNoteRow(row pgx.Row) (*models.NoteDB, error) {
	var note models.NoteDB

	if err := row.Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &note, nil
}

func scanNoteListRow(rows pgx.Rows) ([]models.NoteDB, error) {
	var items []models.NoteDB
	for rows.Next() {
		var note models.NoteDB

		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, note)
	}
	return items, nil
}

func (r *NoteRepository) CreateNote(ctx context.Context, title, content string, userID int) (int, error) {
	query := `
		INSERT INTO notes(title, content, user_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	var taskID int
	err := r.pool.QueryRow(ctx, query, title, content, userID, time.Now()).Scan(&taskID)
	if err != nil {
		return 0, ErrCreateNote
	}
	return taskID, nil
}

func (r *NoteRepository) GetAllNotes(ctx context.Context, userID int) ([]models.NoteDB, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM notes
		WHERE user_id = $1
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, ErrSelect
	}
	notes, err := scanNoteListRow(rows)
	if err != nil {
		return nil, ErrSelect
	}
	return notes, nil
}
func (r *NoteRepository) GetNoteByID(ctx context.Context, noteID int) (*models.NoteDB, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM notes
		WHERE id = $1 
	`
	row := r.pool.QueryRow(ctx, query, noteID)
	note, err := scanNoteRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoteNotFound
		}
		return nil, ErrSelect
	}
	return note, nil
}
func (r *NoteRepository) UpdateNoteByID(ctx context.Context, title, content string, noteID int) error {
	query := `
		UPDATE notes
		SET title = $1,
		    content = $2
		WHERE id = $3;
	`
	tag, err := r.pool.Exec(ctx, query, title, content, noteID)
	if err != nil {
		return ErrUpdateNote
	}
	rows := tag.RowsAffected()
	if rows == 0 {
		return ErrNoteNotFound
	}

	return nil
}
func (r *NoteRepository) DeleteNoteByID(ctx context.Context, noteID int) error {
	query := `
		DELETE FROM notes
		WHERE id = $1;
	`
	tag, err := r.pool.Exec(ctx, query, noteID)
	if err != nil {
		return ErrDeleteNote
	}
	if tag.RowsAffected() == 0 {
		return ErrNoteNotFound
	}

	return nil
}
