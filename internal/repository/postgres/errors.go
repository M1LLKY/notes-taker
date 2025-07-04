package postgres

import "errors"

var (
	ErrSelect     = errors.New("ошибка выборки из базы данных")
	ErrCreateUser = errors.New("ошибка создания пользователя")
)
