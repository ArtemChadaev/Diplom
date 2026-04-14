package dao

import (
	"github.com/ima/diplom-backend/internal/domain"
)

type SystemSettingDAO struct {
	Key   string `gorm:"column:key;primaryKey"`
	Value string `gorm:"column:value"`
}

func (SystemSettingDAO) TableName() string {
	return "system_settings"
}

func (s SystemSettingDAO) ToDomain() domain.SystemSetting {
	return domain.SystemSetting{
		Key:   s.Key,
		Value: s.Value,
	}
}
