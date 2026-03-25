package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
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
	AdminUser      string `env:"ADMIN_USER" env-default:"admin"`
	AdminPassword  string `env:"ADMIN_PASSWORD" env-required:"true"`
	GoogleClientID string `env:"GOOGLE_CLIENT_ID"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		slog.Warn(".env file not found, loading from environment variables")
	}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return &cfg, nil
}
