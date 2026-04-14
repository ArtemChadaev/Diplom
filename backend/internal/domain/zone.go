package domain

import (
	"context"
	"time"
)

type ZoneType string

const (
	ZoneTypeAmbient      ZoneType = "ambient"      // 15-25°C
	ZoneTypeCool         ZoneType = "cool"         // 8-15°C
	ZoneTypeColdStorage  ZoneType = "cold"         // 2-8°C
	ZoneTypeNarcotic     ZoneType = "narcotic"     // Locked storage
	ZoneTypeQuarantine   ZoneType = "quarantine"   // Quality control
)

// Zone — складская зона.
type Zone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`        // Код/название зоны
	Type        ZoneType  `json:"type"`        // Тип (ambient, cold, etc.)
	Description string    `json:"description"`
	TemperatureMin float64 `json:"temp_min"`
	TemperatureMax float64 `json:"temp_max"`
	Capacity    int       `json:"capacity"`    // Вместимость (в условных ед./паллетах)
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ZoneRepository — интерфейс для работы с зонами.
type ZoneRepository interface {
	List(ctx context.Context) ([]Zone, error)
	GetByID(ctx context.Context, id string) (*Zone, error)
	Create(ctx context.Context, z *Zone) error
	Update(ctx context.Context, z *Zone) error
	Delete(ctx context.Context, id string) error
}

// ZoneService — бизнес-логика зон.
type ZoneService interface {
	ListZones(ctx context.Context) ([]Zone, error)
	GetZone(ctx context.Context, id string) (*Zone, error)
	CreateZone(ctx context.Context, callerRole UserRole, z *Zone) (*Zone, error)
	UpdateZone(ctx context.Context, callerRole UserRole, z *Zone) (*Zone, error)
	DeleteZone(ctx context.Context, callerRole UserRole, id string) error
}
