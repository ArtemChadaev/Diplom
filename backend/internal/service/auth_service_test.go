package service_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/service"
	"github.com/ima/diplom-backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	otpRepo domain.OTPRepository,
	tokenSvc domain.TokenService,
	mailer *mocks.MockMailer,
) domain.AuthService {
	return service.NewAuthService(
		userRepo,
		sessionRepo,
		otpRepo,
		tokenSvc,
		15*24*time.Hour, // refreshTTL
		"",              // googleClientID — not tested here
		mailer,
		"test-hmac-secret",
	)
}

func freshMocks() (
	*mocks.MockUserRepository,
	*mocks.MockSessionRepository,
	*mocks.MockOTPRepository,
	*mocks.MockTokenService,
	*mocks.MockMailer,
) {
	return &mocks.MockUserRepository{},
		&mocks.MockSessionRepository{},
		&mocks.MockOTPRepository{},
		&mocks.MockTokenService{},
		&mocks.MockMailer{}
}

var sessionMeta = domain.SessionMeta{UserAgent: "test-agent", IPAddress: "127.0.0.1"}

// ---------------------------------------------------------------------------
// VerifyOTPCode
// ---------------------------------------------------------------------------

// Scenario: OTP code verification gate
//
//	Covers: valid code, wrong code, expired/not found, max attempts, blocked user, user not found
func TestAuthService_VerifyOTPCode(t *testing.T) {
	ctx := context.Background()
	email := "user@example.com"
	correctCode := "123456"

	baseUser := &domain.User{ID: 1, Email: email, IsBlocked: false}
	activeSession := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    1,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	tests := []struct {
		name        string
		setupUser   func(*mocks.MockUserRepository)
		setupOTP    func(*mocks.MockOTPRepository)
		setupSess   func(*mocks.MockSessionRepository)
		setupTok    func(*mocks.MockTokenService)
		inputCode   string
		wantErr     error
		wantNilPair bool
	}{
		{
			name: "valid code → token pair returned",
			setupUser: func(r *mocks.MockUserRepository) {
				r.On("FindByEmail", ctx, email).Return(baseUser, nil)
			},
			setupOTP: func(r *mocks.MockOTPRepository) {
				// Simulate a stored entry that matches "123456" with the service's HMAC.
				// We produce the hash here the same way the service code does.
				hash := hmacForTest(correctCode, "test-hmac-secret")
				r.On("Get", ctx, 1).Return(&domain.OTPCode{UserID: 1, CodeHash: hash, Attempts: 0}, nil)
				r.On("Delete", ctx, 1).Return(nil)
			},
			setupSess: func(r *mocks.MockSessionRepository) {
				r.On("Create", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(activeSession, nil)
			},
			setupTok: func(r *mocks.MockTokenService) {
				r.On("GenerateRefreshToken").Return("raw-rt", "hash-rt", nil)
				r.On("GenerateAccessToken", baseUser, activeSession.ID).Return("access-jwt", nil)
			},
			inputCode:   correctCode,
			wantErr:     nil,
			wantNilPair: false,
		},
		{
			name: "wrong code → ErrOTPInvalid, attempts incremented",
			setupUser: func(r *mocks.MockUserRepository) {
				r.On("FindByEmail", ctx, email).Return(baseUser, nil)
			},
			setupOTP: func(r *mocks.MockOTPRepository) {
				// Store a hash for a DIFFERENT code ("000000")
				hash := hmacForTest("000000", "test-hmac-secret")
				r.On("Get", ctx, 1).Return(&domain.OTPCode{UserID: 1, CodeHash: hash, Attempts: 0}, nil)
				r.On("IncrAttempts", ctx, 1).Return(nil)
			},
			setupSess: func(r *mocks.MockSessionRepository) {},
			setupTok:  func(r *mocks.MockTokenService) {},
			inputCode: "111111", // != "000000" → wrong
			wantErr:   domain.ErrOTPInvalid,
		},
		{
			name: "OTP not found (expired or never sent) → ErrOTPNotFound",
			setupUser: func(r *mocks.MockUserRepository) {
				r.On("FindByEmail", ctx, email).Return(baseUser, nil)
			},
			setupOTP: func(r *mocks.MockOTPRepository) {
				r.On("Get", ctx, 1).Return(nil, domain.ErrOTPNotFound)
			},
			setupSess:   func(r *mocks.MockSessionRepository) {},
			setupTok:    func(r *mocks.MockTokenService) {},
			inputCode:   correctCode,
			wantErr:     domain.ErrOTPNotFound,
			wantNilPair: true,
		},
		{
			name: "max attempts reached → ErrOTPMaxAttempts, IncrAttempts NOT called",
			setupUser: func(r *mocks.MockUserRepository) {
				r.On("FindByEmail", ctx, email).Return(baseUser, nil)
			},
			setupOTP: func(r *mocks.MockOTPRepository) {
				hash := hmacForTest(correctCode, "test-hmac-secret")
				// Attempts == OTPMaxAttempts (5)
				r.On("Get", ctx, 1).Return(&domain.OTPCode{UserID: 1, CodeHash: hash, Attempts: domain.OTPMaxAttempts}, nil)
				// IncrAttempts must NOT be called — no expectation set on purpose
			},
			setupSess:   func(r *mocks.MockSessionRepository) {},
			setupTok:    func(r *mocks.MockTokenService) {},
			inputCode:   correctCode,
			wantErr:     domain.ErrOTPMaxAttempts,
			wantNilPair: true,
		},
		{
			name: "blocked user → ErrUserBlocked",
			setupUser: func(r *mocks.MockUserRepository) {
				blocked := &domain.User{ID: 2, Email: email, IsBlocked: true}
				r.On("FindByEmail", ctx, email).Return(blocked, nil)
			},
			setupOTP:    func(r *mocks.MockOTPRepository) {},
			setupSess:   func(r *mocks.MockSessionRepository) {},
			setupTok:    func(r *mocks.MockTokenService) {},
			inputCode:   correctCode,
			wantErr:     domain.ErrUserBlocked,
			wantNilPair: true,
		},
		{
			name: "user not found → error propagated",
			setupUser: func(r *mocks.MockUserRepository) {
				r.On("FindByEmail", ctx, email).Return(nil, domain.ErrUserNotFound)
			},
			setupOTP:    func(r *mocks.MockOTPRepository) {},
			setupSess:   func(r *mocks.MockSessionRepository) {},
			setupTok:    func(r *mocks.MockTokenService) {},
			inputCode:   correctCode,
			wantErr:     domain.ErrUserNotFound,
			wantNilPair: true,
		},
	}


	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo, sessRepo, otpRepo, tokenSvc, mailer := freshMocks()
			tc.setupUser(userRepo)
			tc.setupOTP(otpRepo)
			tc.setupSess(sessRepo)
			tc.setupTok(tokenSvc)

			svc := newAuthService(userRepo, sessRepo, otpRepo, tokenSvc, mailer)
			pair, err := svc.VerifyOTPCode(ctx, email, tc.inputCode, sessionMeta)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			if tc.wantNilPair {
				assert.Nil(t, pair)
			}

			userRepo.AssertExpectations(t)
			otpRepo.AssertExpectations(t)
			sessRepo.AssertExpectations(t)
			tokenSvc.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// SendOTPCode
// ---------------------------------------------------------------------------

// Scenario: Rate-limited OTP dispatch
//
//	Given: user exists; OTP attempts are checked before sending
//	Covers: success, rate-limit block, user not found
func TestAuthService_SendOTPCode(t *testing.T) {
	ctx := context.Background()
	email := "user@example.com"
	user := &domain.User{ID: 1, Email: email}

	tests := []struct {
		name      string
		setupUser func(*mocks.MockUserRepository)
		setupOTP  func(*mocks.MockOTPRepository)
		setupMail func(*mocks.MockMailer)
		wantErr   error
	}{
		{
			name: "success — OTP stored and email sent",
			setupUser: func(r *mocks.MockUserRepository) {
				r.On("FindByEmail", ctx, email).Return(user, nil)
			},
			setupOTP: func(r *mocks.MockOTPRepository) {
				// No existing OTP entry
				r.On("Get", ctx, 1).Return(nil, domain.ErrOTPNotFound)
				r.On("Store", ctx, 1, mock.AnythingOfType("string")).Return(nil)
			},
			setupMail: func(m *mocks.MockMailer) {
				m.On("SendOTP", ctx, email, mock.AnythingOfType("string")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "rate limit — attempts at max → ErrOTPMaxAttempts",
			setupUser: func(r *mocks.MockUserRepository) {
				r.On("FindByEmail", ctx, email).Return(user, nil)
			},
			setupOTP: func(r *mocks.MockOTPRepository) {
				r.On("Get", ctx, 1).Return(&domain.OTPCode{Attempts: domain.OTPMaxAttempts}, nil)
				// Store must NOT be called — no expectation set
			},
			setupMail: func(m *mocks.MockMailer) {
				// SendOTP must NOT be called
			},
			wantErr: domain.ErrOTPMaxAttempts,
		},
		{
			name: "user not found → error propagated",
			setupUser: func(r *mocks.MockUserRepository) {
				r.On("FindByEmail", ctx, email).Return(nil, domain.ErrUserNotFound)
			},
			setupOTP:  func(r *mocks.MockOTPRepository) {},
			setupMail: func(m *mocks.MockMailer) {},
			wantErr:   domain.ErrUserNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo, sessRepo, otpRepo, tokenSvc, mailer := freshMocks()
			tc.setupUser(userRepo)
			tc.setupOTP(otpRepo)
			tc.setupMail(mailer)

			svc := newAuthService(userRepo, sessRepo, otpRepo, tokenSvc, mailer)
			err := svc.SendOTPCode(ctx, email)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
			otpRepo.AssertExpectations(t)
			mailer.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// RefreshTokens
// ---------------------------------------------------------------------------

// Scenario: Token rotation and theft detection (TC-02 from testing.md)
//
//	Covers: successful rotation, theft (reused revoked token), token not found, expired
func TestAuthService_RefreshTokens(t *testing.T) {
	ctx := context.Background()
	userID := 5
	rawRT := "old-raw-refresh-token"

	activeSession := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour),
		RevokedAt: nil,
	}
	revokedTime := time.Now().Add(-time.Minute)
	revokedSession := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour),
		RevokedAt: &revokedTime,
	}
	expiredSession := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(-time.Hour), // expired
		RevokedAt: nil,
	}
	validUser := &domain.User{ID: userID, Email: "u@example.com", IsBlocked: false}
	newSession := &domain.RefreshToken{ID: uuid.New(), UserID: userID, ExpiresAt: time.Now().Add(time.Hour)}

	tests := []struct {
		name        string
		setupSess   func(*mocks.MockSessionRepository)
		setupUser   func(*mocks.MockUserRepository)
		setupTok    func(*mocks.MockTokenService)
		wantErr     error
		wantNilPair bool
	}{
		{
			name: "valid token — successful rotation",
			setupSess: func(r *mocks.MockSessionRepository) {
				r.On("FindByTokenHash", ctx, mock.AnythingOfType("string")).Return(activeSession, nil)
				r.On("Revoke", ctx, activeSession.ID).Return(nil)
				r.On("Create", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(newSession, nil)
			},
			setupUser: func(r *mocks.MockUserRepository) {
				r.On("FindByID", ctx, userID).Return(validUser, nil)
			},
			setupTok: func(r *mocks.MockTokenService) {
				r.On("HashToken", rawRT).Return("hashed-rt")
				r.On("GenerateRefreshToken").Return("new-raw-rt", "new-hash-rt", nil)
				r.On("GenerateAccessToken", validUser, newSession.ID).Return("new-access-token", nil)
			},
			wantErr: nil,
		},
		{
			name: "theft detected — reused revoked token → all sessions revoked",
			setupSess: func(r *mocks.MockSessionRepository) {
				r.On("FindByTokenHash", ctx, mock.AnythingOfType("string")).Return(revokedSession, nil)
				r.On("RevokeAllForUser", ctx, userID).Return(nil)
			},
			setupUser: func(r *mocks.MockUserRepository) {},
			setupTok: func(r *mocks.MockTokenService) {
				r.On("HashToken", rawRT).Return("hashed-rt")
			},
			wantErr:     domain.ErrSessionNotFound,
			wantNilPair: true,
		},
		{
			name: "token hash not found → ErrInvalidCreds",
			setupSess: func(r *mocks.MockSessionRepository) {
				r.On("FindByTokenHash", ctx, mock.AnythingOfType("string")).Return(nil, errors.New("not found"))
			},
			setupUser: func(r *mocks.MockUserRepository) {},
			setupTok: func(r *mocks.MockTokenService) {
				r.On("HashToken", rawRT).Return("hashed-rt")
			},
			wantErr:     domain.ErrInvalidCreds,
			wantNilPair: true,
		},
		{
			name: "token expired (not revoked, but ExpiresAt in past) → ErrTokenExpired",
			setupSess: func(r *mocks.MockSessionRepository) {
				r.On("FindByTokenHash", ctx, mock.AnythingOfType("string")).Return(expiredSession, nil)
			},
			setupUser: func(r *mocks.MockUserRepository) {},
			setupTok: func(r *mocks.MockTokenService) {
				r.On("HashToken", rawRT).Return("hashed-rt")
			},
			wantErr:     domain.ErrTokenExpired,
			wantNilPair: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo, sessRepo, otpRepo, tokenSvc, mailer := freshMocks()
			tc.setupSess(sessRepo)
			tc.setupUser(userRepo)
			tc.setupTok(tokenSvc)

			svc := newAuthService(userRepo, sessRepo, otpRepo, tokenSvc, mailer)
			pair, err := svc.RefreshTokens(ctx, rawRT, sessionMeta)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, pair)
			}
			if tc.wantNilPair {
				assert.Nil(t, pair)
			}

			sessRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			tokenSvc.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// RevokeSession
// ---------------------------------------------------------------------------

// Scenario: Session revocation with ownership check
//
//	Covers: owner revokes their own, admin revokes any, non-admin revokes another's
func TestAuthService_RevokeSession(t *testing.T) {
	ctx := context.Background()
	sessID := uuid.New()

	tests := []struct {
		name       string
		sessionUID int // UserID on the session
		callerID   int
		callerRole domain.UserRole
		wantErr    error
		wantRevoke bool
	}{
		{
			name:       "owner revokes their own session",
			sessionUID: 1,
			callerID:   1,
			callerRole: domain.RolePharmacist,
			wantErr:    nil,
			wantRevoke: true,
		},
		{
			name:       "admin revokes any user's session",
			sessionUID: 2,
			callerID:   1,
			callerRole: domain.RoleAdmin,
			wantErr:    nil,
			wantRevoke: true,
		},
		{
			name:       "non-admin cannot revoke another user's session",
			sessionUID: 2,
			callerID:   1,
			callerRole: domain.RoleStorekeeper,
			wantErr:    domain.ErrInsufficientPerms,
			wantRevoke: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo, sessRepo, otpRepo, tokenSvc, mailer := freshMocks()
			sess := &domain.RefreshToken{ID: sessID, UserID: tc.sessionUID}
			sessRepo.On("FindByID", ctx, sessID).Return(sess, nil)
			if tc.wantRevoke {
				sessRepo.On("Revoke", ctx, sessID).Return(nil)
			}

			svc := newAuthService(userRepo, sessRepo, otpRepo, tokenSvc, mailer)
			err := svc.RevokeSession(ctx, sessID, tc.callerID, tc.callerRole)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			sessRepo.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// AssignRole
// ---------------------------------------------------------------------------

// Scenario: Admin-only role assignment
//
//	Covers: admin success, non-admin rejection
func TestAuthService_AssignRole(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		callerRole domain.UserRole
		wantErr    error
		wantUpdate bool
	}{
		{
			name:       "admin assigns role successfully",
			callerRole: domain.RoleAdmin,
			wantErr:    nil,
			wantUpdate: true,
		},
		{
			name:       "non-admin gets ErrInsufficientPerms",
			callerRole: domain.RolePharmacist,
			wantErr:    domain.ErrInsufficientPerms,
			wantUpdate: false,
		},
		{
			name:       "QP cannot assign roles",
			callerRole: domain.RoleQP,
			wantErr:    domain.ErrInsufficientPerms,
			wantUpdate: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo, sessRepo, otpRepo, tokenSvc, mailer := freshMocks()
			if tc.wantUpdate {
				userRepo.On("UpdateRole", ctx, 99, domain.RoleWarehouseManager).Return(nil)
			}

			svc := newAuthService(userRepo, sessRepo, otpRepo, tokenSvc, mailer)
			err := svc.AssignRole(ctx, 1, tc.callerRole, 99, domain.RoleWarehouseManager)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			userRepo.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// RegisterByEmail
// ---------------------------------------------------------------------------

// Scenario: New user registration via email
//
//	Covers: success (OTP dispatched), email already taken
func TestAuthService_RegisterByEmail(t *testing.T) {
	ctx := context.Background()
	email := "new@example.com"
	createdUser := &domain.User{ID: 10, Email: email, Role: domain.RolePharmacist}

	tests := []struct {
		name        string
		emailTaken  bool
		takenErr    error
		createErr   error
		storeErr    error
		mailErr     error
		wantErr     error
		wantCreate  bool
		wantStore   bool
		wantMail    bool
	}{
		{
			name:       "success — user created, OTP stored and sent",
			wantCreate: true,
			wantStore:  true,
			wantMail:   true,
			wantErr:    nil,
		},
		{
			name:       "email already taken → ErrEmailTaken",
			emailTaken: true,
			wantCreate: false,
			wantStore:  false,
			wantMail:   false,
			wantErr:    domain.ErrEmailTaken,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo, sessRepo, otpRepo, tokenSvc, mailer := freshMocks()

			userRepo.On("IsEmailTaken", ctx, email).Return(tc.emailTaken, nil)
			if tc.wantCreate {
				userRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(createdUser, nil)
			}
			if tc.wantStore {
				otpRepo.On("Store", ctx, createdUser.ID, mock.AnythingOfType("string")).Return(nil)
			}
			if tc.wantMail {
				mailer.On("SendOTP", ctx, email, mock.AnythingOfType("string")).Return(nil)
			}

			svc := newAuthService(userRepo, sessRepo, otpRepo, tokenSvc, mailer)
			err := svc.RegisterByEmail(ctx, email)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
			otpRepo.AssertExpectations(t)
			mailer.AssertExpectations(t)
		})
	}
}

// ---------------------------------------------------------------------------
// Internal helper — replicates hmacSHA256 for test hash generation
// ---------------------------------------------------------------------------

func hmacForTest(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
