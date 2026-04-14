package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type claimRepository struct {
	db *gorm.DB
}

func NewClaimRepository(db *gorm.DB) domain.ClaimRepository {
	return &claimRepository{db: db}
}

func (r *claimRepository) List(ctx context.Context, limit, offset int) ([]domain.Claim, int, error) {
	var daos []dao.ClaimDAO
	var total int64

	model := r.db.WithContext(ctx).Model(&dao.ClaimDAO{})
	
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

	result := make([]domain.Claim, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}

	return result, int(total), nil
}

func (r *claimRepository) GetByID(ctx context.Context, id string) (*domain.Claim, error) {
	var d dao.ClaimDAO
	if err := r.db.WithContext(ctx).First(&d, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *claimRepository) Create(ctx context.Context, c *domain.Claim) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	d := dao.FromClaimDomain(*c)
	return r.db.WithContext(ctx).Create(&d).Error
}

func (r *claimRepository) Update(ctx context.Context, c *domain.Claim) error {
	d := dao.FromClaimDomain(*c)
	return r.db.WithContext(ctx).Save(&d).Error
}
