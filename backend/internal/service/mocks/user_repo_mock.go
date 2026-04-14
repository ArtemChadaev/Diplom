package mocks

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a testify/mock implementation of domain.UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	args := m.Called(ctx, googleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	args := m.Called(ctx, telegramID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateRole(ctx context.Context, userID int, role domain.UserRole) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *MockUserRepository) LinkGoogle(ctx context.Context, userID int, googleID string) error {
	args := m.Called(ctx, userID, googleID)
	return args.Error(0)
}

func (m *MockUserRepository) LinkTelegram(ctx context.Context, userID int, telegramID int64) error {
	args := m.Called(ctx, userID, telegramID)
	return args.Error(0)
}

func (m *MockUserRepository) SetNsPvAccess(ctx context.Context, userID int, access bool) error {
	args := m.Called(ctx, userID, access)
	return args.Error(0)
}

func (m *MockUserRepository) SetBlocked(ctx context.Context, userID int, blocked bool) error {
	args := m.Called(ctx, userID, blocked)
	return args.Error(0)
}

func (m *MockUserRepository) FindProfileByUserID(ctx context.Context, userID int) (*domain.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserProfile), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, filter domain.UserListFilter) ([]*domain.UserProfile, int, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.UserProfile), args.Int(1), args.Error(2)
}
