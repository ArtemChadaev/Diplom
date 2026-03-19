package bootstrap

import (
	"context"
	"log/slog"

	"github.com/ima/diplom-backend/internal/config"
	"github.com/ima/diplom-backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(ctx context.Context, cfg *config.Config, userRepo domain.UserRepository) error {
	if cfg.AdminUser == "" || cfg.AdminPassword == "" {
		slog.Warn("admin credentials not fully configured, skipping bootstrap")
		return nil
	}

	existingUser, err := userRepo.FindByLogin(ctx, cfg.AdminUser)
	if err == nil && existingUser != nil {
		slog.Info("default admin already exists, skipping creation")
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
		slog.Error("failed to create default admin", slog.Any("error", err))
		return err
	}

	slog.Info("successfully created default admin", slog.String("login", cfg.AdminUser))
	return nil
}
