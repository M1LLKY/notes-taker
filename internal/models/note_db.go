package models

import "time"

type NoteDB struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
