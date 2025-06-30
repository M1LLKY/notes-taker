package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"notes-taker/internal/storage"
	"time"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	pool, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) CreateUser(username string) error {
	const op = "storage.postgres.CreateUser"

	createUserQuery := "INSERT INTO users(username, created_at) values($1, $2)"

	_, err := s.pool.Exec(context.Background(), createUserQuery, username, time.Now())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.UserExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
