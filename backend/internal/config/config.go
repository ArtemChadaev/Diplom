package config

import (
	"context"
	"errors"
	
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
	GoogleClientID string `env:"GOOGLE_CLIENT_ID"` //Не уверен что надо TODO: если не надо то удалить

	// Valkey
	ValkeyHost     string `env:"VALKEY_HOST" env-default:"valkey"`
	ValkeyPort     string `env:"VALKEY_PORT" env-default:"6379"`
	ValkeyPassword string `env:"VALKEY_PASSWORD"`

	// Gmail SMTP
	MAILER_SMTP_SERVER string `env:"SMTP_SERVER" env-default:"smtp.gmail.com"`
	MAILER_SMTP_PORT   string `env:"SMTP_PORT" env-default:"587"`
	MAILER_USERNAME   string `env:"SMTP_USER" env-required:"true"`
	MAILER_PASSWORD   string `env:"SMTP_PASS" env-required:"true"`

	OTPHMACSecret string `env:"OTP_HMAC_SECRET" env-required:"true"`
	UploadDir     string `env:"UPLOAD_DIR" env-default:"./uploads"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		logger.FromContext(context.Background()).Warn(".env file not found, loading from environment variables")
	}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, errors.New("config error: " + err.Error())
	}

	return &cfg, nil
}
