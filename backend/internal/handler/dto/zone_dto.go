package dto

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type ZoneResponse struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Type           domain.ZoneType `json:"type"`
	Description    string          `json:"description"`
	TemperatureMin float64         `json:"temp_min"`
	TemperatureMax float64         `json:"temp_max"`
	Capacity       int             `json:"capacity"`
	IsActive       bool            `json:"is_active"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func ToZoneResponse(z domain.Zone) ZoneResponse {
	return ZoneResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		Description:    z.Description,
		TemperatureMin: z.TemperatureMin,
		TemperatureMax: z.TemperatureMax,
		Capacity:       z.Capacity,
		IsActive:       z.IsActive,
		CreatedAt:      z.CreatedAt,
		UpdatedAt:      z.UpdatedAt,
	}
}

type CreateZoneRequest struct {
	Name           string          `json:"name" validate:"required"`
	Type           domain.ZoneType `json:"type" validate:"required,oneof=ambient cool cold narcotic quarantine"`
	Description    string          `json:"description"`
	TemperatureMin float64         `json:"temp_min"`
	TemperatureMax float64         `json:"temp_max"`
	Capacity       int             `json:"capacity"`
}

func (r CreateZoneRequest) ToDomain() domain.Zone {
	return domain.Zone{
		Name:           r.Name,
		Type:           r.Type,
		Description:    r.Description,
		TemperatureMin: r.TemperatureMin,
		TemperatureMax: r.TemperatureMax,
		Capacity:       r.Capacity,
		IsActive:       true,
	}
}

type UpdateZoneRequest struct {
	Name           *string          `json:"name"`
	Type           *domain.ZoneType `json:"type"`
	Description    *string          `json:"description"`
	TemperatureMin *float64         `json:"temp_min"`
	TemperatureMax *float64         `json:"temp_max"`
	Capacity       *int             `json:"capacity"`
	IsActive       *bool            `json:"is_active"`
}

func (r UpdateZoneRequest) ApplyTo(z *domain.Zone) {
	if r.Name != nil { z.Name = *r.Name }
	if r.Type != nil { z.Type = *r.Type }
	if r.Description != nil { z.Description = *r.Description }
	if r.TemperatureMin != nil { z.TemperatureMin = *r.TemperatureMin }
	if r.TemperatureMax != nil { z.TemperatureMax = *r.TemperatureMax }
	if r.Capacity != nil { z.Capacity = *r.Capacity }
	if r.IsActive != nil { z.IsActive = *r.IsActive }
}
