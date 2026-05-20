package repository

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type batchRepository struct {
	db *gorm.DB
}

func NewBatchRepository(db *gorm.DB) domain.BatchRepository {
	return &batchRepository{db: db}
}

func (r *batchRepository) List(ctx context.Context, filter domain.BatchFilter) ([]domain.Batch, int, error) {
	var daos []dao.BatchDAO
	var total int64

	model := r.db.WithContext(ctx).Model(&dao.BatchDAO{})
	
	if filter.ProductID != "" {
		model = model.Where("product_id = ?", filter.ProductID)
	}
	if filter.ZoneID != "" {
		model = model.Where("zone_id = ?", filter.ZoneID)
	}
	if filter.Status != "" {
		model = model.Where("status = ?", filter.Status)
	}

	if err := model.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Limit > 0 {
		model = model.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		model = model.Offset(filter.Offset)
	}

	if err := model.Order("expiry_date ASC").Find(&daos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]domain.Batch, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}

	return result, int(total), nil
}

func (r *batchRepository) GetByID(ctx context.Context, id string) (*domain.Batch, error) {
	var d dao.BatchDAO
	if err := r.db.WithContext(ctx).First(&d, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *batchRepository) Update(ctx context.Context, b *domain.Batch) error {
	d := dao.FromBatchDomain(*b)
	return r.db.WithContext(ctx).Save(&d).Error
}

func (r *batchRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&dao.BatchDAO{}, "id = ?", id).Error
}

func (r *batchRepository) BlockAllByProductID(ctx context.Context, productID string) error {
	return r.db.WithContext(ctx).Model(&dao.BatchDAO{}).
		Where("product_id = ? AND status = ?", productID, domain.BatchStatusAvailable).
		Update("status", domain.BatchStatusBlocked).Error
}

func (r *batchRepository) ListAvailableSorted(ctx context.Context, productID string) ([]domain.Batch, error) {
	var daos []dao.BatchDAO
	err := r.db.WithContext(ctx).Model(&dao.BatchDAO{}).
		Where("product_id = ? AND status = ? AND quantity > 0", productID, domain.BatchStatusAvailable).
		Order("expiry_date ASC, manufacture_date ASC").
		Find(&daos).Error
	if err != nil {
		return nil, err
	}
	result := make([]domain.Batch, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}
	return result, nil
}

func (r *batchRepository) GetTotalStock(ctx context.Context, productID string) (int, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&dao.BatchDAO{}).
		Where("product_id = ? AND status = ?", productID, domain.BatchStatusAvailable).
		Select("COALESCE(SUM(quantity), 0)").
		Row().Scan(&total)
	if err != nil {
		return 0, err
	}
	return int(total), nil
}
