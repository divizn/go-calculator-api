package utils

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type (
	IConfig struct {
		SERVER_ADDR string `env:"SERVER_ADDR" validate:"required"`
		APP_ENV     string `env:"APP_ENV" validate:"required"`
		DB_HOST     string `env:"DB_HOST" validate:"required"`
		DB_PORT     int    `env:"DB_PORT" validate:"required"`
		DB_DATABASE string `env:"DB_DATABASE" validate:"required"`
		DB_USERNAME string `env:"DB_USERNAME" validate:"required"`
		DB_PASSWORD string `env:"DB_PASSWORD" validate:"required"`
		DB_SCHEMA   string `env:"DB_SCHEMA" validate:"required"`
		DB_URL      string `env:"DB_URL" validate:"required"`
		JWT_SECRET  string `env:"JWT_SECRET" validate:"required"`
	}
)

func NewConfig() (*IConfig, error) {
	cfg := &IConfig{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading environment variables: %v", err)
	}

	return cfg, nil
}
