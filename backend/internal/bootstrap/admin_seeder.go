package bootstrap

import (
	"context"
	"errors"

	"github.com/ima/diplom-backend/internal/config"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

// SeedAdmin ensures the first admin user exists in the database.
// With the new email-based schema, the admin is identified by ADMIN_EMAIL.
// Authentication is not password-based, so this just ensures the record exists
// with the correct role.
func SeedAdmin(ctx context.Context, cfg *config.Config, userRepo domain.UserRepository) error {
	log := logger.FromContext(ctx)

	if cfg.AdminEmail == "" {
		log.Warn("ADMIN_EMAIL not configured, skipping admin seed")
		return nil
	}

	existing, err := userRepo.FindByEmail(ctx, cfg.AdminEmail)
	if err == nil && existing != nil {
		if existing.Role != domain.RoleAdmin {
			// Upgrade existing user to admin
			if updateErr := userRepo.UpdateRole(ctx, existing.ID, domain.RoleAdmin); updateErr != nil {
				log.Error("failed to upgrade user to admin", "email", cfg.AdminEmail, "error", updateErr)
				return updateErr
			}
			log.Info("upgraded existing user to admin", "email", cfg.AdminEmail)
		} else {
			log.Info("default admin already exists, skipping", "email", cfg.AdminEmail)
		}
		return nil
	}

	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}

	// Create user with admin role (they'll authenticate via Google OAuth)
	u := &domain.User{
		Email: cfg.AdminEmail,
		Role:  domain.RoleAdmin,
	}

	_, createErr := userRepo.Create(ctx, u)
	if createErr != nil {
		log.Error("failed to create default admin", "email", cfg.AdminEmail, "error", createErr)
		return createErr
	}

	log.Info("created default admin user", "email", cfg.AdminEmail)
	return nil
}
