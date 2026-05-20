package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type environmentLogRepository struct {
	db *gorm.DB
}

func NewEnvironmentLogRepository(db *gorm.DB) domain.EnvironmentLogRepository {
	return &environmentLogRepository{db: db}
}

func (r *environmentLogRepository) List(ctx context.Context, zoneID string, limit, offset int) ([]domain.EnvironmentLog, int, error) {
	var daos []dao.EnvironmentLogDAO
	var total int64

	model := r.db.WithContext(ctx).Model(&dao.EnvironmentLogDAO{})
	if zoneID != "" {
		model = model.Where("zone_id = ?", zoneID)
	}

	if err := model.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		model = model.Limit(limit)
	}
	if offset > 0 {
		model = model.Offset(offset)
	}

	if err := model.Order("recorded_at DESC").Find(&daos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]domain.EnvironmentLog, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}

	return result, int(total), nil
}

func (r *environmentLogRepository) Create(ctx context.Context, log *domain.EnvironmentLog) error {
	if log.ID == "" {
		log.ID = uuid.NewString()
	}
	d := dao.FromEnvLogDomain(*log)
	return r.db.WithContext(ctx).Create(&d).Error
}

func (r *environmentLogRepository) ExistsByZoneShiftDate(ctx context.Context, zoneID string, shift string, date time.Time) (bool, error) {
	var count int64
	dateStr := date.Format("2006-01-02")
	if err := r.db.WithContext(ctx).
		Model(&dao.EnvironmentLogDAO{}).
		Where("zone_id = ? AND shift = ? AND CAST(recorded_at AT TIME ZONE 'UTC' AS DATE) = ?", zoneID, shift, dateStr).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
