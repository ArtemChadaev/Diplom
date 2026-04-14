package service

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRecalledBatchService_CheckBatch(t *testing.T) {
	mockRepo := new(mocks.MockRecalledBatchRepository)
	svc := NewRecalledBatchService(mockRepo)
	ctx := context.Background()
	serial := "SN-999"

	tests := []struct {
		name       string
		mockSetup  func()
		wantFound  bool
		wantRecall *domain.RecalledBatch
	}{
		{
			name: "batch found in recall list",
			mockSetup: func() {
				mockRepo.On("GetBySerial", mock.Anything, serial).Return(&domain.RecalledBatch{SerialNumber: serial, RecallReason: "Dangerous"}, nil).Once()
			},
			wantFound:  true,
			wantRecall: &domain.RecalledBatch{SerialNumber: serial, RecallReason: "Dangerous"},
		},
		{
			name: "batch not found",
			mockSetup: func() {
				mockRepo.On("GetBySerial", mock.Anything, serial).Return(nil, nil).Once()
			},
			wantFound:  false,
			wantRecall: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			found, recall, err := svc.CheckBatch(ctx, serial)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantFound, found)
			assert.Equal(t, tc.wantRecall, recall)
			mockRepo.AssertExpectations(t)
		})
	}
}
