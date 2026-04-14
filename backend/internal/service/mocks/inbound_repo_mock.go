package mocks

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockInboundRepository is a testify/mock implementation of domain.InboundRepository.
type MockInboundRepository struct {
	mock.Mock
}

func (m *MockInboundRepository) List(ctx context.Context, limit, offset int) ([]domain.InboundReceiving, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]domain.InboundReceiving), args.Int(1), args.Error(2)
}

func (m *MockInboundRepository) GetByID(ctx context.Context, id string) (*domain.InboundReceiving, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InboundReceiving), args.Error(1)
}

func (m *MockInboundRepository) Create(ctx context.Context, i *domain.InboundReceiving) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}

func (m *MockInboundRepository) Update(ctx context.Context, i *domain.InboundReceiving) error {
	args := m.Called(ctx, i)
	return args.Error(0)
}

func (m *MockInboundRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInboundRepository) AddItems(ctx context.Context, inboundID string, items []domain.InboundItem) error {
	args := m.Called(ctx, inboundID, items)
	return args.Error(0)
}

func (m *MockInboundRepository) UpdateItem(ctx context.Context, item *domain.InboundItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockInboundRepository) RemoveItem(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
