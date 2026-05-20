package service

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClaimService_CreateClaim(t *testing.T) {
	// Scenario: Create a claim
	//   Given:  A valid claim object
	//   When:   CreateClaim is called
	//   Then:   Status is set to new and record is saved

	mockRepo := new(mocks.MockClaimRepository)
	mockBatchRepo := new(mocks.MockBatchRepository)
	svc := NewClaimService(mockRepo, mockBatchRepo)
	ctx := context.Background()

	c := &domain.Claim{Title: "Broken Package", Description: "Package arrived crushed"}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(claim *domain.Claim) bool {
		return claim.Status == domain.ClaimStatusNew
	})).Return(nil).Once()

	res, err := svc.CreateClaim(ctx, c)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, domain.ClaimStatusNew, res.Status)
	mockRepo.AssertExpectations(t)
}

func TestClaimService_UpdateStatus(t *testing.T) {
	// Scenario: Update claim status
	//   Given:  An existing claim and a target status
	//   When:   UpdateStatus is called with role-based checks
	//   Then:   Only Admin or WarehouseManager can move to resolved/rejected

	mockRepo := new(mocks.MockClaimRepository)
	mockBatchRepo := new(mocks.MockBatchRepository)
	svc := NewClaimService(mockRepo, mockBatchRepo)
	ctx := context.Background()

	id := "clm-123"

	tests := []struct {
		name      string
		role      domain.UserRole
		status    domain.ClaimStatus
		mockSetup func()
		wantErr   error
	}{
		{
			name:   "admin can resolve",
			role:   domain.RoleAdmin,
			status: domain.ClaimStatusResolved,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Claim{ID: id, Status: domain.ClaimStatusNew}, nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *domain.Claim) bool {
					return c.Status == domain.ClaimStatusResolved
				})).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:   "storekeeper cannot resolve",
			role:   domain.RoleStorekeeper,
			status: domain.ClaimStatusResolved,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Claim{ID: id, Status: domain.ClaimStatusNew}, nil).Once()
			},
			wantErr: domain.ErrInsufficientPerms,
		},
		{
			name:   "not found",
			role:   domain.RoleAdmin,
			status: domain.ClaimStatusResolved,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(nil, nil).Once()
			},
			wantErr: domain.ErrClaimNotFound,
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
