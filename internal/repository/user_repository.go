package repository

import (
	"context"
	"notes-taker/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username, passwordHash string) (int, error)
	GetUserByUsername(ctx context.Context, username string) (*models.UserDB, error)
}
