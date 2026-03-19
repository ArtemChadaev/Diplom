package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AccessTokenClaims — payload embedded in the JWT.
type AccessTokenClaims struct {
	UserID    int       `json:"sub"`
	SessionID uuid.UUID `json:"jti"`
	Role      UserRole  `json:"role"`
	Email     string    `json:"email,omitempty"`
	jwt.RegisteredClaims
}
