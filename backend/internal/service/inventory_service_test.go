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
	svc := NewInventoryService(mockRepo)
	ctx := context.Background()

	mockRepo.On("CreateSession", mock.Anything, mock.MatchedBy(func(s *domain.InventorySession) bool {
		return s.Status == domain.InventoryStatusActive
	})).Return(nil).Once()

	res, err := svc.StartSession(ctx, 1, "zone-123")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, domain.InventoryStatusActive, res.Status)
	mockRepo.AssertExpectations(t)
}

func TestInventoryService_FinishSession(t *testing.T) {
	mockRepo := new(mocks.MockInventoryRepository)
	svc := NewInventoryService(mockRepo)
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
