package mapper

import "notes-taker/internal/models"

func MapNoteDTOFromNoteDb(db models.NoteDB) models.NoteDTO {
	dto := models.NoteDTO{
		ID:        db.ID,
		Title:     db.Title,
		Content:   db.Content,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
	}
	return dto
}
