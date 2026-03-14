package domain

import "github.com/golang-jwt/jwt/v5"

// AccessTokenClaims — payload embedded in the JWT.
type AccessTokenClaims struct {
	UserID int      `json:"sub"`
	Role   UserRole `json:"role"`
	Email  string   `json:"email,omitempty"`
	jwt.RegisteredClaims
}
