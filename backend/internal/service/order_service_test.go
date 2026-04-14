package service

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService_CreateOrder(t *testing.T) {
	// Scenario: Create an order
	//   Given:  A valid order object
	//   When:   CreateOrder is called
	//   Then:   Status is set to new, priority is defaulted to 1 if zero, and record is saved

	mockRepo := new(mocks.MockOrderRepository)
	svc := NewOrderService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name         string
		order        *domain.Order
		mockSetup    func(o *domain.Order)
		wantStatus   domain.OrderStatus
		wantPriority int
	}{
		{
			name: "create order with priority",
			order: &domain.Order{
				OrderNumber: "ORD-001",
				Priority:    2,
			},
			mockSetup: func(o *domain.Order) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(order *domain.Order) bool {
					return order.Status == domain.OrderStatusNew && order.Priority == 2
				})).Return(nil).Once()
			},
			wantStatus:   domain.OrderStatusNew,
			wantPriority: 2,
		},
		{
			name: "create order with default priority",
			order: &domain.Order{
				OrderNumber: "ORD-002",
				Priority:    0,
			},
			mockSetup: func(o *domain.Order) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(order *domain.Order) bool {
					return order.Status == domain.OrderStatusNew && order.Priority == 1
				})).Return(nil).Once()
			},
			wantStatus:   domain.OrderStatusNew,
			wantPriority: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup(tc.order)
			res, err := svc.CreateOrder(ctx, tc.order)
			assert.NoError(t, err)
			assert.NotNil(t, res)
			assert.Equal(t, tc.wantStatus, res.Status)
			assert.Equal(t, tc.wantPriority, res.Priority)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_UpdateStatus(t *testing.T) {
	// Scenario: Update order status
	//   Given:  An existing order and a target status
	//   When:   UpdateStatus is called with role-based checks
	//   Then:   Status is updated if roles allow

	mockRepo := new(mocks.MockOrderRepository)
	svc := NewOrderService(mockRepo)
	ctx := context.Background()

	id := "ord-123"

	tests := []struct {
		name      string
		role      domain.UserRole
		status    domain.OrderStatus
		mockSetup func()
		wantErr   error
	}{
		{
			name:   "storekeeper can start assembling",
			role:   domain.RoleStorekeeper,
			status: domain.OrderStatusAssembling,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Order{ID: id, Status: domain.OrderStatusNew}, nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(o *domain.Order) bool {
					return o.Status == domain.OrderStatusAssembling
				})).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:   "pharmacist cannot start assembling",
			role:   domain.RolePharmacist,
			status: domain.OrderStatusAssembling,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Order{ID: id, Status: domain.OrderStatusNew}, nil).Once()
			},
			wantErr: domain.ErrInsufficientPerms,
		},
		{
			name:   "manager can ship",
			role:   domain.RoleWarehouseManager,
			status: domain.OrderStatusShipped,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Order{ID: id, Status: domain.OrderStatusAssembled}, nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(o *domain.Order) bool {
					return o.Status == domain.OrderStatusShipped
				})).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:   "storekeeper cannot ship",
			role:   domain.RoleStorekeeper,
			status: domain.OrderStatusShipped,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, id).Return(&domain.Order{ID: id, Status: domain.OrderStatusAssembled}, nil).Once()
			},
			wantErr: domain.ErrInsufficientPerms,
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
