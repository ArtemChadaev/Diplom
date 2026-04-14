package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type zoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) domain.ZoneRepository {
	return &zoneRepository{db: db}
}

func (r *zoneRepository) List(ctx context.Context) ([]domain.Zone, error) {
	var daos []dao.ZoneDAO
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&daos).Error; err != nil {
		return nil, err
	}

	result := make([]domain.Zone, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}
	return result, nil
}

func (r *zoneRepository) GetByID(ctx context.Context, id string) (*domain.Zone, error) {
	var d dao.ZoneDAO
	if err := r.db.WithContext(ctx).First(&d, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *zoneRepository) Create(ctx context.Context, z *domain.Zone) error {
	if z.ID == "" {
		z.ID = uuid.NewString()
	}
	d := dao.FromZoneDomain(*z)
	return r.db.WithContext(ctx).Create(&d).Error
}

func (r *zoneRepository) Update(ctx context.Context, z *domain.Zone) error {
	d := dao.FromZoneDomain(*z)
	return r.db.WithContext(ctx).Save(&d).Error
}

func (r *zoneRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&dao.ZoneDAO{}, "id = ?", id).Error
}
