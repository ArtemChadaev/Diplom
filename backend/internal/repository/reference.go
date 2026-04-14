package repository

import (
	"context"
	"strings"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type referenceRepository struct {
	db *gorm.DB
}

func NewReferenceRepository(db *gorm.DB) domain.ReferenceRepository {
	return &referenceRepository{db: db}
}

func (r *referenceRepository) ListCountries(ctx context.Context) ([]domain.Country, error) {
	var daos []dao.CountryDAO
	if err := r.db.WithContext(ctx).Find(&daos).Error; err != nil {
		return nil, err
	}

	result := make([]domain.Country, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}
	return result, nil
}

func (r *referenceRepository) ListATCCodes(ctx context.Context, query string, limit int) ([]domain.ATCCode, error) {
	var daos []dao.ATCCodeDAO
	tx := r.db.WithContext(ctx)

	if query != "" {
		tx = tx.Where("code ILIKE ? OR name_ru ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if limit > 0 {
		tx = tx.Limit(limit)
	}

	if err := tx.Find(&daos).Error; err != nil {
		return nil, err
	}

	result := make([]domain.ATCCode, len(daos))
	for i, d := range daos {
		result[i] = d.ToDomain()
	}
	return result, nil
}

func (r *referenceRepository) GetCountryByCode(ctx context.Context, code string) (*domain.Country, error) {
	var d dao.CountryDAO
	if err := r.db.WithContext(ctx).First(&d, "code = ?", strings.ToUpper(code)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Or specific domain error
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}

func (r *referenceRepository) GetATCByCode(ctx context.Context, code string) (*domain.ATCCode, error) {
	var d dao.ATCCodeDAO
	if err := r.db.WithContext(ctx).First(&d, "code = ?", strings.ToUpper(code)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	res := d.ToDomain()
	return &res, nil
}
