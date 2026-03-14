package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	UserAgent string
	IPAddress string
	Metadata  map[string]any
	CreatedAt time.Time
	RevokedAt *time.Time
}

type SessionRepository interface {
	Create(ctx context.Context, rt *RefreshToken) (*RefreshToken, error)
	FindByID(ctx context.Context, id uuid.UUID) (*RefreshToken, error)
	FindByTokenHash(ctx context.Context, hash string) (*RefreshToken, error)
	FindActiveByUserID(ctx context.Context, userID int) ([]*RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID int) error
	DeleteExpired(ctx context.Context) error // for a cleanup cron/job
}
