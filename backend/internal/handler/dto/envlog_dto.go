package dto

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type EnvLogResponse struct {
	ID          string    `json:"id"`
	ZoneID      string    `json:"zone_id"`
	Shift       string    `json:"shift"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	RecordedBy  int       `json:"recorded_by"`
	RecordedAt  time.Time `json:"recorded_at"`
	Notes       string    `json:"notes"`
}

func ToEnvLogResponse(e domain.EnvironmentLog) EnvLogResponse {
	return EnvLogResponse{
		ID:          e.ID,
		ZoneID:      e.ZoneID,
		Shift:       e.Shift,
		Temperature: e.Temperature,
		Humidity:    e.Humidity,
		RecordedBy:  e.RecordedBy,
		RecordedAt:  e.RecordedAt,
		Notes:       e.Notes,
	}
}

type RecordEnvLogRequest struct {
	ZoneID      string  `json:"zone_id" validate:"required,uuid"`
	Shift       string  `json:"shift" validate:"required,oneof=morning evening"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Notes       string  `json:"notes"`
}

func (r RecordEnvLogRequest) ToDomain() domain.EnvironmentLog {
	return domain.EnvironmentLog{
		ZoneID:      r.ZoneID,
		Shift:       r.Shift,
		Temperature: r.Temperature,
		Humidity:    r.Humidity,
		Notes:       r.Notes,
	}
}

type EnvLogListResponse struct {
	Total int              `json:"total"`
	Logs  []EnvLogResponse `json:"logs"`
}
