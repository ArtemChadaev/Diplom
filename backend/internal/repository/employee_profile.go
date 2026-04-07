package repository

import (
	"context"
	"errors"
	"encoding/json"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/pkg/logger"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type employeeProfileRepository struct {
	db *gorm.DB
}

func NewEmployeeProfileRepository(db *gorm.DB) domain.EmployeeProfileRepository {
	return &employeeProfileRepository{db: db}
}

func (r *employeeProfileRepository) toDomain(d *dao.EmployeeProfileDAO) *domain.EmployeeProfile {
	var gdp []domain.GDPTrainingRecord
	_ = json.Unmarshal(d.GDPTrainingHistory, &gdp)

	return &domain.EmployeeProfile{
		ID:                 d.ID,
		UserID:             d.UserID,
		EmployeeCode:       d.EmployeeCode,
		FullName:           d.FullName,
		CorporateEmail:     d.CorporateEmail,
		Phone:              d.Phone,
		Position:           d.Position,
		Department:         d.Department,
		BirthDate:          d.BirthDate,
		AvatarURL:          d.AvatarURL,
		HireDate:           d.HireDate,
		DismissalDate:      d.DismissalDate,
		MedicalBookScanURL: d.MedicalBookScanURL,
		SpecialZoneAccess:  d.SpecialZoneAccess,
		GDPTrainingHistory: gdp,
	}
}

func (r *employeeProfileRepository) FindByUserID(ctx context.Context, userID int) (*domain.EmployeeProfile, error) {
	var d dao.EmployeeProfileDAO
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeProfileNotFound
		}
		return nil, errors.New("employeeProfileRepo.FindByUserID: " + err.Error())
	}
	return r.toDomain(&d), nil
}

func (r *employeeProfileRepository) FindByID(ctx context.Context, id int) (*domain.EmployeeProfile, error) {
	var d dao.EmployeeProfileDAO
	if err := r.db.WithContext(ctx).First(&d, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeProfileNotFound
		}
		return nil, errors.New("employeeProfileRepo.FindByID: " + err.Error())
	}
	return r.toDomain(&d), nil
}

func (r *employeeProfileRepository) Update(ctx context.Context, id int, input domain.UpdateEmployeeProfileInput) (*domain.EmployeeProfile, error) {
	updates := buildProfileUpdateMap(input)

	if len(updates) == 0 {
		return r.FindByID(ctx, id)
	}

	if err := r.db.WithContext(ctx).
		Model(&dao.EmployeeProfileDAO{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, errors.New("employeeProfileRepo.Update: " + err.Error())
	}

	logger.FromContext(ctx).Info("employee profile updated",
		"profile_id", id,
		"fields_updated", len(updates),
	)

	return r.FindByID(ctx, id)
}

func (r *employeeProfileRepository) List(ctx context.Context, limit, offset int) ([]domain.EmployeeProfile, error) {
	var rows []dao.EmployeeProfileDAO
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, errors.New("employeeProfileRepo.List: " + err.Error())
	}
	result := make([]domain.EmployeeProfile, len(rows))
	for i, row := range rows {
		result[i] = *r.toDomain(&row)
	}
	return result, nil
}

func buildProfileUpdateMap(input domain.UpdateEmployeeProfileInput) map[string]any {
	m := make(map[string]any)
	if input.FullName != nil           { m["full_name"] = *input.FullName }
	if input.CorporateEmail != nil     { m["corporate_email"] = *input.CorporateEmail }
	if input.Phone != nil              { m["phone"] = *input.Phone }
	if input.Position != nil           { m["position"] = *input.Position }
	if input.Department != nil         { m["department"] = *input.Department }
	if input.BirthDate != nil          { m["birth_date"] = *input.BirthDate }
	if input.AvatarURL != nil          { m["avatar_url"] = *input.AvatarURL }
	if input.HireDate != nil           { m["hire_date"] = *input.HireDate }
	if input.DismissalDate != nil      { m["dismissal_date"] = *input.DismissalDate }
	if input.MedicalBookScanURL != nil { m["medical_book_scan_url"] = *input.MedicalBookScanURL }
	if input.SpecialZoneAccess != nil  { m["special_zone_access"] = *input.SpecialZoneAccess }
	if input.GDPTrainingHistory != nil { m["gdp_training_history"] = input.GDPTrainingHistory }
	return m
}
