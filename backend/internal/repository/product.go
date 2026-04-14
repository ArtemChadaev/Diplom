package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
	var daos []dao.ProductDAO
	var total int64

	tx := r.db.WithContext(ctx).Model(&dao.ProductDAO{})

	if filter.Query != "" {
		q := "%" + strings.ToLower(filter.Query) + "%"
		tx = tx.Where("LOWER(name) LIKE ? OR LOWER(generic_name) LIKE ? OR LOWER(sku) LIKE ?", q, q, q)
	}

	if filter.IsJNVLP != nil {
		tx = tx.Where("is_jnvlp = ?", *filter.IsJNVLP)
	}

	if filter.ATCCode != "" {
		tx = tx.Where("atc_code = ?", filter.ATCCode)
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Limit > 0 {
		tx = tx.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		tx = tx.Offset(filter.Offset)
	}

	if err := tx.Order("name ASC").Find(&daos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]domain.Product, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}

	return result, int(total), nil
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var d dao.ProductDAO
	if err := r.db.WithContext(ctx).First(&d, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	var d dao.ProductDAO
	if err := r.db.WithContext(ctx).First(&d, "sku = ?", sku).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *productRepository) Create(ctx context.Context, p *domain.Product) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	d := dao.FromProductDomain(*p)
	return r.db.WithContext(ctx).Create(&d).Error
}

func (r *productRepository) Update(ctx context.Context, p *domain.Product) error {
	d := dao.FromProductDomain(*p)
	// We use Save to update all fields including zero values if provided.
	// But usually, we only update specific fields. 
	// For simplicity, we use Save here.
	return r.db.WithContext(ctx).Save(&d).Error
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&dao.ProductDAO{}, "id = ?", id).Error
}
