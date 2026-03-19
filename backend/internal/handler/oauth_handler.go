package handler

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
)

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
		if errors.Is(err, domain.ErrUserUnverified) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "account pending admin approval"})
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
		Secure:   false, // Set true in production
		SameSite: http.SameSiteStrictMode,
	})

	resp := dto.TokenResponse{
		AccessToken: pair.AccessToken,
		ExpiresIn:   pair.ExpiresIn,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
