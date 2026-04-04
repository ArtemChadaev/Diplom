package domain

import "context"

// OTPCode represents a one-time password code from Valkey
type OTPCode struct {
	UserID   int
	CodeHash string
	Attempts int
}

const (
	OTPMaxAttempts = 5
	OTPTTLSeconds  = 600 // 10 minutes
)

// OTPRepository represents Valkey storage for OTP
type OTPRepository interface {
	Store(ctx context.Context, userID int, codeHash string) error
	Get(ctx context.Context, userID int) (*OTPCode, error)
	IncrAttempts(ctx context.Context, userID int) error
	Delete(ctx context.Context, userID int) error
}
