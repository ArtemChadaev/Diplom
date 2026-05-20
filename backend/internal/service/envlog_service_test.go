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
	// Scenario: Environment logs are recorded successfully
	//   Given: Active shift monitoring with no existing entries for this shift/date
	//   When: RecordLogs is called
	//   Then: Logs are validated, enriched, and successfully saved
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockEnvironmentLogRepository)
		mockZoneRepo := new(mocks.MockZoneRepository)
		mockProfileRepo := new(mocks.MockEmployeeProfileRepository)
		svc := NewEnvironmentLogService(mockRepo, mockZoneRepo, mockProfileRepo)
		ctx := context.Background()

		profile := &domain.EmployeeProfile{
			ID:     1,
			UserID: 1,
			GDPTrainingHistory: []domain.GDPTrainingRecord{
				{Date: "2026-01-01", CourseName: "GDP Guidelines", Result: "pass"},
			},
		}

		logs := []domain.EnvironmentLog{
			{ZoneID: "Z1", Shift: "morning", Temperature: 22.5, Humidity: 45.0},
			{ZoneID: "Z2", Shift: "morning", Temperature: 18.0, Humidity: 50.0},
		}

		mockProfileRepo.On("FindByUserID", mock.Anything, 1).Return(profile, nil).Once()

		mockRepo.On("ExistsByZoneShiftDate", mock.Anything, "Z1", "morning", mock.Anything).Return(false, nil)
		mockRepo.On("ExistsByZoneShiftDate", mock.Anything, "Z2", "morning", mock.Anything).Return(false, nil)

		mockZoneRepo.On("GetByID", mock.Anything, "Z1").Return(&domain.Zone{TemperatureMin: 15.0, TemperatureMax: 25.0}, nil).Once()
		mockZoneRepo.On("GetByID", mock.Anything, "Z2").Return(&domain.Zone{TemperatureMin: 15.0, TemperatureMax: 25.0}, nil).Once()

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(l *domain.EnvironmentLog) bool {
			return l.RecordedBy == 1 && !l.RecordedAt.IsZero() && l.Shift == "morning"
		})).Return(nil).Times(2)

		err := svc.RecordLogs(ctx, 1, logs)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockZoneRepo.AssertExpectations(t)
		mockProfileRepo.AssertExpectations(t)
	})

	// Scenario: Reject log due to duplicate shift entry
	//   Given: A log has already been submitted for Z1 morning shift today
	//   When: RecordLogs is called with the same shift/date
	//   Then: ErrEnvLogDuplicateShift is returned
	t.Run("duplicate shift error", func(t *testing.T) {
		mockRepo := new(mocks.MockEnvironmentLogRepository)
		mockZoneRepo := new(mocks.MockZoneRepository)
		mockProfileRepo := new(mocks.MockEmployeeProfileRepository)
		svc := NewEnvironmentLogService(mockRepo, mockZoneRepo, mockProfileRepo)
		ctx := context.Background()

		profile := &domain.EmployeeProfile{
			ID:     1,
			UserID: 1,
			GDPTrainingHistory: []domain.GDPTrainingRecord{
				{Date: "2026-01-01", CourseName: "GDP Guidelines", Result: "pass"},
			},
		}

		logs := []domain.EnvironmentLog{
			{ZoneID: "Z1", Shift: "morning", Temperature: 22.5, Humidity: 45.0},
		}

		mockProfileRepo.On("FindByUserID", mock.Anything, 1).Return(profile, nil).Once()
		mockRepo.On("ExistsByZoneShiftDate", mock.Anything, "Z1", "morning", mock.Anything).Return(true, nil)

		err := svc.RecordLogs(ctx, 1, logs)
		assert.ErrorIs(t, err, domain.ErrEnvLogDuplicateShift)
		mockRepo.AssertExpectations(t)
		mockProfileRepo.AssertExpectations(t)
	})
}
