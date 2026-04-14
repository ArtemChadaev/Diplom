package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type recalledBatchRepository struct {
	db *gorm.DB
}

func NewRecalledBatchRepository(db *gorm.DB) domain.RecalledBatchRepository {
	return &recalledBatchRepository{db: db}
}

func (r *recalledBatchRepository) List(ctx context.Context, limit, offset int) ([]domain.RecalledBatch, int, error) {
	var daos []dao.RecalledBatchDAO
	var total int64

	model := r.db.WithContext(ctx).Model(&dao.RecalledBatchDAO{})
	
	if err := model.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		model = model.Limit(limit)
	}
	if offset > 0 {
		model = model.Offset(offset)
	}

	if err := model.Order("issued_at DESC").Find(&daos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]domain.RecalledBatch, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}

	return result, int(total), nil
}

func (r *recalledBatchRepository) GetBySerial(ctx context.Context, serial string) (*domain.RecalledBatch, error) {
	var d dao.RecalledBatchDAO
	if err := r.db.WithContext(ctx).First(&d, "serial_number = ?", serial).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *recalledBatchRepository) Create(ctx context.Context, b *domain.RecalledBatch) error {
	if b.ID == "" {
		b.ID = uuid.NewString()
	}
	d := dao.FromRecalledBatchDomain(*b)
	return r.db.WithContext(ctx).Create(&d).Error
}
