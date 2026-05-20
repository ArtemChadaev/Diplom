package mocks

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository is a testify/mock implementation of domain.OrderRepository.
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) List(ctx context.Context, limit, offset int) ([]domain.Order, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]domain.Order), args.Int(1), args.Error(2)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) Create(ctx context.Context, o *domain.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *MockOrderRepository) Update(ctx context.Context, o *domain.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) GetMonthlyTurnover(ctx context.Context, productID string) (int, error) {
	args := m.Called(ctx, productID)
	return args.Int(0), args.Error(1)
}

