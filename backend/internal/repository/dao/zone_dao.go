package dao

import (
	"time"

	"github.com/ima/diplom-backend/internal/domain"
)

type ZoneDAO struct {
	ID             string          `gorm:"column:id;primaryKey"`
	Name           string          `gorm:"column:name;uniqueIndex"`
	Type           domain.ZoneType `gorm:"column:type"`
	Description    string          `gorm:"column:description"`
	TemperatureMin float64         `gorm:"column:temp_min"`
	TemperatureMax float64         `gorm:"column:temp_max"`
	Capacity       int             `gorm:"column:capacity"`
	IsActive       bool            `gorm:"column:is_active;default:true"`
	CreatedAt      time.Time       `gorm:"column:created_at"`
	UpdatedAt      time.Time       `gorm:"column:updated_at"`
}

func (ZoneDAO) TableName() string {
	return "warehouse_zones"
}

func (z ZoneDAO) ToDomain() domain.Zone {
	return domain.Zone{
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

func FromZoneDomain(z domain.Zone) ZoneDAO {
	return ZoneDAO{
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
