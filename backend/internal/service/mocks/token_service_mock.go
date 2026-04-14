package mocks

import (
	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockTokenService is a testify/mock implementation of domain.TokenService.
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateAccessToken(user *domain.User, sessionID uuid.UUID) (string, error) {
	args := m.Called(user, sessionID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken() (string, string, error) {
	args := m.Called()
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockTokenService) ParseAccessToken(tokenStr string) (*domain.AccessTokenClaims, error) {
	args := m.Called(tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AccessTokenClaims), args.Error(1)
}

func (m *MockTokenService) HashToken(raw string) string {
	args := m.Called(raw)
	return args.String(0)
}
