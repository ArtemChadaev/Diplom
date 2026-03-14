package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) domain.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) toDomain(dao *dao.SessionDAO) *domain.RefreshToken {
	if dao == nil {
		return nil
	}
	// For metadata, we just unmarshal JSON to map string any if needed
	// Simplified here since it's just raw JSON bytes usually
	return &domain.RefreshToken{
		ID:        dao.ID,
		UserID:    dao.UserID,
		TokenHash: dao.TokenHash,
		ExpiresAt: dao.ExpiresAt,
		UserAgent: dao.UserAgent,
		IPAddress: dao.IPAddress,
		CreatedAt: dao.CreatedAt,
		RevokedAt: dao.RevokedAt,
	}
}

func (r *sessionRepository) fromDomain(rt *domain.RefreshToken) *dao.SessionDAO {
	if rt == nil {
		return nil
	}
	return &dao.SessionDAO{
		ID:        rt.ID,
		UserID:    rt.UserID,
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt,
		UserAgent: rt.UserAgent,
		IPAddress: rt.IPAddress,
		CreatedAt: rt.CreatedAt,
		RevokedAt: rt.RevokedAt,
	}
}

func (r *sessionRepository) Create(ctx context.Context, rt *domain.RefreshToken) (*domain.RefreshToken, error) {
	d := r.fromDomain(rt)
	if err := r.db.WithContext(ctx).Create(d).Error; err != nil {
		return nil, err
	}
	rt.ID = d.ID
	rt.CreatedAt = d.CreatedAt
	return rt, nil
}

func (r *sessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error) {
	var d dao.SessionDAO
	if err := r.db.WithContext(ctx).First(&d, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}
	return r.toDomain(&d), nil
}

func (r *sessionRepository) FindByTokenHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var d dao.SessionDAO
	if err := r.db.WithContext(ctx).First(&d, "token_hash = ?", hash).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}
	return r.toDomain(&d), nil
}

func (r *sessionRepository) FindActiveByUserID(ctx context.Context, userID int) ([]*domain.RefreshToken, error) {
	var daos []dao.SessionDAO
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("revoked_at IS NULL").
		Where("expires_at > ?", time.Now()).
		Find(&daos).Error

	if err != nil {
		return nil, err
	}

	var results []*domain.RefreshToken
	for _, d := range daos {
		results = append(results, r.toDomain(&d))
	}
	return results, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&dao.SessionDAO{}).Where("id = ?", id).Update("revoked_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrSessionNotFound
	}
	return nil
}

func (r *sessionRepository) RevokeAllForUser(ctx context.Context, userID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&dao.SessionDAO{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&dao.SessionDAO{}).Error
}
