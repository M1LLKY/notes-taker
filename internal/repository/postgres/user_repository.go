package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"notes-taker/internal/models"
	"notes-taker/internal/repository"
	"time"
)

type UserRepository struct {
	Repository
}

func NewUserRepository(connector *DBConnector) repository.UserRepository {
	return &UserRepository{Repository{pool: connector.Pool}}
}

func scanUserRow(row pgx.Row) (*models.UserDB, error) {
	var user models.UserDB

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.UserDB, error) {
	query := `
		SELECT id, username, password_hash, created_at FROM users
		WHERE username = $1;
	`
	row := r.pool.QueryRow(ctx, query, username)
	user, err := scanUserRow(row)
	if err != nil {
		return nil, ErrSelect
	}
	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, username, passwordHash string) (int, error) {
	query := `
		INSERT INTO users(username, password_hash, created_at) values ($1, $2, $3)
		RETURNING id;
	`
	var userID int
	err := r.pool.QueryRow(ctx, query, username, passwordHash, time.Now()).Scan(&userID)
	if err != nil {
		return 0, ErrCreateUser
	}
	return userID, nil
}
