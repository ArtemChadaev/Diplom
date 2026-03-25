package repository

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

// NewPostgresDB создаёт подключение '*gorm.DB' к PostgreSQL
func NewPostgresDB(cfg PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode,
	)

	// Подключение через GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Получаем низкоуровневый *sql.DB для проверки пинга
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	slog.Info("successfully connected to PostgreSQL via GORM")

	return db, nil
}
