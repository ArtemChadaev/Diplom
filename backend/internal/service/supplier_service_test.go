package service

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierService_CreateSupplier(t *testing.T) {
	// Scenario: Create a supplier
	//   Given:  A valid supplier and caller role
	//   When:   CreateSupplier is called
	//   Then:   Admin/Manager can create if INN is unique; Others blocked

	mockRepo := new(mocks.MockSupplierRepository)
	svc := NewSupplierService(mockRepo)
	ctx := context.Background()

	sup := &domain.Supplier{INN: "1234567890", Name: "Test Supplier"}

	tests := []struct {
		name      string
		role      domain.UserRole
		mockSetup func()
		wantErr   error
	}{
		{
			name: "admin can create",
			role: domain.RoleAdmin,
			mockSetup: func() {
				mockRepo.On("GetByINN", mock.Anything, sup.INN).Return(nil, nil).Once()
				mockRepo.On("Create", mock.Anything, sup).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "conflict if INN exists",
			role: domain.RoleAdmin,
			mockSetup: func() {
				mockRepo.On("GetByINN", mock.Anything, sup.INN).Return(&domain.Supplier{ID: "existing"}, nil).Once()
			},
			wantErr: domain.ErrConflict,
		},
		{
			name: "storekeeper cannot create",
			role: domain.RoleStorekeeper,
			mockSetup: func() {
				// Early exit
			},
			wantErr: domain.ErrInsufficientPerms,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			res, err := svc.CreateSupplier(ctx, tc.role, sup)
			assert.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				assert.NotNil(t, res)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSupplierService_DeleteSupplier(t *testing.T) {
	// Scenario: Delete a supplier
	//   Given:  An existing supplier and caller role
	//   When:   DeleteSupplier is called
	//   Then:   Only Admin can delete

	mockRepo := new(mocks.MockSupplierRepository)
	svc := NewSupplierService(mockRepo)
	ctx := context.Background()

	id := "sup-123"

	tests := []struct {
		name      string
		role      domain.UserRole
		mockSetup func()
		wantErr   error
	}{
		{
			name: "admin can delete",
			role: domain.RoleAdmin,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Supplier{ID: id}, nil).Once()
				mockRepo.On("Delete", mock.Anything, id).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "manager cannot delete",
			role: domain.RoleWarehouseManager,
			mockSetup: func() {
				// Early exit
			},
			wantErr: domain.ErrInsufficientPerms,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			err := svc.DeleteSupplier(ctx, tc.role, id)
			assert.ErrorIs(t, err, tc.wantErr)
			mockRepo.AssertExpectations(t)
		})
	}
}
