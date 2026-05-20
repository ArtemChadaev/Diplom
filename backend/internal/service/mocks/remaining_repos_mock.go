package mocks

import (
	"context"
	"time"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockInventoryRepository is a testify/mock implementation of domain.InventoryRepository.
type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) ListSessions(ctx context.Context, limit, offset int) ([]domain.InventorySession, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]domain.InventorySession), args.Int(1), args.Error(2)
}

func (m *MockInventoryRepository) GetSessionByID(ctx context.Context, id string) (*domain.InventorySession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InventorySession), args.Error(1)
}

func (m *MockInventoryRepository) CreateSession(ctx context.Context, s *domain.InventorySession) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockInventoryRepository) UpdateSession(ctx context.Context, s *domain.InventorySession) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockInventoryRepository) AddCount(ctx context.Context, item *domain.InventoryItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

// MockRecalledBatchRepository
type MockRecalledBatchRepository struct {
	mock.Mock
}

func (m *MockRecalledBatchRepository) List(ctx context.Context, limit, offset int) ([]domain.RecalledBatch, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]domain.RecalledBatch), args.Int(1), args.Error(2)
}

func (m *MockRecalledBatchRepository) GetBySerial(ctx context.Context, serial string) (*domain.RecalledBatch, error) {
	args := m.Called(ctx, serial)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RecalledBatch), args.Error(1)
}

func (m *MockRecalledBatchRepository) Create(ctx context.Context, b *domain.RecalledBatch) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

// MockSystemSettingsRepository
type MockSystemSettingsRepository struct {
	mock.Mock
}

func (m *MockSystemSettingsRepository) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockSystemSettingsRepository) Set(ctx context.Context, key, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockSystemSettingsRepository) List(ctx context.Context) ([]domain.SystemSetting, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.SystemSetting), args.Error(1)
}

// MockEnvironmentLogRepository
type MockEnvironmentLogRepository struct {
	mock.Mock
}

func (m *MockEnvironmentLogRepository) List(ctx context.Context, zoneID string, limit, offset int) ([]domain.EnvironmentLog, int, error) {
	args := m.Called(ctx, zoneID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]domain.EnvironmentLog), args.Int(1), args.Error(2)
}

func (m *MockEnvironmentLogRepository) Create(ctx context.Context, log *domain.EnvironmentLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockEnvironmentLogRepository) ExistsByZoneShiftDate(ctx context.Context, zoneID string, shift string, date time.Time) (bool, error) {
	args := m.Called(ctx, zoneID, shift, date)
	return args.Bool(0), args.Error(1)
}

// MockZoneRepository
type MockZoneRepository struct {
	mock.Mock
}

func (m *MockZoneRepository) List(ctx context.Context) ([]domain.Zone, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Zone), args.Error(1)
}

func (m *MockZoneRepository) GetByID(ctx context.Context, id string) (*domain.Zone, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Zone), args.Error(1)
}

func (m *MockZoneRepository) Create(ctx context.Context, z *domain.Zone) error {
	args := m.Called(ctx, z)
	return args.Error(0)
}

func (m *MockZoneRepository) Update(ctx context.Context, z *domain.Zone) error {
	args := m.Called(ctx, z)
	return args.Error(0)
}

func (m *MockZoneRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

