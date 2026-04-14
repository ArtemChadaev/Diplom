package dao

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type EnvironmentLogDAO struct {
	ID          string    `gorm:"column:id;primaryKey"`
	ZoneID      string    `gorm:"column:zone_id;index"`
	Temperature float64   `gorm:"column:temperature"`
	Humidity    float64   `gorm:"column:humidity"`
	RecordedBy  int       `gorm:"column:recorded_by"`
	RecordedAt  time.Time `gorm:"column:recorded_at;index"`
	Notes       string    `gorm:"column:notes"`
}

func (EnvironmentLogDAO) TableName() string {
	return "environment_logs"
}

func (e EnvironmentLogDAO) ToDomain() domain.EnvironmentLog {
	return domain.EnvironmentLog{
		ID:          e.ID,
		ZoneID:      e.ZoneID,
		Temperature: e.Temperature,
		Humidity:    e.Humidity,
		RecordedBy:  e.RecordedBy,
		RecordedAt:  e.RecordedAt,
		Notes:       e.Notes,
	}
}

func FromEnvLogDomain(e domain.EnvironmentLog) EnvironmentLogDAO {
	return EnvironmentLogDAO{
		ID:          e.ID,
		ZoneID:      e.ZoneID,
		Temperature: e.Temperature,
		Humidity:    e.Humidity,
		RecordedBy:  e.RecordedBy,
		RecordedAt:  e.RecordedAt,
		Notes:       e.Notes,
	}
}
