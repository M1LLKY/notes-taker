package config

import (
	"errors"
	"github.com/caarlos0/env/v11"
	"time"
)

type Config struct {
	HTTPListenAddr   string        `env:"HTTP_LISTEN_ADDR" envDefault:":8080"`
	LogLevel         string        `env:"LOG_LEVEL" envDefault:"DEBUG"`
	PostgresEndpoint string        `env:"POSTGRES_ENDPOINT"`
	PostgresDatabase string        `env:"POSTGRES_DATABASE"`
	PostgresUsername string        `env:"POSTGRES_USERNAME"`
	PostgresPassword string        `env:"POSTGRES_PASSWORD"`
	JwtTTLDuration   time.Duration `env:"JWT_TTL" envDefault:"15m"`
	SigningKey       string        `env:"SIGNING_KEY"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, errors.New("failed to parse config")
	}
	return cfg, nil
}
