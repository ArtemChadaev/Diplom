package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type inboundRepository struct {
	db *gorm.DB
}

func NewInboundRepository(db *gorm.DB) domain.InboundRepository {
	return &inboundRepository{db: db}
}

func (r *inboundRepository) List(ctx context.Context, limit, offset int) ([]domain.InboundReceiving, int, error) {
	var daos []dao.InboundReceivingDAO
	var total int64

	model := r.db.WithContext(ctx).Model(&dao.InboundReceivingDAO{})
	
	if err := model.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		model = model.Limit(limit)
	}
	if offset > 0 {
		model = model.Offset(offset)
	}

	if err := model.Order("created_at DESC").Find(&daos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]domain.InboundReceiving, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}

	return result, int(total), nil
}

func (r *inboundRepository) GetByID(ctx context.Context, id string) (*domain.InboundReceiving, error) {
	var d dao.InboundReceivingDAO
	if err := r.db.WithContext(ctx).Preload("Items").First(&d, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *inboundRepository) Create(ctx context.Context, i *domain.InboundReceiving) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	for idx := range i.Items {
		if i.Items[idx].ID == "" {
			i.Items[idx].ID = uuid.NewString()
		}
		i.Items[idx].InboundID = i.ID
	}
	d := dao.FromInboundDomain(*i)
	return r.db.WithContext(ctx).Create(&d).Error
}

func (r *inboundRepository) Update(ctx context.Context, i *domain.InboundReceiving) error {
	d := dao.FromInboundDomain(*i)
	// We use transaction to ensure items are updated correctly
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&d).Error; err != nil {
			return err
		}
		// For items, we might need a more complex sync logic.
		// For now, we assume Save handles everything if Omit isn't used.
		return nil
	})
}

func (r *inboundRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&dao.InboundItemDAO{}, "inbound_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&dao.InboundReceivingDAO{}, "id = ?", id).Error
	})
}

func (r *inboundRepository) AddItems(ctx context.Context, inboundID string, items []domain.InboundItem) error {
	daos := make([]dao.InboundItemDAO, len(items))
	for i, item := range items {
		if item.ID == "" {
			item.ID = uuid.NewString()
		}
		item.InboundID = inboundID
		daos[i] = dao.FromInboundItemDomain(item)
	}
	return r.db.WithContext(ctx).Create(&daos).Error
}

func (r *inboundRepository) UpdateItem(ctx context.Context, item *domain.InboundItem) error {
	d := dao.FromInboundItemDomain(*item)
	return r.db.WithContext(ctx).Save(&d).Error
}

func (r *inboundRepository) RemoveItem(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&dao.InboundItemDAO{}, "id = ?", id).Error
}
