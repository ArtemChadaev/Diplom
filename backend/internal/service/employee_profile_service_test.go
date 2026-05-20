package service_test

import (
	"context"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newEmployeeProfileService(repo domain.EmployeeProfileRepository) domain.EmployeeProfileService {
	return service.NewEmployeeProfileService(repo)
}

// ---------------------------------------------------------------------------
// GetProfile
// ---------------------------------------------------------------------------

// Scenario: Admin-only profile retrieval
//
//	Covers: admin success, non-admin rejection
func TestEmployeeProfileService_GetProfile(t *testing.T) {
	ctx := context.Background()
	targetUserID := 5
	profile := &domain.EmployeeProfile{ID: 10, UserID: uint(targetUserID), FullName: "Jane Doe"}

	tests := []struct {
		name       string
		callerRole domain.UserRole
		repoResult *domain.EmployeeProfile
		repoErr    error
		wantErr    error
		wantCall   bool
	}{
		{
			name:       "admin fetches profile successfully",
			callerRole: domain.RoleAdmin,
			repoResult: profile,
			repoErr:    nil,
			wantErr:    nil,
			wantCall:   true,
		},
		{
			name:       "non-admin (pharmacist) is rejected",
			callerRole: domain.RolePharmacist,
			wantErr:    domain.ErrInsufficientPerms,
			wantCall:   false,
		},
		{
			name:       "non-admin (warehouse_manager) is rejected",
			callerRole: domain.RoleWarehouseManager,
			wantErr:    domain.ErrInsufficientPerms,
			wantCall:   false,
		},
		{
			name:       "admin, profile not found in DB",
			callerRole: domain.RoleAdmin,
			repoErr:    domain.ErrEmployeeProfileNotFound,
			wantErr:    domain.ErrEmployeeProfileNotFound,
			wantCall:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mocks.MockEmployeeProfileRepository{}
			if tc.wantCall {
				repo.On("FindByUserID", ctx, targetUserID).Return(tc.repoResult, tc.repoErr)
			}

			svc := newEmployeeProfileService(repo)
			got, err := svc.GetProfile(ctx, 1, tc.callerRole, targetUserID)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, profile, got)
			}
			repo.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// UpdateProfile
// ---------------------------------------------------------------------------

// Scenario: Admin-only profile update with profile existence check
//
//	Covers: admin success, non-admin rejection, profile not found, employee_code validation
func TestEmployeeProfileService_UpdateProfile(t *testing.T) {
	ctx := context.Background()
	callerID := 1
	targetUserID := 5
	profileID := 10
	profile := &domain.EmployeeProfile{ID: uint(profileID), UserID: uint(targetUserID)}
	updatedProfile := &domain.EmployeeProfile{ID: uint(profileID), UserID: uint(targetUserID), FullName: "Updated Name"}
	input := domain.UpdateEmployeeProfileInput{}

	codeValid := "IT-001"
	inputWithValidCode := domain.UpdateEmployeeProfileInput{EmployeeCode: &codeValid}
	profileWithCode := &domain.EmployeeProfile{ID: uint(profileID), UserID: uint(targetUserID), EmployeeCode: codeValid}

	codeInvalid := "invalid-123"
	inputWithInvalidCode := domain.UpdateEmployeeProfileInput{EmployeeCode: &codeInvalid}

	tests := []struct {
		name       string
		callerRole domain.UserRole
		input      domain.UpdateEmployeeProfileInput
		wantErr    error
		wantResult *domain.EmployeeProfile
		setupRepo  func(*mocks.MockEmployeeProfileRepository)
	}{
		{
			name:       "admin updates profile successfully",
			callerRole: domain.RoleAdmin,
			input:      input,
			wantErr:    nil,
			wantResult: updatedProfile,
			setupRepo: func(r *mocks.MockEmployeeProfileRepository) {
				r.On("FindByUserID", ctx, targetUserID).Return(profile, nil)
				r.On("Update", ctx, profileID, input).Return(updatedProfile, nil)
			},
		},
		{
			name:       "admin updates profile with valid employee_code successfully",
			callerRole: domain.RoleAdmin,
			input:      inputWithValidCode,
			wantErr:    nil,
			wantResult: profileWithCode,
			setupRepo: func(r *mocks.MockEmployeeProfileRepository) {
				r.On("FindByUserID", ctx, targetUserID).Return(profile, nil)
				r.On("Update", ctx, profileID, inputWithValidCode).Return(profileWithCode, nil)
			},
		},
		{
			name:       "admin tries invalid employee_code format → ErrInvalidEmployeeCode",
			callerRole: domain.RoleAdmin,
			input:      inputWithInvalidCode,
			wantErr:    domain.ErrInvalidEmployeeCode,
			setupRepo:  func(r *mocks.MockEmployeeProfileRepository) {},
		},
		{
			name:       "non-admin is rejected before any repo call",
			callerRole: domain.RoleStorekeeper,
			input:      input,
			wantErr:    domain.ErrInsufficientPerms,
			setupRepo:  func(r *mocks.MockEmployeeProfileRepository) {},
		},
		{
			name:       "admin, profile not found → create new profile",
			callerRole: domain.RoleAdmin,
			input:      input,
			wantErr:    nil,
			wantResult: updatedProfile,
			setupRepo: func(r *mocks.MockEmployeeProfileRepository) {
				r.On("FindByUserID", ctx, targetUserID).Return(nil, domain.ErrEmployeeProfileNotFound)
				r.On("Create", ctx, &domain.EmployeeProfile{UserID: uint(targetUserID)}).Return(updatedProfile, nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mocks.MockEmployeeProfileRepository{}
			tc.setupRepo(repo)

			svc := newEmployeeProfileService(repo)
			got, err := svc.UpdateProfile(ctx, callerID, tc.callerRole, targetUserID, tc.input)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantResult, got)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestEmployeeProfileService_PatchSelfProfile(t *testing.T) {
	ctx := context.Background()
	userID := 5
	profileID := 10
	profile := &domain.EmployeeProfile{ID: uint(profileID), UserID: uint(userID)}
	updatedProfile := &domain.EmployeeProfile{ID: uint(profileID), UserID: uint(userID), FullName: "New Self Name"}
	input := domain.UpdateEmployeeProfileInput{FullName: &updatedProfile.FullName}

	t.Run("success", func(t *testing.T) {
		repo := &mocks.MockEmployeeProfileRepository{}
		repo.On("FindByUserID", ctx, userID).Return(profile, nil)
		repo.On("Update", ctx, profileID, input).Return(updatedProfile, nil)

		svc := newEmployeeProfileService(repo)
		got, err := svc.PatchSelfProfile(ctx, userID, input)

		require.NoError(t, err)
		assert.Equal(t, updatedProfile, got)
		repo.AssertExpectations(t)
	})

	t.Run("profile not found", func(t *testing.T) {
		repo := &mocks.MockEmployeeProfileRepository{}
		repo.On("FindByUserID", ctx, userID).Return(nil, domain.ErrEmployeeProfileNotFound)

		svc := newEmployeeProfileService(repo)
		got, err := svc.PatchSelfProfile(ctx, userID, input)

		assert.ErrorIs(t, err, domain.ErrEmployeeProfileNotFound)
		assert.Nil(t, got)
		repo.AssertExpectations(t)
	})
}

// ---------------------------------------------------------------------------
// ListProfiles
// ---------------------------------------------------------------------------

// Scenario: Admin-only profile listing with limit sanitization
//
//	Covers: admin success, non-admin rejection, invalid limit clamped to 20
func TestEmployeeProfileService_ListProfiles(t *testing.T) {
	ctx := context.Background()
	profiles := []domain.EmployeeProfile{{ID: 1}, {ID: 2}}

	tests := []struct {
		name      string
		role      domain.UserRole
		limit     int
		offset    int
		wantErr   error
		wantLimit int // what the repo should receive
		wantCall  bool
	}{
		{
			name:      "admin lists with valid limit",
			role:      domain.RoleAdmin,
			limit:     10,
			offset:    0,
			wantErr:   nil,
			wantLimit: 10,
			wantCall:  true,
		},
		{
			name:      "admin with limit=0 → clamped to 20",
			role:      domain.RoleAdmin,
			limit:     0,
			offset:    0,
			wantErr:   nil,
			wantLimit: 20,
			wantCall:  true,
		},
		{
			name:      "admin with limit=200 (>100) → clamped to 20",
			role:      domain.RoleAdmin,
			limit:     200,
			offset:    0,
			wantErr:   nil,
			wantLimit: 20,
			wantCall:  true,
		},
		{
			name:     "non-admin is rejected",
			role:     domain.RoleWarehouseManager,
			limit:    10,
			wantErr:  domain.ErrInsufficientPerms,
			wantCall: false,
		},
		{
			name:     "QP cannot list profiles",
			role:     domain.RoleQP,
			limit:    10,
			wantErr:  domain.ErrInsufficientPerms,
			wantCall: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mocks.MockEmployeeProfileRepository{}
			if tc.wantCall {
				repo.On("List", ctx, tc.wantLimit, tc.offset).Return(profiles, nil)
			}

			svc := newEmployeeProfileService(repo)
			got, err := svc.ListProfiles(ctx, 1, tc.role, tc.limit, tc.offset)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, profiles, got)
			}
			repo.AssertExpectations(t)
		})
	}
}
