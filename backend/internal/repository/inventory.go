package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) domain.InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) ListSessions(ctx context.Context, limit, offset int) ([]domain.InventorySession, int, error) {
	var daos []dao.InventorySessionDAO
	var total int64

	model := r.db.WithContext(ctx).Model(&dao.InventorySessionDAO{})
	
	if err := model.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		model = model.Limit(limit)
	}
	if offset > 0 {
		model = model.Offset(offset)
	}

	if err := model.Order("started_at DESC").Find(&daos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]domain.InventorySession, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}

	return result, int(total), nil
}

func (r *inventoryRepository) GetSessionByID(ctx context.Context, id string) (*domain.InventorySession, error) {
	var d dao.InventorySessionDAO
	if err := r.db.WithContext(ctx).Preload("Items").First(&d, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *inventoryRepository) CreateSession(ctx context.Context, s *domain.InventorySession) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	d := dao.FromInventorySessionDomain(*s)
	return r.db.WithContext(ctx).Create(&d).Error
}

func (r *inventoryRepository) UpdateSession(ctx context.Context, s *domain.InventorySession) error {
	d := dao.FromInventorySessionDomain(*s)
	return r.db.WithContext(ctx).Save(&d).Error
}

func (r *inventoryRepository) AddCount(ctx context.Context, item *domain.InventoryItem) error {
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	d := dao.FromInventoryItemDomain(*item)
	return r.db.WithContext(ctx).Create(&d).Error
}
