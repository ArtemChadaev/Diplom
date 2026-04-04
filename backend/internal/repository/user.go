package repository

import (
	"context"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/repository/dao"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// toDomain converts DAO to domain model
func (r *userRepository) toDomain(d *dao.UserDAO) *domain.User {
	if d == nil {
		return nil
	}
	return &domain.User{
		ID:         d.ID,
		Email:      d.Email,
		GoogleID:   d.GoogleID,
		TelegramID: d.TelegramID,
		Role:       domain.UserRole(d.Role),
		NsPvAccess: d.NsPvAccess,
		UkepBound:  d.UkepBound,
		IsBlocked:  d.IsBlocked,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}

// fromDomain converts domain model to DAO
func (r *userRepository) fromDomain(u *domain.User) *dao.UserDAO {
	if u == nil {
		return nil
	}
	return &dao.UserDAO{
		ID:         u.ID,
		Email:      u.Email,
		GoogleID:   u.GoogleID,
		TelegramID: u.TelegramID,
		Role:       string(u.Role),
		NsPvAccess: u.NsPvAccess,
		UkepBound:  u.UkepBound,
		IsBlocked:  u.IsBlocked,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (r *userRepository) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&dao.UserDAO{}).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*domain.User, error) {
	var u dao.UserDAO
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.toDomain(&u), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u dao.UserDAO
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.toDomain(&u), nil
}

func (r *userRepository) FindByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	var u dao.UserDAO
	if err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.toDomain(&u), nil
}

func (r *userRepository) FindByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	var u dao.UserDAO
	if err := r.db.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return r.toDomain(&u), nil
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	daoUser := r.fromDomain(u)
	if err := r.db.WithContext(ctx).Create(daoUser).Error; err != nil {
		return nil, err
	}
	u.ID = daoUser.ID
	u.CreatedAt = daoUser.CreatedAt
	u.UpdatedAt = daoUser.UpdatedAt
	return u, nil
}

func (r *userRepository) UpdateRole(ctx context.Context, userID int, role domain.UserRole) error {
	res := r.db.WithContext(ctx).Model(&dao.UserDAO{}).Where("id = ?", userID).Update("role", string(role))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) LinkGoogle(ctx context.Context, userID int, googleID string) error {
	res := r.db.WithContext(ctx).Model(&dao.UserDAO{}).Where("id = ?", userID).Update("google_id", googleID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) LinkTelegram(ctx context.Context, userID int, telegramID int64) error {
	res := r.db.WithContext(ctx).Model(&dao.UserDAO{}).Where("id = ?", userID).Update("telegram_id", telegramID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) SetNsPvAccess(ctx context.Context, userID int, access bool) error {
	res := r.db.WithContext(ctx).Model(&dao.UserDAO{}).Where("id = ?", userID).Update("ns_pv_access", access)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) SetBlocked(ctx context.Context, userID int, blocked bool) error {
	res := r.db.WithContext(ctx).Model(&dao.UserDAO{}).Where("id = ?", userID).Update("is_blocked", blocked)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) FindProfileByUserID(ctx context.Context, userID int) (*domain.UserProfile, error) {
	u, err := r.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var p dao.EmployeeProfileDAO
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Profile not found, return basic user info only
			return &domain.UserProfile{User: *u}, nil
		}
		return nil, err
	}

	return &domain.UserProfile{
		User:         *u,
		EmployeeCode: p.EmployeeCode,
		FullName:     p.FullName,
		Position:     p.Position,
		Department:   p.Department,
	}, nil
}

func (r *userRepository) List(ctx context.Context, filter domain.UserListFilter) ([]*domain.UserProfile, int, error) {
	type Result struct {
		dao.UserDAO
		EmployeeCode string
		FullName     string
		Position     string
		Department   string
	}

	var results []Result
	var total int64

	query := r.db.WithContext(ctx).Table("users u").
		Joins("LEFT JOIN employee_profiles ep ON ep.user_id = u.id")

	if filter.Query != "" {
		likeQuery := "%" + filter.Query + "%"
		query = query.Where("u.email ILIKE ? OR ep.full_name ILIKE ?", likeQuery, likeQuery)
	}
	if filter.Role != "" {
		query = query.Where("u.role = ?", string(filter.Role))
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	err = query.Select("u.*, ep.employee_code, ep.full_name, ep.position, ep.department").
		Order("u.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		return nil, 0, err
	}

	profiles := make([]*domain.UserProfile, len(results))
	for i, res := range results {
		profiles[i] = &domain.UserProfile{
			User:         *r.toDomain(&res.UserDAO),
			EmployeeCode: res.EmployeeCode,
			FullName:     res.FullName,
			Position:     res.Position,
			Department:   res.Department,
		}
	}

	return profiles, int(total), nil
}
