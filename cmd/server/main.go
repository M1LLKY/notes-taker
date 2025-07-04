package main

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"notes-taker/internal/app"
	"notes-taker/internal/config"
	"notes-taker/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if env := os.Getenv("ENV"); env == "" || env == "local" {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
	}
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	logger.Init(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(cfg)
	if err != nil {
		panic(err)
	}

	if err := application.Run(ctx, cfg); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Get().Info(context.Background(), "остановка приложения", logrus.Fields{
				"err": err.Error(),
			})
			return
		}
		panic(err)
	}
}
