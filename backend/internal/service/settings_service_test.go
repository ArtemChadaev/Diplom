package service

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSystemSettingsService_UpdateSetting(t *testing.T) {
	mockRepo := new(mocks.MockSystemSettingsRepository)
	svc := NewSystemSettingsService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		role      domain.UserRole
		mockSetup func()
		wantErr   error
	}{
		{
			name: "admin can update setting",
			role: domain.RoleAdmin,
			mockSetup: func() {
				mockRepo.On("Set", mock.Anything, "maintenance_mode", "true").Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "storekeeper cannot update setting",
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
			err := svc.UpdateSetting(ctx, tc.role, "maintenance_mode", "true")
			assert.ErrorIs(t, err, tc.wantErr)
			mockRepo.AssertExpectations(t)
		})
	}
}
