package mocks

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockOTPRepository is a testify/mock implementation of domain.OTPRepository.
type MockOTPRepository struct {
	mock.Mock
}

func (m *MockOTPRepository) Store(ctx context.Context, userID int, codeHash string) error {
	args := m.Called(ctx, userID, codeHash)
	return args.Error(0)
}

func (m *MockOTPRepository) Get(ctx context.Context, userID int) (*domain.OTPCode, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.OTPCode), args.Error(1)
}

func (m *MockOTPRepository) IncrAttempts(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockOTPRepository) Delete(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
