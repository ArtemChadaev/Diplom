package mocks

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockClaimRepository is a testify/mock implementation of domain.ClaimRepository.
type MockClaimRepository struct {
	mock.Mock
}

func (m *MockClaimRepository) List(ctx context.Context, limit, offset int) ([]domain.Claim, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]domain.Claim), args.Int(1), args.Error(2)
}

func (m *MockClaimRepository) GetByID(ctx context.Context, id string) (*domain.Claim, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Claim), args.Error(1)
}

func (m *MockClaimRepository) Create(ctx context.Context, c *domain.Claim) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockClaimRepository) Update(ctx context.Context, c *domain.Claim) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}
