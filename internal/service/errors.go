package service

import "errors"

var (
	ErrNoteForbidden = errors.New("попытка доступа к чужой заметке")
	ErrNoteNotFound  = errors.New("заметка не найдена")
)
