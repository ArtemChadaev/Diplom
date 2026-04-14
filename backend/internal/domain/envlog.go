package domain

import (
	"context"
	"time"
)

// EnvironmentLog — запись журнала температурного режима и влажности.
type EnvironmentLog struct {
	ID          string    `json:"id"`
	ZoneID      string    `json:"zone_id"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	RecordedBy  int       `json:"recorded_by"` // UserID
	RecordedAt  time.Time `json:"recorded_at"`
	Notes       string    `json:"notes"`
}

// EnvironmentLogRepository — интерфейс для работы с журналом.
type EnvironmentLogRepository interface {
	List(ctx context.Context, zoneID string, limit, offset int) ([]EnvironmentLog, int, error)
	Create(ctx context.Context, log *EnvironmentLog) error
}

// EnvironmentLogService — бизнес-логика журнала.
type EnvironmentLogService interface {
	ListLogs(ctx context.Context, zoneID string, limit, offset int) ([]EnvironmentLog, int, error)
	RecordLogs(ctx context.Context, userID int, logs []EnvironmentLog) error
}
