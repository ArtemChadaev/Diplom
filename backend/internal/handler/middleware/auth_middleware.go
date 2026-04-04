package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

type contextKey string

const (
	CtxUserID contextKey = "userID"
	CtxRole   contextKey = "role"
	CtxEmail  contextKey = "email"
)

// AuthRequired validates the JWT from the Authorization header
func AuthRequired(tokenSvc domain.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			claims, err := tokenSvc.ParseAccessToken(tokenStr)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxRole, claims.Role)
			ctx = context.WithValue(ctx, CtxEmail, claims.Email)

			// Enrich logger context with UserID
			ctx = logger.WithUserID(ctx, claims.UserID)
			
			if logData, ok := r.Context().Value(CtxLogData).(*LogData); ok {
				logData.UserID = claims.UserID
			}

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
