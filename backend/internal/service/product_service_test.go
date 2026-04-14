package service

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_CreateProduct(t *testing.T) {
	// Scenario: Create a product
	//   Given:  A valid product and caller role
	//   When:   CreateProduct is called
	//   Then:   Admin/Manager can create if SKU is unique; Others blocked

	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	p := &domain.Product{SKU: "SKU-123", Name: "Test Product"}

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
				mockRepo.On("GetBySKU", mock.Anything, p.SKU).Return(nil, nil).Once()
				mockRepo.On("Create", mock.Anything, p).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "conflict if SKU exists",
			role: domain.RoleAdmin,
			mockSetup: func() {
				mockRepo.On("GetBySKU", mock.Anything, p.SKU).Return(&domain.Product{ID: "existing"}, nil).Once()
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
			res, err := svc.CreateProduct(ctx, tc.role, p)
			assert.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				assert.NotNil(t, res)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	// Scenario: Delete a product
	//   Given:  An existing product and caller role
	//   When:   DeleteProduct is called
	//   Then:   Only Admin can delete

	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	id := "prod-123"

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
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Product{ID: id}, nil).Once()
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
			err := svc.DeleteProduct(ctx, tc.role, id)
			assert.ErrorIs(t, err, tc.wantErr)
			mockRepo.AssertExpectations(t)
		})
	}
}
