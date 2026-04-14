package service

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBatchService_UpdateStatus(t *testing.T) {
	// Scenario: Update batch status
	//   Given:  An existing batch and a target status
	//   When:   UpdateStatus is called with role-based checks
	//   Then:   Admin or WarehouseManager can update; Pharmacist or Storekeeper cannot

	mockRepo := new(mocks.MockBatchRepository)
	svc := NewBatchService(mockRepo)
	ctx := context.Background()

	id := "bat-123"

	tests := []struct {
		name      string
		role      domain.UserRole
		status    domain.BatchStatus
		mockSetup func()
		wantErr   error
	}{
		{
			name:   "admin can update status",
			role:   domain.RoleAdmin,
			status: domain.BatchStatusAvailable,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Batch{ID: id, Status: domain.BatchStatusQuarantine}, nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(b *domain.Batch) bool {
					return b.Status == domain.BatchStatusAvailable
				})).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:   "storekeeper cannot update status",
			role:   domain.RoleStorekeeper,
			status: domain.BatchStatusAvailable,
			mockSetup: func() {},
			wantErr:   domain.ErrInsufficientPerms,
		},
		{
			name:   "pharmacist cannot update status",
			role:   domain.RolePharmacist,
			status: domain.BatchStatusAvailable,
			mockSetup: func() {},
			wantErr:   domain.ErrInsufficientPerms,
		},
		{
			name:   "not found",
			role:   domain.RoleAdmin,
			status: domain.BatchStatusAvailable,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(nil, nil).Once()
			},
			wantErr: domain.ErrBatchNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			err := svc.UpdateStatus(ctx, tc.role, id, tc.status)
			assert.ErrorIs(t, err, tc.wantErr)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBatchService_TransferBatch(t *testing.T) {
	// Scenario: Transfer batch to different zone
	//   Given:  An existing batch and target zone
	//   When:   TransferBatch is called
	//   Then:   Pharmacist is blocked; Admin, Manager, and Storekeeper are allowed

	mockRepo := new(mocks.MockBatchRepository)
	svc := NewBatchService(mockRepo)
	ctx := context.Background()

	id := "bat-456"
	targetZone := "zone-789"

	tests := []struct {
		name      string
		role      domain.UserRole
		mockSetup func()
		wantErr   error
	}{
		{
			name: "storekeeper can transfer",
			role: domain.RoleStorekeeper,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Batch{ID: id, ZoneID: nil}, nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(b *domain.Batch) bool {
					return *b.ZoneID == targetZone
				})).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "pharmacist cannot transfer (security check)",
			role: domain.RolePharmacist,
			mockSetup: func() {
				// No repository calls expected due to early exit
			},
			wantErr: domain.ErrInsufficientPerms,
		},
		{
			name: "manager can transfer",
			role: domain.RoleWarehouseManager,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Batch{ID: id, ZoneID: nil}, nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(b *domain.Batch) bool {
					return *b.ZoneID == targetZone
				})).Return(nil).Once()
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			err := svc.TransferBatch(ctx, tc.role, id, targetZone)
			assert.ErrorIs(t, err, tc.wantErr)
			mockRepo.AssertExpectations(t)
		})
	}
}
