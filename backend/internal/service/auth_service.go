package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"google.golang.org/api/idtoken"
)

type authService struct {
	userRepo        domain.UserRepository
	sessionRepo     domain.SessionRepository
	tokenSvc        domain.TokenService
	refreshTokenTTL time.Duration
	googleClientID  string
}

func NewAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	tokenSvc domain.TokenService,
	refreshTTL time.Duration,
	googleClientID string,
) domain.AuthService {
	return &authService{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		tokenSvc:        tokenSvc,
		refreshTokenTTL: refreshTTL,
		googleClientID:  googleClientID,
	}
}

func (s *authService) LoginWithGoogle(ctx context.Context, idToken, userAgent, ip string) (*domain.TokenPair, error) {
	payload, err := idtoken.Validate(ctx, idToken, s.googleClientID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	googleID := payload.Subject
	email, _ := payload.Claims["email"].(string)

	// 1. Try finding by Google ID
	u, err := s.userRepo.FindByGoogleID(ctx, googleID)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	if u != nil {
		return s.issueTokens(ctx, u, userAgent, ip)
	}

	// 2. Try finding by Email and link Google ID
	if email != "" {
		u, err = s.userRepo.FindByEmail(ctx, email)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}

		if u != nil {
			if err := s.userRepo.LinkGoogle(ctx, u.ID, googleID); err != nil {
				return nil, err
			}
			u.GoogleID = &googleID
			return s.issueTokens(ctx, u, userAgent, ip)
		}
	}

	// 3. New user — auto-register with default role pharmacist
	newUser := &domain.User{
		Email:    email,
		GoogleID: &googleID,
		Role:     domain.RolePharmacist,
	}

	created, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, created, userAgent, ip)
}

func (s *authService) LoginWithTelegram(ctx context.Context, data domain.TelegramAuthData, userAgent, ip string) (*domain.TokenPair, error) {
	// TODO: verify Telegram hash with bot token
	panic("LoginWithTelegram not implemented: need Telegram hash validator")
}

func (s *authService) RefreshTokens(ctx context.Context, oldRefreshToken string, meta domain.SessionMeta) (*domain.TokenPair, error) {
	hash := s.tokenSvc.HashToken(oldRefreshToken)
	session, err := s.sessionRepo.FindByTokenHash(ctx, hash)
	if err != nil {
		return nil, domain.ErrInvalidCreds
	}

	// Token theft detection: reused revoked token → revoke all sessions
	if session.RevokedAt != nil {
		_ = s.sessionRepo.RevokeAllForUser(ctx, session.UserID)
		return nil, domain.ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, domain.ErrTokenExpired
	}

	// Revoke old session
	if err := s.sessionRepo.Revoke(ctx, session.ID); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user, meta.UserAgent, meta.IPAddress)
}

func (s *authService) RevokeSession(ctx context.Context, sessionID uuid.UUID, callerID int, callerRole domain.UserRole) error {
	sess, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if sess.UserID != callerID && callerRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}
	return s.sessionRepo.Revoke(ctx, sessionID)
}

func (s *authService) AssignRole(ctx context.Context, adminID int, adminRole domain.UserRole, targetUserID int, role domain.UserRole) error {
	if adminRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}
	return s.userRepo.UpdateRole(ctx, targetUserID, role)
}

func (s *authService) SetBlocked(ctx context.Context, adminID int, adminRole domain.UserRole, targetUserID int, blocked bool) error {
	if adminRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}
	return s.userRepo.SetBlocked(ctx, targetUserID, blocked)
}

// issueTokens is a private helper: validates user state, creates refresh session, returns token pair.
func (s *authService) issueTokens(ctx context.Context, u *domain.User, userAgent, ip string) (*domain.TokenPair, error) {
	if u.IsBlocked {
		return nil, domain.ErrUserBlocked
	}

	rawRT, hashRT, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	rt := &domain.RefreshToken{
		UserID:    u.ID,
		TokenHash: hashRT,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
		UserAgent: userAgent,
		IPAddress: ip,
		Metadata:  make(map[string]any),
	}

	createdRT, err := s.sessionRepo.Create(ctx, rt)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(u, createdRT.ID)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		ExpiresIn:    900, // 15 min in seconds
		RefreshToken: rawRT,
	}, nil
}
