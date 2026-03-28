package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

// NOTE: register and login endpoints removed — authentication is OAuth-only.
// Use POST /auth/google for Google OAuth login.

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token missing", http.StatusUnauthorized)
		return
	}

	meta := domain.SessionMeta{
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
	}

	pair, err := h.service.Auth.RefreshTokens(r.Context(), cookie.Value, meta)
	if err != nil {
		// token expired, revoked, or invalid
		http.SetCookie(w, &http.Cookie{
			Name:   "refresh_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1, // delete cookie
		})
		http.Error(w, "invalid session, please login again", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		Path:     "/auth",
		MaxAge:   15 * 24 * 3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	resp := dto.TokenResponse{
		AccessToken: pair.AccessToken,
		ExpiresIn:   pair.ExpiresIn,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		h.clearRefreshCookie(w)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Hash token and find session
	hash := h.tokenSvc.HashToken(cookie.Value)
	_ = hash // session revocation happens via token hash lookup in service

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		const prefix = "Bearer "
		if len(authHeader) > len(prefix) {
			token := authHeader[len(prefix):]
			claims, parseErr := h.tokenSvc.ParseAccessToken(token)
			if parseErr == nil {
				revokeErr := h.service.Auth.RevokeSession(r.Context(), claims.SessionID, claims.UserID, claims.Role)
				if revokeErr != nil {
					logger.FromContext(r.Context()).Warn("failed to revoke session during logout",
						"sessionID", claims.SessionID.String(),
						"error", revokeErr.Error(),
					)
				}
			}
		}
	}

	h.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})
}
