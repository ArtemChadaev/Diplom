package handler

import (
	"encoding/json"
	"errors"
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
		Secure:   h.cfg.Env == "production",
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

func (h *Handler) sendCode(w http.ResponseWriter, r *http.Request) {
	var req dto.SendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}

	if err := h.service.Auth.SendOTPCode(r.Context(), req.Email); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user with this email not found")
			return
		}
		if errors.Is(err, domain.ErrOTPMaxAttempts) {
			writeError(w, http.StatusTooManyRequests, "too many requests, try again later")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to send code")
		return
	}

	resp := dto.SendCodeResponse{
		Message:   "OTP code sent successfully",
		ExpiresIn: 600, // 10 minutes
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) verifyCode(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Code == "" {
		writeError(w, http.StatusBadRequest, "email and code are required")
		return
	}

	meta := domain.SessionMeta{
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
	}

	pair, err := h.service.Auth.VerifyOTPCode(r.Context(), req.Email, req.Code, meta)
	if err != nil {
		if errors.Is(err, domain.ErrOTPNotFound) || errors.Is(err, domain.ErrOTPInvalid) || errors.Is(err, domain.ErrOTPMaxAttempts) { // using ErrOTPNotFound as "invalid or expired"
			writeError(w, http.StatusUnauthorized, "invalid or expired code")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to verify code")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		Path:     "/auth",
		MaxAge:   15 * 24 * 3600,
		HttpOnly: true,
		Secure:   h.cfg.Env == "production",
		SameSite: http.SameSiteStrictMode,
	})

	resp := dto.TokenResponse{
		AccessToken: pair.AccessToken,
		ExpiresIn:   pair.ExpiresIn,
	}

	writeJSON(w, http.StatusOK, resp)
}
