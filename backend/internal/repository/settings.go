package repository

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type systemSettingsRepository struct {
	db *gorm.DB
}

func NewSystemSettingsRepository(db *gorm.DB) domain.SystemSettingsRepository {
	return &systemSettingsRepository{db: db}
}

func (r *systemSettingsRepository) Get(ctx context.Context, key string) (string, error) {
	var d dao.SystemSettingDAO
	if err := r.db.WithContext(ctx).First(&d, "key = ?", key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return d.Value, nil
}

func (r *systemSettingsRepository) Set(ctx context.Context, key, value string) error {
	d := dao.SystemSettingDAO{Key: key, Value: value}
	return r.db.WithContext(ctx).Save(&d).Error
}

func (r *systemSettingsRepository) List(ctx context.Context) ([]domain.SystemSetting, error) {
	var daos []dao.SystemSettingDAO
	if err := r.db.WithContext(ctx).Find(&daos).Error; err != nil {
		return nil, err
	}
	res := make([]domain.SystemSetting, len(daos))
	for i, d := range daos {
		res[i] = d.ToDomain()
	}
	return res, nil
}
