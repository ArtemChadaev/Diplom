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

// toDomain конвертирует DAO в доменную модель
func (r *userRepository) toDomain(dao *dao.UserDAO) *domain.User {
	if dao == nil {
		return nil
	}
	return &domain.User{
		ID:           dao.ID,
		Login:        dao.Login,
		Email:        dao.Email,
		GoogleID:     dao.GoogleID,
		TelegramID:   dao.TelegramID,
		PasswordHash: dao.PasswordHash,
		Role:         domain.UserRole(dao.Role),
		Status:       domain.UserStatus(dao.Status),
		IsBlocked:    dao.IsBlocked,
		CreatedAt:    dao.CreatedAt,
	}
}

// fromDomain конвертирует доменную модель в DAO
func (r *userRepository) fromDomain(u *domain.User) *dao.UserDAO {
	if u == nil {
		return nil
	}
	return &dao.UserDAO{
		ID:           u.ID,
		Login:        u.Login,
		Email:        u.Email,
		GoogleID:     u.GoogleID,
		TelegramID:   u.TelegramID,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
		Status:       string(u.Status),
		IsBlocked:    u.IsBlocked,
		CreatedAt:    u.CreatedAt,
	}
}

// IsLoginTaken — проверяет, занят ли логин.
func (r *userRepository) IsLoginTaken(ctx context.Context, login string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&dao.UserDAO{}).
		Where("login = ?", login).
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

func (r *userRepository) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	var u dao.UserDAO
	if err := r.db.WithContext(ctx).Where("login = ?", login).First(&u).Error; err != nil {
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
		// Basic duplicate check could be mapped here to ErrLoginTaken or ErrEmailTaken
		// For simplicity, we just return the pgerror wrapper or direct error
		return nil, err
	}
	u.ID = daoUser.ID
	u.CreatedAt = daoUser.CreatedAt
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

func (r *userRepository) UpdateStatus(ctx context.Context, userID int, status domain.UserStatus) error {
	res := r.db.WithContext(ctx).Model(&dao.UserDAO{}).Where("id = ?", userID).Update("status", string(status))
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

func (r *userRepository) SetPasswordHash(ctx context.Context, userID int, hash string) error {
	res := r.db.WithContext(ctx).Model(&dao.UserDAO{}).Where("id = ?", userID).Update("password_hash", hash)
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
			// Profile not found, just return basic user
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
