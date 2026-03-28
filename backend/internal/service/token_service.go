package service

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
)

type tokenService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewTokenService(secret string, accessTTL, refreshTTL time.Duration) domain.TokenService {
	return &tokenService{
		secretKey:       []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *tokenService) GenerateAccessToken(user *domain.User, sessionID uuid.UUID) (string, error) {
	claims := &domain.AccessTokenClaims{
		UserID:    user.ID,
		SessionID: sessionID,
		Role:      user.Role,
		Email:     user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        sessionID.String(), // jti
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *tokenService) GenerateRefreshToken() (raw string, hash string, err error) {
	// Generate random 32 byte string (we'll use UUID v4 without hyphens for simplicity/speed)
	rawUUID := uuid.New()
	raw = hex.EncodeToString(rawUUID[:])

	hashStr := s.HashToken(raw)
	return raw, hashStr, nil
}

func (s *tokenService) HashToken(raw string) string {
	b := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(b[:])
}

func (s *tokenService) ParseAccessToken(tokenStr string) (*domain.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, domain.ErrInvalidToken
}
