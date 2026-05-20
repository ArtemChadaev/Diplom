package mocks

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockBatchRepository is a testify/mock implementation of domain.BatchRepository.
type MockBatchRepository struct {
	mock.Mock
}

func (m *MockBatchRepository) List(ctx context.Context, filter domain.BatchFilter) ([]domain.Batch, int, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]domain.Batch), args.Int(1), args.Error(2)
}

func (m *MockBatchRepository) GetByID(ctx context.Context, id string) (*domain.Batch, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Batch), args.Error(1)
}

func (m *MockBatchRepository) Update(ctx context.Context, b *domain.Batch) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBatchRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBatchRepository) BlockAllByProductID(ctx context.Context, productID string) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}

func (m *MockBatchRepository) ListAvailableSorted(ctx context.Context, productID string) ([]domain.Batch, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Batch), args.Error(1)
}

func (m *MockBatchRepository) GetTotalStock(ctx context.Context, productID string) (int, error) {
	args := m.Called(ctx, productID)
	return args.Int(0), args.Error(1)
}

