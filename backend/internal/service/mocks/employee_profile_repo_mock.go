package mocks

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeProfileRepository is a testify/mock implementation of domain.EmployeeProfileRepository.
type MockEmployeeProfileRepository struct {
	mock.Mock
}

func (m *MockEmployeeProfileRepository) FindByUserID(ctx context.Context, userID int) (*domain.EmployeeProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EmployeeProfile), args.Error(1)
}

func (m *MockEmployeeProfileRepository) FindByID(ctx context.Context, id int) (*domain.EmployeeProfile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EmployeeProfile), args.Error(1)
}

func (m *MockEmployeeProfileRepository) Update(ctx context.Context, id int, input domain.UpdateEmployeeProfileInput) (*domain.EmployeeProfile, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EmployeeProfile), args.Error(1)
}

func (m *MockEmployeeProfileRepository) Create(ctx context.Context, input domain.CreateEmployeeProfileInput) (*domain.EmployeeProfile, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EmployeeProfile), args.Error(1)
}

func (m *MockEmployeeProfileRepository) List(ctx context.Context, limit, offset int) ([]domain.EmployeeProfile, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.EmployeeProfile), args.Error(1)
}
