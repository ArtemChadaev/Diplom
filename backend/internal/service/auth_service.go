package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
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

func (s *authService) Register(ctx context.Context, req domain.RegisterInput) (*domain.User, error) {
	// 1. Check if login explicitly taken
	taken, err := s.userRepo.IsLoginTaken(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, domain.ErrLoginTaken
	}

	if req.Email != "" {
		// optionally throw ErrEmailTaken if doing a manual check
		// For now pg err code would handle it or we do explicit:
		if existing, _ := s.userRepo.FindByEmail(ctx, req.Email); existing != nil {
			return nil, domain.ErrEmailTaken
		}
	}

	// 2. Hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}
	hash := string(bytes)
	var em *string
	if req.Email != "" {
		em = &req.Email
	}

	u := &domain.User{
		Login:        req.Login,
		Email:        em,
		PasswordHash: &hash,
		Role:         domain.RoleUnverified,
		Status:       domain.StatusUnverified,
	}

	// 3. Create user
	return s.userRepo.Create(ctx, u)
}

func (s *authService) LoginWithPassword(ctx context.Context, login, password, userAgent, ip string) (*domain.TokenPair, error) {
	u, err := s.userRepo.FindByLogin(ctx, login)
	if err != nil {
		return nil, domain.ErrInvalidCreds // prevent user enumeration
	}

	if u.PasswordHash == nil || bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(password)) != nil {
		return nil, domain.ErrInvalidCreds
	}

	return s.issueTokens(ctx, u, userAgent, ip)
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

	// 2. Try finding by Email and link
	if email != "" {
		u, err = s.userRepo.FindByEmail(ctx, email)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}

		if u != nil {
			if err := s.userRepo.LinkGoogle(ctx, u.ID, googleID); err != nil {
				return nil, err
			}
			return s.issueTokens(ctx, u, userAgent, ip)
		}
	}

	// 3. New user -> Register as unverified
	login := email
	if login == "" {
		login = "google_" + googleID[:8]
	}

	newUser := &domain.User{
		Login:    login,
		Email:    &email,
		GoogleID: &googleID,
		Role:     domain.RoleUnverified,
		Status:   domain.StatusUnverified,
	}

	if _, err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// For unverified users, we don't issue tokens yet
	return nil, domain.ErrUserUnverified
}

func (s *authService) LoginWithTelegram(ctx context.Context, data domain.TelegramAuthData, userAgent, ip string) (*domain.TokenPair, error) {
	// Telegram hash verification not implemented here to save lines
	panic("implement telegram hash validator")
}

func (s *authService) RefreshTokens(ctx context.Context, oldRefreshToken string, meta domain.SessionMeta) (*domain.TokenPair, error) {
	hash := s.tokenSvc.HashToken(oldRefreshToken)
	session, err := s.sessionRepo.FindByTokenHash(ctx, hash)
	if err != nil {
		return nil, domain.ErrInvalidCreds // or specific refresh error
	}

	// Token theft detection: if an old token is reused that was revoked
	if session.RevokedAt != nil {
		_ = s.sessionRepo.RevokeAllForUser(ctx, session.UserID)
		return nil, domain.ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, domain.ErrTokenExpired
	}

	// Revoke old session atomically (if this fails we abort to prevent dup issue)
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

func (s *authService) VerifyUser(ctx context.Context, adminID int, adminRole domain.UserRole, targetUserID int) error {
	if adminRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}
	return s.userRepo.UpdateStatus(ctx, targetUserID, domain.StatusActive)
}

func (s *authService) AssignRole(ctx context.Context, adminID int, adminRole domain.UserRole, targetUserID int, role domain.UserRole) error {
	if adminRole != domain.RoleAdmin {
		return domain.ErrInsufficientPerms
	}
	return s.userRepo.UpdateRole(ctx, targetUserID, role)
}

// issueTokens is a private helper to generate a pair and insert the refresh session
func (s *authService) issueTokens(ctx context.Context, u *domain.User, userAgent, ip string) (*domain.TokenPair, error) {
	if u.IsBlocked || u.Status != domain.StatusActive {
		if u.Status == domain.StatusUnverified {
			return nil, domain.ErrUserUnverified
		}
		return nil, domain.ErrUserBlocked
	}

	rawRT, hashRT, err := s.tokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]any)
	rt := &domain.RefreshToken{
		UserID:    u.ID,
		TokenHash: hashRT,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
		UserAgent: userAgent,
		IPAddress: ip,
		Metadata:  metadata,
	}

	createdRT, err := s.sessionRepo.Create(ctx, rt)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(u, createdRT.ID)
	if err != nil {
		return nil, err
	}

	// 15m default hard code assumption vs Config, we can derive from Config/TTL.
	// For API response:
	return &domain.TokenPair{
		AccessToken:  accessToken,
		ExpiresIn:    900, // 15 min TTL default
		RefreshToken: rawRT,
	}, nil
}
