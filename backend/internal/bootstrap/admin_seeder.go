package bootstrap

import (
	"context"

	"github.com/ima/diplom-backend/internal/config"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(ctx context.Context, cfg *config.Config, userRepo domain.UserRepository) error {
	log := logger.FromContext(ctx)

	if cfg.AdminUser == "" || cfg.AdminPassword == "" {
		log.Warn("admin credentials not fully configured, skipping bootstrap")
		return nil
	}

	existingUser, err := userRepo.FindByLogin(ctx, cfg.AdminUser)
	if err == nil && existingUser != nil {
		log.Info("default admin already exists, skipping creation")
		return nil
	}
	if err != nil && err != domain.ErrUserNotFound {
		return err
	}

	// Not found, so we create
	bytes, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), 12)
	if err != nil {
		return err
	}
	hash := string(bytes)

	u := &domain.User{
		Login:        cfg.AdminUser,
		PasswordHash: &hash,
		Role:         domain.RoleAdmin,
		Status:       domain.StatusActive, // Automatically active
	}

	_, err = userRepo.Create(ctx, u)
	if err != nil {
		log.Error("failed to create default admin", "error", err)
		return err
	}

	log.Info("successfully created default admin", "login", cfg.AdminUser)
	return nil
}
