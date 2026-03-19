package domain

import "github.com/google/uuid"

type TokenService interface {
	GenerateAccessToken(user *User, sessionID uuid.UUID) (string, error)
	GenerateRefreshToken() (raw string, hash string, err error) // raw → cookie, hash → DB
	ParseAccessToken(tokenStr string) (*AccessTokenClaims, error)
	HashToken(raw string) string
}
