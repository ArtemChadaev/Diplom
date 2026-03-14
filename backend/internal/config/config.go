package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// Основные настройки приложения
	Port string `mapstructure:"PORT"`

	// Настройки Postgres
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	// Настройки Auth
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	AdminUser     string `mapstructure:"ADMIN_USER"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`
}

func Load() (*Config, error) {
	v := viper.New()

	// Читаем .env файл напрямую
	v.SetConfigFile(".env")
	v.SetConfigType("env")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ошибка чтения .env файла: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка распаковки конфига: %w", err)
	}

	return &cfg, nil
}
