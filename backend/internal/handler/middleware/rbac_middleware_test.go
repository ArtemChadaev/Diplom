package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/middleware"
	"github.com/stretchr/testify/assert"
)

// passThrough is an http.Handler that records if it was called.
type passThrough struct {
	called bool
}

func (h *passThrough) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
	w.WriteHeader(http.StatusOK)
}

// ctxWithRole returns a request whose context contains the given role,
// simulating what AuthRequired middleware injects.
func ctxWithRole(r *http.Request, role domain.UserRole) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.CtxRole, role)
	return r.WithContext(ctx)
}

// ---------------------------------------------------------------------------
// RequireRole
// ---------------------------------------------------------------------------

// Scenario: RBAC enforcement (TC-03 from testing.md)
//
//	Covers: allowed role passes, disallowed role → 403, missing role in ctx → 401,
//	        multi-role allow-list, all non-admin roles blocked from admin endpoint
func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		ctxRole        *domain.UserRole // nil = no role in context at all
		allowedRoles   []domain.UserRole
		wantStatus     int
		wantNextCalled bool
	}{
		{
			name:           "admin role — allowed",
			ctxRole:        rolePtr(domain.RoleAdmin),
			allowedRoles:   []domain.UserRole{domain.RoleAdmin},
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
		{
			name:           "storekeeper trying admin endpoint → 403",
			ctxRole:        rolePtr(domain.RoleStorekeeper),
			allowedRoles:   []domain.UserRole{domain.RoleAdmin},
			wantStatus:     http.StatusForbidden,
			wantNextCalled: false,
		},
		{
			name:           "pharmacist trying admin endpoint → 403",
			ctxRole:        rolePtr(domain.RolePharmacist),
			allowedRoles:   []domain.UserRole{domain.RoleAdmin},
			wantStatus:     http.StatusForbidden,
			wantNextCalled: false,
		},
		{
			name:           "no role in context (AuthRequired not applied) → 401",
			ctxRole:        nil,
			allowedRoles:   []domain.UserRole{domain.RoleAdmin},
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "multi-role allow list: QP is in [admin, qp] → passes",
			ctxRole:        rolePtr(domain.RoleQP),
			allowedRoles:   []domain.UserRole{domain.RoleAdmin, domain.RoleQP},
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
		{
			name:           "multi-role allow list: pharmacist not in [admin, qp] → 403",
			ctxRole:        rolePtr(domain.RolePharmacist),
			allowedRoles:   []domain.UserRole{domain.RoleAdmin, domain.RoleQP},
			wantStatus:     http.StatusForbidden,
			wantNextCalled: false,
		},
		{
			name:           "warehouse_manager allowed when in list",
			ctxRole:        rolePtr(domain.RoleWarehouseManager),
			allowedRoles:   []domain.UserRole{domain.RoleAdmin, domain.RoleWarehouseManager},
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			next := &passThrough{}
			handler := middleware.RequireRole(tc.allowedRoles...)(next)

			req := httptest.NewRequest(http.MethodGet, "/endpoint", nil)
			if tc.ctxRole != nil {
				req = ctxWithRole(req, *tc.ctxRole)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantStatus, rr.Code)
			assert.Equal(t, tc.wantNextCalled, next.called)
		})
	}
}

// rolePtr is a test helper that returns a pointer to a UserRole value.
func rolePtr(r domain.UserRole) *domain.UserRole {
	return &r
}
