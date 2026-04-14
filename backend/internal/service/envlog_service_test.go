package service

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEnvironmentLogService_RecordLogs(t *testing.T) {
	mockRepo := new(mocks.MockEnvironmentLogRepository)
	svc := NewEnvironmentLogService(mockRepo)
	ctx := context.Background()

	logs := []domain.EnvironmentLog{
		{ZoneID: "Z1", Temperature: 22.5, Humidity: 45.0},
		{ZoneID: "Z2", Temperature: 18.0, Humidity: 50.0},
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(l *domain.EnvironmentLog) bool {
		return l.RecordedBy == 1 && !l.RecordedAt.IsZero()
	})).Return(nil).Times(2)

	err := svc.RecordLogs(ctx, 1, logs)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
