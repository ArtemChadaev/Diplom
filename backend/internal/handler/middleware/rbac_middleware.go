package middleware

import (
	"net/http"

	"github.com/ima/diplom-backend/internal/domain"
)

// RequireRole checks if the caller has any of the allowed roles.
// Must be chained exactly after AuthRequired middleware.
func RequireRole(allowedRoles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleVal := r.Context().Value(CtxRole)
			if roleVal == nil {
				http.Error(w, "unauthorized context", http.StatusUnauthorized)
				return
			}

			userRole, ok := roleVal.(domain.UserRole)
			if !ok {
				http.Error(w, "invalid role in context", http.StatusInternalServerError)
				return
			}

			hasAccess := false
			for _, allowed := range allowedRoles {
				if userRole == allowed {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				http.Error(w, "insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
