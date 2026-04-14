package domain

import (
	"context"
)

// SystemSetting — глобальная настройка системы.
type SystemSetting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SystemSettingsRepository — интерфейс для работы с настройками.
type SystemSettingsRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	List(ctx context.Context) ([]SystemSetting, error)
}

// SystemSettingsService — бизнес-логика настроек.
type SystemSettingsService interface {
	GetSetting(ctx context.Context, key string) (string, error)
	UpdateSetting(ctx context.Context, callerRole UserRole, key, value string) error
	ListSettings(ctx context.Context) ([]SystemSetting, error)
}
