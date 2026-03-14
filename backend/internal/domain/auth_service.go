package domain

import (
	"context"

	"github.com/google/uuid"
)

type RegisterInput struct {
	Login    string
	Email    string
	Password string
}

type TelegramAuthData struct {
	ID        int64
	FirstName string
	Username  string
	PhotoURL  string
	AuthDate  int64
	Hash      string
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"` // in seconds
	RefreshToken string `json:"-"`          // sent in cookie, not body
}

type SessionMeta struct {
	UserAgent string
	IPAddress string
}

type AuthService interface {
	// Email / Password
	Register(ctx context.Context, req RegisterInput) (*User, error)
	LoginWithPassword(ctx context.Context, login, password, userAgent, ip string) (*TokenPair, error)

	// Social auth
	LoginWithGoogle(ctx context.Context, idToken, userAgent, ip string) (*TokenPair, error)
	LoginWithTelegram(ctx context.Context, data TelegramAuthData, userAgent, ip string) (*TokenPair, error)

	// Session management
	RefreshTokens(ctx context.Context, oldRefreshToken string, meta SessionMeta) (*TokenPair, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID, callerID int, callerRole UserRole) error

	// Admin actions
	VerifyUser(ctx context.Context, adminID, targetUserID int) error
	AssignRole(ctx context.Context, adminID, targetUserID int, role UserRole) error
}
