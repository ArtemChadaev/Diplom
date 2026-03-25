package domain

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ima/diplom-backend/internal/pkg/logger"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrLoginTaken        = errors.New("login already taken")
	ErrEmailTaken        = errors.New("email already taken")
	ErrInvalidCreds      = errors.New("invalid login or password")
	ErrUserUnverified    = errors.New("account pending administrator verification")
	ErrUserBlocked       = errors.New("account is blocked")
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidToken      = errors.New("invalid token")
	ErrSessionNotFound   = errors.New("session not found or terminated")
	ErrInvalidTelegram   = errors.New("invalid telegram auth data")
	ErrInsufficientPerms = errors.New("insufficient permissions for this operation")
	ErrEmployeeProfileNotFound = errors.New("employee profile not found")
)


// AppError is a structured application error that carries
// a human-readable message plus arbitrary key-value context.
// It wraps an underlying cause so errors.Is / errors.As keep working.
type AppError struct {
	// Code is a machine-readable error code, e.g. "user_not_found".
	Code string
	// Message is the human-readable message (may be shown to callers).
	Message string
	// Attrs are structured slog attributes attached to this error.
	Attrs []slog.Attr
	// cause is the wrapped error.
	cause error
}

// NewAppError creates a new AppError.
// Example:
//   domain.NewAppError("profile_not_found", "employee profile not found", ErrUserNotFound,
//       slog.Int("user_id", uid))
func NewAppError(code, message string, cause error, attrs ...slog.Attr) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Attrs:   attrs,
		cause:   cause,
	}
}

func (e *AppError) Error() string { return e.Message }

// Unwrap makes errors.Is / errors.As see through AppError to the underlying cause.
func (e *AppError) Unwrap() error { return e.cause }

// LogError logs the AppError via the context-aware logger.
func (e *AppError) LogError(ctx context.Context) {
	args := make([]any, 0, len(e.Attrs)+2)
	args = append(args, slog.String("error_code", e.Code))
	for _, a := range e.Attrs {
		args = append(args, a)
	}
	logger.FromContext(ctx).Error(e.Message, args...)
}

