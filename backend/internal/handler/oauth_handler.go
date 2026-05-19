package handler

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
)

// googleLogin godoc
// @Summary      Google OAuth Login
// @Description  Authenticates a user using Google OAuth ID Token. Returns an access token and sets a HTTP-only refresh token cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.GoogleAuthRequest  true  "Google ID token"
// @Success      200   {object}  dto.TokenResponse
// @Failure      400   {object}  dto.ErrorResponse  "invalid JSON"
// @Failure      401   {object}  dto.ErrorResponse  "google auth failed"
// @Failure      403   {object}  dto.ErrorResponse  "account is blocked"
// @Router       /auth/google [post]
func (h *Handler) googleLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.GoogleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	userAgent := r.UserAgent()
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	pair, err := h.service.Auth.LoginWithGoogle(r.Context(), req.IDToken, userAgent, ip)
	if err != nil {
		if errors.Is(err, domain.ErrUserBlocked) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "account is blocked"})
			return
		}
		http.Error(w, "google auth failed", http.StatusUnauthorized)
		return
	}

	// Set HttpOnly refresh token cookie
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
