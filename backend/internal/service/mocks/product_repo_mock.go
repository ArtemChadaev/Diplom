package mocks

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a testify/mock implementation of domain.ProductRepository.
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]domain.Product), args.Int(1), args.Error(2)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Create(ctx context.Context, p *domain.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductRepository) Update(ctx context.Context, p *domain.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}
