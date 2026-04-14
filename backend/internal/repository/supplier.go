package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) domain.SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) List(ctx context.Context, limit, offset int) ([]domain.Supplier, int, error) {
	var daos []dao.SupplierDAO
	var total int64

	model := r.db.WithContext(ctx).Model(&dao.SupplierDAO{})
	
	if err := model.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		model = model.Limit(limit)
	}
	if offset > 0 {
		model = model.Offset(offset)
	}

	if err := model.Order("name ASC").Find(&daos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]domain.Supplier, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}

	return result, int(total), nil
}

func (r *supplierRepository) GetByID(ctx context.Context, id string) (*domain.Supplier, error) {
	var d dao.SupplierDAO
	if err := r.db.WithContext(ctx).First(&d, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *supplierRepository) GetByINN(ctx context.Context, inn string) (*domain.Supplier, error) {
	var d dao.SupplierDAO
	if err := r.db.WithContext(ctx).First(&d, "inn = ?", inn).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *supplierRepository) Create(ctx context.Context, s *domain.Supplier) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	d := dao.FromSupplierDomain(*s)
	return r.db.WithContext(ctx).Create(&d).Error
}

func (r *supplierRepository) Update(ctx context.Context, s *domain.Supplier) error {
	d := dao.FromSupplierDomain(*s)
	return r.db.WithContext(ctx).Save(&d).Error
}

func (r *supplierRepository) Delete(ctx context.Context, id string) error {
	// Suppliers do not use soft-delete in this current entity spec, but we could add it.
	// For now, hard delete or just set IsActive=false. 
	// The plan says "hard delete for suppliers or set IsActive". 
	// I'll stick to Delete in repo, but it can be changed to IsActive update in service.
	return r.db.WithContext(ctx).Delete(&dao.SupplierDAO{}, "id = ?", id).Error
}
