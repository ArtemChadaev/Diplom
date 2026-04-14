package service

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInboundService_CreateInbound(t *testing.T) {
	// Scenario: Create inbound record
	//   Given:  A valid inbound record and caller role
	//   When:   CreateInbound is called
	//   Then:   If caller is not Storekeeper or Pharmacist, record is created with status draft

	mockRepo := new(mocks.MockInboundRepository)
	svc := NewInboundService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name       string
		role       domain.UserRole
		inbound    *domain.InboundReceiving
		mockSetup  func()
		wantErr    error
		wantStatus domain.InboundStatus
	}{
		{
			name: "admin can create",
			role: domain.RoleAdmin,
			inbound: &domain.InboundReceiving{
				InvoiceNumber: "INV-001",
			},
			mockSetup: func() {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(i *domain.InboundReceiving) bool {
					return i.Status == domain.InboundStatusDraft
				})).Return(nil).Once()
			},
			wantErr:    nil,
			wantStatus: domain.InboundStatusDraft,
		},
		{
			name: "storekeeper cannot create",
			role: domain.RoleStorekeeper,
			inbound: &domain.InboundReceiving{
				InvoiceNumber: "INV-002",
			},
			mockSetup: func() {},
			wantErr:   domain.ErrInsufficientPerms,
		},
		{
			name: "pharmacist cannot create",
			role: domain.RolePharmacist,
			inbound: &domain.InboundReceiving{
				InvoiceNumber: "INV-003",
			},
			mockSetup: func() {},
			wantErr:   domain.ErrInsufficientPerms,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			res, err := svc.CreateInbound(ctx, tc.role, tc.inbound)
			assert.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				assert.NotNil(t, res)
				assert.Equal(t, tc.wantStatus, res.Status)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestInboundService_UpdateStatus(t *testing.T) {
	// Scenario: Update inbound status
	//   Given:  An existing inbound and a target status
	//   When:   UpdateStatus is called with role-based checks
	//   Then:   Status is updated if permissions and logic flow allow

	mockRepo := new(mocks.MockInboundRepository)
	svc := NewInboundService(mockRepo)
	ctx := context.Background()

	id := "inb-123"

	tests := []struct {
		name      string
		role      domain.UserRole
		status    domain.InboundStatus
		mockSetup func()
		wantErr   error
	}{
		{
			name:   "manager can set completed",
			role:   domain.RoleWarehouseManager,
			status: domain.InboundStatusCompleted,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.InboundReceiving{ID: id, Status: domain.InboundStatusReceived}, nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(i *domain.InboundReceiving) bool {
					return i.Status == domain.InboundStatusCompleted
				})).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:   "pharmacist cannot set received",
			role:   domain.RolePharmacist,
			status: domain.InboundStatusReceived,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.InboundReceiving{ID: id, Status: domain.InboundStatusDraft}, nil).Once()
			},
			wantErr: domain.ErrInsufficientPerms,
		},
		{
			name:   "not found",
			role:   domain.RoleAdmin,
			status: domain.InboundStatusCompleted,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(nil, nil).Once()
			},
			wantErr: domain.ErrInboundNotFound,
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
