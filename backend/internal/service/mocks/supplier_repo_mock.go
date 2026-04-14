package mocks

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockSupplierRepository is a testify/mock implementation of domain.SupplierRepository.
type MockSupplierRepository struct {
	mock.Mock
}

func (m *MockSupplierRepository) List(ctx context.Context, limit, offset int) ([]domain.Supplier, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]domain.Supplier), args.Int(1), args.Error(2)
}

func (m *MockSupplierRepository) GetByID(ctx context.Context, id string) (*domain.Supplier, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Supplier), args.Error(1)
}

func (m *MockSupplierRepository) GetByINN(ctx context.Context, inn string) (*domain.Supplier, error) {
	args := m.Called(ctx, inn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Supplier), args.Error(1)
}

func (m *MockSupplierRepository) Create(ctx context.Context, s *domain.Supplier) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockSupplierRepository) Update(ctx context.Context, s *domain.Supplier) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockSupplierRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
