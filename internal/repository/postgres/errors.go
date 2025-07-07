package postgres

import "errors"

var (
	ErrSelect       = errors.New("ошибка выборки из базы данных")
	ErrCreateUser   = errors.New("ошибка создания пользователя")
	ErrCreateNote   = errors.New("ошибка создания заметки")
	ErrNoteNotFound = errors.New("не удалось найти заметку по указанному ID")
	ErrUpdateNote   = errors.New("не удалось обновить заметку по указанному ID")
	ErrDeleteNote   = errors.New("не удалось удалить заметку по указанному ID")
)
