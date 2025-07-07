package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"notes-taker/internal/app/server"
	"notes-taker/internal/config"
	"notes-taker/internal/dependencies"
	"notes-taker/internal/logger"
	"notes-taker/internal/repository/postgres"
	"notes-taker/internal/service"
	"notes-taker/pkg/auth"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Server *server.HTTPServer
	Pool   *pgxpool.Pool
}

func New(cfg *config.Config) (*App, error) {
	connector, err := postgres.NewDBConnector(cfg)
	if err != nil {
		return nil, err
	}
	tokenManager, err := auth.NewManager(cfg.SigningKey)
	if err != nil {
		return nil, err
	}
	userRepository := postgres.NewUserRepository(connector)
	noteRepository := postgres.NewNoteRepository(connector)
	userService := service.NewUserService(userRepository, tokenManager, cfg)
	noteService := service.NewNoteService(noteRepository)
	deps := dependencies.New(
		cfg,
		userRepository,
		noteRepository,
		userService,
		noteService,
		tokenManager,
	)
	srv, err := server.New(deps)
	if err != nil {
		return nil, err
	}
	app := &App{
		Server: srv,
		Pool:   connector.Pool,
	}
	return app, nil
}

func (a *App) Run(ctx context.Context, cfg *config.Config) error {
	serverErrors := make(chan error, 1)

	go func() {
		logger.Get().Info(ctx, fmt.Sprintf("starting server at: %s", cfg.HTTPListenAddr))
		if err := a.Server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Get().Info(ctx, fmt.Sprintf("received shutdown signal: %s", sig))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.Server.ShutDown(ctx); err != nil {
			return err
		}
		return http.ErrServerClosed
	case err := <-serverErrors:
		return err
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	a.Pool.Close()
	return a.Server.ShutDown(ctx)
}
