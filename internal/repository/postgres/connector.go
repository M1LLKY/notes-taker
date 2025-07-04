package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"notes-taker/internal/config"
	"notes-taker/internal/logger"
	"time"
)

type DBConnector struct {
	Pool *pgxpool.Pool
}

func NewDBConnector(cfg *config.Config) (*DBConnector, error) {
	serverConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	serverConfig.ConnConfig.Host = cfg.PostgresEndpoint
	serverConfig.ConnConfig.Port = 5432
	serverConfig.ConnConfig.Database = cfg.PostgresDatabase
	serverConfig.ConnConfig.User = cfg.PostgresUsername
	serverConfig.ConnConfig.Password = cfg.PostgresPassword
	serverConfig.ConnConfig.TLSConfig = nil
	serverConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(context.Background(), serverConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	logger.Get().Info(context.Background(), "db connection established", logrus.Fields{
		"host": cfg.PostgresEndpoint,
		"db":   cfg.PostgresDatabase,
	})

	return &DBConnector{Pool: pool}, nil
}
