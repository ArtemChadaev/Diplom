package service

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInventoryService_StartSession(t *testing.T) {
	mockRepo := new(mocks.MockInventoryRepository)
	mockProfileRepo := new(mocks.MockEmployeeProfileRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	svc := NewInventoryService(mockRepo, mockProfileRepo, mockProductRepo)
	ctx := context.Background()

	profile := &domain.EmployeeProfile{
		ID:     1,
		UserID: 1,
		GDPTrainingHistory: []domain.GDPTrainingRecord{
			{Date: "2026-01-01", CourseName: "GDP Guidelines", Result: "pass"},
		},
	}

	mockProfileRepo.On("FindByUserID", mock.Anything, 1).Return(profile, nil).Once()
	mockRepo.On("CreateSession", mock.Anything, mock.MatchedBy(func(s *domain.InventorySession) bool {
		return s.Status == domain.InventoryStatusActive
	})).Return(nil).Once()

	res, err := svc.StartSession(ctx, 1, "zone-123")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, domain.InventoryStatusActive, res.Status)
	mockRepo.AssertExpectations(t)
	mockProfileRepo.AssertExpectations(t)
}

func TestInventoryService_FinishSession(t *testing.T) {
	mockRepo := new(mocks.MockInventoryRepository)
	mockProfileRepo := new(mocks.MockEmployeeProfileRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	svc := NewInventoryService(mockRepo, mockProfileRepo, mockProductRepo)
	ctx := context.Background()
	id := "sess-123"

	mockRepo.On("GetSessionByID", mock.Anything, id).Return(&domain.InventorySession{ID: id, Status: domain.InventoryStatusActive}, nil).Once()
	mockRepo.On("UpdateSession", mock.Anything, mock.MatchedBy(func(s *domain.InventorySession) bool {
		return s.Status == domain.InventoryStatusCompleted && s.CompletedAt != nil
	})).Return(nil).Once()

	err := svc.FinishSession(ctx, id)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestInventoryService_GetSession_BlindAudit(t *testing.T) {
	t.Run("active session hides expected quantity", func(t *testing.T) {
		mockRepo := new(mocks.MockInventoryRepository)
		mockProfileRepo := new(mocks.MockEmployeeProfileRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		svc := NewInventoryService(mockRepo, mockProfileRepo, mockProductRepo)
		ctx := context.Background()
		id := "sess-123"

		session := &domain.InventorySession{
			ID:     id,
			Status: domain.InventoryStatusActive,
			Items: []domain.InventoryItem{
				{ID: "item-1", SystemQuantity: 42, PhysicalQuantity: 40},
			},
		}

		mockRepo.On("GetSessionByID", mock.Anything, id).Return(session, nil).Once()

		res, err := svc.GetSession(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 0, res.Items[0].SystemQuantity)
		mockRepo.AssertExpectations(t)
	})

	t.Run("completed session shows expected quantity", func(t *testing.T) {
		mockRepo := new(mocks.MockInventoryRepository)
		mockProfileRepo := new(mocks.MockEmployeeProfileRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		svc := NewInventoryService(mockRepo, mockProfileRepo, mockProductRepo)
		ctx := context.Background()
		id := "sess-123"

		session := &domain.InventorySession{
			ID:     id,
			Status: domain.InventoryStatusCompleted,
			Items: []domain.InventoryItem{
				{ID: "item-1", SystemQuantity: 42, PhysicalQuantity: 40},
			},
		}

		mockRepo.On("GetSessionByID", mock.Anything, id).Return(session, nil).Once()

		res, err := svc.GetSession(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 42, res.Items[0].SystemQuantity)
		mockRepo.AssertExpectations(t)
	})
}

func TestInventoryService_CalculateNetting(t *testing.T) {
	mockRepo := new(mocks.MockInventoryRepository)
	mockProfileRepo := new(mocks.MockEmployeeProfileRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	svc := NewInventoryService(mockRepo, mockProfileRepo, mockProductRepo)
	ctx := context.Background()
	id := "sess-123"

	session := &domain.InventorySession{
		ID:     id,
		Status: domain.InventoryStatusCompleted,
		Items: []domain.InventoryItem{
			{ProductID: "prod-A", SystemQuantity: 10, PhysicalQuantity: 15}, // delta: +5
			{ProductID: "prod-B", SystemQuantity: 20, PhysicalQuantity: 12}, // delta: -8
			{ProductID: "prod-A", SystemQuantity: 5, PhysicalQuantity: 5},   // aggregate prod-A -> total delta: +5
		},
	}

	prodA := &domain.Product{ID: "prod-A", Name: "Product A", ATCCode: "A02BC01"}
	prodB := &domain.Product{ID: "prod-B", Name: "Product B", ATCCode: "A02BC02"}

	mockRepo.On("GetSessionByID", mock.Anything, id).Return(session, nil).Once()
	mockProductRepo.On("GetByID", mock.Anything, "prod-A").Return(prodA, nil).Once()
	mockProductRepo.On("GetByID", mock.Anything, "prod-B").Return(prodB, nil).Once()

	res, err := svc.CalculateNetting(ctx, id)
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	// Since they both have ATC code "A02BC01"/"A02BC02", both will have ATCGroup "A02"
	// Order could depend on sorting by ATCGroup, then stable ordering.
	assert.Equal(t, "A02", res[0].ATCGroup)
	assert.Equal(t, "A02", res[1].ATCGroup)
	mockRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}
