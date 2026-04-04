package config

import (
	"context"
	"fmt"
	
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

type Config struct {
	Env  string `env:"ENV" env-default:"development"`
	// env-default задает значение, если переменной нет ни в .env, ни в системе
	Port string `env:"PORT" env-default:"8080"`

	// Настройки Postgres
	// env-required:"true" — приложение упадет с ошибкой, если переменная не задана
	DBHost     string `env:"DB_HOST" env-required:"true"`
	DBPort     string `env:"DB_PORT" env-default:"5432"`
	DBUser     string `env:"DB_USER" env-required:"true"`
	DBName     string `env:"DB_NAME" env-required:"true"`
	DBPassword string `env:"DB_PASSWORD" env-required:"true"`

	JWTSecret      string `env:"JWT_SECRET" env-required:"true"`
	AdminEmail     string `env:"ADMIN_EMAIL" env-required:"true"`
	GoogleClientID string `env:"GOOGLE_CLIENT_ID"`

	// Valkey
	ValkeyHost     string `env:"VALKEY_HOST" env-default:"valkey"`
	ValkeyPort     string `env:"VALKEY_PORT" env-default:"6379"`
	ValkeyPassword string `env:"VALKEY_PASSWORD"`

	// UniSender Mail
	UniSenderAPIKey    string `env:"UNISENDER_API_KEY"`
	UniSenderAPIURL    string `env:"UNISENDER_API_URL" env-default:"https://api.unisender.com/ru/api/sendEmail"`
	UniSenderFromEmail string `env:"UNISENDER_FROM_EMAIL"`
	UniSenderFromName  string `env:"UNISENDER_FROM_NAME"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		logger.FromContext(context.Background()).Warn(".env file not found, loading from environment variables")
	}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return &cfg, nil
}
