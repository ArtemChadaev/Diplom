package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/middleware"
	"github.com/ima/diplom-backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestTokenSvc returns a real tokenService using a fixed test secret.
func newTestTokenSvc() domain.TokenService {
	return service.NewTokenService("middleware-test-secret", 15*time.Minute, 15*24*time.Hour)
}

// nextHandler records whether ServeHTTP was called.
type nextHandler struct {
	called bool
}

func (h *nextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
	w.WriteHeader(http.StatusOK)
}

// ---------------------------------------------------------------------------
// AuthRequired
// ---------------------------------------------------------------------------

// Scenario: JWT Bearer token gate
//
//	Covers: valid token passes through; missing header → 401; wrong prefix → 401;
//	        expired token → 401; garbage token → 401
func TestAuthRequired(t *testing.T) {
	svc := newTestTokenSvc()
	user := &domain.User{ID: 42, Email: "test@example.com", Role: domain.RoleAdmin}

	validToken, err := svc.GenerateAccessToken(user, uuid.New())
	require.NoError(t, err)

	// Generate an already-expired token
	expiredSvc := service.NewTokenService("middleware-test-secret", -1*time.Second, 15*24*time.Hour)
	expiredToken, err := expiredSvc.GenerateAccessToken(user, uuid.New())
	require.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		wantStatus     int
		wantNextCalled bool
	}{
		{
			name:           "valid Bearer token — next handler called with 200",
			authHeader:     "Bearer " + validToken,
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
		{
			name:           "missing Authorization header → 401",
			authHeader:     "",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "wrong scheme (Token prefix) → 401",
			authHeader:     "Token " + validToken,
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "Bearer with no token value → 401",
			authHeader:     "Bearer",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "expired JWT → 401",
			authHeader:     "Bearer " + expiredToken,
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "completely invalid JWT → 401",
			authHeader:     "Bearer not.a.valid.jwt",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			next := &nextHandler{}
			handler := middleware.AuthRequired(svc)(next)

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantStatus, rr.Code)
			assert.Equal(t, tc.wantNextCalled, next.called)
		})
	}
}

// ---------------------------------------------------------------------------
// AuthRequired — context enrichment
// ---------------------------------------------------------------------------

// Scenario: Valid token enriches context with UserID, Role, Email
//
//	Given:  A valid JWT for user ID=42, role=admin, email=test@example.com
//	When:   AuthRequired processes the request
//	Then:   Context values for UserID, Role, Email match the claims
func TestAuthRequired_ContextEnrichment(t *testing.T) {
	svc := newTestTokenSvc()
	user := &domain.User{ID: 42, Email: "ctx@example.com", Role: domain.RoleAdmin}

	token, err := svc.GenerateAccessToken(user, uuid.New())
	require.NoError(t, err)

	var capturedCtx struct {
		userID interface{}
		role   interface{}
		email  interface{}
	}

	captureHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx.userID = r.Context().Value(middleware.CtxUserID)
		capturedCtx.role = r.Context().Value(middleware.CtxRole)
		capturedCtx.email = r.Context().Value(middleware.CtxEmail)
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.AuthRequired(svc)(captureHandler)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 42, capturedCtx.userID)
	assert.Equal(t, domain.RoleAdmin, capturedCtx.role)
	assert.Equal(t, "ctx@example.com", capturedCtx.email)
}
